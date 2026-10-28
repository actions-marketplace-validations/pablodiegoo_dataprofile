package render

import (
	"strings"
	"testing"

	"github.com/pablodiegoo/dataprofile/internal/config"
	"github.com/pablodiegoo/dataprofile/internal/domain"
)

func TestRenderAllSchoolsAndPanels(t *testing.T) {
	manifest := config.DefaultManifest()
	data := DefaultPortfolioData(manifest)

	forbiddenMeta := []string{"Edward Tufte", "Stephen Few", "Alberto Cairo", "Nigel Holmes"}

	for _, school := range domain.AllSchools() {
		for _, panel := range domain.AllPanels() {
			svg, err := Render(school, panel, data)
			if err != nil {
				t.Fatalf("Render(%s, %s) failed: %v", school, panel, err)
			}
			if !strings.HasPrefix(svg, "<svg") || !strings.HasSuffix(strings.TrimSpace(svg), "</svg>") {
				t.Errorf("Render(%s, %s) did not produce a root <svg> element", school, panel)
			}
			for _, meta := range forbiddenMeta {
				if strings.Contains(svg, meta) {
					t.Errorf("Render(%s, %s) contains forbidden School meta-label %q", school, panel, meta)
				}
			}
		}
	}
}

func TestRecruiterFirstGracefulDegradation(t *testing.T) {
	manifest := config.DefaultManifest()
	data := DefaultPortfolioData(manifest)

	// Case 1: 0 Gold medals -> should show "4 Tracks" instead of "0 Gold" across all 4 Schools
	data.Kaggle.TotalGold = 0
	for _, school := range domain.AllSchools() {
		svgKaggle, err := Render(school, domain.PanelKaggle, data)
		if err != nil {
			t.Fatalf("Render(%s, kaggle) failed: %v", school, err)
		}
		if err := ValidateSVG(svgKaggle); err != nil {
			t.Errorf("ValidateSVG(%s, kaggle) failed: %v", school, err)
		}
		if strings.Contains(svgKaggle, "0 Gold") {
			t.Errorf("Render(%s, kaggle): expected Recruiter-First degradation to avoid highlighting '0 Gold'", school)
		}
		if !strings.Contains(svgKaggle, "4 Tracks") {
			t.Errorf("Render(%s, kaggle): expected '4 Tracks' fallback KPI when TotalGold == 0", school)
		}
		for _, cat := range []string{"Notebooks", "Datasets", "Competitions", "Discussions"} {
			if !strings.Contains(svgKaggle, cat) {
				t.Errorf("Render(%s, kaggle): expected category track %q in SVG", school, cat)
			}
		}
	}

	// Case 2: 0 Hugging Face models/datasets/spaces -> should pivot to Engineering Velocity
	data.Impact.HFModels = 0
	data.Impact.HFDatasets = 0
	data.Impact.HFSpaces = 0
	svgImpact, err := Render(domain.SchoolTufte, domain.PanelImpact, data)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	if strings.Contains(svgImpact, "0 / 0 / 0") {
		t.Errorf("expected Recruiter-First degradation to avoid rendering '0 / 0 / 0' when HF is unused")
	}
	if !strings.Contains(svgImpact, "Engineering Velocity") {
		t.Errorf("expected 'Engineering Velocity' fallback section when HF == 0")
	}
}

func TestBilingualPtBRRendering(t *testing.T) {
	manifest := config.DefaultManifest()
	manifest.Subject.Locale = "pt-BR"
	data := DefaultPortfolioData(manifest)

	svg, err := Render(domain.SchoolTufte, domain.PanelLanguages, data)
	if err != nil {
		t.Fatalf("Render pt-BR failed: %v", err)
	}
	if !strings.Contains(svg, "Distribuição de Código") {
		t.Errorf("expected Portuguese title in pt-BR locale, got SVG without it")
	}
}

func TestValidateSVGInvariants(t *testing.T) {
	validSVG := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><text fill="#526071">Data</text></svg>`
	if err := ValidateSVG(validSVG); err != nil {
		t.Fatalf("expected valid SVG to pass ValidateSVG, got %v", err)
	}

	// 1. Malformed XML must fail
	malformedSVG := `<svg xmlns="http://www.w3.org/2000/svg"><text fill="#526071">Unclosed</svg>`
	if err := ValidateSVG(malformedSVG); err == nil {
		t.Errorf("expected malformed XML to fail ValidateSVG")
	}

	// 2. Forbidden low-contrast pale hex codes must fail
	for _, pale := range []string{"#e6edf3", "#e2e8f0", "#f8fafc", "#c9d1d9", "#cbd5e1", "#8b949e", "#94a3b8"} {
		paleSVG := `<svg xmlns="http://www.w3.org/2000/svg"><text fill="` + pale + `">Washed Out</text></svg>`
		if err := ValidateSVG(paleSVG); err == nil {
			t.Errorf("expected pale color %s to fail ValidateSVG", pale)
		}
	}

	// 3. Forbidden School meta-labels must fail
	for _, meta := range []string{"Edward Tufte", "Stephen Few", "Alberto Cairo", "Nigel Holmes"} {
		metaSVG := `<svg xmlns="http://www.w3.org/2000/svg"><text fill="#526071">Inspired by ` + meta + `</text></svg>`
		if err := ValidateSVG(metaSVG); err == nil {
			t.Errorf("expected forbidden School meta-label %q to fail ValidateSVG", meta)
		}
	}
}

