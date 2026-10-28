package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pablodiegoo/dataprofile/internal/config"
	"github.com/pablodiegoo/dataprofile/internal/render"
	"github.com/pablodiegoo/dataprofile/internal/sources/github"
	"github.com/pablodiegoo/dataprofile/internal/sources/huggingface"
	"github.com/pablodiegoo/dataprofile/internal/sources/kaggle"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "validate" || os.Args[1] == "generate") {
		sub := os.Args[1]
		os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
		if sub == "validate" {
			runValidate()
			return
		}
	}
	runGenerate()
}

func runGenerate() {
	var (
		configPath string
		envPath    string
		outputDir  string
		locale     string
		demoMode   bool
	)

	flag.StringVar(&configPath, "config", "profile.yaml", "Path to declarative profile.yaml manifest")
	flag.StringVar(&envPath, "env", ".env", "Path to optional local .env file")
	flag.StringVar(&outputDir, "output", "", "Override output directory from profile.yaml")
	flag.StringVar(&locale, "locale", "", "Override subject locale (en | pt-BR)")
	flag.BoolVar(&demoMode, "demo", false, "Render deterministic reference metrics without querying live network APIs")
	flag.Parse()

	manifest, err := config.LoadManifest(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not read %s (%v), using default manifest\n", configPath, err)
		manifest = config.DefaultManifest()
	}

	if outputDir != "" {
		manifest.Output.Directory = outputDir
	}
	if locale != "" {
		manifest.Subject.Locale = locale
	}

	secrets := config.LoadSecrets(envPath)
	data := render.DefaultPortfolioData(manifest)

	if !demoMode {
		ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
		defer cancel()

		// 1. Collect from GitHub if username configured
		if manifest.Sources.GitHub.Username != "" {
			ghClient := github.NewClient(secrets.GitHubToken)
			if lp, sp, im, err := ghClient.Collect(ctx, manifest.Sources.GitHub, manifest.Normalization); err == nil && lp.TotalRepos > 0 {
				if len(lp.Languages) > 0 {
					data.Languages = lp
				}
				if len(sp.DataAnalytics) > 0 && sp.DataAnalytics[0].RepoCount > 0 {
					data.Stack = sp
				}
				if im.PublicRepos > 0 {
					data.Impact.PublicRepos = im.PublicRepos
				}
				if im.TotalStars > 0 {
					data.Impact.TotalStars = im.TotalStars
				}
				if im.TotalForks > 0 {
					data.Impact.TotalForks = im.TotalForks
				}
				if im.AnnualContributions > 0 {
					data.Impact.AnnualContributions = im.AnnualContributions
				}
			}
		}

		// 2. Collect from Kaggle if enabled
		if manifest.Sources.Kaggle.Enabled && manifest.Sources.Kaggle.Username != "" {
			kgClient := kaggle.NewClient(secrets.KaggleUsername, secrets.KaggleKey)
			if ks, err := kgClient.Collect(ctx, manifest.Sources.Kaggle); err == nil && ks.Available && ks.TotalMedals() > 0 {
				data.Kaggle = ks
			}
		}

		// 3. Collect from Hugging Face Hub if enabled
		if manifest.Sources.HuggingFace.Enabled && manifest.Sources.HuggingFace.Username != "" {
			hfClient := huggingface.NewClient()
			if hf, err := hfClient.Collect(ctx, manifest.Sources.HuggingFace); err == nil && (hf.Models+hf.Datasets+hf.Spaces) > 0 {
				data.Impact.HFModels = hf.Models
				data.Impact.HFDatasets = hf.Datasets
				data.Impact.HFSpaces = hf.Spaces
				data.Impact.HFDownloads30d = hf.Downloads30d
				data.Impact.HFLikes = hf.Likes
			}
		}
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

	fmt.Printf("✓ Generated %d SVG panels across %d schools in %s/\n",
		totalGenerated, len(manifest.Output.Schools), manifest.Output.Directory)
}

func runValidate() {
	var (
		configPath string
		outputDir  string
	)
	flag.StringVar(&configPath, "config", "profile.yaml", "Path to declarative profile.yaml manifest")
	flag.StringVar(&outputDir, "output", "", "Override output directory from profile.yaml")
	flag.Parse()

	manifest, err := config.LoadManifest(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "validation failed loading %s: %v\n", configPath, err)
		os.Exit(1)
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
