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

	// Case 2: 0 Hugging Face models/datasets/spaces -> should pivot to Engineering Velocity across all 4 Schools
	data.Impact.HFModels = 0
	data.Impact.HFDatasets = 0
	data.Impact.HFSpaces = 0
	for _, school := range domain.AllSchools() {
		svgImpact, err := Render(school, domain.PanelImpact, data)
		if err != nil {
			t.Fatalf("Render(%s, impact) failed: %v", school, err)
		}
		if strings.Contains(svgImpact, "0 / 0 / 0") {
			t.Errorf("Render(%s, impact): expected Recruiter-First degradation to avoid rendering '0 / 0 / 0' when HF is unused", school)
		}
		if strings.Contains(svgImpact, "Hugging Face") {
			t.Errorf("Render(%s, impact): expected Recruiter-First degradation to omit 'Hugging Face' section header when HF == 0", school)
		}
		for _, want := range []string{"Engineering Velocity", "Active Streak", "Net Code Volume", "YoY Momentum"} {
			if !strings.Contains(svgImpact, want) {
				t.Errorf("Render(%s, impact): expected %q in fallback velocity section when HF == 0", school, want)
			}
		}
	}
}

func TestImpactPanelRendersDynamicWeeklyTrend(t *testing.T) {
	manifest := config.DefaultManifest()
	dataA := DefaultPortfolioData(manifest)
	dataA.Impact.WeeklyTrend = []int{5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60}

	dataB := DefaultPortfolioData(manifest)
	// 52-week trend ending at minimum (0) instead of maximum
	trend52 := make([]int, 52)
	for i := range trend52 {
		trend52[i] = 52 - i
	}
	dataB.Impact.WeeklyTrend = trend52

	for _, school := range []domain.School{domain.SchoolTufte, domain.SchoolFew, domain.SchoolCairo} {
		svgA, err := Render(school, domain.PanelImpact, dataA)
		if err != nil {
			t.Fatalf("Render(%s, impact) dataA failed: %v", school, err)
		}
		svgB, err := Render(school, domain.PanelImpact, dataB)
		if err != nil {
			t.Fatalf("Render(%s, impact) dataB failed: %v", school, err)
		}
		if svgA == svgB {
			t.Errorf("Render(%s, impact) produced identical SVG for different WeeklyTrend slices; expected dynamic polyline points", school)
		}
	}
}

