package domain

// School represents a classic Data Visualization School of Thought.
// See GLOSSARY.md for canonical definitions.
type School string

const (
	SchoolTufte  School = "tufte"  // Edward Tufte: Maximal data-ink ratio, Dot-Dash plots, Sparklines, Supertables
	SchoolFew    School = "few"    // Stephen Few: Executive telemetry, KPI strips, Bullet Graphs with target markers
	SchoolCairo  School = "cairo"  // Alberto Cairo: Functional Art journalism, Waffle 10x5, Lollipops & Callouts
	SchoolHolmes School = "holmes" // Nigel Holmes: Pictorial Isotypes, stepped podium & Trophy storytelling
)

// AllSchools returns the canonical ordered list of supported visualization schools.
func AllSchools() []School {
	return []School{SchoolTufte, SchoolFew, SchoolCairo, SchoolHolmes}
}

// Panel represents a thematic visual card generated for a Subject.
type Panel string

const (
	PanelSummary   Panel = "summary"   // Consolidated 4-quadrant executive summary
	PanelLanguages Panel = "languages" // Panel 1: Code, Notebooks & Production Engineering
	PanelStack     Panel = "stack"     // Panel 2: Data, ML & Analytics Ecosystem
	PanelKaggle    Panel = "kaggle"    // Panel 3: Kaggle Competitive Standing & Medals
	PanelImpact    Panel = "impact"    // Panel 4: Open-Source & AI Community Impact
	PanelHeader    Panel = "header"    // Animated full-width header banner (900px)
	PanelBio       Panel = "bio"       // Elastic career bio panel (440px mobile/grid width)
	PanelBioFull   Panel = "bio-full"  // Elastic career bio panel (900px full width)
	PanelFull      Panel = "full"      // Full portfolio composite: Animated Header + Dynamic Bio + 4 Complete Panels (900px)
)

// AllPanels returns the canonical ordered list of supported panels.
func AllPanels() []Panel {
	return []Panel{
		PanelSummary,
		PanelLanguages,
		PanelStack,
		PanelKaggle,
		PanelImpact,
		PanelHeader,
		PanelBio,
		PanelBioFull,
		PanelFull,
	}
}

// ProfileManifest represents the declarative profile.yaml configuration file.
type ProfileManifest struct {
	Subject       SubjectConfig       `json:"subject"`
	Sources       SourcesConfig       `json:"sources"`
	Normalization NormalizationConfig `json:"normalization"`
	Output        OutputConfig        `json:"output"`
}

// SubjectConfig identifies the practitioner whose portfolio is rendered.
type SubjectConfig struct {
	Name      string   `json:"name"`
	Role      string   `json:"role"`
	Locale    string   `json:"locale"`    // "en" (default) or "pt-BR"
	Headlines []string `json:"headlines"` // Animated typewriter phrases in header/full panels
	Bio       []string `json:"bio"`       // Career narrative paragraphs/lines with dynamic height adaptation
}

// SourcesConfig groups the external platforms queried for portfolio metrics.
type SourcesConfig struct {
	GitHub      GitHubSourceConfig      `json:"github"`
	Kaggle      KaggleSourceConfig      `json:"kaggle"`
	HuggingFace HuggingFaceSourceConfig `json:"huggingface"`
}

// GitHubSourceConfig configures GitHub GraphQL v4 and REST v3 collection.
type GitHubSourceConfig struct {
	Username            string   `json:"username"`
	IncludeForks        bool     `json:"include_forks"`
	IncludePrivateCount bool     `json:"include_private_count"`
	ExcludeRepos        []string `json:"exclude_repos"`
	ExcludeLangs        []string `json:"exclude_langs"`
}

// KaggleSourceConfig configures hybrid Kaggle collection (headless Chrome CDP + optional REST API v1).
type KaggleSourceConfig struct {
	Username string `json:"username"`
	Enabled  bool   `json:"enabled"`
	Mode     string `json:"mode"` // KaggleModeHybrid | KaggleModeBrowser | KaggleModeAPI
}

// Kaggle collection modes.
const (
	KaggleModeHybrid  = "hybrid"
	KaggleModeBrowser = "browser"
	KaggleModeAPI     = "api"
)

// Kaggle progression tiers.
const (
	KaggleTierNovice      = "Novice"
	KaggleTierContributor = "Contributor"
	KaggleTierExpert      = "Expert"
	KaggleTierMaster      = "Master"
	KaggleTierGrandmaster = "Grandmaster"
)

// HuggingFaceSourceConfig configures public Hugging Face Hub REST collection.
type HuggingFaceSourceConfig struct {
	Username string `json:"username"`
	Enabled  bool   `json:"enabled"`
}

// NormalizationConfig controls notebook deflation and dependency ecosystem detection.
type NormalizationConfig struct {
	DeflateNotebooks bool `json:"deflate_notebooks"`
	DetectDataStack  bool `json:"detect_data_stack"`
}

// OutputConfig controls which Schools and Panels are generated and where they are written.
type OutputConfig struct {
	Directory string   `json:"directory"`
	Schools   []School `json:"schools"`
	Panels    []Panel  `json:"panels"`
}

