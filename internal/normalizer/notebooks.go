package normalizer

import (
	"encoding/json"
	"math"
	"sort"
	"strings"

	"github.com/pablodiegoo/dataprofile/internal/domain"
)

// NotebookStats summarizes the result of deflating a raw .ipynb JSON file.
type NotebookStats struct {
	Language        string
	RawBytes        int64
	CodeBytes       int64
	MarkdownBytes   int64
	StrippedBytes   int64 // Output cells (base64 PNGs, HTML tables, JSON overhead)
	CodeCellCount   int
}

type ipynbRoot struct {
	Cells    []ipynbCell   `json:"cells"`
	Metadata ipynbMetadata `json:"metadata"`
}

type ipynbCell struct {
	CellType string          `json:"cell_type"`
	Source   json.RawMessage `json:"source"`
	Outputs  json.RawMessage `json:"outputs"`
}

type ipynbMetadata struct {
	KernelSpec struct {
		Name     string `json:"name"`
		Language string `json:"language"`
	} `json:"kernelspec"`
	LanguageInfo struct {
		Name string `json:"name"`
	} `json:"language_info"`
}

// DeflateNotebookJSON parses a raw .ipynb JSON document, extracts only executable code cell bytes,
// and computes the exact amount of base64/output bloat stripped.
func DeflateNotebookJSON(raw []byte) (NotebookStats, error) {
	stats := NotebookStats{
		Language: "Python",
		RawBytes: int64(len(raw)),
	}

	var root ipynbRoot
	if err := json.Unmarshal(raw, &root); err != nil {
		return stats, err
	}

	lang := strings.TrimSpace(root.Metadata.LanguageInfo.Name)
	if lang == "" {
		lang = strings.TrimSpace(root.Metadata.KernelSpec.Language)
	}
	if lang == "" {
		lang = strings.TrimSpace(root.Metadata.KernelSpec.Name)
	}
	if lang != "" {
		stats.Language = canonicalLanguageName(lang)
	}

	for _, cell := range root.Cells {
		srcText := extractCellSource(cell.Source)
		switch cell.CellType {
		case "code":
			stats.CodeCellCount++
			stats.CodeBytes += int64(len(srcText))
		case "markdown":
			stats.MarkdownBytes += int64(len(srcText))
		}
	}

	if stats.RawBytes > stats.CodeBytes {
		stats.StrippedBytes = stats.RawBytes - stats.CodeBytes
	}
	return stats, nil
}

// ExtractNotebookCode parses a raw .ipynb JSON document and returns the concatenated source
// of all executable code cells, stripping base64 output payloads and JSON schema noise.
func ExtractNotebookCode(raw []byte) string {
	var root ipynbRoot
	if err := json.Unmarshal(raw, &root); err != nil || len(root.Cells) == 0 {
		return ""
	}
	var parts []string
	for _, cell := range root.Cells {
		if cell.CellType == "code" {
			if src := strings.TrimSpace(extractCellSource(cell.Source)); src != "" {
				parts = append(parts, src)
			}
		}
	}
	return strings.Join(parts, "\n")
}

func extractCellSource(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var lines []string
	if err := json.Unmarshal(raw, &lines); err == nil {
		return strings.Join(lines, "")
	}
	var single string
	if err := json.Unmarshal(raw, &single); err == nil {
		return single
	}
	return ""
}

func canonicalLanguageName(lang string) string {
	switch strings.ToLower(strings.TrimSpace(lang)) {
	case "python", "python3", "py", "cython":
		return "Python"
	case "r", "ir":
		return "R"
	case "julia":
		return "Julia"
	case "sql", "plpgsql", "tsql", "plsql":
		return "SQL"
	case "go", "golang":
		return "Go"
	default:
		if len(lang) > 0 {
			return strings.ToUpper(lang[:1]) + lang[1:]
		}
		return "Python"
	}
}

var distinctTop5Palette = []string{
	"#0284c7", // Sky 600
	"#d97706", // Amber 600
	"#00ADD8", // Cyan / Go
	"#2563eb", // Blue 600
	"#059669", // Emerald 600
	"#7c3aed", // Violet 600
	"#e11d48", // Rose 600
	"#64748b", // Slate 500
}

// LanguageColor returns a high-contrast semantic color for a programming language.
func LanguageColor(name string) string {
	switch strings.ToLower(canonicalLanguageName(name)) {
	case "python":
		return "#0284c7"
	case "sql":
		return "#d97706"
	case "go":
		return "#00ADD8"
	case "r":
		return "#2563eb"
	case "shell", "bash", "sh", "templ":
		return "#059669"
	case "julia", "php":
		return "#7c3aed"
	case "rust", "svelte":
		return "#e11d48"
	case "typescript":
		return "#2563eb"
	case "javascript":
		return "#7c3aed"
	default:
		return "#64748b"
	}
}

