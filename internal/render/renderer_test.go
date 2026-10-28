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

	// Case 1: 0 Gold medals -> should show "4 Tracks" instead of "0 Gold" in KPI
	data.Kaggle.TotalGold = 0
	svgKaggle, err := Render(domain.SchoolFew, domain.PanelKaggle, data)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	if strings.Contains(svgKaggle, ">0 Gold<") {
		t.Errorf("expected Recruiter-First degradation to avoid highlighting '0 Gold' in KPI")
	}
	if !strings.Contains(svgKaggle, "4 Tracks") {
		t.Errorf("expected '4 Tracks' fallback KPI when TotalGold == 0")
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
