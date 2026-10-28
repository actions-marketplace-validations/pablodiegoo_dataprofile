package github

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/pablodiegoo/dataprofile/internal/domain"
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
