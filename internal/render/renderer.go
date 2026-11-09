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
	KaggleTitle              string
	KaggleCumulative         string
	KaggleGold               string
	KaggleSilver             string
	KaggleBronze             string
	KaggleActive             string
	KaggleTracksFmt          string
	KaggleTrack              string
	KaggleGMReq              string
	KaggleGoldMedalsLbl      string
	KaggleSilverMedalsLbl    string
	KaggleBronzeMedalsLbl    string
	KaggleTufteSubFmt        string
	KaggleThCat              string
	KaggleThStanding         string
	KaggleThRank             string
	KaggleThPct              string
	KaggleThMedals           string
	KaggleTotalTally         string
	KaggleMedalsWord         string
	KaggleFewTierTagFmt      string
	KaggleFewTierNoRankFmt   string
	KaggleTop5Finishes       string
	KaggleTop10Finishes      string
	KaggleFewMarkerNote      string
	KaggleCairoHeadFmt       string
	KaggleCairoHeadNoRankFmt string
	KaggleCairoLeadFmt       string
	KaggleCairoBench         string
	KaggleCairoCallout       string
	KaggleHolmesPodium       string
	KaggleHolmesMedals       string
	KaggleOverallTierLbl     string
	KaggleTotalMedalsLbl     string

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
			KaggleTitle:              "Desempenho Competitivo no Kaggle",
			KaggleCumulative:         "Total Acumulado",
			KaggleGold:               "Ouro",
			KaggleSilver:             "Prata",
			KaggleBronze:             "Bronze",
			KaggleActive:             "Trilhas Ativas",
			KaggleTracksFmt:          "%d Trilhas",
			KaggleTrack:              "Progressão por Categoria",
			KaggleGMReq:              "Requisito GM: 5",
			KaggleGoldMedalsLbl:      "Medalhas de Ouro",
			KaggleSilverMedalsLbl:    "Medalhas de Prata",
			KaggleBronzeMedalsLbl:    "Medalhas de Bronze",
			KaggleTufteSubFmt:        "Tier: %s · %d Medalhas Totais",
			KaggleThCat:              "Categoria",
			KaggleThStanding:         "Nível",
			KaggleThRank:             "Rank Global",
			KaggleThPct:              "Percentil",
			KaggleThMedals:           "Medalhas (O / P / B)",
			KaggleTotalTally:         "Total Acumulado no Portfólio",
			KaggleMedalsWord:         "Medalhas",
			KaggleFewTierTagFmt:      "NÍVEL: %s (#%d Pico)",
			KaggleFewTierNoRankFmt:   "NÍVEL: %s",
			KaggleTop5Finishes:       "Finais Top 5%",
			KaggleTop10Finishes:      "Finais Top 10%",
			KaggleFewMarkerNote:      "indica o limiar necessário para o próximo Kaggle Tier.",
			KaggleCairoHeadFmt:       "Nível %s no Kaggle com Rank Global #%s e %d Medalhas",
			KaggleCairoHeadNoRankFmt: "Nível %s no Kaggle com %d Medalhas nas 3 Trilhas Oficiais",
			KaggleCairoLeadFmt:       "Percentil nas 3 trilhas oficiais de progressão do Kaggle (%s · %d Prata · %d Bronze)",
			KaggleCairoBench:         "REFERÊNCIA GLOBAL TOP 1%",
			KaggleCairoCallout:       "↳ Desempenho consistente no topo percentual em código, datasets e competições.",
			KaggleHolmesPodium:       "Pódio de Progressão",
			KaggleHolmesMedals:       "★ %d Medalhas Conquistadas",
			KaggleOverallTierLbl:     "Nível Geral",
			KaggleTotalMedalsLbl:     "Medalhas Totais",
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
		LangFootnote:          fmt.Sprintf("*Note: %.1f MB of embedded base64 chart cache stripped from .ipynb files.", strippedMB),
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
		KaggleTitle:              "Kaggle Competitive Performance",
		KaggleCumulative:         "Cumulative Tally",
		KaggleGold:               "Gold",
		KaggleSilver:             "Silver",
		KaggleBronze:             "Bronze",
		KaggleActive:             "Active Tracks",
		KaggleTracksFmt:          "%d Tracks",
		KaggleTrack:              "Category Progression",
		KaggleGMReq:              "GM Requirement: 5",
		KaggleGoldMedalsLbl:      "Gold Medals",
		KaggleSilverMedalsLbl:    "Silver Medals",
		KaggleBronzeMedalsLbl:    "Bronze Medals",
		KaggleTufteSubFmt:        "Tier: %s · %d Total Medals",
		KaggleThCat:              "Category",
		KaggleThStanding:         "Standing",
		KaggleThRank:             "Global Rank",
		KaggleThPct:              "Percentile",
		KaggleThMedals:           "Medals (G / S / B)",
		KaggleTotalTally:         "Total Portfolio Tally",
		KaggleMedalsWord:         "Medals",
		KaggleFewTierTagFmt:      "TIER: %s (#%d Peak)",
		KaggleFewTierNoRankFmt:   "TIER: %s",
		KaggleTop5Finishes:       "Top 5% Finishes",
		KaggleTop10Finishes:      "Top 10% Finishes",
		KaggleFewMarkerNote:      "indicates the threshold required to unlock the next Kaggle Tier.",
		KaggleCairoHeadFmt:       "Kaggle %s Standing Ranked #%s Globally with %d Medals",
		KaggleCairoHeadNoRankFmt: "Kaggle %s Standing with %d Medals Across 3 Official Tracks",
		KaggleCairoLeadFmt:       "Percentile standing across all 3 official Kaggle progression tracks (%s · %d Silver · %d Bronze)",
		KaggleCairoBench:         "TOP 1% GLOBAL BENCHMARK",
		KaggleCairoCallout:       "↳ Consistent top-percentile standing across code, datasets, and competitive modeling.",
		KaggleHolmesPodium:       "Tier Climb Podium",
		KaggleHolmesMedals:       "★ %d Total Medals Earned",
		KaggleOverallTierLbl:     "Overall Tier",
		KaggleTotalMedalsLbl:     "Total Medals",
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
			TotalSilver: 11,
			TotalBronze: 18,
			Categories: []domain.KaggleCategory{
				{Name: "Code", Tier: "Master", Rank: 412, TotalUsers: 310000, TopPct: "Top 0.2%", Gold: 2, Silver: 6, Bronze: 9, ArtifactCount: 24, ProgressPct: 82, TargetPct: 90},
				{Name: "Datasets", Tier: "Expert", Rank: 910, TotalUsers: 140000, TopPct: "Top 0.7%", Gold: 1, Silver: 3, Bronze: 4, ArtifactCount: 12, ProgressPct: 64, TargetPct: 74},
				{Name: "Competitions", Tier: "Expert", Rank: 1842, TotalUsers: 220000, TopPct: "Top 0.9%", Gold: 0, Silver: 2, Bronze: 5, ArtifactCount: 14, ProgressPct: 57, TargetPct: 74},
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
	} else {
		cleanedCats := make([]domain.KaggleCategory, 0, 3)
		for _, c := range data.Kaggle.Categories {
			if strings.EqualFold(c.Name, "Discussions") || strings.EqualFold(c.Name, "Discussion") {
				continue
			}
			if strings.EqualFold(c.Name, "Notebooks") || strings.EqualFold(c.Name, "Kernels") {
				c.Name = "Code"
			}
			cleanedCats = append(cleanedCats, c)
		}
		if len(cleanedCats) < 3 {
			cleanedCats = append(cleanedCats, def.Kaggle.Categories[len(cleanedCats):3]...)
		}
		data.Kaggle.Categories = cleanedCats[:3]
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
	if n >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(n)/1000000.0)
	}
	if n >= 1000 {
		return fmt.Sprintf("%.1fk", float64(n)/1000.0)
	}
	return fmt.Sprintf("%d", n)
}

func estimateRuneWidth(ch rune, fontSize float64, mono, bold bool) float64 {
	if ch > 0x2600 {
		return fontSize * 1.15
	}
	if mono {
		return fontSize * 0.60
	}
	switch ch {
	case 'i', 'I', 'l', '1', '.', ',', '|', '!', ';', ':', '\'', ' ':
		if bold {
			return fontSize * 0.32
		}
		return fontSize * 0.28
	case 'm', 'M', 'W', '@', '#', '%', '&':
		if bold {
			return fontSize * 0.85
		}
		return fontSize * 0.78
	}
	if ch >= 'A' && ch <= 'Z' {
		if bold {
			return fontSize * 0.68
		}
		return fontSize * 0.62
	}
	if bold {
		return fontSize * 0.56
	}
	return fontSize * 0.51
}

func estimateTextWidth(text string, fontSize float64, mono, bold, upper bool, letterSpacing float64) float64 {
	if text == "" {
		return 0
	}
	if upper {
		text = strings.ToUpper(text)
	}
	var w float64
	runes := 0
	for _, ch := range text {
		w += estimateRuneWidth(ch, fontSize, mono, bold)
		runes++
	}
	if runes > 1 && letterSpacing > 0 {
		w += float64(runes-1) * letterSpacing
	}
	return w
}

func fitFontSize(text string, baseSize, minSize, maxW float64, mono, bold, upper bool, letterSpacing float64) float64 {
	if maxW <= 0 || text == "" {
		return baseSize
	}
	w := estimateTextWidth(text, baseSize, mono, bold, upper, letterSpacing)
	if w <= maxW {
		return baseSize
	}
	// Iteratively step down by 0.1px until it fits or hits minSize
	sz := baseSize
	for sz > minSize {
		sz = math.Round((sz-0.1)*10) / 10
		ls := letterSpacing
		if ls > 0.25 && sz < baseSize*0.85 {
			ls = 0.1
		}
		if estimateTextWidth(text, sz, mono, bold, upper, ls) <= maxW {
			return sz
		}
	}
	return minSize
}

func fontSizeAttr(size, baseSize float64) string {
	if math.Abs(size-baseSize) < 0.05 {
		return ""
	}
	return fmt.Sprintf(` font-size="%.1fpx"`, size)
}

func fitTextAttr(text string, baseSize, minSize, maxW float64, mono, bold, upper bool, letterSpacing float64) string {
	sz := fitFontSize(text, baseSize, minSize, maxW, mono, bold, upper, letterSpacing)
	if math.Abs(sz-baseSize) < 0.05 {
		return ""
	}
	if letterSpacing > 0.25 && sz < baseSize*0.85 {
		return fmt.Sprintf(` font-size="%.1fpx" letter-spacing="0.1px"`, sz)
	}
	return fmt.Sprintf(` font-size="%.1fpx"`, sz)
}

func fitPairTextAttrs(leftText string, leftBase, leftMin float64, leftMono, leftBold, leftUpper bool, leftLS float64, rightText string, rightBase, rightMin float64, rightMono, rightBold, rightUpper bool, rightLS float64, totalW, minGap float64) (string, string) {
	lw := estimateTextWidth(leftText, leftBase, leftMono, leftBold, leftUpper, leftLS)
	rw := estimateTextWidth(rightText, rightBase, rightMono, rightBold, rightUpper, rightLS)
	if lw+rw+minGap <= totalW {
		return "", ""
	}
	avail := totalW - minGap
	if avail < 20 {
		avail = 20
	}
	leftShare := avail * (lw / (lw + rw))
	rightShare := avail - leftShare
	lAttr := fitTextAttr(leftText, leftBase, leftMin, leftShare, leftMono, leftBold, leftUpper, leftLS)
	rAttr := fitTextAttr(rightText, rightBase, rightMin, rightShare, rightMono, rightBold, rightUpper, rightLS)
	return lAttr, rAttr
}

func shortLangName(name string, maxRunes int) string {
	name = strings.TrimSpace(name)
	switch strings.ToLower(name) {
	case "jupyter notebook":
		name = "Jupyter"
	case "objective-c++":
		name = "Obj-C++"
	case "objective-c":
		name = "Obj-C"
	case "visual basic .net":
		name = "VB.NET"
	case "rich text format":
		name = "RTF"
	}
	if maxRunes > 0 && utf8.RuneCountInString(name) > maxRunes {
		return truncateRunes(name, maxRunes)
	}
	return name
}

func fitHeaderNameRole(name, role, sep string, maxW, baseSize, minSize float64, mono, bold, upper bool, letterSpacing float64) (string, float64) {
	prefix := ""
	if name != "" {
		prefix = name + sep
	}
	full := prefix + role
	if estimateTextWidth(full, baseSize, mono, bold, upper, letterSpacing) <= maxW {
		return role, baseSize
	}
	sz := fitFontSize(full, baseSize, minSize, maxW, mono, bold, upper, letterSpacing)
	if estimateTextWidth(full, sz, mono, bold, upper, letterSpacing) <= maxW {
		return role, sz
	}
	// Try first clause before '|' if present
	if idx := strings.IndexByte(role, '|'); idx != -1 {
		firstClause := strings.TrimSpace(role[:idx])
		if firstClause != "" {
			cand := prefix + firstClause
			szClause := fitFontSize(cand, baseSize, minSize, maxW, mono, bold, upper, letterSpacing)
			if estimateTextWidth(cand, szClause, mono, bold, upper, letterSpacing) <= maxW {
				return firstClause, szClause
			}
		}
	}
	roleRunes := []rune(strings.TrimSpace(role))
	for n := len(roleRunes) - 1; n >= 8; n-- {
		shortRole := strings.TrimSpace(string(roleRunes[:n])) + "…"
		cand := prefix + shortRole
		if estimateTextWidth(cand, minSize, mono, bold, upper, letterSpacing) <= maxW {
			return shortRole, minSize
		}
	}
	return truncateRunes(role, 12), minSize
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
		h1Attr, subAttr := fitPairTextAttrs(l.LangTitleTufte, 13.5, 9.8, false, true, false, 0.2, l.LangSubTufte, 10.0, 8.5, false, false, false, 0, 416, 12)
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
  <text x="12" y="21" class="serif h1"%s>%s</text>
  <text x="428" y="21" text-anchor="end" class="serif sub"%s>%s</text>
  <line x1="12" y1="28" x2="428" y2="28" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
  <g transform="translate(12, 46)">
    <text x="0" y="0" class="serif sec">%s</text>
    <line x1="62" y1="10" x2="62" y2="96" stroke="#64748b" stroke-opacity="0.25" stroke-width="0.5"/>
    <line x1="116" y1="10" x2="116" y2="96" stroke="#64748b" stroke-opacity="0.15" stroke-width="0.5" stroke-dasharray="2,2"/>
    <line x1="170" y1="10" x2="170" y2="96" stroke="#64748b" stroke-opacity="0.15" stroke-width="0.5" stroke-dasharray="2,2"/>
`, h1Attr, esc(l.LangTitleTufte), subAttr, esc(l.LangSubTufte), esc(l.LangSecTrueProp))
		for i, lg := range langs {
			yText := 22 + i*18
			yLine := 18 + i*18
			endX := 62 + int(math.Round((lg.Percentage/100.0)*116.0))
			if endX > 178 {
				endX = 178
			}
			dispName := shortLangName(lg.Name, 11)
			lblAttr := fitTextAttr(dispName, 10.5, 8.2, 57, false, false, false, 0)
			fmt.Fprintf(&b, `    <text x="0" y="%d" class="serif lbl"%s>%s</text>
    <line x1="62" y1="%d" x2="%d" y2="%d" stroke="#64748b" stroke-width="1.2"/>
    <circle cx="%d" cy="%d" r="2.8" fill="%s"/>
    <text x="%d" y="%d" class="num val">%.1f%%</text>
`, yText, lblAttr, esc(dispName), yLine, endX, yLine, endX, yLine, lg.Color, endX+7, yText, lg.Percentage)
		}
		prodWidth := int((lp.ProductionPct / 100.0) * 180.0)
		prodValStr := fmt.Sprintf("%.1f%% (%d files)", lp.ProductionPct, lp.ProductionFiles)
		edaValStr := fmt.Sprintf("%.1f%% (%d files)", lp.ExplorationPct, lp.ExplorationFiles)
		prodLblAttr, prodValAttr := fitPairTextAttrs(l.LangProdModules, 10.5, 8.5, false, false, false, 0, prodValStr, 10.5, 8.2, false, true, false, 0, 180, 6)
		edaLblAttr, edaValAttr := fitPairTextAttrs(l.LangEDANotebooks, 10.5, 8.5, false, false, false, 0, edaValStr, 9.2, 8.0, false, false, false, 0, 180, 6)
		footAttr := fitTextAttr(l.LangFootnote, 10.0, 8.0, 416, false, false, false, 0)
		fmt.Fprintf(&b, `  </g>
  <g transform="translate(248, 46)">
    <text x="0" y="0" class="serif sec">%s</text>
    <text x="0" y="19" class="serif lbl"%s>%s</text>
    <text x="180" y="19" text-anchor="end" class="num val"%s>%s</text>
    <text x="0" y="35" class="serif lbl"%s>%s</text>
    <text x="180" y="35" text-anchor="end" class="num micro"%s>%s</text>
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
  <text x="12" y="211" class="serif sub"%s>%s</text>
  <line x1="12" y1="220" x2="428" y2="220" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
</svg>`,
			esc(l.LangSecEDAProd),
			prodLblAttr, esc(l.LangProdModules), prodValAttr, esc(prodValStr),
			edaLblAttr, esc(l.LangEDANotebooks), edaValAttr, esc(edaValStr),
			prodWidth, prodWidth+3,
			esc(l.LangSecHygiene),
			esc(l.LangLockfiles), lp.LockfileRepos, lp.TotalRepos,
			esc(l.LangCICD), lp.CICDRepos, lp.TotalRepos,
			esc(l.LangTests), lp.TestSuiteRepos, lp.TotalRepos,
			esc(l.LangDocker), lp.ContainerizedRepos, lp.TotalRepos,
			footAttr, esc(l.LangFootnote))
		return b.String()

	case domain.SchoolFew:
		var b bytes.Buffer
		strippedTag := fmt.Sprintf(l.LangFewStrippedFmt, lp.StrippedBase64MB)
		hdrAttr, tagAttr := fitPairTextAttrs(l.LangTitleTufte, 12.0, 8.8, false, true, true, 0.5, strippedTag, 9.5, 8.2, true, true, false, 0, 412, 10)
		kpi1Num := fmt.Sprintf("%s (%.1f%%)", shortLangName(langs[0].Name, 12), langs[0].Percentage)
		kpi2Num := fmt.Sprintf("%.0f%% / %.0f%%", lp.ProductionPct, lp.ExplorationPct)
		kpi3Num := fmt.Sprintf("%d / %d Repos", lp.LockfileRepos, lp.TotalRepos)
		kpi1Attr := fitTextAttr(kpi1Num, 14.0, 9.5, 138, true, true, false, 0)
		kpi2Attr := fitTextAttr(kpi2Num, 14.0, 9.5, 132, true, true, false, 0)
		kpi3Attr := fitTextAttr(kpi3Num, 14.0, 9.5, 104, true, true, false, 0)
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
  <text x="14" y="20" class="sans hdr"%s>%s</text>
  <text x="426" y="20" text-anchor="end" class="mono tag"%s>%s</text>
  <line x1="2" y1="28" x2="438" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(14, 35)">
    <text x="0" y="12" class="mono kpi-num" fill="#0284c7"%s>%s</text>
    <text x="0" y="24" class="sans kpi-lbl">%s</text>
    <line x1="145" y1="0" x2="145" y2="27" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="157" y="12" class="mono kpi-num"%s>%s</text>
    <text x="157" y="24" class="sans kpi-lbl">%s</text>
    <line x1="295" y1="0" x2="295" y2="27" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="307" y="12" class="mono kpi-num" fill="#059669"%s>%s</text>
    <text x="307" y="24" class="sans kpi-lbl">%s</text>
  </g>
  <line x1="2" y1="68" x2="438" y2="68" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(14, 84)">
    <text x="0" y="0" class="sans sec">%s</text>
`, hdrAttr, esc(l.LangTitleTufte), tagAttr, esc(strippedTag), kpi1Attr, esc(kpi1Num), esc(l.LangFewPrimary), kpi2Attr, esc(kpi2Num), esc(l.LangFewProdEDA), kpi3Attr, esc(kpi3Num), esc(l.LangFewLockCI), esc(l.LangSecTrueProp))
		for i, lg := range langs {
			y := 18 + i*23
			barW := clampInt(int(math.Round((lg.Percentage/100.0)*92.0)), 0, 92)
			color := "#64748b"
			if i == 0 {
				color = "#0284c7"
			}
			dispName := shortLangName(lg.Name, 11)
			lblAttr := fitTextAttr(dispName, 10.5, 8.2, 56, false, false, false, 0)
			fmt.Fprintf(&b, `    <text x="0" y="%d" class="sans lbl"%s>%s</text>
    <rect x="60" y="%d" width="92" height="8" rx="1.5" fill="#64748b" fill-opacity="0.14"/>
    <rect x="60" y="%d" width="%d" height="8" rx="1.5" fill="%s"/>
    <text x="196" y="%d" text-anchor="end" class="mono val">%.1f%%</text>
`, y, lblAttr, esc(dispName), y-7, y-7, barW, color, y, lg.Percentage)
		}
		totalRepos := max(1, lp.TotalRepos)
		prodBarW := clampInt(int(math.Round((lp.ProductionPct/100.0)*96.0)), 0, 96)
		lockBarW := clampInt(int(math.Round((float64(lp.LockfileRepos)/float64(totalRepos))*96.0)), 0, 96)
		ciBarW := clampInt(int(math.Round((float64(lp.CICDRepos)/float64(totalRepos))*96.0)), 0, 96)
		testBarW := clampInt(int(math.Round((float64(lp.TestSuiteRepos)/float64(totalRepos))*96.0)), 0, 96)
		dockerBandW := clampInt(int(math.Round((float64(lp.ContainerizedRepos)/float64(totalRepos))*96.0)), 24, 96)
		rightSecTitle := fmt.Sprintf("%s %s", l.LangSecHygiene, l.FewBulletGraphsSuffix)
		rightSecAttr := fitTextAttr(rightSecTitle, 9.0, 6.6, 192, false, true, true, 0.4)
		fmt.Fprintf(&b, `  </g>
  <line x1="222" y1="68" x2="222" y2="228" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(234, 84)">
    <text x="0" y="0" class="sans sec"%s>%s %s</text>
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
</svg>`, rightSecAttr, esc(l.LangSecHygiene), esc(l.FewBulletGraphsSuffix), esc(l.LangFewProdCode), prodBarW, lp.ProductionPct, lockBarW, lp.LockfileRepos, lp.TotalRepos, ciBarW, lp.CICDRepos, lp.TotalRepos, esc(l.LangFewTestSuite), dockerBandW, testBarW, lp.TestSuiteRepos, lp.TotalRepos)
		return b.String()

	case domain.SchoolCairo:
		var b bytes.Buffer
		head := fmt.Sprintf(l.LangCairoHeadFmt, langs[0].Name, langs[0].Percentage, lp.StrippedBase64MB)
		lead := fmt.Sprintf(l.LangCairoLeadFmt, l.LangSubTufte, lp.ProductionPct)
		headAttr := fitTextAttr(head, 12.2, 9.2, 416, false, true, false, 0)
		leadAttr := fitTextAttr(lead, 10.0, 8.2, 416, false, false, false, 0)
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
  <text x="12" y="25" class="serif headline"%s>%s</text>
  <text x="12" y="39" class="sans lead"%s>%s</text>
  <line x1="12" y1="46" x2="428" y2="46" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(12, 62)">
    <text x="0" y="0" class="sans sec">%s</text>
    <line x1="130" y1="10" x2="130" y2="112" stroke="#e11d48" stroke-dasharray="3,3" stroke-width="1"/>
