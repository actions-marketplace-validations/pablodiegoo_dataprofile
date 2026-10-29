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
HUGGINGFACE_TOKEN="hf_dotenv_fallback"
`
	if err := os.WriteFile(dotenvPath, []byte(dotenvContent), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("ACCESS_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("KAGGLE_USERNAME", "")
	t.Setenv("KAGGLE_KEY", "")
	t.Setenv("HF_TOKEN", "")
	t.Setenv("HUGGINGFACE_TOKEN", "")

	// 1. Fallback to .env when process env vars are unset
	s := LoadSecrets(dotenvPath)
	if s.GitHubToken != "ghp_from_dotenv" {
		t.Errorf("expected GitHubToken ghp_from_dotenv from .env, got %q", s.GitHubToken)
	}
	if s.KaggleUsername != "kaggle_dotenv_user" || s.KaggleKey != "kaggle_dotenv_key" {
		t.Errorf("unexpected Kaggle secrets from .env: %+v", s)
	}
	if s.HFToken != "hf_dotenv_fallback" {
		t.Errorf("expected HFToken hf_dotenv_fallback from .env HUGGINGFACE_TOKEN, got %q", s.HFToken)
	}

	// 2. Process env overrides .env
	t.Setenv("GITHUB_TOKEN", "ghp_from_env")
	t.Setenv("HUGGINGFACE_TOKEN", "hf_from_env_alias")
	s = LoadSecrets(dotenvPath)
	if s.GitHubToken != "ghp_from_env" {
		t.Errorf("expected env GITHUB_TOKEN to override .env, got %q", s.GitHubToken)
	}
	if s.HFToken != "hf_from_env_alias" {
		t.Errorf("expected env HUGGINGFACE_TOKEN to override .env, got %q", s.HFToken)
	}

	// 3. ACCESS_TOKEN takes precedence over GITHUB_TOKEN, and HF_TOKEN takes precedence over HUGGINGFACE_TOKEN
	t.Setenv("ACCESS_TOKEN", "pat_access_token")
	t.Setenv("HF_TOKEN", "hf_primary_token")
	s = LoadSecrets(dotenvPath)
	if s.GitHubToken != "pat_access_token" {
		t.Errorf("expected ACCESS_TOKEN to take precedence over GITHUB_TOKEN, got %q", s.GitHubToken)
	}
	if s.HFToken != "hf_primary_token" {
		t.Errorf("expected HF_TOKEN to take precedence over HUGGINGFACE_TOKEN, got %q", s.HFToken)
	}
}

func TestDefaultManifestExcludeLangsAndIncludePrivateAlias(t *testing.T) {
	def := DefaultManifest()
	wantExcludeLangs := []string{
		"HTML", "CSS", "SCSS", "TeX", "G-code",
		"Roff", "Liquid", "Blade", "Rich Text Format", "Batchfile",
	}
	if len(def.Sources.GitHub.ExcludeLangs) != len(wantExcludeLangs) {
		t.Fatalf("expected %d default exclude_langs, got %d (%v)",
			len(wantExcludeLangs), len(def.Sources.GitHub.ExcludeLangs), def.Sources.GitHub.ExcludeLangs)
	}
	for i, want := range wantExcludeLangs {
		if def.Sources.GitHub.ExcludeLangs[i] != want {
			t.Errorf("default ExcludeLangs[%d] = %q, want %q", i, def.Sources.GitHub.ExcludeLangs[i], want)
		}
	}

	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "profile.yaml")
	content := `
sources:
  github:
    username: "pablodiegoo"
    include_private: false
`
	if err := os.WriteFile(manifestPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := LoadManifest(manifestPath)
	if err != nil {
		t.Fatalf("LoadManifest failed: %v", err)
	}
	if m.Sources.GitHub.IncludePrivateCount {
		t.Errorf("expected include_private: false alias to set IncludePrivateCount=false")
	}
	if len(m.Sources.GitHub.ExcludeLangs) != len(wantExcludeLangs) {
		t.Errorf("expected omitted exclude_langs to retain %d defaults, got %v", len(wantExcludeLangs), m.Sources.GitHub.ExcludeLangs)
	}
}

func TestLoadManifestHeadlinesAndBioFormats(t *testing.T) {
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "profile.yaml")
	content := `
