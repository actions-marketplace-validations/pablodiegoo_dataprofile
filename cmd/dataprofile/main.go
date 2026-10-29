package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pablodiegoo/dataprofile/internal/config"
	"github.com/pablodiegoo/dataprofile/internal/domain"
	"github.com/pablodiegoo/dataprofile/internal/render"
	"github.com/pablodiegoo/dataprofile/internal/sources/github"
	"github.com/pablodiegoo/dataprofile/internal/sources/huggingface"
	"github.com/pablodiegoo/dataprofile/internal/sources/kaggle"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "validate":
			os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
			runValidate()
			return
		case "sync":
			os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
			runSync()
			return
		case "generate":
			os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
		}
	}
	runGenerate()
}

func runGenerate() {
	var (
		manifestPath string
		envPath      string
		outputDir    string
		locale       string
		demoMode     bool
		noEnv        bool
	)

	flag.StringVar(&manifestPath, "config", "profile.yaml", "Path to declarative profile.yaml manifest")
	flag.StringVar(&envPath, "env", ".env", "Path to optional local .env file")
	flag.StringVar(&outputDir, "output", "", "Override output directory from profile.yaml")
	flag.StringVar(&locale, "locale", "", "Override subject locale (en | pt-BR)")
	flag.BoolVar(&demoMode, "demo", false, "Render deterministic reference metrics without querying live network APIs")
	flag.BoolVar(&noEnv, "no-env", false, "Ignore .env and environment variable overrides (use pure profile.yaml)")
	flag.Parse()

	manifest, err := config.LoadManifest(manifestPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not read %s (%v), using default manifest\n", manifestPath, err)
		manifest = config.DefaultManifest()
	}

	if !noEnv {
		config.ApplyEnvOverrides(&manifest, envPath)
	}

	if outputDir != "" {
		manifest.Output.Directory = outputDir
	}
	if locale != "" {
		manifest.Subject.Locale = locale
	}

	var secrets config.Secrets
	if !noEnv {
		secrets = config.LoadSecrets(envPath)
	}
	data := render.DefaultPortfolioData(manifest)

	if !demoMode {
		ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
		defer cancel()

		ghClient := github.NewClient(secrets.GitHubToken)
		kgClient := kaggle.NewClient(secrets.KaggleUsername, secrets.KaggleKey)
		hfClient := huggingface.NewClientWithToken(secrets.HFToken)
		collectLiveSources(ctx, manifest, &data, ghClient, kgClient, hfClient)
	}

	totalGenerated := 0
	for _, school := range manifest.Output.Schools {
		schoolDir := filepath.Join(manifest.Output.Directory, string(school))
		if err := os.MkdirAll(schoolDir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "error creating directory %s: %v\n", schoolDir, err)
			os.Exit(1)
		}

		for _, panel := range manifest.Output.Panels {
			svg, err := render.Render(school, panel, data)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error rendering %s/%s: %v\n", school, panel, err)
				os.Exit(1)
			}

			outPath := filepath.Join(schoolDir, fmt.Sprintf("%s.svg", panel))
			if err := os.WriteFile(outPath, []byte(svg), 0o644); err != nil {
				fmt.Fprintf(os.Stderr, "error writing %s: %v\n", outPath, err)
				os.Exit(1)
			}
			totalGenerated++
		}
	}

	readmePath := filepath.Join(manifest.Output.Directory, "README.md")
	if err := os.WriteFile(readmePath, []byte(render.RenderOutputREADME(manifest)), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing %s: %v\n", readmePath, err)
		os.Exit(1)
	}

	fmt.Printf("✓ Generated %d SVG panels across %d schools in %s/ (with %s)\n",
		totalGenerated, len(manifest.Output.Schools), manifest.Output.Directory, readmePath)
}

func runValidate() {
	var (
		manifestPath string
		envPath      string
		outputDir    string
		noEnv        bool
	)
	flag.StringVar(&manifestPath, "config", "profile.yaml", "Path to declarative profile.yaml manifest")
	flag.StringVar(&envPath, "env", ".env", "Path to optional local .env file")
	flag.StringVar(&outputDir, "output", "", "Override output directory from profile.yaml")
	flag.BoolVar(&noEnv, "no-env", false, "Ignore .env and environment variable overrides (use pure profile.yaml)")
	flag.Parse()

	manifest, err := config.LoadManifest(manifestPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "validation failed loading %s: %v\n", manifestPath, err)
		os.Exit(1)
	}
	if !noEnv {
		config.ApplyEnvOverrides(&manifest, envPath)
	}
	if outputDir != "" {
		manifest.Output.Directory = outputDir
	}

	checked := 0
	for _, school := range manifest.Output.Schools {
		for _, panel := range manifest.Output.Panels {
			p := filepath.Join(manifest.Output.Directory, string(school), fmt.Sprintf("%s.svg", panel))
			raw, err := os.ReadFile(p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "missing generated SVG %s: %v\n", p, err)
				os.Exit(1)
			}
			if err := render.ValidateSVG(string(raw)); err != nil {
				fmt.Fprintf(os.Stderr, "invalid SVG %s: %v\n", p, err)
				os.Exit(1)
			}
			checked++
		}
	}
	fmt.Printf("✓ Validated %d SVG files in %s/ (XML well-formed, high-contrast, zero meta-labels)\n", checked, manifest.Output.Directory)
}

