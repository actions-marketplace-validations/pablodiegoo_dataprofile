package render

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/pablodiegoo/dataprofile/internal/domain"
)

type labels struct {
	IsPtBR                bool
	FewBulletGraphsSuffix string

	// Panel 1: Languages
	LangTitleTufte      string
	LangSubTufte        string
	LangSecTrueProp     string
	LangSecEDAProd      string
	LangProdModules     string
	LangEDANotebooks    string
	LangSecHygiene      string
	LangLockfiles       string
	LangCICD            string
	LangTests           string
	LangDocker          string
	LangFootnote        string
	LangFewStrippedFmt  string
	LangFewPrimary      string
	LangFewProdEDA      string
	LangFewLockCI       string
	LangFewProdCode     string
	LangFewTestSuite    string
	LangCairoHeadFmt    string
	LangCairoLeadFmt    string
	LangCairoCallout    string
	LangCairoWaffle     string
	LangHolmesPictorial string
	LangHolmesFilter    string
	LangHolmesStripped  string
	LangHolmesProdSplit string

	// Panel 2: Stack
	StackTitle        string
	StackSub          string
	StackPillarData   string
	StackPillarML     string
	StackPillarViz    string
	StackFootnote     string
	StackCairoHeadFmt string
	StackCairoLeadFmt string
	StackCairoCallout string

	// Panel 3: Kaggle
	KaggleTitle          string
	KaggleCumulative     string
	KaggleGold           string
	KaggleSilver         string
	KaggleBronze         string
	KaggleActive         string
	KaggleTracksFmt      string
	KaggleTrack          string
	KaggleGMReq          string
	KaggleTufteSubFmt    string
	KaggleThCat          string
	KaggleThStanding     string
	KaggleThRank         string
	KaggleThPct          string
	KaggleThMedals       string
	KaggleTotalTally     string
	KaggleMedalsWord     string
	KaggleFewTierTagFmt  string
	KaggleTop5Finishes   string
	KaggleTop10Finishes  string
	KaggleFewMarkerNote  string
	KaggleCairoHeadFmt   string
	KaggleCairoLeadFmt   string
	KaggleCairoBench     string
	KaggleCairoCallout   string
	KaggleHolmesPodium   string
	KaggleHolmesMedals   string
	KaggleOverallTierLbl string
	KaggleTotalMedalsLbl string

	// Panel 4: Impact
	ImpactTitle          string
	ImpactContribs       string
	ImpactGitHubSec      string
	ImpactAISec          string
	ImpactStars          string
	ImpactRepos          string
	ImpactForks          string
	ImpactViews          string
	ImpactModels         string
	ImpactDownloads      string
	ImpactLikes          string
	ImpactVelocitySec    string
	ImpactStreak         string
	ImpactCommits        string
	ImpactYoYMomentum    string
	ImpactDaysCapFmt     string
	ImpactGrowthFmt      string
	ImpactTufteHorizon   string
	ImpactFewVelocityFmt string
	ImpactCairoHeadFmt   string
	ImpactCairoLeadFmt   string
	ImpactCairoCallout   string
	ImpactHolmesBanner   string

	// Header, Bio & Full Panels
	BioTitle                   string
	BioSubFmt                  string
	BioCairoCalloutFmt         string
	BioFallbackP1Fmt           string
	BioFallbackP2Fmt           string
	BioFallbackKaggleFmt       string
	BioFallbackKaggleMedalsFmt string
	HeaderCairoCalloutFmt      string
	HeaderTrendLbl             string
}

