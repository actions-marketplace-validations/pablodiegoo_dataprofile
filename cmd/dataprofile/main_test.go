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
  schools: ["tufte", "few", "cairo", "holmes"]
  panels: ["languages", "header", "bio", "bio-full", "full"]
`
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	binPath := filepath.Join(dir, "dataprofile-bin")
	buildCmd := exec.Command("go", "build", "-o", binPath, ".")
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build CLI binary: %v\n%s", err, string(out))
	}

	// 1. Run generate with --config, --output, --demo, and --no-env
	genCmd := exec.Command(binPath, "generate", "--config", manifestPath, "--output", outDir, "--demo", "--no-env")
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

	for _, school := range []string{"tufte", "few", "cairo", "holmes"} {
		hdrPath := filepath.Join(outDir, school, "header.svg")
		hdrBytes, err := os.ReadFile(hdrPath)
		if err != nil {
			t.Fatalf("expected generated header.svg at %s: %v", hdrPath, err)
		}
		hdrStr := string(hdrBytes)
		for _, tok := range []string{
			`viewBox="0 0 900 115"`,
			`<textPath xlink:href="#hdr-path0">Welcome! I&#39;m Ada Lovelace 👋</textPath>`,
			`attributeName="stroke-dasharray"`,
		} {
			if !strings.Contains(hdrStr, tok) {
				t.Errorf("expected %s in %s", tok, hdrPath)
			}
		}

		for _, bc := range []struct {
			file      string
			wantWidth string
		}{
			{"bio.svg", `width="440"`},
			{"bio-full.svg", `width="900"`},
		} {
			p := filepath.Join(outDir, school, bc.file)
			raw, err := os.ReadFile(p)
			if err != nil {
				t.Fatalf("expected generated %s at %s: %v", bc.file, p, err)
			}
			s := string(raw)
			if !strings.Contains(s, bc.wantWidth) || !strings.Contains(s, "Ada Lovelace is a Data Scientist") {
				t.Errorf("expected %s and fallback bio narrative in %s", bc.wantWidth, p)
			}
		}

		fullPath := filepath.Join(outDir, school, "full.svg")
		fullBytes, err := os.ReadFile(fullPath)
		if err != nil {
			t.Fatalf("expected generated full.svg at %s: %v", fullPath, err)
		}
		fullStr := string(fullBytes)
		for _, tok := range []string{
			`viewBox="0 0 900 595" width="900" height="595"`,
			`<textPath xlink:href="#full-hdr-path0">Welcome! I&#39;m Ada Lovelace 👋</textPath>`,
			`<g transform="translate(8, 121)">`,
			`<g transform="translate(452, 121)">`,
			`<g transform="translate(8, 359)">`,
			`<g transform="translate(452, 359)">`,
			"Pandas / Polars",
			"Notebooks",
		} {
			if !strings.Contains(fullStr, tok) {
				t.Errorf("expected %s in %s", tok, fullPath)
			}
		}
	}

	outReadme, err := os.ReadFile(filepath.Join(outDir, "README.md"))
	if err != nil {
		t.Fatalf("expected generated README.md in %s: %v", outDir, err)
	}
	for _, tok := range []string{
		"generated/tufte/header.svg",
		"generated/tufte/bio.svg",
		"generated/tufte/bio-full.svg",
		"generated/tufte/full.svg",
		"https://raw.githubusercontent.com/adalovelace/dataprofile/output/generated/tufte/full.svg",
	} {
		if !strings.Contains(string(outReadme), tok) {
			t.Errorf("expected %s in generated README.md", tok)
		}
	}

	// 2. Run validate with --config, --output, and --no-env
	valCmd := exec.Command(binPath, "validate", "--config", manifestPath, "--output", outDir, "--no-env")
	if out, err := valCmd.CombinedOutput(); err != nil {
		t.Fatalf("validate with --output failed: %v\n%s", err, string(out))
	}

	// 3. Run generate with --locale=pt-BR overriding manifest locale: "en"
	ptOutDir := filepath.Join(dir, "pt-out")
	genPtCmd := exec.Command(binPath, "generate", "--config", manifestPath, "--output", ptOutDir, "--locale=pt-BR", "--demo", "--no-env")
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
	ptHdrSVG, err := os.ReadFile(filepath.Join(ptOutDir, "cairo", "header.svg"))
	if err != nil {
		t.Fatalf("expected pt-BR generated cairo/header.svg: %v", err)
	}
	if !strings.Contains(string(ptHdrSVG), "Olá! Eu sou Ada Lovelace 👋") || !strings.Contains(string(ptHdrSVG), "↳ 1,428 contribuições anuais") {
		t.Errorf("expected --locale=pt-BR CLI flag to localize fallback headlines and callout in cairo/header.svg")
	}
	ptBioSVG, err := os.ReadFile(filepath.Join(ptOutDir, "cairo", "bio.svg"))
	if err != nil {
		t.Fatalf("expected pt-BR generated cairo/bio.svg: %v", err)
	}
	if !strings.Contains(string(ptBioSVG), "Ada Lovelace atua como Data Scientist") || !strings.Contains(string(ptBioSVG), "Trajetória Profissional &amp; Resumo Executivo") {
		t.Errorf("expected --locale=pt-BR CLI flag to localize fallback bio and title in cairo/bio.svg")
	}
	ptFullSVG, err := os.ReadFile(filepath.Join(ptOutDir, "cairo", "full.svg"))
	if err != nil {
		t.Fatalf("expected pt-BR generated cairo/full.svg: %v", err)
	}
	if !strings.Contains(string(ptFullSVG), "Olá! Eu sou Ada Lovelace 👋") || !strings.Contains(string(ptFullSVG), "Ecossistema Completo de Dados &amp; ML Ancorado por") {
		t.Errorf("expected --locale=pt-BR CLI flag to localize header and analytical panels in cairo/full.svg")
	}

	// 4. Verify .env override with [item 1; item 2] syntax, process env & CLI flag precedence, --no-env isolation, and `sync --dry-run`
	envFile := filepath.Join(dir, "custom.env")
	envContent := `
SUBJECT_NAME="Pablo Diego"
SUBJECT_LOCALE="en"
SUBJECT_HEADLINES=[Welcome people and mystical beings.; I'm Pablo! 👋; Working with data 10+ years]
SUBJECT_BIO=[Economist and Data Scientist.; Open-source tooling in Go and Python.]
GITHUB_USERNAME="pablodiegoo"
OUTPUT_DIRECTORY="dotenv-ignored-dir"
ACCESS_TOKEN="ghp_test_token"
`
	if err := os.WriteFile(envFile, []byte(envContent), 0o600); err != nil {
		t.Fatal(err)
	}
	envOutDir := filepath.Join(dir, "env-out")
	genEnvCmd := exec.Command(binPath, "generate", "--config", manifestPath, "--env", envFile, "--output", envOutDir, "--demo")
	if out, err := genEnvCmd.CombinedOutput(); err != nil {
		t.Fatalf("generate with --env failed: %v\n%s", err, string(out))
	}
	envHdr, err := os.ReadFile(filepath.Join(envOutDir, "tufte", "header.svg"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(envHdr), "I&#39;m Pablo! 👋") || strings.Contains(string(envHdr), "Ada Lovelace") {
		t.Errorf("expected --env override to replace Ada Lovelace with Pablo's headlines in header.svg")
	}
	valEnvCmd := exec.Command(binPath, "validate", "--config", manifestPath, "--env", envFile, "--output", envOutDir)
	if out, err := valEnvCmd.CombinedOutput(); err != nil {
		t.Fatalf("validate with --env failed: %v\n%s", err, string(out))
	}

	// Verify process env overrides .env while explicit CLI flags (--output, --locale) override process env and .env
	precOutDir := filepath.Join(dir, "prec-out")
	precCmd := exec.Command(binPath, "generate", "--config", manifestPath, "--env", envFile, "--output", precOutDir, "--locale=pt-BR", "--demo")
	precCmd.Env = append(os.Environ(),
		"SUBJECT_HEADLINES=[Headline from process env]",
		"SUBJECT_LOCALE=en",
		"OUTPUT_DIRECTORY="+filepath.Join(dir, "wrong-process-env-dir"),
	)
	if out, err := precCmd.CombinedOutput(); err != nil {
		t.Fatalf("generate precedence check failed: %v\n%s", err, string(out))
	}
	precHdr, err := os.ReadFile(filepath.Join(precOutDir, "cairo", "header.svg"))
	if err != nil {
		t.Fatalf("expected --output CLI flag to take precedence over OUTPUT_DIRECTORY: %v", err)
	}
	if !strings.Contains(string(precHdr), "Headline from process env") || !strings.Contains(string(precHdr), "contribuições anuais") {
		t.Errorf("expected process env SUBJECT_HEADLINES to override .env and --locale=pt-BR CLI flag to override SUBJECT_LOCALE=en")
	}

	// Verify --no-env isolates Showcase renders from both --env and process env overrides on generate and validate
	noEnvOutDir := filepath.Join(dir, "no-env-showcase")
	noEnvGenCmd := exec.Command(binPath, "generate", "--config", manifestPath, "--env", envFile, "--output", noEnvOutDir, "--demo", "--no-env")
	noEnvGenCmd.Env = append(os.Environ(),
		"SUBJECT_NAME=Process Env Intruder",
		"SUBJECT_HEADLINES=[Intruder Headline]",
		"OUTPUT_SCHOOLS=[tufte]",
		"OUTPUT_PANELS=[header]",
	)
	if out, err := noEnvGenCmd.CombinedOutput(); err != nil {
		t.Fatalf("generate --no-env failed: %v\n%s", err, string(out))
	}
	noEnvHdr, err := os.ReadFile(filepath.Join(noEnvOutDir, "holmes", "header.svg"))
	if err != nil {
		t.Fatalf("expected --no-env to ignore OUTPUT_SCHOOLS=[tufte] and render holmes/header.svg: %v", err)
	}
	if !strings.Contains(string(noEnvHdr), "Ada Lovelace") || strings.Contains(string(noEnvHdr), "I&#39;m Pablo! 👋") || strings.Contains(string(noEnvHdr), "Intruder") {
		t.Errorf("expected --no-env to render pure profile.yaml (Ada Lovelace) without .env or process env overrides")
	}
	noEnvValCmd := exec.Command(binPath, "validate", "--config", manifestPath, "--env", envFile, "--output", noEnvOutDir, "--no-env")
	noEnvValCmd.Env = noEnvGenCmd.Env
	if out, err := noEnvValCmd.CombinedOutput(); err != nil {
		t.Fatalf("validate --no-env failed: %v\n%s", err, string(out))
	}

	syncCmd := exec.Command(binPath, "sync", "--env", envFile, "--dry-run")
	syncOut, err := syncCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("sync --dry-run failed: %v\n%s", err, string(syncOut))
	}
	for _, want := range []string{
		`[dry-run] gh variable set GH_USERNAME --body "pablodiegoo"`,
		`[dry-run] gh variable set SUBJECT_NAME --body "Pablo Diego"`,
		`[dry-run] gh secret set ACCESS_TOKEN --body "***"`,
		`[dry-run] gh workflow run generate.yml`,
	} {
		if !strings.Contains(string(syncOut), want) {
			t.Errorf("expected %q in sync --dry-run output, got:\n%s", want, string(syncOut))
		}
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

func assertFileContains(t *testing.T, root, relPath string, requiredTokens []string) {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(root, relPath))
	if err != nil {
		t.Fatalf("missing %s: %v", relPath, err)
	}
	content := string(raw)
	for _, tok := range requiredTokens {
		if !strings.Contains(content, tok) {
			t.Errorf("%s missing expected token %q", relPath, tok)
		}
	}
}

