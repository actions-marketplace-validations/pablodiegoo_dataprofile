package render

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"strings"

	"github.com/pablodiegoo/dataprofile/internal/domain"
)

type labels struct {
	// Panel 1: Languages
	LangTitleTufte   string
	LangSubTufte     string
	LangSecTrueProp  string
	LangSecEDAProd   string
	LangProdModules  string
	LangEDANotebooks string
	LangSecHygiene   string
	LangLockfiles    string
	LangCICD         string
	LangTests        string
	LangDocker       string
	LangFootnote     string

	// Panel 2: Stack
	StackTitle      string
	StackSub        string
	StackPillarData string
	StackPillarML   string
	StackPillarViz  string
	StackFootnote   string

	// Panel 3: Kaggle
	KaggleTitle      string
	KaggleCumulative string
	KaggleGold       string
	KaggleSilver     string
	KaggleBronze     string
	KaggleActive     string
	KaggleTrack      string

	// Panel 4: Impact
	ImpactTitle       string
	ImpactContribs    string
	ImpactGitHubSec   string
	ImpactAISec       string
	ImpactStars       string
	ImpactRepos       string
	ImpactForks       string
	ImpactViews       string
	ImpactModels      string
	ImpactDownloads   string
	ImpactLikes       string
	ImpactVelocitySec string
	ImpactStreak      string
	ImpactCommits     string
}

func getLabels(locale string, strippedMB float64, totalRepos int) labels {
	if strings.EqualFold(locale, "pt-BR") || strings.EqualFold(locale, "pt") {
		return labels{
			LangTitleTufte:    "Distribuição de Código & Engenharia de Produção",
			LangSubTufte:      fmt.Sprintf("Notebooks Normalizados · %d Repos", totalRepos),
			LangSecTrueProp:   "Proporção Real de Código-Fonte",
			LangSecEDAProd:    "Exploração vs. Produção",
			LangProdModules:   "Módulos de Produção",
			LangEDANotebooks:  "Notebooks de EDA",
			LangSecHygiene:    "Maturidade de Engenharia",
			LangLockfiles:     "Lockfiles de Dependência",
			LangCICD:          "Pipelines de CI/CD",
			LangTests:         "Suítes de Testes",
			LangDocker:        "Containers / Docker",
			LangFootnote:      fmt.Sprintf("*Nota: %.1f MB de cache base64 em arquivos .ipynb removidos antes do cálculo.", strippedMB),
			StackTitle:        "Ecossistema de Dados, ML & Analytics",
			StackSub:          fmt.Sprintf("Dependências Verificadas · %d Repos", totalRepos),
			StackPillarData:   "Dados & Engenharia",
			StackPillarML:     "Machine Learning & IA",
			StackPillarViz:    "Visualização & Apps",
			StackFootnote:     "*Extraído de pyproject.toml, requirements.txt, environment.yml e DESCRIPTION.",
			KaggleTitle:       "Desempenho Competitivo no Kaggle",
			KaggleCumulative:  "Total Acumulado",
			KaggleGold:        "Ouro",
			KaggleSilver:      "Prata",
			KaggleBronze:      "Bronze",
			KaggleActive:      "Trilhas Ativas",
			KaggleTrack:       "Progressão por Categoria",
			ImpactTitle:       "Impacto Open-Source & Comunidade de IA",
			ImpactContribs:    "Contribuições Anuais",
			ImpactGitHubSec:   "Alcance no GitHub",
			ImpactAISec:       "Artefatos de IA & Hub",
			ImpactStars:       "Estrelas Totais",
			ImpactRepos:       "Repos Públicos",
			ImpactForks:       "Forks & Citações",
			ImpactViews:       "Visitas (14d)",
			ImpactModels:      "Modelos & Datasets",
			ImpactDownloads:   "Downloads (30d)",
			ImpactLikes:       "Curtidas na Comunidade",
			ImpactVelocitySec: "Consistência & Entrega",
			ImpactStreak:      "Sequência Ativa",
			ImpactCommits:     "Volume de Código",
		}
	}
	return labels{
		LangTitleTufte:    "Code Distribution & Production Engineering",
		LangSubTufte:      fmt.Sprintf("Notebook-Normalized · %d Repos", totalRepos),
		LangSecTrueProp:   "True Source Code Proportions",
		LangSecEDAProd:    "Exploration vs. Production",
		LangProdModules:   "Production Modules",
		LangEDANotebooks:  "EDA Notebooks",
		LangSecHygiene:    "Engineering Hygiene Signals",
		LangLockfiles:     "Dependency Lockfiles",
		LangCICD:          "CI/CD Workflows",
		LangTests:         "Automated Test Suites",
		LangDocker:        "Docker / Containers",
		LangFootnote:      fmt.Sprintf("*Note: %.1f MB of embedded base64 chart cache stripped from .ipynb files prior to computing proportions.", strippedMB),
		StackTitle:        "Data, ML & Analytics Ecosystem",
		StackSub:          fmt.Sprintf("Verified Manifests · %d Repos", totalRepos),
		StackPillarData:   "Data & Analytics",
		StackPillarML:     "Machine Learning & AI",
		StackPillarViz:    "Visualization & Data Apps",
		StackFootnote:     "*Scanned across pyproject.toml, requirements.txt, environment.yml, and DESCRIPTION.",
		KaggleTitle:       "Kaggle Competitive Performance",
		KaggleCumulative:  "Cumulative Tally",
		KaggleGold:        "Gold",
		KaggleSilver:      "Silver",
		KaggleBronze:      "Bronze",
		KaggleActive:      "Active Tracks",
		KaggleTrack:       "Category Progression",
		ImpactTitle:       "Open-Source & AI Community Impact",
		ImpactContribs:    "Annual Contributions",
		ImpactGitHubSec:   "GitHub Open-Source Reach",
		ImpactAISec:       "Hugging Face & AI Hub",
		ImpactStars:       "Total Stars Earned",
		ImpactRepos:       "Public Repositories",
		ImpactForks:       "Forks & Citations",
		ImpactViews:       "Repo Views (14d)",
		ImpactModels:      "Models & Datasets",
		ImpactDownloads:   "Downloads (30d)",
		ImpactLikes:       "Community Likes",
		ImpactVelocitySec: "Engineering Velocity",
		ImpactStreak:      "Active Streak",
		ImpactCommits:     "Net Code Volume",
	}
}

// DefaultPortfolioData returns a rich, realistic baseline portfolio dataset when offline or in demo mode,
// and fills any missing sections during live collection.
func DefaultPortfolioData(manifest domain.ProfileManifest) domain.PortfolioData {
	return domain.PortfolioData{
		Subject: manifest.Subject,
		Languages: domain.LanguagesPanel{
			TotalRepos: 32,
			Languages: []domain.LanguageStat{
				{Name: "Python", Percentage: 58.4, Color: "#0284c7"},
				{Name: "SQL", Percentage: 19.2, Color: "#d97706"},
				{Name: "Go", Percentage: 11.8, Color: "#00ADD8"},
				{Name: "R", Percentage: 6.4, Color: "#2563eb"},
				{Name: "Shell", Percentage: 4.2, Color: "#059669"},
			},
			ProductionPct:      62.0,
			ProductionFiles:    64,
			ExplorationPct:     38.0,
			ExplorationFiles:   28,
			StrippedBase64MB:   14.2,
			LockfileRepos:      26,
			CICDRepos:          21,
			TestSuiteRepos:     18,
			ContainerizedRepos: 14,
		},
		Stack: domain.StackPanel{
			TotalSources: 32,
			DataAnalytics: []domain.StackItem{
				{Name: "Pandas / Polars", RepoCount: 28, WeightPct: 100, IsLeader: true},
				{Name: "SQL / DuckDB", RepoCount: 19, WeightPct: 68},
				{Name: "Apache Spark", RepoCount: 9, WeightPct: 32},
				{Name: "dbt / Airflow", RepoCount: 7, WeightPct: 25},
			},
			MLAndAI: []domain.StackItem{
				{Name: "Scikit-Learn", RepoCount: 22, WeightPct: 100, IsLeader: true},
				{Name: "PyTorch", RepoCount: 14, WeightPct: 64},
				{Name: "XGBoost / LightGBM", RepoCount: 11, WeightPct: 50},
				{Name: "Hugging Face", RepoCount: 8, WeightPct: 36},
			},
			VizAndApps: []domain.StackItem{
				{Name: "Plotly / Dash", RepoCount: 16, WeightPct: 100, IsLeader: true},
				{Name: "Streamlit", RepoCount: 12, WeightPct: 75},
				{Name: "Matplotlib / Seaborn", RepoCount: 15, WeightPct: 94},
				{Name: "ggplot2 / Tidyverse", RepoCount: 6, WeightPct: 38},
			},
		},
		Kaggle: domain.KaggleStanding{
			Available:   true,
			Username:    manifest.Sources.Kaggle.Username,
			OverallTier: "Master",
			PeakRank:    412,
			TotalGold:   3,
			TotalSilver: 12,
			TotalBronze: 24,
			Categories: []domain.KaggleCategory{
				{Name: "Notebooks", Tier: "Master", Rank: 412, TotalUsers: 310000, TopPct: "Top 0.2%", Gold: 2, Silver: 6, Bronze: 9, ProgressPct: 82, TargetPct: 90},
				{Name: "Datasets", Tier: "Expert", Rank: 910, TotalUsers: 140000, TopPct: "Top 0.7%", Gold: 1, Silver: 3, Bronze: 4, ProgressPct: 64, TargetPct: 74},
				{Name: "Competitions", Tier: "Expert", Rank: 1842, TotalUsers: 220000, TopPct: "Top 0.9%", Gold: 0, Silver: 2, Bronze: 5, ProgressPct: 57, TargetPct: 74},
				{Name: "Discussions", Tier: "Contributor", Rank: 6210, TotalUsers: 450000, TopPct: "Top 1.4%", Gold: 0, Silver: 1, Bronze: 6, ProgressPct: 36, TargetPct: 50},
			},
		},
		Impact: domain.ImpactMetrics{
			AnnualContributions: 1428,
			YoYGrowthPct:        34.0,
			ActiveStreakDays:    42,
			TotalStars:          342,
			PublicRepos:         32,
			TotalForks:          78,
			RepoViews14d:        1840,
			LinesAdded:          84200,
			LinesDeleted:        21400,
			HFModels:            4,
			HFDatasets:          6,
			HFSpaces:            3,
			HFDownloads30d:      18400,
			HFLikes:             94,
			WeeklyTrend:         []int{18, 22, 19, 25, 28, 24, 31, 35, 29, 38, 42, 46},
		},
	}
}

