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
	"github.com/pablodiegoo/dataprofile/internal/sources/kaggle"
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

	// 3. Run generate with --locale=pt-BR overriding manifest locale: "en"
	ptOutDir := filepath.Join(dir, "pt-out")
	genPtCmd := exec.Command(binPath, "generate", "--config", manifestPath, "--output", ptOutDir, "--locale=pt-BR", "--demo")
	if out, err := genPtCmd.CombinedOutput(); err != nil {
		t.Fatalf("generate --locale=pt-BR failed: %v\n%s", err, string(out))
	}
	ptSVG, err := os.ReadFile(filepath.Join(ptOutDir, "tufte", "languages.svg"))
	if err != nil {
		t.Fatalf("expected pt-BR generated SVG: %v", err)
	}
	if !strings.Contains(string(ptSVG), "Distribuição de Código") {
		t.Errorf("expected --locale=pt-BR CLI flag to render Portuguese title in languages.svg")
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

	hfClient := huggingface.NewClientWithToken("hf_live_token")
	hfClient.BaseURL = hfZeroSrv.URL

	data := render.DefaultPortfolioData(manifest)
	collectLiveSources(context.Background(), manifest, &data, ghClient, nil, hfClient)

	if data.Impact.AnnualContributions != 2150 {
		t.Errorf("expected AnnualContributions=2150, got %d", data.Impact.AnnualContributions)
	}
	if data.Impact.ContributedRepos != 2 || data.Impact.PublicRepos != 3 || data.Languages.TotalRepos != 1 {
		t.Errorf("expected ContributedRepos=2, Impact.PublicRepos=3 (1 Owned + 2 Contributed), Languages.TotalRepos=1, got ContributedRepos=%d PublicRepos=%d TotalRepos=%d",
			data.Impact.ContributedRepos, data.Impact.PublicRepos, data.Languages.TotalRepos)
	}
	if len(data.Impact.WeeklyTrend) != 4 || data.Impact.WeeklyTrend[3] != 55 {
		t.Errorf("expected WeeklyTrend=[10 25 40 55], got %v", data.Impact.WeeklyTrend)
	}
	if data.Impact.RepoViews14d != 480*5+62*8 {
		t.Errorf("expected live RepoViews14d=%d, got %d", 480*5+62*8, data.Impact.RepoViews14d)
	}
	if data.Impact.ActiveStreakDays != 4 {
		t.Errorf("expected live ActiveStreakDays=4, got %d", data.Impact.ActiveStreakDays)
	}
	if data.Impact.YoYGrowthPct < 100 {
		t.Errorf("expected live YoYGrowthPct > 100 from [10,25,40,55], got %.1f", data.Impact.YoYGrowthPct)
	}
	if data.Impact.LinesAdded <= 0 || data.Impact.LinesDeleted <= 0 || data.Impact.LinesAdded == 84200 {
		t.Errorf("expected live LinesAdded/LinesDeleted copied instead of demo 84200, got +%d/-%d", data.Impact.LinesAdded, data.Impact.LinesDeleted)
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

func TestLiveCollectionZeroGoldKaggleAndZeroStarGitHubDoesNotLeakDemoValues(t *testing.T) {
	ghMux := http.NewServeMux()
	ghMux.HandleFunc("/users/newuser/repos", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode([]map[string]interface{}{
			{
				"name":             "starter-repo",
				"full_name":        "newuser/starter-repo",
				"fork":             false,
				"stargazers_count": 0,
				"forks_count":      0,
				"language":         "Python",
			},
		})
	})
	ghMux.HandleFunc("/repos/newuser/starter-repo/languages", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]int64{"Python": 25000})
	})
	ghSrv := httptest.NewServer(ghMux)
	defer ghSrv.Close()

	manifest := config.DefaultManifest()
	manifest.Sources.GitHub.Username = "newuser"
	manifest.Sources.Kaggle.Enabled = true
	manifest.Sources.Kaggle.Username = "newuser"
	manifest.Sources.Kaggle.Mode = domain.KaggleModeBrowser
	manifest.Sources.HuggingFace.Enabled = false

	ghClient := github.NewClient("")
	ghClient.BaseURL = ghSrv.URL

	kgClient := kaggle.NewClient("", "")
	kgClient.DOMDumper = func(ctx context.Context, username string) (string, error) {
		return `Notebooks Contributor #9,800 0 Gold 0 Silver 0 Bronze Datasets Novice Competitions Novice Discussions Novice`, nil
	}

	data := render.DefaultPortfolioData(manifest)
	collectLiveSources(context.Background(), manifest, &data, ghClient, kgClient, nil)

	if data.Impact.TotalStars != 0 || data.Impact.TotalForks != 0 {
		t.Errorf("expected live zero-star/zero-fork GitHub profile to have 0 stars and 0 forks (not demo 342/78), got stars=%d forks=%d",
			data.Impact.TotalStars, data.Impact.TotalForks)
	}
	if data.Impact.PublicRepos != 1 {
		t.Errorf("expected PublicRepos=1, got %d", data.Impact.PublicRepos)
	}
	if data.Impact.RepoViews14d != 0 {
		t.Errorf("expected RepoViews14d=0 for zero-star/zero-fork profile, got %d", data.Impact.RepoViews14d)
	}
	if data.Kaggle.OverallTier != domain.KaggleTierContributor || data.Kaggle.TotalGold != 0 || data.Kaggle.TotalMedals() != 0 {
		t.Errorf("expected live zero-gold Kaggle profile to overwrite demo Master standing, got %+v", data.Kaggle)
	}

	svgKaggle, err := render.Render(domain.SchoolFew, domain.PanelKaggle, data)
	if err != nil {
		t.Fatalf("Render kaggle failed: %v", err)
	}
	if !strings.Contains(svgKaggle, "4 Tracks") || strings.Contains(svgKaggle, "3 Gold") {
		t.Errorf("expected live zero-gold Kaggle profile to trigger '4 Tracks' degradation instead of leaking demo '3 Gold'")
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("failed to resolve repo root: %v", err)
	}
	return root
}

