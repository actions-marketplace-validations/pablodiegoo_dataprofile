package github

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/pablodiegoo/dataprofile/internal/domain"
	"github.com/pablodiegoo/dataprofile/internal/render"
)

func TestCollectRepositoriesLanguagesAndHygiene(t *testing.T) {
	var activeLangReqs int32
	var maxConcurrentLangReqs int32

	mux := http.NewServeMux()
	mux.HandleFunc("/users/adalovelace/repos", func(w http.ResponseWriter, r *http.Request) {
		repos := []restRepo{
			{
				Name:       "ml-pipeline",
				FullName:   "adalovelace/ml-pipeline",
				Fork:       false,
				Stargazers: 120,
				Forks:      18,
				Language:   "Python",
			},
			{
				Name:       "eda-experiments",
				FullName:   "adalovelace/eda-experiments",
				Fork:       false,
				Stargazers: 45,
				Forks:      7,
				Language:   "Jupyter Notebook",
			},
			{
				Name:       "quarto-reports",
				FullName:   "adalovelace/quarto-reports",
				Fork:       false,
				Stargazers: 30,
				Forks:      4,
				Language:   "Quarto",
			},
			{
				Name:       "forked-repo",
				FullName:   "adalovelace/forked-repo",
				Fork:       true,
				Stargazers: 999,
				Forks:      100,
				Language:   "JavaScript",
			},
			{
				Name:       "legacy-archive",
				FullName:   "adalovelace/legacy-archive",
				Fork:       false,
				Stargazers: 500,
				Forks:      50,
				Language:   "Java",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(repos)
	})

	langMaps := map[string]map[string]int64{
		"/repos/adalovelace/ml-pipeline/languages": {
			"Python": 2 * 1024 * 1024,
			"SQL":    512 * 1024,
			"HTML":   256 * 1024,
		},
		"/repos/adalovelace/eda-experiments/languages": {
			"Jupyter Notebook": 10 * 1024 * 1024,
			"Python":           512 * 1024,
		},
		"/repos/adalovelace/quarto-reports/languages": {
			"Quarto": 1 * 1024 * 1024,
			"R":      256 * 1024,
		},
	}

	for path, lmap := range langMaps {
		lmapCopy := lmap
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			cur := atomic.AddInt32(&activeLangReqs, 1)
			for {
				prev := atomic.LoadInt32(&maxConcurrentLangReqs)
				if cur <= prev || atomic.CompareAndSwapInt32(&maxConcurrentLangReqs, prev, cur) {
					break
				}
			}
			defer atomic.AddInt32(&activeLangReqs, -1)

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(lmapCopy)
		})
	}

	rootContents := map[string][]map[string]string{
		"/repos/adalovelace/ml-pipeline/contents": {
			{"name": "poetry.lock", "type": "file"},
			{"name": ".github", "type": "dir"},
			{"name": "tests", "type": "dir"},
			{"name": "Dockerfile", "type": "file"},
		},
		"/repos/adalovelace/eda-experiments/contents": {
			{"name": "requirements.txt", "type": "file"},
			{"name": ".github", "type": "dir"},
		},
		"/repos/adalovelace/quarto-reports/contents": {
			{"name": "go.sum", "type": "file"},
			{"name": "pipeline_test.go", "type": "file"},
		},
	}
	for path, entries := range rootContents {
		entriesCopy := entries
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(entriesCopy)
		})
	}

	srv := httptest.NewServer(mux)
	defer srv.Close()

	client := NewClient("")
	client.BaseURL = srv.URL

	ghCfg := domain.GitHubSourceConfig{
		Username:     "adalovelace",
		IncludeForks: false,
		ExcludeRepos: []string{"legacy-archive"},
		ExcludeLangs: []string{"HTML", "CSS"},
	}
	normCfg := domain.NormalizationConfig{
		DeflateNotebooks: true,
		DetectDataStack:  true,
	}

	langPanel, _, impact, err := client.Collect(context.Background(), ghCfg, normCfg)
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	// 3 active repos (forked-repo and legacy-archive excluded)
	if langPanel.TotalRepos != 3 || impact.PublicRepos != 3 {
		t.Errorf("expected 3 filtered repos, got TotalRepos=%d PublicRepos=%d", langPanel.TotalRepos, impact.PublicRepos)
	}
	if impact.TotalStars != 195 || impact.TotalForks != 29 {
		t.Errorf("expected TotalStars=195 TotalForks=29, got stars=%d forks=%d", impact.TotalStars, impact.TotalForks)
	}
	if atomic.LoadInt32(&maxConcurrentLangReqs) > 8 {
		t.Errorf("semaphore bound violated: maxConcurrentLangReqs=%d", maxConcurrentLangReqs)
	}

	// 1 pure production repo (ml-pipeline) and 2 notebook/Quarto repos (eda-experiments, quarto-reports)
	if langPanel.ProductionFiles != 1 || langPanel.ExplorationFiles != 2 {
		t.Errorf("expected ProductionFiles=1 ExplorationFiles=2, got prod=%d eda=%d", langPanel.ProductionFiles, langPanel.ExplorationFiles)
	}

	// Stripped base64 MB from 10 MB Jupyter Notebook + 1 MB Quarto should be > 8 MB
	if langPanel.StrippedBase64MB < 8.5 {
		t.Errorf("expected StrippedBase64MB >= 8.5, got %.1f", langPanel.StrippedBase64MB)
	}
	if langPanel.ProductionPct <= 0 || langPanel.ExplorationPct <= 0 {
		t.Errorf("expected positive ProductionPct and ExplorationPct, got %.1f%% / %.1f%%", langPanel.ProductionPct, langPanel.ExplorationPct)
	}

	// Excluded language (HTML) and raw Jupyter Notebook/Quarto must not appear in normalized languages
	for _, lg := range langPanel.Languages {
		if lg.Name == "HTML" || lg.Name == "Jupyter Notebook" || lg.Name == "Quarto" {
			t.Errorf("unexpected language %q in normalized output", lg.Name)
		}
	}
	if len(langPanel.Languages) == 0 || langPanel.Languages[0].Name != "Python" {
		t.Errorf("expected Python as top language, got %+v", langPanel.Languages)
	}

	// Verified engineering hygiene signals from /repos/{owner}/{repo}/contents:
	// Lockfiles: 3 (poetry.lock, requirements.txt, go.sum)
	// CI/CD: 2 (.github in ml-pipeline and eda-experiments)
	// Tests: 2 (tests dir in ml-pipeline, pipeline_test.go in quarto-reports)
	// Docker: 1 (Dockerfile in ml-pipeline)
	if langPanel.LockfileRepos != 3 || langPanel.CICDRepos != 2 || langPanel.TestSuiteRepos != 2 || langPanel.ContainerizedRepos != 1 {
		t.Errorf("expected verified hygiene counts Lockfiles=3 CI=2 Tests=2 Docker=1, got Lockfiles=%d CI=%d Tests=%d Docker=%d",
			langPanel.LockfileRepos, langPanel.CICDRepos, langPanel.TestSuiteRepos, langPanel.ContainerizedRepos)
	}
}

