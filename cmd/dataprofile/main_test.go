package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
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