func TestGitHubActionWorkflowAndSecretHygiene(t *testing.T) {
	root := repoRoot(t)

	// 1. Verify action.yml composite GitHub Action contract
	actionRaw, err := os.ReadFile(filepath.Join(root, "action.yml"))
	if err != nil {
		t.Fatalf("missing action.yml: %v", err)
	}
	actionStr := string(actionRaw)
	requiredActionTokens := []string{
		`using: "composite"`,
		`actions/setup-go@v5`,
		`go-version: "1.23"`,
		`config:`,
		`output:`,
		`locale:`,
		`go run ./cmd/dataprofile generate`,
		`--config=${{ inputs.config }}`,
		`--output=${{ inputs.output }}`,
		`--locale=${{ inputs.locale }}`,
		`go run ./cmd/dataprofile validate --config="${{ inputs.config }}" --output="${{ inputs.output }}"`,
	}
	for _, tok := range requiredActionTokens {
		if !strings.Contains(actionStr, tok) {
			t.Errorf("action.yml missing required token %q", tok)
		}
	}

	// 2. Verify .github/workflows/generate.yml CI/CD & orphan output branch workflow
	wfRaw, err := os.ReadFile(filepath.Join(root, ".github", "workflows", "generate.yml"))
	if err != nil {
		t.Fatalf("missing .github/workflows/generate.yml: %v", err)
	}
	wfStr := string(wfRaw)
	requiredWorkflowTokens := []string{
		`cron:`,
		`workflow_dispatch:`,
		`branches:`,
		`- main`,
		`contents: write`,
		`go test -v ./...`,
		`ACCESS_TOKEN: ${{ secrets.ACCESS_TOKEN }}`,
		`GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}`,
		`KAGGLE_USERNAME: ${{ secrets.KAGGLE_USERNAME }}`,
		`KAGGLE_KEY: ${{ secrets.KAGGLE_KEY }}`,
		`HF_TOKEN: ${{ secrets.HF_TOKEN }}`,
		`go run ./cmd/dataprofile generate`,
		`go run ./cmd/dataprofile validate`,
		`git checkout --orphan`,
		`git push -f origin output`,
	}
	for _, tok := range requiredWorkflowTokens {
		if !strings.Contains(wfStr, tok) {
			t.Errorf(".github/workflows/generate.yml missing required token %q", tok)
		}
	}

	// 3. Verify .gitignore secret exclusion rules
	gitignoreRaw, err := os.ReadFile(filepath.Join(root, ".gitignore"))
	if err != nil {
		t.Fatalf("missing .gitignore: %v", err)
	}
	gitignoreStr := string(gitignoreRaw)
	for _, pattern := range []string{".env", ".env.*", "!.env.example", "kaggle.json"} {
		if !strings.Contains(gitignoreStr, pattern) {
			t.Errorf(".gitignore missing required secret exclusion rule %q", pattern)
		}
	}

	// 4. Verify .env.example documents all optional secrets without leaking values
	envExRaw, err := os.ReadFile(filepath.Join(root, ".env.example"))
	if err != nil {
		t.Fatalf("missing .env.example: %v", err)
	}
	envExStr := string(envExRaw)
	for _, key := range []string{"ACCESS_TOKEN=", "GITHUB_TOKEN=", "KAGGLE_USERNAME=", "KAGGLE_KEY=", "HF_TOKEN="} {
		if !strings.Contains(envExStr, key) {
			t.Errorf(".env.example missing documented variable %q", key)
		}
	}
}