func TestCollectScansDependencyManifestsAndSQLArtifactsForStack(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/adalovelace/repos", func(w http.ResponseWriter, r *http.Request) {
		repos := []restRepo{
			{
				Name:       "lakehouse-etl",
				FullName:   "adalovelace/lakehouse-etl",
				Language:   "Python",
				Stargazers: 50,
			},
			{
				Name:       "vision-trainer",
				FullName:   "adalovelace/vision-trainer",
				Language:   "Python",
				Stargazers: 80,
			},
			{
				Name:       "biostat-models",
				FullName:   "adalovelace/biostat-models",
				Language:   "R",
				Stargazers: 35,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(repos)
	})

	// Repo 1 has Python + SQL in its language map and a pyproject.toml + requirements.txt
	mux.HandleFunc("/repos/adalovelace/lakehouse-etl/languages", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]int64{"Python": 120000, "SQL": 45000})
	})
	mux.HandleFunc("/repos/adalovelace/lakehouse-etl/contents/pyproject.toml", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[project]
dependencies = ["polars>=0.20", "dbt-core>=1.7", "scikit-learn>=1.4", "plotly>=5.18", "streamlit>=1.31"]`))
	})

	// Repo 2 has environment.yml and requirements.txt
	mux.HandleFunc("/repos/adalovelace/vision-trainer/languages", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]int64{"Python": 200000})
	})
	mux.HandleFunc("/repos/adalovelace/vision-trainer/contents/environment.yml", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`name: vision
dependencies:
  - pandas
  - pytorch
  - transformers
  - xgboost
  - matplotlib
  - plotly`))
	})
	mux.HandleFunc("/repos/adalovelace/vision-trainer/contents/requirements.txt", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("scikit-learn==1.4.0\nduckdb==0.10.0\n"))
	})

	// Repo 3 has R DESCRIPTION file
	mux.HandleFunc("/repos/adalovelace/biostat-models/languages", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]int64{"R": 95000})
	})
	mux.HandleFunc("/repos/adalovelace/biostat-models/contents/DESCRIPTION", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`Package: biostat
Imports:
    ggplot2,
    tidyverse,
    dplyr`))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	client := NewClient("")
	client.BaseURL = srv.URL

	ghCfg := domain.GitHubSourceConfig{Username: "adalovelace"}
	normCfg := domain.NormalizationConfig{DeflateNotebooks: true, DetectDataStack: true}

	_, stack, _, err := client.Collect(context.Background(), ghCfg, normCfg)
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if stack.TotalSources != 3 {
		t.Errorf("expected TotalSources=3, got %d", stack.TotalSources)
	}

	// Data & Analytics: Pandas/Polars (2 repos), SQL/DuckDB (2 repos: SQL lang in repo 1 + duckdb in repo 2), dbt/Airflow (1 repo)
	if len(stack.DataAnalytics) != 4 {
		t.Fatalf("expected 4 items in DataAnalytics, got %d", len(stack.DataAnalytics))
	}
	if stack.DataAnalytics[0].Name != "Pandas / Polars" || stack.DataAnalytics[0].RepoCount != 2 || !stack.DataAnalytics[0].IsLeader {
		t.Errorf("unexpected DataAnalytics[0]: %+v", stack.DataAnalytics[0])
	}
	if stack.DataAnalytics[1].Name != "SQL / DuckDB" || stack.DataAnalytics[1].RepoCount != 2 || stack.DataAnalytics[1].IsLeader {
		t.Errorf("unexpected DataAnalytics[1]: %+v", stack.DataAnalytics[1])
	}

	// ML & AI: Scikit-Learn (2 repos -> leader), PyTorch (1), XGBoost/LightGBM (1), Hugging Face (1)
	if len(stack.MLAndAI) != 4 {
		t.Fatalf("expected 4 items in MLAndAI, got %d", len(stack.MLAndAI))
	}
	if stack.MLAndAI[0].Name != "Scikit-Learn" || stack.MLAndAI[0].RepoCount != 2 || !stack.MLAndAI[0].IsLeader {
		t.Errorf("unexpected MLAndAI[0]: %+v", stack.MLAndAI[0])
	}

	// Viz & Data Apps: Plotly/Dash (2 repos -> leader), Streamlit (1), Matplotlib/Seaborn (1), ggplot2/Tidyverse (1)
	if len(stack.VizAndApps) != 4 {
		t.Fatalf("expected 4 items in VizAndApps, got %d", len(stack.VizAndApps))
	}
	if stack.VizAndApps[0].Name != "Plotly / Dash" || stack.VizAndApps[0].RepoCount != 2 || !stack.VizAndApps[0].IsLeader {
		t.Errorf("unexpected VizAndApps[0]: %+v", stack.VizAndApps[0])
	}
}

func TestCollectEnrichesGraphQLContributionsAnd52WeekTrend(t *testing.T) {
	var gotAuthHeader string
	var gotLoginVar string

	mux := http.NewServeMux()
	mux.HandleFunc("/users/adalovelace/repos", func(w http.ResponseWriter, r *http.Request) {
		repos := []restRepo{
			{Name: "core-ml", FullName: "adalovelace/core-ml", Stargazers: 210, Forks: 34, Language: "Python"},
			{Name: "data-lake", FullName: "adalovelace/data-lake", Stargazers: 90, Forks: 12, Language: "SQL"},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(repos)
	})
	mux.HandleFunc("/repos/adalovelace/core-ml/languages", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]int64{"Python": 120000})
	})
	mux.HandleFunc("/repos/adalovelace/data-lake/languages", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]int64{"SQL": 80000})
	})
	mux.HandleFunc("/repos/adalovelace/core-ml/traffic/views", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]int{"count": 640, "uniques": 190})
	})
	mux.HandleFunc("/repos/adalovelace/data-lake/traffic/views", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]int{"count": 280, "uniques": 85})
	})
	mux.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST /graphql, got %s", r.Method)
		}
		gotAuthHeader = r.Header.Get("Authorization")
		var reqBody struct {
			Query     string            `json:"query"`
			Variables map[string]string `json:"variables"`
		}
		_ = json.NewDecoder(r.Body).Decode(&reqBody)
		gotLoginVar = reqBody.Variables["login"]

		type day struct {
			ContributionCount int `json:"contributionCount"`
		}
		type week struct {
			ContributionDays []day `json:"contributionDays"`
		}
		weeks := make([]week, 52)
		total := 0
		for i := 0; i < 52; i++ {
			c1 := (i % 7) + 2
			c2 := (i % 5) + 1
			weeks[i] = week{
				ContributionDays: []day{{ContributionCount: c1}, {ContributionCount: c2}},
			}
			total += c1 + c2
		}

		resp := map[string]interface{}{
			"data": map[string]interface{}{
				"user": map[string]interface{}{
					"contributionsCollection": map[string]interface{}{
						"contributionCalendar": map[string]interface{}{
							"totalContributions": total,
							"weeks":              weeks,
						},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	client := NewClient("ghp_secret_token")
	client.BaseURL = srv.URL

	_, _, impact, err := client.Collect(
		context.Background(),
		domain.GitHubSourceConfig{Username: "adalovelace"},
		domain.NormalizationConfig{DeflateNotebooks: true},
	)
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if gotAuthHeader != "Bearer ghp_secret_token" {
		t.Errorf("expected Authorization header 'Bearer ghp_secret_token', got %q", gotAuthHeader)
	}
	if gotLoginVar != "adalovelace" {
		t.Errorf("expected GraphQL login variable 'adalovelace', got %q", gotLoginVar)
	}
	if impact.PublicRepos != 2 || impact.TotalStars != 300 || impact.TotalForks != 46 {
		t.Errorf("expected PublicRepos=2 TotalStars=300 TotalForks=46, got %+v", impact)
	}
	if len(impact.WeeklyTrend) != 52 {
		t.Fatalf("expected 52 weekly contribution points, got %d", len(impact.WeeklyTrend))
	}
	if impact.WeeklyTrend[0] != 3 || impact.WeeklyTrend[1] != 5 {
		t.Errorf("unexpected weekly sums: week[0]=%d week[1]=%d", impact.WeeklyTrend[0], impact.WeeklyTrend[1])
	}
	expectedTotal := 0
	for _, w := range impact.WeeklyTrend {
		expectedTotal += w
	}
	if impact.AnnualContributions != expectedTotal {
		t.Errorf("expected AnnualContributions=%d, got %d", expectedTotal, impact.AnnualContributions)
	}
	if impact.RepoViews14d != 920 {
		t.Errorf("expected RepoViews14d=920 from /traffic/views, got %d", impact.RepoViews14d)
	}
	if impact.ActiveStreakDays <= 0 || impact.YoYGrowthPct <= 0 || impact.LinesAdded <= 0 || impact.LinesDeleted <= 0 {
		t.Errorf("expected velocity metrics populated (ActiveStreakDays>0, YoYGrowthPct>0, LinesAdded>0, LinesDeleted>0), got %+v", impact)
	}

	// Verify non-200 GraphQL status is rejected even if body contains JSON
	errSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"data":{"user":{"contributionsCollection":{"contributionCalendar":{"totalContributions":999}}}}}`))
	}))
	defer errSrv.Close()

	errClient := NewClient("bad_token")
	errClient.BaseURL = errSrv.URL
	if _, _, err := errClient.fetchGraphQLContributions(context.Background(), "adalovelace"); err == nil {
		t.Errorf("expected fetchGraphQLContributions to fail on HTTP 401 status")
	}
}