func runSync() {
	var (
		envPath  string
		repo     string
		workflow string
		trigger  bool
		dryRun   bool
	)
	flag.StringVar(&envPath, "env", ".env", "Path to local .env file to sync with GitHub Actions")
	flag.StringVar(&repo, "repo", "", "Optional target GitHub repository (owner/repo)")
	flag.StringVar(&workflow, "workflow", "generate.yml", "GitHub Actions workflow file to trigger after syncing")
	flag.BoolVar(&trigger, "trigger", true, "Trigger workflow_dispatch after syncing secrets and variables")
	flag.BoolVar(&dryRun, "dry-run", false, "Print planned gh CLI commands without executing them")
	flag.Parse()

	envMap := config.LoadDotenvMap(envPath)
	targets := config.PlanGitHubSync(envMap)
	if len(targets) == 0 {
		fmt.Fprintf(os.Stderr, "error: no non-empty variables or secrets found in %s\n", envPath)
		os.Exit(1)
	}

	varsSynced, secretsSynced := 0, 0
	for _, t := range targets {
		kind := "variable"
		displayVal := t.Value
		if t.IsSecret {
			kind = "secret"
			displayVal = "***"
		}
		args := []string{kind, "set", t.Name, "--body", t.Value}
		if repo != "" {
			args = append(args, "--repo", repo)
		}
		if dryRun {
			fmt.Printf("[dry-run] gh %s set %s --body %q\n", kind, t.Name, displayVal)
		} else {
			cmd := exec.Command("gh", args...)
			if out, err := cmd.CombinedOutput(); err != nil {
				fmt.Fprintf(os.Stderr, "error syncing %s %s: %v\n%s", kind, t.Name, err, string(out))
				os.Exit(1)
			}
			fmt.Printf("✓ Synced GitHub %s: %s\n", kind, t.Name)
		}
		if t.IsSecret {
			secretsSynced++
		} else {
			varsSynced++
		}
	}

	if trigger {
		wfArgs := []string{"workflow", "run", workflow}
		if repo != "" {
			wfArgs = append(wfArgs, "--repo", repo)
		}
		if dryRun {
			fmt.Printf("[dry-run] gh workflow run %s\n", workflow)
		} else {
			cmd := exec.Command("gh", wfArgs...)
			if out, err := cmd.CombinedOutput(); err != nil {
				fmt.Fprintf(os.Stderr, "error triggering workflow %s: %v\n%s", workflow, err, string(out))
				os.Exit(1)
			}
			fmt.Printf("✓ Triggered GitHub Actions workflow: %s\n", workflow)
		}
	}

	fmt.Printf("✓ Sync complete (%d variables, %d secrets from %s)\n", varsSynced, secretsSynced, envPath)
}

func collectLiveSources(
	ctx context.Context,
	manifest domain.ProfileManifest,
	data *domain.PortfolioData,
	ghClient *github.Client,
	kgClient *kaggle.Client,
	hfClient *huggingface.Client,
) {
	// 1. Collect from GitHub if username configured
	if manifest.Sources.GitHub.Username != "" && ghClient != nil {
		if lp, sp, im, err := ghClient.Collect(ctx, manifest.Sources.GitHub, manifest.Normalization); err == nil {
			if len(lp.Languages) > 0 {
				data.Languages = lp
			}
			if len(sp.DataAnalytics) > 0 && sp.DataAnalytics[0].RepoCount > 0 {
				data.Stack = sp
			}
			data.Impact.PublicRepos = im.PublicRepos
			data.Impact.ContributedRepos = im.ContributedRepos
			data.Impact.TotalStars = im.TotalStars
			data.Impact.TotalForks = im.TotalForks
			data.Impact.AnnualContributions = im.AnnualContributions
			data.Impact.RepoViews14d = im.RepoViews14d
			data.Impact.ActiveStreakDays = im.ActiveStreakDays
			data.Impact.YoYGrowthPct = im.YoYGrowthPct
			data.Impact.LinesAdded = im.LinesAdded
			data.Impact.LinesDeleted = im.LinesDeleted
			if len(im.WeeklyTrend) > 0 {
				data.Impact.WeeklyTrend = im.WeeklyTrend
			}
		}
	}

	// 2. Collect from Kaggle if enabled
	if manifest.Sources.Kaggle.Enabled && manifest.Sources.Kaggle.Username != "" && kgClient != nil {
		if ks, err := kgClient.Collect(ctx, manifest.Sources.Kaggle); err == nil && ks.Available {
			data.Kaggle = ks
		}
	}

	// 3. Collect from Hugging Face Hub if enabled
	if !manifest.Sources.HuggingFace.Enabled || manifest.Sources.HuggingFace.Username == "" {
		data.Impact.HFModels = 0
		data.Impact.HFDatasets = 0
		data.Impact.HFSpaces = 0
		data.Impact.HFDownloads30d = 0
		data.Impact.HFLikes = 0
	} else if hfClient != nil {
		if hf, err := hfClient.Collect(ctx, manifest.Sources.HuggingFace); err == nil {
			data.Impact.HFModels = hf.Models
			data.Impact.HFDatasets = hf.Datasets
			data.Impact.HFSpaces = hf.Spaces
			data.Impact.HFDownloads30d = hf.Downloads30d
			data.Impact.HFLikes = hf.Likes
		}
	}
}
