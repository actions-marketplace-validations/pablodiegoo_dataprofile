package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadManifest(t *testing.T) {
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "profile.yaml")
	content := `
subject:
  name: "Pablo Diego"
  role: "Data Scientist"
  locale: "pt-BR"

sources:
  github:
    username: "pablodiegoo"
    include_forks: false
    include_private_count: true
    exclude_langs:
      - "HTML"
      - "CSS"
  kaggle:
    username: "pablodiegoo"
    enabled: true
    mode: "hybrid"
  huggingface:
    username: "pablodiegoo"
    enabled: true

normalization:
  deflate_notebooks: true
  detect_data_stack: true

output:
  directory: "generated"
  schools:
    - "tufte"
    - "few"
  panels:
    - "languages"
    - "kaggle"
`
	if err := os.WriteFile(manifestPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	m, err := LoadManifest(manifestPath)
	if err != nil {
		t.Fatalf("LoadManifest failed: %v", err)
	}

	if m.Subject.Name != "Pablo Diego" || m.Subject.Locale != "pt-BR" {
		t.Errorf("unexpected subject: %+v", m.Subject)
	}
	if m.Sources.GitHub.Username != "pablodiegoo" {
		t.Errorf("expected github username pablodiegoo, got %q", m.Sources.GitHub.Username)
	}
	if len(m.Output.Schools) != 2 || m.Output.Schools[0] != "tufte" || m.Output.Schools[1] != "few" {
		t.Errorf("unexpected schools: %+v", m.Output.Schools)
	}
	if len(m.Output.Panels) != 2 || m.Output.Panels[0] != "languages" || m.Output.Panels[1] != "kaggle" {
		t.Errorf("unexpected panels: %+v", m.Output.Panels)
	}
}

func TestLoadManifestInlineArrays(t *testing.T) {
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "profile.yaml")
	content := `
subject:
  name: "Ada Lovelace"
  role: "ML Engineer"
  locale: "en"
sources:
  github:
    username: "adalovelace"
    include_forks: false
    include_private_count: true
    exclude_repos: ["legacy-fork", "archive-repo"]
    exclude_langs: ["HTML", "CSS", "TeX"]
normalization:
  deflate_notebooks: true
output:
  directory: "generated"
  schools: ["tufte"]
  panels: ["languages"]
`
	if err := os.WriteFile(manifestPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	m, err := LoadManifest(manifestPath)
	if err != nil {
		t.Fatalf("LoadManifest failed: %v", err)
	}

	if len(m.Sources.GitHub.ExcludeRepos) != 2 || m.Sources.GitHub.ExcludeRepos[0] != "legacy-fork" || m.Sources.GitHub.ExcludeRepos[1] != "archive-repo" {
		t.Errorf("expected inline exclude_repos [legacy-fork archive-repo], got %+v", m.Sources.GitHub.ExcludeRepos)
	}
	if len(m.Sources.GitHub.ExcludeLangs) != 3 || m.Sources.GitHub.ExcludeLangs[2] != "TeX" {
		t.Errorf("expected inline exclude_langs [HTML CSS TeX], got %+v", m.Sources.GitHub.ExcludeLangs)
	}
	if len(m.Output.Schools) != 1 || m.Output.Schools[0] != "tufte" {
		t.Errorf("expected inline schools [tufte], got %+v", m.Output.Schools)
	}
	if len(m.Output.Panels) != 1 || m.Output.Panels[0] != "languages" {
		t.Errorf("expected inline panels [languages], got %+v", m.Output.Panels)
	}
}

func TestLoadSecrets(t *testing.T) {
	dir := t.TempDir()
	dotenvPath := filepath.Join(dir, ".env")
	dotenvContent := `
# Local secrets file
export GITHUB_TOKEN="ghp_from_dotenv" # inline comment
KAGGLE_USERNAME='kaggle_dotenv_user'
KAGGLE_KEY=kaggle_dotenv_key
`
	if err := os.WriteFile(dotenvPath, []byte(dotenvContent), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("ACCESS_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("KAGGLE_USERNAME", "")
	t.Setenv("KAGGLE_KEY", "")

	// 1. Fallback to .env when process env vars are unset
	s := LoadSecrets(dotenvPath)
	if s.GitHubToken != "ghp_from_dotenv" {
		t.Errorf("expected GitHubToken ghp_from_dotenv from .env, got %q", s.GitHubToken)
	}
	if s.KaggleUsername != "kaggle_dotenv_user" || s.KaggleKey != "kaggle_dotenv_key" {
		t.Errorf("unexpected Kaggle secrets from .env: %+v", s)
	}

	// 2. Process env overrides .env
	t.Setenv("GITHUB_TOKEN", "ghp_from_env")
	s = LoadSecrets(dotenvPath)
	if s.GitHubToken != "ghp_from_env" {
		t.Errorf("expected env GITHUB_TOKEN to override .env, got %q", s.GitHubToken)
	}

	// 3. ACCESS_TOKEN takes precedence over GITHUB_TOKEN
	t.Setenv("ACCESS_TOKEN", "pat_access_token")
	s = LoadSecrets(dotenvPath)
	if s.GitHubToken != "pat_access_token" {
		t.Errorf("expected ACCESS_TOKEN to take precedence over GITHUB_TOKEN, got %q", s.GitHubToken)
	}
}


