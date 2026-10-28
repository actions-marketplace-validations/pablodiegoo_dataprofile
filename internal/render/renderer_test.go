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

func TestRenderStackPanelAllSchoolsAndLeaderHighlighting(t *testing.T) {
	manifest := config.DefaultManifest()
	data := DefaultPortfolioData(manifest)

	canonicalPillars := []string{
		"Data &amp; Analytics",
		"Machine Learning &amp; AI",
		"Visualization &amp; Data Apps",
	}
	canonicalTools := []string{
		"Pandas / Polars", "SQL / DuckDB", "Apache Spark", "dbt / Airflow",
		"Scikit-Learn", "PyTorch", "XGBoost / LightGBM", "Hugging Face",
		"Plotly / Dash", "Streamlit", "Matplotlib / Seaborn", "ggplot2 / Tidyverse",
	}

	for _, school := range domain.AllSchools() {
		svg, err := Render(school, domain.PanelStack, data)
		if err != nil {
			t.Fatalf("Render(%s, stack) failed: %v", school, err)
		}
		if err := ValidateSVG(svg); err != nil {
			t.Errorf("ValidateSVG(%s, stack) failed: %v", school, err)
		}
		for _, pillar := range canonicalPillars {
			if !strings.Contains(svg, pillar) {
				t.Errorf("Render(%s, stack) missing pillar %q", school, pillar)
			}
		}
		for _, tool := range canonicalTools {
			if !strings.Contains(svg, tool) {
				t.Errorf("Render(%s, stack) missing canonical tool %q", school, tool)
			}
		}
	}

	// Verify IsLeader flag controls pre-attentive leader highlighting in Few stack panel:
	// When all items in DataAnalytics have RepoCount == 0 and IsLeader == false, none should get the #0284c7 bar fill.
	zeroData := DefaultPortfolioData(manifest)
	zeroData.Stack.DataAnalytics = []domain.StackItem{
		{Name: "Pandas / Polars", RepoCount: 0, WeightPct: 0, IsLeader: false},
		{Name: "SQL / DuckDB", RepoCount: 0, WeightPct: 0, IsLeader: false},
		{Name: "Apache Spark", RepoCount: 0, WeightPct: 0, IsLeader: false},
		{Name: "dbt / Airflow", RepoCount: 0, WeightPct: 0, IsLeader: false},
	}
	svgFewZero, err := Render(domain.SchoolFew, domain.PanelStack, zeroData)
	if err != nil {
		t.Fatalf("Render(few, stack) failed: %v", err)
	}
	if strings.Contains(svgFewZero, `height="6" rx="1.5" fill="#0284c7"`) {
		t.Errorf("expected zero-count non-leader items not to receive leader accent bar fill #0284c7")
	}
}