// NormalizeLanguageMap takes raw language byte counts (as reported by GitHub Linguist),
// deflates Jupyter Notebook / RMarkdown bloat when enabled, filters excluded languages,
// and returns the top languages with normalized percentages and EDA vs. Production split.
func NormalizeLanguageMap(rawBytes map[string]int64, deflate bool, exclude []string) (langs []domain.LanguageStat, prodPct, edaPct, strippedMB float64) {
	exMap := make(map[string]bool, len(exclude))
	for _, e := range exclude {
		exMap[strings.ToLower(strings.TrimSpace(e))] = true
	}

	cleaned := make(map[string]int64)
	var notebookRaw int64
	var notebookDeflated int64
	var prodBytes int64

	for lang, b := range rawBytes {
		if b <= 0 || exMap[strings.ToLower(lang)] {
			continue
		}
		switch lang {
		case "Jupyter Notebook":
			notebookRaw += b
			if deflate {
				// Empirical deflation: on average 12% of raw .ipynb bytes are actual executable Python source code;
				// 88% is JSON schema overhead, execution logs, and embedded base64 PNG/HTML plots.
				codeOnly := int64(math.Max(float64(b)*0.12, 1))
				notebookDeflated += codeOnly
				cleaned["Python"] += codeOnly
			} else {
				notebookDeflated += b
				cleaned["Jupyter Notebook"] += b
			}
		case "RMarkdown", "Quarto":
			notebookRaw += b
			if deflate {
				codeOnly := int64(math.Max(float64(b)*0.35, 1))
				notebookDeflated += codeOnly
				cleaned["R"] += codeOnly
			} else {
				notebookDeflated += b
				cleaned[lang] += b
			}
		case "PLpgSQL", "TSQL", "PLSQL":
			if !exMap["sql"] {
				cleaned["SQL"] += b
				prodBytes += b
			}
		case "Cython":
			if !exMap["python"] {
				cleaned["Python"] += b
				prodBytes += b
			}
		default:
			cleaned[lang] += b
			prodBytes += b
		}
	}

	if notebookRaw > notebookDeflated {
		strippedMB = float64(notebookRaw-notebookDeflated) / (1024.0 * 1024.0)
	}

	var total int64
	for _, b := range cleaned {
		total += b
	}
	if total == 0 {
		return nil, 100, 0, 0
	}

	type kv struct {
		name  string
		bytes int64
	}
	var items []kv
	for k, v := range cleaned {
		items = append(items, kv{k, v})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].bytes == items[j].bytes {
			return items[i].name < items[j].name
		}
		return items[i].bytes > items[j].bytes
	})

	if len(items) > 5 {
		items = items[:5]
	}

	var topTotal int64
	for _, it := range items {
		topTotal += it.bytes
	}
	if topTotal == 0 {
		topTotal = total
	}

	usedColors := make(map[string]bool, len(items))
	for _, it := range items {
		pct := math.Round((float64(it.bytes)/float64(topTotal))*1000) / 10
		color := LanguageColor(it.name)
		if usedColors[color] {
			for _, candidate := range distinctTop5Palette {
				if !usedColors[candidate] {
					color = candidate
					break
				}
			}
		}
		usedColors[color] = true
		langs = append(langs, domain.LanguageStat{
			Name:       it.name,
			Bytes:      it.bytes,
			Percentage: pct,
			Color:      color,
		})
	}

	totalCode := prodBytes + notebookDeflated
	if totalCode > 0 {
		prodPct = math.Round((float64(prodBytes)/float64(totalCode))*1000) / 10
		edaPct = math.Round((100.0-prodPct)*10) / 10
	} else {
		prodPct = 100.0
		edaPct = 0.0
	}

	return langs, prodPct, edaPct, math.Round(strippedMB*10) / 10
}

type stackRule struct {
	canonical string
	pillar    string // "data", "ml", "viz"
	tokens    []string
}