func TestCollectIncludePrivateCountAndTrafficFallback(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/adalovelace/repos", func(w http.ResponseWriter, r *http.Request) {
		repos := []restRepo{
			{Name: "repo-a", FullName: "adalovelace/repo-a", Stargazers: 40, Forks: 10, Language: "Go"},
			{Name: "repo-b", FullName: "adalovelace/repo-b", Stargazers: 20, Forks: 5, Language: "Python"},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(repos)
	})
	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer ghp_priv_token" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]int{
			"total_private_repos": 5,
			"owned_private_repos": 4,
		})
	})
	mux.HandleFunc("/repos/adalovelace/repo-a/languages", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]int64{"Go": 80000})
	})
	mux.HandleFunc("/repos/adalovelace/repo-b/languages", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]int64{"Python": 60000})
	})
	// Simulate 403 Forbidden on traffic/views (e.g. token without push scope) -> should estimate from stars/forks
	mux.HandleFunc("/repos/adalovelace/repo-a/traffic/views", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})
	mux.HandleFunc("/repos/adalovelace/repo-b/traffic/views", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	client := NewClient("ghp_priv_token")
	client.BaseURL = srv.URL

	// 1. With IncludePrivateCount: true -> 2 public + 5 private = 7 repos
	lpPriv, _, imPriv, err := client.Collect(
		context.Background(),
		domain.GitHubSourceConfig{Username: "adalovelace", IncludePrivateCount: true},
		domain.NormalizationConfig{DeflateNotebooks: true},
	)
	if err != nil {
		t.Fatalf("Collect with IncludePrivateCount=true failed: %v", err)
	}
	if lpPriv.TotalRepos != 7 || imPriv.PublicRepos != 7 {
		t.Errorf("expected TotalRepos=7 and PublicRepos=7 when IncludePrivateCount=true, got TotalRepos=%d PublicRepos=%d",
			lpPriv.TotalRepos, imPriv.PublicRepos)
	}
	if imPriv.RepoViews14d <= 0 {
		t.Errorf("expected positive estimated RepoViews14d when /traffic/views returns 403, got %d", imPriv.RepoViews14d)
	}
	// Since /contents returned 404, baseline ratios should be used on 7 repos
	if lpPriv.LockfileRepos <= 0 || lpPriv.CICDRepos <= 0 {
		t.Errorf("expected fallback baseline hygiene counts when /contents is 404, got %+v", lpPriv)
	}

	// 2. With IncludePrivateCount: false -> 2 public repos only
	lpPub, _, imPub, err := client.Collect(
		context.Background(),
		domain.GitHubSourceConfig{Username: "adalovelace", IncludePrivateCount: false},
		domain.NormalizationConfig{DeflateNotebooks: true},
	)
	if err != nil {
		t.Fatalf("Collect with IncludePrivateCount=false failed: %v", err)
	}
	if lpPub.TotalRepos != 2 || imPub.PublicRepos != 2 {
		t.Errorf("expected TotalRepos=2 and PublicRepos=2 when IncludePrivateCount=false, got TotalRepos=%d PublicRepos=%d",
			lpPub.TotalRepos, imPub.PublicRepos)
	}
}