func getLabels(locale string, strippedMB float64, totalRepos int) labels {
	if strings.EqualFold(locale, "pt-BR") || strings.EqualFold(locale, "pt") {
		return labels{
			IsPtBR:                true,
			FewBulletGraphsSuffix: "(Gráficos de Bala)",
			LangTitleTufte:        "Distribuição de Código & Engenharia de Produção",
			LangSubTufte:          fmt.Sprintf("Notebooks Normalizados · %d Repos", totalRepos),
			LangSecTrueProp:       "Proporção Real de Código-Fonte",
			LangSecEDAProd:        "Exploração vs. Produção",
			LangProdModules:       "Módulos de Produção",
			LangEDANotebooks:      "Notebooks de EDA",
			LangSecHygiene:        "Maturidade de Engenharia",
			LangLockfiles:         "Lockfiles de Dependência",
			LangCICD:              "Pipelines de CI/CD",
			LangTests:             "Suítes de Testes",
			LangDocker:            "Containers / Docker",
			LangFootnote:          fmt.Sprintf("*Nota: %.1f MB de cache base64 em arquivos .ipynb removidos antes do cálculo.", strippedMB),
			LangFewStrippedFmt:    "%.1f MB Removidos",
			LangFewPrimary:        "Linguagem Principal",
			LangFewProdEDA:        "Produção vs. EDA",
			LangFewLockCI:         "Lockfile & CI Prontos",
			LangFewProdCode:       "Cód. Prod",
			LangFewTestSuite:      "Suíte Testes",
			LangCairoHeadFmt:      "%s Lidera com %.1f%% Após Deflacionar %.1f MB de Cache",
			LangCairoLeadFmt:      "%s · Módulos de produção representam %.0f%% da base de código",
			LangCairoCallout:      "↳ Saídas .ipynb removidas",
			LangCairoWaffle:       "Estrutura de Código Waffle (10×5)",
			LangHolmesPictorial:   "Blocos Pictóricos",
			LangHolmesFilter:      "🧹 Filtro de Cache Jupyter",
			LangHolmesStripped:    "Removidos %.1f MB de PNGs base64",
			LangHolmesProdSplit:   "✓ %.0f%% Produção · %.0f%% Notebooks",
			StackTitle:            "Ecossistema de Dados, ML & Analytics",
			StackSub:              fmt.Sprintf("Dependências Verificadas · %d Repos", totalRepos),
			StackPillarData:       "Dados & Engenharia",
			StackPillarML:         "Machine Learning & IA",
			StackPillarViz:        "Visualização & Apps",
			StackFootnote:         "*Extraído de pyproject.toml, requirements.txt, environment.yml e DESCRIPTION.",
			StackCairoHeadFmt:     "Ecossistema Completo de Dados & ML Ancorado por %s e %s",
			StackCairoLeadFmt:     "%s · Cobertura ponta a ponta da ingestão a aplicações interativas",
			StackCairoCallout:     "↳ Ecossistema principal verificado nos manifestos dos repositórios.",
			KaggleTitle:           "Desempenho Competitivo no Kaggle",
			KaggleCumulative:      "Total Acumulado",
			KaggleGold:            "Ouro",
			KaggleSilver:          "Prata",
			KaggleBronze:          "Bronze",
			KaggleActive:          "Trilhas Ativas",
			KaggleTracksFmt:       "%d Trilhas",
			KaggleTrack:           "Progressão por Categoria",
			KaggleGMReq:           "Requisito GM: 5",
			KaggleTufteSubFmt:     "Tier: %s · %d Medalhas Totais",
			KaggleThCat:           "Categoria",
			KaggleThStanding:      "Nível",
			KaggleThRank:          "Rank Global",
			KaggleThPct:           "Percentil",
			KaggleThMedals:        "Medalhas (O / P / B)",
			KaggleTotalTally:      "Total Acumulado no Portfólio",
			KaggleMedalsWord:      "Medalhas",
			KaggleFewTierTagFmt:   "NÍVEL: %s (#%d Pico)",
			KaggleTop5Finishes:    "Finais Top 5%",
			KaggleTop10Finishes:   "Finais Top 10%",
			KaggleFewMarkerNote:   "indica o limiar necessário para o próximo Kaggle Tier.",
			KaggleCairoHeadFmt:    "Nível %s no Kaggle com Rank Global #%s e %d Medalhas",
			KaggleCairoLeadFmt:    "Percentil nas 4 trilhas oficiais de progressão do Kaggle (%s · %d Prata · %d Bronze)",
			KaggleCairoBench:      "REFERÊNCIA GLOBAL TOP 1%",
			KaggleCairoCallout:    "↳ Desempenho consistente no topo percentual em notebooks, datasets e competições.",
			KaggleHolmesPodium:    "Pódio de Progressão",
			KaggleHolmesMedals:    "★ %d Medalhas Conquistadas",
			KaggleOverallTierLbl:  "Nível Geral",
			KaggleTotalMedalsLbl:  "Medalhas Totais",
			ImpactTitle:           "Impacto Open-Source & Comunidade de IA",
			ImpactContribs:        "Contribuições Anuais",
			ImpactGitHubSec:       "Alcance no GitHub",
			ImpactAISec:           "Artefatos de IA & Hub",
			ImpactStars:           "Estrelas Totais",
			ImpactRepos:           "Repos Públicos",
			ImpactForks:           "Forks & Citações",
			ImpactViews:           "Visitas (14d)",
			ImpactModels:          "Modelos & Datasets",
			ImpactDownloads:       "Downloads (30d)",
			ImpactLikes:           "Curtidas na Comunidade",
			ImpactVelocitySec:     "Consistência & Entrega",
			ImpactStreak:          "Sequência Ativa",
			ImpactCommits:         "Volume de Código",
			ImpactYoYMomentum:     "Impulso Anual (YoY)",
			ImpactDaysCapFmt:      "%d Dias",
			ImpactGrowthFmt:       "+%.0f%% Cresc.",
			ImpactTufteHorizon:    "Horizonte de 52 Semanas",
			ImpactFewVelocityFmt:  "+%.0f%% VELOCIDADE ANUAL",
			ImpactCairoHeadFmt:    "%s Contribuições Anuais (+%.0f%% YoY) & Alcance Open-Source",
			ImpactCairoLeadFmt:    "Produção contínua de engenharia em %d repositórios públicos e artefatos da comunidade",
			ImpactCairoCallout:    "↳ Pico de velocidade atingido no Q4 com sequência ativa de %d dias.",
			ImpactHolmesBanner:    "🚀 %s — %s CONTRIBUIÇÕES",
			BioTitle:                   "Trajetória Profissional & Resumo Executivo",
			BioSubFmt:                  "%d Repositórios · ★ %s Estrelas",
			BioCairoCalloutFmt:         "↳ Portfólio analítico verificado em %d repositórios com %.1f MB de cache Jupyter removido.",
			BioFallbackP1Fmt:           "%s atua como %s, combinando engenharia de software reprodutível, modelagem estatística e arquitetura de dados em %d repositórios verificados.",
			BioFallbackP2Fmt:           "Especialista em %s (%.1f%%), %s (%.1f%%) e %s, com ecossistema ancorado em %s e %s, %s e %s contribuições anuais.",
			BioFallbackKaggleFmt:       "nível %s no Kaggle",
			BioFallbackKaggleMedalsFmt: "nível %s no Kaggle (%d medalhas)",
			HeaderCairoCalloutFmt:      "↳ %s contribuições anuais (+%.0f%% YoY) · Kaggle %s",
			HeaderTrendLbl:             "Tendência 52 Semanas",
		}
	}
	return labels{
		IsPtBR:                false,
		FewBulletGraphsSuffix: "(Bullet Graphs)",
		LangTitleTufte:        "Code Distribution & Production Engineering",
		LangSubTufte:          fmt.Sprintf("Notebook-Normalized · %d Repos", totalRepos),
		LangSecTrueProp:       "True Source Code Proportions",
		LangSecEDAProd:        "Exploration vs. Production",
		LangProdModules:       "Production Modules",
		LangEDANotebooks:      "EDA Notebooks",
		LangSecHygiene:        "Engineering Hygiene Signals",
		LangLockfiles:         "Dependency Lockfiles",
		LangCICD:              "CI/CD Workflows",
		LangTests:             "Automated Test Suites",
		LangDocker:            "Docker / Containers",
		LangFootnote:          fmt.Sprintf("*Note: %.1f MB of embedded base64 chart cache stripped from .ipynb files prior to computing proportions.", strippedMB),
		LangFewStrippedFmt:    "%.1f MB Stripped",
		LangFewPrimary:        "Primary Language",
		LangFewProdEDA:        "Production vs. EDA",
		LangFewLockCI:         "Lockfile & CI Ready",
		LangFewProdCode:       "Prod Code",
		LangFewTestSuite:      "Test Suite",
		LangCairoHeadFmt:      "%s Leads at %.1f%% After Deflating %.1f MB of Cache",
		LangCairoLeadFmt:      "%s · Production modules represent %.0f%% of total codebase",
		LangCairoCallout:      "↳ Stripped .ipynb outputs",
		LangCairoWaffle:       "Code Structure Waffle (10×5)",
		LangHolmesPictorial:   "Pictorial Blocks",
		LangHolmesFilter:      "🧹 Jupyter Bloat Filter",
		LangHolmesStripped:    "Stripped %.1f MB of base64 PNGs",
		LangHolmesProdSplit:   "✓ %.0f%% Production · %.0f%% Notebooks",
		StackTitle:            "Data, ML & Analytics Ecosystem",
		StackSub:              fmt.Sprintf("Verified Manifests · %d Repos", totalRepos),
		StackPillarData:       "Data & Analytics",
		StackPillarML:         "Machine Learning & AI",
		StackPillarViz:        "Visualization & Data Apps",
		StackFootnote:         "*Scanned across pyproject.toml, requirements.txt, environment.yml, and DESCRIPTION.",
		StackCairoHeadFmt:     "Full-Cycle Data & ML Stack Anchored by %s and %s",
		StackCairoLeadFmt:     "%s · End-to-end coverage from ingestion to interactive apps",
		StackCairoCallout:     "↳ Core production stack verified across repository manifests.",
		KaggleTitle:           "Kaggle Competitive Performance",
		KaggleCumulative:      "Cumulative Tally",
		KaggleGold:            "Gold",
		KaggleSilver:          "Silver",
		KaggleBronze:          "Bronze",
		KaggleActive:          "Active Tracks",
		KaggleTracksFmt:       "%d Tracks",
		KaggleTrack:           "Category Progression",
		KaggleGMReq:           "GM Requirement: 5",
		KaggleTufteSubFmt:     "Tier: %s · %d Total Medals",
		KaggleThCat:           "Category",
		KaggleThStanding:      "Standing",
		KaggleThRank:          "Global Rank",
		KaggleThPct:           "Percentile",
		KaggleThMedals:        "Medals (G / S / B)",
		KaggleTotalTally:      "Total Portfolio Tally",
		KaggleMedalsWord:      "Medals",
		KaggleFewTierTagFmt:   "TIER: %s (#%d Peak)",
		KaggleTop5Finishes:    "Top 5% Finishes",
		KaggleTop10Finishes:   "Top 10% Finishes",
		KaggleFewMarkerNote:   "indicates the threshold required to unlock the next Kaggle Tier.",
		KaggleCairoHeadFmt:    "Kaggle %s Standing Ranked #%s Globally with %d Medals",
		KaggleCairoLeadFmt:    "Percentile standing across all 4 official Kaggle progression tracks (%s · %d Silver · %d Bronze)",
		KaggleCairoBench:      "TOP 1% GLOBAL BENCHMARK",
		KaggleCairoCallout:    "↳ Consistent top-percentile standing across notebooks, datasets, and competitive modeling.",
		KaggleHolmesPodium:    "Tier Climb Podium",
		KaggleHolmesMedals:    "★ %d Total Medals Earned",
		KaggleOverallTierLbl:  "Overall Tier",
		KaggleTotalMedalsLbl:  "Total Medals",
		ImpactTitle:           "Open-Source & AI Community Impact",
		ImpactContribs:        "Annual Contributions",
		ImpactGitHubSec:       "GitHub Open-Source Reach",
		ImpactAISec:           "Hugging Face & AI Hub",
		ImpactStars:           "Total Stars Earned",
		ImpactRepos:           "Public Repositories",
		ImpactForks:           "Forks & Citations",
		ImpactViews:           "Repo Views (14d)",
		ImpactModels:          "Models & Datasets",
		ImpactDownloads:       "Downloads (30d)",
		ImpactLikes:           "Community Likes",
		ImpactVelocitySec:     "Engineering Velocity",
		ImpactStreak:          "Active Streak",
		ImpactCommits:         "Net Code Volume",
		ImpactYoYMomentum:     "YoY Momentum",
		ImpactDaysCapFmt:      "%d Days",
		ImpactGrowthFmt:       "+%.0f%% Growth",
		ImpactTufteHorizon:    "52-Week Horizon",
		ImpactFewVelocityFmt:  "+%.0f%% YoY VELOCITY",
		ImpactCairoHeadFmt:    "%s Annual Contributions (+%.0f%% YoY) & Growing Open-Source Reach",
		ImpactCairoLeadFmt:    "Sustained engineering output across %d public repositories and community artifacts",
		ImpactCairoCallout:    "↳ Peak velocity achieved in Q4 with %d-day active engineering streak.",
		ImpactHolmesBanner:    "🚀 %s — %s CONTRIBUTIONS",
		BioTitle:                   "Career Narrative & Executive Bio",
		BioSubFmt:                  "%d Repositories · ★ %s Stars",
		BioCairoCalloutFmt:         "↳ Verified analytical portfolio across %d repositories with %.1f MB of Jupyter cache stripped.",
		BioFallbackP1Fmt:           "%s is a %s combining reproducible software engineering, statistical modeling, and analytics architecture across %d verified repositories.",
		BioFallbackP2Fmt:           "Specializing in %s (%.1f%%), %s (%.1f%%), and %s with a production stack anchored by %s and %s, %s, and %s annual contributions.",
		BioFallbackKaggleFmt:       "Kaggle %s standing",
		BioFallbackKaggleMedalsFmt: "Kaggle %s standing (%d medals)",
		HeaderCairoCalloutFmt:      "↳ %s annual contributions (+%.0f%% YoY) · Kaggle %s",
		HeaderTrendLbl:             "52-Week Velocity",
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

func normalizePortfolioData(data domain.PortfolioData) domain.PortfolioData {
	def := DefaultPortfolioData(domain.ProfileManifest{Subject: data.Subject})
	if len(data.Languages.Languages) == 0 {
		data.Languages = def.Languages
	}
	if len(data.Stack.DataAnalytics) == 0 {
		data.Stack = def.Stack
	} else {
		data.Stack.DataAnalytics = padStack(data.Stack.DataAnalytics, def.Stack.DataAnalytics)
		data.Stack.MLAndAI = padStack(data.Stack.MLAndAI, def.Stack.MLAndAI)
		data.Stack.VizAndApps = padStack(data.Stack.VizAndApps, def.Stack.VizAndApps)
	}
	if len(data.Kaggle.Categories) == 0 {
		data.Kaggle = def.Kaggle
	} else if len(data.Kaggle.Categories) < 4 {
		data.Kaggle.Categories = append(data.Kaggle.Categories, def.Kaggle.Categories[len(data.Kaggle.Categories):4]...)
	}
	if data.Impact.AnnualContributions == 0 && data.Impact.PublicRepos == 0 {
		data.Impact = def.Impact
	}
	return data
}

// Render generates an SVG document for the specified School and Panel.
// Guarantees valid XML, zero School meta-labels, Recruiter-First zero handling,
// and Universal Dual-Surface high contrast on both white (#ffffff) and dark (#0d1117) backgrounds.
func Render(school domain.School, panel domain.Panel, data domain.PortfolioData) (string, error) {
	data = normalizePortfolioData(data)

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
	case domain.PanelHeader:
		svg = renderHeaderSection(school, data, lbl, "hdr-", true)
	case domain.PanelBio:
		svg = renderBio(school, data, lbl, false)
	case domain.PanelBioFull:
		svg = renderBio(school, data, lbl, true)
	case domain.PanelFull:
		svg = renderFull(school, data, lbl)
	default:
		return "", fmt.Errorf("unsupported panel %q", panel)
	}

	const attributionComment = "\n  <!-- Generated by DataProfile (https://github.com/pablodiegoo/dataprofile) | Copyright (c) Pablo Diego | Licensed under Apache-2.0 -->"
	if idx := strings.IndexByte(svg, '>'); idx != -1 && !strings.Contains(svg, "https://github.com/pablodiegoo/dataprofile") {
		svg = svg[:idx+1] + attributionComment + svg[idx+1:]
	}

	if err := ValidateSVG(svg); err != nil {
		return "", fmt.Errorf("invalid SVG generated for %s/%s: %w", school, panel, err)
	}
	return svg, nil
}

// ValidateSVG checks that the SVG string is well-formed XML, contains no external scripts or <foreignObject> tags (ADR-0004),
// contains no low-contrast washed-out pale greys, and contains zero School meta-labels inside the artwork.
func ValidateSVG(svg string) error {
	decoder := xml.NewDecoder(strings.NewReader(svg))
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if se, ok := tok.(xml.StartElement); ok {
			tag := strings.ToLower(se.Name.Local)
			if tag == "script" || tag == "foreignobject" {
				return fmt.Errorf("SVG contains forbidden <%s> element (violates ADR-0004)", se.Name.Local)
			}
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
		strippedTag := fmt.Sprintf(l.LangFewStrippedFmt, lp.StrippedBase64MB)
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
  <text x="426" y="20" text-anchor="end" class="mono tag">%s</text>
  <line x1="2" y1="28" x2="438" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(14, 35)">
    <text x="0" y="12" class="mono kpi-num" fill="#0284c7">%s (%.1f%%)</text>
    <text x="0" y="24" class="sans kpi-lbl">%s</text>
    <line x1="145" y1="0" x2="145" y2="27" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="157" y="12" class="mono kpi-num">%.0f%% / %.0f%%</text>
    <text x="157" y="24" class="sans kpi-lbl">%s</text>
    <line x1="295" y1="0" x2="295" y2="27" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="307" y="12" class="mono kpi-num" fill="#059669">%d / %d Repos</text>
    <text x="307" y="24" class="sans kpi-lbl">%s</text>
  </g>
  <line x1="2" y1="68" x2="438" y2="68" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(14, 84)">
    <text x="0" y="0" class="sans sec">%s</text>
`, esc(l.LangTitleTufte), esc(strippedTag), esc(langs[0].Name), langs[0].Percentage, esc(l.LangFewPrimary), lp.ProductionPct, lp.ExplorationPct, esc(l.LangFewProdEDA), lp.LockfileRepos, lp.TotalRepos, esc(l.LangFewLockCI), esc(l.LangSecTrueProp))
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
		totalRepos := max(1, lp.TotalRepos)
		prodBarW := clampInt(int(math.Round((lp.ProductionPct/100.0)*96.0)), 0, 96)
		lockBarW := clampInt(int(math.Round((float64(lp.LockfileRepos)/float64(totalRepos))*96.0)), 0, 96)
		ciBarW := clampInt(int(math.Round((float64(lp.CICDRepos)/float64(totalRepos))*96.0)), 0, 96)
		testBarW := clampInt(int(math.Round((float64(lp.TestSuiteRepos)/float64(totalRepos))*96.0)), 0, 96)
		dockerBandW := clampInt(int(math.Round((float64(lp.ContainerizedRepos)/float64(totalRepos))*96.0)), 24, 96)
		fmt.Fprintf(&b, `  </g>
  <line x1="222" y1="68" x2="222" y2="228" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(234, 84)">
    <text x="0" y="0" class="sans sec">%s %s</text>
    <g transform="translate(0, 11)">
      <text x="0" y="9" class="sans lbl">%s</text>
      <rect x="64" y="1" width="96" height="10" fill="#64748b" fill-opacity="0.14"/>
      <rect x="64" y="1" width="68" height="10" fill="#64748b" fill-opacity="0.25"/>
      <rect x="64" y="3.5" width="%d" height="5" fill="#059669"/>
      <line x1="126" y1="-1" x2="126" y2="13" stroke="#0284c7" stroke-width="2"/>
      <text x="194" y="9" text-anchor="end" class="mono val">%.0f%%</text>
    </g>
    <g transform="translate(0, 39)">
      <text x="0" y="9" class="sans lbl">Lockfiles</text>
      <rect x="64" y="1" width="96" height="10" fill="#64748b" fill-opacity="0.14"/>
      <rect x="64" y="1" width="68" height="10" fill="#64748b" fill-opacity="0.25"/>
      <rect x="64" y="3.5" width="%d" height="5" fill="#0284c7"/>
      <line x1="136" y1="-1" x2="136" y2="13" stroke="#0284c7" stroke-width="2"/>
      <text x="194" y="9" text-anchor="end" class="mono val">%d/%d</text>
    </g>
    <g transform="translate(0, 67)">
      <text x="0" y="9" class="sans lbl">CI/CD</text>
      <rect x="64" y="1" width="96" height="10" fill="#64748b" fill-opacity="0.14"/>
      <rect x="64" y="1" width="68" height="10" fill="#64748b" fill-opacity="0.25"/>
      <rect x="64" y="3.5" width="%d" height="5" fill="#64748b"/>
      <line x1="126" y1="-1" x2="126" y2="13" stroke="#0284c7" stroke-width="2"/>
      <text x="194" y="9" text-anchor="end" class="mono val">%d/%d</text>
    </g>
    <g transform="translate(0, 95)">
      <text x="0" y="9" class="sans lbl">%s</text>
      <rect x="64" y="1" width="96" height="10" fill="#64748b" fill-opacity="0.14"/>
      <rect x="64" y="1" width="%d" height="10" fill="#64748b" fill-opacity="0.25"/>
      <rect x="64" y="3.5" width="%d" height="5" fill="#64748b"/>
      <line x1="118" y1="-1" x2="118" y2="13" stroke="#0284c7" stroke-width="2"/>
      <text x="194" y="9" text-anchor="end" class="mono val">%d/%d</text>
    </g>
  </g>
</svg>`, esc(l.LangSecHygiene), esc(l.FewBulletGraphsSuffix), esc(l.LangFewProdCode), prodBarW, lp.ProductionPct, lockBarW, lp.LockfileRepos, lp.TotalRepos, ciBarW, lp.CICDRepos, lp.TotalRepos, esc(l.LangFewTestSuite), dockerBandW, testBarW, lp.TestSuiteRepos, lp.TotalRepos)
		return b.String()

	case domain.SchoolCairo:
		var b bytes.Buffer
		head := fmt.Sprintf(l.LangCairoHeadFmt, langs[0].Name, langs[0].Percentage, lp.StrippedBase64MB)
		lead := fmt.Sprintf(l.LangCairoLeadFmt, l.LangSubTufte, lp.ProductionPct)
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
  <text x="12" y="25" class="serif headline">%s</text>
  <text x="12" y="39" class="sans lead">%s</text>
  <line x1="12" y1="46" x2="428" y2="46" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(12, 62)">
    <text x="0" y="0" class="sans sec">%s</text>
    <line x1="144" y1="10" x2="144" y2="118" stroke="#e11d48" stroke-dasharray="3,3" stroke-width="1"/>
`, esc(head), esc(lead), esc(l.LangSecTrueProp))
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
		waffleCells := buildWaffleGridSVG(langs)
		fmt.Fprintf(&b, `    <path d="M 150,24 L 150,32 L 164,32" fill="none" stroke="#0284c7" stroke-width="1"/>
    <text x="168" y="35" class="serif callout">%s</text>
  </g>
  <g transform="translate(244, 62)">
    <text x="0" y="0" class="sans sec">%s</text>
    <g transform="translate(0, 10)">
%s    </g>
    <text x="0" y="75" class="sans val"><tspan fill="#059669" font-weight="700">■ %.0f%% Prod</tspan> (%d) · <tspan fill="#d97706" font-weight="700">■ %.0f%% EDA</tspan> (%d)</text>
    <text x="0" y="95" class="sans sec">%s</text>
    <text x="0" y="112" class="sans lbl">%s</text><text x="184" y="112" text-anchor="end" class="mono val">%d/%d</text>
    <text x="0" y="128" class="sans lbl">%s</text><text x="184" y="128" text-anchor="end" class="mono val">%d/%d</text>
    <text x="0" y="144" class="sans lbl">%s</text><text x="184" y="144" text-anchor="end" class="mono val">%d/%d</text>
  </g>
</svg>`, esc(l.LangCairoCallout), esc(l.LangCairoWaffle), waffleCells, lp.ProductionPct, lp.ProductionFiles, lp.ExplorationPct, lp.ExplorationFiles, esc(l.LangSecHygiene), esc(l.LangLockfiles), lp.LockfileRepos, lp.TotalRepos, esc(l.LangCICD), lp.CICDRepos, lp.TotalRepos, esc(l.LangTests), lp.TestSuiteRepos, lp.TotalRepos)
		return b.String()

	default: // SchoolHolmes
		strippedMsg := fmt.Sprintf(l.LangHolmesStripped, lp.StrippedBase64MB)
		prodSplitMsg := fmt.Sprintf(l.LangHolmesProdSplit, lp.ProductionPct, lp.ExplorationPct)
		isotypeRows := buildHolmesIsotypeRowsSVG(langs)
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
    <text x="0" y="0" class="sans sec">%s (%s)</text>
%s    <rect x="0" y="106" width="196" height="58" rx="8" fill="#0284c7" fill-opacity="0.08" stroke="#0284c7" stroke-opacity="0.35" stroke-width="1.2"/>
    <text x="10" y="123" class="sans sec" fill="#0284c7">%s</text>
    <text x="10" y="140" class="sans badge">%s</text>
    <text x="10" y="154" class="sans badge" fill="#059669">%s</text>
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
			esc(l.LangTitleTufte), esc(l.LangSecTrueProp), esc(l.LangHolmesPictorial),
			isotypeRows,
			esc(l.LangHolmesFilter), esc(strippedMsg), esc(prodSplitMsg),
			esc(l.LangSecHygiene),
			esc(l.LangLockfiles), lp.LockfileRepos, lp.TotalRepos,
			esc(l.LangCICD), lp.CICDRepos, lp.TotalRepos,
			esc(l.LangTests), lp.TestSuiteRepos, lp.TotalRepos,
			esc(l.LangDocker), lp.ContainerizedRepos, lp.TotalRepos)
	}
}

func buildWaffleGridSVG(langs []domain.LanguageStat) string {
	colors := make([]string, 0, 50)
	for _, lg := range langs {
		c := lg.Color
		if c == "" {
			c = "#64748b"
		}
		n := int(math.Round(lg.Percentage / 2.0))
		for k := 0; k < n && len(colors) < 50; k++ {
			colors = append(colors, c)
		}
	}
	fallbackColor := "#059669"
	if len(langs) > 0 && langs[len(langs)-1].Color != "" {
		fallbackColor = langs[len(langs)-1].Color
	}
	for len(colors) < 50 {
		colors = append(colors, fallbackColor)
	}

	var b bytes.Buffer
	for idx := 0; idx < 50; idx++ {
		row := idx / 10
		col := idx % 10
		x := col * 18
		y := row * 11
		fmt.Fprintf(&b, `      <rect x="%d" y="%d" width="15" height="8" rx="1.5" fill="%s"/>`+"\n", x, y, colors[idx])
	}
	return b.String()
}

func buildHolmesIsotypeRowsSVG(langs []domain.LanguageStat) string {
	type rowSpec struct {
		label string
		pct   float64
		color string
	}
	rows := []rowSpec{
		{langs[0].Name, langs[0].Percentage, nonEmptyColor(langs[0].Color, "#0284c7")},
		{langs[1].Name, langs[1].Percentage, nonEmptyColor(langs[1].Color, "#d97706")},
		{langs[2].Name, langs[2].Percentage, nonEmptyColor(langs[2].Color, "#00ADD8")},
		{langs[3].Name + " + " + langs[4].Name, langs[3].Percentage + langs[4].Percentage, "#7c3aed"},
	}

	counts := make([]int, len(rows))
	sum := 0
	maxIdx := 0
	for i, r := range rows {
		c := int(math.Round(r.pct / 10.0))
		if c < 1 && r.pct > 0 {
			c = 1
		}
		counts[i] = c
		sum += c
		if r.pct > rows[maxIdx].pct {
			maxIdx = i
		}
	}
	if sum != 10 && counts[maxIdx]+(10-sum) >= 1 {
		counts[maxIdx] += 10 - sum
	}

	var b bytes.Buffer
	for i, r := range rows {
		yText := 22 + i*23
		yBlocks := 12 + i*23
		fmt.Fprintf(&b, `    <text x="0" y="%d" class="sans lbl">%s</text>`+"\n", yText, esc(r.label))
		fmt.Fprintf(&b, `    <g transform="translate(56, %d)" fill="%s">`+"\n      ", yBlocks, r.color)
		for k := 0; k < counts[i]; k++ {
			fmt.Fprintf(&b, `<rect x="%d" y="0" width="14" height="11" rx="2"/>`, k*17)
		}
		fmt.Fprintf(&b, "\n    </g>\n")
		valX := 56 + counts[i]*17 + 6
		fmt.Fprintf(&b, `    <text x="%d" y="%d" class="mono val" fill="%s">%.1f%%</text>`+"\n", valX, yText, r.color, r.pct)
	}
	return b.String()
}

