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
		"pandas==2.1.0\nscikit-learn>=1.3\nplotly\nstreamlit",
		"polars\npytorch\nxgboost\nduckdb\nmatplotlib",
		"pandas\nscikit-learn\ntransformers\nseaborn",
	}

	stack := DetectStackFromManifests(manifests)
	if len(stack.DataAnalytics) == 0 || stack.DataAnalytics[0].Name != "Pandas / Polars" || stack.DataAnalytics[0].RepoCount != 3 {
		t.Errorf("unexpected DataAnalytics leader: %+v", stack.DataAnalytics)
	}
	if len(stack.MLAndAI) == 0 || stack.MLAndAI[0].Name != "Scikit-Learn" || stack.MLAndAI[0].RepoCount != 2 {
		t.Errorf("unexpected MLAndAI leader: %+v", stack.MLAndAI)
	}
}
