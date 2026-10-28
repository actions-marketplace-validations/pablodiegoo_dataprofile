package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pablodiegoo/dataprofile/internal/config"
	"github.com/pablodiegoo/dataprofile/internal/domain"
	"github.com/pablodiegoo/dataprofile/internal/render"
	"github.com/pablodiegoo/dataprofile/internal/sources/github"
	"github.com/pablodiegoo/dataprofile/internal/sources/huggingface"
)

func TestCLIGenerateAndValidate(t *testing.T) {
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "profile.yaml")
	outDir := filepath.Join(dir, "custom-out")

	manifest := `
subject:
  name: "Ada Lovelace"
  role: "Data Scientist"
  locale: "en"
sources:
  github:
    username: "adalovelace"
    include_forks: false
    include_private_count: true
    exclude_repos: ["old-fork"]
    exclude_langs: ["HTML", "CSS"]
normalization:
  deflate_notebooks: true
output:
  directory: "default-ignored-dir"
  schools: ["tufte"]
  panels: ["languages"]
`
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	binPath := filepath.Join(dir, "dataprofile-bin")
	buildCmd := exec.Command("go", "build", "-o", binPath, ".")
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build CLI binary: %v\n%s", err, string(out))
	}

	// 1. Run generate with --config, --output, and --demo
	genCmd := exec.Command(binPath, "generate", "--config", manifestPath, "--output", outDir, "--demo")
	if out, err := genCmd.CombinedOutput(); err != nil {
		t.Fatalf("generate failed: %v\n%s", err, string(out))
	}

	svgPath := filepath.Join(outDir, "tufte", "languages.svg")
	rawSVG, err := os.ReadFile(svgPath)
	if err != nil {
		t.Fatalf("expected generated SVG at %s: %v", svgPath, err)
	}
	svgStr := string(rawSVG)
	if !strings.Contains(svgStr, "#526071") || !strings.Contains(svgStr, "#64748b") {
		t.Errorf("expected Universal Dual-Surface Palette (#526071 and #64748b) in %s", svgPath)
	}
	if strings.Contains(svgStr, "<rect") {
		t.Errorf("expected Tufte languages.svg to have zero outer box/cage (<rect>), found one")
	}

	// 2. Run validate with --config and --output
	valCmd := exec.Command(binPath, "validate", "--config", manifestPath, "--output", outDir)
	if out, err := valCmd.CombinedOutput(); err != nil {
		t.Fatalf("validate with --output failed: %v\n%s", err, string(out))
	}
}

func TestLiveCollectionImpactWiringAndZeroHFFallback(t *testing.T) {
	ghMux := http.NewServeMux()
	ghMux.HandleFunc("/users/adalovelace/repos", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode([]map[string]interface{}{
			{
				"name":             "forecast-engine",
				"full_name":        "adalovelace/forecast-engine",
				"fork":             false,
				"stargazers_count": 480,
				"forks_count":      62,
				"language":         "Python",
			},
		})
	})
	ghMux.HandleFunc("/repos/adalovelace/forecast-engine/languages", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]int64{"Python": 500000})
	})
	ghMux.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"user": map[string]interface{}{
					"contributionsCollection": map[string]interface{}{
						"contributionCalendar": map[string]interface{}{
							"totalContributions": 2150,
							"weeks": []map[string]interface{}{
								{"contributionDays": []map[string]int{{"contributionCount": 10}}},
								{"contributionDays": []map[string]int{{"contributionCount": 25}}},
								{"contributionDays": []map[string]int{{"contributionCount": 40}}},
								{"contributionDays": []map[string]int{{"contributionCount": 55}}},
							},
						},
					},
				},
			},
		})
	})
	ghSrv := httptest.NewServer(ghMux)
	defer ghSrv.Close()

	// Hugging Face server returning zero artifacts
	hfZeroSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer hfZeroSrv.Close()

	manifest := config.DefaultManifest()
	manifest.Sources.GitHub.Username = "adalovelace"
	manifest.Sources.Kaggle.Enabled = false
	manifest.Sources.HuggingFace.Enabled = true
	manifest.Sources.HuggingFace.Username = "adalovelace"

	ghClient := github.NewClient("ghp_live_token")
	ghClient.BaseURL = ghSrv.URL

	hfClient := huggingface.NewClient()
	hfClient.BaseURL = hfZeroSrv.URL

	data := render.DefaultPortfolioData(manifest)
	collectLiveSources(context.Background(), manifest, &data, ghClient, nil, hfClient)

	if data.Impact.AnnualContributions != 2150 {
		t.Errorf("expected AnnualContributions=2150, got %d", data.Impact.AnnualContributions)
	}
	if len(data.Impact.WeeklyTrend) != 4 || data.Impact.WeeklyTrend[3] != 55 {
		t.Errorf("expected WeeklyTrend=[10 25 40 55], got %v", data.Impact.WeeklyTrend)
	}
	if data.Impact.HFModels != 0 || data.Impact.HFDatasets != 0 || data.Impact.HFSpaces != 0 {
		t.Errorf("expected zero HF counts when live HF API returns [], got %d/%d/%d",
			data.Impact.HFModels, data.Impact.HFDatasets, data.Impact.HFSpaces)
	}

	svgImpact, err := render.Render(domain.SchoolTufte, domain.PanelImpact, data)
	if err != nil {
		t.Fatalf("Render impact failed: %v", err)
	}
	if !strings.Contains(svgImpact, "Engineering Velocity") || strings.Contains(svgImpact, "0 / 0 / 0") {
		t.Errorf("expected live zero-HF profile to render Engineering Velocity fallback in impact.svg")
	}
}