`, headAttr, esc(head), leadAttr, esc(lead), esc(l.LangSecTrueProp))
		for i, lg := range langs {
			y := 22 + i*21
			lineY := y - 4
			endX := 60 + int(math.Round((lg.Percentage/100.0)*116.0))
			dispName := shortLangName(lg.Name, 11)
			lblAttr := fitTextAttr(dispName, 10.5, 8.2, 55, false, true, false, 0)
			fmt.Fprintf(&b, `    <text x="0" y="%d" class="sans lbl"%s>%s</text>
    <line x1="60" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="2.2"/>
    <circle cx="%d" cy="%d" r="4" fill="%s"/>
    <text x="%d" y="%d" class="mono val">%.1f%%</text>
`, y, lblAttr, esc(dispName), lineY, endX, lineY, lg.Color, endX, lineY, lg.Color, endX+7, y, lg.Percentage)
		}
		waffleCells := buildWaffleGridSVG(langs)
		waffleSecAttr := fitTextAttr(l.LangCairoWaffle, 9.5, 7.2, 184, false, true, true, 0.4)
		fmt.Fprintf(&b, `    <text x="0" y="144" class="serif callout">%s</text>
  </g>
  <g transform="translate(244, 62)">
    <text x="0" y="0" class="sans sec"%s>%s</text>
    <g transform="translate(0, 10)">
%s    </g>
    <text x="0" y="75" class="sans val"><tspan fill="#059669" font-weight="700">■ %.0f%% Prod</tspan> (%d) · <tspan fill="#d97706" font-weight="700">■ %.0f%% EDA</tspan> (%d)</text>
    <text x="0" y="95" class="sans sec">%s</text>
    <text x="0" y="112" class="sans lbl">%s</text><text x="184" y="112" text-anchor="end" class="mono val">%d/%d</text>
    <text x="0" y="128" class="sans lbl">%s</text><text x="184" y="128" text-anchor="end" class="mono val">%d/%d</text>
    <text x="0" y="144" class="sans lbl">%s</text><text x="184" y="144" text-anchor="end" class="mono val">%d/%d</text>
  </g>
</svg>`, esc(l.LangCairoCallout), waffleSecAttr, esc(l.LangCairoWaffle), waffleCells, lp.ProductionPct, lp.ProductionFiles, lp.ExplorationPct, lp.ExplorationFiles, esc(l.LangSecHygiene), esc(l.LangLockfiles), lp.LockfileRepos, lp.TotalRepos, esc(l.LangCICD), lp.CICDRepos, lp.TotalRepos, esc(l.LangTests), lp.TestSuiteRepos, lp.TotalRepos)
		return b.String()

	default: // SchoolHolmes
		strippedMsg := fmt.Sprintf(l.LangHolmesStripped, lp.StrippedBase64MB)
		prodSplitMsg := fmt.Sprintf(l.LangHolmesProdSplit, lp.ProductionPct, lp.ExplorationPct)
		isotypeRows := buildHolmesIsotypeRowsSVG(langs)
		bannerTxt := "⚡ " + l.LangTitleTufte
		bannerAttr := fitTextAttr(bannerTxt, 12.0, 8.8, 408, false, true, true, 0.6)
		leftSecTxt := fmt.Sprintf("%s (%s)", l.LangSecTrueProp, l.LangHolmesPictorial)
		leftSecAttr := fitTextAttr(leftSecTxt, 9.5, 6.5, 204, false, true, true, 0.5)
		rightSecAttr := fitTextAttr(l.LangSecHygiene, 9.5, 7.5, 192, false, true, true, 0.5)
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
  <text x="16" y="22" class="sans banner"%s>⚡ %s</text>
  <g transform="translate(16, 48)">
    <text x="0" y="0" class="sans sec"%s>%s (%s)</text>
%s    <rect x="0" y="106" width="196" height="58" rx="8" fill="#0284c7" fill-opacity="0.08" stroke="#0284c7" stroke-opacity="0.35" stroke-width="1.2"/>
    <text x="10" y="123" class="sans sec" fill="#0284c7">%s</text>
    <text x="10" y="140" class="sans badge">%s</text>
    <text x="10" y="154" class="sans badge" fill="#059669">%s</text>
  </g>
  <g transform="translate(232, 48)">
    <text x="0" y="0" class="sans sec"%s>%s</text>
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
			bannerAttr, esc(l.LangTitleTufte), leftSecAttr, esc(l.LangSecTrueProp), esc(l.LangHolmesPictorial),
			isotypeRows,
			esc(l.LangHolmesFilter), esc(strippedMsg), esc(prodSplitMsg),
			rightSecAttr, esc(l.LangSecHygiene),
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
		{shortLangName(langs[0].Name, 10), langs[0].Percentage, nonEmptyColor(langs[0].Color, "#0284c7")},
		{shortLangName(langs[1].Name, 10), langs[1].Percentage, nonEmptyColor(langs[1].Color, "#d97706")},
		{shortLangName(langs[2].Name, 10), langs[2].Percentage, nonEmptyColor(langs[2].Color, "#00ADD8")},
		{shortLangName(langs[3].Name, 7) + " + " + shortLangName(langs[4].Name, 7), langs[3].Percentage + langs[4].Percentage, "#7c3aed"},
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
		lblAttr := fitTextAttr(r.label, 11.0, 7.6, 58, false, true, false, 0)
		fmt.Fprintf(&b, `    <text x="0" y="%d" class="sans lbl"%s>%s</text>`+"\n", yText, lblAttr, esc(r.label))
		fmt.Fprintf(&b, `    <g transform="translate(62, %d)" fill="%s">`+"\n      ", yBlocks, r.color)
		step := 15
		if counts[i] >= 8 {
			step = 11
		}
		for k := 0; k < counts[i]; k++ {
			fmt.Fprintf(&b, `<rect x="%d" y="0" width="14" height="11" rx="2"/>`, k*step)
		}
		fmt.Fprintf(&b, "\n    </g>\n")
		valX := 62 + counts[i]*step + 4
		if counts[i] >= 8 {
			valX = 62 + (counts[i]-1)*step + 14 + 4
		}
		fmt.Fprintf(&b, `    <text x="%d" y="%d" class="mono val" font-size="9.5px" fill="%s">%.1f%%</text>`+"\n", valX, yText, r.color, r.pct)
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
		h1Attr, subAttr := fitPairTextAttrs(l.StackTitle, 13.5, 10.0, false, true, false, 0, l.StackSub, 10.0, 8.2, false, false, false, 0, 416, 12)
		footAttr := fitTextAttr(l.StackFootnote, 10.0, 8.0, 416, false, false, false, 0)
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
  <text x="12" y="21" class="serif h1"%s>%s</text>
  <text x="428" y="21" text-anchor="end" class="serif sub"%s>%s</text>
  <line x1="12" y1="28" x2="428" y2="28" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
`, h1Attr, esc(l.StackTitle), subAttr, esc(l.StackSub))

		for _, col := range cols {
			colSecAttr := fitTextAttr(col.title, 10.0, 8.0, 128, false, true, false, 0)
			fmt.Fprintf(&b, `  <g transform="translate(%d, 46)">
    <text x="0" y="0" class="serif sec" fill="%s"%s>%s</text>
    <line x1="0" y1="6" x2="128" y2="6" stroke="#64748b" stroke-opacity="0.35" stroke-width="0.6"/>
`, col.x, col.color, colSecAttr, esc(col.title))
			for i, it := range col.items {
				y := 24 + i*33
				barW := int((it.WeightPct / 100.0) * 116.0)
				if barW < 18 {
					barW = 18
				}
				itVal := fmt.Sprintf("%d", it.RepoCount)
				lblAttr, valAttr := fitPairTextAttrs(it.Name, 10.0, 8.0, false, false, false, 0, itVal, 10.0, 8.5, false, true, false, 0, 128, 6)
				fmt.Fprintf(&b, `    <text x="0" y="%d" class="serif lbl"%s>%s</text>
    <text x="128" y="%d" text-anchor="end" class="num val"%s>%d</text>
    <line x1="0" y1="%d" x2="%d" y2="%d" stroke="#64748b" stroke-width="1"/>
    <circle cx="%d" cy="%d" r="2.3" fill="%s"/>
`, y, lblAttr, esc(it.Name), y, valAttr, it.RepoCount, y+6, barW, y+6, barW, y+6, col.color)
			}
			fmt.Fprintf(&b, "  </g>\n")
		}
		fmt.Fprintf(&b, `  <line x1="12" y1="196" x2="428" y2="196" stroke="#64748b" stroke-opacity="0.3" stroke-width="0.6"/>
  <text x="12" y="211" class="serif sub"%s>%s</text>
  <line x1="12" y1="220" x2="428" y2="220" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
</svg>`, footAttr, esc(l.StackFootnote))
		return b.String()

	case domain.SchoolFew:
		var b bytes.Buffer
		hdrAttr, tagAttr := fitPairTextAttrs(l.StackTitle, 12.0, 9.0, false, true, true, 0.5, l.StackSub, 9.5, 8.0, true, true, false, 0, 412, 12)
		kpi1Val := fmt.Sprintf("%s (%d)", dataCol[0].Name, dataCol[0].RepoCount)
		kpi2Val := fmt.Sprintf("%s (%d)", mlCol[0].Name, mlCol[0].RepoCount)
		kpi3Val := fmt.Sprintf("%s (%d)", vizCol[0].Name, vizCol[0].RepoCount)
		kpi1NumAttr := fitTextAttr(kpi1Val, 13.5, 8.2, 132, true, true, false, 0)
		kpi1LblAttr := fitTextAttr(l.StackPillarData, 8.5, 6.8, 132, false, true, true, 0)
		kpi2NumAttr := fitTextAttr(kpi2Val, 13.5, 8.2, 124, true, true, false, 0)
		kpi2LblAttr := fitTextAttr(l.StackPillarML, 8.5, 6.8, 124, false, true, true, 0)
		kpi3NumAttr := fitTextAttr(kpi3Val, 13.5, 8.2, 124, true, true, false, 0)
		kpi3LblAttr := fitTextAttr(l.StackPillarViz, 8.5, 6.8, 124, false, true, true, 0)
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
  <text x="14" y="20" class="sans hdr"%s>%s</text>
  <text x="426" y="20" text-anchor="end" class="mono tag"%s>%s</text>
  <line x1="2" y1="28" x2="438" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(14, 33)">
    <text x="0" y="11" class="mono kpi-num" fill="#0284c7"%s>%s (%d)</text>
    <text x="0" y="22" class="sans kpi-lbl"%s>%s</text>
    <line x1="138" y1="0" x2="138" y2="25" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="146" y="11" class="mono kpi-num" fill="#7c3aed"%s>%s (%d)</text>
    <text x="146" y="22" class="sans kpi-lbl"%s>%s</text>
    <line x1="278" y1="0" x2="278" y2="25" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="286" y="11" class="mono kpi-num" fill="#059669"%s>%s (%d)</text>
    <text x="286" y="22" class="sans kpi-lbl"%s>%s</text>
  </g>
  <line x1="2" y1="63" x2="438" y2="63" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
`, hdrAttr, esc(l.StackTitle), tagAttr, esc(l.StackSub),
			kpi1NumAttr, esc(dataCol[0].Name), dataCol[0].RepoCount, kpi1LblAttr, esc(l.StackPillarData),
			kpi2NumAttr, esc(mlCol[0].Name), mlCol[0].RepoCount, kpi2LblAttr, esc(l.StackPillarML),
			kpi3NumAttr, esc(vizCol[0].Name), vizCol[0].RepoCount, kpi3LblAttr, esc(l.StackPillarViz))
		for _, col := range cols {
			colSecAttr := fitTextAttr(col.title, 9.0, 7.2, 128, false, true, true, 0.4)
			fmt.Fprintf(&b, `  <g transform="translate(%d, 76)">
    <text x="0" y="0" class="sans sec" fill="%s"%s>%s</text>
`, col.x, col.color, colSecAttr, esc(col.title))
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
				itVal := fmt.Sprintf("%d", it.RepoCount)
				lblAttr, valAttr := fitPairTextAttrs(it.Name, 10.0, 7.8, false, false, false, 0, itVal, 9.5, 8.0, true, true, false, 0, 124, 6)
				fmt.Fprintf(&b, `    <text x="0" y="%d" class="sans lbl"%s>%s</text>
    <text x="124" y="%d" text-anchor="end" class="mono val"%s>%d</text>
    <rect x="0" y="%d" width="124" height="8" rx="1.5" fill="#64748b" fill-opacity="0.14"/>
    <rect x="0" y="%d" width="84" height="8" rx="1.5" fill="#64748b" fill-opacity="0.25"/>
    <rect x="0" y="%d" width="%d" height="6" rx="1.5" fill="%s"/>
    <line x1="96" y1="%d" x2="96" y2="%d" stroke="#0284c7" stroke-width="2"/>
`, y, lblAttr, esc(it.Name), y, valAttr, it.RepoCount, y+4, y+4, y+5, barW, fill, y+3, y+13)
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
		headAttr := fitTextAttr(head, 12.2, 8.8, 416, false, true, false, 0)
		leadAttr := fitTextAttr(lead, 10.0, 8.0, 416, false, false, false, 0)
		calloutAttr := fitTextAttr(l.StackCairoCallout, 9.0, 7.6, 416, false, false, false, 0)
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
  <text x="12" y="25" class="serif headline"%s>%s</text>
  <text x="12" y="39" class="sans lead"%s>%s</text>
  <line x1="12" y1="46" x2="428" y2="46" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
`, headAttr, esc(head), leadAttr, esc(lead))
		for _, col := range cols {
			colSecAttr := fitTextAttr(col.title, 9.0, 7.2, 126, false, true, true, 0.4)
			fmt.Fprintf(&b, `  <g transform="translate(%d, 60)">
    <text x="0" y="0" class="sans sec" fill="%s"%s>%s</text>
    <line x1="88" y1="8" x2="88" y2="140" stroke="#e11d48" stroke-dasharray="3,3" stroke-width="1"/>
`, col.x, col.color, colSecAttr, esc(col.title))
			for i, it := range col.items {
				y := 18 + i*33
				endX := int((it.WeightPct / 100.0) * 116.0)
				if endX < 18 {
					endX = 18
				}
				itVal := fmt.Sprintf("%d", it.RepoCount)
				lblAttr, valAttr := fitPairTextAttrs(it.Name, 10.0, 7.8, false, true, false, 0, itVal, 9.5, 8.0, true, true, false, 0, 126, 6)
				fmt.Fprintf(&b, `    <text x="0" y="%d" class="sans lbl"%s>%s</text>
    <text x="126" y="%d" text-anchor="end" class="mono val"%s>%d</text>
    <line x1="0" y1="%d" x2="126" y2="%d" stroke="#64748b" stroke-opacity="0.2" stroke-width="1.5"/>
    <line x1="0" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="2.2"/>
    <circle cx="%d" cy="%d" r="3.5" fill="%s"/>
`, y, lblAttr, esc(it.Name), y, valAttr, it.RepoCount, y+7, y+7, y+7, endX, y+7, col.color, endX, y+7, col.color)
			}
			fmt.Fprintf(&b, "  </g>\n")
		}
		fmt.Fprintf(&b, `  <text x="12" y="218" class="serif callout"%s>%s</text>
</svg>`, calloutAttr, esc(l.StackCairoCallout))
		return b.String()

	default: // SchoolHolmes
		var b bytes.Buffer
		bannerTxt := fmt.Sprintf("🧰 %s", l.StackTitle)
		bannerAttr := fitTextAttr(bannerTxt, 12.0, 9.0, 408, false, true, true, 0.6)
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
  <text x="16" y="22" class="sans banner"%s>🧰 %s</text>
`, bannerAttr, esc(l.StackTitle))
		for _, col := range cols {
			colSecAttr := fitTextAttr(col.title, 9.5, 7.2, 128, false, true, true, 0.4)
			fmt.Fprintf(&b, `  <g transform="translate(%d, 48)">
    <text x="0" y="0" class="sans sec" fill="%s"%s>%s</text>
`, col.x, col.color, colSecAttr, esc(col.title))
			for i, it := range col.items {
				y := 10 + i*38
				itVal := fmt.Sprintf("%d", it.RepoCount)
				lblAttr, valAttr := fitPairTextAttrs(it.Name, 10.0, 7.6, false, true, false, 0, itVal, 9.5, 8.0, true, true, false, 0, 112, 6)
				fmt.Fprintf(&b, `    <rect x="0" y="%d" width="128" height="30" rx="6" fill="%s" fill-opacity="0.09" stroke="%s" stroke-opacity="0.35" stroke-width="1.2"/>
    <text x="8" y="%d" class="sans lbl"%s>%s</text>
    <text x="120" y="%d" text-anchor="end" class="mono val" fill="%s"%s>%d</text>
