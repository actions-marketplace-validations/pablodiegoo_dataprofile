package normalizer

import (
	"strings"
	"testing"
)

func TestDeflateNotebookJSON(t *testing.T) {
	// Synthetic .ipynb with small Python code cell and huge base64 PNG output cell
	fakeBase64 := strings.Repeat("iVBORw0KGgoAAAANSUhEUgAAAAUA", 400)
	ipynb := `{
		"metadata": {
			"kernelspec": {"language": "python"}
		},
		"cells": [
			{
				"cell_type": "code",
				"source": ["import pandas as pd\n", "df = pd.read_parquet('data.parquet')\n"],
				"outputs": [{"data": {"image/png": "` + fakeBase64 + `"}}]
			},
			{
				"cell_type": "markdown",
				"source": ["# Exploratory Analysis\n"]
			}
		]
	}`

	stats, err := DeflateNotebookJSON([]byte(ipynb))
	if err != nil {
		t.Fatalf("DeflateNotebookJSON returned error: %v", err)
	}
	if stats.Language != "Python" {
		t.Errorf("expected language Python, got %q", stats.Language)
	}
	if stats.CodeCellCount != 1 {
		t.Errorf("expected 1 code cell, got %d", stats.CodeCellCount)
	}
	if stats.CodeBytes <= 0 || stats.CodeBytes > 120 {
		t.Errorf("expected small code byte count (~57 bytes), got %d", stats.CodeBytes)
	}
	if stats.StrippedBytes < 10000 {
		t.Errorf("expected >10KB stripped base64 bloat, got %d", stats.StrippedBytes)
	}

	// R notebook identified via kernelspec.name = "ir" and single-string cell source
	rNotebook := `{
		"metadata": {
			"kernelspec": {"name": "ir"}
		},
		"cells": [
			{
				"cell_type": "code",
				"source": "library(tidyverse)\ndf <- read_csv('data.csv')\n",
				"outputs": []
			}
		]
	}`
	rStats, err := DeflateNotebookJSON([]byte(rNotebook))
	if err != nil {
		t.Fatalf("DeflateNotebookJSON (R) returned error: %v", err)
	}
	if rStats.Language != "R" {
		t.Errorf("expected language R for ir kernelspec, got %q", rStats.Language)
	}
	if rStats.CodeBytes != int64(len("library(tidyverse)\ndf <- read_csv('data.csv')\n")) {
		t.Errorf("unexpected R code bytes: %d", rStats.CodeBytes)
	}
}

func TestNormalizeLanguageMapDeflatesNotebooks(t *testing.T) {
	raw := map[string]int64{
		"Jupyter Notebook": 10 * 1024 * 1024, // 10 MB raw .ipynb
		"Python":           2 * 1024 * 1024,  // 2 MB .py modules
		"SQL":              1 * 1024 * 1024,  // 1 MB .sql
		"HTML":             500 * 1024,       // excluded
	}

	langs, prodPct, edaPct, strippedMB := NormalizeLanguageMap(raw, true, []string{"HTML"})
	if len(langs) == 0 {
		t.Fatal("expected non-empty normalized languages")
	}
	if langs[0].Name != "Python" {
		t.Errorf("expected Python to lead after notebook deflation, got %s", langs[0].Name)
	}
	for _, l := range langs {
		if l.Name == "Jupyter Notebook" || l.Name == "HTML" {
			t.Errorf("unexpected language %q in normalized output", l.Name)
		}
	}
	if strippedMB < 8.0 {
		t.Errorf("expected >8.0 MB stripped from 10 MB Jupyter Notebook, got %.1f MB", strippedMB)
	}
	if prodPct <= 50.0 || edaPct <= 0.0 {
		t.Errorf("unexpected prod/eda split: prod=%.1f%% eda=%.1f%%", prodPct, edaPct)
	}
}