func TestGitHubActionWorkflowAndSecretHygiene(t *testing.T) {
	root := repoRoot(t)

	// 1. Verify action.yml composite GitHub Action contract
	assertFileContains(t, root, "action.yml", []string{
		`using: "composite"`,
		`actions/setup-go@v5`,
		`go-version: "1.23"`,
		`config:`,
		`output:`,
		`locale:`,
		`demo:`,
		`no-env:`,
		`go run ./cmd/dataprofile generate`,
		`--config=${{ inputs.config }}`,
		`--output=${{ inputs.output }}`,
		`--locale=${{ inputs.locale }}`,
		`--no-env`,
		`VAL_ARGS="--no-env"`,
		`go run ./cmd/dataprofile validate --config="${{ inputs.config }}" --output="${{ inputs.output }}" $VAL_ARGS`,
	})

	// 2. Verify .github/workflows/generate.yml CI/CD & orphan output branch workflow
	assertFileContains(t, root, filepath.Join(".github", "workflows", "generate.yml"), []string{
		`cron:`,
		`workflow_dispatch:`,
		`branches:`,
		`- main`,
		`contents: write`,
		`go test -v ./...`,
		`ACCESS_TOKEN: ${{ secrets.ACCESS_TOKEN }}`,
		`GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}`,
		`KAGGLE_USERNAME: ${{ secrets.KAGGLE_USERNAME || vars.KAGGLE_USERNAME }}`,
		`KAGGLE_KEY: ${{ secrets.KAGGLE_KEY }}`,
		`HF_TOKEN: ${{ secrets.HF_TOKEN }}`,
		`SUBJECT_NAME: ${{ vars.SUBJECT_NAME }}`,
		`SUBJECT_ROLE: ${{ vars.SUBJECT_ROLE }}`,
		`SUBJECT_LOCALE: ${{ vars.SUBJECT_LOCALE }}`,
		`SUBJECT_HEADLINES: ${{ vars.SUBJECT_HEADLINES }}`,
		`SUBJECT_BIO: ${{ vars.SUBJECT_BIO }}`,
		`GH_USERNAME: ${{ vars.GH_USERNAME }}`,
		`GH_INCLUDE_FORKS: ${{ vars.GH_INCLUDE_FORKS }}`,
		`GH_INCLUDE_PRIVATE_COUNT: ${{ vars.GH_INCLUDE_PRIVATE_COUNT }}`,
		`GH_EXCLUDE_REPOS: ${{ vars.GH_EXCLUDE_REPOS }}`,
		`GH_EXCLUDE_LANGS: ${{ vars.GH_EXCLUDE_LANGS }}`,
		`KAGGLE_ENABLED: ${{ vars.KAGGLE_ENABLED }}`,
		`KAGGLE_MODE: ${{ vars.KAGGLE_MODE }}`,
		`HF_USERNAME: ${{ vars.HF_USERNAME }}`,
		`HF_ENABLED: ${{ vars.HF_ENABLED }}`,
		`DEFLATE_NOTEBOOKS: ${{ vars.DEFLATE_NOTEBOOKS }}`,
		`DETECT_DATA_STACK: ${{ vars.DETECT_DATA_STACK }}`,
		`OUTPUT_DIRECTORY: ${{ vars.OUTPUT_DIRECTORY }}`,
		`OUTPUT_SCHOOLS: ${{ vars.OUTPUT_SCHOOLS }}`,
		`OUTPUT_PANELS: ${{ vars.OUTPUT_PANELS }}`,
		`go run ./cmd/dataprofile generate --config=profile.yaml --output=generated`,
		`go run ./cmd/dataprofile validate --config=profile.yaml --output=generated`,
		`go run ./cmd/dataprofile generate --config=profile.yaml --demo --no-env --output=showcase`,
		`go run ./cmd/dataprofile validate --config=profile.yaml --no-env --output=showcase`,
		`cp -r "$TMP_DIR/generated" ./generated`,
		`cp -r "$TMP_DIR/showcase" ./showcase`,
		`git add generated showcase README.md`,
		`git checkout --orphan`,
		`git push -f origin output`,
	})

	// 3. Verify .gitignore secret and generated output exclusion rules and confirm generated/ and showcase/ are untracked on main
	assertFileContains(t, root, ".gitignore", []string{
		".env",
		".env.*",
		"!.env.example",
		"kaggle.json",
		"/generated/",
		"/showcase/",
	})
	lsCmd := exec.Command("git", "-C", root, "ls-files", "generated", "showcase")
	if trackedOut, err := lsCmd.CombinedOutput(); err != nil {
		t.Fatalf("git ls-files failed: %v\n%s", err, string(trackedOut))
	} else if strings.TrimSpace(string(trackedOut)) != "" {
		t.Errorf("expected zero tracked files under generated/ or showcase/ on main, found:\n%s", string(trackedOut))
	}

	// 4. Verify .env.example documents all optional secrets and [item 1; item 2] overrides without leaking values
	assertFileContains(t, root, ".env.example", []string{
		"ACCESS_TOKEN=",
		"GITHUB_TOKEN=",
		"KAGGLE_USERNAME=",
		"KAGGLE_KEY=",
		"HF_TOKEN=",
		"SUBJECT_HEADLINES=",
		"SUBJECT_BIO=",
	})
}