func nonEmptyColor(c, fallback string) string {
	if strings.TrimSpace(c) == "" {
		return fallback
	}
	return c
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
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
	dataCol := d.Stack.DataAnalytics
	mlCol := d.Stack.MLAndAI
	vizCol := d.Stack.VizAndApps

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
    .kpi-num { font-size: 12.5px; font-weight: 700; fill: #526071; }
    .kpi-lbl { font-size: 8px; font-weight: 600; fill: #64748b; text-transform: uppercase; }
  </style>
  <rect x="2" y="2" width="436" height="226" rx="6" fill="none" stroke="#64748b" stroke-opacity="0.28" stroke-width="1"/>
  <text x="14" y="20" class="sans hdr">%s</text>
  <text x="426" y="20" text-anchor="end" class="mono tag">%s</text>
  <line x1="2" y1="28" x2="438" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(14, 33)">
    <text x="0" y="11" class="mono kpi-num" fill="#0284c7">%s (%d)</text>
    <text x="0" y="22" class="sans kpi-lbl">%s</text>
    <line x1="145" y1="0" x2="145" y2="25" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="157" y="11" class="mono kpi-num" fill="#7c3aed">%s (%d)</text>
    <text x="157" y="22" class="sans kpi-lbl">%s</text>
    <line x1="295" y1="0" x2="295" y2="25" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="307" y="11" class="mono kpi-num" fill="#059669">%s (%d)</text>
    <text x="307" y="22" class="sans kpi-lbl">%s</text>
  </g>
  <line x1="2" y1="63" x2="438" y2="63" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
`, esc(l.StackTitle), esc(l.StackSub),
			esc(dataCol[0].Name), dataCol[0].RepoCount, esc(l.StackPillarData),
			esc(mlCol[0].Name), mlCol[0].RepoCount, esc(l.StackPillarML),
			esc(vizCol[0].Name), vizCol[0].RepoCount, esc(l.StackPillarViz))
		for _, col := range cols {
			fmt.Fprintf(&b, `  <g transform="translate(%d, 76)">
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
				y := 18 + i*33
				barW := int((it.WeightPct / 100.0) * 124.0)
				fill := "#64748b"
				if it.IsLeader || (!hasLeader && i == 0 && it.RepoCount > 0) {
					fill = col.color
				}
				fmt.Fprintf(&b, `    <text x="0" y="%d" class="sans lbl">%s</text>
    <text x="124" y="%d" text-anchor="end" class="mono val">%d</text>
    <rect x="0" y="%d" width="124" height="8" rx="1.5" fill="#64748b" fill-opacity="0.14"/>
    <rect x="0" y="%d" width="84" height="8" rx="1.5" fill="#64748b" fill-opacity="0.25"/>
    <rect x="0" y="%d" width="%d" height="6" rx="1.5" fill="%s"/>
    <line x1="96" y1="%d" x2="96" y2="%d" stroke="#0284c7" stroke-width="2"/>
`, y, esc(it.Name), y, it.RepoCount, y+4, y+4, y+5, barW, fill, y+3, y+13)
			}
			fmt.Fprintf(&b, "  </g>\n")
		}
		fmt.Fprintf(&b, `  <line x1="148" y1="63" x2="148" y2="228" stroke="#64748b" stroke-opacity="0.2" stroke-width="1"/>
  <line x1="292" y1="63" x2="292" y2="228" stroke="#64748b" stroke-opacity="0.2" stroke-width="1"/>
</svg>`)
		return b.String()

	case domain.SchoolCairo:
		var b bytes.Buffer
		head := fmt.Sprintf(l.StackCairoHeadFmt, dataCol[0].Name, mlCol[0].Name)
		lead := fmt.Sprintf(l.StackCairoLeadFmt, l.StackSub)
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
    .callout { font-size: 9px; font-style: italic; fill: #0284c7; }
  </style>
  <rect x="12" y="6" width="416" height="3" fill="#e11d48"/>
  <text x="12" y="25" class="serif headline">%s</text>
  <text x="12" y="39" class="sans lead">%s</text>
  <line x1="12" y1="46" x2="428" y2="46" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
`, esc(head), esc(lead))
		for _, col := range cols {
			fmt.Fprintf(&b, `  <g transform="translate(%d, 60)">
    <text x="0" y="0" class="sans sec" fill="%s">%s</text>
    <line x1="88" y1="8" x2="88" y2="140" stroke="#e11d48" stroke-dasharray="3,3" stroke-width="1"/>
`, col.x, col.color, esc(col.title))
			for i, it := range col.items {
				y := 18 + i*33
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
		fmt.Fprintf(&b, `  <text x="12" y="218" class="serif callout">%s</text>
</svg>`, esc(l.StackCairoCallout))
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
	cats := k.Categories

	// Recruiter-First KPI logic: never highlight "0 Gold" if Gold == 0
	kpi2Val := fmt.Sprintf("%d %s", k.TotalGold, l.KaggleGold)
	kpi2Sub := l.KaggleGMReq
	if k.TotalGold == 0 {
		activeCount := len(cats)
		if activeCount == 0 {
			activeCount = 4
		}
		kpi2Val = fmt.Sprintf(l.KaggleTracksFmt, activeCount)
		kpi2Sub = l.KaggleActive
	}

	switch school {
	case domain.SchoolTufte:
		var b bytes.Buffer
		sub := fmt.Sprintf(l.KaggleTufteSubFmt, k.OverallTier, k.TotalMedals())
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
  <text x="428" y="21" text-anchor="end" class="serif sub">%s</text>
  <line x1="12" y1="28" x2="428" y2="28" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
  <g transform="translate(12, 48)">
    <text x="0" y="0" class="serif th">%s</text>
    <text x="105" y="0" class="serif th">%s</text>
    <text x="200" y="0" text-anchor="end" class="serif th">%s</text>
    <text x="278" y="0" text-anchor="end" class="serif th">%s</text>
    <text x="416" y="0" text-anchor="end" class="serif th">%s</text>
    <line x1="0" y1="7" x2="416" y2="7" stroke="#64748b" stroke-opacity="0.4" stroke-width="0.6"/>
`, esc(l.KaggleTitle), esc(sub), esc(l.KaggleThCat), esc(l.KaggleThStanding), esc(l.KaggleThRank), esc(l.KaggleThPct), esc(l.KaggleThMedals))
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
    <text x="0" y="144" class="serif val">%s</text>
    <text x="416" y="144" text-anchor="end" class="num val" fill="#d97706">%d %s (%s · %d %s · %d %s)</text>
  </g>
  <line x1="12" y1="220" x2="428" y2="220" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
</svg>`, esc(l.KaggleTotalTally), k.TotalMedals(), esc(l.KaggleMedalsWord), esc(kpi2Val), k.TotalSilver, esc(l.KaggleSilver), k.TotalBronze, esc(l.KaggleBronze))
		return b.String()

	case domain.SchoolFew:
		var b bytes.Buffer
		markerLabel := "│ Vertical marker"
		if l.IsPtBR {
			markerLabel = "│ Marcador vertical"
		}
		tierTag := fmt.Sprintf(l.KaggleFewTierTagFmt, strings.ToUpper(k.OverallTier), k.PeakRank)
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
  <text x="426" y="20" text-anchor="end" class="mono tag">%s</text>
  <line x1="2" y1="28" x2="438" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(14, 35)">
    <text x="0" y="13" class="mono kpi-num" fill="#d97706">%d %s</text>
    <text x="0" y="25" class="sans kpi-lbl">%s</text>
    <line x1="108" y1="0" x2="108" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="118" y="13" class="mono kpi-num" fill="#d97706">%s</text>
    <text x="118" y="25" class="sans kpi-lbl">%s</text>
    <line x1="214" y1="0" x2="214" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="224" y="13" class="mono kpi-num">%d %s</text>
    <text x="224" y="25" class="sans kpi-lbl">%s</text>
    <line x1="316" y1="0" x2="316" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="326" y="13" class="mono kpi-num">%d %s</text>
    <text x="326" y="25" class="sans kpi-lbl">%s</text>
  </g>
  <line x1="2" y1="69" x2="438" y2="69" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(14, 85)">
    <text x="0" y="0" class="sans sec">%s %s</text>
    <text x="412" y="0" text-anchor="end" class="sans sec">Rank &amp; %s</text>
`, esc(l.KaggleTitle), esc(tierTag), k.TotalMedals(), esc(l.KaggleMedalsWord), esc(l.KaggleCumulative), esc(kpi2Val), esc(kpi2Sub), k.TotalSilver, esc(l.KaggleSilver), esc(l.KaggleTop5Finishes), k.TotalBronze, esc(l.KaggleBronze), esc(l.KaggleTop10Finishes), esc(l.KaggleTrack), esc(l.FewBulletGraphsSuffix), esc(l.KaggleMedalsWord))
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
		fmt.Fprintf(&b, `    <text x="0" y="124" class="sans" font-size="8.5" fill="#64748b"><tspan fill="#0284c7" font-weight="700">%s</tspan> %s</text>
  </g>
</svg>`, esc(markerLabel), esc(l.KaggleFewMarkerNote))
		return b.String()

	case domain.SchoolCairo:
		var b bytes.Buffer
		head := fmt.Sprintf(l.KaggleCairoHeadFmt, k.OverallTier, formatInt(k.PeakRank), k.TotalMedals())
		lead := fmt.Sprintf(l.KaggleCairoLeadFmt, kpi2Val, k.TotalSilver, k.TotalBronze)
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
  <text x="12" y="25" class="serif headline">%s</text>
  <text x="12" y="39" class="sans lead">%s</text>
  <line x1="12" y1="46" x2="428" y2="46" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(12, 66)">
    <line x1="230" y1="0" x2="230" y2="128" stroke="#e11d48" stroke-dasharray="3,3" stroke-width="1.2"/>
    <text x="230" y="-4" text-anchor="middle" class="sans" font-size="8.5" font-weight="700" fill="#e11d48">%s</text>
`, esc(head), esc(lead), esc(l.KaggleCairoBench))
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
		fmt.Fprintf(&b, `    <text x="0" y="146" class="serif callout">%s</text>
  </g>
</svg>`, esc(l.KaggleCairoCallout))
		return b.String()

	default: // SchoolHolmes
		medalsEarned := fmt.Sprintf(l.KaggleHolmesMedals, k.TotalMedals())
		trophyX, trophyY := holmesTrophyPosition(k.OverallTier)
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
    <text x="0" y="0" class="sans sec">%s</text>
    <rect x="0" y="96" width="34" height="28" rx="3" fill="#64748b" fill-opacity="0.25"/>
    <rect x="38" y="76" width="34" height="48" rx="3" fill="#059669" fill-opacity="0.4"/>
    <rect x="76" y="54" width="34" height="70" rx="3" fill="#0284c7" fill-opacity="0.65"/>
    <rect x="114" y="30" width="34" height="94" rx="3" fill="#d97706"/>
    <rect x="152" y="10" width="34" height="114" rx="3" fill="none" stroke="#d97706" stroke-dasharray="3,3" stroke-width="1.5"/>
    <text x="%d" y="%d" text-anchor="middle" font-size="14">🏆</text>
    <text x="0" y="144" class="sans lbl" fill="#d97706">%s</text>
    <text x="0" y="160" class="mono val">%s · %d %s · %d %s</text>
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
			esc(l.KaggleHolmesPodium), trophyX, trophyY, esc(medalsEarned), esc(kpi2Val), k.TotalSilver, esc(l.KaggleSilver), k.TotalBronze, esc(l.KaggleBronze),
			esc(l.KaggleTrack),
			esc(cats[0].Name), esc(cats[0].Tier), formatInt(cats[0].Rank),
			esc(cats[1].Name), esc(cats[1].Tier), formatInt(cats[1].Rank),
			esc(cats[2].Name), esc(cats[2].Tier), formatInt(cats[2].Rank),
			esc(cats[3].Name), esc(cats[3].Tier), formatInt(cats[3].Rank))
	}
}

func holmesTrophyPosition(tier string) (int, int) {
	switch strings.ToLower(strings.TrimSpace(tier)) {
	case strings.ToLower(domain.KaggleTierNovice):
		return 17, 90
	case strings.ToLower(domain.KaggleTierContributor):
		return 55, 70
	case strings.ToLower(domain.KaggleTierExpert):
		return 93, 48
	case strings.ToLower(domain.KaggleTierGrandmaster):
		return 169, 4
	default: // Master
		return 131, 24
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
		rRow1Val = fmt.Sprintf(l.ImpactDaysCapFmt, max(14, im.ActiveStreakDays))
		rRow2Lbl = l.ImpactCommits
		rRow2Val = fmt.Sprintf("+%s / -%s", formatCompact(max(12400, im.LinesAdded)), formatCompact(max(3200, im.LinesDeleted)))
		rRow3Lbl = l.ImpactYoYMomentum
		rRow3Val = fmt.Sprintf(l.ImpactGrowthFmt, maxFloat(24.0, im.YoYGrowthPct))
	}

	switch school {
	case domain.SchoolTufte:
		pts, lastX, lastY := buildTrendPolyline(im.WeeklyTrend, 240, 1, 176, 16)
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
  <text x="428" y="21" text-anchor="end" class="serif sub">%s</text>
  <line x1="12" y1="28" x2="428" y2="28" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
  <g transform="translate(12, 42)">
    <text x="0" y="16" class="num big-num">%s</text>
    <text x="58" y="16" class="serif lbl">%s (+%.0f%% YoY)</text>
    <rect x="240" y="3" width="176" height="14" fill="#64748b" fill-opacity="0.12"/>
    <polyline fill="none" stroke="#0284c7" stroke-width="1.3" points="%s"/>
    <circle cx="%d" cy="%d" r="2.5" fill="#059669"/>
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
			esc(l.ImpactTitle), esc(l.ImpactTufteHorizon),
			formatInt(im.AnnualContributions), esc(l.ImpactContribs), im.YoYGrowthPct,
			pts, lastX, lastY,
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
		pts, _, _ := buildTrendPolyline(im.WeeklyTrend, 0, 88, 194, 27)
		velTag := fmt.Sprintf(l.ImpactFewVelocityFmt, im.YoYGrowthPct)
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
  <text x="426" y="20" text-anchor="end" class="mono tag">%s</text>
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
    <rect x="0" y="106" width="192" height="8" rx="1.5" fill="#64748b" fill-opacity="0.14"/>
    <rect x="0" y="106" width="128" height="8" rx="1.5" fill="#64748b" fill-opacity="0.25"/>
    <rect x="0" y="108" width="144" height="4" rx="1" fill="#059669"/>
    <line x1="150" y1="104" x2="150" y2="116" stroke="#0284c7" stroke-width="2"/>
  </g>
  <line x1="220" y1="69" x2="220" y2="228" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(232, 85)">
    <text x="0" y="0" class="sans sec">%s</text>
    <text x="0" y="24" class="sans lbl">%s</text><text x="194" y="24" text-anchor="end" class="mono val">%s</text>
    <text x="0" y="48" class="sans lbl">%s</text><text x="194" y="48" text-anchor="end" class="mono val" fill="#0284c7">%s</text>
    <text x="0" y="72" class="sans lbl">%s</text><text x="194" y="72" text-anchor="end" class="mono val">%s</text>
    <polyline fill="none" stroke="#0284c7" stroke-width="2" points="%s"/>
  </g>
</svg>`,
			esc(l.ImpactTitle), esc(velTag),
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
			esc(rRow3Lbl), esc(rRow3Val),
			pts)

	case domain.SchoolCairo:
		pts, lastX, lastY := buildTrendPolyline(im.WeeklyTrend, 0, 2, 416, 32)
		head := fmt.Sprintf(l.ImpactCairoHeadFmt, formatInt(im.AnnualContributions), im.YoYGrowthPct)
		lead := fmt.Sprintf(l.ImpactCairoLeadFmt, im.PublicRepos)
		callout := fmt.Sprintf(l.ImpactCairoCallout, im.ActiveStreakDays)
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
  <text x="12" y="25" class="serif headline">%s</text>
  <text x="12" y="39" class="sans lead">%s</text>
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
    <line x1="0" y1="18" x2="416" y2="18" stroke="#e11d48" stroke-dasharray="3,3" stroke-width="1"/>
    <polyline fill="none" stroke="#0284c7" stroke-width="2" points="%s"/>
    <circle cx="%d" cy="%d" r="3.5" fill="#059669"/>
    <text x="0" y="50" class="serif callout">%s</text>
  </g>