func TestDetectStackFromManifests(t *testing.T) {
	manifests := []string{
		// Repo 1: pyproject.toml + SQL artifacts
		`[project]
dependencies = ["pandas>=2.1.0", "scikit-learn>=1.3", "plotly>=5.18", "streamlit>=1.30"]
-- analytics_queries.sql
SELECT * FROM events;`,
		// Repo 2: requirements.txt
		"polars==0.20.0\npytorch>=2.1\nxgboost>=2.0\nduckdb>=0.10\nmatplotlib>=3.8\ndbt\n",
		// Repo 3: environment.yml
		`name: ml-env
dependencies:
  - pandas
  - scikit-learn
  - transformers
  - seaborn
  - pyspark
  - plotly`,
		// Repo 4: R DESCRIPTION file + SQL
		`Package: customerchurn
Title: Customer Churn Modeling
Imports:
    ggplot2,
    tidyverse,
    dplyr
Suggests: duckdb`,
	}

	stack := DetectStackFromManifests(manifests)
	if stack.TotalSources != 4 {
		t.Errorf("expected TotalSources=4, got %d", stack.TotalSources)
	}

	// Verify Data & Analytics pillar (all 4 canonical items ranked by RepoCount)
	// Pandas/Polars: repos 1, 2, 3 -> 3 (100%, IsLeader=true)
	// SQL/DuckDB: repos 1, 2, 4 -> 3 (100%, IsLeader=false due to canonical tie-break)
	// Apache Spark: repo 3 -> 1 (33%, IsLeader=false)
	// dbt/Airflow: repo 2 -> 1 (33%, IsLeader=false)
	if len(stack.DataAnalytics) != 4 {
		t.Fatalf("expected 4 items in DataAnalytics, got %d", len(stack.DataAnalytics))
	}
	if stack.DataAnalytics[0].Name != "Pandas / Polars" || stack.DataAnalytics[0].RepoCount != 3 || stack.DataAnalytics[0].WeightPct != 100 || !stack.DataAnalytics[0].IsLeader {
		t.Errorf("unexpected DataAnalytics[0]: %+v", stack.DataAnalytics[0])
	}
	if stack.DataAnalytics[1].Name != "SQL / DuckDB" || stack.DataAnalytics[1].RepoCount != 3 || stack.DataAnalytics[1].WeightPct != 100 || stack.DataAnalytics[1].IsLeader {
		t.Errorf("unexpected DataAnalytics[1]: %+v", stack.DataAnalytics[1])
	}
	if stack.DataAnalytics[2].Name != "Apache Spark" || stack.DataAnalytics[2].RepoCount != 1 || stack.DataAnalytics[2].WeightPct != 33 || stack.DataAnalytics[2].IsLeader {
		t.Errorf("unexpected DataAnalytics[2]: %+v", stack.DataAnalytics[2])
	}
	if stack.DataAnalytics[3].Name != "dbt / Airflow" || stack.DataAnalytics[3].RepoCount != 1 || stack.DataAnalytics[3].WeightPct != 33 || stack.DataAnalytics[3].IsLeader {
		t.Errorf("unexpected DataAnalytics[3]: %+v", stack.DataAnalytics[3])
	}

	// Verify Machine Learning & AI pillar
	// Scikit-Learn: repos 1, 3 -> 2 (100%, IsLeader=true)
	// PyTorch: repo 2 -> 1 (50%, IsLeader=false)
	// XGBoost / LightGBM: repo 2 -> 1 (50%, IsLeader=false)
	// Hugging Face: repo 3 -> 1 (50%, IsLeader=false)
	if len(stack.MLAndAI) != 4 {
		t.Fatalf("expected 4 items in MLAndAI, got %d", len(stack.MLAndAI))
	}
	if stack.MLAndAI[0].Name != "Scikit-Learn" || stack.MLAndAI[0].RepoCount != 2 || stack.MLAndAI[0].WeightPct != 100 || !stack.MLAndAI[0].IsLeader {
		t.Errorf("unexpected MLAndAI[0]: %+v", stack.MLAndAI[0])
	}
	for i := 1; i < 4; i++ {
		if stack.MLAndAI[i].RepoCount != 1 || stack.MLAndAI[i].WeightPct != 50 || stack.MLAndAI[i].IsLeader {
			t.Errorf("unexpected MLAndAI[%d]: %+v", i, stack.MLAndAI[i])
		}
	}

	// Verify Visualization & Data Apps pillar
	// Plotly / Dash: repos 1, 3 -> 2 (100%, IsLeader=true)
	// Matplotlib / Seaborn: repos 2, 3 -> 2 (100%, IsLeader=false)
	// Streamlit: repo 1 -> 1 (50%, IsLeader=false)
	// ggplot2 / Tidyverse: repo 4 -> 1 (50%, IsLeader=false)
	if len(stack.VizAndApps) != 4 {
		t.Fatalf("expected 4 items in VizAndApps, got %d", len(stack.VizAndApps))
	}
	if stack.VizAndApps[0].Name != "Plotly / Dash" || stack.VizAndApps[0].RepoCount != 2 || stack.VizAndApps[0].WeightPct != 100 || !stack.VizAndApps[0].IsLeader {
		t.Errorf("unexpected VizAndApps[0]: %+v", stack.VizAndApps[0])
	}
	if stack.VizAndApps[1].Name != "Matplotlib / Seaborn" || stack.VizAndApps[1].RepoCount != 2 || stack.VizAndApps[1].WeightPct != 100 || stack.VizAndApps[1].IsLeader {
		t.Errorf("unexpected VizAndApps[1]: %+v", stack.VizAndApps[1])
	}

	// Verify empty manifests do not flag any leader and have 0% weight
	emptyStack := DetectStackFromManifests(nil)
	for _, item := range emptyStack.DataAnalytics {
		if item.IsLeader || item.WeightPct != 0 || item.RepoCount != 0 {
			t.Errorf("expected zeroed item for empty manifest, got %+v", item)
		}
	}
}