var canonicalStackRules = []stackRule{
	// Data & Analytics (first 4 are default baseline order)
	{canonical: "Pandas / Polars", pillar: "data", tokens: []string{"pandas", "polars", "geopandas"}},
	{canonical: "SQL / DuckDB", pillar: "data", tokens: []string{"duckdb", "sqlalchemy", "psycopg", "sql", "plpgsql", "bigquery", "snowflake", "sqlite"}},
	{canonical: "Apache Spark", pillar: "data", tokens: []string{"pyspark", "spark", "databricks", "delta-spark"}},
	{canonical: "dbt / Airflow", pillar: "data", tokens: []string{"dbt-", "dbt_core", "dbt", "apache-airflow", "airflow", "prefect", "dagster"}},
	{canonical: "NumPy / SciPy", pillar: "data", tokens: []string{"numpy", "scipy", "pyarrow", "parquet", "openpyxl"}},
	{canonical: "GeoPandas / Spatial", pillar: "data", tokens: []string{"geopandas", "shapely", "osmnx", "import pyproj", "from pyproj", "pyproj==", "pyproj>=", "fiona", "rasterio", "postgis"}},
	{canonical: "ETL & CLI Pipelines", pillar: "data", tokens: []string{"typer", "fastapi", "pydantic", "httpx", "celery", "alembic"}},

	// Machine Learning & AI (first 4 are default baseline order)
	{canonical: "Scikit-Learn", pillar: "ml", tokens: []string{"scikit-learn", "sklearn", "logisticregression", "linearregression", "randomforest"}},
	{canonical: "PyTorch", pillar: "ml", tokens: []string{"torch", "pytorch", "torchvision", "torchaudio", "lightning"}},
	{canonical: "XGBoost / LightGBM", pillar: "ml", tokens: []string{"xgboost", "lightgbm", "catboost", "optuna"}},
	{canonical: "Hugging Face", pillar: "ml", tokens: []string{"transformers", "datasets", "huggingface", "accelerate", "peft", "diffusers"}},
	{canonical: "Statsmodels / GLM", pillar: "ml", tokens: []string{"statsmodels", "pingouin", "linearmodels", "factor_analysis", "factor_analyzer", "chi2", "sm.ols", "smf.ols", "ols(", "glm"}},
	{canonical: "Clustering & PCA", pillar: "ml", tokens: []string{"kmeans", "dbscan", "gower", "umap", "hdbscan", "pca", "clustering"}},
	{canonical: "LLMs & Speech AI", pillar: "ml", tokens: []string{"whisper", "openai", "google-generativeai", "gemini", "langchain", "llamaindex", "vibevoice"}},

	// Visualization & Data Apps (first 4 are default baseline order)
	{canonical: "Plotly / Dash", pillar: "viz", tokens: []string{"plotly", "dash", "kaleido"}},
	{canonical: "Streamlit", pillar: "viz", tokens: []string{"streamlit", "gradio", "shiny"}},
	{canonical: "Matplotlib / Seaborn", pillar: "viz", tokens: []string{"matplotlib", "seaborn", "altair", "bokeh"}},
	{canonical: "ggplot2 / Tidyverse", pillar: "viz", tokens: []string{"ggplot2", "tidyverse", "dplyr", "tidymodels"}},
	{canonical: "Geospatial Maps", pillar: "viz", tokens: []string{"folium", "contextily", "leafmap", "osmnx", "voronoi", "choropleth", "shapefile"}},
	{canonical: "Report Automation", pillar: "viz", tokens: []string{"pypandoc", "python-pptx", "reportlab", "weasyprint", "jinja2", "xlsxwriter"}},
}

// DetectStackFromManifests scans repository and Analytical Workspace dependency manifests and code hints
// and aggregates counts across the 3 pillars, dynamically selecting the Top 4 ecosystems per pillar.
func DetectStackFromManifests(repoManifests []string) domain.StackPanel {
	counts := make(map[string]int)
	for _, content := range repoManifests {
		lower := strings.ToLower(content)
		for _, rule := range canonicalStackRules {
			for _, tok := range rule.tokens {
				if strings.Contains(lower, tok) {
					counts[rule.canonical]++
					break
				}
			}
		}
	}

	buildPillar := func(pillar string) []domain.StackItem {
		var list []domain.StackItem
		for _, rule := range canonicalStackRules {
			if rule.pillar != pillar {
				continue
			}
			c := counts[rule.canonical]
			list = append(list, domain.StackItem{
				Name:      rule.canonical,
				RepoCount: c,
			})
		}
		sort.SliceStable(list, func(i, j int) bool {
			return list[i].RepoCount > list[j].RepoCount
		})
		if len(list) > 4 {
			list = list[:4]
		}
		maxCount := 0
		for _, it := range list {
			if it.RepoCount > maxCount {
				maxCount = it.RepoCount
			}
		}
		for i := range list {
			if maxCount > 0 {
				list[i].WeightPct = math.Round((float64(list[i].RepoCount) / float64(maxCount)) * 100)
				if i == 0 && list[i].RepoCount > 0 {
					list[i].IsLeader = true
				}
			}
		}
		return list
	}

	return domain.StackPanel{
		DataAnalytics: buildPillar("data"),
		MLAndAI:       buildPillar("ml"),
		VizAndApps:    buildPillar("viz"),
		TotalSources:  len(repoManifests),
	}
}