func TestBilingualPtBRRendering(t *testing.T) {
	manifest := config.DefaultManifest()
	manifest.Subject.Locale = "pt-BR"
	data := DefaultPortfolioData(manifest)

	// 1. Verify all 20 School x Panel combinations render valid SVGs in pt-BR
	for _, school := range domain.AllSchools() {
		for _, panel := range domain.AllPanels() {
			svg, err := Render(school, panel, data)
			if err != nil {
				t.Fatalf("Render(%s, %s) pt-BR failed: %v", school, panel, err)
			}
			if err := ValidateSVG(svg); err != nil {
				t.Errorf("ValidateSVG(%s, %s) pt-BR failed: %v", school, panel, err)
			}
		}
	}

	// 2. Verify School-specific Portuguese localization across all 5 Panels
	checks := []struct {
		school       domain.School
		panel        domain.Panel
		wantContains []string
		forbid       []string
	}{
		{
			school:       domain.SchoolTufte,
			panel:        domain.PanelLanguages,
			wantContains: []string{"Distribuição de Código", "Proporção Real de Código-Fonte", "Módulos de Produção", "Maturidade de Engenharia", "*Nota:"},
			forbid:       []string{"True Source Code Proportions", "Engineering Hygiene Signals"},
		},
		{
			school:       domain.SchoolFew,
			panel:        domain.PanelLanguages,
			wantContains: []string{"Distribuição de Código", "MB Removidos", "Linguagem Principal", "Produção vs. EDA", "Lockfile &amp; CI Prontos"},
			forbid:       []string{"MB Stripped", "Primary Language", "Lockfile &amp; CI Ready"},
		},
		{
			school:       domain.SchoolCairo,
			panel:        domain.PanelLanguages,
			wantContains: []string{"Lidera com", "Módulos de produção representam", "↳ Saídas .ipynb removidas", "Estrutura de Código Waffle"},
			forbid:       []string{"Leads at", "↳ Stripped .ipynb outputs", "Code Structure Waffle"},
		},
		{
			school:       domain.SchoolHolmes,
			panel:        domain.PanelLanguages,
			wantContains: []string{"Distribuição de Código", "Filtro de Cache Jupyter", "Removidos"},
			forbid:       []string{"Jupyter Bloat Filter", "Stripped "},
		},
		{
			school:       domain.SchoolCairo,
			panel:        domain.PanelStack,
			wantContains: []string{"Ecossistema Completo de Dados &amp; ML Ancorado por", "Cobertura ponta a ponta", "↳ Ecossistema principal verificado"},
			forbid:       []string{"Full-Cycle Data", "End-to-end coverage"},
		},
		{
			school:       domain.SchoolTufte,
			panel:        domain.PanelKaggle,
			wantContains: []string{"Desempenho Competitivo no Kaggle", "Medalhas Totais", "Categoria", "Nível", "Rank Global", "Percentil", "Total Acumulado no Portfólio"},
			forbid:       []string{">Category<", ">Global Rank<", "Total Portfolio Tally"},
		},
		{
			school:       domain.SchoolFew,
			panel:        domain.PanelKaggle,
			wantContains: []string{"Desempenho Competitivo no Kaggle", "Total Acumulado", "Marcador vertical", "indica o limiar necessário"},
			forbid:       []string{"Vertical marker", "indicates the threshold"},
		},
		{
			school:       domain.SchoolCairo,
			panel:        domain.PanelKaggle,
			wantContains: []string{"Nível", "Percentil nas 4 trilhas oficiais", "REFERÊNCIA GLOBAL TOP 1%", "↳ Desempenho consistente"},
			forbid:       []string{"TOP 1% GLOBAL BENCHMARK", "↳ Consistent top-percentile"},
		},
		{
			school:       domain.SchoolHolmes,
			panel:        domain.PanelKaggle,
			wantContains: []string{"Desempenho Competitivo no Kaggle", "Pódio de Progressão", "Medalhas Conquistadas"},
			forbid:       []string{"Tier Climb Podium", "Total Medals Earned"},
		},
		{
			school:       domain.SchoolTufte,
			panel:        domain.PanelImpact,
			wantContains: []string{"Impacto Open-Source &amp; Comunidade de IA", "Horizonte de 52 Semanas", "Contribuições Anuais", "Alcance no GitHub"},
			forbid:       []string{"52-Week Horizon"},
		},
		{
			school:       domain.SchoolFew,
			panel:        domain.PanelImpact,
			wantContains: []string{"Impacto Open-Source &amp; Comunidade de IA", "VELOCIDADE ANUAL", "Contribuições Anuais"},
			forbid:       []string{"YoY VELOCITY"},
		},
		{
			school:       domain.SchoolCairo,
			panel:        domain.PanelImpact,
			wantContains: []string{"Contribuições Anuais", "Alcance Open-Source", "Produção contínua de engenharia", "↳ Pico de velocidade atingido"},
			forbid:       []string{"Growing Open-Source Reach", "↳ Peak velocity achieved"},
		},
		{
			school:       domain.SchoolHolmes,
			panel:        domain.PanelImpact,
			wantContains: []string{"CONTRIBUIÇÕES", "Alcance no GitHub", "Artefatos de IA &amp; Hub"},
			forbid:       []string{"CONTRIBUTIONS"},
		},
		{
			school:       domain.SchoolTufte,
			panel:        domain.PanelSummary,
			wantContains: []string{"1. Proporção Real de Código-Fonte", "2. Ecossistema de Dados, ML &amp; Analytics", "3. Desempenho Competitivo no Kaggle", "4. Impacto Open-Source &amp; Comunidade de IA", "Nível Geral", "Medalhas Totais"},
			forbid:       []string{">Overall Tier<", ">Total Medals<"},
		},
		{
			school:       domain.SchoolCairo,
			panel:        domain.PanelSummary,
			wantContains: []string{"1. Proporção Real de Código-Fonte", "2. Ecossistema de Dados, ML &amp; Analytics", "3. Desempenho Competitivo no Kaggle", "4. Impacto Open-Source &amp; Comunidade de IA", "Nível Geral", "Medalhas Totais", "Síntese executiva do portfólio", "↳ Código-fonte normalizado"},
			forbid:       []string{">Overall Tier<", ">Total Medals<", "Executive portfolio synthesis", "Notebook-deflated source code"},
		},
	}

	for _, tc := range checks {
		svg, err := Render(tc.school, tc.panel, data)
		if err != nil {
			t.Fatalf("Render(%s, %s) failed: %v", tc.school, tc.panel, err)
		}
		for _, w := range tc.wantContains {
			if !strings.Contains(svg, w) {
				t.Errorf("Render(%s, %s) pt-BR missing expected localized string %q", tc.school, tc.panel, w)
			}
		}
		for _, f := range tc.forbid {
			if strings.Contains(svg, f) {
				t.Errorf("Render(%s, %s) pt-BR still contains unlocalized English string %q", tc.school, tc.panel, f)
			}
		}
	}

	// 3. Verify Recruiter-First Velocity fallback in pt-BR
	zeroHF := data
	zeroHF.Impact.HFModels = 0
	zeroHF.Impact.HFDatasets = 0
	zeroHF.Impact.HFSpaces = 0
	svgZeroHF, err := Render(domain.SchoolTufte, domain.PanelImpact, zeroHF)
	if err != nil {
		t.Fatalf("Render(tufte, impact) zero-HF pt-BR failed: %v", err)
	}
	for _, w := range []string{"Consistência &amp; Entrega", "Sequência Ativa", "Volume de Código", "Impulso Anual (YoY)"} {
		if !strings.Contains(svgZeroHF, w) {
			t.Errorf("Render(tufte, impact) zero-HF pt-BR missing %q", w)
		}
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

func TestRenderSummarySchoolGrammars(t *testing.T) {
	manifest := config.DefaultManifest()
	data := DefaultPortfolioData(manifest)

	// 1. Tufte summary.svg: zero outer border/cage, serif + tabular-nums, dot-dash plot, supertable, and sparkline with interquartile band
	tufteSVG, err := Render(domain.SchoolTufte, domain.PanelSummary, data)
	if err != nil {
		t.Fatalf("Render(tufte, summary) failed: %v", err)
	}
	if strings.Contains(tufteSVG, `<rect x="2" y="2"`) || strings.Contains(tufteSVG, `<rect x="3" y="3"`) {
		t.Errorf("Render(tufte, summary): expected zero outer cages/borders, found outer <rect>")
	}
	for _, want := range []string{"tabular-nums", "<polyline", "<circle", "1. True Source Code Proportions", "2. Data, ML &amp; Analytics Ecosystem", "3. Kaggle Competitive Performance", "4. Open-Source &amp; AI Community Impact"} {
		if !strings.Contains(tufteSVG, want) {
			t.Errorf("Render(tufte, summary): missing expected Tufte grammar element %q", want)
		}
	}

	// 2. Few summary.svg: top KPI strip, bullet graph qualitative bands + vertical target marker
	fewSVG, err := Render(domain.SchoolFew, domain.PanelSummary, data)
	if err != nil {
		t.Fatalf("Render(few, summary) failed: %v", err)
	}
	for _, want := range []string{"kpi-num", "kpi-lbl", `fill-opacity="0.14"`, `fill-opacity="0.25"`, `stroke="#0284c7" stroke-width="2"`} {
		if !strings.Contains(fewSVG, want) {
			t.Errorf("Render(few, summary): missing expected Few grammar element %q", want)
		}
	}

	// 3. Cairo summary.svg: top crimson rule (#e11d48), serif headline + lead deck, dashed benchmark line, and ↳ callout
	cairoSVG, err := Render(domain.SchoolCairo, domain.PanelSummary, data)
	if err != nil {
		t.Fatalf("Render(cairo, summary) failed: %v", err)
	}
	for _, want := range []string{`fill="#e11d48"`, `class="serif headline"`, `class="sans lead"`, `stroke-dasharray="3,3"`, "↳"} {
		if !strings.Contains(cairoSVG, want) {
			t.Errorf("Render(cairo, summary): missing expected Cairo grammar element %q", want)
		}
	}

	// 4. Holmes summary.svg: magazine header ribbon, isotype blocks, stepped podium with trophy 🏆, and visual badges
	holmesSVG, err := Render(domain.SchoolHolmes, domain.PanelSummary, data)
	if err != nil {
		t.Fatalf("Render(holmes, summary) failed: %v", err)
	}
	for _, want := range []string{`<path d="M 3,13 Q 3,3 13,3`, "🏆", `class="sans banner"`} {
		if !strings.Contains(holmesSVG, want) {
			t.Errorf("Render(holmes, summary): missing expected Holmes grammar element %q", want)
		}
	}
}

func TestAllFivePanelsHonorSchoolVisualGrammars(t *testing.T) {
	manifest := config.DefaultManifest()
	data := DefaultPortfolioData(manifest)

	for _, panel := range domain.AllPanels() {
		// Tufte: zero outer cages/borders, serif + tabular-nums across all 5 panels
		tufteSVG, err := Render(domain.SchoolTufte, panel, data)
		if err != nil {
			t.Fatalf("Render(tufte, %s) failed: %v", panel, err)
		}
		if strings.Contains(tufteSVG, `<rect x="2" y="2"`) || strings.Contains(tufteSVG, `<rect x="3" y="3"`) {
			t.Errorf("Render(tufte, %s) has outer cage/border <rect>", panel)
		}
		if !strings.Contains(tufteSVG, "tabular-nums") {
			t.Errorf("Render(tufte, %s) missing tabular-nums", panel)
		}

		// Few: top KPI strip + Bullet Graph qualitative bands + vertical target marker across all 5 panels
		fewSVG, err := Render(domain.SchoolFew, panel, data)
		if err != nil {
			t.Fatalf("Render(few, %s) failed: %v", panel, err)
		}
		for _, want := range []string{"kpi-num", "kpi-lbl", `fill-opacity="0.1`, `fill-opacity="0.2`, `stroke="#0284c7" stroke-width="2`} {
			if !strings.Contains(fewSVG, want) {
				t.Errorf("Render(few, %s) missing Few grammar element %q", panel, want)
			}
		}

		// Cairo: top crimson rule (#e11d48), serif headline + lead deck, dashed benchmark line, and ↳ callout across all 5 panels
		cairoSVG, err := Render(domain.SchoolCairo, panel, data)
		if err != nil {
			t.Fatalf("Render(cairo, %s) failed: %v", panel, err)
		}
		for _, want := range []string{`fill="#e11d48"`, `class="serif headline"`, `class="sans lead"`, `stroke-dasharray="3,3"`, "↳"} {
			if !strings.Contains(cairoSVG, want) {
				t.Errorf("Render(cairo, %s) missing Cairo grammar element %q", panel, want)
			}
		}

		// Holmes: magazine header ribbon across all 5 panels
		holmesSVG, err := Render(domain.SchoolHolmes, panel, data)
		if err != nil {
			t.Fatalf("Render(holmes, %s) failed: %v", panel, err)
		}
		for _, want := range []string{`<path d="M 3,13 Q 3,3 13,3`, `class="sans banner"`} {
			if !strings.Contains(holmesSVG, want) {
				t.Errorf("Render(holmes, %s) missing Holmes grammar element %q", panel, want)
			}
		}
	}
}



