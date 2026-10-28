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