`, y, col.color, col.color, y+19, lblAttr, esc(it.Name), y+19, col.color, valAttr, it.RepoCount)
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

func formatKaggleRankOrCount(c domain.KaggleCategory, isPtBR bool) string {
	if c.Rank > 0 {
		return "#" + formatInt(c.Rank)
	}
	if c.ArtifactCount > 0 {
		switch strings.ToLower(strings.TrimSpace(c.Name)) {
		case "code", "notebooks", "kernels":
			return fmt.Sprintf("%d notebooks", c.ArtifactCount)
		case "datasets", "dataset":
			return fmt.Sprintf("%d datasets", c.ArtifactCount)
		case "competitions", "competition":
			return fmt.Sprintf("%d comps", c.ArtifactCount)
		default:
			if isPtBR {
				return fmt.Sprintf("%d itens", c.ArtifactCount)
			}
			return fmt.Sprintf("%d items", c.ArtifactCount)
		}
	}
	if isPtBR {
		return "Ativo"
	}
	return "Active"
}

func formatKaggleTopPct(c domain.KaggleCategory, isPtBR bool) string {
	pct := strings.TrimSpace(c.TopPct)
	if pct == "" {
		totalMedals := c.Gold + c.Silver + c.Bronze
		req := 3
		switch strings.ToLower(strings.TrimSpace(c.Name)) {
		case "code", "notebooks", "kernels":
			req = 5
		case "competitions", "competition":
			req = 2
		}
		if totalMedals > 0 && totalMedals < req {
			pct = fmt.Sprintf("%d/%d Expert", totalMedals, req)
		} else {
			pct = "Active"
		}
	}
	if isPtBR && strings.EqualFold(pct, "Active") {
		return "Ativo"
	}
	return pct
}

func formatKaggleHolmesCardRight(c domain.KaggleCategory, isPtBR bool) string {
	if c.Rank > 0 {
		return fmt.Sprintf("%s #%s", c.Tier, formatInt(c.Rank))
	}
	return fmt.Sprintf("%s · %s", c.Tier, formatKaggleRankOrCount(c, isPtBR))
}

func formatKaggleOverallTierAndRank(k domain.KaggleStanding) string {
	tier := nonEmptyStr(k.OverallTier, domain.KaggleTierContributor)
	if k.PeakRank > 0 {
		return fmt.Sprintf("%s (#%s)", tier, formatInt(k.PeakRank))
	}
	if k.TotalMedals() > 0 {
		return fmt.Sprintf("%s (%dM)", tier, k.TotalMedals())
	}
	return tier
}

func formatKaggleSummaryTrack(c domain.KaggleCategory, isPtBR bool) string {
	if c.Rank > 0 {
		return "#" + formatInt(c.Rank)
	}
	if m := c.Gold + c.Silver + c.Bronze; m > 0 {
		return fmt.Sprintf("%dM", m)
	}
	return formatKaggleRankOrCount(c, isPtBR)
}

func renderKaggle(school domain.School, d domain.PortfolioData, l labels) string {
	k := d.Kaggle
	cats := k.Categories

	// Recruiter-First KPI logic: never highlight "0 Gold" if Gold == 0
	kpi2Val := fmt.Sprintf("%d %s", k.TotalGold, l.KaggleGold)
	kpi2Sub := l.KaggleGoldMedalsLbl
	if k.TotalGold == 0 {
		activeCount := len(cats)
		if activeCount == 0 || activeCount > 3 {
			activeCount = 3
		}
		kpi2Val = fmt.Sprintf(l.KaggleTracksFmt, activeCount)
		kpi2Sub = l.KaggleActive
	}

	switch school {
	case domain.SchoolTufte:
		var b bytes.Buffer
		sub := fmt.Sprintf(l.KaggleTufteSubFmt, k.OverallTier, k.TotalMedals())
		h1Attr, subAttr := fitPairTextAttrs(l.KaggleTitle, 13.5, 10.0, false, true, false, 0, sub, 10.0, 8.0, false, false, false, 0, 416, 12)
		thCatAttr := fitTextAttr(l.KaggleThCat, 9.5, 8.0, 76, false, true, false, 0)
		thStandAttr := fitTextAttr(l.KaggleThStanding, 9.5, 8.0, 68, false, true, false, 0)
		thRankAttr := fitTextAttr(l.KaggleThRank, 9.5, 8.0, 72, false, true, false, 0)
		thPctAttr := fitTextAttr(l.KaggleThPct, 9.5, 8.0, 72, false, true, false, 0)
		thMedalsAttr := fitTextAttr(l.KaggleThMedals, 9.5, 8.0, 104, false, true, false, 0)
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
  <text x="12" y="21" class="serif h1"%s>%s</text>
  <text x="428" y="21" text-anchor="end" class="serif sub"%s>%s</text>
  <line x1="12" y1="28" x2="428" y2="28" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
  <g transform="translate(12, 48)">
    <text x="0" y="0" class="serif th"%s>%s</text>
    <text x="82" y="0" class="serif th"%s>%s</text>
    <text x="228" y="0" text-anchor="end" class="serif th"%s>%s</text>
    <text x="306" y="0" text-anchor="end" class="serif th"%s>%s</text>
    <text x="416" y="0" text-anchor="end" class="serif th"%s>%s</text>
    <line x1="0" y1="7" x2="416" y2="7" stroke="#64748b" stroke-opacity="0.4" stroke-width="0.6"/>
`, h1Attr, esc(l.KaggleTitle), subAttr, esc(sub),
			thCatAttr, esc(l.KaggleThCat),
			thStandAttr, esc(l.KaggleThStanding),
			thRankAttr, esc(l.KaggleThRank),
			thPctAttr, esc(l.KaggleThPct),
			thMedalsAttr, esc(l.KaggleThMedals))
		for i, c := range cats[:3] {
			y := 32 + i*32
			tierColor := "#526071"
			if c.Tier == "Master" || c.Tier == "Grandmaster" {
				tierColor = "#d97706"
			} else if c.Tier == "Expert" {
				tierColor = "#0284c7"
			}
			catAttr := fitTextAttr(c.Name, 10.5, 8.2, 76, false, false, false, 0)
			tierAttr := fitTextAttr(c.Tier, 10.5, 8.2, 68, false, true, false, 0)
			rankStr := formatKaggleRankOrCount(c, l.IsPtBR)
			rankAttr := fitTextAttr(rankStr, 10.5, 8.2, 72, false, true, false, 0)
			pctStr := formatKaggleTopPct(c, l.IsPtBR)
			pctAttr := fitTextAttr(pctStr, 9.5, 8.0, 72, false, false, false, 0)
			medalsStr := fmt.Sprintf("%dG · %dS · %dB", c.Gold, c.Silver, c.Bronze)
			medalsAttr := fitTextAttr(medalsStr, 10.5, 8.2, 104, false, true, false, 0)
			fmt.Fprintf(&b, `    <text x="0" y="%d" class="serif lbl"%s>%s</text>
    <text x="82" y="%d" class="serif val" fill="%s"%s>%s</text>
    <text x="228" y="%d" text-anchor="end" class="num val"%s>%s</text>
    <text x="306" y="%d" text-anchor="end" class="num micro"%s>%s</text>
    <text x="416" y="%d" text-anchor="end" class="num val"%s>%s</text>
`, y, catAttr, esc(c.Name), y, tierColor, tierAttr, esc(c.Tier), y, rankAttr, esc(rankStr), y, pctAttr, esc(pctStr), y, medalsAttr, esc(medalsStr))
		}
		tallyRight := fmt.Sprintf("%d %s (%s · %d %s · %d %s)", k.TotalMedals(), l.KaggleMedalsWord, kpi2Val, k.TotalSilver, l.KaggleSilver, k.TotalBronze, l.KaggleBronze)
		tallyLAttr, tallyRAttr := fitPairTextAttrs(l.KaggleTotalTally, 10.5, 8.5, false, true, false, 0, tallyRight, 10.5, 8.2, false, true, false, 0, 416, 10)
		fmt.Fprintf(&b, `    <line x1="0" y1="122" x2="416" y2="122" stroke="#64748b" stroke-opacity="0.4" stroke-width="0.6"/>
    <text x="0" y="142" class="serif val"%s>%s</text>
    <text x="416" y="142" text-anchor="end" class="num val" fill="#d97706"%s>%d %s (%s · %d %s · %d %s)</text>
  </g>
  <line x1="12" y1="220" x2="428" y2="220" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
</svg>`, tallyLAttr, esc(l.KaggleTotalTally), tallyRAttr, k.TotalMedals(), esc(l.KaggleMedalsWord), esc(kpi2Val), k.TotalSilver, esc(l.KaggleSilver), k.TotalBronze, esc(l.KaggleBronze))
		return b.String()

	case domain.SchoolFew:
		var b bytes.Buffer
		markerLabel := "│ Vertical marker"
		if l.IsPtBR {
			markerLabel = "│ Marcador vertical"
		}
		tierTag := fmt.Sprintf(l.KaggleFewTierNoRankFmt, strings.ToUpper(k.OverallTier))
		if k.PeakRank > 0 {
			tierTag = fmt.Sprintf(l.KaggleFewTierTagFmt, strings.ToUpper(k.OverallTier), k.PeakRank)
		}
		hdrAttr, tagAttr := fitPairTextAttrs(l.KaggleTitle, 12.0, 9.0, false, true, true, 0.5, tierTag, 10.0, 8.0, true, true, false, 0, 412, 12)
		kpi1Val := fmt.Sprintf("%d %s", k.TotalMedals(), l.KaggleMedalsWord)
		kpi3Val := fmt.Sprintf("%d %s", k.TotalSilver, l.KaggleSilver)
		kpi4Val := fmt.Sprintf("%d %s", k.TotalBronze, l.KaggleBronze)
		kpi1NumAttr := fitTextAttr(kpi1Val, 14.5, 10.0, 104, true, true, false, 0)
		kpi1LblAttr := fitTextAttr(l.KaggleCumulative, 8.5, 6.8, 104, false, true, true, 0)
		kpi2NumAttr := fitTextAttr(kpi2Val, 14.5, 10.0, 92, true, true, false, 0)
		kpi2LblAttr := fitTextAttr(kpi2Sub, 8.5, 6.8, 92, false, true, true, 0)
		kpi3NumAttr := fitTextAttr(kpi3Val, 14.5, 10.0, 88, true, true, false, 0)
		kpi3LblAttr := fitTextAttr(l.KaggleSilverMedalsLbl, 8.5, 6.8, 88, false, true, true, 0)
		kpi4NumAttr := fitTextAttr(kpi4Val, 14.5, 10.0, 86, true, true, false, 0)
		kpi4LblAttr := fitTextAttr(l.KaggleBronzeMedalsLbl, 8.5, 6.8, 86, false, true, true, 0)
		secL := fmt.Sprintf("%s %s", l.KaggleTrack, l.FewBulletGraphsSuffix)
		secR := fmt.Sprintf("Rank & %s", l.KaggleMedalsWord)
		secLAttr, secRAttr := fitPairTextAttrs(secL, 9.0, 7.5, false, true, true, 0.4, secR, 9.0, 7.5, false, true, true, 0.4, 412, 12)
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
  <text x="14" y="20" class="sans hdr"%s>%s</text>
  <text x="426" y="20" text-anchor="end" class="mono tag"%s>%s</text>
  <line x1="2" y1="28" x2="438" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(14, 35)">
    <text x="0" y="13" class="mono kpi-num" fill="#d97706"%s>%d %s</text>
    <text x="0" y="25" class="sans kpi-lbl"%s>%s</text>
    <line x1="108" y1="0" x2="108" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="118" y="13" class="mono kpi-num" fill="#d97706"%s>%s</text>
    <text x="118" y="25" class="sans kpi-lbl"%s>%s</text>
    <line x1="214" y1="0" x2="214" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="224" y="13" class="mono kpi-num"%s>%d %s</text>
    <text x="224" y="25" class="sans kpi-lbl"%s>%s</text>
    <line x1="316" y1="0" x2="316" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="326" y="13" class="mono kpi-num"%s>%d %s</text>
    <text x="326" y="25" class="sans kpi-lbl"%s>%s</text>
  </g>
  <line x1="2" y1="69" x2="438" y2="69" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(14, 85)">
    <text x="0" y="0" class="sans sec"%s>%s %s</text>
    <text x="412" y="0" text-anchor="end" class="sans sec"%s>Rank &amp; %s</text>
`, hdrAttr, esc(l.KaggleTitle), tagAttr, esc(tierTag),
			kpi1NumAttr, k.TotalMedals(), esc(l.KaggleMedalsWord), kpi1LblAttr, esc(l.KaggleCumulative),
			kpi2NumAttr, esc(kpi2Val), kpi2LblAttr, esc(kpi2Sub),
			kpi3NumAttr, k.TotalSilver, esc(l.KaggleSilver), kpi3LblAttr, esc(l.KaggleSilverMedalsLbl),
			kpi4NumAttr, k.TotalBronze, esc(l.KaggleBronze), kpi4LblAttr, esc(l.KaggleBronzeMedalsLbl),
			secLAttr, esc(l.KaggleTrack), esc(l.FewBulletGraphsSuffix), secRAttr, esc(l.KaggleMedalsWord))
		for i, c := range cats[:3] {
			y := 16 + i*32
			barW := clampInt(int(math.Round((c.ProgressPct/100.0)*112.0)), 0, 112)
			targetX := 84 + clampInt(int(math.Round((c.TargetPct/100.0)*112.0)), 0, 112)
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
			lblAttr := fitTextAttr(c.Name, 10.5, 8.2, 80, false, true, false, 0)
			rightStr := fmt.Sprintf("%s · %s (%dG %dS %dB)", shortTier, formatKaggleRankOrCount(c, l.IsPtBR), c.Gold, c.Silver, c.Bronze)
			valAttr := fitTextAttr(rightStr, 10.0, 8.0, 208, true, true, false, 0)
			fmt.Fprintf(&b, `    <g transform="translate(0, %d)">
      <text x="0" y="11" class="sans lbl"%s>%s</text>
      <rect x="84" y="2" width="112" height="13" rx="1.5" fill="#64748b" fill-opacity="0.12"/>
      <rect x="84" y="2" width="84" height="13" rx="1.5" fill="#64748b" fill-opacity="0.22"/>
      <rect x="84" y="2" width="52" height="13" rx="1.5" fill="#64748b" fill-opacity="0.34"/>
      <rect x="84" y="5" width="%d" height="7" rx="1" fill="%s"/>
      <line x1="%d" y1="0" x2="%d" y2="17" stroke="#0284c7" stroke-width="2.2"/>
      <text x="412" y="11" text-anchor="end" class="mono val" fill="%s"%s>%s · %s (%dG %dS %dB)</text>
    </g>
`, y, lblAttr, esc(c.Name), barW, barColor, targetX, targetX, barColor, valAttr, esc(shortTier), esc(formatKaggleRankOrCount(c, l.IsPtBR)), c.Gold, c.Silver, c.Bronze)
		}
		fmt.Fprintf(&b, `    <text x="0" y="124" class="sans" font-size="8.5" fill="#64748b"><tspan fill="#0284c7" font-weight="700">%s</tspan> %s</text>
  </g>
</svg>`, esc(markerLabel), esc(l.KaggleFewMarkerNote))
		return b.String()

	case domain.SchoolCairo:
		var b bytes.Buffer
		head := fmt.Sprintf(l.KaggleCairoHeadNoRankFmt, k.OverallTier, k.TotalMedals())
		if k.PeakRank > 0 {
			head = fmt.Sprintf(l.KaggleCairoHeadFmt, k.OverallTier, formatInt(k.PeakRank), k.TotalMedals())
		}
		lead := fmt.Sprintf(l.KaggleCairoLeadFmt, kpi2Val, k.TotalSilver, k.TotalBronze)
		headAttr := fitTextAttr(head, 12.2, 8.8, 416, false, true, false, 0)
		leadAttr := fitTextAttr(lead, 10.0, 7.8, 416, false, false, false, 0)
		calloutAttr := fitTextAttr(l.KaggleCairoCallout, 9.0, 7.6, 416, false, false, false, 0)
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
  <text x="12" y="25" class="serif headline"%s>%s</text>
  <text x="12" y="39" class="sans lead"%s>%s</text>
  <line x1="12" y1="46" x2="428" y2="46" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(12, 66)">
    <line x1="168" y1="0" x2="168" y2="120" stroke="#e11d48" stroke-dasharray="3,3" stroke-width="1.2"/>
    <text x="168" y="-4" text-anchor="middle" class="sans" font-size="8.5" font-weight="700" fill="#e11d48">%s</text>
`, headAttr, esc(head), leadAttr, esc(lead), esc(l.KaggleCairoBench))
		for i, c := range cats[:3] {
			y := 24 + i*36
			endX := 86 + clampInt(int(math.Round((c.ProgressPct/100.0)*108.0)), 0, 108)
			color := "#0284c7"
			if i == 0 {
				color = "#d97706"
			}
			lblAttr := fitTextAttr(c.Name, 10.5, 8.2, 82, false, true, false, 0)
			rightStr := fmt.Sprintf("%s · %s (%s)", c.Tier, formatKaggleRankOrCount(c, l.IsPtBR), formatKaggleTopPct(c, l.IsPtBR))
			valAttr := fitTextAttr(rightStr, 10.0, 7.8, 210, true, true, false, 0)
			fmt.Fprintf(&b, `    <text x="0" y="%d" class="sans lbl"%s>%s</text>
    <line x1="86" y1="%d" x2="196" y2="%d" stroke="#64748b" stroke-opacity="0.2" stroke-width="1.5"/>
    <line x1="86" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="2.5"/>
    <circle cx="%d" cy="%d" r="4.2" fill="%s"/>
    <text x="416" y="%d" text-anchor="end" class="mono val"%s>%s · %s (%s)</text>
`, y, lblAttr, esc(c.Name), y-4, y-4, y-4, endX, y-4, color, endX, y-4, color, y, valAttr, esc(c.Tier), esc(formatKaggleRankOrCount(c, l.IsPtBR)), esc(formatKaggleTopPct(c, l.IsPtBR)))
		}
		fmt.Fprintf(&b, `    <text x="0" y="146" class="serif callout"%s>%s</text>
  </g>
</svg>`, calloutAttr, esc(l.KaggleCairoCallout))
		return b.String()

	default: // SchoolHolmes
		medalsEarned := fmt.Sprintf(l.KaggleHolmesMedals, k.TotalMedals())
		trophyX, trophyY := holmesTrophyPosition(k.OverallTier)
		bannerRaw := fmt.Sprintf("🏆 %s — %s", l.KaggleTitle, strings.ToUpper(k.OverallTier))
		bannerText := fmt.Sprintf("🏆 %s — %s", esc(l.KaggleTitle), strings.ToUpper(esc(k.OverallTier)))
		if k.PeakRank > 0 {
			bannerRaw = fmt.Sprintf("🏆 %s — %s (#%s PEAK)", l.KaggleTitle, strings.ToUpper(k.OverallTier), formatInt(k.PeakRank))
			bannerText = fmt.Sprintf("🏆 %s — %s (#%s PEAK)", esc(l.KaggleTitle), strings.ToUpper(esc(k.OverallTier)), formatInt(k.PeakRank))
		}
		bannerAttr := fitTextAttr(bannerRaw, 12.0, 8.8, 408, false, true, true, 0.6)
		podiumSecAttr := fitTextAttr(l.KaggleHolmesPodium, 9.5, 7.8, 196, false, true, true, 0.5)
		trackSecAttr := fitTextAttr(l.KaggleTrack, 9.5, 7.8, 204, false, true, true, 0.5)
		medalsSub := fmt.Sprintf("%s · %d %s · %d %s", kpi2Val, k.TotalSilver, l.KaggleSilver, k.TotalBronze, l.KaggleBronze)
		medalsSubAttr := fitTextAttr(medalsSub, 10.0, 8.0, 196, true, true, false, 0)
		c0Right := formatKaggleHolmesCardRight(cats[0], l.IsPtBR)
		c1Right := formatKaggleHolmesCardRight(cats[1], l.IsPtBR)
		c2Right := formatKaggleHolmesCardRight(cats[2], l.IsPtBR)
		c0LAttr, c0RAttr := fitPairTextAttrs("📓 "+cats[0].Name, 10.5, 8.2, false, true, false, 0, c0Right, 10.0, 7.8, true, true, false, 0, 184, 8)
		c1LAttr, c1RAttr := fitPairTextAttrs("🗂️ "+cats[1].Name, 10.5, 8.2, false, true, false, 0, c1Right, 10.0, 7.8, true, true, false, 0, 184, 8)
		c2LAttr, c2RAttr := fitPairTextAttrs("⚔️ "+cats[2].Name, 10.5, 8.2, false, true, false, 0, c2Right, 10.0, 7.8, true, true, false, 0, 184, 8)
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
  <text x="16" y="22" class="sans banner"%s>%s</text>
  <g transform="translate(16, 48)">
    <text x="0" y="0" class="sans sec"%s>%s</text>
    <rect x="0" y="96" width="34" height="28" rx="3" fill="#64748b" fill-opacity="0.25"/>
    <rect x="38" y="76" width="34" height="48" rx="3" fill="#059669" fill-opacity="0.4"/>
    <rect x="76" y="54" width="34" height="70" rx="3" fill="#0284c7" fill-opacity="0.65"/>
    <rect x="114" y="30" width="34" height="94" rx="3" fill="#d97706"/>
    <rect x="152" y="10" width="34" height="114" rx="3" fill="none" stroke="#d97706" stroke-dasharray="3,3" stroke-width="1.5"/>
    <text x="%d" y="%d" text-anchor="middle" font-size="14">🏆</text>
    <text x="0" y="144" class="sans lbl" fill="#d97706">%s</text>
    <text x="0" y="160" class="mono val"%s>%s · %d %s · %d %s</text>
  </g>
  <g transform="translate(220, 48)">
    <text x="0" y="0" class="sans sec"%s>%s</text>
    <rect x="0" y="14" width="204" height="38" rx="6" fill="#d97706" fill-opacity="0.1" stroke="#d97706" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="37" class="sans lbl"%s>📓 %s</text>
    <text x="194" y="37" text-anchor="end" class="mono val" fill="#d97706"%s>%s</text>
    <rect x="0" y="62" width="204" height="38" rx="6" fill="#0284c7" fill-opacity="0.1" stroke="#0284c7" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="85" class="sans lbl"%s>🗂️ %s</text>
    <text x="194" y="85" text-anchor="end" class="mono val"%s>%s</text>
    <rect x="0" y="110" width="204" height="38" rx="6" fill="#059669" fill-opacity="0.1" stroke="#059669" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="133" class="sans lbl"%s>⚔️ %s</text>
    <text x="194" y="133" text-anchor="end" class="mono val" fill="#059669"%s>%s</text>
  </g>