func TestNormalizeLanguageMapTop5WithoutFakeShellAndSQLCanonicalization(t *testing.T) {
	raw := map[string]int64{
		"Python":     450000,
		"Cython":     50000, // should merge into Python (500,000 total)
		"Go":         300000,
		"JavaScript": 200000,
		"TypeScript": 150000,
		"PLpgSQL":    50000, // should merge into SQL (100,000 total)
		"TSQL":       30000,
		"PLSQL":      20000,
		"R":          80000,
		"PHP":        40000,
	}

	langs, _, _, _ := NormalizeLanguageMap(raw, true, nil)
	if len(langs) != 5 {
		t.Fatalf("expected Top 5 languages, got %d", len(langs))
	}
	wantOrder := []string{"Python", "Go", "JavaScript", "TypeScript", "SQL"}
	wantBytes := []int64{500000, 300000, 200000, 150000, 100000}
	seenColors := make(map[string]string)
	for i, want := range wantOrder {
		if langs[i].Name != want {
			t.Errorf("expected langs[%d].Name=%q, got %q", i, want, langs[i].Name)
		}
		if langs[i].Bytes != wantBytes[i] {
			t.Errorf("expected langs[%d] (%s) Bytes=%d, got %d", i, want, wantBytes[i], langs[i].Bytes)
		}
		if langs[i].Color == "" {
			t.Errorf("expected non-empty color for %s", langs[i].Name)
		}
		if prevLang, exists := seenColors[langs[i].Color]; exists {
			t.Errorf("expected distinct hues for all Top 5 languages, but %s and %s both got %s",
				prevLang, langs[i].Name, langs[i].Color)
		}
		seenColors[langs[i].Color] = langs[i].Name
	}
}

func TestDetectStackDynamicCatalogReplacesZeroSparkWithRealEcosystems(t *testing.T) {
	workspaces := []string{
		"import pandas as pd\nimport numpy as np\nimport geopandas as gpd\nimport statsmodels.api as sm\nimport folium\n",
		"import pandas as pd\nimport scipy.stats as stats\nfrom sklearn.cluster import KMeans, DBSCAN\nimport matplotlib.pyplot as plt\nimport seaborn as sns\n",
		"import duckdb\nimport typer\nimport openai\nimport whisper\nimport pypandoc\n",
	}

	stack := DetectStackFromManifests(workspaces)
	for _, item := range stack.DataAnalytics {
		if item.RepoCount == 0 {
			t.Errorf("expected dynamic Top 4 DataAnalytics to have non-zero counts, got %+v", stack.DataAnalytics)
		}
	}
	for _, item := range stack.MLAndAI {
		if item.RepoCount == 0 {
			t.Errorf("expected dynamic Top 4 MLAndAI to have non-zero counts, got %+v", stack.MLAndAI)
		}
	}
}