func TestPrototypesShowcaseAndGeneratedSVGs(t *testing.T) {
	root := repoRoot(t)

	// 1. Verify all 20 generated/<school>/<panel>.svg exist and pass ValidateSVG
	for _, school := range domain.AllSchools() {
		for _, panel := range domain.AllPanels() {
			genPath := filepath.Join(root, "generated", string(school), string(panel)+".svg")
			genBytes, err := os.ReadFile(genPath)
			if err != nil {
				t.Fatalf("missing generated SVG %s: %v", genPath, err)
			}
			if err := render.ValidateSVG(string(genBytes)); err != nil {
				t.Errorf("generated SVG %s failed ValidateSVG: %v", genPath, err)
			}

			// For the 4 analytical panels (languages, stack, kaggle, impact), verify prototypes/schools/<school>/<panel>.svg matches
			if panel != domain.PanelSummary {
				protoPath := filepath.Join(root, "prototypes", "schools", string(school), string(panel)+".svg")
				protoBytes, err := os.ReadFile(protoPath)
				if err != nil {
					t.Fatalf("missing prototype SVG %s: %v", protoPath, err)
				}
				if err := render.ValidateSVG(string(protoBytes)); err != nil {
					t.Errorf("prototype SVG %s failed ValidateSVG: %v", protoPath, err)
				}
				if string(protoBytes) != string(genBytes) {
					t.Errorf("prototype SVG %s is out of sync with %s", protoPath, genPath)
				}
			}
		}
	}

	// 2. Verify all remaining SVGs in prototypes/ pass ValidateSVG
	err := filepath.Walk(filepath.Join(root, "prototypes"), func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".svg") {
			return err
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if vErr := render.ValidateSVG(string(raw)); vErr != nil {
			t.Errorf("prototype SVG %s failed ValidateSVG: %v", path, vErr)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("error walking prototypes directory: %v", err)
	}

	// 3. Verify interactive prototype viewer (prototypes/index.html)
	protoHTMLRaw, err := os.ReadFile(filepath.Join(root, "prototypes", "index.html"))
	if err != nil {
		t.Fatalf("missing prototypes/index.html: %v", err)
	}
	protoHTML := string(protoHTMLRaw)
	requiredHTMLTokens := []string{
		"#0d1117",
		"#ffffff",
		"🌙 Dark GitHub",
		"☀️ Light GitHub",
		"ArrowRight",
		"ArrowLeft",
		"←",
		"→",
		"schools/tufte",
		"schools/few",
		"schools/cairo",
		"schools/holmes",
		"languages.svg",
		"stack.svg",
		"kaggle.svg",
		"impact.svg",
	}
	for _, tok := range requiredHTMLTokens {
		if !strings.Contains(protoHTML, tok) {
			t.Errorf("prototypes/index.html missing required token %q", tok)
		}
	}

	// 4. Verify flagship README.md showcase & documentation
	readmeRaw, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatalf("missing README.md: %v", err)
	}
	readmeStr := string(readmeRaw)
	requiredReadmeTokens := []string{
		"languages.svg",
		"stack.svg",
		"kaggle.svg",
		"impact.svg",
		"Edward Tufte",
		"Stephen Few",
		"Alberto Cairo",
		"Nigel Holmes",
		"generated/tufte/languages.svg",
		"generated/few/languages.svg",
		"generated/cairo/languages.svg",
		"generated/holmes/languages.svg",
		"2-Minute Setup Guide",
		"profile.yaml",
		"go run ./cmd/dataprofile generate",
		"go run ./cmd/dataprofile validate",
		"--demo",
		"--locale=pt-BR",
		"go test -v ./...",
		"GLOSSARY.md",
		"docs/adr/",
	}
	for _, tok := range requiredReadmeTokens {
		if !strings.Contains(readmeStr, tok) {
			t.Errorf("README.md missing required section or reference %q", tok)
		}
	}

	// 5. Verify Open Source Guide & GitHub Community Standards files
	communityFiles := map[string][]string{
		"LICENSE":                           {"Apache License", "Version 2.0"},
		"NOTICE":                            {"https://github.com/pablodiegoo/dataprofile", "Pablo Diego", "Apache License, Version 2.0"},
		"CODE_OF_CONDUCT.md":                {"Contributor Covenant Code of Conduct", "Our pledge"},
		"CONTRIBUTING.md":                   {"GLOSSARY.md", "go test -v ./...", "Apache License 2.0"},
		"SECURITY.md":                       {"Security Policy", "Reporting a vulnerability", "ACCESS_TOKEN"},
		"CITATION.cff":                      {"cff-version: 1.2.0", "Apache-2.0", "https://github.com/pablodiegoo/dataprofile"},
		"CHANGELOG.md":                      {"v1.0.0", "Apache License 2.0"},
		".github/PULL_REQUEST_TEMPLATE.md":  {"GLOSSARY.md", "NOTICE", "go test -v ./..."},
		".github/ISSUE_TEMPLATE/config.yml": {"https://github.com/pablodiegoo/dataprofile/discussions"},
	}
	for relPath, requiredTokens := range communityFiles {
		raw, err := os.ReadFile(filepath.Join(root, relPath))
		if err != nil {
			t.Fatalf("missing community standard file %s: %v", relPath, err)
		}
		content := string(raw)
		for _, tok := range requiredTokens {
			if !strings.Contains(content, tok) {
				t.Errorf("%s missing expected token %q", relPath, tok)
			}
		}
	}
}