func TestPrototypesShowcaseAndGeneratedSVGs(t *testing.T) {
	root := repoRoot(t)

	// 1. Render all 36 Showcase SVGs (4 Schools x 9 Panels) from profile.yaml in demo mode and verify against prototypes/schools/
	manifest, err := config.LoadManifest(filepath.Join(root, "profile.yaml"))
	if err != nil {
		t.Fatalf("failed loading root profile.yaml: %v", err)
	}
	if len(domain.AllSchools()) != 4 || len(domain.AllPanels()) != 9 {
		t.Fatalf("expected 4 Schools and 9 Panels (36 SVGs), got %d Schools and %d Panels", len(domain.AllSchools()), len(domain.AllPanels()))
	}
	data := render.DefaultPortfolioData(manifest)

	renderedCount := 0
	for _, school := range domain.AllSchools() {
		for _, panel := range domain.AllPanels() {
			svgStr, err := render.Render(school, panel, data)
			if err != nil {
				t.Fatalf("failed rendering %s/%s: %v", school, panel, err)
			}
			if err := render.ValidateSVG(svgStr); err != nil {
				t.Errorf("rendered SVG %s/%s failed ValidateSVG: %v", school, panel, err)
			}
			renderedCount++

			// For the 4 analytical panels (languages, stack, kaggle, impact), verify prototypes/schools/<school>/<panel>.svg matches
			if panel == domain.PanelLanguages || panel == domain.PanelStack || panel == domain.PanelKaggle || panel == domain.PanelImpact {
				protoPath := filepath.Join(root, "prototypes", "schools", string(school), string(panel)+".svg")
				protoBytes, err := os.ReadFile(protoPath)
				if err != nil {
					t.Fatalf("missing prototype SVG %s: %v", protoPath, err)
				}
				if err := render.ValidateSVG(string(protoBytes)); err != nil {
					t.Errorf("prototype SVG %s failed ValidateSVG: %v", protoPath, err)
				}
				if string(protoBytes) != svgStr {
					t.Errorf("prototype SVG %s is out of sync with rendered %s/%s", protoPath, school, panel)
				}
			}
		}
	}
	if renderedCount != 36 {
		t.Errorf("expected 36 rendered Showcase SVGs across 4 Schools and 9 Panels, got %d", renderedCount)
	}

	// 2. Verify all remaining SVGs in prototypes/ pass ValidateSVG
	err = filepath.Walk(filepath.Join(root, "prototypes"), func(path string, info os.FileInfo, err error) error {
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
	assertFileContains(t, root, filepath.Join("prototypes", "index.html"), []string{
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
	})

	// 4. Verify flagship README.md showcase & documentation
	assertFileContains(t, root, "README.md", []string{
		"36 SVGs",
		"languages.svg",
		"stack.svg",
		"kaggle.svg",
		"impact.svg",
		"summary.svg",
		"header.svg",
		"bio.svg",
		"bio-full.svg",
		"full.svg",
		"Edward Tufte",
		"Stephen Few",
		"Alberto Cairo",
		"Nigel Holmes",
		"output/showcase/few/summary.svg",
		"output/showcase/tufte/summary.svg",
		"output/showcase/tufte/languages.svg",
		"output/showcase/few/languages.svg",
		"output/showcase/cairo/languages.svg",
		"output/showcase/holmes/languages.svg",
		"output/generated/",
		"2-Minute Setup Guide",
		"profile.yaml",
		"[item 1; item 2]",
		"go run ./cmd/dataprofile generate",
		"go run ./cmd/dataprofile sync",
		"go run ./cmd/dataprofile validate",
		"--demo",
		"--no-env",
		"--output=showcase",
		"--locale=pt-BR",
		"go test -v ./...",
		"GLOSSARY.md",
		"docs/adr/",
		"ADR-0005",
	})

	// 5. Verify Open Source Guide, Domain Docs, ADRs & GitHub Community Standards files
	communityFiles := map[string][]string{
		"LICENSE":            {"Apache License", "Version 2.0"},
		"NOTICE":             {"https://github.com/pablodiegoo/dataprofile", "Pablo Diego", "Apache License, Version 2.0"},
		"CODE_OF_CONDUCT.md": {"Contributor Covenant Code of Conduct", "Our pledge"},
		"CONTRIBUTING.md": {
			"GLOSSARY.md",
			"ADR-0005",
			"36 Showcase SVG panels (4 Schools x 9 Panels)",
			"header", "bio", "bio-full", "full",
			"--demo --no-env --output=showcase",
			"output/showcase/",
			"output/generated/",
			"go run ./cmd/dataprofile sync --dry-run",
			"go test -v ./...",
			"Apache License 2.0",
		},
		"SECURITY.md": {
			"Security Policy",
			"Reporting a vulnerability",
			"ACCESS_TOKEN",
			"[item 1; item 2]",
			"PlanGitHubSync",
			"/generated/",
			"/showcase/",
			"ADR-0005",
		},
		"CITATION.cff": {"cff-version: 1.2.0", "Apache-2.0", "https://github.com/pablodiegoo/dataprofile"},
		"CHANGELOG.md": {
			"v1.0.0",
			"Apache License 2.0",
			"ADR-0005",
			"[item 1; item 2]",
			"--no-env",
			"dataprofile sync",
			"output/showcase/",
			"output/generated/",
			"header", "bio", "bio-full", "full",
		},
		".github/PULL_REQUEST_TEMPLATE.md": {
			"GLOSSARY.md",
			"NOTICE",
			"go test -v ./...",
			"--demo --no-env --output=showcase",
			"Showcase",
			"ADR-0005",
		},
		".github/ISSUE_TEMPLATE/config.yml": {"https://github.com/pablodiegoo/dataprofile/discussions"},
		"GLOSSARY.md": {
			"Subject",
			"Source",
			"Panel",
			"School",
			"Palette",
			"Profile Manifest",
			"Showcase",
			"Owned Repository",
			"Passive Fork",
			"Contributed Repository",
			"Notebook Deflation",
			"Analytical Workspace",
		},
		"docs/adr/0005-env-overrides-and-orphan-output-branch.md": {
			"[item 1; item 2]",
			"dataprofile sync",
			"output/showcase/",
			"output/generated/",
			"--demo --no-env",
		},
	}
	for relPath, requiredTokens := range communityFiles {
		assertFileContains(t, root, relPath, requiredTokens)
	}
}