func TestCollectThreeTierPrivateOwnedMonorepoAndContributedRepos(t *testing.T) {
	mux := http.NewServeMux()
	var gotAffiliation string

	// Authenticated /user/repos returns both public and private Owned Repositories + a passive fork
	mux.HandleFunc("/user/repos", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer pat_token" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		gotAffiliation = r.URL.Query().Get("affiliation")
		repos := []restRepo{
			{
				Name:          "AGP",
				FullName:      "pablodiegoo/AGP",
				Private:       true,
				Fork:          false,
				Stargazers:    5,
				Forks:         1,
				Language:      "Jupyter Notebook",
				DefaultBranch: "main",
			},
			{
				Name:          "secret-client-audit",
				FullName:      "pablodiegoo/secret-client-audit",
				Private:       true,
				Fork:          false,
				Stargazers:    0,
				Forks:         0,
				Language:      "", // Empty primary language and empty /languages map; files counted strictly from recursive git tree
				DefaultBranch: "main",
			},
			{
				Name:          "dataprofile",
				FullName:      "pablodiegoo/dataprofile",
				Private:       false,
				Fork:          false,
				Stargazers:    12,
				Forks:         2,
				Language:      "Go",
				DefaultBranch: "main",
			},
			{
				Name:          "iota-sdk",
				FullName:      "pablodiegoo/iota-sdk",
				Private:       false,
				Fork:          true,
				Stargazers:    900,
				Forks:         200,
				Language:      "Go",
				DefaultBranch: "main",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(repos)
	})

	mux.HandleFunc("/repos/pablodiegoo/AGP/languages", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]int64{
			"Jupyter Notebook": 20 * 1024 * 1024,
			"Python":           4 * 1024 * 1024,
		})
	})
	mux.HandleFunc("/repos/pablodiegoo/secret-client-audit/languages", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]int64{})
	})
	mux.HandleFunc("/repos/pablodiegoo/dataprofile/languages", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]int64{
			"Go": 300 * 1024,
		})
	})

	// Monorepo recursive git tree with multiple Analytical Workspaces + nested hygiene signals
	mux.HandleFunc("/repos/pablodiegoo/AGP/git/trees/main", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"tree": []gitTreeEntry{
				{Path: "Reports/Rio Bonito/DataAnalysis.ipynb", Type: "blob", Size: 42000},
				{Path: "Reports/Rio Bonito/Wrangling.ipynb", Type: "blob", Size: 38000},
				{Path: "Reports/Rio Bonito/regressao.py", Type: "blob", Size: 4200},
				{Path: "Reports/Empetur/generate_maps.py", Type: "blob", Size: 5100},
				{Path: "Reports/Empetur/test_models.py", Type: "blob", Size: 3400},
				{Path: "Reports/Marica/duckdb_cleaner.py", Type: "blob", Size: 6100},
				{Path: "Utils/pyproject.toml", Type: "blob", Size: 1200},
				{Path: "Utils/uv.lock", Type: "blob", Size: 18000},
				{Path: "Utils/agp_media/cli.py", Type: "blob", Size: 3200},
				{Path: "Utils/ZapSender/adLog.ipynb", Type: "blob", Size: 25000},
				{Path: ".github/workflows/ci.yml", Type: "blob", Size: 900},
				{Path: "deploy/Dockerfile", Type: "blob", Size: 650},
			},
		})
	})
	// Second private repo with 2 .py files and 1 .qmd file in its git tree (even though /languages was empty)
	mux.HandleFunc("/repos/pablodiegoo/secret-client-audit/git/trees/main", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"tree": []gitTreeEntry{
				{Path: "audit/pipeline.py", Type: "blob", Size: 2400},
				{Path: "audit/test_pipeline.py", Type: "blob", Size: 1800},
				{Path: "audit/report.qmd", Type: "blob", Size: 8400},
			},
		})
	})
	mux.HandleFunc("/repos/pablodiegoo/AGP/contents/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/pablodiegoo/AGP/contents/Utils/pyproject.toml":
			_, _ = w.Write([]byte(`dependencies = ["pandas", "geopandas", "torch", "transformers", "openai-whisper", "typer"]`))
		case "/repos/pablodiegoo/AGP/contents/Reports/Rio Bonito/regressao.py":
			_, _ = w.Write([]byte("SECRET_PROPRIETARY_FORMULA = 42\nimport pandas as pd\nimport numpy as np\nimport statsmodels.api as sm\nfrom sklearn.linear_model import LogisticRegression\nimport seaborn as sns\n"))
		case "/repos/pablodiegoo/AGP/contents/Reports/Empetur/generate_maps.py":
			_, _ = w.Write([]byte("import geopandas as gpd\nimport osmnx as ox\nimport matplotlib.pyplot as plt\n"))
		case "/repos/pablodiegoo/AGP/contents/Reports/Marica/duckdb_cleaner.py":
			_, _ = w.Write([]byte("import duckdb\nimport pandas as pd\nfrom sklearn.cluster import DBSCAN\n"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	// GraphQL returns multi-year contributions + 2 Contributed Repositories
	mux.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		var reqBody struct {
			Query string `json:"query"`
		}
		_ = json.NewDecoder(r.Body).Decode(&reqBody)
		w.Header().Set("Content-Type", "application/json")
		if len(reqBody.Query) > 0 && reqBody.Query[0] == 'q' && !strings.Contains(reqBody.Query, "y2026:") {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"user": map[string]interface{}{
						"contributionsCollection": map[string]interface{}{
							"contributionYears": []int{2026, 2025},
							"contributionCalendar": map[string]interface{}{
								"totalContributions": 8331,
								"weeks": []map[string]interface{}{
									{"contributionDays": []map[string]int{{"contributionCount": 20}}},
									{"contributionDays": []map[string]int{{"contributionCount": 40}}},
								},
							},
						},
						"repositoriesContributedTo": map[string]interface{}{
							"totalCount": 2,
							"nodes": []map[string]string{
								{"nameWithOwner": "iota-uz/iota-sdk"},
								{"nameWithOwner": "Agora-Pesquisa/agora-pesquisa.github.io"},
							},
						},
					},
				},
			})
			return
		}
		// Multi-year query response
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"user": map[string]interface{}{
					"y2026": map[string]interface{}{"contributionCalendar": map[string]int{"totalContributions": 8144}},
					"y2025": map[string]interface{}{"contributionCalendar": map[string]int{"totalContributions": 933}},
				},
			},
		})
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	client := NewClient("pat_token")
	client.BaseURL = srv.URL

	lp, sp, im, err := client.Collect(
		context.Background(),
		domain.GitHubSourceConfig{Username: "pablodiegoo", IncludeForks: false, IncludePrivateCount: true},
		domain.NormalizationConfig{DeflateNotebooks: true, DetectDataStack: true},
	)
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if gotAffiliation != "owner" {
		t.Errorf("expected /user/repos query param affiliation=owner, got %q", gotAffiliation)
	}
	// 3 Owned Repositories (AGP private + secret-client-audit private + dataprofile public; iota-sdk fork excluded)
	if lp.TotalRepos != 3 {
		t.Errorf("expected 3 Owned Repositories in lp.TotalRepos, got %d", lp.TotalRepos)
	}
	// Impact reach includes 3 Owned + 2 Contributed Repositories = 5
	if im.ContributedRepos != 2 || im.PublicRepos != 5 {
		t.Errorf("expected ContributedRepos=2 and PublicRepos=5, got ContributedRepos=%d PublicRepos=%d", im.ContributedRepos, im.PublicRepos)
	}
	// Multi-year all-time contributions = 8144 + 933 = 9077
	if im.AnnualContributions != 9077 {
		t.Errorf("expected multi-year contributions=9077, got %d", im.AnnualContributions)
	}
	// True file counts across trees: AGP (5 .py + 3 .ipynb) + secret-client-audit (2 .py + 1 .qmd) + dataprofile fallback (1 prod) = 8 prod, 4 eda
	if lp.ExplorationFiles != 4 || lp.ProductionFiles != 8 {
		t.Errorf("expected ExplorationFiles=4 and ProductionFiles=8 from recursive trees, got eda=%d prod=%d", lp.ExplorationFiles, lp.ProductionFiles)
	}
	// Nested hygiene signals detected across Analytical Workspaces in AGP and secret-client-audit
	if lp.LockfileRepos != 1 || lp.CICDRepos != 1 || lp.TestSuiteRepos != 2 || lp.ContainerizedRepos != 1 {
		t.Errorf("expected recursive tree hygiene Lockfiles=1 CI=1 Tests=2 Docker=1, got Lockfiles=%d CI=%d Tests=%d Docker=%d",
			lp.LockfileRepos, lp.CICDRepos, lp.TestSuiteRepos, lp.ContainerizedRepos)
	}
	// Multiple Analytical Workspaces detected inside AGP
	if sp.TotalSources < 4 {
		t.Errorf("expected >=4 Analytical Workspaces scanned across Owned Repositories, got %d", sp.TotalSources)
	}

	// Privacy-preserving guarantee: zero private repo names, Analytical Workspace paths, or code leak into rendered SVGs
	portfolio := render.DefaultPortfolioData(domain.ProfileManifest{
		Subject: domain.SubjectConfig{Name: "Pablo Diego", Locale: "en"},
	})
	portfolio.Languages = lp
	portfolio.Stack = sp
	portfolio.Impact = im

	privateIdentifiers := []string{
		"AGP",
		"secret-client-audit",
		"Rio Bonito",
		"Empetur",
		"regressao.py",
		"SECRET_PROPRIETARY_FORMULA",
	}
	for _, school := range domain.AllSchools() {
		for _, panel := range domain.AllPanels() {
			svg, err := render.Render(school, panel, portfolio)
			if err != nil {
				t.Fatalf("Render(%s, %s) failed: %v", school, panel, err)
			}
			for _, secret := range privateIdentifiers {
				if strings.Contains(svg, secret) {
					t.Errorf("privacy leak: Render(%s, %s) leaked private identifier %q into SVG", school, panel, secret)
				}
			}
		}
	}
}