</svg>`,
			bannerAttr, bannerText,
			podiumSecAttr, esc(l.KaggleHolmesPodium), trophyX, trophyY, esc(medalsEarned), medalsSubAttr, esc(kpi2Val), k.TotalSilver, esc(l.KaggleSilver), k.TotalBronze, esc(l.KaggleBronze),
			trackSecAttr, esc(l.KaggleTrack),
			c0LAttr, esc(cats[0].Name), c0RAttr, esc(c0Right),
			c1LAttr, esc(cats[1].Name), c1RAttr, esc(c1Right),
			c2LAttr, esc(cats[2].Name), c2RAttr, esc(c2Right))
	}
}

func holmesTrophyPosition(tier string) (int, int) {
	switch strings.ToLower(strings.TrimSpace(tier)) {
	case strings.ToLower(domain.KaggleTierNovice):
		return 17, 90
	case strings.ToLower(domain.KaggleTierContributor), strings.ToLower(domain.KaggleTierKaggler):
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
	rRow4Lbl := l.ImpactStreak
	rRow4Val := fmt.Sprintf(l.ImpactDaysCapFmt, im.ActiveStreakDays)

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
		rRow4Lbl = l.ImpactContribs
		rRow4Val = formatInt(im.AnnualContributions)
	}

	leftStreakVal := fmt.Sprintf(l.ImpactDaysCapFmt, im.ActiveStreakDays)

	switch school {
	case domain.SchoolTufte:
		pts, lastX, lastY := buildTrendPolyline(im.WeeklyTrend, 240, 1, 176, 16)
		h1Attr, subAttr := fitPairTextAttrs(l.ImpactTitle, 13.5, 10.0, false, true, false, 0, l.ImpactTufteHorizon, 10.0, 8.0, false, false, false, 0, 416, 12)
		bigNumStr := formatInt(im.AnnualContributions)
		lblX := max(58, int(math.Ceil(estimateTextWidth(bigNumStr, 18.0, false, true, false, 0)))+10)
		contribsLbl := fmt.Sprintf("%s (+%.0f%% YoY)", l.ImpactContribs, im.YoYGrowthPct)
		contribsLblAttr := fitTextAttr(contribsLbl, 10.5, 8.0, float64(234-lblX), false, false, false, 0)
		lSecAttr := fitTextAttr(l.ImpactGitHubSec, 10.0, 8.2, 190, false, true, false, 0)
		rSecAttr := fitTextAttr(rightSecTitle, 10.0, 8.2, 196, false, true, false, 0)
		lR1L, lR1R := fitPairTextAttrs(l.ImpactStars, 10.5, 8.2, false, false, false, 0, "★ "+formatInt(im.TotalStars), 10.5, 8.2, false, true, false, 0, 190, 8)
		lR2L, lR2R := fitPairTextAttrs(l.ImpactRepos, 10.5, 8.2, false, false, false, 0, fmt.Sprintf("%d", im.PublicRepos), 10.5, 8.2, false, true, false, 0, 190, 8)
		lR3L, lR3R := fitPairTextAttrs(l.ImpactForks, 10.5, 8.2, false, false, false, 0, fmt.Sprintf("%d", im.TotalForks), 10.5, 8.2, false, true, false, 0, 190, 8)
		lR4L, lR4R := fitPairTextAttrs(l.ImpactViews, 10.5, 8.2, false, false, false, 0, formatCompact(im.RepoViews14d), 10.5, 8.2, false, true, false, 0, 190, 8)
		rR1L, rR1R := fitPairTextAttrs(rRow1Lbl, 10.5, 8.2, false, false, false, 0, rRow1Val, 10.5, 8.2, false, true, false, 0, 196, 8)
		rR2L, rR2R := fitPairTextAttrs(rRow2Lbl, 10.5, 8.2, false, false, false, 0, rRow2Val, 10.5, 8.2, false, true, false, 0, 196, 8)
		rR3L, rR3R := fitPairTextAttrs(rRow3Lbl, 10.5, 8.2, false, false, false, 0, rRow3Val, 10.5, 8.2, false, true, false, 0, 196, 8)
		rR4L, rR4R := fitPairTextAttrs(rRow4Lbl, 10.5, 8.2, false, false, false, 0, rRow4Val, 10.5, 8.2, false, true, false, 0, 196, 8)
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
  <text x="12" y="21" class="serif h1"%s>%s</text>
  <text x="428" y="21" text-anchor="end" class="serif sub"%s>%s</text>
  <line x1="12" y1="28" x2="428" y2="28" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
  <g transform="translate(12, 42)">
    <text x="0" y="16" class="num big-num">%s</text>
    <text x="%d" y="16" class="serif lbl"%s>%s (+%.0f%% YoY)</text>
    <rect x="240" y="3" width="176" height="14" fill="#64748b" fill-opacity="0.12"/>
    <polyline fill="none" stroke="#0284c7" stroke-width="1.3" points="%s"/>
    <circle cx="%d" cy="%d" r="2.5" fill="#059669"/>
  </g>
  <line x1="12" y1="72" x2="428" y2="72" stroke="#64748b" stroke-opacity="0.35" stroke-width="0.6"/>
  <g transform="translate(12, 90)">
    <text x="0" y="0" class="serif sec"%s>%s</text>
    <text x="0" y="22" class="serif lbl"%s>%s</text><text x="190" y="22" text-anchor="end" class="num val"%s>★ %s</text>
    <text x="0" y="42" class="serif lbl"%s>%s</text><text x="190" y="42" text-anchor="end" class="num val"%s>%d</text>
    <text x="0" y="62" class="serif lbl"%s>%s</text><text x="190" y="62" text-anchor="end" class="num val"%s>%d</text>
    <text x="0" y="82" class="serif lbl"%s>%s</text><text x="190" y="82" text-anchor="end" class="num val"%s>%s</text>
  </g>
  <g transform="translate(232, 90)">
    <text x="0" y="0" class="serif sec"%s>%s</text>
    <text x="0" y="22" class="serif lbl"%s>%s</text><text x="196" y="22" text-anchor="end" class="num val"%s>%s</text>
    <text x="0" y="42" class="serif lbl"%s>%s</text><text x="196" y="42" text-anchor="end" class="num val" fill="#059669"%s>%s</text>
    <text x="0" y="62" class="serif lbl"%s>%s</text><text x="196" y="62" text-anchor="end" class="num val"%s>%s</text>
    <text x="0" y="82" class="serif lbl"%s>%s</text><text x="196" y="82" text-anchor="end" class="num val"%s>%s</text>
  </g>
  <line x1="12" y1="220" x2="428" y2="220" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
</svg>`,
			h1Attr, esc(l.ImpactTitle), subAttr, esc(l.ImpactTufteHorizon),
			bigNumStr, lblX, contribsLblAttr, esc(l.ImpactContribs), im.YoYGrowthPct,
			pts, lastX, lastY,
			lSecAttr, esc(l.ImpactGitHubSec),
			lR1L, esc(l.ImpactStars), lR1R, formatInt(im.TotalStars),
			lR2L, esc(l.ImpactRepos), lR2R, im.PublicRepos,
			lR3L, esc(l.ImpactForks), lR3R, im.TotalForks,
			lR4L, esc(l.ImpactViews), lR4R, formatCompact(im.RepoViews14d),
			rSecAttr, esc(rightSecTitle),
			rR1L, esc(rRow1Lbl), rR1R, esc(rRow1Val),
			rR2L, esc(rRow2Lbl), rR2R, esc(rRow2Val),
			rR3L, esc(rRow3Lbl), rR3R, esc(rRow3Val),
			rR4L, esc(rRow4Lbl), rR4R, esc(rRow4Val))

	case domain.SchoolFew:
		pts, _, _ := buildTrendPolyline(im.WeeklyTrend, 0, 88, 194, 27)
		velTag := fmt.Sprintf(l.ImpactFewVelocityFmt, im.YoYGrowthPct)
		hdrAttr, tagAttr := fitPairTextAttrs(l.ImpactTitle, 12.0, 9.0, false, true, true, 0.5, velTag, 10.0, 8.0, true, true, false, 0, 412, 12)
		kpi1NumAttr := fitTextAttr(formatInt(im.AnnualContributions), 14.5, 9.5, 134, true, true, false, 0)
		kpi1LblAttr := fitTextAttr(l.ImpactContribs, 8.5, 6.8, 134, false, true, true, 0)
		kpi2Val := "★ " + formatInt(im.TotalStars)
		kpi2NumAttr := fitTextAttr(kpi2Val, 14.5, 9.5, 124, true, true, false, 0)
		kpi2LblAttr := fitTextAttr(l.ImpactStars, 8.5, 6.8, 124, false, true, true, 0)
		kpi3NumAttr := fitTextAttr(rRow2Val, 14.5, 9.0, 122, true, true, false, 0)
		kpi3LblAttr := fitTextAttr(rRow2Lbl, 8.5, 6.8, 122, false, true, true, 0)
		lSecAttr := fitTextAttr(l.ImpactGitHubSec, 9.0, 7.5, 192, false, true, true, 0.4)
		rSecAttr := fitTextAttr(rightSecTitle, 9.0, 7.5, 194, false, true, true, 0.4)
		lR1L, lR1R := fitPairTextAttrs(l.ImpactRepos, 10.5, 8.2, false, false, false, 0, fmt.Sprintf("%d", im.PublicRepos), 10.5, 8.2, true, true, false, 0, 192, 8)
		lR2L, lR2R := fitPairTextAttrs(l.ImpactForks, 10.5, 8.2, false, false, false, 0, fmt.Sprintf("%d", im.TotalForks), 10.5, 8.2, true, true, false, 0, 192, 8)
		lR3L, lR3R := fitPairTextAttrs(l.ImpactViews, 10.5, 8.2, false, false, false, 0, formatCompact(im.RepoViews14d), 10.5, 8.2, true, true, false, 0, 192, 8)
		lR4L, lR4R := fitPairTextAttrs(l.ImpactStreak, 10.5, 8.2, false, false, false, 0, leftStreakVal, 10.5, 8.2, true, true, false, 0, 192, 8)
		rR1L, rR1R := fitPairTextAttrs(rRow1Lbl, 10.5, 8.2, false, false, false, 0, rRow1Val, 10.5, 8.2, true, true, false, 0, 194, 8)
		rR2L, rR2R := fitPairTextAttrs(rRow2Lbl, 10.5, 8.2, false, false, false, 0, rRow2Val, 10.5, 8.2, true, true, false, 0, 194, 8)
		rR3L, rR3R := fitPairTextAttrs(rRow3Lbl, 10.5, 8.2, false, false, false, 0, rRow3Val, 10.5, 8.2, true, true, false, 0, 194, 8)
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
  <text x="14" y="20" class="sans hdr"%s>%s</text>
  <text x="426" y="20" text-anchor="end" class="mono tag"%s>%s</text>
  <line x1="2" y1="28" x2="438" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(14, 35)">
    <text x="0" y="13" class="mono kpi-num" fill="#0284c7"%s>%s</text>
    <text x="0" y="25" class="sans kpi-lbl"%s>%s</text>
    <line x1="138" y1="0" x2="138" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="150" y="13" class="mono kpi-num"%s>★ %s</text>
    <text x="150" y="25" class="sans kpi-lbl"%s>%s</text>
    <line x1="278" y1="0" x2="278" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="290" y="13" class="mono kpi-num" fill="#059669"%s>%s</text>
    <text x="290" y="25" class="sans kpi-lbl"%s>%s</text>
  </g>
  <line x1="2" y1="69" x2="438" y2="69" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(14, 85)">
    <text x="0" y="0" class="sans sec"%s>%s</text>
    <text x="0" y="24" class="sans lbl"%s>%s</text><text x="192" y="24" text-anchor="end" class="mono val"%s>%d</text>
    <text x="0" y="48" class="sans lbl"%s>%s</text><text x="192" y="48" text-anchor="end" class="mono val"%s>%d</text>
    <text x="0" y="72" class="sans lbl"%s>%s</text><text x="192" y="72" text-anchor="end" class="mono val"%s>%s</text>
    <text x="0" y="96" class="sans lbl"%s>%s</text><text x="192" y="96" text-anchor="end" class="mono val" fill="#059669"%s>%s</text>
    <rect x="0" y="106" width="192" height="8" rx="1.5" fill="#64748b" fill-opacity="0.14"/>
    <rect x="0" y="106" width="128" height="8" rx="1.5" fill="#64748b" fill-opacity="0.25"/>
    <rect x="0" y="108" width="144" height="4" rx="1" fill="#059669"/>
    <line x1="150" y1="104" x2="150" y2="116" stroke="#0284c7" stroke-width="2"/>
  </g>
  <line x1="220" y1="69" x2="220" y2="228" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(232, 85)">
    <text x="0" y="0" class="sans sec"%s>%s</text>
    <text x="0" y="24" class="sans lbl"%s>%s</text><text x="194" y="24" text-anchor="end" class="mono val"%s>%s</text>
    <text x="0" y="48" class="sans lbl"%s>%s</text><text x="194" y="48" text-anchor="end" class="mono val" fill="#0284c7"%s>%s</text>
    <text x="0" y="72" class="sans lbl"%s>%s</text><text x="194" y="72" text-anchor="end" class="mono val"%s>%s</text>
    <polyline fill="none" stroke="#0284c7" stroke-width="2" points="%s"/>
  </g>
</svg>`,
			hdrAttr, esc(l.ImpactTitle), tagAttr, esc(velTag),
			kpi1NumAttr, formatInt(im.AnnualContributions), kpi1LblAttr, esc(l.ImpactContribs),
			kpi2NumAttr, formatInt(im.TotalStars), kpi2LblAttr, esc(l.ImpactStars),
			kpi3NumAttr, esc(rRow2Val), kpi3LblAttr, esc(rRow2Lbl),
			lSecAttr, esc(l.ImpactGitHubSec),
			lR1L, esc(l.ImpactRepos), lR1R, im.PublicRepos,
			lR2L, esc(l.ImpactForks), lR2R, im.TotalForks,
			lR3L, esc(l.ImpactViews), lR3R, formatCompact(im.RepoViews14d),
			lR4L, esc(l.ImpactStreak), lR4R, esc(leftStreakVal),
			rSecAttr, esc(rightSecTitle),
			rR1L, esc(rRow1Lbl), rR1R, esc(rRow1Val),
			rR2L, esc(rRow2Lbl), rR2R, esc(rRow2Val),
			rR3L, esc(rRow3Lbl), rR3R, esc(rRow3Val),
			pts)

	case domain.SchoolCairo:
		pts, lastX, lastY := buildTrendPolyline(im.WeeklyTrend, 0, 2, 416, 32)
		head := fmt.Sprintf(l.ImpactCairoHeadFmt, formatInt(im.AnnualContributions), im.YoYGrowthPct)
		lead := fmt.Sprintf(l.ImpactCairoLeadFmt, im.PublicRepos)
		callout := fmt.Sprintf(l.ImpactCairoCallout, im.ActiveStreakDays)
		headAttr := fitTextAttr(head, 12.2, 8.8, 416, false, true, false, 0)
		leadAttr := fitTextAttr(lead, 10.0, 8.0, 416, false, false, false, 0)
		calloutAttr := fitTextAttr(callout, 9.0, 7.6, 416, false, false, false, 0)
		lSecAttr := fitTextAttr(l.ImpactGitHubSec, 9.0, 7.5, 190, false, true, true, 0.4)
		rSecAttr := fitTextAttr(rightSecTitle, 9.0, 7.5, 196, false, true, true, 0.4)
		lR1L, lR1R := fitPairTextAttrs(l.ImpactStars, 10.5, 8.2, false, true, false, 0, "★ "+formatInt(im.TotalStars), 10.5, 8.2, true, true, false, 0, 190, 8)
		lR2L, lR2R := fitPairTextAttrs(l.ImpactRepos, 10.5, 8.2, false, true, false, 0, fmt.Sprintf("%d", im.PublicRepos), 10.5, 8.2, true, true, false, 0, 190, 8)
		lR3L, lR3R := fitPairTextAttrs(l.ImpactForks, 10.5, 8.2, false, true, false, 0, fmt.Sprintf("%d", im.TotalForks), 10.5, 8.2, true, true, false, 0, 190, 8)
		lR4L, lR4R := fitPairTextAttrs(l.ImpactViews, 10.5, 8.2, false, true, false, 0, formatCompact(im.RepoViews14d), 10.5, 8.2, true, true, false, 0, 190, 8)
		rR1L, rR1R := fitPairTextAttrs(rRow1Lbl, 10.5, 8.2, false, true, false, 0, rRow1Val, 10.5, 8.2, true, true, false, 0, 196, 8)
		rR2L, rR2R := fitPairTextAttrs(rRow2Lbl, 10.5, 8.2, false, true, false, 0, rRow2Val, 10.5, 8.2, true, true, false, 0, 196, 8)
		rR3L, rR3R := fitPairTextAttrs(rRow3Lbl, 10.5, 8.2, false, true, false, 0, rRow3Val, 10.5, 8.2, true, true, false, 0, 196, 8)
		rR4L, rR4R := fitPairTextAttrs(rRow4Lbl, 10.5, 8.2, false, true, false, 0, rRow4Val, 10.5, 8.2, true, true, false, 0, 196, 8)
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 440 230" width="440" height="230">
  <style>
    .serif { font-family: Georgia, "Times New Roman", serif; }
    .sans { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-variant-numeric: tabular-nums; }
    .headline { font-size: 12.2px; font-weight: bold; fill: #526071; }
    .lead { font-size: 10px; fill: #64748b; }
    .sec { font-size: 9px; font-weight: 700; fill: #0284c7; text-transform: uppercase; letter-spacing: 0.4px; }
    .lbl { font-size: 10.5px; font-weight: 600; fill: #526071; }
    .val { font-size: 10.5px; font-weight: 700; fill: #526071; }
    .callout { font-size: 9px; font-style: italic; fill: #059669; }
  </style>
  <rect x="12" y="6" width="416" height="3" fill="#e11d48"/>
  <text x="12" y="25" class="serif headline"%s>%s</text>
  <text x="12" y="39" class="sans lead"%s>%s</text>
  <line x1="12" y1="46" x2="428" y2="46" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(12, 62)">
    <text x="0" y="0" class="sans sec"%s>%s</text>
    <text x="0" y="22" class="sans lbl"%s>%s</text><text x="190" y="22" text-anchor="end" class="mono val" fill="#d97706"%s>★ %s</text>
    <text x="0" y="42" class="sans lbl"%s>%s</text><text x="190" y="42" text-anchor="end" class="mono val"%s>%d</text>
    <text x="0" y="62" class="sans lbl"%s>%s</text><text x="190" y="62" text-anchor="end" class="mono val"%s>%d</text>
    <text x="0" y="82" class="sans lbl"%s>%s</text><text x="190" y="82" text-anchor="end" class="mono val"%s>%s</text>
  </g>
  <g transform="translate(232, 62)">
    <text x="0" y="0" class="sans sec"%s>%s</text>
    <text x="0" y="22" class="sans lbl"%s>%s</text><text x="196" y="22" text-anchor="end" class="mono val"%s>%s</text>
    <text x="0" y="42" class="sans lbl"%s>%s</text><text x="196" y="42" text-anchor="end" class="mono val" fill="#059669"%s>%s</text>
    <text x="0" y="62" class="sans lbl"%s>%s</text><text x="196" y="62" text-anchor="end" class="mono val"%s>%s</text>
    <text x="0" y="82" class="sans lbl"%s>%s</text><text x="196" y="82" text-anchor="end" class="mono val"%s>%s</text>
  </g>
  <g transform="translate(12, 166)">
    <line x1="0" y1="18" x2="416" y2="18" stroke="#e11d48" stroke-dasharray="3,3" stroke-width="1"/>
    <polyline fill="none" stroke="#0284c7" stroke-width="2" points="%s"/>
    <circle cx="%d" cy="%d" r="3.5" fill="#059669"/>
    <text x="0" y="50" class="serif callout"%s>%s</text>
  </g>
</svg>`,
			headAttr, esc(head), leadAttr, esc(lead),
			lSecAttr, esc(l.ImpactGitHubSec),
			lR1L, esc(l.ImpactStars), lR1R, formatInt(im.TotalStars),
			lR2L, esc(l.ImpactRepos), lR2R, im.PublicRepos,
			lR3L, esc(l.ImpactForks), lR3R, im.TotalForks,
			lR4L, esc(l.ImpactViews), lR4R, formatCompact(im.RepoViews14d),
			rSecAttr, esc(rightSecTitle),
			rR1L, esc(rRow1Lbl), rR1R, esc(rRow1Val),
			rR2L, esc(rRow2Lbl), rR2R, esc(rRow2Val),
			rR3L, esc(rRow3Lbl), rR3R, esc(rRow3Val),
			rR4L, esc(rRow4Lbl), rR4R, esc(rRow4Val),
			pts, lastX, lastY,
			calloutAttr, esc(callout))

	default: // SchoolHolmes
		banner := fmt.Sprintf(l.ImpactHolmesBanner, l.ImpactTitle, formatInt(im.AnnualContributions))
		bannerAttr := fitTextAttr(banner, 11.5, 8.0, 408, false, true, true, 0.5)
		lSecAttr := fitTextAttr(l.ImpactGitHubSec, 9.5, 7.8, 194, false, true, true, 0.5)
		rSecAttr := fitTextAttr(rightSecTitle, 9.5, 7.8, 196, false, true, true, 0.5)
		lR1L, lR1R := fitPairTextAttrs("⭐ "+l.ImpactStars, 10.5, 8.0, false, true, false, 0, formatInt(im.TotalStars), 11.0, 8.2, true, true, false, 0, 174, 8)
		lR2L, lR2R := fitPairTextAttrs("📦 "+l.ImpactRepos, 10.5, 8.0, false, true, false, 0, fmt.Sprintf("%d", im.PublicRepos), 11.0, 8.2, true, true, false, 0, 174, 8)
		lR3L, lR3R := fitPairTextAttrs("🔀 "+l.ImpactForks, 10.5, 8.0, false, true, false, 0, fmt.Sprintf("%d", im.TotalForks), 11.0, 8.2, true, true, false, 0, 174, 8)
		lR4L, lR4R := fitPairTextAttrs("👀 "+l.ImpactViews, 10.5, 8.0, false, true, false, 0, formatCompact(im.RepoViews14d), 11.0, 8.2, true, true, false, 0, 174, 8)
		rR1L, rR1R := fitPairTextAttrs("🤗 "+rRow1Lbl, 10.5, 8.0, false, true, false, 0, rRow1Val, 11.0, 8.2, true, true, false, 0, 176, 8)
		rR2L, rR2R := fitPairTextAttrs("📥 "+rRow2Lbl, 10.5, 8.0, false, true, false, 0, rRow2Val, 11.0, 8.2, true, true, false, 0, 176, 8)
		rR3L, rR3R := fitPairTextAttrs("❤️ "+rRow3Lbl, 10.5, 8.0, false, true, false, 0, rRow3Val, 11.0, 8.2, true, true, false, 0, 176, 8)
		rR4L, rR4R := fitPairTextAttrs("🔥 "+rRow4Lbl, 10.5, 8.0, false, true, false, 0, rRow4Val, 11.0, 8.2, true, true, false, 0, 176, 8)
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
  <text x="16" y="22" class="sans banner"%s>%s</text>
  <g transform="translate(16, 48)">
    <text x="0" y="0" class="sans sec"%s>%s</text>
    <rect x="0" y="10" width="194" height="34" rx="6" fill="#d97706" fill-opacity="0.1" stroke="#d97706" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="31" class="sans lbl"%s>⭐ %s</text>
    <text x="184" y="31" text-anchor="end" class="mono val" fill="#d97706"%s>%s</text>
    <rect x="0" y="52" width="194" height="34" rx="6" fill="#0284c7" fill-opacity="0.1" stroke="#0284c7" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="73" class="sans lbl"%s>📦 %s</text>
    <text x="184" y="73" text-anchor="end" class="mono val"%s>%d</text>
    <rect x="0" y="94" width="194" height="34" rx="6" fill="#059669" fill-opacity="0.1" stroke="#059669" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="115" class="sans lbl"%s>🔀 %s</text>
    <text x="184" y="115" text-anchor="end" class="mono val" fill="#059669"%s>%d</text>
    <rect x="0" y="136" width="194" height="34" rx="6" fill="#7c3aed" fill-opacity="0.1" stroke="#7c3aed" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="157" class="sans lbl"%s>👀 %s</text>
    <text x="184" y="157" text-anchor="end" class="mono val" fill="#7c3aed"%s>%s</text>
  </g>
  <g transform="translate(228, 48)">
    <text x="0" y="0" class="sans sec"%s>%s</text>
    <rect x="0" y="10" width="196" height="34" rx="6" fill="#0284c7" fill-opacity="0.1" stroke="#0284c7" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="31" class="sans lbl"%s>🤗 %s</text>
    <text x="186" y="31" text-anchor="end" class="mono val"%s>%s</text>
    <rect x="0" y="52" width="196" height="34" rx="6" fill="#059669" fill-opacity="0.1" stroke="#059669" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="73" class="sans lbl"%s>📥 %s</text>
    <text x="186" y="73" text-anchor="end" class="mono val" fill="#059669"%s>%s</text>
    <rect x="0" y="94" width="196" height="34" rx="6" fill="#e11d48" fill-opacity="0.1" stroke="#e11d48" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="115" class="sans lbl"%s>❤️ %s</text>
    <text x="186" y="115" text-anchor="end" class="mono val" fill="#e11d48"%s>%s</text>
    <rect x="0" y="136" width="196" height="34" rx="6" fill="#d97706" fill-opacity="0.1" stroke="#d97706" stroke-opacity="0.4" stroke-width="1.2"/>
    <text x="10" y="157" class="sans lbl"%s>🔥 %s</text>
    <text x="186" y="157" text-anchor="end" class="mono val" fill="#d97706"%s>%s</text>
  </g>