</svg>`,
			esc(head), esc(lead),
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
			pts, lastX, lastY,
			esc(callout))

	default: // SchoolHolmes
		banner := fmt.Sprintf(l.ImpactHolmesBanner, l.ImpactTitle, formatInt(im.AnnualContributions))
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
  <text x="16" y="22" class="sans banner">%s</text>
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
			esc(banner),
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
	dataCol := d.Stack.DataAnalytics
	mlCol := d.Stack.MLAndAI
	k := d.Kaggle
	cats := k.Categories
	im := d.Impact

	switch school {
	case domain.SchoolTufte:
		pts, lastX, lastY := buildTrendPolyline(im.WeeklyTrend, 118, 46, 82, 12)
		l0End := 58 + int(langs[0].Percentage*1.4)
		l1End := 58 + int(langs[1].Percentage*1.4)
		l2End := 58 + int(langs[2].Percentage*1.4)
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 465 255" width="465" height="255">
  <style>
    .serif { font-family: Georgia, "Palatino Linotype", "Book Antiqua", Palatino, serif; }
    .num { font-family: Georgia, serif; font-variant-numeric: tabular-nums; }
    .h1 { font-size: 13.5px; font-weight: bold; fill: #526071; }
    .sub { font-size: 10px; font-style: italic; fill: #64748b; }
    .sec { font-size: 10px; font-style: italic; font-weight: bold; fill: #0284c7; }
    .lbl { font-size: 10px; fill: #526071; }
    .val { font-size: 10px; font-weight: bold; fill: #526071; }
  </style>
  <text x="14" y="22" class="serif h1">%s — %s</text>
  <text x="451" y="22" text-anchor="end" class="num sub">%d Repos · ★ %s</text>
  <line x1="14" y1="30" x2="451" y2="30" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>

  <!-- Q1: Languages (Dot-Dash Plot) -->
  <g transform="translate(14, 46)">
    <text x="0" y="0" class="serif sec">1. %s</text>
    <text x="0" y="18" class="serif lbl">%s</text>
    <line x1="58" y1="14" x2="%d" y2="14" stroke="#64748b" stroke-width="1.1"/>
    <circle cx="%d" cy="14" r="2.5" fill="#0284c7"/>
    <text x="200" y="18" text-anchor="end" class="num val" fill="#0284c7">%.1f%%</text>
    <text x="0" y="34" class="serif lbl">%s</text>
    <line x1="58" y1="30" x2="%d" y2="30" stroke="#64748b" stroke-width="1.1"/>
    <circle cx="%d" cy="30" r="2.5" fill="#d97706"/>
    <text x="200" y="34" text-anchor="end" class="num val">%.1f%%</text>
    <text x="0" y="50" class="serif lbl">%s</text>
    <line x1="58" y1="46" x2="%d" y2="46" stroke="#64748b" stroke-width="1.1"/>
    <circle cx="%d" cy="46" r="2.5" fill="#00ADD8"/>
    <text x="200" y="50" text-anchor="end" class="num val">%.1f%%</text>
    <text x="0" y="66" class="serif lbl">%s</text>
    <text x="200" y="66" text-anchor="end" class="num val" fill="#059669">%.0f%% Prod / %.0f%% EDA</text>
  </g>

  <!-- Q2: Stack (Supertable) -->
  <g transform="translate(248, 46)">
    <text x="0" y="0" class="serif sec">2. %s</text>
    <text x="0" y="18" class="serif lbl">%s</text><text x="203" y="18" text-anchor="end" class="num val">%d repos</text>
    <text x="0" y="34" class="serif lbl">%s</text><text x="203" y="34" text-anchor="end" class="num val">%d repos</text>
    <text x="0" y="50" class="serif lbl">%s</text><text x="203" y="50" text-anchor="end" class="num val">%d repos</text>
    <text x="0" y="66" class="serif lbl">%s</text><text x="203" y="66" text-anchor="end" class="num val">%d / %d</text>
  </g>

  <line x1="14" y1="132" x2="451" y2="132" stroke="#64748b" stroke-opacity="0.35" stroke-width="0.6"/>

  <!-- Q3: Kaggle (Supertable) -->
  <g transform="translate(14, 150)">
    <text x="0" y="0" class="serif sec">3. %s</text>
    <text x="0" y="18" class="serif lbl">%s</text><text x="200" y="18" text-anchor="end" class="num val" fill="#d97706">%s (#%s)</text>
    <text x="0" y="34" class="serif lbl">%s</text><text x="200" y="34" text-anchor="end" class="num val">%d (%dG · %dS · %dB)</text>
    <text x="0" y="50" class="serif lbl">Notebooks / Datasets</text><text x="200" y="50" text-anchor="end" class="num val">#%s / #%s</text>
  </g>

  <!-- Q4: Impact (Supertable + Inline Sparkline with Interquartile Band) -->
  <g transform="translate(248, 150)">
    <text x="0" y="0" class="serif sec">4. %s</text>
    <text x="0" y="18" class="serif lbl">%s</text><text x="203" y="18" text-anchor="end" class="num val" fill="#059669">%s (+%.0f%%)</text>
    <text x="0" y="34" class="serif lbl">%s</text><text x="203" y="34" text-anchor="end" class="num val">★ %s (%d forks)</text>
    <text x="0" y="54" class="serif lbl">%s</text>
    <rect x="118" y="46" width="82" height="12" fill="#64748b" fill-opacity="0.12"/>
    <polyline fill="none" stroke="#0284c7" stroke-width="1.2" points="%s"/>
    <circle cx="%d" cy="%d" r="2.2" fill="#059669"/>
  </g>

  <line x1="14" y1="242" x2="451" y2="242" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
</svg>`,
			esc(d.Subject.Name), esc(d.Subject.Role),
			d.Languages.TotalRepos, formatInt(im.TotalStars),
			esc(l.LangSecTrueProp),
			esc(langs[0].Name), l0End, l0End, langs[0].Percentage,
			esc(langs[1].Name), l1End, l1End, langs[1].Percentage,
			esc(langs[2].Name), l2End, l2End, langs[2].Percentage,
			esc(l.LangSecEDAProd), d.Languages.ProductionPct, d.Languages.ExplorationPct,
			esc(l.StackTitle),
			esc(dataCol[0].Name), dataCol[0].RepoCount,
			esc(mlCol[0].Name), mlCol[0].RepoCount,
			esc(mlCol[1].Name), mlCol[1].RepoCount,
			esc(l.LangLockfiles), d.Languages.LockfileRepos, d.Languages.TotalRepos,
			esc(l.KaggleTitle),
			esc(l.KaggleOverallTierLbl), esc(k.OverallTier), formatInt(k.PeakRank),
			esc(l.KaggleTotalMedalsLbl), k.TotalMedals(), k.TotalGold, k.TotalSilver, k.TotalBronze,
			formatInt(cats[0].Rank), formatInt(cats[1].Rank),
			esc(l.ImpactTitle),
			esc(l.ImpactContribs), formatInt(im.AnnualContributions), im.YoYGrowthPct,
			esc(l.ImpactStars), formatInt(im.TotalStars), im.TotalForks,
			esc(l.ImpactDownloads), pts, lastX, lastY)

	case domain.SchoolFew:
		l0Bar := int(langs[0].Percentage * 0.9)
		prodBar := int(d.Languages.ProductionPct * 0.9)
		kBar := int((cats[0].ProgressPct / 100.0) * 90.0)
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 465 255" width="465" height="255">
  <style>
    .sans { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-variant-numeric: tabular-nums; }
    .hdr { font-size: 12px; font-weight: 700; fill: #526071; text-transform: uppercase; letter-spacing: 0.5px; }
    .sec { font-size: 9px; font-weight: 700; fill: #0284c7; text-transform: uppercase; letter-spacing: 0.4px; }
    .lbl { font-size: 10px; font-weight: 500; fill: #526071; }
    .val { font-size: 10px; font-weight: 700; fill: #64748b; }
    .kpi-num { font-size: 13px; font-weight: 700; fill: #526071; }
    .kpi-lbl { font-size: 8px; font-weight: 600; fill: #64748b; text-transform: uppercase; }
  </style>
  <rect x="2" y="2" width="461" height="251" rx="6" fill="none" stroke="#64748b" stroke-opacity="0.28" stroke-width="1"/>
  <text x="14" y="20" class="sans hdr">%s — %s</text>
  <text x="451" y="20" text-anchor="end" class="mono val">%d Repos · ★ %s</text>
  <line x1="2" y1="27" x2="463" y2="27" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>

  <!-- Top KPI Strip -->
  <g transform="translate(14, 33)">
    <text x="0" y="11" class="mono kpi-num" fill="#0284c7">%s %.1f%%</text>
    <text x="0" y="22" class="sans kpi-lbl">%s</text>
    <line x1="145" y1="0" x2="145" y2="25" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="157" y="11" class="mono kpi-num" fill="#d97706">%s (#%s)</text>
    <text x="157" y="22" class="sans kpi-lbl">%s</text>
    <line x1="305" y1="0" x2="305" y2="25" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="317" y="11" class="mono kpi-num" fill="#059669">%s (+%.0f%%)</text>
    <text x="317" y="22" class="sans kpi-lbl">%s</text>
  </g>
  <line x1="2" y1="63" x2="463" y2="63" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>

  <!-- Q1: Languages (Bullet Graphs) -->
  <g transform="translate(14, 78)">
    <text x="0" y="0" class="sans sec">1. %s</text>
    <text x="0" y="18" class="sans lbl">%s</text>
    <rect x="62" y="10" width="90" height="9" fill="#64748b" fill-opacity="0.14"/>
    <rect x="62" y="10" width="64" height="9" fill="#64748b" fill-opacity="0.25"/>
    <rect x="62" y="12" width="%d" height="5" fill="#0284c7"/>
    <line x1="126" y1="8" x2="126" y2="21" stroke="#0284c7" stroke-width="2"/>
    <text x="200" y="18" text-anchor="end" class="mono val" fill="#0284c7">%.1f%%</text>

    <text x="0" y="38" class="sans lbl">Prod vs EDA</text>
    <rect x="62" y="30" width="90" height="9" fill="#64748b" fill-opacity="0.14"/>
    <rect x="62" y="30" width="64" height="9" fill="#64748b" fill-opacity="0.25"/>
    <rect x="62" y="32" width="%d" height="5" fill="#059669"/>
    <line x1="120" y1="28" x2="120" y2="41" stroke="#0284c7" stroke-width="2"/>
    <text x="200" y="38" text-anchor="end" class="mono val">%.0f%%</text>

    <text x="0" y="56" class="sans lbl">%s / %s</text>
    <text x="200" y="56" text-anchor="end" class="mono val">%.1f%% / %.1f%%</text>
  </g>

  <!-- Q2: Stack -->
  <g transform="translate(248, 78)">
    <text x="0" y="0" class="sans sec">2. %s</text>
    <text x="0" y="18" class="sans lbl">%s</text><text x="203" y="18" text-anchor="end" class="mono val" fill="#0284c7">%d repos</text>
    <text x="0" y="38" class="sans lbl">%s</text><text x="203" y="38" text-anchor="end" class="mono val" fill="#7c3aed">%d repos</text>
    <text x="0" y="56" class="sans lbl">%s</text><text x="203" y="56" text-anchor="end" class="mono val">%d / %d</text>
  </g>

  <line x1="232" y1="63" x2="232" y2="253" stroke="#64748b" stroke-opacity="0.22" stroke-width="1"/>
  <line x1="2" y1="154" x2="463" y2="154" stroke="#64748b" stroke-opacity="0.22" stroke-width="1"/>

  <!-- Q3: Kaggle (Bullet Graph) -->
  <g transform="translate(14, 170)">
    <text x="0" y="0" class="sans sec">3. %s</text>
    <text x="0" y="18" class="sans lbl">%s</text>
    <rect x="62" y="10" width="90" height="9" fill="#64748b" fill-opacity="0.14"/>
    <rect x="62" y="10" width="64" height="9" fill="#64748b" fill-opacity="0.25"/>
    <rect x="62" y="12" width="%d" height="5" fill="#d97706"/>
    <line x1="142" y1="8" x2="142" y2="21" stroke="#0284c7" stroke-width="2"/>
    <text x="200" y="18" text-anchor="end" class="mono val" fill="#d97706">#%s</text>
    <text x="0" y="38" class="sans lbl">%s</text><text x="200" y="38" text-anchor="end" class="mono val">%d (%dG · %dS · %dB)</text>
    <text x="0" y="56" class="sans lbl">Notebooks / Datasets</text><text x="200" y="56" text-anchor="end" class="mono val">#%s / #%s</text>
  </g>

  <!-- Q4: Impact -->
  <g transform="translate(248, 170)">
    <text x="0" y="0" class="sans sec">4. %s</text>
    <text x="0" y="18" class="sans lbl">%s</text><text x="203" y="18" text-anchor="end" class="mono val" fill="#059669">%s (+%.0f%%)</text>
    <text x="0" y="38" class="sans lbl">%s</text><text x="203" y="38" text-anchor="end" class="mono val">★ %s (%d forks)</text>
    <text x="0" y="56" class="sans lbl">%s</text><text x="203" y="56" text-anchor="end" class="mono val">%s (30d)</text>
  </g>
</svg>`,
			esc(d.Subject.Name), esc(d.Subject.Role),
			d.Languages.TotalRepos, formatInt(im.TotalStars),
			esc(langs[0].Name), langs[0].Percentage, esc(l.LangSecTrueProp),
			esc(k.OverallTier), formatInt(k.PeakRank), esc(l.KaggleTitle),
			formatInt(im.AnnualContributions), im.YoYGrowthPct, esc(l.ImpactContribs),
			esc(l.LangSecTrueProp),
			esc(langs[0].Name), l0Bar, langs[0].Percentage,
			prodBar, d.Languages.ProductionPct,
			esc(langs[1].Name), esc(langs[2].Name), langs[1].Percentage, langs[2].Percentage,
			esc(l.StackTitle),
			esc(dataCol[0].Name), dataCol[0].RepoCount,
			esc(mlCol[0].Name), mlCol[0].RepoCount,
			esc(l.LangLockfiles), d.Languages.LockfileRepos, d.Languages.TotalRepos,
			esc(l.KaggleTitle),
			esc(cats[0].Name), kBar, formatInt(cats[0].Rank),
			esc(l.KaggleTotalMedalsLbl), k.TotalMedals(), k.TotalGold, k.TotalSilver, k.TotalBronze,
			formatInt(cats[0].Rank), formatInt(cats[1].Rank),
			esc(l.ImpactTitle),
			esc(l.ImpactContribs), formatInt(im.AnnualContributions), im.YoYGrowthPct,
			esc(l.ImpactStars), formatInt(im.TotalStars), im.TotalForks,
			esc(l.ImpactDownloads), formatCompact(im.HFDownloads30d))

	case domain.SchoolCairo:
		l0End := 62 + int(langs[0].Percentage*1.5)
		l1End := 62 + int(langs[1].Percentage*1.5)
		sumLead := "Executive portfolio synthesis across source code, ecosystem stack, Kaggle standing, and community reach"
		sumCallout := fmt.Sprintf("↳ Notebook-deflated source code (%.1f MB stripped) combined with verified multi-platform telemetry.", d.Languages.StrippedBase64MB)
		if l.IsPtBR {
			sumLead = "Síntese executiva do portfólio em código-fonte, ecossistema, Kaggle e alcance na comunidade"
			sumCallout = fmt.Sprintf("↳ Código-fonte normalizado (%.1f MB removidos) combinado com telemetria multiplataforma verificada.", d.Languages.StrippedBase64MB)
		}
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 465 255" width="465" height="255">
  <style>
    .serif { font-family: Georgia, "Times New Roman", serif; }
    .sans { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-variant-numeric: tabular-nums; }
    .headline { font-size: 12.5px; font-weight: bold; fill: #526071; }
    .lead { font-size: 9.8px; fill: #64748b; }
    .sec { font-size: 9px; font-weight: 700; fill: #e11d48; text-transform: uppercase; letter-spacing: 0.4px; }
    .lbl { font-size: 10px; font-weight: 600; fill: #526071; }
    .val { font-size: 10px; font-weight: 700; fill: #64748b; }
    .callout { font-size: 9px; font-style: italic; fill: #0284c7; }
  </style>
  <rect x="12" y="6" width="441" height="3" fill="#e11d48"/>
  <text x="12" y="24" class="serif headline">%s — %s (%d Repos · ★ %s)</text>
  <text x="12" y="38" class="sans lead">%s</text>
  <line x1="12" y1="45" x2="453" y2="45" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>

  <!-- Q1: Languages (Lollipop + Waffle) -->
  <g transform="translate(12, 60)">
    <text x="0" y="0" class="sans sec">1. %s</text>
    <text x="0" y="18" class="sans lbl">%s</text>
    <line x1="62" y1="14" x2="%d" y2="14" stroke="#0284c7" stroke-width="2.2"/>
    <circle cx="%d" cy="14" r="3.5" fill="#0284c7"/>
    <text x="202" y="18" text-anchor="end" class="mono val" fill="#0284c7">%.1f%%</text>

    <text x="0" y="36" class="sans lbl">%s</text>
    <line x1="62" y1="32" x2="%d" y2="32" stroke="#d97706" stroke-width="2.2"/>
    <circle cx="%d" cy="32" r="3.5" fill="#d97706"/>
    <text x="202" y="36" text-anchor="end" class="mono val">%.1f%%</text>

    <text x="0" y="54" class="sans lbl">%s</text>
    <text x="202" y="54" text-anchor="end" class="mono val" fill="#059669">%.0f%% Prod / %.0f%% EDA</text>
  </g>

  <!-- Q2: Stack (Lollipop + Dashed Benchmark) -->
  <g transform="translate(248, 60)">
    <text x="0" y="0" class="sans sec">2. %s</text>
    <line x1="150" y1="8" x2="150" y2="58" stroke="#e11d48" stroke-dasharray="3,3" stroke-width="1"/>
    <text x="0" y="18" class="sans lbl">%s</text><text x="203" y="18" text-anchor="end" class="mono val">%d repos</text>
    <text x="0" y="36" class="sans lbl">%s</text><text x="203" y="36" text-anchor="end" class="mono val">%d repos</text>
    <text x="0" y="54" class="sans lbl">%s</text><text x="203" y="54" text-anchor="end" class="mono val">%d / %d</text>
  </g>

  <line x1="12" y1="132" x2="453" y2="132" stroke="#64748b" stroke-opacity="0.22" stroke-width="1"/>

  <!-- Q3: Kaggle -->
  <g transform="translate(12, 148)">
    <text x="0" y="0" class="sans sec">3. %s</text>
    <text x="0" y="18" class="sans lbl">%s</text><text x="202" y="18" text-anchor="end" class="mono val" fill="#d97706">%s (#%s)</text>
    <text x="0" y="36" class="sans lbl">%s</text><text x="202" y="36" text-anchor="end" class="mono val">%d (%dG · %dS · %dB)</text>
    <text x="0" y="54" class="sans lbl">Notebooks / Datasets</text><text x="202" y="54" text-anchor="end" class="mono val">#%s / #%s</text>
  </g>

  <!-- Q4: Impact -->
  <g transform="translate(248, 148)">
    <text x="0" y="0" class="sans sec">4. %s</text>
    <text x="0" y="18" class="sans lbl">%s</text><text x="203" y="18" text-anchor="end" class="mono val" fill="#059669">%s (+%.0f%%)</text>
    <text x="0" y="36" class="sans lbl">%s</text><text x="203" y="36" text-anchor="end" class="mono val">★ %s (%d forks)</text>
    <text x="0" y="54" class="sans lbl">%s</text><text x="203" y="54" text-anchor="end" class="mono val">%s (30d)</text>
  </g>

  <text x="12" y="238" class="serif callout">%s</text>
</svg>`,
			esc(d.Subject.Name), esc(d.Subject.Role),
			d.Languages.TotalRepos, formatInt(im.TotalStars),
			esc(sumLead),
			esc(l.LangSecTrueProp),
			esc(langs[0].Name), l0End, l0End, langs[0].Percentage,
			esc(langs[1].Name), l1End, l1End, langs[1].Percentage,
			esc(l.LangSecEDAProd), d.Languages.ProductionPct, d.Languages.ExplorationPct,
			esc(l.StackTitle),
			esc(dataCol[0].Name), dataCol[0].RepoCount,
			esc(mlCol[0].Name), mlCol[0].RepoCount,
			esc(l.LangLockfiles), d.Languages.LockfileRepos, d.Languages.TotalRepos,
			esc(l.KaggleTitle),
			esc(l.KaggleOverallTierLbl), esc(k.OverallTier), formatInt(k.PeakRank),
			esc(l.KaggleTotalMedalsLbl), k.TotalMedals(), k.TotalGold, k.TotalSilver, k.TotalBronze,
			formatInt(cats[0].Rank), formatInt(cats[1].Rank),
			esc(l.ImpactTitle),
			esc(l.ImpactContribs), formatInt(im.AnnualContributions), im.YoYGrowthPct,
			esc(l.ImpactStars), formatInt(im.TotalStars), im.TotalForks,
			esc(l.ImpactDownloads), formatCompact(im.HFDownloads30d),
			esc(sumCallout))

	default: // SchoolHolmes
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 465 255" width="465" height="255">
  <style>
    .sans { font-family: "Trebuchet MS", Inter, -apple-system, BlinkMacSystemFont, sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-weight: 700; }
    .banner { font-size: 12px; font-weight: 800; fill: #ffffff; letter-spacing: 0.6px; text-transform: uppercase; }
    .sec { font-size: 9.5px; font-weight: 800; fill: #d97706; text-transform: uppercase; letter-spacing: 0.4px; }
    .lbl { font-size: 10px; font-weight: 700; fill: #526071; }
    .val { font-size: 10px; font-weight: 800; fill: #0284c7; }
  </style>
  <rect x="3" y="3" width="459" height="249" rx="10" fill="none" stroke="#64748b" stroke-opacity="0.35" stroke-width="2"/>
  <path d="M 3,13 Q 3,3 13,3 L 452,3 Q 462,3 462,13 L 462,32 L 3,32 Z" fill="#0284c7"/>
  <text x="16" y="22" class="sans banner">⚡ %s — %s (%d REPOS)</text>

  <!-- Q1: Languages (Pictorial Unit Blocks) -->
  <g transform="translate(16, 48)">
    <text x="0" y="0" class="sans sec">1. %s</text>
    <text x="0" y="20" class="sans lbl">%s</text>
    <g transform="translate(54, 10)" fill="#0284c7">
      <rect x="0" y="0" width="12" height="10" rx="2"/><rect x="15" y="0" width="12" height="10" rx="2"/><rect x="30" y="0" width="12" height="10" rx="2"/><rect x="45" y="0" width="12" height="10" rx="2"/><rect x="60" y="0" width="12" height="10" rx="2"/>
    </g>
    <text x="200" y="20" text-anchor="end" class="mono val">%.1f%%</text>
    <text x="0" y="38" class="sans lbl">%s / %s</text>
    <text x="200" y="38" text-anchor="end" class="mono val" fill="#d97706">%.1f%% / %.1f%%</text>
    <rect x="0" y="46" width="200" height="22" rx="5" fill="#059669" fill-opacity="0.1" stroke="#059669" stroke-opacity="0.4" stroke-width="1"/>
    <text x="8" y="61" class="sans lbl" fill="#059669">✓ %.0f%% Prod · %.0f%% EDA</text>
  </g>

  <!-- Q2: Stack (Visual Engineering Badges) -->
  <g transform="translate(244, 48)">
    <text x="0" y="0" class="sans sec">2. %s</text>
    <rect x="0" y="8" width="204" height="20" rx="5" fill="#0284c7" fill-opacity="0.1" stroke="#0284c7" stroke-opacity="0.4" stroke-width="1"/>
    <text x="8" y="22" class="sans lbl">🧰 %s</text><text x="196" y="22" text-anchor="end" class="mono val">%d repos</text>
    <rect x="0" y="32" width="204" height="20" rx="5" fill="#7c3aed" fill-opacity="0.1" stroke="#7c3aed" stroke-opacity="0.4" stroke-width="1"/>
    <text x="8" y="46" class="sans lbl">🤖 %s</text><text x="196" y="46" text-anchor="end" class="mono val" fill="#7c3aed">%d repos</text>
    <rect x="0" y="56" width="204" height="20" rx="5" fill="#059669" fill-opacity="0.1" stroke="#059669" stroke-opacity="0.4" stroke-width="1"/>
    <text x="8" y="70" class="sans lbl">🔒 %s</text><text x="196" y="70" text-anchor="end" class="mono val" fill="#059669">%d/%d</text>
  </g>

  <!-- Q3: Kaggle (Stepped Podium + Trophy) -->
  <g transform="translate(16, 146)">
    <text x="0" y="0" class="sans sec">3. %s</text>
    <rect x="0" y="44" width="18" height="16" rx="2" fill="#64748b" fill-opacity="0.25"/>
    <rect x="21" y="34" width="18" height="26" rx="2" fill="#059669" fill-opacity="0.4"/>
    <rect x="42" y="22" width="18" height="38" rx="2" fill="#0284c7" fill-opacity="0.65"/>
    <rect x="63" y="10" width="18" height="50" rx="2" fill="#d97706"/>
    <text x="72" y="8" text-anchor="middle" font-size="10">🏆</text>
    <text x="92" y="26" class="sans lbl">Tier: %s</text>
    <text x="92" y="42" class="mono val" fill="#d97706">#%s · %d Medals</text>
    <text x="92" y="58" class="mono val">%dG · %dS · %dB</text>
  </g>

  <!-- Q4: Impact (Visual Badges) -->
  <g transform="translate(244, 146)">
    <text x="0" y="0" class="sans sec">4. %s</text>
    <rect x="0" y="8" width="204" height="24" rx="5" fill="#059669" fill-opacity="0.1" stroke="#059669" stroke-opacity="0.4" stroke-width="1"/>
    <text x="8" y="24" class="sans lbl">🚀 %s</text><text x="196" y="24" text-anchor="end" class="mono val" fill="#059669">%s (+%.0f%%)</text>
    <rect x="0" y="36" width="204" height="24" rx="5" fill="#d97706" fill-opacity="0.1" stroke="#d97706" stroke-opacity="0.4" stroke-width="1"/>
    <text x="8" y="52" class="sans lbl">⭐ %s</text><text x="196" y="52" text-anchor="end" class="mono val" fill="#d97706">★ %s</text>
    <rect x="0" y="64" width="204" height="24" rx="5" fill="#0284c7" fill-opacity="0.1" stroke="#0284c7" stroke-opacity="0.4" stroke-width="1"/>
    <text x="8" y="80" class="sans lbl">📥 %s</text><text x="196" y="80" text-anchor="end" class="mono val">%s (30d)</text>
  </g>
</svg>`,
			esc(d.Subject.Name), esc(d.Subject.Role), d.Languages.TotalRepos,
			esc(l.LangSecTrueProp),
			esc(langs[0].Name), langs[0].Percentage,
			esc(langs[1].Name), esc(langs[2].Name), langs[1].Percentage, langs[2].Percentage,
			d.Languages.ProductionPct, d.Languages.ExplorationPct,
			esc(l.StackTitle),
			esc(dataCol[0].Name), dataCol[0].RepoCount,
			esc(mlCol[0].Name), mlCol[0].RepoCount,
			esc(l.LangLockfiles), d.Languages.LockfileRepos, d.Languages.TotalRepos,
			esc(l.KaggleTitle),
			esc(k.OverallTier), formatInt(k.PeakRank), k.TotalMedals(),
			k.TotalGold, k.TotalSilver, k.TotalBronze,
			esc(l.ImpactTitle),
			esc(l.ImpactContribs), formatInt(im.AnnualContributions), im.YoYGrowthPct,
			esc(l.ImpactStars), formatInt(im.TotalStars),
			esc(l.ImpactDownloads), formatCompact(im.HFDownloads30d))
	}
}

func buildTrendPolyline(weekly []int, x0, y0, width, height int) (string, int, int) {
	if len(weekly) < 2 {
		weekly = []int{18, 22, 19, 25, 28, 24, 31, 35, 29, 38, 42, 46}
	}
	minVal, maxVal := weekly[0], weekly[0]
	for _, v := range weekly[1:] {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}
	span := maxVal - minVal
	if span <= 0 {
		span = 1
	}

	pts := make([]string, len(weekly))
	lastX, lastY := x0, y0
	n := len(weekly) - 1
	for i, v := range weekly {
		x := x0 + (i*width)/n
		y := (y0 + height) - ((v-minVal)*height)/span
		pts[i] = fmt.Sprintf("%d,%d", x, y)
		lastX, lastY = x, y
	}
	return strings.Join(pts, " "), lastX, lastY
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

// ============================================================================
// ANIMATED HEADER, ELASTIC BIO (440px & 900px), AND FULL COMPOSITE PANEL (900px)
// ============================================================================

func nonEmptyStr(s, fallback string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return fallback
	}
	return s
}

func resolveHeadlines(d domain.PortfolioData, l labels) []string {
	var cleaned []string
	for _, h := range d.Subject.Headlines {
		if s := strings.TrimSpace(h); s != "" {
			cleaned = append(cleaned, s)
		}
	}
	if len(cleaned) > 0 {
		return cleaned
	}

	name := nonEmptyStr(d.Subject.Name, "Data Practitioner")
	langs := padLangs(d.Languages.Languages)
	topStack := "Pandas / Polars"
	if len(d.Stack.DataAnalytics) > 0 && d.Stack.DataAnalytics[0].Name != "" {
		topStack = d.Stack.DataAnalytics[0].Name
	}

	if l.IsPtBR {
		return []string{
			fmt.Sprintf("Olá! Eu sou %s 👋", name),
			fmt.Sprintf("%s (%.1f%%) · %s · %d Repositórios", langs[0].Name, langs[0].Percentage, topStack, d.Languages.TotalRepos),
			fmt.Sprintf("Kaggle %s (#%s) · %s Contribuições (+%.0f%% YoY)", d.Kaggle.OverallTier, formatInt(d.Kaggle.PeakRank), formatInt(d.Impact.AnnualContributions), d.Impact.YoYGrowthPct),
		}
	}
	return []string{
		fmt.Sprintf("Welcome! I'm %s 👋", name),
		fmt.Sprintf("%s (%.1f%%) · %s · %d Repositories", langs[0].Name, langs[0].Percentage, topStack, d.Languages.TotalRepos),
		fmt.Sprintf("Kaggle %s (#%s) · %s Contributions (+%.0f%% YoY)", d.Kaggle.OverallTier, formatInt(d.Kaggle.PeakRank), formatInt(d.Impact.AnnualContributions), d.Impact.YoYGrowthPct),
	}
}

func buildAnimatedHeadlinesSVG(headlines []string, idPrefix string, x, y, maxW int, textClass, fillColor string, fontSizePx float64) string {
	if len(headlines) == 0 {
		return ""
	}
	var b bytes.Buffer
	n := len(headlines)
	for i, msg := range headlines {
		pathID := fmt.Sprintf("%spath%d", idPrefix, i)
		animID := fmt.Sprintf("%sd%d", idPrefix, i)
		var beginAttr string
		if n == 1 {
			beginAttr = fmt.Sprintf("0s;%s.end", animID)
		} else if i == 0 {
			beginAttr = fmt.Sprintf("0s;%sd%d.end", idPrefix, n-1)
		} else {
			beginAttr = fmt.Sprintf("%sd%d.end", idPrefix, i-1)
		}
		fmt.Fprintf(&b, `  <path id="%s">
    <animate id="%s" attributeName="d" begin="%s" dur="6000ms" fill="remove" values="m0,%d h0 ; m0,%d h%d ; m0,%d h%d ; m0,%d h0" keyTimes="0;0.666667;0.833333;1"/>
  </path>
  <text x="%d" class="%s" font-weight="bold" fill="%s" font-size="%.1fpx"><textPath xlink:href="#%s">%s</textPath></text>
`, pathID, animID, beginAttr, y, y, maxW, y, maxW, y, x, textClass, fillColor, fontSizePx, pathID, esc(msg))
	}
	return b.String()
}

func buildTrendPathD(weekly []int, x0, y0, width, height int) (string, int, int) {
	ptsStr, lastX, lastY := buildTrendPolyline(weekly, x0, y0, width, height)
	parts := strings.Fields(ptsStr)
	if len(parts) == 0 {
		return fmt.Sprintf("M %d,%d L %d,%d", x0, y0+height, x0+width, y0), x0 + width, y0
	}
	var b strings.Builder
	for i, pt := range parts {
		if i == 0 {
			b.WriteString("M ")
		} else {
			b.WriteString(" L ")
		}
		b.WriteString(pt)
	}
	return b.String(), lastX, lastY
}

func renderHeaderSection(school domain.School, d domain.PortfolioData, l labels, idPrefix string, standalone bool) string {
	name := nonEmptyStr(d.Subject.Name, "Data Practitioner")
	role := nonEmptyStr(d.Subject.Role, "Data Scientist & Analytics Engineer")
	headlines := resolveHeadlines(d, l)
	langs := padLangs(d.Languages.Languages)
	k := d.Kaggle
	im := d.Impact

	switch school {
	case domain.SchoolTufte:
		pathD, lastX, lastY := buildTrendPathD(im.WeeklyTrend, 596, 36, 280, 46)
		animText := buildAnimatedHeadlinesSVG(headlines, idPrefix, 14, 78, 555, "serif", "#526071", 21.0)
		rightSub := fmt.Sprintf("%d Repos · ★ %s · Kaggle %s", d.Languages.TotalRepos, formatInt(im.TotalStars), k.OverallTier)
		inner := fmt.Sprintf(`  <text x="14" y="24" class="serif h1">%s — <tspan class="sub">%s</tspan></text>
  <text x="886" y="24" text-anchor="end" class="num sub">%s</text>
  <line x1="14" y1="31" x2="886" y2="31" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
%s  <rect x="596" y="42" width="280" height="32" fill="#64748b" fill-opacity="0.12"/>
  <path d="%s" fill="none" stroke="#64748b" stroke-dasharray="2.5,2.5" stroke-width="1"/>
  <path d="%s" fill="none" stroke="#0284c7" stroke-width="2">
    <animate attributeName="stroke-dasharray" values="0 100%%; 10%% 100%%; 100%% 10%%; 10%% 100%%" keyTimes="0; 0.5; 0.9; 1" dur="4s" begin="0s" repeatCount="indefinite"/>
  </path>
  <circle cx="%d" cy="%d" r="3.2" fill="#059669"/>
  <text x="886" y="96" text-anchor="end" class="num micro">%s · %s (+%.0f%% YoY)</text>
  <line x1="14" y1="106" x2="886" y2="106" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>`,
			esc(name), esc(role), esc(rightSub),
			animText, pathD, pathD, lastX, lastY,
			esc(l.HeaderTrendLbl), formatInt(im.AnnualContributions), im.YoYGrowthPct)
		if !standalone {
			return inner
		}
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 900 115" width="900" height="115">
  <style>
    .serif { font-family: Georgia, "Palatino Linotype", "Book Antiqua", Palatino, serif; }
    .num { font-family: Georgia, serif; font-variant-numeric: tabular-nums; }
    .h1 { font-size: 13.5px; font-weight: bold; fill: #526071; letter-spacing: 0.2px; }
    .sub { font-size: 11px; font-style: italic; fill: #64748b; }
    .micro { font-size: 9.5px; fill: #64748b; }
  </style>
%s
</svg>`, inner)

	case domain.SchoolFew:
		pathD, lastX, lastY := buildTrendPathD(im.WeeklyTrend, 600, 44, 276, 38)
		animText := buildAnimatedHeadlinesSVG(headlines, idPrefix, 16, 86, 555, "sans", "#526071", 20.0)
		bulletW := clampInt(int(math.Round((d.Languages.ProductionPct/100.0)*120.0)), 18, 120)
		inner := fmt.Sprintf(`  <text x="16" y="21" class="sans hdr">%s — %s</text>
  <text x="884" y="21" text-anchor="end" class="mono tag">+%.0f%% YoY · %d REPOS</text>
  <line x1="2" y1="29" x2="898" y2="29" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(16, 34)">
    <text x="0" y="11" class="mono kpi-num" fill="#0284c7">%s %.1f%%</text>
    <text x="0" y="21" class="sans kpi-lbl">%s</text>
    <line x1="150" y1="0" x2="150" y2="23" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="162" y="11" class="mono kpi-num" fill="#d97706">%s (#%s)</text>
    <text x="162" y="21" class="sans kpi-lbl">%s</text>
    <line x1="330" y1="0" x2="330" y2="23" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="342" y="11" class="mono kpi-num" fill="#059669">%s</text>
    <text x="342" y="21" class="sans kpi-lbl">%s</text>
  </g>
%s  <line x1="584" y1="29" x2="584" y2="113" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <path d="%s" fill="none" stroke="#64748b" stroke-dasharray="3,2" stroke-width="1.1"/>
  <path d="%s" fill="none" stroke="#059669" stroke-width="2.3">
    <animate attributeName="stroke-dasharray" values="0 100%%; 10%% 100%%; 100%% 10%%; 10%% 100%%" keyTimes="0; 0.5; 0.9; 1" dur="4s" begin="0s" repeatCount="indefinite"/>
  </path>
  <circle cx="%d" cy="%d" r="3" fill="#0284c7"/>
  <g transform="translate(600, 92)">
    <rect x="0" y="1" width="120" height="8" rx="1.5" fill="#64748b" fill-opacity="0.14"/>
    <rect x="0" y="1" width="84" height="8" rx="1.5" fill="#64748b" fill-opacity="0.25"/>
    <rect x="0" y="3" width="%d" height="4" rx="1" fill="#0284c7"/>
    <line x1="96" y1="-1" x2="96" y2="11" stroke="#0284c7" stroke-width="2"/>
    <text x="276" y="8" text-anchor="end" class="mono kpi-lbl">%s</text>
  </g>`,
			esc(name), esc(role), im.YoYGrowthPct, d.Languages.TotalRepos,
			esc(langs[0].Name), langs[0].Percentage, esc(l.LangFewPrimary),
			esc(k.OverallTier), formatInt(k.PeakRank), esc(l.KaggleOverallTierLbl),
			formatInt(im.AnnualContributions), esc(l.ImpactContribs),
			animText, pathD, pathD, lastX, lastY, bulletW, esc(l.HeaderTrendLbl))
		if !standalone {
			return inner + "\n  <line x1=\"2\" y1=\"113\" x2=\"898\" y2=\"113\" stroke=\"#64748b\" stroke-opacity=\"0.25\" stroke-width=\"1\"/>"
		}
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 900 115" width="900" height="115">
  <style>
    .sans { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-variant-numeric: tabular-nums; }
    .hdr { font-size: 12px; font-weight: 700; fill: #526071; text-transform: uppercase; letter-spacing: 0.5px; }
    .tag { font-size: 9.5px; font-weight: 700; fill: #0284c7; }
    .kpi-num { font-size: 12px; font-weight: 700; fill: #526071; }
    .kpi-lbl { font-size: 8px; font-weight: 600; fill: #64748b; text-transform: uppercase; }
  </style>
  <rect x="2" y="2" width="896" height="111" rx="6" fill="none" stroke="#64748b" stroke-opacity="0.28" stroke-width="1"/>
%s
</svg>`, inner)

	case domain.SchoolCairo:
		pathD, lastX, lastY := buildTrendPathD(im.WeeklyTrend, 602, 26, 272, 56)
		animText := buildAnimatedHeadlinesSVG(headlines, idPrefix, 14, 84, 555, "serif", "#526071", 21.0)
		callout := fmt.Sprintf(l.HeaderCairoCalloutFmt, formatInt(im.AnnualContributions), im.YoYGrowthPct, k.OverallTier)
		inner := fmt.Sprintf(`  <defs>
    <marker id="%sarrow" overflow="visible" markerHeight="4" markerWidth="6" orient="auto-start-reverse" viewBox="0 0 8.75 5">
      <path transform="scale(-.5)" d="m0 0 5-5-17.5 5 17.5 5z" fill="#64748b" fill-rule="evenodd"/>
    </marker>
  </defs>
  <rect x="12" y="6" width="876" height="3" fill="#e11d48"/>
  <text x="14" y="26" class="serif headline">%s</text>
  <text x="14" y="41" class="sans lead">%s</text>
  <line x1="14" y1="48" x2="572" y2="48" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
%s  <path d="M 598,20 V 86 H 880" fill="none" marker-start="url(#%sarrow)" marker-end="url(#%sarrow)" stroke="#64748b" stroke-width="1.2"/>
  <line x1="598" y1="52" x2="880" y2="52" stroke="#e11d48" stroke-dasharray="3,3" stroke-width="1"/>
  <path d="%s" fill="none" stroke="#64748b" stroke-dasharray="3,2" stroke-width="1.1"/>
  <path d="%s" fill="none" stroke="#e11d48" stroke-width="2.4">
    <animate attributeName="stroke-dasharray" values="0 100%%; 10%% 100%%; 100%% 10%%; 10%% 100%%" keyTimes="0; 0.5; 0.9; 1" dur="4s" begin="0s" repeatCount="indefinite"/>
  </path>
  <circle cx="%d" cy="%d" r="3.5" fill="#e11d48"/>
  <text x="598" y="102" class="serif callout">%s</text>`,
			idPrefix, esc(name), esc(role), animText, idPrefix, idPrefix, pathD, pathD, lastX, lastY, esc(callout))
		if !standalone {
			return inner + "\n  <line x1=\"12\" y1=\"112\" x2=\"888\" y2=\"112\" stroke=\"#64748b\" stroke-opacity=\"0.25\" stroke-width=\"1\"/>"
		}
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 900 115" width="900" height="115">
  <style>
    .serif { font-family: Georgia, "Times New Roman", serif; }
    .sans { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; }
    .headline { font-size: 13.5px; font-weight: bold; fill: #526071; }
    .lead { font-size: 10.5px; fill: #64748b; }
    .callout { font-size: 9.2px; font-style: italic; fill: #0284c7; }
  </style>
%s
</svg>`, inner)

	default: // SchoolHolmes
		pathD, lastX, lastY := buildTrendPathD(im.WeeklyTrend, 606, 44, 268, 42)
		animText := buildAnimatedHeadlinesSVG(headlines, idPrefix, 18, 82, 555, "sans", "#526071", 21.0)
		inner := fmt.Sprintf(`  <path d="M 3,13 Q 3,3 13,3 L 887,3 Q 897,3 897,13 L 897,32 L 3,32 Z" fill="#0284c7"/>
  <text x="16" y="22" class="sans banner">⚡ %s — %s</text>
  <text x="884" y="22" text-anchor="end" class="sans banner">🏆 %s · ★ %s</text>
%s  <g transform="translate(602, 42)">
    <rect x="0" y="34" width="24" height="16" rx="2" fill="#64748b" fill-opacity="0.22"/>
    <rect x="28" y="24" width="24" height="26" rx="2" fill="#059669" fill-opacity="0.32"/>
    <rect x="56" y="14" width="24" height="36" rx="2" fill="#0284c7" fill-opacity="0.42"/>
    <rect x="84" y="4" width="24" height="46" rx="2" fill="#d97706" fill-opacity="0.45"/>
  </g>
  <path d="%s" fill="none" stroke="#64748b" stroke-dasharray="3,2" stroke-width="1.2"/>
  <path d="%s" fill="none" stroke="#d97706" stroke-width="2.6">
    <animate attributeName="stroke-dasharray" values="0 100%%; 10%% 100%%; 100%% 10%%; 10%% 100%%" keyTimes="0; 0.5; 0.9; 1" dur="4s" begin="0s" repeatCount="indefinite"/>
  </path>
  <circle cx="%d" cy="%d" r="3.8" fill="#059669"/>
  <text x="880" y="103" text-anchor="end" class="mono val">🚀 %s (+%.0f%% YoY)</text>`,
			esc(name), esc(role), strings.ToUpper(esc(k.OverallTier)), formatInt(im.TotalStars),
			animText, pathD, pathD, lastX, lastY, formatInt(im.AnnualContributions), im.YoYGrowthPct)
		if !standalone {
			return inner + "\n  <line x1=\"16\" y1=\"112\" x2=\"884\" y2=\"112\" stroke=\"#64748b\" stroke-opacity=\"0.28\" stroke-width=\"1.2\"/>"
		}
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 900 115" width="900" height="115">
  <style>
    .sans { font-family: "Trebuchet MS", Inter, -apple-system, BlinkMacSystemFont, sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-weight: 700; }
    .banner { font-size: 11.5px; font-weight: 800; fill: #ffffff; letter-spacing: 0.5px; text-transform: uppercase; }
    .val { font-size: 9.5px; font-weight: 800; fill: #059669; }
  </style>
  <rect x="3" y="3" width="894" height="109" rx="10" fill="none" stroke="#64748b" stroke-opacity="0.35" stroke-width="2"/>
%s
</svg>`, inner)
	}
}

func fallbackBioParagraphs(d domain.PortfolioData, l labels) []string {
	name := nonEmptyStr(d.Subject.Name, "Data Practitioner")
	role := nonEmptyStr(d.Subject.Role, "Data Scientist & Analytics Engineer")
	langs := padLangs(d.Languages.Languages)
	topData := "Pandas / Polars"
	if len(d.Stack.DataAnalytics) > 0 && d.Stack.DataAnalytics[0].Name != "" {
		topData = d.Stack.DataAnalytics[0].Name
	}
	topML := "Scikit-Learn"
	if len(d.Stack.MLAndAI) > 0 && d.Stack.MLAndAI[0].Name != "" {
		topML = d.Stack.MLAndAI[0].Name
	}
	tier := nonEmptyStr(d.Kaggle.OverallTier, domain.KaggleTierContributor)
	kaggleClause := fmt.Sprintf(l.BioFallbackKaggleFmt, tier)
	if d.Kaggle.TotalMedals() > 0 {
		kaggleClause = fmt.Sprintf(l.BioFallbackKaggleMedalsFmt, tier, d.Kaggle.TotalMedals())
	}
	return []string{
		fmt.Sprintf(l.BioFallbackP1Fmt, name, role, d.Languages.TotalRepos),
		fmt.Sprintf(l.BioFallbackP2Fmt, langs[0].Name, langs[0].Percentage, langs[1].Name, langs[1].Percentage, langs[2].Name, topData, topML, kaggleClause, formatInt(d.Impact.AnnualContributions)),
	}
}

func wrapBioParagraphs(paragraphs []string, maxCharsPerLine int) [][]string {
	if maxCharsPerLine < 20 {
		maxCharsPerLine = 20
	}
	var out [][]string
	for _, rawPara := range paragraphs {
		for _, subLine := range strings.Split(strings.ReplaceAll(rawPara, `\n`, "\n"), "\n") {
			words := strings.Fields(subLine)
			if len(words) == 0 {
				continue
			}
			var lines []string
			cur := words[0]
			for _, w := range words[1:] {
				if utf8.RuneCountInString(cur)+1+utf8.RuneCountInString(w) <= maxCharsPerLine {
					cur += " " + w
				} else {
					lines = append(lines, cur)
					cur = w
				}
			}
			lines = append(lines, cur)
			out = append(out, lines)
		}
	}
	return out
}

func computeBioLayout(wrapped [][]string, lineHeight, paraGap int) (totalLines int, bodyHeight int) {
	for _, p := range wrapped {
		totalLines += len(p)
	}
	if totalLines == 0 {
		return 0, 0
	}
	bodyHeight = totalLines * lineHeight
	if len(wrapped) > 1 {
		bodyHeight += (len(wrapped) - 1) * paraGap
	}
	return totalLines, bodyHeight
}

func renderBioLinesSVG(wrapped [][]string, x, startY, lineHeight, paraGap int, textClass string) string {
	var b bytes.Buffer
	y := startY
	for pIdx, para := range wrapped {
		for _, line := range para {
			fmt.Fprintf(&b, `  <text x="%d" y="%d" class="%s">%s</text>`+"\n", x, y, textClass, esc(line))
			y += lineHeight
		}
		if pIdx < len(wrapped)-1 {
			y += paraGap
		}
	}
	return b.String()
}

func renderBio(school domain.School, d domain.PortfolioData, l labels, fullWidth bool) string {
	width := 440
	maxChars := 62
	if fullWidth {
		width = 900
		maxChars = 126
	}

	rawParas := d.Subject.Bio
	if len(rawParas) == 0 {
		rawParas = fallbackBioParagraphs(d, l)
	}
	wrapped := wrapBioParagraphs(rawParas, maxChars)
	if len(wrapped) == 0 {
		wrapped = wrapBioParagraphs(fallbackBioParagraphs(d, l), maxChars)
	}

	const (
		lineHeight = 17
		paraGap    = 8
	)
	_, bodyHeight := computeBioLayout(wrapped, lineHeight, paraGap)

	name := nonEmptyStr(d.Subject.Name, "Data Practitioner")
	role := nonEmptyStr(d.Subject.Role, "Data Scientist & Analytics Engineer")
	langs := padLangs(d.Languages.Languages)
	k := d.Kaggle
	im := d.Impact
	subRight := fmt.Sprintf(l.BioSubFmt, d.Languages.TotalRepos, formatInt(im.TotalStars))

	switch school {
	case domain.SchoolTufte:
		startY := 48
		totalHeight := startY + bodyHeight + 16
		linesSVG := renderBioLinesSVG(wrapped, 14, startY, lineHeight, paraGap, "serif body")
		headerTitle := fmt.Sprintf("%s — %s", name, role)
		if !fullWidth && utf8.RuneCountInString(headerTitle) > 42 {
			headerTitle = name
		}
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">
  <style>
    .serif { font-family: Georgia, "Palatino Linotype", "Book Antiqua", Palatino, serif; }
    .num { font-family: Georgia, serif; font-variant-numeric: tabular-nums; }
    .h1 { font-size: 13.5px; font-weight: bold; fill: #526071; letter-spacing: 0.2px; }
    .sub { font-size: 10px; font-style: italic; fill: #64748b; }
    .body { font-size: 11px; fill: #526071; }
  </style>
  <text x="14" y="21" class="serif h1">%s</text>
  <text x="%d" y="21" text-anchor="end" class="num sub">%s</text>
  <line x1="14" y1="28" x2="%d" y2="28" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
%s  <line x1="14" y1="%d" x2="%d" y2="%d" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
</svg>`, width, totalHeight, width, totalHeight, esc(headerTitle), width-14, esc(subRight), width-14, linesSVG, totalHeight-8, width-14, totalHeight-8)

	case domain.SchoolFew:
		startY := 82
		totalHeight := startY + bodyHeight + 14
		linesSVG := renderBioLinesSVG(wrapped, 14, startY, lineHeight, paraGap, "sans body")
		prodBar := clampInt(int(math.Round((d.Languages.ProductionPct/100.0)*84.0)), 14, 84)
		tierKPI := esc(nonEmptyStr(k.OverallTier, domain.KaggleTierContributor))
		if k.PeakRank > 0 {
			tierKPI = fmt.Sprintf("%s (#%s)", tierKPI, formatInt(k.PeakRank))
		} else if k.TotalMedals() > 0 {
			tierKPI = fmt.Sprintf("%s (%dM)", tierKPI, k.TotalMedals())
		}
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">
  <style>
    .sans { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-variant-numeric: tabular-nums; }
    .hdr { font-size: 12px; font-weight: 700; fill: #526071; text-transform: uppercase; letter-spacing: 0.5px; }
    .tag { font-size: 9.5px; font-weight: 600; fill: #0284c7; }
    .body { font-size: 10.8px; font-weight: 500; fill: #526071; }
    .kpi-num { font-size: 12.5px; font-weight: 700; fill: #526071; }
    .kpi-lbl { font-size: 8px; font-weight: 600; fill: #64748b; text-transform: uppercase; }
  </style>
  <rect x="2" y="2" width="%d" height="%d" rx="6" fill="none" stroke="#64748b" stroke-opacity="0.28" stroke-width="1"/>
  <text x="14" y="20" class="sans hdr">%s</text>
  <text x="%d" y="20" text-anchor="end" class="mono tag">%s</text>
  <line x1="2" y1="28" x2="%d" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(14, 33)">
    <text x="0" y="11" class="mono kpi-num" fill="#0284c7">%s (%.1f%%)</text>
    <text x="0" y="22" class="sans kpi-lbl">%s</text>
    <line x1="138" y1="0" x2="138" y2="25" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="150" y="11" class="mono kpi-num" fill="#d97706">%s</text>
    <text x="150" y="22" class="sans kpi-lbl">%s</text>
    <line x1="292" y1="0" x2="292" y2="25" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <rect x="304" y="3" width="84" height="8" rx="1.5" fill="#64748b" fill-opacity="0.14"/>
    <rect x="304" y="3" width="60" height="8" rx="1.5" fill="#64748b" fill-opacity="0.25"/>
    <rect x="304" y="5" width="%d" height="4" rx="1" fill="#059669"/>
    <line x1="370" y1="1" x2="370" y2="13" stroke="#0284c7" stroke-width="2"/>
    <text x="304" y="22" class="sans kpi-lbl">%s</text>
  </g>
  <line x1="2" y1="63" x2="%d" y2="63" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
%s</svg>`,
			width, totalHeight, width, totalHeight,
			width-4, totalHeight-4,
			esc(name), width-14, esc(l.BioTitle), width-2,
			esc(langs[0].Name), langs[0].Percentage, esc(l.LangFewPrimary),
			tierKPI, esc(l.KaggleOverallTierLbl),
			prodBar, esc(l.LangFewProdEDA), width-2, linesSVG)

	case domain.SchoolCairo:
		startY := 66
		totalHeight := startY + bodyHeight + 28
		linesSVG := renderBioLinesSVG(wrapped, 20, startY, lineHeight, paraGap, "serif body")
		callout := fmt.Sprintf(l.BioCairoCalloutFmt, d.Languages.TotalRepos, d.Languages.StrippedBase64MB)
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">
  <style>
    .serif { font-family: Georgia, "Times New Roman", serif; }
    .sans { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; }
    .headline { font-size: 13px; font-weight: bold; fill: #526071; }
    .lead { font-size: 10px; fill: #64748b; }
    .body { font-size: 11px; fill: #526071; }
    .callout { font-size: 9px; font-style: italic; fill: #0284c7; }
  </style>
  <rect x="12" y="6" width="%d" height="3" fill="#e11d48"/>
  <text x="12" y="25" class="serif headline">%s — %s</text>
  <text x="12" y="39" class="sans lead">%s · %s</text>
  <line x1="12" y1="46" x2="%d" y2="46" stroke="#e11d48" stroke-dasharray="3,3" stroke-width="1"/>
  <line x1="12" y1="54" x2="12" y2="%d" stroke="#0284c7" stroke-width="2.2"/>
%s  <text x="12" y="%d" class="serif callout">%s</text>
</svg>`,
			width, totalHeight, width, totalHeight,
			width-24, esc(name), esc(l.BioTitle), esc(role), esc(subRight),
			width-12, startY+bodyHeight-10, linesSVG, totalHeight-10, esc(callout))

	default: // SchoolHolmes
		startY := 54
		totalHeight := startY + bodyHeight + 16
		linesSVG := renderBioLinesSVG(wrapped, 16, startY, lineHeight, paraGap, "sans body")
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">
  <style>
    .sans { font-family: "Trebuchet MS", Inter, -apple-system, BlinkMacSystemFont, sans-serif; }
    .banner { font-size: 11.5px; font-weight: 800; fill: #ffffff; letter-spacing: 0.5px; text-transform: uppercase; }
    .body { font-size: 11px; font-weight: 600; fill: #526071; }
  </style>
  <rect x="3" y="3" width="%d" height="%d" rx="10" fill="none" stroke="#64748b" stroke-opacity="0.35" stroke-width="2"/>
  <path d="M 3,13 Q 3,3 13,3 L %d,3 Q %d,3 %d,13 L %d,32 L 3,32 Z" fill="#0284c7"/>
  <text x="16" y="22" class="sans banner">⚡ %s — %s</text>
%s</svg>`,
			width, totalHeight, width, totalHeight,
			width-6, totalHeight-6,
			width-13, width-3, width-3, width-3,
			esc(name), esc(l.BioTitle), linesSVG)
	}
}

func extractInnerPanelBody(svgStr string, tagColorOverride, calloutColorOverride, secColorOverride string) string {
	start := strings.Index(svgStr, "</style>")
	if start != -1 {
		svgStr = svgStr[start+len("</style>"):]
	} else if idx := strings.IndexByte(svgStr, '>'); idx != -1 {
		svgStr = svgStr[idx+1:]
	}
	if end := strings.LastIndex(svgStr, "</svg>"); end != -1 {
		svgStr = svgStr[:end]
	}
	if tagColorOverride != "" {
		svgStr = strings.ReplaceAll(svgStr, `class="mono tag"`, fmt.Sprintf(`class="mono tag" fill="%s"`, tagColorOverride))
	}
	if calloutColorOverride != "" {
		svgStr = strings.ReplaceAll(svgStr, `class="serif callout"`, fmt.Sprintf(`class="serif callout" fill="%s"`, calloutColorOverride))
	}
	if secColorOverride != "" {
		svgStr = strings.ReplaceAll(svgStr, `class="sans sec"`, fmt.Sprintf(`class="sans sec" fill="%s"`, secColorOverride))
	}
	// Add subtle SMIL entrance animation on sparkline polylines inside the full composite view
	svgStr = strings.ReplaceAll(svgStr, `<polyline fill="none"`, `<polyline stroke-dasharray="450" fill="none"`)
	return svgStr
}

func renderFull(school domain.School, d domain.PortfolioData, l labels) string {
	const (
		width        = 900
		headerHeight = 115
		lineHeight   = 17
		paraGap      = 8
	)

	// 1. Dynamic Bio section: collapses to height 0 when Subject.Bio is empty; expands to exact line count when present
	wrappedBio := wrapBioParagraphs(d.Subject.Bio, 124)
	_, bioBodyHeight := computeBioLayout(wrappedBio, lineHeight, paraGap)
	bioSectionHeight := 0
	if bioBodyHeight > 0 {
		bioSectionHeight = bioBodyHeight + 18
	}

	gridY := headerHeight + bioSectionHeight + 6
	row2Y := gridY + 238
	totalHeight := row2Y + 236

	headerInner := renderHeaderSection(school, d, l, "full-hdr-", false)
	var bioInner string
	if bioSectionHeight > 0 {
		bioStartY := headerHeight + 18
		bioLines := renderBioLinesSVG(wrappedBio, 18, bioStartY, lineHeight, paraGap, "body-txt")
		bioSepY := headerHeight + bioSectionHeight
		switch school {
		case domain.SchoolTufte:
			bioInner = fmt.Sprintf("%s  <line x1=\"14\" y1=\"%d\" x2=\"886\" y2=\"%d\" stroke=\"#64748b\" stroke-opacity=\"0.4\" stroke-width=\"0.6\"/>\n", bioLines, bioSepY, bioSepY)
		case domain.SchoolFew:
			bioInner = fmt.Sprintf("%s  <line x1=\"2\" y1=\"%d\" x2=\"898\" y2=\"%d\" stroke=\"#64748b\" stroke-opacity=\"0.25\" stroke-width=\"1\"/>\n", bioLines, bioSepY, bioSepY)
		case domain.SchoolCairo:
			bioInner = fmt.Sprintf("  <line x1=\"12\" y1=\"%d\" x2=\"12\" y2=\"%d\" stroke=\"#e11d48\" stroke-width=\"2.2\"/>\n%s  <line x1=\"12\" y1=\"%d\" x2=\"888\" y2=\"%d\" stroke=\"#64748b\" stroke-opacity=\"0.25\" stroke-width=\"1\"/>\n", headerHeight+6, bioSepY-6, bioLines, bioSepY, bioSepY)
		default: // SchoolHolmes
			bioInner = fmt.Sprintf("  <rect x=\"14\" y=\"%d\" width=\"872\" height=\"%d\" rx=\"8\" fill=\"#0284c7\" fill-opacity=\"0.06\" stroke=\"#0284c7\" stroke-opacity=\"0.3\" stroke-width=\"1\"/>\n%s", headerHeight+3, bioSectionHeight-4, bioLines)
		}
	}

	p1 := extractInnerPanelBody(renderLanguages(school, d, l), "#0284c7", "#0284c7", "")
	p2 := extractInnerPanelBody(renderStack(school, d, l), "#64748b", "#0284c7", "")
	p3 := extractInnerPanelBody(renderKaggle(school, d, l), "#d97706", "#d97706", "")
	p4 := extractInnerPanelBody(renderImpact(school, d, l), "#059669", "#059669", "#059669")

	quadrants := fmt.Sprintf(`  <g transform="translate(8, %d)">%s  </g>
  <g transform="translate(452, %d)">%s  </g>
  <g transform="translate(8, %d)">%s  </g>
  <g transform="translate(452, %d)">%s  </g>`, gridY, p1, gridY, p2, row2Y, p3, row2Y, p4)

	switch school {
	case domain.SchoolTufte:
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 %d %d" width="%d" height="%d">
  <style>
    .serif { font-family: Georgia, "Palatino Linotype", "Book Antiqua", Palatino, serif; }
    .num { font-family: Georgia, serif; font-variant-numeric: tabular-nums; }
    .h1 { font-size: 13.5px; font-weight: bold; fill: #526071; letter-spacing: 0.2px; }
    .sub { font-size: 10px; font-style: italic; fill: #64748b; }
    .sec { font-size: 10px; font-style: italic; font-weight: bold; fill: #0284c7; }
    .th { font-size: 9.5px; font-style: italic; font-weight: bold; fill: #0284c7; }
    .lbl { font-size: 10.5px; fill: #526071; }
    .val { font-size: 10.5px; font-weight: bold; fill: #526071; }
    .micro { font-size: 9.2px; fill: #64748b; }
    .big-num { font-size: 18px; font-weight: bold; fill: #526071; }
    .body-txt { font-family: Georgia, "Palatino Linotype", serif; font-size: 11px; fill: #526071; }
  </style>
%s
%s%s
</svg>`, width, totalHeight, width, totalHeight, headerInner, bioInner, quadrants)

	case domain.SchoolFew:
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 %d %d" width="%d" height="%d">
  <style>
    .sans { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-variant-numeric: tabular-nums; }
    .hdr { font-size: 12px; font-weight: 700; fill: #526071; text-transform: uppercase; letter-spacing: 0.5px; }
    .tag { font-size: 9.5px; font-weight: 600; fill: #0284c7; }
    .sec { font-size: 9px; font-weight: 700; fill: #0284c7; text-transform: uppercase; letter-spacing: 0.4px; }
    .lbl { font-size: 10.5px; font-weight: 500; fill: #526071; }
    .val { font-size: 10px; font-weight: 600; fill: #64748b; }
    .kpi-num { font-size: 13.5px; font-weight: 700; fill: #526071; }
    .kpi-lbl { font-size: 8.5px; font-weight: 600; fill: #64748b; text-transform: uppercase; }
    .body-txt { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; font-size: 11px; font-weight: 500; fill: #526071; }
  </style>
  <rect x="2" y="2" width="%d" height="%d" rx="6" fill="none" stroke="#64748b" stroke-opacity="0.28" stroke-width="1"/>
%s
%s%s
</svg>`, width, totalHeight, width, totalHeight, width-4, totalHeight-4, headerInner, bioInner, quadrants)

	case domain.SchoolCairo:
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 %d %d" width="%d" height="%d">
  <style>
    .serif { font-family: Georgia, "Times New Roman", serif; }
    .sans { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-variant-numeric: tabular-nums; }
    .headline { font-size: 12.2px; font-weight: bold; fill: #526071; }
    .lead { font-size: 10px; fill: #64748b; }
    .sec { font-size: 9.5px; font-weight: 700; fill: #0284c7; text-transform: uppercase; letter-spacing: 0.4px; }
    .lbl { font-size: 10.5px; font-weight: 600; fill: #526071; }
    .val { font-size: 10px; font-weight: 600; fill: #64748b; }
    .callout { font-size: 9px; font-style: italic; fill: #0284c7; }
    .body-txt { font-family: Georgia, "Times New Roman", serif; font-size: 11px; fill: #526071; }
  </style>
%s
%s%s
</svg>`, width, totalHeight, width, totalHeight, headerInner, bioInner, quadrants)

	default: // SchoolHolmes
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 %d %d" width="%d" height="%d">
  <style>
    .sans { font-family: "Trebuchet MS", Inter, -apple-system, BlinkMacSystemFont, sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-weight: 700; }
    .banner { font-size: 11.5px; font-weight: 800; fill: #ffffff; letter-spacing: 0.5px; text-transform: uppercase; }
    .sec { font-size: 9.5px; font-weight: 800; fill: #d97706; text-transform: uppercase; letter-spacing: 0.5px; }
    .lbl { font-size: 10.5px; font-weight: 700; fill: #526071; }
    .val { font-size: 10.5px; font-weight: 800; fill: #0284c7; }
    .badge { font-size: 9.5px; font-weight: 700; fill: #526071; }
    .body-txt { font-family: "Trebuchet MS", Inter, sans-serif; font-size: 11px; font-weight: 600; fill: #526071; }
  </style>
  <rect x="3" y="3" width="%d" height="%d" rx="10" fill="none" stroke="#64748b" stroke-opacity="0.35" stroke-width="2"/>
%s
%s%s
</svg>`, width, totalHeight, width, totalHeight, width-6, totalHeight-6, headerInner, bioInner, quadrants)
	}
}

// RenderOutputREADME generates a bilingual (en / pt-BR) showcase README for the orphan output branch,
// presenting all 4 visualization Schools, explaining their visual differences, and providing ready-to-copy
// Markdown/HTML embed snippets interpolated with the Subject's GitHub username.
func RenderOutputREADME(manifest domain.ProfileManifest) string {
	user := strings.TrimSpace(manifest.Sources.GitHub.Username)
	if user == "" {
		user = "your-username"
	}
	repoSlug := user + "/dataprofile"
	name := strings.TrimSpace(manifest.Subject.Name)
	if name == "" {
		name = user
	}
	isPt := strings.EqualFold(strings.TrimSpace(manifest.Subject.Locale), "pt-BR") ||
		strings.EqualFold(strings.TrimSpace(manifest.Subject.Locale), "pt")

	var b bytes.Buffer
	if isPt {
		fmt.Fprintf(&b, `<div align="center">

# 📊 DataProfile — Catálogo de SVGs (%s)

Painéis de telemetria de portfólio gerados automaticamente na branch `+"`output`"+`.
Todos os SVGs utilizam a **Paleta Universal Dual-Surface** (fundo transparente `+"`fill=\"none\"`"+` e contraste WCAG AA em **GitHub Light `+"`#ffffff`"+`** e **GitHub Dark `+"`#0d1117`"+`**) e **zero meta-labels** com o nome do autor/escola dentro da arte.

</div>

---

## 🧭 Diferenças entre as 4 Escolas de Visualização (`+"`School`"+`)

| Escola (`+"`generated/<school>/`"+`) | Referência Teórica | Gramática Visual & Quando Escolher |
| :--- | :--- | :--- |
| **1. `+"`tufte`"+`** | **Edward Tufte** (*The Visual Display of Quantitative Information*, 1983) | **Máxima razão tinta-dados (*data-ink ratio*).** Sem bordas externas ("sem gaiolas"), tipografia clássica serifada com numerais tabulares, *Dot-Dash Plots*, *Supertables* e *Sparklines* com banda interquartil. Ideal para perfis acadêmicos, pesquisa quantitativa e engenharia rigorosa. |
| **2. `+"`few`"+`** | **Stephen Few** (*Information Dashboard Design*, 2006) | **Telemetria executiva de leitura em 3 segundos.** Faixas superiores de KPIs, *Bullet Graphs* (faixas qualitativas de fundo + barra de desempenho + marcador vertical de meta) e cor pré-atentiva restrita aos líderes. Ideal para leitura rápida por recrutadores técnicos e gestores. |
| **3. `+"`cairo`"+`** | **Alberto Cairo** (*The Functional Art*, 2012; *The Truthful Art*, 2016) | **Jornalismo de dados editorial.** Filete superior carmesim (`+"`#e11d48`"+`), manchete serifada orientada a tese, *Waffle Charts 10×5*, *Lollipop Charts* com linha tracejada de benchmark global e camada de anotação (`+"`↳`"+`). Ideal para contar uma história analítica contextualizada. |
| **4. `+"`holmes`"+`** | **Nigel Holmes** (*Designer's Guide to Creating Charts*, 1984) & ISOTYPE | **Storytelling visual e infografia de impacto.** Faixas de cabeçalho estilo revista, blocos pictóricos (*Isotypes*), pódio em degraus rumo a Kaggle Grandmaster com troféu (`+"`🏆`"+`) e badges visuais de engenharia. Ideal para destaque visual imediato. |

---
`, name)
	} else {
		fmt.Fprintf(&b, `<div align="center">

# 📊 DataProfile — Generated SVG Catalog (%s)

Automated portfolio telemetry panels published to the orphan `+"`output`"+` branch.
Every SVG uses the **Universal Dual-Surface Palette** (transparent canvas `+"`fill=\"none\"`"+` with WCAG AA legibility on both **GitHub Light `+"`#ffffff`"+`** and **GitHub Dark `+"`#0d1117`"+`**) and contains **zero meta-labels** mentioning the school or author inside the artwork.

</div>

---

## 🧭 Differences Across the 4 Visualization Schools (`+"`School`"+`)

| School (`+"`generated/<school>/`"+`) | Theoretical Foundation | Visual Grammar & When to Choose |
| :--- | :--- | :--- |
| **1. `+"`tufte`"+`** | **Edward Tufte** (*The Visual Display of Quantitative Information*, 1983) | **Maximum data-ink ratio.** Zero outer borders ("no cages"), classical serif typography with tabular numerals, *Dot-Dash Plots*, *Supertables*, and inline *Sparklines* with interquartile bands. Ideal for quantitative researchers, academics, and systems-minded engineers. |
| **2. `+"`few`"+`** | **Stephen Few** (*Information Dashboard Design*, 2006) | **3-second executive telemetry.** Top KPI strips, *Stephen Few Bullet Graphs* (qualitative background range bands + inner performance bar + vertical target threshold marker), and pre-attentive color strictly on category leaders. Ideal for fast scanning by technical recruiters and hiring managers. |
| **3. `+"`cairo`"+`** | **Alberto Cairo** (*The Functional Art*, 2012; *The Truthful Art*, 2016) | **Editorial data journalism.** Top crimson rule (`+"`#e11d48`"+`), thesis-driven serif headline + lead deck, *10×5 Waffle Charts*, *Lollipop Charts* with dashed global benchmark lines, and explanatory callouts (`+"`↳`"+`). Ideal for narrative, context-rich portfolios. |
| **4. `+"`holmes`"+`** | **Nigel Holmes** (*Designer's Guide to Creating Charts*, 1984) & ISOTYPE | **High-impact pictorial storytelling.** Bold magazine header ribbons, *Pictorial Unit Blocks (Isotypes)*, stepped podium climb toward Kaggle Grandmaster with trophy iconography (`+"`🏆`"+`), and visual engineering badges. Ideal for memorable visual presence. |

---
`, name)
	}

	type schoolDoc struct {
		code    string
		titleEN string
		descEN  string
		titlePT string
		descPT  string
	}

	schools := []schoolDoc{
		{
			code:    "tufte",
			titleEN: "I. `tufte` — Analytical Minimalism & Maximum Data-Ink Ratio (Edward Tufte)",
			descEN:  "Removes all decorative boxes and grid cages so every drop of ink encodes data: **Dot-Dash Plots** (`languages`), **3-Pillar Dot Plots** (`stack`), **Tabular Supertables** (`kaggle`), and **Inline Sparklines** (`impact`).",
			titlePT: "I. `tufte` — Minimalismo Analítico & Razão Tinta-Dados Máxima (Edward Tufte)",
			descPT:  "Elimina caixas e grades decorativas para que cada traço represente dados: **Dot-Dash Plots** (`languages`), **Dot Plots em 3 Pilares** (`stack`), **Supertables Tabulares** (`kaggle`) e **Sparklines Inline** (`impact`).",
		},
		{
			code:    "few",
			titleEN: "II. `few` — Executive Dashboard Telemetry & Bullet Graphs (Stephen Few)",
			descEN:  "Designed for 3-second executive scanning: top **KPI Telemetry Strips**, **Stephen Few Bullet Graphs** comparing current standing against vertical target thresholds (`│`), and pre-attentive color highlighting the #1 tool in each stack pillar.",
			titlePT: "II. `few` — Telemetria Executiva & Bullet Graphs (Stephen Few)",
			descPT:  "Projetado para leitura executiva em 3 segundos: **Faixas Superiores de KPIs**, **Bullet Graphs de Stephen Few** comparando o desempenho atual contra marcadores verticais de meta (`│`) e destaque pré-atentivo apenas nos líderes de cada pilar.",
		},
		{
			code:    "cairo",
			titleEN: "III. `cairo` — Functional Art & Editorial Data Journalism (Alberto Cairo)",
			descEN:  "Frames portfolio telemetry like a data-journalism piece: top **Crimson Editorial Rule (`#e11d48`)**, thesis-driven **Serif Headlines**, **10×5 Waffle Charts**, **Lollipop Charts** with dashed benchmark lines, and **`↳` Callout Annotations**.",
			titlePT: "III. `cairo` — Arte Funcional & Jornalismo de Dados Editorial (Alberto Cairo)",
			descPT:  "Estrutura o portfólio como uma reportagem de dados: **Filete Editorial Carmesim (`#e11d48`)**, **Manchetes Serifadas** com tese analítica, **Waffle Charts 10×5**, **Lollipop Charts** com linha tracejada de referência e **Anotações Explicativas (`↳`)**.",
		},
		{
			code:    "holmes",
			titleEN: "IV. `holmes` — Pictorial Isotypes & Visual Storytelling (Nigel Holmes)",
			descEN:  "Combines bold **Magazine Header Ribbons** with **10-Block Pictorial Isotypes**, a **Stepped Podium Climb** with dynamic trophy placement (`🏆`) for Kaggle progression, and **Visual Engineering Badges**.",
			titlePT: "IV. `holmes` — Isotypes Pictóricos & Storytelling Visual (Nigel Holmes)",
			descPT:  "Combina **Faixas de Cabeçalho de Revista** com **Blocos Pictóricos (Isotypes)**, **Pódio em Degraus** com troféu dinâmico (`🏆`) na progressão do Kaggle e **Badges Visuais de Engenharia**.",
		},
	}

	for _, s := range schools {
		title := s.titleEN
		desc := s.descEN
		copyGridLabel := "📋 **Copy & Paste — 4-Panel Grid (2×2) for your Profile `README.md`:**"
		copySumLabel := "📋 **Copy & Paste — Single Consolidated Card (`summary.svg`):**"
		copyFullLabel := "📋 **Copy & Paste — Full Portfolio Suite (`full.svg` / `header.svg` / `bio.svg` / `bio-full.svg`):**"
		sumPreviewLabel := "Consolidated Executive Card (`summary.svg`):"
		fullPreviewLabel := "Full Portfolio Composite (`full.svg`) & Animated Header (`header.svg`):"
		if isPt {
			title = s.titlePT
			desc = s.descPT
			copyGridLabel = "📋 **Copiar e Colar — Grid 2×2 (4 Painéis) para o seu `README.md` de perfil:**"
			copySumLabel = "📋 **Copiar e Colar — Card Executivo Único (`summary.svg`):**"
			copyFullLabel = "📋 **Copiar e Colar — Portfólio Completo (`full.svg` / `header.svg` / `bio.svg` / `bio-full.svg`):**"
			sumPreviewLabel = "Card Executivo Consolidado (`summary.svg`):"
			fullPreviewLabel = "Portfólio Completo (`full.svg`) & Header Animado (`header.svg`):"
		}

		fmt.Fprintf(&b, `
## %s

%s

<p align="center">
  <img src="generated/%s/header.svg" width="96%%" alt="%s header" />
</p>
<p align="center">
  <img src="generated/%s/bio.svg" width="48%%" alt="%s bio" />
  <img src="generated/%s/bio-full.svg" width="96%%" alt="%s bio-full" />
</p>
<p align="center">
  <img src="generated/%s/languages.svg" width="48%%" alt="%s languages" />
  <img src="generated/%s/stack.svg" width="48%%" alt="%s stack" />
</p>
<p align="center">
  <img src="generated/%s/kaggle.svg" width="48%%" alt="%s kaggle" />
  <img src="generated/%s/impact.svg" width="48%%" alt="%s impact" />
</p>
<p align="center">
  <em>%s</em><br/>
  <img src="generated/%s/summary.svg" width="65%%" alt="%s summary" />
</p>
<p align="center">
  <em>%s</em><br/>
  <img src="generated/%s/full.svg" width="96%%" alt="%s full" />
</p>

%s

`+"```markdown"+`
<p align="center">
  <img src="https://raw.githubusercontent.com/%s/output/generated/%s/languages.svg" width="48%%" alt="Languages & Engineering" />
  <img src="https://raw.githubusercontent.com/%s/output/generated/%s/stack.svg" width="48%%" alt="Data & ML Stack" />
</p>
<p align="center">
  <img src="https://raw.githubusercontent.com/%s/output/generated/%s/kaggle.svg" width="48%%" alt="Kaggle Standing" />
  <img src="https://raw.githubusercontent.com/%s/output/generated/%s/impact.svg" width="48%%" alt="Community Impact" />
</p>
`+"```"+`

%s

`+"```markdown"+`
<p align="center">
  <img src="https://raw.githubusercontent.com/%s/output/generated/%s/summary.svg" width="100%%" alt="DataProfile Executive Summary" />
</p>
`+"```"+`

%s

`+"```markdown"+`
<p align="center">
  <img src="https://raw.githubusercontent.com/%s/output/generated/%s/header.svg" width="100%%" alt="DataProfile Animated Header" />
</p>
<p align="center">
  <img src="https://raw.githubusercontent.com/%s/output/generated/%s/bio.svg" width="48%%" alt="DataProfile Career Bio (Compact)" />
  <img src="https://raw.githubusercontent.com/%s/output/generated/%s/bio-full.svg" width="100%%" alt="DataProfile Career Bio (Full Width)" />
</p>
<p align="center">
  <img src="https://raw.githubusercontent.com/%s/output/generated/%s/full.svg" width="100%%" alt="DataProfile Full Portfolio" />
</p>
`+"```"+`

---
`,
			title,
			desc,
			s.code, s.code,
			s.code, s.code,
			s.code, s.code,
			s.code, s.code,
			s.code, s.code,
			s.code, s.code,
			s.code, s.code,
			sumPreviewLabel,
			s.code, s.code,
			fullPreviewLabel,
			s.code, s.code,
			copyGridLabel,
			repoSlug, s.code,
			repoSlug, s.code,
			repoSlug, s.code,
			repoSlug, s.code,
			copySumLabel,
			repoSlug, s.code,
			copyFullLabel,
			repoSlug, s.code,
			repoSlug, s.code,
			repoSlug, s.code,
			repoSlug, s.code,
		)
	}

	if isPt {
		fmt.Fprintf(&b, `
## 🎛️ Mix & Match (Combinando Escolas Diferentes)

Como nenhum SVG exibe rótulos com o nome da Escola dentro da imagem, você pode combinar livremente painéis de diferentes escolas no seu perfil:

`+"```markdown"+`
<p align="center">
  <img src="https://raw.githubusercontent.com/%s/output/generated/tufte/languages.svg" width="48%%" alt="Languages (Tufte)" />
  <img src="https://raw.githubusercontent.com/%s/output/generated/cairo/stack.svg" width="48%%" alt="Stack (Cairo)" />
</p>
<p align="center">
  <img src="https://raw.githubusercontent.com/%s/output/generated/few/kaggle.svg" width="48%%" alt="Kaggle (Few)" />
  <img src="https://raw.githubusercontent.com/%s/output/generated/holmes/impact.svg" width="48%%" alt="Impact (Holmes)" />
</p>
`+"```"+`

---

<p align="center">
  <sub>Gerado por <a href="https://github.com/pablodiegoo/dataprofile">DataProfile</a> (Apache-2.0). Ao reutilizar ou redistribuir este projeto, mantenha a referência ao repositório original.</sub>
</p>
`, repoSlug, repoSlug, repoSlug, repoSlug)
	} else {
		fmt.Fprintf(&b, `
## 🎛️ Mix & Match Across Schools

Because no SVG contains meta-labels mentioning the School name inside the graphic, you can freely mix and match panels from different schools in your profile `+"`README.md`"+`:

`+"```markdown"+`
<p align="center">
  <img src="https://raw.githubusercontent.com/%s/output/generated/tufte/languages.svg" width="48%%" alt="Languages (Tufte)" />
  <img src="https://raw.githubusercontent.com/%s/output/generated/cairo/stack.svg" width="48%%" alt="Stack (Cairo)" />
</p>
<p align="center">
  <img src="https://raw.githubusercontent.com/%s/output/generated/few/kaggle.svg" width="48%%" alt="Kaggle (Few)" />
  <img src="https://raw.githubusercontent.com/%s/output/generated/holmes/impact.svg" width="48%%" alt="Impact (Holmes)" />
</p>
`+"```"+`

---

<p align="center">
  <sub>Generated by <a href="https://github.com/pablodiegoo/dataprofile">DataProfile</a> under the Apache-2.0 License. Please retain attribution to the original project when reusing or redistributing.</sub>
</p>
`, repoSlug, repoSlug, repoSlug, repoSlug)
	}

	return b.String()
}