func TestCollectMonorepoAnalyticalWorkspacesVsSingleProjectRepo(t *testing.T) {
	mux := http.NewServeMux()
	var monorepoContentFetches int32

	mux.HandleFunc("/users/adalovelace/repos", func(w http.ResponseWriter, r *http.Request) {
		repos := []restRepo{
			{
				Name:          "client-monorepo",
				FullName:      "adalovelace/client-monorepo",
				Language:      "Python",
				DefaultBranch: "main",
			},
			{
				Name:          "single-project-ml",
				FullName:      "adalovelace/single-project-ml",
				Language:      "Python",
				DefaultBranch: "main",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(repos)
	})

	mux.HandleFunc("/repos/adalovelace/client-monorepo/languages", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]int64{"Jupyter Notebook": 500000, "Python": 120000})
	})
	mux.HandleFunc("/repos/adalovelace/single-project-ml/languages", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]int64{"Python": 180000})
	})

	// Repo 1: Multi-project monorepo with 5 analytical scripts/notebooks (+ 1 nested uv.lock) across 3 workspaces
	mux.HandleFunc("/repos/adalovelace/client-monorepo/git/trees/main", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"tree": []gitTreeEntry{
				{Path: "Reports/RetailChurn/eda.ipynb", Type: "blob", Size: 48000},
				{Path: "Reports/RetailChurn/glm_model.R", Type: "blob", Size: 3200},
				{Path: "Reports/SpatialDemand/choropleth.py", Type: "blob", Size: 4100},
				{Path: "Reports/SpatialDemand/spatial_queries.sql", Type: "blob", Size: 2100},
				{Path: "Utils/DataCLI/uv.lock", Type: "blob", Size: 14000},
				{Path: "Utils/DataCLI/cli.py", Type: "blob", Size: 2800},
			},
		})
	})

	// Repo 2: Standard single-project repository with 7 analytical entries in standard layout (root, src/, tests/, notebooks/)
	mux.HandleFunc("/repos/adalovelace/single-project-ml/git/trees/main", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"tree": []gitTreeEntry{
				{Path: "pyproject.toml", Type: "blob", Size: 900},
				{Path: "src/ingest.py", Type: "blob", Size: 3100},
				{Path: "src/features.py", Type: "blob", Size: 4200},
				{Path: "src/train.py", Type: "blob", Size: 5400},
				{Path: "src/evaluate.py", Type: "blob", Size: 2900},
				{Path: "tests/test_train.py", Type: "blob", Size: 1800},
				{Path: "notebooks/eda_overview.ipynb", Type: "blob", Size: 15000},
			},
		})
	})

	// Notebook with >32KB base64 PNG in cell 0 and imports in cell 1 (single-line JSON)
	bigB64 := strings.Repeat("iVBORw0KGgoAAAANSUhEUgAAAAUA", 1400)
	compactNotebookJSON := `{"metadata":{"kernelspec":{"language":"python"}},"cells":[{"cell_type":"code","source":["plt.plot([1,2,3])\n"],"outputs":[{"data":{"image/png":"` + bigB64 + `"}}]},{"cell_type":"code","source":["import statsmodels.api as sm\n","from sklearn.cluster import KMeans\n","import seaborn as sns\n"],"outputs":[]}]}`

	mux.HandleFunc("/repos/adalovelace/client-monorepo/contents/", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&monorepoContentFetches, 1)
		switch r.URL.Path {
		case "/repos/adalovelace/client-monorepo/contents/Reports/RetailChurn/eda.ipynb":
			_, _ = w.Write([]byte(compactNotebookJSON))
		case "/repos/adalovelace/client-monorepo/contents/Reports/RetailChurn/glm_model.R":
			_, _ = w.Write([]byte("library(tidyverse)\nlibrary(ggplot2)\n"))
		case "/repos/adalovelace/client-monorepo/contents/Reports/SpatialDemand/choropleth.py":
			_, _ = w.Write([]byte("import geopandas as gpd\nimport folium\nimport numpy as np\n"))
		case "/repos/adalovelace/client-monorepo/contents/Reports/SpatialDemand/spatial_queries.sql":
			_, _ = w.Write([]byte("SELECT ST_Area(geom) FROM postgis_table;\n"))
		case "/repos/adalovelace/client-monorepo/contents/Utils/DataCLI/uv.lock":
			_, _ = w.Write([]byte("name = \"typer\"\nname = \"pydantic\"\nname = \"openai-whisper\"\nname = \"pypandoc\"\n"))
		case "/repos/adalovelace/client-monorepo/contents/Utils/DataCLI/cli.py":
			_, _ = w.Write([]byte("import typer\nimport pydantic\n"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	mux.HandleFunc("/repos/adalovelace/single-project-ml/contents/pyproject.toml", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[project]
dependencies = ["geopandas>=0.14", "scipy>=1.12", "typer>=0.9", "scikit-learn>=1.4", "pypandoc>=1.13"]`))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	client := NewClient("")
	client.BaseURL = srv.URL

	_, stack, _, err := client.Collect(
		context.Background(),
		domain.GitHubSourceConfig{Username: "adalovelace"},
		domain.NormalizationConfig{DeflateNotebooks: true, DetectDataStack: true},
	)
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	// 3 Analytical Workspaces from client-monorepo + 1 single Analytical Workspace from single-project-ml = 4
	if stack.TotalSources != 4 {
		t.Errorf("expected TotalSources=4 (3 monorepo workspaces + 1 single-project repo), got %d", stack.TotalSources)
	}
	if fetches := atomic.LoadInt32(&monorepoContentFetches); fetches == 0 || fetches > 12 {
		t.Errorf("expected monorepo content fetches in [1..12], got %d", fetches)
	}

	// Verify nested uv.lock, compact .ipynb imports after 32KB base64, and Analytical Workspace scripts fed into stack detection
	foundData := make(map[string]int)
	for _, it := range stack.DataAnalytics {
		foundData[it.Name] = it.RepoCount
	}
	if foundData["GeoPandas / Spatial"] < 1 || foundData["ETL & CLI Pipelines"] < 1 {
		t.Errorf("expected GeoPandas / Spatial and ETL & CLI Pipelines detected from monorepo workspaces, got %+v", stack.DataAnalytics)
	}

	foundML := make(map[string]int)
	for _, it := range stack.MLAndAI {
		foundML[it.Name] = it.RepoCount
	}
	if foundML["Statsmodels / GLM"] < 1 || foundML["Clustering & PCA"] < 1 || foundML["LLMs & Speech AI"] < 1 {
		t.Errorf("expected Statsmodels / GLM, Clustering & PCA (from .ipynb cell after >32KB base64), and LLMs & Speech AI (from uv.lock), got %+v", stack.MLAndAI)
	}

	foundViz := make(map[string]int)
	for _, it := range stack.VizAndApps {
		foundViz[it.Name] = it.RepoCount
	}
	if foundViz["Geospatial Maps"] < 1 || foundViz["Report Automation"] < 1 {
		t.Errorf("expected Geospatial Maps and Report Automation detected in VizAndApps, got %+v", stack.VizAndApps)
	}
}

func TestCollectContributedRepositoriesReachAndMultiYearAllTimeContributions(t *testing.T) {
	mux := http.NewServeMux()
	var capturedQueries []string

	mux.HandleFunc("/user/repos", func(w http.ResponseWriter, r *http.Request) {
		repos := []restRepo{
			{
				Name:          "owned-public",
				FullName:      "pablodiegoo/owned-public",
				Private:       false,
				Fork:          false,
				Stargazers:    25,
				Forks:         4,
				Language:      "Python",
				DefaultBranch: "main",
			},
			{
				Name:          "owned-private",
				FullName:      "pablodiegoo/owned-private",
				Private:       true,
				Fork:          false,
				Stargazers:    0,
				Forks:         0,
				Language:      "R",
				DefaultBranch: "main",
			},
			{
				Name:          "iota-sdk",
				FullName:      "pablodiegoo/iota-sdk",
				Private:       false,
				Fork:          true,
				Stargazers:    1200,
				Forks:         310,
				Language:      "Go",
				DefaultBranch: "main",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(repos)
	})

	mux.HandleFunc("/repos/pablodiegoo/owned-public/languages", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]int64{"Python": 90000})
	})
	mux.HandleFunc("/repos/pablodiegoo/owned-private/languages", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]int64{"R": 45000})
	})

	mux.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		var reqBody struct {
			Query string `json:"query"`
		}
		_ = json.NewDecoder(r.Body).Decode(&reqBody)
		capturedQueries = append(capturedQueries, reqBody.Query)
		w.Header().Set("Content-Type", "application/json")

		if !strings.Contains(reqBody.Query, "y2026:") {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"user": map[string]interface{}{
						"contributionsCollection": map[string]interface{}{
							"contributionYears": []int{2026, 2025, 2024},
							"contributionCalendar": map[string]interface{}{
								"totalContributions": 8144,
								"weeks": []map[string]interface{}{
									{"contributionDays": []map[string]int{{"contributionCount": 15}}},
									{"contributionDays": []map[string]int{{"contributionCount": 30}}},
								},
							},
						},
						"repositoriesContributedTo": map[string]interface{}{
							"totalCount": 3,
							"nodes": []map[string]string{
								{"nameWithOwner": "iota-uz/iota-sdk"},
								{"nameWithOwner": "Agora-Pesquisa/agora-pesquisa.github.io"},
								{"nameWithOwner": "some-org/noisy-sandbox"},
							},
						},
					},
				},
			})
			return
		}

		// Multi-year aliased query response
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"user": map[string]interface{}{
					"y2026": map[string]interface{}{"contributionCalendar": map[string]int{"totalContributions": 8144}},
					"y2025": map[string]interface{}{"contributionCalendar": map[string]int{"totalContributions": 933}},
					"y2024": map[string]interface{}{"contributionCalendar": map[string]int{"totalContributions": 423}},
				},
			},
		})
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	client := NewClient("pat_token")
	client.BaseURL = srv.URL

	lp, _, im, err := client.Collect(
		context.Background(),
		domain.GitHubSourceConfig{
			Username:            "pablodiegoo",
			IncludeForks:        false,
			IncludePrivateCount: true,
			ExcludeRepos:        []string{"noisy-sandbox"},
		},
		domain.NormalizationConfig{DeflateNotebooks: true},
	)
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if len(capturedQueries) < 2 {
		t.Fatalf("expected 2 GraphQL queries (initial + multi-year), got %d", len(capturedQueries))
	}
	for _, token := range []string{
		"contributionYears",
		"repositoriesContributedTo",
		"includeUserRepositories: false",
		"contributionTypes: [COMMIT, PULL_REQUEST, REPOSITORY, PULL_REQUEST_REVIEW]",
	} {
		if !strings.Contains(capturedQueries[0], token) {
			t.Errorf("initial GraphQL query missing required clause %q:\n%s", token, capturedQueries[0])
		}
	}
	for _, alias := range []string{"y2026:", "y2025:", "y2024:"} {
		if !strings.Contains(capturedQueries[1], alias) {
			t.Errorf("multi-year GraphQL query missing year alias %q:\n%s", alias, capturedQueries[1])
		}
	}

	// All-time contributions = 8144 + 933 + 423 = 9500
	if im.AnnualContributions != 9500 {
		t.Errorf("expected multi-year AnnualContributions=9500, got %d", im.AnnualContributions)
	}

	// 2 Owned Repositories (owned-public + owned-private; iota-sdk fork excluded)
	if lp.TotalRepos != 2 {
		t.Errorf("expected lp.TotalRepos=2 Owned Repositories, got %d", lp.TotalRepos)
	}
	if im.TotalStars != 25 || im.TotalForks != 4 {
		t.Errorf("expected Owned Repository stars=25 forks=4 (without claiming fork/upstream stars), got stars=%d forks=%d",
			im.TotalStars, im.TotalForks)
	}

	// 2 Contributed Repositories after excluding "noisy-sandbox" by short name -> impact.PublicRepos = 2 Owned + 2 Contributed = 4
	if im.ContributedRepos != 2 || im.PublicRepos != 4 {
		t.Errorf("expected ContributedRepos=2 and PublicRepos=4 after excluding noisy-sandbox, got ContributedRepos=%d PublicRepos=%d",
			im.ContributedRepos, im.PublicRepos)
	}
}

