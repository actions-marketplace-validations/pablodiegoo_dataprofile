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