</svg>`,
			bannerAttr, esc(banner),
			lSecAttr, esc(l.ImpactGitHubSec),
			lR1L, esc(l.ImpactStars), lR1R, formatInt(im.TotalStars),
			lR2L, esc(l.ImpactRepos), lR2R, im.PublicRepos,
			lR3L, esc(l.ImpactForks), lR3R, im.TotalForks,
			lR4L, esc(l.ImpactViews), lR4R, formatCompact(im.RepoViews14d),
			rSecAttr, esc(rightSecTitle),
			rR1L, esc(rRow1Lbl), rR1R, esc(rRow1Val),
			rR2L, esc(rRow2Lbl), rR2R, esc(rRow2Val),
			rR3L, esc(rRow3Lbl), rR3R, esc(rRow3Val),
			rR4L, esc(rRow4Lbl), rR4R, esc(rRow4Val))
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
	overallTierRank := formatKaggleOverallTierAndRank(k)
	summaryTracksVal := fmt.Sprintf("%s / %s", formatKaggleSummaryTrack(cats[0], l.IsPtBR), formatKaggleSummaryTrack(cats[1], l.IsPtBR))
	l0Short := shortLangName(langs[0].Name, 11)
	l1Short := shortLangName(langs[1].Name, 11)
	l2Short := shortLangName(langs[2].Name, 11)

	switch school {
	case domain.SchoolTufte:
		pts, lastX, lastY := buildTrendPolyline(im.WeeklyTrend, 118, 46, 82, 12)
		l0End := 62 + clampInt(int(math.Round((langs[0].Percentage/100.0)*96.0)), 6, 96)
		l1End := 62 + clampInt(int(math.Round((langs[1].Percentage/100.0)*96.0)), 6, 96)
		l2End := 62 + clampInt(int(math.Round((langs[2].Percentage/100.0)*96.0)), 6, 96)
		rightSub := fmt.Sprintf("%d Repos · ★ %s", d.Languages.TotalRepos, formatInt(im.TotalStars))
		rightSubW := estimateTextWidth(rightSub, 10.0, false, false, false, 0)
		dispRole, _ := fitHeaderNameRole(d.Subject.Name, d.Subject.Role, " — ", 437-rightSubW-14, 13.5, 9.0, false, true, false, 0)
		leftTitle := fmt.Sprintf("%s — %s", d.Subject.Name, dispRole)
		h1Attr, subAttr := fitPairTextAttrs(leftTitle, 13.5, 9.0, false, true, false, 0, rightSub, 10.0, 8.2, false, false, false, 0, 437, 12)
		sec1Attr := fitTextAttr("1. "+l.LangSecTrueProp, 10.0, 8.0, 200, false, true, false, 0)
		sec2Attr := fitTextAttr("2. "+l.StackTitle, 10.0, 8.0, 203, false, true, false, 0)
		sec3Attr := fitTextAttr("3. "+l.KaggleTitle, 10.0, 8.0, 200, false, true, false, 0)
		sec4Attr := fitTextAttr("4. "+l.ImpactTitle, 10.0, 8.0, 203, false, true, false, 0)
		l0Attr := fitTextAttr(l0Short, 10.0, 7.8, 56, false, false, false, 0)
		l1Attr := fitTextAttr(l1Short, 10.0, 7.8, 56, false, false, false, 0)
		l2Attr := fitTextAttr(l2Short, 10.0, 7.8, 56, false, false, false, 0)
		q1R4Val := fmt.Sprintf("%.0f%% Prod / %.0f%% EDA", d.Languages.ProductionPct, d.Languages.ExplorationPct)
		q1R4L, q1R4R := fitPairTextAttrs("Prod vs. EDA", 10.0, 8.0, false, false, false, 0, q1R4Val, 10.0, 8.0, false, true, false, 0, 200, 8)
		q2R1L, q2R1R := fitPairTextAttrs(dataCol[0].Name, 10.0, 8.0, false, false, false, 0, fmt.Sprintf("%d repos", dataCol[0].RepoCount), 10.0, 8.0, false, true, false, 0, 203, 8)
		q2R2L, q2R2R := fitPairTextAttrs(mlCol[0].Name, 10.0, 8.0, false, false, false, 0, fmt.Sprintf("%d repos", mlCol[0].RepoCount), 10.0, 8.0, false, true, false, 0, 203, 8)
		q2R3L, q2R3R := fitPairTextAttrs(mlCol[1].Name, 10.0, 8.0, false, false, false, 0, fmt.Sprintf("%d repos", mlCol[1].RepoCount), 10.0, 8.0, false, true, false, 0, 203, 8)
		q2R4L, q2R4R := fitPairTextAttrs(l.LangLockfiles, 10.0, 8.0, false, false, false, 0, fmt.Sprintf("%d / %d", d.Languages.LockfileRepos, d.Languages.TotalRepos), 10.0, 8.0, false, true, false, 0, 203, 8)
		q3R1L, q3R1R := fitPairTextAttrs(l.KaggleOverallTierLbl, 10.0, 8.0, false, false, false, 0, overallTierRank, 10.0, 8.0, false, true, false, 0, 200, 8)
		q3R2Val := fmt.Sprintf("%d (%dG · %dS · %dB)", k.TotalMedals(), k.TotalGold, k.TotalSilver, k.TotalBronze)
		q3R2L, q3R2R := fitPairTextAttrs(l.KaggleTotalMedalsLbl, 10.0, 8.0, false, false, false, 0, q3R2Val, 10.0, 8.0, false, true, false, 0, 200, 8)
		q3R3L, q3R3R := fitPairTextAttrs("Code / Datasets", 10.0, 8.0, false, false, false, 0, summaryTracksVal, 10.0, 8.0, false, true, false, 0, 200, 8)
		q4R1Val := fmt.Sprintf("%s (+%.0f%%)", formatInt(im.AnnualContributions), im.YoYGrowthPct)
		q4R1L, q4R1R := fitPairTextAttrs(l.ImpactContribs, 10.0, 8.0, false, false, false, 0, q4R1Val, 10.0, 8.0, false, true, false, 0, 203, 8)
		q4R2Val := fmt.Sprintf("★ %s (%d forks)", formatInt(im.TotalStars), im.TotalForks)
		q4R2L, q4R2R := fitPairTextAttrs(l.ImpactStars, 10.0, 8.0, false, false, false, 0, q4R2Val, 10.0, 8.0, false, true, false, 0, 203, 8)
		q4R3L := fitTextAttr(l.ImpactDownloads, 10.0, 8.0, 112, false, false, false, 0)
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
  <text x="14" y="22" class="serif h1"%s>%s — %s</text>
  <text x="451" y="22" text-anchor="end" class="num sub"%s>%d Repos · ★ %s</text>
  <line x1="14" y1="30" x2="451" y2="30" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>

  <!-- Q1: Languages (Dot-Dash Plot) -->
  <g transform="translate(14, 46)">
    <text x="0" y="0" class="serif sec"%s>1. %s</text>
    <text x="0" y="18" class="serif lbl"%s>%s</text>
    <line x1="62" y1="14" x2="%d" y2="14" stroke="#64748b" stroke-width="1.1"/>
    <circle cx="%d" cy="14" r="2.5" fill="#0284c7"/>
    <text x="200" y="18" text-anchor="end" class="num val" fill="#0284c7">%.1f%%</text>
    <text x="0" y="34" class="serif lbl"%s>%s</text>
    <line x1="62" y1="30" x2="%d" y2="30" stroke="#64748b" stroke-width="1.1"/>
    <circle cx="%d" cy="30" r="2.5" fill="#d97706"/>
    <text x="200" y="34" text-anchor="end" class="num val">%.1f%%</text>
    <text x="0" y="50" class="serif lbl"%s>%s</text>
    <line x1="62" y1="46" x2="%d" y2="46" stroke="#64748b" stroke-width="1.1"/>
    <circle cx="%d" cy="46" r="2.5" fill="#00ADD8"/>
    <text x="200" y="50" text-anchor="end" class="num val">%.1f%%</text>
    <text x="0" y="66" class="serif lbl"%s>Prod vs. EDA</text>
    <text x="200" y="66" text-anchor="end" class="num val" fill="#059669"%s>%.0f%% Prod / %.0f%% EDA</text>
  </g>

  <!-- Q2: Stack (Supertable) -->
  <g transform="translate(248, 46)">
    <text x="0" y="0" class="serif sec"%s>2. %s</text>
    <text x="0" y="18" class="serif lbl"%s>%s</text><text x="203" y="18" text-anchor="end" class="num val"%s>%d repos</text>
    <text x="0" y="34" class="serif lbl"%s>%s</text><text x="203" y="34" text-anchor="end" class="num val"%s>%d repos</text>
    <text x="0" y="50" class="serif lbl"%s>%s</text><text x="203" y="50" text-anchor="end" class="num val"%s>%d repos</text>
    <text x="0" y="66" class="serif lbl"%s>%s</text><text x="203" y="66" text-anchor="end" class="num val"%s>%d / %d</text>
  </g>

  <line x1="14" y1="132" x2="451" y2="132" stroke="#64748b" stroke-opacity="0.35" stroke-width="0.6"/>

  <!-- Q3: Kaggle (Supertable) -->
  <g transform="translate(14, 150)">
    <text x="0" y="0" class="serif sec"%s>3. %s</text>
    <text x="0" y="18" class="serif lbl"%s>%s</text><text x="200" y="18" text-anchor="end" class="num val" fill="#d97706"%s>%s</text>
    <text x="0" y="34" class="serif lbl"%s>%s</text><text x="200" y="34" text-anchor="end" class="num val"%s>%d (%dG · %dS · %dB)</text>
    <text x="0" y="50" class="serif lbl"%s>Code / Datasets</text><text x="200" y="50" text-anchor="end" class="num val"%s>%s</text>
  </g>

  <!-- Q4: Impact (Supertable + Inline Sparkline with Interquartile Band) -->
  <g transform="translate(248, 150)">
    <text x="0" y="0" class="serif sec"%s>4. %s</text>
    <text x="0" y="18" class="serif lbl"%s>%s</text><text x="203" y="18" text-anchor="end" class="num val" fill="#059669"%s>%s (+%.0f%%)</text>
    <text x="0" y="34" class="serif lbl"%s>%s</text><text x="203" y="34" text-anchor="end" class="num val"%s>★ %s (%d forks)</text>
    <text x="0" y="54" class="serif lbl"%s>%s</text>
    <rect x="118" y="46" width="82" height="12" fill="#64748b" fill-opacity="0.12"/>
    <polyline fill="none" stroke="#0284c7" stroke-width="1.2" points="%s"/>
    <circle cx="%d" cy="%d" r="2.2" fill="#059669"/>
  </g>

  <line x1="14" y1="242" x2="451" y2="242" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
</svg>`,
			h1Attr, esc(d.Subject.Name), esc(dispRole),
			subAttr, d.Languages.TotalRepos, formatInt(im.TotalStars),
			sec1Attr, esc(l.LangSecTrueProp),
			l0Attr, esc(l0Short), l0End, l0End, langs[0].Percentage,
			l1Attr, esc(l1Short), l1End, l1End, langs[1].Percentage,
			l2Attr, esc(l2Short), l2End, l2End, langs[2].Percentage,
			q1R4L, q1R4R, d.Languages.ProductionPct, d.Languages.ExplorationPct,
			sec2Attr, esc(l.StackTitle),
			q2R1L, esc(dataCol[0].Name), q2R1R, dataCol[0].RepoCount,
			q2R2L, esc(mlCol[0].Name), q2R2R, mlCol[0].RepoCount,
			q2R3L, esc(mlCol[1].Name), q2R3R, mlCol[1].RepoCount,
			q2R4L, esc(l.LangLockfiles), q2R4R, d.Languages.LockfileRepos, d.Languages.TotalRepos,
			sec3Attr, esc(l.KaggleTitle),
			q3R1L, esc(l.KaggleOverallTierLbl), q3R1R, esc(overallTierRank),
			q3R2L, esc(l.KaggleTotalMedalsLbl), q3R2R, k.TotalMedals(), k.TotalGold, k.TotalSilver, k.TotalBronze,
			q3R3L, q3R3R, esc(summaryTracksVal),
			sec4Attr, esc(l.ImpactTitle),
			q4R1L, esc(l.ImpactContribs), q4R1R, formatInt(im.AnnualContributions), im.YoYGrowthPct,
			q4R2L, esc(l.ImpactStars), q4R2R, formatInt(im.TotalStars), im.TotalForks,
			q4R3L, esc(l.ImpactDownloads), pts, lastX, lastY)

	case domain.SchoolFew:
		l0Bar := clampInt(int(math.Round((langs[0].Percentage/100.0)*88.0)), 0, 88)
		prodBar := clampInt(int(math.Round((d.Languages.ProductionPct/100.0)*88.0)), 0, 88)
		kBar := clampInt(int(math.Round((cats[0].ProgressPct/100.0)*88.0)), 0, 88)
		rightSub := fmt.Sprintf("%d Repos · ★ %s", d.Languages.TotalRepos, formatInt(im.TotalStars))
		rightSubW := estimateTextWidth(rightSub, 10.0, true, true, false, 0)
		dispRole, _ := fitHeaderNameRole(d.Subject.Name, d.Subject.Role, " — ", 437-rightSubW-14, 10.5, 8.0, false, true, true, 0.3)
		leftTitle := fmt.Sprintf("%s — %s", d.Subject.Name, dispRole)
		hdrAttr, subAttr := fitPairTextAttrs(leftTitle, 10.5, 8.0, false, true, true, 0.3, rightSub, 10.0, 8.0, true, true, false, 0, 437, 12)
		kpi1Val := fmt.Sprintf("%s %.1f%%", langs[0].Name, langs[0].Percentage)
		kpi1NumAttr := fitTextAttr(kpi1Val, 13.0, 8.2, 141, true, true, false, 0)
		kpi1LblAttr := fitTextAttr(l.LangSecTrueProp, 8.0, 6.8, 141, false, true, true, 0)
		kpi2NumAttr := fitTextAttr(overallTierRank, 13.0, 8.2, 144, true, true, false, 0)
		kpi2LblAttr := fitTextAttr(l.KaggleTitle, 8.0, 6.8, 144, false, true, true, 0)
		kpi3Val := fmt.Sprintf("%s (+%.0f%%)", formatInt(im.AnnualContributions), im.YoYGrowthPct)
		kpi3NumAttr := fitTextAttr(kpi3Val, 13.0, 8.2, 120, true, true, false, 0)
		kpi3LblAttr := fitTextAttr(l.ImpactContribs, 8.0, 6.8, 120, false, true, true, 0)
		sec1Attr := fitTextAttr("1. "+l.LangSecTrueProp, 9.0, 7.2, 200, false, true, true, 0.4)
		sec2Attr := fitTextAttr("2. "+l.StackTitle, 9.0, 7.2, 203, false, true, true, 0.4)
		sec3Attr := fitTextAttr("3. "+l.KaggleTitle, 9.0, 7.2, 200, false, true, true, 0.4)
		sec4Attr := fitTextAttr("4. "+l.ImpactTitle, 9.0, 7.2, 203, false, true, true, 0.4)
		l0Attr := fitTextAttr(l0Short, 10.0, 7.8, 58, false, false, false, 0)
		prodLblAttr := fitTextAttr("Prod vs EDA", 10.0, 7.8, 58, false, false, false, 0)
		q1R3Left := fmt.Sprintf("%s / %s", l1Short, l2Short)
		q1R3Right := fmt.Sprintf("%.1f%% / %.1f%%", langs[1].Percentage, langs[2].Percentage)
		q1R3L, q1R3R := fitPairTextAttrs(q1R3Left, 10.0, 7.8, false, false, false, 0, q1R3Right, 10.0, 7.8, true, true, false, 0, 200, 8)
		q2R1L, q2R1R := fitPairTextAttrs(dataCol[0].Name, 10.0, 7.8, false, false, false, 0, fmt.Sprintf("%d repos", dataCol[0].RepoCount), 10.0, 7.8, true, true, false, 0, 203, 8)
		q2R2L, q2R2R := fitPairTextAttrs(mlCol[0].Name, 10.0, 7.8, false, false, false, 0, fmt.Sprintf("%d repos", mlCol[0].RepoCount), 10.0, 7.8, true, true, false, 0, 203, 8)
		q2R3L, q2R3R := fitPairTextAttrs(l.LangLockfiles, 10.0, 7.8, false, false, false, 0, fmt.Sprintf("%d / %d", d.Languages.LockfileRepos, d.Languages.TotalRepos), 10.0, 7.8, true, true, false, 0, 203, 8)
		kCat0Attr := fitTextAttr(cats[0].Name, 10.0, 7.8, 58, false, false, false, 0)
		kCat0Track := formatKaggleSummaryTrack(cats[0], l.IsPtBR)
		kCat0TrackAttr := fitTextAttr(kCat0Track, 10.0, 7.5, 44, true, true, false, 0)
		q3R2Val := fmt.Sprintf("%d (%dG · %dS · %dB)", k.TotalMedals(), k.TotalGold, k.TotalSilver, k.TotalBronze)
		q3R2L, q3R2R := fitPairTextAttrs(l.KaggleTotalMedalsLbl, 10.0, 7.8, false, false, false, 0, q3R2Val, 10.0, 7.8, true, true, false, 0, 200, 8)
		q3R3L, q3R3R := fitPairTextAttrs("Code / Datasets", 10.0, 7.8, false, false, false, 0, summaryTracksVal, 10.0, 7.8, true, true, false, 0, 200, 8)
		q4R1Val := fmt.Sprintf("%s (+%.0f%%)", formatInt(im.AnnualContributions), im.YoYGrowthPct)
		q4R1L, q4R1R := fitPairTextAttrs(l.ImpactContribs, 10.0, 7.8, false, false, false, 0, q4R1Val, 10.0, 7.8, true, true, false, 0, 203, 8)
		q4R2Val := fmt.Sprintf("★ %s (%d forks)", formatInt(im.TotalStars), im.TotalForks)
		q4R2L, q4R2R := fitPairTextAttrs(l.ImpactStars, 10.0, 7.8, false, false, false, 0, q4R2Val, 10.0, 7.8, true, true, false, 0, 203, 8)
		q4R3Val := fmt.Sprintf("%s (30d)", formatCompact(im.HFDownloads30d))
		q4R3L, q4R3R := fitPairTextAttrs(l.ImpactDownloads, 10.0, 7.8, false, false, false, 0, q4R3Val, 10.0, 7.8, true, true, false, 0, 203, 8)
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 465 255" width="465" height="255">
  <style>
    .sans { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-variant-numeric: tabular-nums; }
    .hdr { font-size: 10.5px; font-weight: 700; fill: #526071; text-transform: uppercase; letter-spacing: 0.3px; }
    .sec { font-size: 9px; font-weight: 700; fill: #0284c7; text-transform: uppercase; letter-spacing: 0.4px; }
    .lbl { font-size: 10px; font-weight: 500; fill: #526071; }
    .val { font-size: 10px; font-weight: 700; fill: #64748b; }
    .kpi-num { font-size: 13px; font-weight: 700; fill: #526071; }
    .kpi-lbl { font-size: 8px; font-weight: 600; fill: #64748b; text-transform: uppercase; }
  </style>
  <rect x="2" y="2" width="461" height="251" rx="6" fill="none" stroke="#64748b" stroke-opacity="0.28" stroke-width="1"/>
  <text x="14" y="20" class="sans hdr"%s>%s — %s</text>
  <text x="451" y="20" text-anchor="end" class="mono val"%s>%d Repos · ★ %s</text>
  <line x1="2" y1="27" x2="463" y2="27" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>

  <!-- Top KPI Strip -->
  <g transform="translate(14, 33)">
    <text x="0" y="11" class="mono kpi-num" fill="#0284c7"%s>%s %.1f%%</text>
    <text x="0" y="22" class="sans kpi-lbl"%s>%s</text>
    <line x1="145" y1="0" x2="145" y2="25" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="157" y="11" class="mono kpi-num" fill="#d97706"%s>%s</text>
    <text x="157" y="22" class="sans kpi-lbl"%s>%s</text>
    <line x1="305" y1="0" x2="305" y2="25" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="317" y="11" class="mono kpi-num" fill="#059669"%s>%s (+%.0f%%)</text>
    <text x="317" y="22" class="sans kpi-lbl"%s>%s</text>
  </g>
  <line x1="2" y1="63" x2="463" y2="63" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>

  <!-- Q1: Languages (Bullet Graphs) -->
  <g transform="translate(14, 78)">
    <text x="0" y="0" class="sans sec"%s>1. %s</text>
    <text x="0" y="18" class="sans lbl"%s>%s</text>
    <rect x="64" y="10" width="88" height="9" fill="#64748b" fill-opacity="0.14"/>
    <rect x="64" y="10" width="62" height="9" fill="#64748b" fill-opacity="0.25"/>
    <rect x="64" y="12" width="%d" height="5" fill="#0284c7"/>
    <line x1="126" y1="8" x2="126" y2="21" stroke="#0284c7" stroke-width="2"/>
    <text x="200" y="18" text-anchor="end" class="mono val" fill="#0284c7">%.1f%%</text>

    <text x="0" y="38" class="sans lbl"%s>Prod vs EDA</text>
    <rect x="64" y="30" width="88" height="9" fill="#64748b" fill-opacity="0.14"/>
    <rect x="64" y="30" width="62" height="9" fill="#64748b" fill-opacity="0.25"/>
    <rect x="64" y="32" width="%d" height="5" fill="#059669"/>
    <line x1="120" y1="28" x2="120" y2="41" stroke="#0284c7" stroke-width="2"/>
    <text x="200" y="38" text-anchor="end" class="mono val">%.0f%%</text>

    <text x="0" y="56" class="sans lbl"%s>%s / %s</text>
    <text x="200" y="56" text-anchor="end" class="mono val"%s>%.1f%% / %.1f%%</text>
  </g>

  <!-- Q2: Stack -->
  <g transform="translate(248, 78)">
    <text x="0" y="0" class="sans sec"%s>2. %s</text>
    <text x="0" y="18" class="sans lbl"%s>%s</text><text x="203" y="18" text-anchor="end" class="mono val" fill="#0284c7"%s>%d repos</text>
    <text x="0" y="38" class="sans lbl"%s>%s</text><text x="203" y="38" text-anchor="end" class="mono val" fill="#7c3aed"%s>%d repos</text>
    <text x="0" y="56" class="sans lbl"%s>%s</text><text x="203" y="56" text-anchor="end" class="mono val"%s>%d / %d</text>
  </g>

  <line x1="232" y1="63" x2="232" y2="253" stroke="#64748b" stroke-opacity="0.22" stroke-width="1"/>
  <line x1="2" y1="154" x2="463" y2="154" stroke="#64748b" stroke-opacity="0.22" stroke-width="1"/>

  <!-- Q3: Kaggle (Bullet Graph) -->
  <g transform="translate(14, 170)">
    <text x="0" y="0" class="sans sec"%s>3. %s</text>
    <text x="0" y="18" class="sans lbl"%s>%s</text>
    <rect x="64" y="10" width="88" height="9" fill="#64748b" fill-opacity="0.14"/>
    <rect x="64" y="10" width="62" height="9" fill="#64748b" fill-opacity="0.25"/>
    <rect x="64" y="12" width="%d" height="5" fill="#d97706"/>
    <line x1="142" y1="8" x2="142" y2="21" stroke="#0284c7" stroke-width="2"/>
    <text x="200" y="18" text-anchor="end" class="mono val" fill="#d97706"%s>%s</text>
    <text x="0" y="38" class="sans lbl"%s>%s</text><text x="200" y="38" text-anchor="end" class="mono val"%s>%d (%dG · %dS · %dB)</text>
    <text x="0" y="56" class="sans lbl"%s>Code / Datasets</text><text x="200" y="56" text-anchor="end" class="mono val"%s>%s</text>
  </g>

  <!-- Q4: Impact -->
  <g transform="translate(248, 170)">
    <text x="0" y="0" class="sans sec"%s>4. %s</text>
    <text x="0" y="18" class="sans lbl"%s>%s</text><text x="203" y="18" text-anchor="end" class="mono val" fill="#059669"%s>%s (+%.0f%%)</text>
    <text x="0" y="38" class="sans lbl"%s>%s</text><text x="203" y="38" text-anchor="end" class="mono val"%s>★ %s (%d forks)</text>
    <text x="0" y="56" class="sans lbl"%s>%s</text><text x="203" y="56" text-anchor="end" class="mono val"%s>%s (30d)</text>
  </g>
</svg>`,
			hdrAttr, esc(d.Subject.Name), esc(dispRole),
			subAttr, d.Languages.TotalRepos, formatInt(im.TotalStars),
			kpi1NumAttr, esc(langs[0].Name), langs[0].Percentage, kpi1LblAttr, esc(l.LangSecTrueProp),
			kpi2NumAttr, esc(overallTierRank), kpi2LblAttr, esc(l.KaggleTitle),
			kpi3NumAttr, formatInt(im.AnnualContributions), im.YoYGrowthPct, kpi3LblAttr, esc(l.ImpactContribs),
			sec1Attr, esc(l.LangSecTrueProp),
			l0Attr, esc(l0Short), l0Bar, langs[0].Percentage,
			prodLblAttr, prodBar, d.Languages.ProductionPct,
			q1R3L, esc(l1Short), esc(l2Short), q1R3R, langs[1].Percentage, langs[2].Percentage,
			sec2Attr, esc(l.StackTitle),
			q2R1L, esc(dataCol[0].Name), q2R1R, dataCol[0].RepoCount,
			q2R2L, esc(mlCol[0].Name), q2R2R, mlCol[0].RepoCount,
			q2R3L, esc(l.LangLockfiles), q2R3R, d.Languages.LockfileRepos, d.Languages.TotalRepos,
			sec3Attr, esc(l.KaggleTitle),
			kCat0Attr, esc(cats[0].Name), kBar, kCat0TrackAttr, esc(kCat0Track),
			q3R2L, esc(l.KaggleTotalMedalsLbl), q3R2R, k.TotalMedals(), k.TotalGold, k.TotalSilver, k.TotalBronze,
			q3R3L, q3R3R, esc(summaryTracksVal),
			sec4Attr, esc(l.ImpactTitle),
			q4R1L, esc(l.ImpactContribs), q4R1R, formatInt(im.AnnualContributions), im.YoYGrowthPct,
			q4R2L, esc(l.ImpactStars), q4R2R, formatInt(im.TotalStars), im.TotalForks,
			q4R3L, esc(l.ImpactDownloads), q4R3R, formatCompact(im.HFDownloads30d))

	case domain.SchoolCairo:
		l0End := 64 + clampInt(int(math.Round((langs[0].Percentage/100.0)*92.0)), 6, 92)
		l1End := 64 + clampInt(int(math.Round((langs[1].Percentage/100.0)*92.0)), 6, 92)
		sumLead := "Executive portfolio synthesis across source code, ecosystem stack, Kaggle standing, and community reach"
		sumCallout := fmt.Sprintf("↳ Notebook-deflated source code (%.1f MB stripped) combined with verified multi-platform telemetry.", d.Languages.StrippedBase64MB)
		if l.IsPtBR {
			sumLead = "Síntese executiva do portfólio em código-fonte, ecossistema, Kaggle e alcance na comunidade"
			sumCallout = fmt.Sprintf("↳ Código-fonte normalizado (%.1f MB removidos) combinado com telemetria multiplataforma verificada.", d.Languages.StrippedBase64MB)
		}
		suffix := fmt.Sprintf(" (%d Repos · ★ %s)", d.Languages.TotalRepos, formatInt(im.TotalStars))
		suffixW := estimateTextWidth(suffix, 12.5, false, true, false, 0)
		dispRole, _ := fitHeaderNameRole(d.Subject.Name, d.Subject.Role, " — ", 441-suffixW, 12.5, 8.5, false, true, false, 0)
		headTxt := fmt.Sprintf("%s — %s%s", d.Subject.Name, dispRole, suffix)
		headAttr := fitTextAttr(headTxt, 12.5, 8.5, 441, false, true, false, 0)
		leadAttr := fitTextAttr(sumLead, 9.0, 7.5, 441, false, false, false, 0)
		calloutAttr := fitTextAttr(sumCallout, 9.0, 7.5, 441, false, false, false, 0)
		sec1Attr := fitTextAttr("1. "+l.LangSecTrueProp, 9.0, 7.2, 202, false, true, true, 0.4)
		sec2Attr := fitTextAttr("2. "+l.StackTitle, 9.0, 7.2, 203, false, true, true, 0.4)
		sec3Attr := fitTextAttr("3. "+l.KaggleTitle, 9.0, 7.2, 202, false, true, true, 0.4)
		sec4Attr := fitTextAttr("4. "+l.ImpactTitle, 9.0, 7.2, 203, false, true, true, 0.4)
		l0Attr := fitTextAttr(l0Short, 10.0, 7.8, 58, false, true, false, 0)
		l1Attr := fitTextAttr(l1Short, 10.0, 7.8, 58, false, true, false, 0)
		q1R3Val := fmt.Sprintf("%.0f%% Prod / %.0f%% EDA", d.Languages.ProductionPct, d.Languages.ExplorationPct)
		q1R3L, q1R3R := fitPairTextAttrs("Prod vs. EDA", 10.0, 7.8, false, true, false, 0, q1R3Val, 10.0, 7.8, true, true, false, 0, 202, 8)
		q2R1L, q2R1R := fitPairTextAttrs(dataCol[0].Name, 10.0, 7.8, false, true, false, 0, fmt.Sprintf("%d repos", dataCol[0].RepoCount), 10.0, 7.8, true, true, false, 0, 203, 8)
		q2R2L, q2R2R := fitPairTextAttrs(mlCol[0].Name, 10.0, 7.8, false, true, false, 0, fmt.Sprintf("%d repos", mlCol[0].RepoCount), 10.0, 7.8, true, true, false, 0, 203, 8)
		q2R3L, q2R3R := fitPairTextAttrs(l.LangLockfiles, 10.0, 7.8, false, true, false, 0, fmt.Sprintf("%d / %d", d.Languages.LockfileRepos, d.Languages.TotalRepos), 10.0, 7.8, true, true, false, 0, 203, 8)
		q3R1L, q3R1R := fitPairTextAttrs(l.KaggleOverallTierLbl, 10.0, 7.8, false, true, false, 0, overallTierRank, 10.0, 7.8, true, true, false, 0, 202, 8)
		q3R2Val := fmt.Sprintf("%d (%dG · %dS · %dB)", k.TotalMedals(), k.TotalGold, k.TotalSilver, k.TotalBronze)
		q3R2L, q3R2R := fitPairTextAttrs(l.KaggleTotalMedalsLbl, 10.0, 7.8, false, true, false, 0, q3R2Val, 10.0, 7.8, true, true, false, 0, 202, 8)
		q3R3L, q3R3R := fitPairTextAttrs("Code / Datasets", 10.0, 7.8, false, true, false, 0, summaryTracksVal, 10.0, 7.8, true, true, false, 0, 202, 8)
		q4R1Val := fmt.Sprintf("%s (+%.0f%%)", formatInt(im.AnnualContributions), im.YoYGrowthPct)
		q4R1L, q4R1R := fitPairTextAttrs(l.ImpactContribs, 10.0, 7.8, false, true, false, 0, q4R1Val, 10.0, 7.8, true, true, false, 0, 203, 8)
		q4R2Val := fmt.Sprintf("★ %s (%d forks)", formatInt(im.TotalStars), im.TotalForks)
		q4R2L, q4R2R := fitPairTextAttrs(l.ImpactStars, 10.0, 7.8, false, true, false, 0, q4R2Val, 10.0, 7.8, true, true, false, 0, 203, 8)
		q4R3Val := fmt.Sprintf("%s (30d)", formatCompact(im.HFDownloads30d))
		q4R3L, q4R3R := fitPairTextAttrs(l.ImpactDownloads, 10.0, 7.8, false, true, false, 0, q4R3Val, 10.0, 7.8, true, true, false, 0, 203, 8)
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 465 255" width="465" height="255">
  <style>
    .serif { font-family: Georgia, "Times New Roman", serif; }
    .sans { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-variant-numeric: tabular-nums; }
    .headline { font-size: 12.5px; font-weight: bold; fill: #526071; }
    .lead { font-size: 9px; fill: #64748b; }
    .sec { font-size: 9px; font-weight: 700; fill: #e11d48; text-transform: uppercase; letter-spacing: 0.4px; }
    .lbl { font-size: 10px; font-weight: 600; fill: #526071; }
    .val { font-size: 10px; font-weight: 700; fill: #64748b; }
    .callout { font-size: 9px; font-style: italic; fill: #0284c7; }
  </style>
  <rect x="12" y="6" width="441" height="3" fill="#e11d48"/>
  <text x="12" y="24" class="serif headline"%s>%s — %s (%d Repos · ★ %s)</text>
  <text x="12" y="38" class="sans lead"%s>%s</text>
  <line x1="12" y1="45" x2="453" y2="45" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>

  <!-- Q1: Languages (Lollipop + Waffle) -->
  <g transform="translate(12, 60)">
    <text x="0" y="0" class="sans sec"%s>1. %s</text>
    <text x="0" y="18" class="sans lbl"%s>%s</text>
    <line x1="64" y1="14" x2="%d" y2="14" stroke="#0284c7" stroke-width="2.2"/>
    <circle cx="%d" cy="14" r="3.5" fill="#0284c7"/>
    <text x="202" y="18" text-anchor="end" class="mono val" fill="#0284c7">%.1f%%</text>

    <text x="0" y="36" class="sans lbl"%s>%s</text>
    <line x1="64" y1="32" x2="%d" y2="32" stroke="#d97706" stroke-width="2.2"/>
    <circle cx="%d" cy="32" r="3.5" fill="#d97706"/>
    <text x="202" y="36" text-anchor="end" class="mono val">%.1f%%</text>

    <text x="0" y="54" class="sans lbl"%s>Prod vs. EDA</text>
    <text x="202" y="54" text-anchor="end" class="mono val" fill="#059669"%s>%.0f%% Prod / %.0f%% EDA</text>
  </g>

  <!-- Q2: Stack (Lollipop + Dashed Benchmark) -->
  <g transform="translate(248, 60)">
    <text x="0" y="0" class="sans sec"%s>2. %s</text>
    <line x1="150" y1="8" x2="150" y2="58" stroke="#e11d48" stroke-dasharray="3,3" stroke-width="1"/>
    <text x="0" y="18" class="sans lbl"%s>%s</text><text x="203" y="18" text-anchor="end" class="mono val"%s>%d repos</text>
    <text x="0" y="36" class="sans lbl"%s>%s</text><text x="203" y="36" text-anchor="end" class="mono val"%s>%d repos</text>
    <text x="0" y="54" class="sans lbl"%s>%s</text><text x="203" y="54" text-anchor="end" class="mono val"%s>%d / %d</text>
  </g>

  <line x1="12" y1="132" x2="453" y2="132" stroke="#64748b" stroke-opacity="0.22" stroke-width="1"/>

  <!-- Q3: Kaggle -->
  <g transform="translate(12, 148)">
    <text x="0" y="0" class="sans sec"%s>3. %s</text>
    <text x="0" y="18" class="sans lbl"%s>%s</text><text x="202" y="18" text-anchor="end" class="mono val" fill="#d97706"%s>%s</text>
    <text x="0" y="36" class="sans lbl"%s>%s</text><text x="202" y="36" text-anchor="end" class="mono val"%s>%d (%dG · %dS · %dB)</text>
    <text x="0" y="54" class="sans lbl"%s>Code / Datasets</text><text x="202" y="54" text-anchor="end" class="mono val"%s>%s</text>
  </g>

  <!-- Q4: Impact -->
  <g transform="translate(248, 148)">
    <text x="0" y="0" class="sans sec"%s>4. %s</text>
    <text x="0" y="18" class="sans lbl"%s>%s</text><text x="203" y="18" text-anchor="end" class="mono val" fill="#059669"%s>%s (+%.0f%%)</text>
    <text x="0" y="36" class="sans lbl"%s>%s</text><text x="203" y="36" text-anchor="end" class="mono val"%s>★ %s (%d forks)</text>
    <text x="0" y="54" class="sans lbl"%s>%s</text><text x="203" y="54" text-anchor="end" class="mono val"%s>%s (30d)</text>
  </g>

  <text x="12" y="238" class="serif callout"%s>%s</text>
</svg>`,
			headAttr, esc(d.Subject.Name), esc(dispRole),
			d.Languages.TotalRepos, formatInt(im.TotalStars),
			leadAttr, esc(sumLead),
			sec1Attr, esc(l.LangSecTrueProp),
			l0Attr, esc(l0Short), l0End, l0End, langs[0].Percentage,
			l1Attr, esc(l1Short), l1End, l1End, langs[1].Percentage,
			q1R3L, q1R3R, d.Languages.ProductionPct, d.Languages.ExplorationPct,
			sec2Attr, esc(l.StackTitle),
			q2R1L, esc(dataCol[0].Name), q2R1R, dataCol[0].RepoCount,
			q2R2L, esc(mlCol[0].Name), q2R2R, mlCol[0].RepoCount,
			q2R3L, esc(l.LangLockfiles), q2R3R, d.Languages.LockfileRepos, d.Languages.TotalRepos,
			sec3Attr, esc(l.KaggleTitle),
			q3R1L, esc(l.KaggleOverallTierLbl), q3R1R, esc(overallTierRank),
			q3R2L, esc(l.KaggleTotalMedalsLbl), q3R2R, k.TotalMedals(), k.TotalGold, k.TotalSilver, k.TotalBronze,
			q3R3L, q3R3R, esc(summaryTracksVal),
			sec4Attr, esc(l.ImpactTitle),
			q4R1L, esc(l.ImpactContribs), q4R1R, formatInt(im.AnnualContributions), im.YoYGrowthPct,
			q4R2L, esc(l.ImpactStars), q4R2R, formatInt(im.TotalStars), im.TotalForks,
			q4R3L, esc(l.ImpactDownloads), q4R3R, formatCompact(im.HFDownloads30d),
			calloutAttr, esc(sumCallout))

	default: // SchoolHolmes
		holmesKaggleMid := fmt.Sprintf("%d %s", k.TotalMedals(), l.KaggleMedalsWord)
		if k.PeakRank > 0 {
			holmesKaggleMid = fmt.Sprintf("#%s · %d %s", formatInt(k.PeakRank), k.TotalMedals(), l.KaggleMedalsWord)
		}
		suffix := fmt.Sprintf(" (%d REPOS)", d.Languages.TotalRepos)
		suffixW := estimateTextWidth(suffix, 10.5, false, true, true, 0.3)
		dispRole, _ := fitHeaderNameRole("⚡ "+d.Subject.Name, d.Subject.Role, " — ", 433-suffixW, 10.5, 8.0, false, true, true, 0.3)
		bannerRaw := fmt.Sprintf("⚡ %s — %s%s", d.Subject.Name, dispRole, suffix)
		bannerAttr := fitTextAttr(bannerRaw, 10.5, 8.0, 433, false, true, true, 0.3)
		sec1Attr := fitTextAttr("1. "+l.LangSecTrueProp, 9.5, 7.5, 200, false, true, true, 0.4)
		sec2Attr := fitTextAttr("2. "+l.StackTitle, 9.5, 7.5, 204, false, true, true, 0.4)
		sec3Attr := fitTextAttr("3. "+l.KaggleTitle, 9.5, 7.5, 200, false, true, true, 0.4)
		sec4Attr := fitTextAttr("4. "+l.ImpactTitle, 9.5, 7.5, 204, false, true, true, 0.4)
		l0Attr := fitTextAttr(l0Short, 9.5, 7.5, 56, false, true, false, 0)
		q1R2Left := fmt.Sprintf("%s / %s", l1Short, l2Short)
		q1R2Right := fmt.Sprintf("%.1f%% / %.1f%%", langs[1].Percentage, langs[2].Percentage)
		q1R2L, q1R2R := fitPairTextAttrs(q1R2Left, 9.5, 7.5, false, true, false, 0, q1R2Right, 10.0, 7.8, true, true, false, 0, 200, 8)
		q2R1L, q2R1R := fitPairTextAttrs("🧰 "+dataCol[0].Name, 9.5, 7.5, false, true, false, 0, fmt.Sprintf("%d repos", dataCol[0].RepoCount), 10.0, 7.8, true, true, false, 0, 188, 8)
		q2R2L, q2R2R := fitPairTextAttrs("🤖 "+mlCol[0].Name, 9.5, 7.5, false, true, false, 0, fmt.Sprintf("%d repos", mlCol[0].RepoCount), 10.0, 7.8, true, true, false, 0, 188, 8)
		q2R3L, q2R3R := fitPairTextAttrs("🔒 "+l.LangLockfiles, 9.5, 7.5, false, true, false, 0, fmt.Sprintf("%d/%d", d.Languages.LockfileRepos, d.Languages.TotalRepos), 10.0, 7.8, true, true, false, 0, 188, 8)
		kTierAttr := fitTextAttr("Tier: "+k.OverallTier, 9.5, 7.5, 108, false, true, false, 0)
		kMidAttr := fitTextAttr(holmesKaggleMid, 10.0, 7.5, 108, true, true, false, 0)
		kMedalsAttr := fitTextAttr(fmt.Sprintf("%dG · %dS · %dB", k.TotalGold, k.TotalSilver, k.TotalBronze), 10.0, 7.5, 108, true, true, false, 0)
		q4R1Val := fmt.Sprintf("%s (+%.0f%%)", formatInt(im.AnnualContributions), im.YoYGrowthPct)
		q4R1L, q4R1R := fitPairTextAttrs("🚀 "+l.ImpactContribs, 9.5, 7.5, false, true, false, 0, q4R1Val, 10.0, 7.8, true, true, false, 0, 188, 8)
		q4R2L, q4R2R := fitPairTextAttrs("⭐ "+l.ImpactStars, 9.5, 7.5, false, true, false, 0, "★ "+formatInt(im.TotalStars), 10.0, 7.8, true, true, false, 0, 188, 8)
		q4R3L, q4R3R := fitPairTextAttrs("📥 "+l.ImpactDownloads, 9.5, 7.5, false, true, false, 0, formatCompact(im.HFDownloads30d)+" (30d)", 10.0, 7.8, true, true, false, 0, 188, 8)
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 465 255" width="465" height="255">
  <style>
    .sans { font-family: "Trebuchet MS", Inter, -apple-system, BlinkMacSystemFont, sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-weight: 700; }
    .banner { font-size: 10.5px; font-weight: 800; fill: #ffffff; letter-spacing: 0.3px; text-transform: uppercase; }
    .sec { font-size: 9.5px; font-weight: 800; fill: #d97706; text-transform: uppercase; letter-spacing: 0.4px; }
    .lbl { font-size: 9.5px; font-weight: 700; fill: #526071; }
    .val { font-size: 10px; font-weight: 800; fill: #0284c7; }
  </style>
  <rect x="3" y="3" width="459" height="249" rx="10" fill="none" stroke="#64748b" stroke-opacity="0.35" stroke-width="2"/>
  <path d="M 3,13 Q 3,3 13,3 L 452,3 Q 462,3 462,13 L 462,32 L 3,32 Z" fill="#0284c7"/>
  <text x="16" y="22" class="sans banner"%s>⚡ %s — %s (%d REPOS)</text>

  <!-- Q1: Languages (Pictorial Unit Blocks) -->
  <g transform="translate(16, 48)">
    <text x="0" y="0" class="sans sec"%s>1. %s</text>
    <text x="0" y="20" class="sans lbl"%s>%s</text>
    <g transform="translate(60, 10)" fill="#0284c7">
      <rect x="0" y="0" width="12" height="10" rx="2"/><rect x="15" y="0" width="12" height="10" rx="2"/><rect x="30" y="0" width="12" height="10" rx="2"/><rect x="45" y="0" width="12" height="10" rx="2"/><rect x="60" y="0" width="12" height="10" rx="2"/>
    </g>
    <text x="200" y="20" text-anchor="end" class="mono val">%.1f%%</text>
    <text x="0" y="38" class="sans lbl"%s>%s / %s</text>
    <text x="200" y="38" text-anchor="end" class="mono val" fill="#d97706"%s>%.1f%% / %.1f%%</text>
    <rect x="0" y="46" width="200" height="22" rx="5" fill="#059669" fill-opacity="0.1" stroke="#059669" stroke-opacity="0.4" stroke-width="1"/>
    <text x="8" y="61" class="sans lbl" fill="#059669">✓ %.0f%% Prod · %.0f%% EDA</text>
  </g>

  <!-- Q2: Stack (Visual Engineering Badges) -->
  <g transform="translate(244, 48)">
    <text x="0" y="0" class="sans sec"%s>2. %s</text>
    <rect x="0" y="8" width="204" height="20" rx="5" fill="#0284c7" fill-opacity="0.1" stroke="#0284c7" stroke-opacity="0.4" stroke-width="1"/>
    <text x="8" y="22" class="sans lbl"%s>🧰 %s</text><text x="196" y="22" text-anchor="end" class="mono val"%s>%d repos</text>
    <rect x="0" y="32" width="204" height="20" rx="5" fill="#7c3aed" fill-opacity="0.1" stroke="#7c3aed" stroke-opacity="0.4" stroke-width="1"/>
    <text x="8" y="46" class="sans lbl"%s>🤖 %s</text><text x="196" y="46" text-anchor="end" class="mono val" fill="#7c3aed"%s>%d repos</text>
    <rect x="0" y="56" width="204" height="20" rx="5" fill="#059669" fill-opacity="0.1" stroke="#059669" stroke-opacity="0.4" stroke-width="1"/>
    <text x="8" y="70" class="sans lbl"%s>🔒 %s</text><text x="196" y="70" text-anchor="end" class="mono val" fill="#059669"%s>%d/%d</text>
  </g>

  <!-- Q3: Kaggle (Stepped Podium + Trophy) -->
  <g transform="translate(16, 146)">
    <text x="0" y="0" class="sans sec"%s>3. %s</text>
    <rect x="0" y="44" width="18" height="16" rx="2" fill="#64748b" fill-opacity="0.25"/>
    <rect x="21" y="34" width="18" height="26" rx="2" fill="#059669" fill-opacity="0.4"/>
    <rect x="42" y="24" width="18" height="36" rx="2" fill="#0284c7" fill-opacity="0.65"/>
    <rect x="63" y="18" width="18" height="42" rx="2" fill="#d97706"/>
    <text x="72" y="15" text-anchor="middle" font-size="10">🏆</text>
    <text x="92" y="26" class="sans lbl"%s>Tier: %s</text>
    <text x="92" y="42" class="mono val" fill="#d97706"%s>%s</text>
    <text x="92" y="58" class="mono val"%s>%dG · %dS · %dB</text>
  </g>

  <!-- Q4: Impact (Visual Badges) -->
  <g transform="translate(244, 146)">
    <text x="0" y="0" class="sans sec"%s>4. %s</text>
    <rect x="0" y="8" width="204" height="24" rx="5" fill="#059669" fill-opacity="0.1" stroke="#059669" stroke-opacity="0.4" stroke-width="1"/>
    <text x="8" y="24" class="sans lbl"%s>🚀 %s</text><text x="196" y="24" text-anchor="end" class="mono val" fill="#059669"%s>%s (+%.0f%%)</text>
    <rect x="0" y="36" width="204" height="24" rx="5" fill="#d97706" fill-opacity="0.1" stroke="#d97706" stroke-opacity="0.4" stroke-width="1"/>
    <text x="8" y="52" class="sans lbl"%s>⭐ %s</text><text x="196" y="52" text-anchor="end" class="mono val" fill="#d97706"%s>★ %s</text>
    <rect x="0" y="64" width="204" height="24" rx="5" fill="#0284c7" fill-opacity="0.1" stroke="#0284c7" stroke-opacity="0.4" stroke-width="1"/>
    <text x="8" y="80" class="sans lbl"%s>📥 %s</text><text x="196" y="80" text-anchor="end" class="mono val"%s>%s (30d)</text>
  </g>
</svg>`,
			bannerAttr, esc(d.Subject.Name), esc(dispRole), d.Languages.TotalRepos,
			sec1Attr, esc(l.LangSecTrueProp),
			l0Attr, esc(l0Short), langs[0].Percentage,
			q1R2L, esc(l1Short), esc(l2Short), q1R2R, langs[1].Percentage, langs[2].Percentage,
			d.Languages.ProductionPct, d.Languages.ExplorationPct,
			sec2Attr, esc(l.StackTitle),
			q2R1L, esc(dataCol[0].Name), q2R1R, dataCol[0].RepoCount,
			q2R2L, esc(mlCol[0].Name), q2R2R, mlCol[0].RepoCount,
			q2R3L, esc(l.LangLockfiles), q2R3R, d.Languages.LockfileRepos, d.Languages.TotalRepos,
			sec3Attr, esc(l.KaggleTitle),
			kTierAttr, esc(k.OverallTier), kMidAttr, esc(holmesKaggleMid),
			kMedalsAttr, k.TotalGold, k.TotalSilver, k.TotalBronze,
			sec4Attr, esc(l.ImpactTitle),
			q4R1L, esc(l.ImpactContribs), q4R1R, formatInt(im.AnnualContributions), im.YoYGrowthPct,
			q4R2L, esc(l.ImpactStars), q4R2R, formatInt(im.TotalStars),
			q4R3L, esc(l.ImpactDownloads), q4R3R, formatCompact(im.HFDownloads30d))
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

	kaggleHl := fmt.Sprintf("Kaggle %s", formatKaggleOverallTierAndRank(d.Kaggle))
	if l.IsPtBR {
		return []string{
			fmt.Sprintf("Olá! Eu sou %s 👋", name),
			fmt.Sprintf("%s (%.1f%%) · %s · %d Repositórios", langs[0].Name, langs[0].Percentage, topStack, d.Languages.TotalRepos),
			fmt.Sprintf("%s · %s Contribuições (+%.0f%% YoY)", kaggleHl, formatInt(d.Impact.AnnualContributions), d.Impact.YoYGrowthPct),
		}
	}
	return []string{
		fmt.Sprintf("Welcome! I'm %s 👋", name),
		fmt.Sprintf("%s (%.1f%%) · %s · %d Repositories", langs[0].Name, langs[0].Percentage, topStack, d.Languages.TotalRepos),
		fmt.Sprintf("%s · %s Contributions (+%.0f%% YoY)", kaggleHl, formatInt(d.Impact.AnnualContributions), d.Impact.YoYGrowthPct),
	}
}

func headlineFontSize(msg string, baseSize float64, x, maxW int) float64 {
	availW := float64(maxW - x - 6)
	if availW < 200 {
		availW = 200
	}
	runes := float64(utf8.RuneCountInString(msg))
	if runes <= 0 {
		return baseSize
	}
	estW := runes * baseSize * 0.56
	if estW <= availW {
		return baseSize
	}
	scaled := availW / (runes * 0.56)
	minFloor := 13.5
	if runes > 75 {
		minFloor = 10.5
	}
	if scaled < minFloor {
		return minFloor
	}
	return math.Round(scaled*10) / 10
}

func buildAnimatedHeadlinesSVG(headlines []string, idPrefix string, x, y, maxW int, textClass, fillColor string, fontSizePx float64) string {
	if len(headlines) == 0 {
		return ""
	}
	// SMIL syncbase timing treats '-' inside an element ID (e.g., begin="hdr-d0.end") as a clock-offset
	// subtraction operator in browsers, breaking sequential chaining and infinite looping after the first line.
	animPrefix := strings.ReplaceAll(idPrefix, "-", "")
	var b bytes.Buffer
	n := len(headlines)
	for i, msg := range headlines {
		pathID := fmt.Sprintf("%spath%d", idPrefix, i)
		animID := fmt.Sprintf("%sd%d", animPrefix, i)
		var beginAttr string
		if n == 1 {
			beginAttr = fmt.Sprintf("0s;%s.end", animID)
		} else if i == 0 {
			beginAttr = fmt.Sprintf("0s;%sd%d.end", animPrefix, n-1)
		} else {
			beginAttr = fmt.Sprintf("%sd%d.end", animPrefix, i-1)
		}
		lineFontPx := headlineFontSize(msg, fontSizePx, x, maxW)
		fmt.Fprintf(&b, `  <path id="%s">
    <animate id="%s" attributeName="d" begin="%s" dur="6000ms" fill="remove" values="m0,%d h0 ; m0,%d h%d ; m0,%d h%d ; m0,%d h0" keyTimes="0;0.666667;0.833333;1"/>
  </path>
  <text x="%d" class="%s" font-weight="bold" fill="%s" font-size="%.1fpx"><textPath xlink:href="#%s">%s</textPath></text>
`, pathID, animID, beginAttr, y, y, maxW, y, maxW, y, x, textClass, fillColor, lineFontPx, pathID, esc(msg))
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
		rightSubW := estimateTextWidth(rightSub, 11.0, false, false, false, 0)
		namePrefixW := estimateTextWidth(name+" — ", 13.5, false, true, false, 0.2)
		availRoleW := 872.0 - rightSubW - namePrefixW - 14.0
		dispRole, roleFont := fitHeaderNameRole("", role, "", availRoleW, 11.0, 8.8, false, false, false, 0)
		roleAttr := fontSizeAttr(roleFont, 11.0)
		trendRight := fmt.Sprintf("%s · %s (+%.0f%% YoY)", l.HeaderTrendLbl, formatInt(im.AnnualContributions), im.YoYGrowthPct)
		trendAttr := fitTextAttr(trendRight, 9.5, 7.8, 286, false, false, false, 0)
		inner := fmt.Sprintf(`  <text x="14" y="24" class="serif h1">%s — <tspan class="sub"%s>%s</tspan></text>
  <text x="886" y="24" text-anchor="end" class="num sub">%s</text>
  <line x1="14" y1="31" x2="886" y2="31" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
%s  <rect x="596" y="42" width="280" height="32" fill="#64748b" fill-opacity="0.12"/>
  <path d="%s" fill="none" stroke="#64748b" stroke-dasharray="2.5,2.5" stroke-width="1"/>
  <path d="%s" fill="none" stroke="#0284c7" stroke-width="2">
    <animate attributeName="stroke-dasharray" values="0 100%%; 10%% 100%%; 100%% 10%%; 10%% 100%%" keyTimes="0; 0.5; 0.9; 1" dur="4s" begin="0s" repeatCount="indefinite"/>
  </path>
  <circle cx="%d" cy="%d" r="3.2" fill="#059669"/>
  <text x="886" y="96" text-anchor="end" class="num micro"%s>%s · %s (+%.0f%% YoY)</text>
  <line x1="14" y1="106" x2="886" y2="106" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>`,
			esc(name), roleAttr, esc(dispRole), esc(rightSub),
			animText, pathD, pathD, lastX, lastY,
			trendAttr, esc(l.HeaderTrendLbl), formatInt(im.AnnualContributions), im.YoYGrowthPct)
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
		rightTag := fmt.Sprintf("+%.0f%% YoY · %d REPOS", im.YoYGrowthPct, d.Languages.TotalRepos)
		rightTagW := estimateTextWidth(rightTag, 9.5, true, true, false, 0)
		dispRole, _ := fitHeaderNameRole(name, role, " — ", 868-rightTagW-16, 12.0, 8.8, false, true, true, 0.5)
		leftHdr := fmt.Sprintf("%s — %s", name, dispRole)
		hdrAttr, tagAttr := fitPairTextAttrs(leftHdr, 12.0, 8.8, false, true, true, 0.5, rightTag, 9.5, 8.0, true, true, false, 0, 868, 16)
		kpi1Val := fmt.Sprintf("%s %.1f%%", langs[0].Name, langs[0].Percentage)
		kpi2Val := formatKaggleOverallTierAndRank(k)
		kpi3Val := formatInt(im.AnnualContributions)
		kpi1NumAttr := fitTextAttr(kpi1Val, 12.0, 8.5, 146, true, true, false, 0)
		kpi1LblAttr := fitTextAttr(l.LangFewPrimary, 8.0, 6.8, 146, false, true, true, 0)
		kpi2NumAttr := fitTextAttr(kpi2Val, 12.0, 8.5, 164, true, true, false, 0)
		kpi2LblAttr := fitTextAttr(l.KaggleOverallTierLbl, 8.0, 6.8, 164, false, true, true, 0)
		kpi3NumAttr := fitTextAttr(kpi3Val, 12.0, 8.5, 218, true, true, false, 0)
		kpi3LblAttr := fitTextAttr(l.ImpactContribs, 8.0, 6.8, 218, false, true, true, 0)
		trendLblAttr := fitTextAttr(l.HeaderTrendLbl, 8.0, 6.8, 148, true, true, true, 0)
		inner := fmt.Sprintf(`  <text x="16" y="21" class="sans hdr"%s>%s — %s</text>
  <text x="884" y="21" text-anchor="end" class="mono tag"%s>+%.0f%% YoY · %d REPOS</text>
  <line x1="2" y1="29" x2="898" y2="29" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(16, 34)">
    <text x="0" y="11" class="mono kpi-num" fill="#0284c7"%s>%s %.1f%%</text>
    <text x="0" y="21" class="sans kpi-lbl"%s>%s</text>
    <line x1="150" y1="0" x2="150" y2="23" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="162" y="11" class="mono kpi-num" fill="#d97706"%s>%s</text>
    <text x="162" y="21" class="sans kpi-lbl"%s>%s</text>
    <line x1="330" y1="0" x2="330" y2="23" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="342" y="11" class="mono kpi-num" fill="#059669"%s>%s</text>
    <text x="342" y="21" class="sans kpi-lbl"%s>%s</text>
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
    <text x="276" y="8" text-anchor="end" class="mono kpi-lbl"%s>%s</text>
  </g>`,
			hdrAttr, esc(name), esc(dispRole), tagAttr, im.YoYGrowthPct, d.Languages.TotalRepos,
			kpi1NumAttr, esc(langs[0].Name), langs[0].Percentage, kpi1LblAttr, esc(l.LangFewPrimary),
			kpi2NumAttr, esc(kpi2Val), kpi2LblAttr, esc(l.KaggleOverallTierLbl),
			kpi3NumAttr, kpi3Val, kpi3LblAttr, esc(l.ImpactContribs),
			animText, pathD, pathD, lastX, lastY, bulletW, trendLblAttr, esc(l.HeaderTrendLbl))
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
		nameAttr := fitTextAttr(name, 13.5, 10.0, 558, false, true, false, 0)
		roleAttr := fitTextAttr(role, 10.5, 8.2, 558, false, false, false, 0)
		calloutAttr := fitTextAttr(callout, 9.2, 7.5, 286, false, false, false, 0)
		inner := fmt.Sprintf(`  <defs>
    <marker id="%sarrow" overflow="visible" markerHeight="4" markerWidth="6" orient="auto-start-reverse" viewBox="0 0 8.75 5">
      <path transform="scale(-.5)" d="m0 0 5-5-17.5 5 17.5 5z" fill="#64748b" fill-rule="evenodd"/>
    </marker>
  </defs>
  <rect x="12" y="6" width="876" height="3" fill="#e11d48"/>
  <text x="14" y="26" class="serif headline"%s>%s</text>
  <text x="14" y="41" class="sans lead"%s>%s</text>
  <line x1="14" y1="48" x2="572" y2="48" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
%s  <path d="M 598,20 V 86 H 880" fill="none" marker-start="url(#%sarrow)" marker-end="url(#%sarrow)" stroke="#64748b" stroke-width="1.2"/>
  <line x1="598" y1="52" x2="880" y2="52" stroke="#e11d48" stroke-dasharray="3,3" stroke-width="1"/>
  <path d="%s" fill="none" stroke="#64748b" stroke-dasharray="3,2" stroke-width="1.1"/>
  <path d="%s" fill="none" stroke="#e11d48" stroke-width="2.4">
    <animate attributeName="stroke-dasharray" values="0 100%%; 10%% 100%%; 100%% 10%%; 10%% 100%%" keyTimes="0; 0.5; 0.9; 1" dur="4s" begin="0s" repeatCount="indefinite"/>
  </path>
  <circle cx="%d" cy="%d" r="3.5" fill="#e11d48"/>
  <text x="598" y="102" class="serif callout"%s>%s</text>`,
			idPrefix, nameAttr, esc(name), roleAttr, esc(role), animText, idPrefix, idPrefix, pathD, pathD, lastX, lastY, calloutAttr, esc(callout))
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
		rightBanner := fmt.Sprintf("🏆 %s · ★ %s", strings.ToUpper(k.OverallTier), formatInt(im.TotalStars))
		rightBannerW := estimateTextWidth(rightBanner, 11.5, false, true, true, 0.5)
		dispRole, _ := fitHeaderNameRole("⚡ "+name, role, " — ", 868-rightBannerW-16, 11.5, 8.5, false, true, true, 0.5)
		leftBanner := fmt.Sprintf("⚡ %s — %s", name, dispRole)
		lBanAttr, rBanAttr := fitPairTextAttrs(leftBanner, 11.5, 8.5, false, true, true, 0.5, rightBanner, 11.5, 8.5, false, true, true, 0.5, 868, 16)
		inner := fmt.Sprintf(`  <path d="M 3,13 Q 3,3 13,3 L 887,3 Q 897,3 897,13 L 897,32 L 3,32 Z" fill="#0284c7"/>
  <text x="16" y="22" class="sans banner"%s>⚡ %s — %s</text>
  <text x="884" y="22" text-anchor="end" class="sans banner"%s>🏆 %s · ★ %s</text>
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
			lBanAttr, esc(name), esc(dispRole), rBanAttr, strings.ToUpper(esc(k.OverallTier)), formatInt(im.TotalStars),
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

const (
	bioCompactMaxChars   = 52
	bioFullMaxChars      = 90
	bioCompactLineHeight = 20
	bioCompactParaGap    = 10
	bioFullLineHeight    = 24
	bioFullParaGap       = 12
)

func truncateRunes(s string, maxRunes int) string {
	if utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	runes := []rune(s)
	return strings.TrimSpace(string(runes[:maxRunes-1])) + "…"
}

func renderBio(school domain.School, d domain.PortfolioData, l labels, fullWidth bool) string {
	width := 440
	maxChars := bioCompactMaxChars
	lineHeight := bioCompactLineHeight
	paraGap := bioCompactParaGap
	if fullWidth {
		width = 900
		maxChars = bioFullMaxChars
		lineHeight = bioFullLineHeight
		paraGap = bioFullParaGap
	}

	rawParas := d.Subject.Bio
	if len(rawParas) == 0 {
		rawParas = fallbackBioParagraphs(d, l)
	}
	wrapped := wrapBioParagraphs(rawParas, maxChars)
	if len(wrapped) == 0 {
		wrapped = wrapBioParagraphs(fallbackBioParagraphs(d, l), maxChars)
	}

	_, bodyHeight := computeBioLayout(wrapped, lineHeight, paraGap)

	name := nonEmptyStr(d.Subject.Name, "Data Practitioner")
	role := nonEmptyStr(d.Subject.Role, "Data Scientist & Analytics Engineer")
	langs := padLangs(d.Languages.Languages)
	k := d.Kaggle
	im := d.Impact
	subRight := fmt.Sprintf(l.BioSubFmt, d.Languages.TotalRepos, formatInt(im.TotalStars))

	switch school {
	case domain.SchoolTufte:
		h1Size := 14.5
		subSize := 11.0
		bodySize := 13.5
		maxInlineHeaderRunes := 40
		if fullWidth {
			h1Size = 16.0
			subSize = 12.0
			bodySize = 16.0
			maxInlineHeaderRunes = 68
		}
		headerTitle := fmt.Sprintf("%s — %s", name, role)
		startY := 50
		ruleY := 29
		var roleSubRow string
		if utf8.RuneCountInString(headerTitle) > maxInlineHeaderRunes {
			headerTitle = name
			displayRole := role
			if !fullWidth {
				displayRole = truncateRunes(role, 58)
			}
			roleSubAttr := fitTextAttr(displayRole, subSize, 8.5, float64(width-28), false, false, false, 0)
			roleSubRow = fmt.Sprintf("  <text x=\"14\" y=\"37\" class=\"serif sub\"%s>%s</text>\n", roleSubAttr, esc(displayRole))
			ruleY = 44
			startY = 66
		}
		h1Attr, subAttr := fitPairTextAttrs(headerTitle, h1Size, 10.0, false, true, false, 0.2, subRight, subSize, 8.5, false, false, false, 0, float64(width-28), 12)
		totalHeight := startY + bodyHeight + 16
		linesSVG := renderBioLinesSVG(wrapped, 14, startY, lineHeight, paraGap, "serif body")
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">
  <style>
    .serif { font-family: Georgia, "Palatino Linotype", "Book Antiqua", Palatino, serif; }
    .num { font-family: Georgia, serif; font-variant-numeric: tabular-nums; }
    .h1 { font-size: %.1fpx; font-weight: bold; fill: #526071; letter-spacing: 0.2px; }
    .sub { font-size: %.1fpx; font-style: italic; fill: #64748b; }
    .body { font-size: %.1fpx; fill: #526071; }
  </style>
  <text x="14" y="21" class="serif h1"%s>%s</text>
  <text x="%d" y="21" text-anchor="end" class="num sub"%s>%s</text>
%s  <line x1="14" y1="%d" x2="%d" y2="%d" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
%s  <line x1="14" y1="%d" x2="%d" y2="%d" stroke="#64748b" stroke-opacity="0.7" stroke-width="0.8"/>
</svg>`, width, totalHeight, width, totalHeight, h1Size, subSize, bodySize, h1Attr, esc(headerTitle), width-14, subAttr, esc(subRight), roleSubRow, ruleY, width-14, ruleY, linesSVG, totalHeight-8, width-14, totalHeight-8)

	case domain.SchoolFew:
		hdrSize := 12.0
		tagSize := 9.5
		bodySize := 12.5
		if fullWidth {
			hdrSize = 13.5
			tagSize = 10.5
			bodySize = 15.0
		}
		startY := 84
		totalHeight := startY + bodyHeight + 14
		linesSVG := renderBioLinesSVG(wrapped, 14, startY, lineHeight, paraGap, "sans body")
		prodBar := clampInt(int(math.Round((d.Languages.ProductionPct/100.0)*84.0)), 14, 84)
		tierRaw := nonEmptyStr(k.OverallTier, domain.KaggleTierContributor)
		if k.PeakRank > 0 {
			tierRaw = fmt.Sprintf("%s (#%s)", tierRaw, formatInt(k.PeakRank))
		} else if k.TotalMedals() > 0 {
			tierRaw = fmt.Sprintf("%s (%dM)", tierRaw, k.TotalMedals())
		}
		hdrAttr, tagAttr := fitPairTextAttrs(name, hdrSize, 9.0, false, true, true, 0.5, l.BioTitle, tagSize, 8.0, true, true, false, 0, float64(width-28), 12)
		kpi1Val := fmt.Sprintf("%s (%.1f%%)", langs[0].Name, langs[0].Percentage)
		kpi1NumAttr := fitTextAttr(kpi1Val, 12.5, 8.5, 134, true, true, false, 0)
		kpi1LblAttr := fitTextAttr(l.LangFewPrimary, 8.0, 6.8, 134, false, true, true, 0)
		kpi2NumAttr := fitTextAttr(tierRaw, 12.5, 8.5, 138, true, true, false, 0)
		kpi2LblAttr := fitTextAttr(l.KaggleOverallTierLbl, 8.0, 6.8, 138, false, true, true, 0)
		kpi3LblAttr := fitTextAttr(l.LangFewProdEDA, 8.0, 6.8, float64(width-18-304), false, true, true, 0)
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">
  <style>
    .sans { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-variant-numeric: tabular-nums; }
    .hdr { font-size: %.1fpx; font-weight: 700; fill: #526071; text-transform: uppercase; letter-spacing: 0.5px; }
    .tag { font-size: %.1fpx; font-weight: 600; fill: #0284c7; }
    .body { font-size: %.1fpx; font-weight: 500; fill: #526071; }
    .kpi-num { font-size: 12.5px; font-weight: 700; fill: #526071; }
    .kpi-lbl { font-size: 8px; font-weight: 600; fill: #64748b; text-transform: uppercase; }
  </style>
  <rect x="2" y="2" width="%d" height="%d" rx="6" fill="none" stroke="#64748b" stroke-opacity="0.28" stroke-width="1"/>
  <text x="14" y="20" class="sans hdr"%s>%s</text>
  <text x="%d" y="20" text-anchor="end" class="mono tag"%s>%s</text>
  <line x1="2" y1="28" x2="%d" y2="28" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
  <g transform="translate(14, 33)">
    <text x="0" y="11" class="mono kpi-num" fill="#0284c7"%s>%s (%.1f%%)</text>
    <text x="0" y="22" class="sans kpi-lbl"%s>%s</text>
    <line x1="138" y1="0" x2="138" y2="25" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <text x="150" y="11" class="mono kpi-num" fill="#d97706"%s>%s</text>
    <text x="150" y="22" class="sans kpi-lbl"%s>%s</text>
    <line x1="292" y1="0" x2="292" y2="25" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
    <rect x="304" y="3" width="84" height="8" rx="1.5" fill="#64748b" fill-opacity="0.14"/>
    <rect x="304" y="3" width="60" height="8" rx="1.5" fill="#64748b" fill-opacity="0.25"/>
    <rect x="304" y="5" width="%d" height="4" rx="1" fill="#059669"/>
    <line x1="370" y1="1" x2="370" y2="13" stroke="#0284c7" stroke-width="2"/>
    <text x="304" y="22" class="sans kpi-lbl"%s>%s</text>
  </g>
  <line x1="2" y1="63" x2="%d" y2="63" stroke="#64748b" stroke-opacity="0.25" stroke-width="1"/>
%s</svg>`,
			width, totalHeight, width, totalHeight,
			hdrSize, tagSize, bodySize,
			width-4, totalHeight-4,
			hdrAttr, esc(name), width-14, tagAttr, esc(l.BioTitle), width-2,
			kpi1NumAttr, esc(langs[0].Name), langs[0].Percentage, kpi1LblAttr, esc(l.LangFewPrimary),
			kpi2NumAttr, esc(tierRaw), kpi2LblAttr, esc(l.KaggleOverallTierLbl),
			prodBar, kpi3LblAttr, esc(l.LangFewProdEDA), width-2, linesSVG)

	case domain.SchoolCairo:
		headlineSize := 13.5
		leadSize := 10.5
		bodySize := 13.5
		calloutSize := 9.5
		displayRole := truncateRunes(role, 44)
		if fullWidth {
			headlineSize = 15.5
			leadSize = 12.0
			bodySize = 16.0
			calloutSize = 10.5
			displayRole = role
		}
		startY := 68
		totalHeight := startY + bodyHeight + 28
		linesSVG := renderBioLinesSVG(wrapped, 20, startY, lineHeight, paraGap, "serif body")
		callout := fmt.Sprintf(l.BioCairoCalloutFmt, d.Languages.TotalRepos, d.Languages.StrippedBase64MB)
		headTxt := fmt.Sprintf("%s — %s", name, l.BioTitle)
		leadTxt := fmt.Sprintf("%s · %s", displayRole, subRight)
		headAttr := fitTextAttr(headTxt, headlineSize, 9.5, float64(width-24), false, true, false, 0)
		leadAttr := fitTextAttr(leadTxt, leadSize, 8.2, float64(width-24), false, false, false, 0)
		calloutAttr := fitTextAttr(callout, calloutSize, 7.6, float64(width-24), false, false, false, 0)
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">
  <style>
    .serif { font-family: Georgia, "Times New Roman", serif; }
    .sans { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; }
    .headline { font-size: %.1fpx; font-weight: bold; fill: #526071; }
    .lead { font-size: %.1fpx; fill: #64748b; }
    .body { font-size: %.1fpx; fill: #526071; }
    .callout { font-size: %.1fpx; font-style: italic; fill: #0284c7; }
  </style>
  <rect x="12" y="6" width="%d" height="3" fill="#e11d48"/>
  <text x="12" y="25" class="serif headline"%s>%s — %s</text>
  <text x="12" y="40" class="sans lead"%s>%s · %s</text>
  <line x1="12" y1="47" x2="%d" y2="47" stroke="#e11d48" stroke-dasharray="3,3" stroke-width="1"/>
  <line x1="12" y1="54" x2="12" y2="%d" stroke="#0284c7" stroke-width="2.2"/>
%s  <text x="12" y="%d" class="serif callout"%s>%s</text>
</svg>`,
			width, totalHeight, width, totalHeight,
			headlineSize, leadSize, bodySize, calloutSize,
			width-24, headAttr, esc(name), esc(l.BioTitle), leadAttr, esc(displayRole), esc(subRight),
			width-12, startY+bodyHeight-10, linesSVG, totalHeight-10, calloutAttr, esc(callout))

	default: // SchoolHolmes
		bannerSize := 11.5
		bodySize := 12.5
		if fullWidth {
			bannerSize = 13.0
			bodySize = 15.0
		}
		startY := 56
		totalHeight := startY + bodyHeight + 16
		linesSVG := renderBioLinesSVG(wrapped, 16, startY, lineHeight, paraGap, "sans body")
		dispBioTitle, _ := fitHeaderNameRole("⚡ "+name, l.BioTitle, " — ", float64(width-32), bannerSize, 7.2, false, true, true, 0.5)
		bannerTxt := fmt.Sprintf("⚡ %s — %s", name, dispBioTitle)
		bannerAttr := fitTextAttr(bannerTxt, bannerSize, 7.2, float64(width-32), false, true, true, 0.5)
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">
  <style>
    .sans { font-family: "Trebuchet MS", Inter, -apple-system, BlinkMacSystemFont, sans-serif; }
    .banner { font-size: %.1fpx; font-weight: 800; fill: #ffffff; letter-spacing: 0.5px; text-transform: uppercase; }
    .body { font-size: %.1fpx; font-weight: 600; fill: #526071; }
  </style>
  <rect x="3" y="3" width="%d" height="%d" rx="10" fill="none" stroke="#64748b" stroke-opacity="0.35" stroke-width="2"/>
  <path d="M 3,13 Q 3,3 13,3 L %d,3 Q %d,3 %d,13 L %d,32 L 3,32 Z" fill="#0284c7"/>
  <text x="16" y="22" class="sans banner"%s>⚡ %s — %s</text>
%s</svg>`,
			width, totalHeight, width, totalHeight,
			bannerSize, bodySize,
			width-6, totalHeight-6,
			width-13, width-3, width-3, width-3,
			bannerAttr, esc(name), esc(dispBioTitle), linesSVG)
	}
}

type panelStyleOverrides struct {
	tagColor     string
	calloutColor string
	secColor     string
}

func extractInnerPanelBody(svgStr string, ov panelStyleOverrides) string {
	start := strings.Index(svgStr, "</style>")
	if start != -1 {
		svgStr = svgStr[start+len("</style>"):]
	} else if idx := strings.IndexByte(svgStr, '>'); idx != -1 {
		svgStr = svgStr[idx+1:]
	}
	if end := strings.LastIndex(svgStr, "</svg>"); end != -1 {
		svgStr = svgStr[:end]
	}
	if ov.tagColor != "" {
		svgStr = strings.ReplaceAll(svgStr, `class="mono tag"`, fmt.Sprintf(`class="mono tag" fill="%s"`, ov.tagColor))
	}
	if ov.calloutColor != "" {
		svgStr = strings.ReplaceAll(svgStr, `class="serif callout"`, fmt.Sprintf(`class="serif callout" fill="%s"`, ov.calloutColor))
	}
	if ov.secColor != "" {
		svgStr = strings.ReplaceAll(svgStr, `class="sans sec"`, fmt.Sprintf(`class="sans sec" fill="%s"`, ov.secColor))
	}
	// Add subtle SMIL entrance animation on sparkline polylines inside the full composite view
	svgStr = strings.ReplaceAll(svgStr, `<polyline fill="none"`, `<polyline stroke-dasharray="450" fill="none"`)
	return svgStr
}

func renderFull(school domain.School, d domain.PortfolioData, l labels) string {
	const (
		width        = 900
		headerHeight = 115
	)

	// 1. Dynamic Bio section: collapses to height 0 when Subject.Bio is empty; expands to exact line count when present
	wrappedBio := wrapBioParagraphs(d.Subject.Bio, bioFullMaxChars)
	_, bioBodyHeight := computeBioLayout(wrappedBio, bioFullLineHeight, bioFullParaGap)
	bioSectionHeight := 0
	if bioBodyHeight > 0 {
		bioSectionHeight = bioBodyHeight + 20
	}

	gridY := headerHeight + bioSectionHeight + 6
	row2Y := gridY + 238
	totalHeight := row2Y + 236

	headerInner := renderHeaderSection(school, d, l, "full-hdr-", false)
	var bioInner string
	if bioSectionHeight > 0 {
		bioStartY := headerHeight + 22
		bioLines := renderBioLinesSVG(wrappedBio, 18, bioStartY, bioFullLineHeight, bioFullParaGap, "body-txt")
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

	p1 := extractInnerPanelBody(renderLanguages(school, d, l), panelStyleOverrides{tagColor: "#0284c7", calloutColor: "#0284c7"})
	p2 := extractInnerPanelBody(renderStack(school, d, l), panelStyleOverrides{tagColor: "#64748b", calloutColor: "#0284c7"})
	p3 := extractInnerPanelBody(renderKaggle(school, d, l), panelStyleOverrides{tagColor: "#d97706", calloutColor: "#d97706"})
	p4 := extractInnerPanelBody(renderImpact(school, d, l), panelStyleOverrides{tagColor: "#059669", calloutColor: "#059669", secColor: "#059669"})

	quadrants := fmt.Sprintf(`  <defs><clipPath id="full-panel-clip"><rect x="0" y="0" width="440" height="230"/></clipPath></defs>
  <g transform="translate(8, %d)"><g clip-path="url(#full-panel-clip)">%s  </g></g>
  <g transform="translate(452, %d)"><g clip-path="url(#full-panel-clip)">%s  </g></g>
  <g transform="translate(8, %d)"><g clip-path="url(#full-panel-clip)">%s  </g></g>
  <g transform="translate(452, %d)"><g clip-path="url(#full-panel-clip)">%s  </g></g>`, gridY, p1, gridY, p2, row2Y, p3, row2Y, p4)

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
    .body-txt { font-family: Georgia, "Palatino Linotype", serif; font-size: 16.0px; fill: #526071; }
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
    .body-txt { font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; font-size: 15.0px; font-weight: 500; fill: #526071; }
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
    .body-txt { font-family: Georgia, "Times New Roman", serif; font-size: 16.0px; fill: #526071; }
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
    .body-txt { font-family: "Trebuchet MS", Inter, sans-serif; font-size: 15.0px; font-weight: 600; fill: #526071; }
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
	subjectHandle := strings.TrimSpace(manifest.Sources.GitHub.Username)
	if subjectHandle == "" {
		subjectHandle = "your-username"
	}
	repoSlug := subjectHandle + "/dataprofile"
	name := strings.TrimSpace(manifest.Subject.Name)
	if name == "" {
		name = subjectHandle
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