subject:
  name: "Pablo Diego"
  role: "Economist | Data Scientist"
  locale: "pt-BR"
  headlines:
    - "Welcome people and mystical beings."
    - "I'm Pablo! 👋"
    - "Working with data 10+ years"
  bio:
    - "Economista e Cientista de Dados com mais de 10 anos de experiência."
    - "Focado em modelagem estatística, machine learning e engenharia analítica."
output:
  panels:
    - "header"
    - "bio"
    - "bio-full"
    - "full"
`
	if err := os.WriteFile(manifestPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := LoadManifest(manifestPath)
	if err != nil {
		t.Fatalf("LoadManifest failed: %v", err)
	}
	if len(m.Subject.Headlines) != 3 || m.Subject.Headlines[1] != "I'm Pablo! 👋" {
		t.Errorf("unexpected Subject.Headlines: %+v", m.Subject.Headlines)
	}
	if len(m.Subject.Bio) != 2 || m.Subject.Bio[0] != "Economista e Cientista de Dados com mais de 10 anos de experiência." {
		t.Errorf("unexpected Subject.Bio: %+v", m.Subject.Bio)
	}
	if len(m.Output.Panels) != 4 || m.Output.Panels[2] != "bio-full" || m.Output.Panels[3] != "full" {
		t.Errorf("unexpected Output.Panels: %+v", m.Output.Panels)
	}

	// Test block scalar (|) for bio (including a line starting with a literal quote)
	blockPath := filepath.Join(dir, "block.yaml")
	blockContent := `
subject:
  name: "Ada Lovelace"
  headlines: ["Analytical Engine", "First Programmer"]
  bio: |
    "Data is key," she wrote in her notes.
    Second line of career narrative.
`
	if err := os.WriteFile(blockPath, []byte(blockContent), 0o644); err != nil {
		t.Fatal(err)
	}
	mb, err := LoadManifest(blockPath)
	if err != nil {
		t.Fatalf("LoadManifest block failed: %v", err)
	}
	if len(mb.Subject.Headlines) != 2 || mb.Subject.Headlines[0] != "Analytical Engine" {
		t.Errorf("unexpected inline Subject.Headlines: %+v", mb.Subject.Headlines)
	}
	if len(mb.Subject.Bio) != 2 || mb.Subject.Bio[0] != `"Data is key," she wrote in her notes.` || mb.Subject.Bio[1] != "Second line of career narrative." {
		t.Errorf("unexpected block Subject.Bio: %+v", mb.Subject.Bio)
	}

	// Test block scalar (|- and >) with literal '#' and \n-escaped scalar for headlines
	headlinesBlockPath := filepath.Join(dir, "headlines-block.yaml")
	headlinesBlockContent := `