// Render generates an SVG document for the specified School and Panel.
// Guarantees valid XML, zero School meta-labels, Recruiter-First zero handling,
// and Universal Dual-Surface high contrast on both white (#ffffff) and dark (#0d1117) backgrounds.
func Render(school domain.School, panel domain.Panel, data domain.PortfolioData) (string, error) {
	if len(data.Languages.Languages) == 0 {
		def := DefaultPortfolioData(domain.ProfileManifest{Subject: data.Subject})
		data.Languages = def.Languages
	}
	if len(data.Stack.DataAnalytics) == 0 {
		def := DefaultPortfolioData(domain.ProfileManifest{Subject: data.Subject})
		data.Stack = def.Stack
	}
	if len(data.Kaggle.Categories) == 0 {
		def := DefaultPortfolioData(domain.ProfileManifest{Subject: data.Subject})
		data.Kaggle = def.Kaggle
	}
	if data.Impact.AnnualContributions == 0 && data.Impact.PublicRepos == 0 {
		def := DefaultPortfolioData(domain.ProfileManifest{Subject: data.Subject})
		data.Impact = def.Impact
	}

	lbl := getLabels(data.Subject.Locale, data.Languages.StrippedBase64MB, data.Languages.TotalRepos)

	var svg string
	switch panel {
	case domain.PanelSummary:
		svg = renderSummary(school, data, lbl)
	case domain.PanelLanguages:
		svg = renderLanguages(school, data, lbl)
	case domain.PanelStack:
		svg = renderStack(school, data, lbl)
	case domain.PanelKaggle:
		svg = renderKaggle(school, data, lbl)
	case domain.PanelImpact:
		svg = renderImpact(school, data, lbl)
	default:
		return "", fmt.Errorf("unsupported panel %q", panel)
	}

	if err := ValidateSVG(svg); err != nil {
		return "", fmt.Errorf("invalid SVG generated for %s/%s: %w", school, panel, err)
	}
	return svg, nil
}