// Palette defines the Universal Dual-Surface high-contrast semantic colors used across SVGs
// to guarantee WCAG AA legibility on both white (#ffffff) and dark (#0d1117) GitHub themes
// without relying on fragile OS-level @media queries.
type Palette struct {
	PrimaryInk   string // #526071 (Slate 550: 5.8:1 on white, 3.4:1 on dark)
	SecondaryInk string // #64748b (Slate 500: 4.76:1 on white, 4.15:1 on dark)
	AccentSky    string // #0284c7 (Sky 600)
	AccentBlue   string // #2563eb (Blue 600)
	AccentAmber  string // #d97706 (Amber 600 - Gold/Master)
	AccentGreen  string // #059669 (Emerald 600 - Production/Growth)
	AccentViolet string // #7c3aed (Violet 600 - ML/AI)
	AccentRose   string // #e11d48 (Rose 600 - Editorial/Cairo)
}

// DefaultPalette returns the canonical Universal Dual-Surface Palette.
func DefaultPalette() Palette {
	return Palette{
		PrimaryInk:   "#526071",
		SecondaryInk: "#64748b",
		AccentSky:    "#0284c7",
		AccentBlue:   "#2563eb",
		AccentAmber:  "#d97706",
		AccentGreen:  "#059669",
		AccentViolet: "#7c3aed",
		AccentRose:   "#e11d48",
	}
}

// PortfolioData aggregates normalized metrics across all enabled Sources for rendering.
type PortfolioData struct {
	Subject   SubjectConfig     `json:"subject"`
	Languages LanguagesPanel    `json:"languages"`
	Stack     StackPanel        `json:"stack"`
	Kaggle    KaggleStanding    `json:"kaggle"`
	Impact    ImpactMetrics     `json:"impact"`
}

// LanguagesPanel holds normalized code distribution and engineering maturity metrics (Panel 1).
type LanguagesPanel struct {
	TotalRepos          int            `json:"total_repos"`
	Languages           []LanguageStat `json:"languages"`
	ProductionPct       float64        `json:"production_pct"`
	ProductionFiles     int            `json:"production_files"`
	ExplorationPct      float64        `json:"exploration_pct"`
	ExplorationFiles    int            `json:"exploration_files"`
	StrippedBase64MB    float64        `json:"stripped_base64_mb"`
	LockfileRepos       int            `json:"lockfile_repos"`
	CICDRepos           int            `json:"cicd_repos"`
	TestSuiteRepos      int            `json:"test_suite_repos"`
	ContainerizedRepos  int            `json:"containerized_repos"`
}

// LanguageStat represents a single programming language's normalized proportion.
type LanguageStat struct {
	Name       string  `json:"name"`
	Percentage float64 `json:"percentage"`
	Bytes      int64   `json:"bytes"`
	Color      string  `json:"color"`
}

// StackPanel holds detected data science, ML, and visualization libraries (Panel 2).
type StackPanel struct {
	DataAnalytics []StackItem `json:"data_analytics"`
	MLAndAI       []StackItem `json:"ml_and_ai"`
	VizAndApps    []StackItem `json:"viz_and_apps"`
	TotalSources  int         `json:"total_sources"`
}

// StackItem represents a detected library or tool across repositories.
type StackItem struct {
	Name      string  `json:"name"`
	RepoCount int     `json:"repo_count"`
	WeightPct float64 `json:"weight_pct"` // Relative prominence 0..100
	IsLeader  bool    `json:"is_leader"`
}

// KaggleStanding holds competitive standing, medals, and category progressions (Panel 3).
type KaggleStanding struct {
	Available   bool             `json:"available"`
	Username    string           `json:"username"`
	OverallTier string           `json:"overall_tier"` // Novice, Contributor, Expert, Master, Grandmaster
	PeakRank    int              `json:"peak_rank"`
	TotalGold   int              `json:"total_gold"`
	TotalSilver int              `json:"total_silver"`
	TotalBronze int              `json:"total_bronze"`
	Categories  []KaggleCategory `json:"categories"`
}

// TotalMedals returns the sum of Gold, Silver, and Bronze medals.
func (k KaggleStanding) TotalMedals() int {
	return k.TotalGold + k.TotalSilver + k.TotalBronze
}

// KaggleCategory represents one of the 4 official Kaggle progression tracks.
type KaggleCategory struct {
	Name        string  `json:"name"` // Notebooks, Datasets, Competitions, Discussions
	Tier        string  `json:"tier"`
	Rank        int     `json:"rank"`
	TotalUsers  int     `json:"total_users"`
	TopPct      string  `json:"top_pct"`
	Gold        int     `json:"gold"`
	Silver      int     `json:"silver"`
	Bronze      int     `json:"bronze"`
	ProgressPct float64 `json:"progress_pct"` // 0..100 progress along bullet/lollipop scale
	TargetPct   float64 `json:"target_pct"`   // Next tier threshold marker 0..100
}

// ImpactMetrics holds combined GitHub and Hugging Face community reach metrics (Panel 4).
type ImpactMetrics struct {
	AnnualContributions int     `json:"annual_contributions"`
	YoYGrowthPct        float64 `json:"yoy_growth_pct"`
	ActiveStreakDays    int     `json:"active_streak_days"`
	TotalStars          int     `json:"total_stars"`
	PublicRepos         int     `json:"public_repos"`
	ContributedRepos    int     `json:"contributed_repos"`
	TotalForks          int     `json:"total_forks"`
	RepoViews14d        int     `json:"repo_views_14d"`
	LinesAdded          int     `json:"lines_added"`
	LinesDeleted        int     `json:"lines_deleted"`
	HFModels            int     `json:"hf_models"`
	HFDatasets          int     `json:"hf_datasets"`
	HFSpaces            int     `json:"hf_spaces"`
	HFDownloads30d      int     `json:"hf_downloads_30d"`
	HFLikes             int     `json:"hf_likes"`
	WeeklyTrend         []int   `json:"weekly_trend"` // 12 or 52 points for sparkline/curve
}
