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

	// Engineering hygiene signals populated
	if langPanel.LockfileRepos <= 0 || langPanel.CICDRepos <= 0 || langPanel.TestSuiteRepos <= 0 || langPanel.ContainerizedRepos <= 0 {
		t.Errorf("expected positive engineering hygiene counts, got %+v", langPanel)
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