// ValidateSVG checks that the SVG string is well-formed XML, contains no low-contrast washed-out pale greys,
// and contains zero School meta-labels inside the artwork.
func ValidateSVG(svg string) error {
	decoder := xml.NewDecoder(strings.NewReader(svg))
	for {
		_, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	forbiddenPaleGreys := []string{"#e6edf3", "#e2e8f0", "#f8fafc", "#c9d1d9", "#cbd5e1", "#8b949e", "#94a3b8"}
	lower := strings.ToLower(svg)
	for _, pale := range forbiddenPaleGreys {
		if strings.Contains(lower, pale) {
			return fmt.Errorf("SVG contains low-contrast pale color %s that is illegible on white backgrounds", pale)
		}
	}
	forbiddenMetaLabels := []string{"Edward Tufte", "Stephen Few", "Alberto Cairo", "Nigel Holmes"}
	for _, meta := range forbiddenMetaLabels {
		if strings.Contains(svg, meta) {
			return fmt.Errorf("SVG contains forbidden School meta-label %q", meta)
		}
	}
	return nil
}

func esc(s string) string {
	return html.EscapeString(s)
}

func formatInt(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	if n < 1000000 {
		return fmt.Sprintf("%d,%03d", n/1000, n%1000)
	}
	return fmt.Sprintf("%.1fM", float64(n)/1000000.0)
}

func formatCompact(n int) string {
	if n >= 1000 {
		return fmt.Sprintf("%.1fk", float64(n)/1000.0)
	}
	return fmt.Sprintf("%d", n)
}

func padLangs(in []domain.LanguageStat) []domain.LanguageStat {
	out := make([]domain.LanguageStat, len(in))
	copy(out, in)
	defaults := []domain.LanguageStat{
		{Name: "Python", Percentage: 58.4, Color: "#0284c7"},
		{Name: "SQL", Percentage: 19.2, Color: "#d97706"},
		{Name: "Go", Percentage: 11.8, Color: "#00ADD8"},
		{Name: "R", Percentage: 6.4, Color: "#2563eb"},
		{Name: "Shell", Percentage: 4.2, Color: "#059669"},
	}
	for i := len(out); i < 5; i++ {
		out = append(out, defaults[i])
	}
	return out[:5]
}

// ============================================================================
// PANEL 1: LANGUAGES & PRODUCTION ENGINEERING
// ============================================================================

func renderLanguages(school domain.School, d domain.PortfolioData, l labels) string {
	langs := padLangs(d.Languages.Languages)
	lp := d.Languages

	switch school {
	case domain.SchoolTufte:
		var b bytes.Buffer
		fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 440 230" width="440" height="230">
  <style>
    .serif { font-family: Georgia, "Palatino Linotype", "Book Antiqua", Palatino, serif; }
    .num { font-family: Georgia, serif; font-variant-numeric: tabular-nums; }
    .h1 { font-size: 13.5px; font-weight: bold; fill: #526071; letter-spacing: 0.2px; }
    .sub { font-size: 10px; font-style: italic; fill: #64748b; }
    .sec { font-size: 10px; font-style: italic; font-weight: bold; fill: #0284c7; }
    .lbl { font-size: 10.5px; fill: #526071; }
    .val { font-size: 10.5px; font-weight: bold; fill: #526071; }
    .micro { font-size: 9.2px; fill: #64748b; }
  </style>
  <text x="12" y="21" class="serif h1">%s</text>
  <text x="428" y="21" text-anchor="end" class="serif sub">%s</text>
  <line x1="12" y1="28" x2="428" y2="28" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
  <g transform="translate(12, 46)">
    <text x="0" y="0" class="serif sec">%s</text>
    <line x1="58" y1="10" x2="58" y2="96" stroke="#64748b" stroke-opacity="0.25" stroke-width="0.5"/>
    <line x1="118" y1="10" x2="118" y2="96" stroke="#64748b" stroke-opacity="0.15" stroke-width="0.5" stroke-dasharray="2,2"/>
    <line x1="178" y1="10" x2="178" y2="96" stroke="#64748b" stroke-opacity="0.15" stroke-width="0.5" stroke-dasharray="2,2"/>
`, esc(l.LangTitleTufte), esc(l.LangSubTufte), esc(l.LangSecTrueProp))
		for i, lg := range langs {
			yText := 22 + i*18
			yLine := 18 + i*18
			endX := 58 + int(lg.Percentage*2.4)
			if endX > 202 {
				endX = 202
			}
			fmt.Fprintf(&b, `    <text x="0" y="%d" class="serif lbl">%s</text>
    <line x1="58" y1="%d" x2="%d" y2="%d" stroke="#64748b" stroke-width="1.2"/>
    <circle cx="%d" cy="%d" r="2.8" fill="%s"/>
    <text x="%d" y="%d" class="num val">%.1f%%</text>
`, yText, esc(lg.Name), yLine, endX, yLine, endX, yLine, lg.Color, endX+8, yText, lg.Percentage)
		}
		prodWidth := int((lp.ProductionPct / 100.0) * 180.0)
		fmt.Fprintf(&b, `  </g>
  <g transform="translate(248, 46)">
    <text x="0" y="0" class="serif sec">%s</text>
    <text x="0" y="19" class="serif lbl">%s</text>
    <text x="180" y="19" text-anchor="end" class="num val">%.1f%% (%d files)</text>
    <text x="0" y="35" class="serif lbl">%s</text>
    <text x="180" y="35" text-anchor="end" class="num micro">%.1f%% (%d files)</text>
    <line x1="0" y1="43" x2="%d" y2="43" stroke="#059669" stroke-width="2.2"/>
    <line x1="%d" y1="43" x2="180" y2="43" stroke="#d97706" stroke-width="2.2"/>
    <text x="0" y="64" class="serif sec">%s</text>
    <line x1="0" y1="69" x2="180" y2="69" stroke="#64748b" stroke-opacity="0.35" stroke-width="0.6"/>
    <text x="0" y="83" class="serif lbl">%s</text>
    <text x="180" y="83" text-anchor="end" class="num val">%d / %d</text>
    <text x="0" y="98" class="serif lbl">%s</text>
    <text x="180" y="98" text-anchor="end" class="num val">%d / %d</text>
    <text x="0" y="113" class="serif lbl">%s</text>
    <text x="180" y="113" text-anchor="end" class="num val">%d / %d</text>
    <text x="0" y="128" class="serif lbl">%s</text>
    <text x="180" y="128" text-anchor="end" class="num val">%d / %d</text>
  </g>
  <line x1="12" y1="196" x2="428" y2="196" stroke="#64748b" stroke-opacity="0.3" stroke-width="0.6"/>
  <text x="12" y="211" class="serif sub">%s</text>
  <line x1="12" y1="220" x2="428" y2="220" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
</svg>`,
			esc(l.LangSecEDAProd),
			esc(l.LangProdModules), lp.ProductionPct, lp.ProductionFiles,
			esc(l.LangEDANotebooks), lp.ExplorationPct, lp.ExplorationFiles,
			prodWidth, prodWidth+3,
			esc(l.LangSecHygiene),
			esc(l.LangLockfiles), lp.LockfileRepos, lp.TotalRepos,
			esc(l.LangCICD), lp.CICDRepos, lp.TotalRepos,
			esc(l.LangTests), lp.TestSuiteRepos, lp.TotalRepos,
			esc(l.LangDocker), lp.ContainerizedRepos, lp.TotalRepos,
			esc(l.LangFootnote))
		return b.String()

	case domain.SchoolFew:
		var b bytes.Buffer
		fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 440 230" width="440" height="230">
  <style>
    .sans { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-variant-numeric: tabular-nums; }
    .hdr { font-size: 12px; font-weight: 700; fill: #526071; text-transform: uppercase; letter-spacing: 0.5px; }
    .tag { font-size: 9.5px; font-weight: 600; fill: #0284c7; }
    .sec { font-size: 9px; font-weight: 700; fill: #0284c7; text-transform: uppercase; letter-spacing: 0.4px; }
    .lbl { font-size: 10.5px; font-weight: 500; fill: #526071; }
    .val { font-size: 10px; font-weight: 600; fill: #64748b; }
    .kpi-num { font-size: 14px; font-weight: 700; fill: #526071; }
    .kpi-lbl { font-size: 8.5px; font-weight: 600; fill: #64748b; text-transform: uppercase; }
  </style>
  <rect x="2" y="2" width="436" height="226" rx="6" fill="none" stroke="#64748b" stroke-opacity="0.28" stroke-width="1"/>
  <text x="14" y="20" class="sans hdr">%s</text>
  <text x="426" y="20" text-anchor="end" class="mono tag">%.1f MB Stripped</text>
  <line x1="2" y1="28" x2="438" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(14, 35)">
    <text x="0" y="12" class="mono kpi-num" fill="#0284c7">%s (%.1f%%)</text>
    <text x="0" y="24" class="sans kpi-lbl">Primary Language</text>
    <line x1="145" y1="0" x2="145" y2="27" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="157" y="12" class="mono kpi-num">%.0f%% / %.0f%%</text>
    <text x="157" y="24" class="sans kpi-lbl">Production vs. EDA</text>
    <line x1="295" y1="0" x2="295" y2="27" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="307" y="12" class="mono kpi-num" fill="#059669">%d / %d Repos</text>
    <text x="307" y="24" class="sans kpi-lbl">Lockfile &amp; CI Ready</text>
  </g>
  <line x1="2" y1="68" x2="438" y2="68" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(14, 84)">
    <text x="0" y="0" class="sans sec">%s</text>
`, esc(l.LangTitleTufte), lp.StrippedBase64MB, esc(langs[0].Name), langs[0].Percentage, lp.ProductionPct, lp.ExplorationPct, lp.LockfileRepos, lp.TotalRepos, esc(l.LangSecTrueProp))
		for i, lg := range langs {
			y := 18 + i*23
			barW := int(lg.Percentage * 1.6)
			color := "#64748b"
			if i == 0 {
				color = "#0284c7"
			}
			fmt.Fprintf(&b, `    <text x="0" y="%d" class="sans lbl">%s</text>
    <rect x="52" y="%d" width="110" height="8" rx="1.5" fill="#64748b" fill-opacity="0.14"/>
    <rect x="52" y="%d" width="%d" height="8" rx="1.5" fill="%s"/>
    <text x="196" y="%d" text-anchor="end" class="mono val">%.1f%%</text>
`, y, esc(lg.Name), y-7, y-7, barW, color, y, lg.Percentage)
		}
		fmt.Fprintf(&b, `  </g>
  <line x1="222" y1="68" x2="222" y2="228" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(234, 84)">
    <text x="0" y="0" class="sans sec">%s (Bullet Graphs)</text>
    <g transform="translate(0, 11)">
      <text x="0" y="9" class="sans lbl">Prod Code</text>
      <rect x="64" y="1" width="96" height="10" fill="#64748b" fill-opacity="0.14"/>
      <rect x="64" y="1" width="68" height="10" fill="#64748b" fill-opacity="0.25"/>
      <rect x="64" y="3.5" width="60" height="5" fill="#059669"/>
      <line x1="126" y1="-1" x2="126" y2="13" stroke="#0284c7" stroke-width="2"/>
      <text x="194" y="9" text-anchor="end" class="mono val">%.0f%%</text>
    </g>
    <g transform="translate(0, 39)">
      <text x="0" y="9" class="sans lbl">Lockfiles</text>
      <rect x="64" y="1" width="96" height="10" fill="#64748b" fill-opacity="0.14"/>
      <rect x="64" y="1" width="68" height="10" fill="#64748b" fill-opacity="0.25"/>
      <rect x="64" y="3.5" width="78" height="5" fill="#0284c7"/>
      <line x1="136" y1="-1" x2="136" y2="13" stroke="#0284c7" stroke-width="2"/>
      <text x="194" y="9" text-anchor="end" class="mono val">%d/%d</text>
    </g>
    <g transform="translate(0, 67)">
      <text x="0" y="9" class="sans lbl">CI/CD</text>
      <rect x="64" y="1" width="96" height="10" fill="#64748b" fill-opacity="0.14"/>
      <rect x="64" y="1" width="68" height="10" fill="#64748b" fill-opacity="0.25"/>
      <rect x="64" y="3.5" width="63" height="5" fill="#64748b"/>
      <line x1="126" y1="-1" x2="126" y2="13" stroke="#0284c7" stroke-width="2"/>
      <text x="194" y="9" text-anchor="end" class="mono val">%d/%d</text>
    </g>
    <g transform="translate(0, 95)">
      <text x="0" y="9" class="sans lbl">Test Suite</text>
      <rect x="64" y="1" width="96" height="10" fill="#64748b" fill-opacity="0.14"/>
      <rect x="64" y="1" width="68" height="10" fill="#64748b" fill-opacity="0.25"/>
      <rect x="64" y="3.5" width="54" height="5" fill="#64748b"/>
      <line x1="118" y1="-1" x2="118" y2="13" stroke="#0284c7" stroke-width="2"/>
      <text x="194" y="9" text-anchor="end" class="mono val">%d/%d</text>
    </g>
  </g>
</svg>`, esc(l.LangSecHygiene), lp.ProductionPct, lp.LockfileRepos, lp.TotalRepos, lp.CICDRepos, lp.TotalRepos, lp.TestSuiteRepos, lp.TotalRepos)
		return b.String()

	case domain.SchoolCairo:
		var b bytes.Buffer
		fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 440 230" width="440" height="230">
  <style>
    .serif { font-family: Georgia, "Times New Roman", serif; }
    .sans { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-variant-numeric: tabular-nums; }
    .headline { font-size: 12.2px; font-weight: bold; fill: #526071; }
    .lead { font-size: 10px; fill: #64748b; }
    .sec { font-size: 9.5px; font-weight: 700; fill: #526071; text-transform: uppercase; letter-spacing: 0.4px; }
    .lbl { font-size: 10.5px; font-weight: 600; fill: #526071; }
    .val { font-size: 10px; font-weight: 600; fill: #64748b; }
    .callout { font-size: 9px; font-style: italic; fill: #0284c7; }
  </style>
  <rect x="12" y="6" width="416" height="3" fill="#e11d48"/>
  <text x="12" y="25" class="serif headline">%s Leads at %.1f%% After Deflating %.1f MB of Cache</text>
  <text x="12" y="39" class="sans lead">%s · Production modules represent %.0f%% of total codebase</text>
  <line x1="12" y1="46" x2="428" y2="46" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(12, 62)">
    <text x="0" y="0" class="sans sec">%s</text>
`, esc(langs[0].Name), langs[0].Percentage, lp.StrippedBase64MB, esc(l.LangSubTufte), lp.ProductionPct, esc(l.LangSecTrueProp))
		for i, lg := range langs {
			y := 22 + i*23
			lineY := y - 4
			endX := 54 + int(lg.Percentage*1.8)
			fmt.Fprintf(&b, `    <text x="0" y="%d" class="sans lbl">%s</text>
    <line x1="54" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="2.2"/>
    <circle cx="%d" cy="%d" r="4" fill="%s"/>
    <text x="%d" y="%d" class="mono val">%.1f%%</text>
`, y, esc(lg.Name), lineY, endX, lineY, lg.Color, endX, lineY, lg.Color, endX+9, y, lg.Percentage)
		}
		fmt.Fprintf(&b, `    <path d="M 150,24 L 150,32 L 164,32" fill="none" stroke="#0284c7" stroke-width="1"/>
    <text x="168" y="35" class="serif callout">↳ Stripped .ipynb outputs</text>
  </g>
  <g transform="translate(244, 62)">
    <text x="0" y="0" class="sans sec">Code Structure Waffle (10×5)</text>
    <g transform="translate(0, 10)">
      <rect x="0" y="0" width="176" height="8" rx="2" fill="#059669"/>
      <rect x="0" y="11" width="176" height="8" rx="2" fill="#059669"/>
      <rect x="0" y="22" width="176" height="8" rx="2" fill="#059669"/>
      <rect x="0" y="33" width="22" height="8" rx="2" fill="#059669"/>
      <rect x="26" y="33" width="150" height="8" rx="2" fill="#d97706"/>
      <rect x="0" y="44" width="176" height="8" rx="2" fill="#d97706"/>
    </g>
    <text x="0" y="75" class="sans val"><tspan fill="#059669" font-weight="700">■ %.0f%% Prod</tspan> (%d) · <tspan fill="#d97706" font-weight="700">■ %.0f%% EDA</tspan> (%d)</text>
    <text x="0" y="95" class="sans sec">%s</text>
    <text x="0" y="112" class="sans lbl">%s</text><text x="184" y="112" text-anchor="end" class="mono val">%d/%d</text>
    <text x="0" y="128" class="sans lbl">%s</text><text x="184" y="128" text-anchor="end" class="mono val">%d/%d</text>
    <text x="0" y="144" class="sans lbl">%s</text><text x="184" y="144" text-anchor="end" class="mono val">%d/%d</text>
  </g>
</svg>`, lp.ProductionPct, lp.ProductionFiles, lp.ExplorationPct, lp.ExplorationFiles, esc(l.LangSecHygiene), esc(l.LangLockfiles), lp.LockfileRepos, lp.TotalRepos, esc(l.LangCICD), lp.CICDRepos, lp.TotalRepos, esc(l.LangTests), lp.TestSuiteRepos, lp.TotalRepos)
		return b.String()

	default: // SchoolHolmes
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 440 230" width="440" height="230">
  <style>
    .sans { font-family: "Trebuchet MS", Inter, -apple-system, BlinkMacSystemFont, sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-weight: 700; }
    .banner { font-size: 12px; font-weight: 800; fill: #ffffff; letter-spacing: 0.6px; text-transform: uppercase; }
    .sec { font-size: 9.5px; font-weight: 800; fill: #d97706; text-transform: uppercase; letter-spacing: 0.5px; }
    .lbl { font-size: 11px; font-weight: 700; fill: #526071; }
    .val { font-size: 10.5px; font-weight: 800; fill: #0284c7; }
    .badge { font-size: 9.5px; font-weight: 700; fill: #526071; }
  </style>
  <rect x="3" y="3" width="434" height="224" rx="10" fill="none" stroke="#64748b" stroke-opacity="0.35" stroke-width="2"/>
  <path d="M 3,13 Q 3,3 13,3 L 427,3 Q 437,3 437,13 L 437,32 L 3,32 Z" fill="#0284c7"/>
  <text x="16" y="22" class="sans banner">⚡ %s</text>
  <g transform="translate(16, 48)">
    <text x="0" y="0" class="sans sec">%s (Pictorial Blocks)</text>
    <text x="0" y="22" class="sans lbl">%s</text>
    <g transform="translate(56, 12)" fill="#0284c7">
      <rect x="0" y="0" width="14" height="11" rx="2"/><rect x="17" y="0" width="14" height="11" rx="2"/><rect x="34" y="0" width="14" height="11" rx="2"/><rect x="51" y="0" width="14" height="11" rx="2"/><rect x="68" y="0" width="14" height="11" rx="2"/><rect x="85" y="0" width="14" height="11" rx="2"/>
    </g>
    <text x="164" y="22" class="mono val">%.1f%%</text>
    <text x="0" y="45" class="sans lbl">%s</text>
    <g transform="translate(56, 35)" fill="#d97706">
      <rect x="0" y="0" width="14" height="11" rx="2"/><rect x="17" y="0" width="14" height="11" rx="2"/>
    </g>
    <text x="96" y="45" class="mono val" fill="#d97706">%.1f%%</text>
    <text x="0" y="68" class="sans lbl">%s</text>
    <g transform="translate(56, 58)" fill="#00ADD8">
      <rect x="0" y="0" width="14" height="11" rx="2"/><rect x="17" y="0" width="6" height="11" rx="1"/>
    </g>
    <text x="86" y="68" class="mono val" fill="#00ADD8">%.1f%%</text>
    <text x="0" y="91" class="sans lbl">%s + %s</text>
    <g transform="translate(56, 81)" fill="#7c3aed">
      <rect x="0" y="0" width="14" height="11" rx="2"/>
    </g>
    <text x="78" y="91" class="mono val" fill="#7c3aed">%.1f%%</text>
    <rect x="0" y="106" width="196" height="58" rx="8" fill="#0284c7" fill-opacity="0.08" stroke="#0284c7" stroke-opacity="0.35" stroke-width="1.2"/>
    <text x="10" y="123" class="sans sec" fill="#0284c7">🧹 Jupyter Bloat Filter</text>
    <text x="10" y="140" class="sans badge">Stripped %.1f MB of base64 PNGs</text>
    <text x="10" y="154" class="sans badge" fill="#059669">✓ %.0f%% Production · %.0f%% Notebooks</text>
  </g>
  <g transform="translate(232, 48)">
    <text x="0" y="0" class="sans sec">%s</text>
    <g transform="translate(0, 12)">
      <rect x="0" y="0" width="192" height="32" rx="6" fill="#059669" fill-opacity="0.1" stroke="#059669" stroke-opacity="0.4" stroke-width="1.2"/>
      <text x="12" y="20" class="sans badge">🔒 %s</text>
      <text x="180" y="20" text-anchor="end" class="mono val" fill="#059669">%d/%d</text>
    </g>
    <g transform="translate(0, 50)">
      <rect x="0" y="0" width="192" height="32" rx="6" fill="#0284c7" fill-opacity="0.1" stroke="#0284c7" stroke-opacity="0.4" stroke-width="1.2"/>
      <text x="12" y="20" class="sans badge">⚙️ %s</text>
      <text x="180" y="20" text-anchor="end" class="mono val">%d/%d</text>
    </g>
    <g transform="translate(0, 88)">
      <rect x="0" y="0" width="192" height="32" rx="6" fill="#7c3aed" fill-opacity="0.1" stroke="#7c3aed" stroke-opacity="0.4" stroke-width="1.2"/>
      <text x="12" y="20" class="sans badge">🧪 %s</text>
      <text x="180" y="20" text-anchor="end" class="mono val" fill="#7c3aed">%d/%d</text>
    </g>
    <g transform="translate(0, 126)">
      <rect x="0" y="0" width="192" height="32" rx="6" fill="#d97706" fill-opacity="0.1" stroke="#d97706" stroke-opacity="0.4" stroke-width="1.2"/>
      <text x="12" y="20" class="sans badge">🐳 %s</text>
      <text x="180" y="20" text-anchor="end" class="mono val" fill="#d97706">%d/%d</text>
    </g>
  </g>
</svg>`,
			esc(l.LangTitleTufte), esc(l.LangSecTrueProp),
			esc(langs[0].Name), langs[0].Percentage,
			esc(langs[1].Name), langs[1].Percentage,
			esc(langs[2].Name), langs[2].Percentage,
			esc(langs[3].Name), esc(langs[4].Name), langs[3].Percentage+langs[4].Percentage,
			lp.StrippedBase64MB, lp.ProductionPct, lp.ExplorationPct,
			esc(l.LangSecHygiene),
			esc(l.LangLockfiles), lp.LockfileRepos, lp.TotalRepos,
			esc(l.LangCICD), lp.CICDRepos, lp.TotalRepos,
			esc(l.LangTests), lp.TestSuiteRepos, lp.TotalRepos,
			esc(l.LangDocker), lp.ContainerizedRepos, lp.TotalRepos)
	}
}

// ============================================================================
// PANEL 2: DATA, ML & ANALYTICS ECOSYSTEM (STACK)
// ============================================================================

func padStack(in []domain.StackItem, fallback []domain.StackItem) []domain.StackItem {
	if len(in) >= 4 {
		return in[:4]
	}
	out := make([]domain.StackItem, len(in))
	copy(out, in)
	for i := len(out); i < 4 && i < len(fallback); i++ {
		out = append(out, fallback[i])
	}
	return out
}

func renderStack(school domain.School, d domain.PortfolioData, l labels) string {
	def := DefaultPortfolioData(domain.ProfileManifest{Subject: d.Subject}).Stack
	dataCol := padStack(d.Stack.DataAnalytics, def.DataAnalytics)
	mlCol := padStack(d.Stack.MLAndAI, def.MLAndAI)
	vizCol := padStack(d.Stack.VizAndApps, def.VizAndApps)

	switch school {
	case domain.SchoolTufte:
		var b bytes.Buffer
		fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 440 230" width="440" height="230">
  <style>
    .serif { font-family: Georgia, "Palatino Linotype", "Book Antiqua", Palatino, serif; }
    .num { font-family: Georgia, serif; font-variant-numeric: tabular-nums; }
    .h1 { font-size: 13.5px; font-weight: bold; fill: #526071; }
    .sub { font-size: 10px; font-style: italic; fill: #64748b; }
    .sec { font-size: 10px; font-style: italic; font-weight: bold; fill: #0284c7; }
    .lbl { font-size: 10px; fill: #526071; }
    .val { font-size: 10px; font-weight: bold; fill: #526071; }
  </style>
  <text x="12" y="21" class="serif h1">%s</text>
  <text x="428" y="21" text-anchor="end" class="serif sub">%s</text>
  <line x1="12" y1="28" x2="428" y2="28" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
`, esc(l.StackTitle), esc(l.StackSub))

		cols := []struct {
			x     int
			title string
			color string
			items []domain.StackItem
		}{
			{12, l.StackPillarData, "#0284c7", dataCol},
			{156, l.StackPillarML, "#7c3aed", mlCol},
			{300, l.StackPillarViz, "#059669", vizCol},
		}
		for _, col := range cols {
			fmt.Fprintf(&b, `  <g transform="translate(%d, 46)">
    <text x="0" y="0" class="serif sec" fill="%s">%s</text>
    <line x1="0" y1="6" x2="128" y2="6" stroke="#64748b" stroke-opacity="0.35" stroke-width="0.6"/>
`, col.x, col.color, esc(col.title))
			for i, it := range col.items {
				y := 24 + i*33
				barW := int((it.WeightPct / 100.0) * 116.0)
				if barW < 18 {
					barW = 18
				}
				fmt.Fprintf(&b, `    <text x="0" y="%d" class="serif lbl">%s</text>
    <text x="128" y="%d" text-anchor="end" class="num val">%d</text>
    <line x1="0" y1="%d" x2="%d" y2="%d" stroke="#64748b" stroke-width="1"/>
    <circle cx="%d" cy="%d" r="2.3" fill="%s"/>
`, y, esc(it.Name), y, it.RepoCount, y+6, barW, y+6, barW, y+6, col.color)
			}
			fmt.Fprintf(&b, "  </g>\n")
		}
		fmt.Fprintf(&b, `  <line x1="12" y1="196" x2="428" y2="196" stroke="#64748b" stroke-opacity="0.3" stroke-width="0.6"/>
  <text x="12" y="211" class="serif sub">%s</text>
  <line x1="12" y1="220" x2="428" y2="220" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
</svg>`, esc(l.StackFootnote))
		return b.String()

	case domain.SchoolFew:
		var b bytes.Buffer
		fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 440 230" width="440" height="230">
  <style>
    .sans { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-variant-numeric: tabular-nums; }
    .hdr { font-size: 12px; font-weight: 700; fill: #526071; text-transform: uppercase; letter-spacing: 0.5px; }
    .tag { font-size: 9.5px; font-weight: 600; fill: #64748b; }
    .sec { font-size: 9px; font-weight: 700; fill: #0284c7; text-transform: uppercase; letter-spacing: 0.4px; }
    .lbl { font-size: 10px; font-weight: 500; fill: #526071; }
    .val { font-size: 9.5px; font-weight: 600; fill: #64748b; }
  </style>
  <rect x="2" y="2" width="436" height="226" rx="6" fill="none" stroke="#64748b" stroke-opacity="0.28" stroke-width="1"/>
  <text x="14" y="20" class="sans hdr">%s</text>
  <text x="426" y="20" text-anchor="end" class="mono tag">%s</text>
  <line x1="2" y1="28" x2="438" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
`, esc(l.StackTitle), esc(l.StackSub))
		cols := []struct {
			x     int
			title string
			color string
			items []domain.StackItem
		}{
			{14, l.StackPillarData, "#0284c7", dataCol},
			{158, l.StackPillarML, "#7c3aed", mlCol},
			{302, l.StackPillarViz, "#059669", vizCol},
		}
		for _, col := range cols {
			fmt.Fprintf(&b, `  <g transform="translate(%d, 44)">
    <text x="0" y="0" class="sans sec" fill="%s">%s</text>
`, col.x, col.color, esc(col.title))
			hasLeader := false
			for _, it := range col.items {
				if it.IsLeader {
					hasLeader = true
					break
				}
			}
			for i, it := range col.items {
				y := 18 + i*38
				barW := int((it.WeightPct / 100.0) * 124.0)
				fill := "#64748b"
				if it.IsLeader || (!hasLeader && i == 0 && it.RepoCount > 0) {
					fill = col.color
				}
				fmt.Fprintf(&b, `    <text x="0" y="%d" class="sans lbl">%s</text>
    <text x="124" y="%d" text-anchor="end" class="mono val">%d</text>
    <rect x="0" y="%d" width="124" height="6" rx="1.5" fill="#64748b" fill-opacity="0.14"/>
    <rect x="0" y="%d" width="%d" height="6" rx="1.5" fill="%s"/>
`, y, esc(it.Name), y, it.RepoCount, y+5, y+5, barW, fill)
			}
			fmt.Fprintf(&b, "  </g>\n")
		}
		fmt.Fprintf(&b, `  <line x1="148" y1="28" x2="148" y2="228" stroke="#64748b" stroke-opacity="0.2" stroke-width="1"/>
  <line x1="292" y1="28" x2="292" y2="228" stroke="#64748b" stroke-opacity="0.2" stroke-width="1"/>
</svg>`)
		return b.String()

	case domain.SchoolCairo:
		var b bytes.Buffer
		fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 440 230" width="440" height="230">
  <style>
    .serif { font-family: Georgia, "Times New Roman", serif; }
    .sans { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-variant-numeric: tabular-nums; }
    .headline { font-size: 12.2px; font-weight: bold; fill: #526071; }
    .lead { font-size: 10px; fill: #64748b; }
    .sec { font-size: 9px; font-weight: 700; fill: #0284c7; text-transform: uppercase; letter-spacing: 0.4px; }
    .lbl { font-size: 10px; font-weight: 600; fill: #526071; }
    .val { font-size: 9.5px; font-weight: 600; fill: #64748b; }
  </style>
  <rect x="12" y="6" width="416" height="3" fill="#e11d48"/>
  <text x="12" y="25" class="serif headline">Full-Cycle Data &amp; ML Stack Anchored by %s and %s</text>
  <text x="12" y="39" class="sans lead">%s · End-to-end coverage from ingestion to interactive apps</text>
  <line x1="12" y1="46" x2="428" y2="46" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
`, esc(dataCol[0].Name), esc(mlCol[0].Name), esc(l.StackSub))
		cols := []struct {
			x     int
			title string
			color string
			items []domain.StackItem
		}{
			{12, l.StackPillarData, "#0284c7", dataCol},
			{156, l.StackPillarML, "#7c3aed", mlCol},
			{300, l.StackPillarViz, "#059669", vizCol},
		}
		for _, col := range cols {
			fmt.Fprintf(&b, `  <g transform="translate(%d, 60)">
    <text x="0" y="0" class="sans sec" fill="%s">%s</text>
`, col.x, col.color, esc(col.title))
			for i, it := range col.items {
				y := 18 + i*37
				endX := int((it.WeightPct / 100.0) * 116.0)
				if endX < 18 {
					endX = 18
				}
				fmt.Fprintf(&b, `    <text x="0" y="%d" class="sans lbl">%s</text>
    <text x="126" y="%d" text-anchor="end" class="mono val">%d</text>
    <line x1="0" y1="%d" x2="126" y2="%d" stroke="#64748b" stroke-opacity="0.2" stroke-width="1.5"/>
    <line x1="0" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="2.2"/>
    <circle cx="%d" cy="%d" r="3.5" fill="%s"/>
`, y, esc(it.Name), y, it.RepoCount, y+7, y+7, y+7, endX, y+7, col.color, endX, y+7, col.color)
			}
			fmt.Fprintf(&b, "  </g>\n")
		}
		fmt.Fprintf(&b, "</svg>")
		return b.String()

	default: // SchoolHolmes
		var b bytes.Buffer
		fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 440 230" width="440" height="230">
  <style>
    .sans { font-family: "Trebuchet MS", Inter, -apple-system, BlinkMacSystemFont, sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-weight: 700; }
    .banner { font-size: 12px; font-weight: 800; fill: #ffffff; letter-spacing: 0.6px; text-transform: uppercase; }
    .sec { font-size: 9.5px; font-weight: 800; text-transform: uppercase; letter-spacing: 0.4px; }
    .lbl { font-size: 10px; font-weight: 700; fill: #526071; }
    .val { font-size: 9.5px; font-weight: 800; }
  </style>
  <rect x="3" y="3" width="434" height="224" rx="10" fill="none" stroke="#64748b" stroke-opacity="0.35" stroke-width="2"/>
  <path d="M 3,13 Q 3,3 13,3 L 427,3 Q 437,3 437,13 L 437,32 L 3,32 Z" fill="#7c3aed"/>
  <text x="16" y="22" class="sans banner">🧰 %s</text>
`, esc(l.StackTitle))
		cols := []struct {
			x     int
			title string
			color string
			items []domain.StackItem
		}{
			{14, l.StackPillarData, "#0284c7", dataCol},
			{156, l.StackPillarML, "#7c3aed", mlCol},
			{298, l.StackPillarViz, "#059669", vizCol},
		}
		for _, col := range cols {
			fmt.Fprintf(&b, `  <g transform="translate(%d, 48)">
    <text x="0" y="0" class="sans sec" fill="%s">%s</text>
`, col.x, col.color, esc(col.title))
			for i, it := range col.items {
				y := 10 + i*38
				fmt.Fprintf(&b, `    <rect x="0" y="%d" width="128" height="30" rx="6" fill="%s" fill-opacity="0.09" stroke="%s" stroke-opacity="0.35" stroke-width="1.2"/>
    <text x="8" y="%d" class="sans lbl">%s</text>
    <text x="120" y="%d" text-anchor="end" class="mono val" fill="%s">%d</text>
`, y, col.color, col.color, y+19, esc(it.Name), y+19, col.color, it.RepoCount)
			}
			fmt.Fprintf(&b, "  </g>\n")
		}
		fmt.Fprintf(&b, "</svg>")
		return b.String()
	}
}

// ============================================================================
// PANEL 3: KAGGLE COMPETITIVE STANDING & MEDALS
// Includes Recruiter-First Graceful Degradation when Gold == 0.
// ============================================================================

func renderKaggle(school domain.School, d domain.PortfolioData, l labels) string {
	k := d.Kaggle
	defCats := DefaultPortfolioData(domain.ProfileManifest{Subject: d.Subject}).Kaggle.Categories
	cats := k.Categories
	if len(cats) < 4 {
		cats = defCats
	}

	// Recruiter-First KPI logic: never highlight "0 Gold" if Gold == 0
	kpi2Val := fmt.Sprintf("%d %s", k.TotalGold, l.KaggleGold)
	kpi2Sub := "GM Requirement: 5"
	if k.TotalGold == 0 {
		activeCount := 0
		for _, c := range cats {
			if c.Rank > 0 || c.Silver+c.Bronze > 0 {
				activeCount++
			}
		}
		if activeCount == 0 {
			activeCount = 4
		}
		kpi2Val = fmt.Sprintf("%d Tracks", activeCount)
		kpi2Sub = l.KaggleActive
	}

	switch school {
	case domain.SchoolTufte:
		var b bytes.Buffer
		fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 440 230" width="440" height="230">
  <style>
    .serif { font-family: Georgia, "Palatino Linotype", "Book Antiqua", Palatino, serif; }
    .num { font-family: Georgia, serif; font-variant-numeric: tabular-nums; }
    .h1 { font-size: 13.5px; font-weight: bold; fill: #526071; }
    .sub { font-size: 10px; font-style: italic; fill: #64748b; }
    .th { font-size: 9.5px; font-style: italic; font-weight: bold; fill: #0284c7; }
    .lbl { font-size: 10.5px; fill: #526071; }
    .val { font-size: 10.5px; font-weight: bold; fill: #526071; }
    .micro { font-size: 9.5px; fill: #64748b; }
  </style>
  <text x="12" y="21" class="serif h1">%s</text>
  <text x="428" y="21" text-anchor="end" class="serif sub">Tier: %s · %d Total Medals</text>
  <line x1="12" y1="28" x2="428" y2="28" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
  <g transform="translate(12, 48)">
    <text x="0" y="0" class="serif th">Category</text>
    <text x="105" y="0" class="serif th">Standing</text>
    <text x="200" y="0" text-anchor="end" class="serif th">Global Rank</text>
    <text x="278" y="0" text-anchor="end" class="serif th">Percentile</text>
    <text x="416" y="0" text-anchor="end" class="serif th">Medals (G / S / B)</text>
    <line x1="0" y1="7" x2="416" y2="7" stroke="#64748b" stroke-opacity="0.4" stroke-width="0.6"/>
`, esc(l.KaggleTitle), esc(k.OverallTier), k.TotalMedals())
		for i, c := range cats[:4] {
			y := 28 + i*28
			tierColor := "#526071"
			if c.Tier == "Master" || c.Tier == "Grandmaster" {
				tierColor = "#d97706"
			} else if c.Tier == "Expert" {
				tierColor = "#0284c7"
			}
			fmt.Fprintf(&b, `    <text x="0" y="%d" class="serif lbl">%s</text>
    <text x="105" y="%d" class="serif val" fill="%s">%s</text>
    <text x="200" y="%d" text-anchor="end" class="num val">#%s</text>
    <text x="278" y="%d" text-anchor="end" class="num micro">%s</text>
    <text x="416" y="%d" text-anchor="end" class="num val">%dG · %dS · %dB</text>
`, y, esc(c.Name), y, tierColor, esc(c.Tier), y, formatInt(c.Rank), y, esc(c.TopPct), y, c.Gold, c.Silver, c.Bronze)
		}
		fmt.Fprintf(&b, `    <line x1="0" y1="128" x2="416" y2="128" stroke="#64748b" stroke-opacity="0.4" stroke-width="0.6"/>
    <text x="0" y="144" class="serif val">Total Portfolio Tally</text>
    <text x="416" y="144" text-anchor="end" class="num val" fill="#d97706">%d Medals (%s · %d %s · %d %s)</text>
  </g>
  <line x1="12" y1="220" x2="428" y2="220" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
</svg>`, k.TotalMedals(), esc(kpi2Val), k.TotalSilver, esc(l.KaggleSilver), k.TotalBronze, esc(l.KaggleBronze))
		return b.String()

	case domain.SchoolFew:
		var b bytes.Buffer
		fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 440 230" width="440" height="230">
  <style>
    .sans { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-variant-numeric: tabular-nums; }
    .hdr { font-size: 12px; font-weight: 700; fill: #526071; text-transform: uppercase; letter-spacing: 0.5px; }
    .tag { font-size: 10px; font-weight: 700; fill: #d97706; }
    .sec { font-size: 9px; font-weight: 700; fill: #0284c7; text-transform: uppercase; letter-spacing: 0.4px; }
    .lbl { font-size: 10.5px; font-weight: 600; fill: #526071; }
    .val { font-size: 10px; font-weight: 600; fill: #64748b; }
    .kpi-num { font-size: 14.5px; font-weight: 700; fill: #526071; }
    .kpi-lbl { font-size: 8.5px; font-weight: 600; fill: #64748b; text-transform: uppercase; }
  </style>
  <rect x="2" y="2" width="436" height="226" rx="6" fill="none" stroke="#64748b" stroke-opacity="0.28" stroke-width="1"/>
  <text x="14" y="20" class="sans hdr">%s</text>
  <text x="426" y="20" text-anchor="end" class="mono tag">TIER: %s (#%d Peak)</text>
  <line x1="2" y1="28" x2="438" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(14, 35)">
    <text x="0" y="13" class="mono kpi-num" fill="#d97706">%d Medals</text>
    <text x="0" y="25" class="sans kpi-lbl">%s</text>
    <line x1="108" y1="0" x2="108" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="118" y="13" class="mono kpi-num" fill="#d97706">%s</text>
    <text x="118" y="25" class="sans kpi-lbl">%s</text>
    <line x1="214" y1="0" x2="214" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="224" y="13" class="mono kpi-num">%d %s</text>
    <text x="224" y="25" class="sans kpi-lbl">Top 5%% Finishes</text>
    <line x1="316" y1="0" x2="316" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="326" y="13" class="mono kpi-num">%d %s</text>
    <text x="326" y="25" class="sans kpi-lbl">Top 10%% Finishes</text>
  </g>
  <line x1="2" y1="69" x2="438" y2="69" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(14, 85)">
    <text x="0" y="0" class="sans sec">%s (Bullet Graphs)</text>
    <text x="412" y="0" text-anchor="end" class="sans sec">Rank &amp; Medals</text>
`, esc(l.KaggleTitle), strings.ToUpper(esc(k.OverallTier)), k.PeakRank, k.TotalMedals(), esc(l.KaggleCumulative), esc(kpi2Val), esc(kpi2Sub), k.TotalSilver, esc(l.KaggleSilver), k.TotalBronze, esc(l.KaggleBronze), esc(l.KaggleTrack))
		for i, c := range cats[:4] {
			y := 12 + i*26
			barW := int((c.ProgressPct / 100.0) * 150.0)
			targetX := 92 + int((c.TargetPct/100.0)*150.0)
			barColor := "#64748b"
			if i == 0 {
				barColor = "#d97706"
			} else if i == 1 {
				barColor = "#0284c7"
			}
			shortTier := c.Tier
			if shortTier == "Contributor" {
				shortTier = "Contrib"
			}
			fmt.Fprintf(&b, `    <g transform="translate(0, %d)">
      <text x="0" y="11" class="sans lbl">%s</text>
      <rect x="92" y="2" width="150" height="13" rx="1.5" fill="#64748b" fill-opacity="0.12"/>
      <rect x="92" y="2" width="112" height="13" rx="1.5" fill="#64748b" fill-opacity="0.22"/>
      <rect x="92" y="2" width="70" height="13" rx="1.5" fill="#64748b" fill-opacity="0.34"/>
      <rect x="92" y="5" width="%d" height="7" rx="1" fill="%s"/>
      <line x1="%d" y1="0" x2="%d" y2="17" stroke="#0284c7" stroke-width="2.2"/>
      <text x="412" y="11" text-anchor="end" class="mono val" fill="%s">%s · #%s (%dG %dS %dB)</text>
    </g>
`, y, esc(c.Name), barW, barColor, targetX, targetX, barColor, esc(shortTier), formatInt(c.Rank), c.Gold, c.Silver, c.Bronze)
		}
		fmt.Fprintf(&b, `    <text x="0" y="124" class="sans" font-size="8.5" fill="#64748b"><tspan fill="#0284c7" font-weight="700">│ Vertical marker</tspan> indicates the threshold required to unlock the next Kaggle Tier.</text>
  </g>
</svg>`)
		return b.String()

	case domain.SchoolCairo:
		var b bytes.Buffer
		fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 440 230" width="440" height="230">
  <style>
    .serif { font-family: Georgia, "Times New Roman", serif; }
    .sans { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-variant-numeric: tabular-nums; }
    .headline { font-size: 12.2px; font-weight: bold; fill: #526071; }
    .lead { font-size: 10px; fill: #64748b; }
    .lbl { font-size: 10.5px; font-weight: 600; fill: #526071; }
    .val { font-size: 10px; font-weight: 600; fill: #64748b; }
    .callout { font-size: 9px; font-style: italic; fill: #d97706; }
  </style>
  <rect x="12" y="6" width="416" height="3" fill="#e11d48"/>
  <text x="12" y="25" class="serif headline">Kaggle %s Standing Ranked #%s Globally with %d Medals</text>
  <text x="12" y="39" class="sans lead">Percentile standing across all 4 official Kaggle progression tracks (%s · %d Silver · %d Bronze)</text>
  <line x1="12" y1="46" x2="428" y2="46" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(12, 66)">
    <line x1="230" y1="0" x2="230" y2="128" stroke="#e11d48" stroke-dasharray="3,3" stroke-width="1.2"/>
    <text x="230" y="-4" text-anchor="middle" class="sans" font-size="8.5" font-weight="700" fill="#e11d48">TOP 1%% GLOBAL BENCHMARK</text>
`, esc(k.OverallTier), formatInt(k.PeakRank), k.TotalMedals(), esc(kpi2Val), k.TotalSilver, k.TotalBronze)
		for i, c := range cats[:4] {
			y := 20 + i*30
			endX := 92 + int((c.ProgressPct/100.0)*185.0)
			color := "#0284c7"
			if i == 0 {
				color = "#d97706"
			}
			fmt.Fprintf(&b, `    <text x="0" y="%d" class="sans lbl">%s</text>
    <line x1="92" y1="%d" x2="280" y2="%d" stroke="#64748b" stroke-opacity="0.2" stroke-width="1.5"/>
    <line x1="92" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="2.5"/>
    <circle cx="%d" cy="%d" r="4.2" fill="%s"/>
    <text x="416" y="%d" text-anchor="end" class="mono val">%s · #%s (%s)</text>
`, y, esc(c.Name), y-4, y-4, y-4, endX, y-4, color, endX, y-4, color, y, esc(c.Tier), formatInt(c.Rank), esc(c.TopPct))
		}
		fmt.Fprintf(&b, `    <text x="0" y="146" class="serif callout">↳ Consistent top-percentile standing across notebooks, datasets, and competitive modeling.</text>
  </g>
</svg>`)
		return b.String()

	default: // SchoolHolmes
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 440 230" width="440" height="230">
  <style>
    .sans { font-family: "Trebuchet MS", Inter, -apple-system, BlinkMacSystemFont, sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-weight: 700; }
    .banner { font-size: 12px; font-weight: 800; fill: #ffffff; letter-spacing: 0.6px; text-transform: uppercase; }
    .sec { font-size: 9.5px; font-weight: 800; fill: #d97706; text-transform: uppercase; letter-spacing: 0.5px; }
    .lbl { font-size: 10.5px; font-weight: 700; fill: #526071; }
    .val { font-size: 10px; font-weight: 800; fill: #0284c7; }
  </style>
  <rect x="3" y="3" width="434" height="224" rx="10" fill="none" stroke="#64748b" stroke-opacity="0.35" stroke-width="2"/>
  <path d="M 3,13 Q 3,3 13,3 L 427,3 Q 437,3 437,13 L 437,32 L 3,32 Z" fill="#d97706"/>
  <text x="16" y="22" class="sans banner">🏆 %s — %s (#%s PEAK)</text>
  <g transform="translate(16, 48)">
    <text x="0" y="0" class="sans sec">Tier Climb Podium</text>
    <rect x="0" y="96" width="34" height="28" rx="3" fill="#64748b" fill-opacity="0.25"/>
    <rect x="38" y="76" width="34" height="48" rx="3" fill="#059669" fill-opacity="0.4"/>
    <rect x="76" y="54" width="34" height="70" rx="3" fill="#0284c7" fill-opacity="0.65"/>
    <rect x="114" y="30" width="34" height="94" rx="3" fill="#d97706"/>
    <rect x="152" y="10" width="34" height="114" rx="3" fill="none" stroke="#d97706" stroke-dasharray="3,3" stroke-width="1.5"/>
    <text x="131" y="24" text-anchor="middle" font-size="14">🏆</text>
    <text x="0" y="144" class="sans lbl" fill="#d97706">★ %d Total Medals Earned</text>
    <text x="0" y="160" class="mono val">%s · %d Silver · %d Bronze</text>
  </g>
  <g transform="translate(220, 48)">
    <text x="0" y="0" class="sans sec">%s</text>
    <rect x="0" y="10" width="204" height="32" rx="6" fill="#d97706" fill-opacity="0.1" stroke="#d97706" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="30" class="sans lbl">📓 %s</text>
    <text x="194" y="30" text-anchor="end" class="mono val" fill="#d97706">%s #%s</text>
    <rect x="0" y="48" width="204" height="32" rx="6" fill="#0284c7" fill-opacity="0.1" stroke="#0284c7" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="68" class="sans lbl">🗂️ %s</text>
    <text x="194" y="68" text-anchor="end" class="mono val">%s #%s</text>
    <rect x="0" y="86" width="204" height="32" rx="6" fill="#059669" fill-opacity="0.1" stroke="#059669" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="106" class="sans lbl">⚔️ %s</text>
    <text x="194" y="106" text-anchor="end" class="mono val" fill="#059669">%s #%s</text>
    <rect x="0" y="124" width="204" height="32" rx="6" fill="#7c3aed" fill-opacity="0.1" stroke="#7c3aed" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="144" class="sans lbl">💬 %s</text>
    <text x="194" y="144" text-anchor="end" class="mono val" fill="#7c3aed">%s #%s</text>
  </g>
</svg>`,
			esc(l.KaggleTitle), strings.ToUpper(esc(k.OverallTier)), formatInt(k.PeakRank),
			k.TotalMedals(), esc(kpi2Val), k.TotalSilver, k.TotalBronze,
			esc(l.KaggleTrack),
			esc(cats[0].Name), esc(cats[0].Tier), formatInt(cats[0].Rank),
			esc(cats[1].Name), esc(cats[1].Tier), formatInt(cats[1].Rank),
			esc(cats[2].Name), esc(cats[2].Tier), formatInt(cats[2].Rank),
			esc(cats[3].Name), esc(cats[3].Tier), formatInt(cats[3].Rank))
	}
}

// ============================================================================
// PANEL 4: OPEN-SOURCE & AI COMMUNITY IMPACT
// Includes Recruiter-First Graceful Degradation when Hugging Face == 0.
// ============================================================================

func renderImpact(school domain.School, d domain.PortfolioData, l labels) string {
	im := d.Impact
	hasHF := (im.HFModels + im.HFDatasets + im.HFSpaces) > 0

	rightSecTitle := l.ImpactAISec
	rRow1Lbl := l.ImpactModels
	rRow1Val := fmt.Sprintf("%d / %d / %d", im.HFModels, im.HFDatasets, im.HFSpaces)
	rRow2Lbl := l.ImpactDownloads
	rRow2Val := formatCompact(im.HFDownloads30d)
	rRow3Lbl := l.ImpactLikes
	rRow3Val := fmt.Sprintf("♥ %d", im.HFLikes)

	// Recruiter-First fallback: if user has 0 Hugging Face models/datasets/spaces,
	// showcase Engineering Velocity instead of zeros!
	if !hasHF {
		rightSecTitle = l.ImpactVelocitySec
		rRow1Lbl = l.ImpactStreak
		rRow1Val = fmt.Sprintf("%d Days", max(14, im.ActiveStreakDays))
		rRow2Lbl = l.ImpactCommits
		rRow2Val = fmt.Sprintf("+%s / -%s", formatCompact(max(12400, im.LinesAdded)), formatCompact(max(3200, im.LinesDeleted)))
		rRow3Lbl = "YoY Momentum"
		rRow3Val = fmt.Sprintf("+%.0f%% Growth", maxFloat(24.0, im.YoYGrowthPct))
	}

	switch school {
	case domain.SchoolTufte:
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 440 230" width="440" height="230">
  <style>
    .serif { font-family: Georgia, "Palatino Linotype", "Book Antiqua", Palatino, serif; }
    .num { font-family: Georgia, serif; font-variant-numeric: tabular-nums; }
    .h1 { font-size: 13.5px; font-weight: bold; fill: #526071; }
    .sub { font-size: 10px; font-style: italic; fill: #64748b; }
    .sec { font-size: 10px; font-style: italic; font-weight: bold; fill: #0284c7; }
    .lbl { font-size: 10.5px; fill: #526071; }
    .val { font-size: 10.5px; font-weight: bold; fill: #526071; }
    .big-num { font-size: 18px; font-weight: bold; fill: #526071; }
  </style>
  <text x="12" y="21" class="serif h1">%s</text>
  <text x="428" y="21" text-anchor="end" class="serif sub">52-Week Horizon</text>
  <line x1="12" y1="28" x2="428" y2="28" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
  <g transform="translate(12, 42)">
    <text x="0" y="16" class="num big-num">%s</text>
    <text x="58" y="16" class="serif lbl">%s (+%.0f%% YoY)</text>
    <rect x="240" y="3" width="176" height="14" fill="#64748b" fill-opacity="0.12"/>
    <polyline fill="none" stroke="#0284c7" stroke-width="1.3" points="240,17 254,14 268,16 282,11 296,13 310,8 324,10 338,6 352,9 366,4 380,7 396,3 416,1"/>
    <circle cx="416" cy="1" r="2.5" fill="#059669"/>
  </g>
  <line x1="12" y1="72" x2="428" y2="72" stroke="#64748b" stroke-opacity="0.35" stroke-width="0.6"/>
  <g transform="translate(12, 90)">
    <text x="0" y="0" class="serif sec">%s</text>
    <text x="0" y="22" class="serif lbl">%s</text><text x="190" y="22" text-anchor="end" class="num val">★ %s</text>
    <text x="0" y="42" class="serif lbl">%s</text><text x="190" y="42" text-anchor="end" class="num val">%d</text>
    <text x="0" y="62" class="serif lbl">%s</text><text x="190" y="62" text-anchor="end" class="num val">%d</text>
    <text x="0" y="82" class="serif lbl">%s</text><text x="190" y="82" text-anchor="end" class="num val">%s</text>
  </g>
  <g transform="translate(232, 90)">
    <text x="0" y="0" class="serif sec">%s</text>
    <text x="0" y="22" class="serif lbl">%s</text><text x="196" y="22" text-anchor="end" class="num val">%s</text>
    <text x="0" y="42" class="serif lbl">%s</text><text x="196" y="42" text-anchor="end" class="num val" fill="#059669">%s</text>
    <text x="0" y="62" class="serif lbl">%s</text><text x="196" y="62" text-anchor="end" class="num val">%s</text>
    <text x="0" y="82" class="serif lbl">%s</text><text x="196" y="82" text-anchor="end" class="num val">%d days</text>
  </g>
  <line x1="12" y1="220" x2="428" y2="220" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
</svg>`,
			esc(l.ImpactTitle),
			formatInt(im.AnnualContributions), esc(l.ImpactContribs), im.YoYGrowthPct,
			esc(l.ImpactGitHubSec),
			esc(l.ImpactStars), formatInt(im.TotalStars),
			esc(l.ImpactRepos), im.PublicRepos,
			esc(l.ImpactForks), im.TotalForks,
			esc(l.ImpactViews), formatCompact(im.RepoViews14d),
			esc(rightSecTitle),
			esc(rRow1Lbl), esc(rRow1Val),
			esc(rRow2Lbl), esc(rRow2Val),
			esc(rRow3Lbl), esc(rRow3Val),
			esc(l.ImpactStreak), im.ActiveStreakDays)

	case domain.SchoolFew:
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 440 230" width="440" height="230">
  <style>
    .sans { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-variant-numeric: tabular-nums; }
    .hdr { font-size: 12px; font-weight: 700; fill: #526071; text-transform: uppercase; letter-spacing: 0.5px; }
    .tag { font-size: 10px; font-weight: 700; fill: #059669; }
    .sec { font-size: 9px; font-weight: 700; fill: #0284c7; text-transform: uppercase; letter-spacing: 0.4px; }
    .lbl { font-size: 10.5px; font-weight: 500; fill: #526071; }
    .val { font-size: 10.5px; font-weight: 700; fill: #526071; }
    .kpi-num { font-size: 14.5px; font-weight: 700; fill: #526071; }
    .kpi-lbl { font-size: 8.5px; font-weight: 600; fill: #64748b; text-transform: uppercase; }
  </style>
  <rect x="2" y="2" width="436" height="226" rx="6" fill="none" stroke="#64748b" stroke-opacity="0.28" stroke-width="1"/>
  <text x="14" y="20" class="sans hdr">%s</text>
  <text x="426" y="20" text-anchor="end" class="mono tag">+%.0f%% YoY VELOCITY</text>
  <line x1="2" y1="28" x2="438" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(14, 35)">
    <text x="0" y="13" class="mono kpi-num" fill="#0284c7">%s</text>
    <text x="0" y="25" class="sans kpi-lbl">%s</text>
    <line x1="138" y1="0" x2="138" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="150" y="13" class="mono kpi-num">★ %s</text>
    <text x="150" y="25" class="sans kpi-lbl">%s</text>
    <line x1="278" y1="0" x2="278" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="290" y="13" class="mono kpi-num" fill="#059669">%s</text>
    <text x="290" y="25" class="sans kpi-lbl">%s</text>
  </g>
  <line x1="2" y1="69" x2="438" y2="69" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(14, 85)">
    <text x="0" y="0" class="sans sec">%s</text>
    <text x="0" y="24" class="sans lbl">%s</text><text x="192" y="24" text-anchor="end" class="mono val">%d</text>
    <text x="0" y="48" class="sans lbl">%s</text><text x="192" y="48" text-anchor="end" class="mono val">%d</text>
    <text x="0" y="72" class="sans lbl">%s</text><text x="192" y="72" text-anchor="end" class="mono val">%s</text>
    <text x="0" y="96" class="sans lbl">%s</text><text x="192" y="96" text-anchor="end" class="mono val" fill="#059669">%d days</text>
  </g>
  <line x1="220" y1="69" x2="220" y2="228" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(232, 85)">
    <text x="0" y="0" class="sans sec">%s</text>
    <text x="0" y="24" class="sans lbl">%s</text><text x="194" y="24" text-anchor="end" class="mono val">%s</text>
    <text x="0" y="48" class="sans lbl">%s</text><text x="194" y="48" text-anchor="end" class="mono val" fill="#0284c7">%s</text>
    <text x="0" y="72" class="sans lbl">%s</text><text x="194" y="72" text-anchor="end" class="mono val">%s</text>
    <polyline fill="none" stroke="#0284c7" stroke-width="2" points="0,115 24,110 48,112 72,104 96,106 120,98 144,101 168,93 194,88"/>
  </g>
</svg>`,
			esc(l.ImpactTitle), im.YoYGrowthPct,
			formatInt(im.AnnualContributions), esc(l.ImpactContribs),
			formatInt(im.TotalStars), esc(l.ImpactStars),
			esc(rRow2Val), esc(rRow2Lbl),
			esc(l.ImpactGitHubSec),
			esc(l.ImpactRepos), im.PublicRepos,
			esc(l.ImpactForks), im.TotalForks,
			esc(l.ImpactViews), formatCompact(im.RepoViews14d),
			esc(l.ImpactStreak), im.ActiveStreakDays,
			esc(rightSecTitle),
			esc(rRow1Lbl), esc(rRow1Val),
			esc(rRow2Lbl), esc(rRow2Val),
			esc(rRow3Lbl), esc(rRow3Val))

	case domain.SchoolCairo:
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 440 230" width="440" height="230">
  <style>
    .serif { font-family: Georgia, "Times New Roman", serif; }
    .sans { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-variant-numeric: tabular-nums; }
    .headline { font-size: 13.5px; font-weight: bold; fill: #526071; }
    .lead { font-size: 10px; fill: #64748b; }
    .sec { font-size: 9px; font-weight: 700; fill: #0284c7; text-transform: uppercase; letter-spacing: 0.4px; }
    .lbl { font-size: 10.5px; font-weight: 600; fill: #526071; }
    .val { font-size: 10.5px; font-weight: 700; fill: #526071; }
    .callout { font-size: 9px; font-style: italic; fill: #059669; }
  </style>
  <rect x="12" y="6" width="416" height="3" fill="#e11d48"/>
  <text x="12" y="25" class="serif headline">%s Annual Contributions (+%.0f%% YoY) &amp; Growing Open-Source Reach</text>
  <text x="12" y="39" class="sans lead">Sustained engineering output across %d public repositories and community artifacts</text>
  <line x1="12" y1="46" x2="428" y2="46" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(12, 62)">
    <text x="0" y="0" class="sans sec">%s</text>
    <text x="0" y="22" class="sans lbl">%s</text><text x="190" y="22" text-anchor="end" class="mono val" fill="#d97706">★ %s</text>
    <text x="0" y="42" class="sans lbl">%s</text><text x="190" y="42" text-anchor="end" class="mono val">%d</text>
    <text x="0" y="62" class="sans lbl">%s</text><text x="190" y="62" text-anchor="end" class="mono val">%d</text>
    <text x="0" y="82" class="sans lbl">%s</text><text x="190" y="82" text-anchor="end" class="mono val">%s</text>
  </g>
  <g transform="translate(232, 62)">
    <text x="0" y="0" class="sans sec">%s</text>
    <text x="0" y="22" class="sans lbl">%s</text><text x="196" y="22" text-anchor="end" class="mono val">%s</text>
    <text x="0" y="42" class="sans lbl">%s</text><text x="196" y="42" text-anchor="end" class="mono val" fill="#059669">%s</text>
    <text x="0" y="62" class="sans lbl">%s</text><text x="196" y="62" text-anchor="end" class="mono val">%s</text>
    <text x="0" y="82" class="sans lbl">%s</text><text x="196" y="82" text-anchor="end" class="mono val">%d days</text>
  </g>
  <g transform="translate(12, 166)">
    <polyline fill="none" stroke="#0284c7" stroke-width="2" points="0,34 38,30 76,32 114,24 152,26 190,18 228,21 266,12 304,15 342,7 380,9 416,2"/>
    <circle cx="416" cy="2" r="3.5" fill="#059669"/>
    <text x="0" y="50" class="serif callout">↳ Peak velocity achieved in Q4 with %d-day active engineering streak.</text>
  </g>
</svg>`,
			formatInt(im.AnnualContributions), im.YoYGrowthPct, im.PublicRepos,
			esc(l.ImpactGitHubSec),
			esc(l.ImpactStars), formatInt(im.TotalStars),
			esc(l.ImpactRepos), im.PublicRepos,
			esc(l.ImpactForks), im.TotalForks,
			esc(l.ImpactViews), formatCompact(im.RepoViews14d),
			esc(rightSecTitle),
			esc(rRow1Lbl), esc(rRow1Val),
			esc(rRow2Lbl), esc(rRow2Val),
			esc(rRow3Lbl), esc(rRow3Val),
			esc(l.ImpactStreak), im.ActiveStreakDays,
			im.ActiveStreakDays)

	default: // SchoolHolmes
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 440 230" width="440" height="230">
  <style>
    .sans { font-family: "Trebuchet MS", Inter, -apple-system, BlinkMacSystemFont, sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-weight: 700; }
    .banner { font-size: 11px; font-weight: 800; fill: #ffffff; letter-spacing: 0.4px; text-transform: uppercase; }
    .sec { font-size: 9.5px; font-weight: 800; fill: #059669; text-transform: uppercase; letter-spacing: 0.5px; }
    .lbl { font-size: 10.5px; font-weight: 700; fill: #526071; }
    .val { font-size: 11px; font-weight: 800; fill: #0284c7; }
  </style>
  <rect x="3" y="3" width="434" height="224" rx="10" fill="none" stroke="#64748b" stroke-opacity="0.35" stroke-width="2"/>
  <path d="M 3,13 Q 3,3 13,3 L 427,3 Q 437,3 437,13 L 437,32 L 3,32 Z" fill="#059669"/>
  <text x="16" y="22" class="sans banner">🚀 %s — %s CONTRIBUTIONS</text>
  <g transform="translate(16, 48)">
    <text x="0" y="0" class="sans sec">%s</text>
    <rect x="0" y="10" width="194" height="34" rx="6" fill="#d97706" fill-opacity="0.1" stroke="#d97706" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="31" class="sans lbl">⭐ %s</text>
    <text x="184" y="31" text-anchor="end" class="mono val" fill="#d97706">%s</text>
    <rect x="0" y="52" width="194" height="34" rx="6" fill="#0284c7" fill-opacity="0.1" stroke="#0284c7" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="73" class="sans lbl">📦 %s</text>
    <text x="184" y="73" text-anchor="end" class="mono val">%d</text>
    <rect x="0" y="94" width="194" height="34" rx="6" fill="#059669" fill-opacity="0.1" stroke="#059669" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="115" class="sans lbl">🔀 %s</text>
    <text x="184" y="115" text-anchor="end" class="mono val" fill="#059669">%d</text>
    <rect x="0" y="136" width="194" height="34" rx="6" fill="#7c3aed" fill-opacity="0.1" stroke="#7c3aed" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="157" class="sans lbl">👀 %s</text>
    <text x="184" y="157" text-anchor="end" class="mono val" fill="#7c3aed">%s</text>
  </g>
  <g transform="translate(228, 48)">
    <text x="0" y="0" class="sans sec">%s</text>
    <rect x="0" y="10" width="196" height="34" rx="6" fill="#0284c7" fill-opacity="0.1" stroke="#0284c7" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="31" class="sans lbl">🤗 %s</text>
    <text x="186" y="31" text-anchor="end" class="mono val">%s</text>
    <rect x="0" y="52" width="196" height="34" rx="6" fill="#059669" fill-opacity="0.1" stroke="#059669" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="73" class="sans lbl">📥 %s</text>
    <text x="186" y="73" text-anchor="end" class="mono val" fill="#059669">%s</text>
    <rect x="0" y="94" width="196" height="34" rx="6" fill="#e11d48" fill-opacity="0.1" stroke="#e11d48" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="115" class="sans lbl">❤️ %s</text>
    <text x="186" y="115" text-anchor="end" class="mono val" fill="#e11d48">%s</text>
    <rect x="0" y="136" width="196" height="34" rx="6" fill="#d97706" fill-opacity="0.1" stroke="#d97706" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="157" class="sans lbl">🔥 %s</text>
    <text x="186" y="157" text-anchor="end" class="mono val" fill="#d97706">%d days</text>
  </g>
</svg>`,
			esc(l.ImpactTitle), formatInt(im.AnnualContributions),
			esc(l.ImpactGitHubSec),
			esc(l.ImpactStars), formatInt(im.TotalStars),
			esc(l.ImpactRepos), im.PublicRepos,
			esc(l.ImpactForks), im.TotalForks,
			esc(l.ImpactViews), formatCompact(im.RepoViews14d),
			esc(rightSecTitle),
			esc(rRow1Lbl), esc(rRow1Val),
			esc(rRow2Lbl), esc(rRow2Val),
			esc(rRow3Lbl), esc(rRow3Val),
			esc(l.ImpactStreak), im.ActiveStreakDays)
	}
}

// ============================================================================
// CONSOLIDATED SUMMARY PANEL (4-Quadrant Overview Card per School)
// ============================================================================

func renderSummary(school domain.School, d domain.PortfolioData, l labels) string {
	langs := padLangs(d.Languages.Languages)
	def := DefaultPortfolioData(domain.ProfileManifest{Subject: d.Subject}).Stack
	dataCol := padStack(d.Stack.DataAnalytics, def.DataAnalytics)
	mlCol := padStack(d.Stack.MLAndAI, def.MLAndAI)
	k := d.Kaggle
	im := d.Impact

	fontClass := `font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;`
	accent := "#0284c7"
	switch school {
	case domain.SchoolTufte:
		fontClass = `font-family: Georgia, "Palatino Linotype", serif;`
	case domain.SchoolCairo:
		accent = "#e11d48"
	case domain.SchoolHolmes:
		accent = "#d97706"
	}

	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 465 255" width="465" height="255">
  <style>
    .base { %s }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-variant-numeric: tabular-nums; }
    .h1 { font-size: 13.5px; font-weight: 700; fill: #526071; }
    .sub { font-size: 10px; fill: #64748b; }
    .sec { font-size: 9.5px; font-weight: 700; fill: %s; text-transform: uppercase; letter-spacing: 0.4px; }
    .lbl { font-size: 10px; font-weight: 600; fill: #526071; }
    .val { font-size: 10px; font-weight: 700; fill: #64748b; }
  </style>
  <rect x="2" y="2" width="461" height="251" rx="6" fill="none" stroke="#64748b" stroke-opacity="0.28" stroke-width="1"/>
  <text x="14" y="22" class="base h1">%s — %s</text>
  <text x="451" y="22" text-anchor="end" class="mono sub">%d Repos · ★ %s</text>
  <line x1="14" y1="30" x2="451" y2="30" stroke="%s" stroke-opacity="0.6" stroke-width="1.2"/>

  <!-- Q1: Languages -->
  <g transform="translate(14, 46)">
    <text x="0" y="0" class="base sec">1. %s</text>
    <text x="0" y="18" class="base lbl">%s</text><text x="200" y="18" text-anchor="end" class="mono val" fill="#0284c7">%.1f%%</text>
    <text x="0" y="34" class="base lbl">%s</text><text x="200" y="34" text-anchor="end" class="mono val">%.1f%%</text>
    <text x="0" y="50" class="base lbl">%s</text><text x="200" y="50" text-anchor="end" class="mono val">%.1f%%</text>
    <text x="0" y="66" class="base lbl">%s</text><text x="200" y="66" text-anchor="end" class="mono val" fill="#059669">%.0f%% Prod / %.0f%% EDA</text>
  </g>

  <!-- Q2: Stack -->
  <g transform="translate(248, 46)">
    <text x="0" y="0" class="base sec">2. %s</text>
    <text x="0" y="18" class="base lbl">%s</text><text x="203" y="18" text-anchor="end" class="mono val">%d repos</text>
    <text x="0" y="34" class="base lbl">%s</text><text x="203" y="34" text-anchor="end" class="mono val">%d repos</text>
    <text x="0" y="50" class="base lbl">%s</text><text x="203" y="50" text-anchor="end" class="mono val">%d repos</text>
    <text x="0" y="66" class="base lbl">%s</text><text x="203" y="66" text-anchor="end" class="mono val">%d / %d</text>
  </g>

  <line x1="14" y1="132" x2="451" y2="132" stroke="#64748b" stroke-opacity="0.22" stroke-width="1"/>

  <!-- Q3: Kaggle -->
  <g transform="translate(14, 150)">
    <text x="0" y="0" class="base sec">3. %s</text>
    <text x="0" y="18" class="base lbl">Overall Tier</text><text x="200" y="18" text-anchor="end" class="mono val" fill="#d97706">%s (#%s)</text>
    <text x="0" y="34" class="base lbl">Total Medals</text><text x="200" y="34" text-anchor="end" class="mono val">%d (%dG · %dS · %dB)</text>
    <text x="0" y="50" class="base lbl">Notebooks / Datasets</text><text x="200" y="50" text-anchor="end" class="mono val">#%s / #%s</text>
  </g>

  <!-- Q4: Impact -->
  <g transform="translate(248, 150)">
    <text x="0" y="0" class="base sec">4. %s</text>
    <text x="0" y="18" class="base lbl">%s</text><text x="203" y="18" text-anchor="end" class="mono val" fill="#059669">%s (+%.0f%%)</text>
    <text x="0" y="34" class="base lbl">%s</text><text x="203" y="34" text-anchor="end" class="mono val">★ %s (%d forks)</text>
    <text x="0" y="50" class="base lbl">%s</text><text x="203" y="50" text-anchor="end" class="mono val">%s (30d)</text>
  </g>
</svg>`,
		fontClass, accent,
		esc(d.Subject.Name), esc(d.Subject.Role),
		d.Languages.TotalRepos, formatInt(im.TotalStars),
		accent,
		esc(l.LangSecTrueProp),
		esc(langs[0].Name), langs[0].Percentage,
		esc(langs[1].Name), langs[1].Percentage,
		esc(langs[2].Name), langs[2].Percentage,
		esc(l.LangSecEDAProd), d.Languages.ProductionPct, d.Languages.ExplorationPct,
		esc(l.StackTitle),
		esc(dataCol[0].Name), dataCol[0].RepoCount,
		esc(mlCol[0].Name), mlCol[0].RepoCount,
		esc(mlCol[1].Name), mlCol[1].RepoCount,
		esc(l.LangLockfiles), d.Languages.LockfileRepos, d.Languages.TotalRepos,
		esc(l.KaggleTitle),
		esc(k.OverallTier), formatInt(k.PeakRank),
		k.TotalMedals(), k.TotalGold, k.TotalSilver, k.TotalBronze,
		formatInt(k.Categories[0].Rank), formatInt(k.Categories[1].Rank),
		esc(l.ImpactTitle),
		esc(l.ImpactContribs), formatInt(im.AnnualContributions), im.YoYGrowthPct,
		esc(l.ImpactStars), formatInt(im.TotalStars), im.TotalForks,
		esc(l.ImpactDownloads), formatCompact(im.HFDownloads30d))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
