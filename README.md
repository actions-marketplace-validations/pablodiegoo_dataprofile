<div align="center">

# 📊 DataProfile

**Publication-Grade, Notebook-Normalized Portfolio Telemetry for Data Scientists, ML Engineers, Analysts & Researchers**

[![Go 1.23+](https://img.shields.io/badge/Go-1.23%2B-00ADD8?logo=go&logoColor=white)](go.mod)
[![Zero External Dependencies](https://img.shields.io/badge/Dependencies-Zero%20Stdlib-059669)](#architecture--engineering)
[![Visualization Schools](https://img.shields.io/badge/Schools-Tufte%20%7C%20Few%20%7C%20Cairo%20%7C%20Holmes-0284c7)](#-the-4-visualization-schools-of-thought)
[![Sources](https://img.shields.io/badge/Sources-GitHub%20%7C%20Kaggle%20%7C%20Hugging%20Face-d97706)](#tri-source-telemetry)

<p align="center">
  <img src="generated/few/summary.svg" width="48%" alt="Executive Summary (Stephen Few Grammar)" />
  <img src="generated/tufte/summary.svg" width="48%" alt="Executive Summary (Edward Tufte Grammar)" />
</p>

</div>

---

## Why `dataprofile` Exists

Traditional GitHub profile stat generators were built for general web developers. When applied to a **Data Scientist, Machine Learning Engineer, Data Analyst, or Quantitative Researcher**, they fail in three critical ways:

1. **The Jupyter `.ipynb` Base64 Distortion:** GitHub Linguist counts the raw byte size of `.ipynb` files—including embedded base64 PNG plots, HTML DataFrame dumps, and JSON cell metadata. A single notebook with three `matplotlib` charts easily outweighs 15,000 lines of production `.py`, `.sql`, or `.go` code, falsely branding the practitioner's profile as `"94% Jupyter Notebook"`.
2. **Invisible Data & ML Ecosystem:** Knowing someone writes `"Python"` or `"R"` tells a recruiter or fellow researcher almost nothing. What matters is the verified ecosystem inside `pyproject.toml`, `requirements.txt`, `environment.yml`, and `DESCRIPTION` (**Polars, DuckDB, PyTorch, Scikit-Learn, dbt, Spark, Hugging Face, Streamlit**), alongside production engineering hygiene (**Lockfiles, CI/CD, Automated Tests, Docker**).
3. **Fragmented Career Proof & Generic Aesthetics:** A data practitioner's real impact spans **GitHub** (code & open-source), **Kaggle** (competitions, datasets, notebooks, global ranks & medals), and **Hugging Face Hub** (models, datasets, spaces & downloads)—and someone who works with data deserves visuals rooted in the **classic schools of information design**, not neon gamer cards.

---

## 🧩 The 4 Analytical Panels (`Panel`)

`dataprofile` synthesizes your portfolio into **4 modular 2×2 panels** (plus an optional `summary.svg`), all rendered with a **Universal Dual-Surface Palette** (`fill="none"`) engineered for crisp WCAG AA contrast on **both White (`#ffffff`) and Dark (`#0d1117`) GitHub themes**:

| Panel File | Domain Focus | Key Engineering & Data Signals |
| :--- | :--- | :--- |
| **`languages.svg`** | **Code, Notebooks & Production Engineering** | **Jupyter Deflator** (strips embedded base64 chart bloat to reveal true language proportions), **Exploration (EDA) vs. Production** split, and **Engineering Hygiene** (`Lockfiles`, `CI/CD`, `Test Suites`, `Docker`). |
| **`stack.svg`** | **Data, ML & Analytics Ecosystem** | Automated dependency scanner across **3 Pillars**: *Data & Analytics* (`Pandas/Polars`, `SQL/DuckDB`, `Spark`, `dbt`), *Machine Learning & AI* (`Scikit-Learn`, `PyTorch`, `XGBoost`, `Hugging Face`), and *Viz & Apps* (`Plotly`, `Streamlit`, `Matplotlib`, `ggplot2`). |
| **`kaggle.svg`** | **Kaggle Competitive Standing & Medals** | Hybrid Headless Chrome CDP + API v1 collector displaying Overall Tier (*Novice* → *Grandmaster*), global peak rank, cumulative medals (**Gold, Silver, Bronze**), and progression across *Notebooks, Datasets, Competitions, and Discussions*. |
| **`impact.svg`** | **Open-Source & AI Community Impact** | Unifies **GitHub Reach** (Annual Contributions, YoY velocity, Stars, Public Repos, Forks, 14d Views) with **Hugging Face Hub Reach** (Published Models, Datasets, Spaces, 30d Downloads, Community Likes). |

> **Recruiter-First Graceful Degradation:** `dataprofile` never embarrasses the user by highlighting zeros. If a practitioner has `0 Gold` medals on Kaggle or `0` Hugging Face models, the layout automatically pivots its KPIs to highlight active progression tracks, total medals, or GitHub engineering velocity.

---

## 🎨 The 4 Visualization Schools of Thought (`School`)

Every run generates all panels across **4 classic Data Visualization Schools** in dedicated directories (`generated/<school>/<panel>.svg`). Notably, **none of the generated SVGs contain meta-labels mentioning the school's name**—the aesthetic speaks entirely for itself. You can adopt a single school or mix-and-match across panels.

---

### I. `generated/tufte/` — Analytical Minimalism & Maximum Data-Ink Ratio
*Inspired by **Edward Tufte** (The Visual Display of Quantitative Information, 1983).*
Zero outer borders or "cages", classical serif typography with tabular numerals, **Dot-Dash Plots** for language proportions, **Supertables** for Kaggle standing, and inline **Sparklines** with interquartile bands.

<p align="center">
  <img src="generated/tufte/languages.svg" width="48%" alt="Tufte Languages Panel" />
  <img src="generated/tufte/stack.svg" width="48%" alt="Tufte Stack Panel" />
</p>
<p align="center">
  <img src="generated/tufte/kaggle.svg" width="48%" alt="Tufte Kaggle Panel" />
  <img src="generated/tufte/impact.svg" width="48%" alt="Tufte Impact Panel" />
</p>

---

### II. `generated/few/` — Executive Telemetry & Bullet Graphs
*Inspired by **Stephen Few** (Information Dashboard Design, 2006).*
Built for 3-second executive scanning by technical recruiters and engineering managers. Features top **KPI telemetry strips**, **Stephen Few Bullet Graphs** (qualitative background range bands + inner performance bar + vertical target threshold marker), and pre-attentive color highlighting strictly on category leaders.

<p align="center">
  <img src="generated/few/languages.svg" width="48%" alt="Few Languages Panel" />
  <img src="generated/few/stack.svg" width="48%" alt="Few Stack Panel" />
</p>
<p align="center">
  <img src="generated/few/kaggle.svg" width="48%" alt="Few Kaggle Panel" />
  <img src="generated/few/impact.svg" width="48%" alt="Few Impact Panel" />
</p>

---

### III. `generated/cairo/` — Functional Art & Editorial Data Journalism
*Inspired by **Alberto Cairo** (The Functional Art, 2012; The Truthful Art, 2016).*
Combines a bold crimson editorial header rule with a **thesis-driven serif headline**, **10×5 Waffle Charts** for Exploration vs. Production structure, **Lollipop Charts** with global benchmark reference lines, and an explicit **Annotation Layer (`↳`)** explaining context to non-specialist readers.

<p align="center">
  <img src="generated/cairo/languages.svg" width="48%" alt="Cairo Languages Panel" />
  <img src="generated/cairo/stack.svg" width="48%" alt="Cairo Stack Panel" />
</p>
<p align="center">
  <img src="generated/cairo/kaggle.svg" width="48%" alt="Cairo Kaggle Panel" />
  <img src="generated/cairo/impact.svg" width="48%" alt="Cairo Impact Panel" />
</p>

---

### IV. `generated/holmes/` — Pictorial Isotypes & Visual Storytelling
*Inspired by **Nigel Holmes** (Designer's Guide to Creating Charts and Diagrams, 1984) & Otto Neurath's ISOTYPE.*
High-impact magazine header ribbons, **Pictorial Unit Blocks** (1 block = 10% codebase share), **Stepped Podium Climb** toward Kaggle Grandmaster with trophy iconography, and visual engineering badges.

<p align="center">
  <img src="generated/holmes/languages.svg" width="48%" alt="Holmes Languages Panel" />
  <img src="generated/holmes/stack.svg" width="48%" alt="Holmes Stack Panel" />
</p>
<p align="center">
  <img src="generated/holmes/kaggle.svg" width="48%" alt="Holmes Kaggle Panel" />
  <img src="generated/holmes/impact.svg" width="48%" alt="Holmes Impact Panel" />
</p>

---

## ⚡ 2-Minute Setup Guide (Zero Code Changes Needed)

`dataprofile` is **100% repository-agnostic**. You never need to edit Go source code to use it for your own profile:

### Step 1: Fork this repository & edit `profile.yaml`
Open [`profile.yaml`](profile.yaml) in the root of your fork and set your public usernames and preferred locale (`"en"` or `"pt-BR"`):

```yaml
subject:
  name: "Your Name"
  role: "Data Scientist & ML Engineer"
  locale: "en" # "en" or "pt-BR"

sources:
  github:
    username: "your-github-username"
    include_forks: false
    exclude_langs: ["HTML", "CSS", "TeX"]
  kaggle:
    username: "your-kaggle-username"
    enabled: true
    mode: "hybrid"
  huggingface:
    username: "your-hf-username"
    enabled: true
```

### Step 2: Configure Optional GitHub Secrets
In your forked repository, go to **Settings → Secrets and variables → Actions** and add (optional, for private repo counts, traffic views, and Kaggle API v1 enrichment):
- `ACCESS_TOKEN`: A GitHub Personal Access Token (`read:user`, `repo`). If omitted, the workflow automatically falls back to the built-in `GITHUB_TOKEN`.
- `KAGGLE_USERNAME` & `KAGGLE_KEY`: Optional Kaggle API credentials. Even without them, the `hybrid` collector extracts your public Kaggle tiers, ranks, and medals via headless Chrome on the GitHub Actions runner.

> 🔒 **Security Guarantee:** Never commit tokens to `profile.yaml`. Local `.env` files are strictly gitignored via [`.gitignore`](.gitignore); see [`.env.example`](.env.example) for local CLI testing.

### Step 3: Embed in your Profile `README.md`
The GitHub Action ([`.github/workflows/generate.yml`](.github/workflows/generate.yml)) runs daily and pushes fresh SVGs to the orphan `output` branch. Copy the snippet for your favorite school into your `<username>/<username>` profile `README.md`:

```markdown
<p align="center">
  <img src="https://raw.githubusercontent.com/<your-username>/dataprofile/output/generated/few/languages.svg" width="48%" />
  <img src="https://raw.githubusercontent.com/<your-username>/dataprofile/output/generated/few/stack.svg" width="48%" />
</p>
<p align="center">
  <img src="https://raw.githubusercontent.com/<your-username>/dataprofile/output/generated/few/kaggle.svg" width="48%" />
  <img src="https://raw.githubusercontent.com/<your-username>/dataprofile/output/generated/few/impact.svg" width="48%" />
</p>
```
*(Replace `/few/` with `/tufte/`, `/cairo/`, or `/holmes/`—or mix different schools across panels!)*

---

## 🛠️ Local CLI Usage & Architecture

Built as a fast, zero-external-dependency Go static binary:

```bash
# Generate all 20 SVGs (4 Schools × 5 Panels) using live APIs + profile.yaml
go run ./cmd/dataprofile generate --config=profile.yaml

# Generate deterministic showcase SVGs in offline/demo mode
go run ./cmd/dataprofile generate --demo

# Render in Brazilian Portuguese (pt-BR)
go run ./cmd/dataprofile generate --demo --locale=pt-BR

# Validate XML well-formedness, WCAG contrast rules & absence of meta-labels
go run ./cmd/dataprofile validate

# Run the automated unit & integration test suite
go test -v ./...
```

### Project Layout & Domain Documentation

```text
dataprofile/
├── cmd/dataprofile/main.go            # CLI entrypoint (generate & validate subcommands)
├── internal/
│   ├── domain/model.go                # Canonical domain model (School, Panel, Palette, PortfolioData)
│   ├── config/                        # Declarative profile.yaml & .env secret loader
│   ├── normalizer/                    # Jupyter .ipynb base64 deflator & 3-pillar Stack detector
│   ├── sources/{github,kaggle,huggingface}/ # Concurrent Tri-Source collectors
│   └── render/                        # Parameterized SVG engine (4 Schools × 5 Panels, EN/PT-BR)
├── prototypes/                        # Interactive HTML/SVG design prototypes (prototypes/index.html)
├── generated/{tufte,few,cairo,holmes}/# Output SVG panels ready for README embedding
├── GLOSSARY.md                        # Ubiquitous domain language definitions
└── docs/adr/                          # Architectural Decision Records (ADR-0001 Go, ADR-0002 Kaggle)
```

---

## 🤝 Contributing

Contributions from data practitioners, information designers, and Go engineers are warmly welcomed!
1. Read [`GLOSSARY.md`](GLOSSARY.md) and [`docs/adr/`](docs/adr/) to familiarize yourself with our canonical domain terms (`Subject`, `Source`, `Panel`, `School`, `Palette`, `Profile Manifest`).
2. **Adding a new Visualization School (`School`):** Add the constant in [`internal/domain/model.go`](internal/domain/model.go) and implement its 5 panel renderers in [`internal/render/renderer.go`](internal/render/renderer.go). All generated SVGs must pass `go run ./cmd/dataprofile validate` (well-formed XML, zero meta-labels, and high-contrast legibility on both `#ffffff` and `#0d1117`).
3. **Adding a new Ecosystem Library:** Extend `canonicalStackRules` in [`internal/normalizer/notebooks.go`](internal/normalizer/notebooks.go) and add a test case in [`internal/normalizer/notebooks_test.go`](internal/normalizer/notebooks_test.go).