subject:
  name: "Grace Hopper"
  headlines: |-
    First compiler pioneer
    Kaggle Master (#412) & COBOL
    Debugging since 1947
`
	if err := os.WriteFile(headlinesBlockPath, []byte(headlinesBlockContent), 0o644); err != nil {
		t.Fatal(err)
	}
	mhb, err := LoadManifest(headlinesBlockPath)
	if err != nil {
		t.Fatalf("LoadManifest headlines-block failed: %v", err)
	}
	wantHeadlines := []string{
		"First compiler pioneer",
		"Kaggle Master (#412) & COBOL",
		"Debugging since 1947",
	}
	if len(mhb.Subject.Headlines) != len(wantHeadlines) {
		t.Fatalf("expected %d block headlines, got %d (%+v)", len(wantHeadlines), len(mhb.Subject.Headlines), mhb.Subject.Headlines)
	}
	for i, w := range wantHeadlines {
		if mhb.Subject.Headlines[i] != w {
			t.Errorf("Subject.Headlines[%d] = %q, want %q", i, mhb.Subject.Headlines[i], w)
		}
	}

	headlinesEscapedPath := filepath.Join(dir, "headlines-escaped.yaml")
	if err := os.WriteFile(headlinesEscapedPath, []byte("subject:\n  headlines: \"Line one\\nLine two\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	mhe, err := LoadManifest(headlinesEscapedPath)
	if err != nil {
		t.Fatalf("LoadManifest headlines-escaped failed: %v", err)
	}
	if len(mhe.Subject.Headlines) != 2 || mhe.Subject.Headlines[0] != "Line one" || mhe.Subject.Headlines[1] != "Line two" {
		t.Errorf("expected 2 escaped headlines [Line one, Line two], got %+v", mhe.Subject.Headlines)
	}

	// Test inline array with commas inside quotes for subject.bio
	bioInlineCommasPath := filepath.Join(dir, "bio-inline-commas.yaml")
	bioInlineCommasContent := `
subject:
  name: "Ada Lovelace"
  bio: ["Economist, Data Scientist, and Engineer.", 'Focused on Python, SQL, Go, and R (#412).']
`
	if err := os.WriteFile(bioInlineCommasPath, []byte(bioInlineCommasContent), 0o644); err != nil {
		t.Fatal(err)
	}
	mbic, err := LoadManifest(bioInlineCommasPath)
	if err != nil {
		t.Fatalf("LoadManifest bio-inline-commas failed: %v", err)
	}
	wantBioCommas := []string{
		"Economist, Data Scientist, and Engineer.",
		"Focused on Python, SQL, Go, and R (#412).",
	}
	if len(mbic.Subject.Bio) != len(wantBioCommas) {
		t.Fatalf("expected %d inline bio items with commas, got %d (%+v)", len(wantBioCommas), len(mbic.Subject.Bio), mbic.Subject.Bio)
	}
	for i, w := range wantBioCommas {
		if mbic.Subject.Bio[i] != w {
			t.Errorf("Subject.Bio[%d] = %q, want %q", i, mbic.Subject.Bio[i], w)
		}
	}

	// Test multi-line string scalar (both \n-escaped and physical multi-line continuation) and |- block scalar with '#' for bio
	bioScalarPath := filepath.Join(dir, "bio-scalar.yaml")
	bioScalarContent := `
subject:
  name: "Grace Hopper"
  bio: "First paragraph with Python, SQL, and Go.\nSecond paragraph with Kaggle Master (#412)."
`
	if err := os.WriteFile(bioScalarPath, []byte(bioScalarContent), 0o644); err != nil {
		t.Fatal(err)
	}
	mbs, err := LoadManifest(bioScalarPath)
	if err != nil {
		t.Fatalf("LoadManifest bio-scalar failed: %v", err)
	}
	if len(mbs.Subject.Bio) != 2 || mbs.Subject.Bio[0] != "First paragraph with Python, SQL, and Go." || mbs.Subject.Bio[1] != "Second paragraph with Kaggle Master (#412)." {
		t.Errorf("unexpected scalar Subject.Bio: %+v", mbs.Subject.Bio)
	}

	bioMultilinePath := filepath.Join(dir, "bio-multiline.yaml")
	bioMultilineContent := `
subject:
  name: "Katherine Johnson"
  bio: "Orbital mechanics pioneer at NASA.
    Computed trajectories for Project Mercury and Apollo 11."
`
	if err := os.WriteFile(bioMultilinePath, []byte(bioMultilineContent), 0o644); err != nil {
		t.Fatal(err)
	}
	mbm, err := LoadManifest(bioMultilinePath)
	if err != nil {
		t.Fatalf("LoadManifest bio-multiline failed: %v", err)
	}
	if len(mbm.Subject.Bio) != 2 || mbm.Subject.Bio[0] != "Orbital mechanics pioneer at NASA." || mbm.Subject.Bio[1] != "Computed trajectories for Project Mercury and Apollo 11." {
		t.Errorf("unexpected physical multiline Subject.Bio: %+v", mbm.Subject.Bio)
	}
}

func TestParseEnvListAndApplyEnvOverrides(t *testing.T) {
	items, ok := ParseEnvList(`[Welcome people and mystical beings.; I'm Pablo! 👋; Working with data 10+ years\; and counting; Kaggle Master (#412)] # trailing comment`)
	if !ok || len(items) != 4 {
		t.Fatalf("expected 4 items from ParseEnvList, got ok=%v items=%+v", ok, items)
	}
	if items[1] != "I'm Pablo! 👋" || items[2] != "Working with data 10+ years; and counting" || items[3] != "Kaggle Master (#412)" {
		t.Errorf("unexpected parsed items: %+v", items)
	}

	emptyItems, okEmpty := ParseEnvList("[] # empty list comment")
	if !okEmpty || len(emptyItems) != 0 {
		t.Errorf("expected empty slice for [], got ok=%v items=%+v", okEmpty, emptyItems)
	}

	if _, okScalar := ParseEnvList("plain string without brackets"); okScalar {
		t.Errorf("expected ParseEnvList to return false for plain scalar")
	}

	dir := t.TempDir()
	dotenvPath := filepath.Join(dir, ".env")
	dotenvContent := `
SUBJECT_NAME="Pablo Diego"
SUBJECT_ROLE=I'm a Data Scientist & Engineers' Lead # inline comment after apostrophes
SUBJECT_LOCALE="pt-BR"
SUBJECT_HEADLINES=[Welcome people and mystical beings.; I'm Pablo! 👋; Kaggle Master (#412)] # inline comment
SUBJECT_BIO='[Economist and Data Scientist with 10+ years.; I'm focused on Python, SQL, Go, and R (#1).]' # single-quoted bracket list
GITHUB_USERNAME="pablodiegoo"
GITHUB_INCLUDE_FORKS=true
GITHUB_INCLUDE_PRIVATE_COUNT=false
GITHUB_EXCLUDE_REPOS=[old-fork; scratch-repo]
GITHUB_EXCLUDE_LANGS=[HTML; CSS]
KAGGLE_USERNAME="pablodiegoo"
KAGGLE_ENABLED=false
KAGGLE_MODE="api"
HF_USERNAME="pablodiegoo-hf"
HF_ENABLED=false
DEFLATE_NOTEBOOKS=false
DETECT_DATA_STACK=false
OUTPUT_DIRECTORY="custom-generated"
OUTPUT_SCHOOLS=[tufte; cairo]
OUTPUT_PANELS=[summary; full]
ACCESS_TOKEN="ghp_secret_pat"
GITHUB_TOKEN="ghp_ignored_by_sync"
`
	if err := os.WriteFile(dotenvPath, []byte(dotenvContent), 0o600); err != nil {
		t.Fatal(err)
	}

	// Ensure process env doesn't interfere
	for _, k := range []string{
		"SUBJECT_NAME", "SUBJECT_ROLE", "SUBJECT_LOCALE", "SUBJECT_HEADLINES", "SUBJECT_BIO",
		"GITHUB_USERNAME", "GH_USERNAME", "GITHUB_INCLUDE_FORKS", "GH_INCLUDE_FORKS",
		"GITHUB_INCLUDE_PRIVATE_COUNT", "GH_INCLUDE_PRIVATE_COUNT",
		"GITHUB_EXCLUDE_REPOS", "GH_EXCLUDE_REPOS", "GITHUB_EXCLUDE_LANGS", "GH_EXCLUDE_LANGS",
		"KAGGLE_USERNAME", "KAGGLE_ENABLED", "KAGGLE_MODE", "HF_USERNAME", "HF_ENABLED",
		"DEFLATE_NOTEBOOKS", "DETECT_DATA_STACK", "OUTPUT_DIRECTORY", "OUTPUT_SCHOOLS", "OUTPUT_PANELS",
	} {
		t.Setenv(k, "")
	}

	m := DefaultManifest()
	ApplyEnvOverrides(&m, dotenvPath)

	if m.Subject.Name != "Pablo Diego" || m.Subject.Role != "I'm a Data Scientist & Engineers' Lead" || m.Subject.Locale != "pt-BR" {
		t.Errorf("unexpected Subject after .env override: %+v", m.Subject)
	}
	if len(m.Subject.Headlines) != 3 || m.Subject.Headlines[2] != "Kaggle Master (#412)" {
		t.Errorf("unexpected Headlines after .env override: %+v", m.Subject.Headlines)
	}
	if len(m.Subject.Bio) != 2 || m.Subject.Bio[1] != "I'm focused on Python, SQL, Go, and R (#1)." {
		t.Errorf("unexpected Bio after .env override: %+v", m.Subject.Bio)
	}
	if m.Sources.GitHub.Username != "pablodiegoo" || !m.Sources.GitHub.IncludeForks || m.Sources.GitHub.IncludePrivateCount {
		t.Errorf("unexpected GitHub config after .env override: %+v", m.Sources.GitHub)
	}
	if len(m.Sources.GitHub.ExcludeRepos) != 2 || m.Sources.GitHub.ExcludeRepos[1] != "scratch-repo" {
		t.Errorf("unexpected ExcludeRepos after .env override: %+v", m.Sources.GitHub.ExcludeRepos)
	}
	if len(m.Sources.GitHub.ExcludeLangs) != 2 || m.Sources.GitHub.ExcludeLangs[0] != "HTML" {
		t.Errorf("unexpected ExcludeLangs after .env override: %+v", m.Sources.GitHub.ExcludeLangs)
	}
	if m.Sources.Kaggle.Username != "pablodiegoo" || m.Sources.Kaggle.Enabled || m.Sources.Kaggle.Mode != "api" {
		t.Errorf("unexpected Kaggle config after .env override: %+v", m.Sources.Kaggle)
	}
	if m.Sources.HuggingFace.Username != "pablodiegoo-hf" || m.Sources.HuggingFace.Enabled {
		t.Errorf("unexpected HuggingFace config after .env override: %+v", m.Sources.HuggingFace)
	}
	if m.Normalization.DeflateNotebooks || m.Normalization.DetectDataStack {
		t.Errorf("unexpected Normalization config after .env override: %+v", m.Normalization)
	}
	if m.Output.Directory != "custom-generated" || len(m.Output.Schools) != 2 || len(m.Output.Panels) != 2 {
		t.Errorf("unexpected Output config after .env override: %+v", m.Output)
	}

	// Verify GH_* and HUGGINGFACE_* process env aliases override .env values and empty [] clears slices
	t.Setenv("GH_USERNAME", "ci-variable-user")
	t.Setenv("GH_INCLUDE_FORKS", "false")
	t.Setenv("GH_INCLUDE_PRIVATE", "true")
	t.Setenv("GH_EXCLUDE_REPOS", "[ci-excluded-repo]")
	t.Setenv("GH_EXCLUDE_LANGS", "[]")
	t.Setenv("HUGGINGFACE_USERNAME", "ci-hf-user")
	t.Setenv("HUGGINGFACE_ENABLED", "true")
	ApplyEnvOverrides(&m, dotenvPath)
	if m.Sources.GitHub.Username != "ci-variable-user" {
		t.Errorf("expected GH_USERNAME process env to override .env, got %q", m.Sources.GitHub.Username)
	}
	if m.Sources.GitHub.IncludeForks || !m.Sources.GitHub.IncludePrivateCount {
		t.Errorf("expected GH_INCLUDE_FORKS=false and GH_INCLUDE_PRIVATE=true from process env, got %+v", m.Sources.GitHub)
	}
	if len(m.Sources.GitHub.ExcludeRepos) != 1 || m.Sources.GitHub.ExcludeRepos[0] != "ci-excluded-repo" {
		t.Errorf("expected GH_EXCLUDE_REPOS=[ci-excluded-repo] from process env, got %+v", m.Sources.GitHub.ExcludeRepos)
	}
	if len(m.Sources.GitHub.ExcludeLangs) != 0 {
		t.Errorf("expected GH_EXCLUDE_LANGS=[] to clear ExcludeLangs, got %+v", m.Sources.GitHub.ExcludeLangs)
	}
	if m.Sources.HuggingFace.Username != "ci-hf-user" || !m.Sources.HuggingFace.Enabled {
		t.Errorf("expected HUGGINGFACE_* process env overrides, got %+v", m.Sources.HuggingFace)
	}

	// Verify PlanGitHubSync maps GITHUB_* to GH_* and classifies secrets vs variables
	targets := PlanGitHubSync(LoadDotenvMap(dotenvPath))
	foundGHUser := false
	foundSecret := false
	for _, tgt := range targets {
		if tgt.Name == "GITHUB_TOKEN" || tgt.Name == "GITHUB_USERNAME" {
			t.Errorf("unexpected reserved GITHUB_ target name in PlanGitHubSync: %+v", tgt)
		}
		if tgt.Name == "GH_USERNAME" && !tgt.IsSecret && tgt.Value == "pablodiegoo" {
			foundGHUser = true
		}
		if tgt.Name == "ACCESS_TOKEN" && tgt.IsSecret && tgt.Value == "ghp_secret_pat" {
			foundSecret = true
		}
	}
	if !foundGHUser || !foundSecret {
		t.Errorf("expected GH_USERNAME variable and ACCESS_TOKEN secret in sync plan, got %+v", targets)
	}
}




