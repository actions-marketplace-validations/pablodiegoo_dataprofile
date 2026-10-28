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
			wantContains: []string{"Distribuição de Código", "MB Removidos", "Linguagem Principal", "Produção vs. EDA", "Lockfile &amp; CI Prontos", "(Gráficos de Bala)", "Cód. Prod", "Suíte Testes"},
			forbid:       []string{"MB Stripped", "Primary Language", "Lockfile &amp; CI Ready", "(Bullet Graphs)", "Prod Code", "Test Suite"},
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
			wantContains: []string{"Desempenho Competitivo no Kaggle", "NÍVEL: MASTER (#412 Pico)", "Total Acumulado", "Finais Top 5%", "(Gráficos de Bala)", "Marcador vertical", "indica o limiar necessário"},
			forbid:       []string{"TIER: MASTER (#412 Peak)", "Top 5% Finishes", "(Bullet Graphs)", "Vertical marker", "indicates the threshold"},
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

	// 3. Verify Recruiter-First Velocity fallback and Zero-Gold Kaggle fallback in pt-BR
	zeroHF := data
	zeroHF.Impact.HFModels = 0
	zeroHF.Impact.HFDatasets = 0
	zeroHF.Impact.HFSpaces = 0
	svgZeroHF, err := Render(domain.SchoolTufte, domain.PanelImpact, zeroHF)
	if err != nil {
		t.Fatalf("Render(tufte, impact) zero-HF pt-BR failed: %v", err)
	}
	for _, w := range []string{"Consistência &amp; Entrega", "Sequência Ativa", "Volume de Código", "Impulso Anual (YoY)", "42 Dias", "+34% Cresc."} {
		if !strings.Contains(svgZeroHF, w) {
			t.Errorf("Render(tufte, impact) zero-HF pt-BR missing %q", w)
		}
	}
	for _, f := range []string{"42 Days", "+34% Growth"} {
		if strings.Contains(svgZeroHF, f) {
			t.Errorf("Render(tufte, impact) zero-HF pt-BR still contains English %q", f)
		}
	}

	zeroGoldPt := data
	zeroGoldPt.Kaggle.TotalGold = 0
	svgZeroGoldPt, err := Render(domain.SchoolFew, domain.PanelKaggle, zeroGoldPt)
	if err != nil {
		t.Fatalf("Render(few, kaggle) zero-gold pt-BR failed: %v", err)
	}
	if !strings.Contains(svgZeroGoldPt, "4 Trilhas") || strings.Contains(svgZeroGoldPt, "4 Tracks") {
		t.Errorf("Render(few, kaggle) zero-gold pt-BR expected '4 Trilhas' instead of '4 Tracks'")
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

func TestDynamicSVGGeometriesAcrossSchools(t *testing.T) {
	manifest := config.DefaultManifest()

	// 1. Few Languages: changing ProductionPct and hygiene counts must dynamically change bullet bar widths
	fewA := DefaultPortfolioData(manifest)
	fewA.Languages.TotalRepos = 20
	fewA.Languages.ProductionPct = 80.0
	fewA.Languages.LockfileRepos = 18
	fewA.Languages.CICDRepos = 16
	fewA.Languages.TestSuiteRepos = 14
	fewA.Languages.ContainerizedRepos = 12

	fewB := DefaultPortfolioData(manifest)
	fewB.Languages.TotalRepos = 20
	fewB.Languages.ProductionPct = 25.0
	fewB.Languages.LockfileRepos = 4
	fewB.Languages.CICDRepos = 2
	fewB.Languages.TestSuiteRepos = 2
	fewB.Languages.ContainerizedRepos = 2

	svgFewA, err := Render(domain.SchoolFew, domain.PanelLanguages, fewA)
	if err != nil {
		t.Fatalf("Render(few, languages) fewA failed: %v", err)
	}
	svgFewB, err := Render(domain.SchoolFew, domain.PanelLanguages, fewB)
	if err != nil {
		t.Fatalf("Render(few, languages) fewB failed: %v", err)
	}
	// 80% of 96 = 77 vs 25% of 96 = 24
	if !strings.Contains(svgFewA, `width="77" height="5" fill="#059669"`) {
		t.Errorf("expected Few Prod Code bullet bar width=77 for 80%%, got SVG:\n%s", svgFewA)
	}
	if !strings.Contains(svgFewB, `width="24" height="5" fill="#059669"`) {
		t.Errorf("expected Few Prod Code bullet bar width=24 for 25%%, got SVG:\n%s", svgFewB)
	}

	// 2. Cairo Languages: 10x5 Waffle Chart (50 unit cells) colored proportionally to TopLanguages
	cairoA := DefaultPortfolioData(manifest)
	cairoA.Languages.Languages = []domain.LanguageStat{
		{Name: "Python", Percentage: 80.0, Color: "#0284c7"},
		{Name: "SQL", Percentage: 10.0, Color: "#d97706"},
		{Name: "Go", Percentage: 4.0, Color: "#00ADD8"},
		{Name: "R", Percentage: 4.0, Color: "#2563eb"},
		{Name: "Shell", Percentage: 2.0, Color: "#059669"},
	}
	cairoB := DefaultPortfolioData(manifest)
	cairoB.Languages.Languages = []domain.LanguageStat{
		{Name: "Python", Percentage: 20.0, Color: "#0284c7"},
		{Name: "SQL", Percentage: 60.0, Color: "#d97706"},
		{Name: "Go", Percentage: 10.0, Color: "#00ADD8"},
		{Name: "R", Percentage: 6.0, Color: "#2563eb"},
		{Name: "Shell", Percentage: 4.0, Color: "#059669"},
	}
	svgCairoA, err := Render(domain.SchoolCairo, domain.PanelLanguages, cairoA)
	if err != nil {
		t.Fatalf("Render(cairo, languages) cairoA failed: %v", err)
	}
	svgCairoB, err := Render(domain.SchoolCairo, domain.PanelLanguages, cairoB)
	if err != nil {
		t.Fatalf("Render(cairo, languages) cairoB failed: %v", err)
	}
	// Verify 50 waffle unit cells (width="15" height="8") and dynamic color counts
	if gotCells := strings.Count(svgCairoA, `width="15" height="8"`); gotCells != 50 {
		t.Errorf("expected 50 unit cells in Cairo 10x5 Waffle Chart, got %d", gotCells)
	}
	if strings.Count(svgCairoA, `fill="#0284c7"`) <= strings.Count(svgCairoB, `fill="#0284c7"`) {
		t.Errorf("expected Cairo waffle chart to have more #0284c7 cells at 80%% Python than at 20%% Python")
	}

	// 3. Holmes Languages: 10 Pictorial Unit Blocks (Isotypes) colored dynamically by language percentages
	svgHolmesA, err := Render(domain.SchoolHolmes, domain.PanelLanguages, cairoA)
	if err != nil {
		t.Fatalf("Render(holmes, languages) cairoA failed: %v", err)
	}
	svgHolmesB, err := Render(domain.SchoolHolmes, domain.PanelLanguages, cairoB)
	if err != nil {
		t.Fatalf("Render(holmes, languages) cairoB failed: %v", err)
	}
	if strings.Count(svgHolmesA, `width="14" height="11"`) != 10 || strings.Count(svgHolmesB, `width="14" height="11"`) != 10 {
		t.Errorf("expected 10 Pictorial Unit Blocks in Holmes languages panel, got A=%d B=%d",
			strings.Count(svgHolmesA, `width="14" height="11"`), strings.Count(svgHolmesB, `width="14" height="11"`))
	}
	if svgHolmesA == svgHolmesB {
		t.Errorf("expected Holmes Pictorial Unit Blocks to change geometry when language percentages change")
	}

	// 4. Holmes Kaggle: trophy (🏆) x/y position on stepped podium computed dynamically from OverallTier
	tiers := []string{
		domain.KaggleTierNovice,
		domain.KaggleTierContributor,
		domain.KaggleTierExpert,
		domain.KaggleTierMaster,
		domain.KaggleTierGrandmaster,
	}
	seenTrophyTags := make(map[string]bool)
	for _, tier := range tiers {
		kd := DefaultPortfolioData(manifest)
		kd.Kaggle.OverallTier = tier
		svgK, err := Render(domain.SchoolHolmes, domain.PanelKaggle, kd)
		if err != nil {
			t.Fatalf("Render(holmes, kaggle) for tier %s failed: %v", tier, err)
		}
		idx := strings.Index(svgK, `text-anchor="middle" font-size="14">🏆</text>`)
		if idx == -1 {
			t.Fatalf("expected trophy element in Holmes kaggle SVG for tier %s", tier)
		}
		start := strings.LastIndex(svgK[:idx], "<text ")
		trophyTag := svgK[start : idx+len(`text-anchor="middle" font-size="14">🏆</text>`)]
		if seenTrophyTags[trophyTag] {
			t.Errorf("duplicate trophy position %q for tier %s; expected distinct podium step position per tier", trophyTag, tier)
		}
		seenTrophyTags[trophyTag] = true
	}
}



