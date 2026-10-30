<div align="center">

# DataProfile

**Notebook-normalized portfolio telemetry for data scientists, ML engineers, analysts, and researchers**

[![Go 1.23+](https://img.shields.io/badge/Go-1.23%2B-00ADD8?logo=go&logoColor=white)](go.mod)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Visualization Schools](https://img.shields.io/badge/Schools-Tufte%20%7C%20Few%20%7C%20Cairo%20%7C%20Holmes-0284c7)](#the-four-visualization-schools-school)
[![Sources](https://img.shields.io/badge/Sources-GitHub%20%7C%20Kaggle%20%7C%20Hugging%20Face-d97706)](#the-nine-portfolio-panels-panel)

<p align="center">
  <img src="https://raw.githubusercontent.com/pablodiegoo/dataprofile/output/showcase/few/summary.svg" width="48%" alt="Executive Summary (Stephen Few Grammar)" />
  <img src="https://raw.githubusercontent.com/pablodiegoo/dataprofile/output/showcase/tufte/summary.svg" width="48%" alt="Executive Summary (Edward Tufte Grammar)" />
</p>

</div>

> [!WARNING]
> Upvote my post please: https://www.kaggle.com/discussions/general/754125

## Contents

1. [Why DataProfile exists](#why-dataprofile-exists)
2. [How DataProfile works and keeps your images updated](#how-dataprofile-works-and-keeps-your-images-updated)
3. [The nine portfolio panels (`Panel`)](#the-nine-portfolio-panels-panel)
4. [The four visualization schools (`School`)](#the-four-visualization-schools-school)
5. [Choosing how to configure and run DataProfile](#choosing-how-to-configure-and-run-dataprofile)
   - [Mode A: Fork + local `.env` + `dataprofile sync` (clean `main` via CLI)](#mode-a-fork--local-env--dataprofile-sync-clean-main-via-cli)
   - [Mode B: Fork + GitHub web UI variables and secrets (clean `main` without local tools)](#mode-b-fork--github-web-ui-variables-and-secrets-clean-main-without-local-tools)
   - [Mode C: Fork + editing `profile.yaml` on `main`](#mode-c-fork--editing-profileyaml-on-main)
   - [Mode D: GitHub Marketplace Action in `<username>/<username>` (no fork)](#mode-d-github-marketplace-action-in-usernameusername-no-fork)
6. [Authentication tokens, private metrics, and Kaggle modes](#authentication-tokens-private-metrics-and-kaggle-modes)
7. [Embedding the updated SVGs in your profile `README.md`](#embedding-the-updated-svgs-in-your-profile-readmemd)
8. [Local CLI reference and repository layout](#local-cli-reference-and-repository-layout)
9. [Contributing, security, and license](#contributing-security-and-license)

## Why DataProfile exists

Standard GitHub profile generators were built for general software engineers. When applied to someone who works with data, machine learning, or quantitative research, those generators misrepresent the portfolio for three concrete reasons.

First, GitHub Linguist measures language share by raw file bytes, which includes embedded base64 PNG plots, HTML table dumps, and execution metadata stored inside `.ipynb` files. Because a single Jupyter notebook with three `matplotlib` charts can weigh 4 MB while 15,000 lines of `.py`, `.sql`, and `.go` modules weigh less than 500 KB, the profile ends up labeled as 94% Jupyter Notebook. That byte distortion hides the actual engineering work, so DataProfile parses `.ipynb` JSON trees and strips cell output payloads before counting bytes (*Notebook Deflation*), separating exploratory notebooks (EDA) from production modules.

Second, knowing that a practitioner writes Python or R says very little about their domain focus. What distinguishes a data engineer from a computer vision researcher is the set of libraries declared inside `pyproject.toml`, `requirements.txt`, `environment.yml`, and `DESCRIPTION` (such as Polars, DuckDB, PyTorch, Scikit-Learn, dbt, or Streamlit), together with engineering habits like dependency lockfiles, CI workflows, automated tests, and Dockerfiles. Scanning those manifests across every *Analytical Workspace* (both standalone repositories and subdirectories inside monorepos) surfaces the tools actually used in projects.

Third, a data professional's public work rarely lives on GitHub alone. Competition rankings, datasets, and notebooks often sit on Kaggle, while pretrained weights and interactive demos live on Hugging Face Hub. DataProfile collects signals from all three *Sources* concurrently and renders them through four classic schools of information design (`tufte`, `few`, `cairo`, and `holmes`).

## How DataProfile works and keeps your images updated

DataProfile runs as a zero-dependency Go binary inside GitHub Actions (or on your local machine) so that your profile `README.md` always displays current metrics without cluttering your commit history on `main`.

Every execution follows four steps:

1. Read the declarative [`profile.yaml`](profile.yaml) (*Profile Manifest*) and layer any environment variables or secrets configured in GitHub Actions (or in a local `.env` file) on top of it.
2. Query GitHub, Kaggle, and Hugging Face Hub concurrently within a 40-second timeout, applying *Notebook Deflation* to `.ipynb` files and scanning dependency manifests across your repositories.
3. Render and validate 36 SVG files (9 *Panels* across 4 *Schools*), verifying XML well-formedness, WCAG AA contrast across light and dark backgrounds, and the absence of internal school meta-labels.
4. Force-push the generated files to an orphan branch named `output`, completely separated from your `main` code branch.

That orphan `output` branch is where your updated images live. In a fork of this repository, the workflow ([`.github/workflows/generate.yml`](.github/workflows/generate.yml)) publishes two directories and a visual catalog to `output`: `output/generated/<school>/<panel>.svg` holds the *Subject*'s personal portfolio panels rendered with your live metrics, usernames, and bio overrides; `output/showcase/<school>/<panel>.svg` holds the deterministic reference *Showcase* rendered from the default `profile.yaml` (`--demo --no-env`); and `output/README.md` at the root of the branch previews all 36 panels alongside ready-to-copy Markdown snippets already filled in with your repository URL.

Once your `<username>/<username>` profile `README.md` points its `<img src="https://raw.githubusercontent.com/.../output/generated/...">` tags to that branch, you never need to edit or commit to your profile README again just to refresh numbers. The `output` branch updates automatically on four occasions: every day at `06:00 UTC` via the scheduled cron job in GitHub Actions; right after you run `go run ./cmd/dataprofile sync` locally to push `.env` changes; on any commit that modifies `profile.yaml` on `main`; or on demand when you open the **Actions** tab on GitHub, select **Generate DataProfile SVGs**, and click **Run workflow** (`workflow_dispatch`).

Two GitHub platform rules affect scheduled runs on forks. First, GitHub disables workflows on fresh forks until you open the **Actions** tab and click **"I understand my workflows, go ahead and enable them"**. Second, GitHub pauses scheduled cron workflows if a repository has no activity for 60 consecutive days; if that happens, clicking **Run workflow** (or re-enabling the workflow in the Actions tab) resumes daily updates immediately. After each run pushes to `output`, GitHub's image proxy (Camo) refreshes the SVG in your profile as soon as its short cache window expires.

## The nine portfolio panels (`Panel`)

Each run compiles metrics from GitHub, Kaggle, and Hugging Face Hub into nine modular panels (four analytical panels, three identity/narrative panels, and two composite panels, totaling 36 SVGs across the four schools). Every file uses a transparent canvas (`fill="none"`) and a dual-surface *Palette* calibrated for WCAG AA contrast on both GitHub Light (`#ffffff`) and GitHub Dark (`#0d1117`).

| Panel File | Domain Focus | Key Engineering & Data Signals |
| :--- | :--- | :--- |
| **`languages.svg`** | Code, Notebooks & Engineering Hygiene | Strips embedded base64 chart bloat from `.ipynb` files to recover true source code proportions, measures the Exploration (EDA) vs. Production split, and counts repositories with lockfiles, CI/CD, test suites, and Docker. |
| **`stack.svg`** | Data, ML & Analytics Ecosystem | Scans dependency manifests across three pillars: Data & Analytics (`Pandas/Polars`, `SQL/DuckDB`, `Apache Spark`, `dbt/Airflow`), Machine Learning & AI (`Scikit-Learn`, `PyTorch`, `XGBoost/LightGBM`, `Hugging Face`), and Visualization & Data Apps (`Plotly/Dash`, `Streamlit`, `Matplotlib/Seaborn`, `ggplot2`). |
| **`kaggle.svg`** | Kaggle Standing & Medals | Combines headless Chrome CDP extraction with optional Kaggle REST API v1 calls to display Overall Tier (`Novice` to `Grandmaster`), global peak rank, cumulative medals (Gold, Silver, Bronze), and percentile standing across Notebooks, Datasets, Competitions, and Discussions. |
| **`impact.svg`** | Open-Source & AI Community Impact | Joins GitHub reach (annual contributions, 52-week trend, stars, public and contributed repos, forks, 14-day views) with Hugging Face Hub metrics (published models, datasets, spaces, 30-day downloads, and community likes). |
| **`summary.svg`** | Consolidated 4-Quadrant Overview | Synthesizes code distribution, stack leaders, Kaggle standing, and community reach into a single executive panel (`440px` grid or `100%` width). |
| **`header.svg`** | Animated Portfolio Header (`900px`) | Pure SVG 1.1 SMIL typing/headline rotation (`subject.headlines`), role subtitle, and live telemetry strip across the top of a profile README. |
| **`bio.svg` / `bio-full.svg`** | Elastic Career Bio (`440px` / `900px`) | Deterministic Go word-wrapping and dynamic `viewBox` height reflow for multi-paragraph career narratives (`subject.bio`), with zero blank space on short bios and zero overlap on long ones. |
| **`full.svg`** | Full Portfolio Composite (`900px`) | Combines the animated `header`, elastic `bio-full`, and the 2×2 grid of all four analytical panels into a single self-contained SVG. |

When a metric is zero, the renderer reallocates that space to active signals. For instance, if a practitioner has zero Gold medals on Kaggle or no published models on Hugging Face Hub, the panel shifts those coordinates to active progression tracks or GitHub engineering velocity instead of printing empty rows.

## The four visualization schools (`School`)

None of the generated SVGs print the name of the school or its author inside the artwork. That omission allows you to adopt a single visual grammar across your entire profile or mix panels from different schools in the same `README.md`.

### I. `tufte` (Edward Tufte)

Inspired by **Edward Tufte** (*The Visual Display of Quantitative Information*, 1983), this grammar removes outer boxes and grid cages so that every line serves the data. It uses serif typography with tabular numerals, dot-dash plots for language share, supertables for Kaggle percentiles, and inline sparklines with interquartile shading.

<p align="center">
  <img src="https://raw.githubusercontent.com/pablodiegoo/dataprofile/output/showcase/tufte/languages.svg" width="48%" alt="Tufte Languages Panel" />
  <img src="https://raw.githubusercontent.com/pablodiegoo/dataprofile/output/showcase/tufte/stack.svg" width="48%" alt="Tufte Stack Panel" />
</p>
<p align="center">
  <img src="https://raw.githubusercontent.com/pablodiegoo/dataprofile/output/showcase/tufte/kaggle.svg" width="48%" alt="Tufte Kaggle Panel" />
  <img src="https://raw.githubusercontent.com/pablodiegoo/dataprofile/output/showcase/tufte/impact.svg" width="48%" alt="Tufte Impact Panel" />
</p>

### II. `few` (Stephen Few)

Inspired by **Stephen Few** (*Information Dashboard Design*, 2006), this layout is organized for three-second scanning by technical recruiters and engineering managers. It pairs top KPI strips with bullet graphs (qualitative background bands, an inner performance bar, and a vertical target threshold marker) and reserves accent color for category leaders.

<p align="center">
  <img src="https://raw.githubusercontent.com/pablodiegoo/dataprofile/output/showcase/few/languages.svg" width="48%" alt="Few Languages Panel" />
  <img src="https://raw.githubusercontent.com/pablodiegoo/dataprofile/output/showcase/few/stack.svg" width="48%" alt="Few Stack Panel" />
</p>
<p align="center">
  <img src="https://raw.githubusercontent.com/pablodiegoo/dataprofile/output/showcase/few/kaggle.svg" width="48%" alt="Few Kaggle Panel" />
  <img src="https://raw.githubusercontent.com/pablodiegoo/dataprofile/output/showcase/few/impact.svg" width="48%" alt="Few Impact Panel" />
</p>

### III. `cairo` (Alberto Cairo)

Inspired by **Alberto Cairo** (*The Functional Art*, 2012; *The Truthful Art*, 2016), this school treats the portfolio like an editorial data-journalism piece. It combines a crimson top rule with a thesis-driven serif headline, a 10×5 waffle chart for code structure, lollipop charts with dashed benchmark lines, and an annotation layer (`↳`) that explains what the numbers mean in context.

<p align="center">
  <img src="https://raw.githubusercontent.com/pablodiegoo/dataprofile/output/showcase/cairo/languages.svg" width="48%" alt="Cairo Languages Panel" />
  <img src="https://raw.githubusercontent.com/pablodiegoo/dataprofile/output/showcase/cairo/stack.svg" width="48%" alt="Cairo Stack Panel" />
</p>
<p align="center">
  <img src="https://raw.githubusercontent.com/pablodiegoo/dataprofile/output/showcase/cairo/kaggle.svg" width="48%" alt="Cairo Kaggle Panel" />
  <img src="https://raw.githubusercontent.com/pablodiegoo/dataprofile/output/showcase/cairo/impact.svg" width="48%" alt="Cairo Impact Panel" />
</p>

### IV. `holmes` (Nigel Holmes & ISOTYPE)

Inspired by **Nigel Holmes** (*Designer's Guide to Creating Charts and Diagrams*, 1984) and Otto Neurath's ISOTYPE system, this grammar uses magazine header ribbons, 10-block pictorial unit charts (where each block equals 10% of the codebase), a stepped podium showing progression toward Kaggle Grandmaster, and visual badges for engineering practices.

<p align="center">
  <img src="https://raw.githubusercontent.com/pablodiegoo/dataprofile/output/showcase/holmes/languages.svg" width="48%" alt="Holmes Languages Panel" />
  <img src="https://raw.githubusercontent.com/pablodiegoo/dataprofile/output/showcase/holmes/stack.svg" width="48%" alt="Holmes Stack Panel" />
</p>
<p align="center">
  <img src="https://raw.githubusercontent.com/pablodiegoo/dataprofile/output/showcase/holmes/kaggle.svg" width="48%" alt="Holmes Kaggle Panel" />
  <img src="https://raw.githubusercontent.com/pablodiegoo/dataprofile/output/showcase/holmes/impact.svg" width="48%" alt="Holmes Impact Panel" />
</p>

## Choosing how to configure and run DataProfile

Because the engine layers environment variables over [`profile.yaml`](profile.yaml) (with precedence: process environment variables > local `.env` file > `profile.yaml`), you can operate DataProfile in four ways depending on whether you want a fork, whether you have Go installed locally, and whether you want your `main` branch free of personal curriculum commits.

| Configuration Mode | Requires Fork? | Requires Local Go & `gh` CLI? | Keeps `main` Free of Personal Commits? | Where Your Profile Manifest & Overrides Live |
| :--- | :--- | :--- | :--- | :--- |
| **[Mode A](#mode-a-fork--local-env--dataprofile-sync-clean-main-via-cli)** (Recommended for local CLI setup) | Yes | Yes | Yes (avoids merge conflicts when syncing upstream) | Local gitignored `.env` synced to GitHub Actions via `dataprofile sync` |
| **[Mode B](#mode-b-fork--github-web-ui-variables-and-secrets-clean-main-without-local-tools)** (Recommended for browser-only setup) | Yes | No | Yes | GitHub Actions Repository Variables and Secrets in the web UI |
| **[Mode C](#mode-c-fork--editing-profileyaml-on-main)** | Yes | No | No (commits changes to `profile.yaml` on `main`) | Tracked [`profile.yaml`](profile.yaml) file in your fork |
| **[Mode D](#mode-d-github-marketplace-action-in-usernameusername-no-fork)** | No | No | Tracked inside `<username>/<username>` | `profile.yaml` and `.github/workflows/dataprofile.yml` in your profile repo |

### Mode A: Fork + local `.env` + `dataprofile sync` (clean `main` via CLI)

This mode lets you manage your bio, headlines, usernames, and tokens in a single local `.env` file while leaving `profile.yaml` on `main` untouched, which prevents merge conflicts whenever you pull upstream updates.

1. Click **Fork** at the top-right of this repository and clone your fork locally.
2. Open the **Actions** tab in your fork on GitHub and click **"I understand my workflows, go ahead and enable them"**.
3. Copy [`.env.example`](.env.example) to `.env` (which is ignored by Git) and fill in your overrides. For lists such as `SUBJECT_HEADLINES`, `SUBJECT_BIO`, `GITHUB_EXCLUDE_REPOS`, or `GITHUB_EXCLUDE_LANGS`, use bracketed semicolon-delimited syntax `[item 1; item 2]` (escaping literal semicolons as `\;`):

```dotenv
SUBJECT_NAME="Your Name"
SUBJECT_ROLE="Data Scientist & Analytics Engineer"
SUBJECT_LOCALE="en" # "en" or "pt-BR"
SUBJECT_HEADLINES=[Welcome people and mystical beings.; I'm Your Name! 👋; Working with data 10+ years]
SUBJECT_BIO=[First paragraph of your career bio.; Second paragraph of your career bio.]

GITHUB_USERNAME="your-github-username"
GITHUB_INCLUDE_FORKS=false
GITHUB_INCLUDE_PRIVATE_COUNT=true

KAGGLE_USERNAME="your-kaggle-username"
KAGGLE_ENABLED=true
KAGGLE_MODE="hybrid"

HF_USERNAME="your-hf-username"
HF_ENABLED=true

# Optional secrets (leave empty if not needed)
ACCESS_TOKEN=
KAGGLE_KEY=
HF_TOKEN=
```

4. Authenticate the GitHub CLI (`gh auth login` if not already logged in) and run the sync command from your terminal:

```bash
# Preview the exact gh variable/secret commands with secret values masked as ***
go run ./cmd/dataprofile sync --dry-run

# Push variables and secrets to GitHub Actions and trigger generate.yml immediately
go run ./cmd/dataprofile sync
```

Because GitHub Actions reserves the `GITHUB_` prefix for internal runner variables, `dataprofile sync` automatically maps `GITHUB_USERNAME` and other `GITHUB_*` keys from your `.env` file to `GH_USERNAME` and `GH_*` repository variables on GitHub, omits `GITHUB_TOKEN` (since Actions provides it natively), uploads credential tokens (`ACCESS_TOKEN`, `KAGGLE_KEY`, `HF_TOKEN`) via `gh secret set`, and triggers `generate.yml`. Within roughly a minute, your 36 SVGs and personalized `README.md` catalog will be live on the `output` branch of your fork.

### Mode B: Fork + GitHub web UI variables and secrets (clean `main` without local tools)

If you want to keep `main` clean without cloning the repository or installing Go and `gh` locally, you can set the exact same overrides directly in the GitHub web interface.

1. Click **Fork** at the top-right of this repository.
2. Open the **Actions** tab in your fork and click **"I understand my workflows, go ahead and enable them"**.
3. Open **Settings -> Secrets and variables -> Actions**, switch to the **Variables** tab, and click **New repository variable** for the fields you wish to override: identity variables (`SUBJECT_NAME`, `SUBJECT_ROLE`, `SUBJECT_LOCALE`, plus `SUBJECT_HEADLINES` and `SUBJECT_BIO` formatted as `[Paragraph 1; Paragraph 2]`); GitHub source variables using the `GH_` prefix required by the GitHub Actions UI (`GH_USERNAME`, `GH_INCLUDE_FORKS`, `GH_INCLUDE_PRIVATE_COUNT`, `GH_EXCLUDE_REPOS`, `GH_EXCLUDE_LANGS`); Kaggle and Hugging Face variables (`KAGGLE_USERNAME`, `KAGGLE_ENABLED`, `KAGGLE_MODE`, `HF_USERNAME`, `HF_ENABLED`); and optional output selectors (`OUTPUT_SCHOOLS`, `OUTPUT_PANELS`).
4. Add `ACCESS_TOKEN`, `KAGGLE_KEY`, or `HF_TOKEN` under the **Secrets** tab if you want private repository counts or authenticated API enrichment.
5. Open the **Actions** tab, select **Generate DataProfile SVGs** on the left sidebar, and click **Run workflow**.

### Mode C: Fork + editing `profile.yaml` on `main`

If you prefer keeping your public configuration in a single YAML file inside your fork and do not mind commits on `main`, edit [`profile.yaml`](profile.yaml) (*Profile Manifest*) directly.

1. Fork this repository and enable workflows in the **Actions** tab.
2. Edit [`profile.yaml`](profile.yaml) on `main` (either in the GitHub web editor or via git) with your name, role, bio, and usernames across `github`, `kaggle`, and `huggingface`.
3. Commit the change to `main`. Because `profile.yaml` is listed in the `on.push.paths` trigger of [`.github/workflows/generate.yml`](.github/workflows/generate.yml), committing that file automatically starts a run that publishes your updated SVGs to the `output` branch.

### Mode D: GitHub Marketplace Action in `<username>/<username>` (no fork)

If you do not want a separate `dataprofile` fork, you can run `pablodiegoo/dataprofile@v1` directly inside your `<username>/<username>` profile repository and publish the generated files to an orphan `output` branch in that same repository.

1. Create a `profile.yaml` file at the root of your `<username>/<username>` repository:

```yaml
subject:
  name: "Your Name"
  role: "Data Scientist & ML Engineer"
  locale: "en" # "en" or "pt-BR"
  headlines:
    - "Data Scientist & ML Engineer"
    - "Building production ML systems and analytical pipelines"
  bio:
    - "Transforming complex datasets into predictive models and reproducible data tooling."

sources:
  github:
    username: "your-github-username"
    include_forks: false
    include_private_count: true
    exclude_langs: ["HTML", "CSS", "TeX"]
  kaggle:
    username: "your-kaggle-username"
    enabled: true
    mode: "hybrid"
  huggingface:
    username: "your-hf-username"
    enabled: true
```

2. Create `.github/workflows/dataprofile.yml` in your `<username>/<username>` repository, including both the generator step and the step that publishes the `generated/` folder to the `output` branch:

```yaml
name: Generate DataProfile SVGs

on:
  schedule:
    - cron: "0 6 * * *"
  workflow_dispatch:
  push:
    branches: [main]
    paths: ["profile.yaml", ".github/workflows/dataprofile.yml"]

permissions:
  contents: write

jobs:
  dataprofile:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Generate and validate DataProfile SVGs
        uses: pablodiegoo/dataprofile@v1
        with:
          config: "profile.yaml"
          output: "generated"
          github-token: ${{ secrets.ACCESS_TOKEN || github.token }}
          kaggle-username: ${{ secrets.KAGGLE_USERNAME }}
          kaggle-key: ${{ secrets.KAGGLE_KEY }}

      - name: Publish SVGs to orphan output branch
        run: |
          git config --global user.name "github-actions[bot]"
          git config --global user.email "41898282+github-actions[bot]@users.noreply.github.com"
          TMP_DIR=$(mktemp -d)
          cp -r generated "$TMP_DIR/generated"
          git checkout --orphan output-temp
          git rm -rf . >/dev/null 2>&1 || true
          cp -r "$TMP_DIR/generated" ./generated
          if [ -f ./generated/README.md ]; then
            cp ./generated/README.md ./README.md
          fi
          git add generated README.md
          git commit -m "chore(svg): update generated DataProfile panels"
          git branch -M output-temp output
          git push -f origin output
```

When using Mode D, your SVGs live inside your `<username>/<username>` repository on the `output` branch, so your image URLs will follow the pattern `https://raw.githubusercontent.com/<your-username>/<your-username>/output/generated/<school>/<panel>.svg`.

## Authentication tokens, private metrics, and Kaggle modes

DataProfile works out of the box using only the automatic `GITHUB_TOKEN` provided by GitHub Actions. Adding optional personal tokens unlocks private repository counts, 14-day repository traffic views, and authenticated API enrichment while preserving strict privacy boundaries.

| Secret Name | Required? | Recommended Scopes / Format | What It Unlocks |
| :--- | :--- | :--- | :--- |
| **`GITHUB_TOKEN`** | Automatic in Actions | Provided natively by GitHub Actions | Public repositories, public language distributions, dependency manifests, stars, forks, and contribution calendar. |
| **`ACCESS_TOKEN`** | Optional | GitHub Personal Access Token (`read:user`, `repo`) | Includes each private *Owned Repository* and *Contributed Repository* in aggregate counts, language shares, and stack detection, plus 14-day repository traffic views (`RepoViews14d`). |
| **`KAGGLE_KEY`** | Optional | `KGAT_...` Bearer token or classic `kaggle.json` key (paired with `KAGGLE_USERNAME`) | Enriches Kaggle metrics through the official REST API v1 (dataset/notebook votes, downloads, and competition standings). |
| **`HF_TOKEN`** | Optional | Hugging Face Hub read token (`hf_...`) | Includes private models, datasets, and spaces in aggregate Hugging Face counts. |

To protect proprietary work when `ACCESS_TOKEN`, `KAGGLE_KEY`, or `HF_TOKEN` are supplied, collection follows a three-tier boundary (detailed in [`docs/adr/0003-privacy-preserving-three-tier-collection.md`](docs/adr/0003-privacy-preserving-three-tier-collection.md)): every *Owned Repository* (a public or private non-fork repository belonging to the *Subject*) undergoes deep tree and manifest inspection in memory across all *Analytical Workspaces*; every *Contributed Repository* (an external or organization repository where the *Subject* authored verified commits, pull requests, or reviews) counts toward repository reach and contribution velocity without claiming upstream code bytes or stars; and every *Passive Fork* is excluded by default (`include_forks: false`) so third-party codebases do not distort language shares or star totals. All private signals are reduced to aggregate percentages and counts before rendering, so repository names, file paths, and source code never appear in the generated SVGs.

For Kaggle (`sources.kaggle.mode` or `KAGGLE_MODE`), three collection modes are available:

| Kaggle Mode | Credentials Needed | Collection Mechanism |
| :--- | :--- | :--- |
| **`hybrid`** (default) | None required (`KAGGLE_KEY` optional) | Extracts public Overall Tier, global peak rank, and medal counts via headless Chrome speaking the Chrome DevTools Protocol (using the browser pre-installed on GitHub Actions runners to bypass Google reCAPTCHA challenges on plain HTTP requests), and merges official REST API v1 metrics whenever `KAGGLE_KEY` is configured. |
| **`browser`** | None | Runs only headless Chrome CDP extraction against the public Kaggle profile page. |
| **`api`** | `KAGGLE_USERNAME` and `KAGGLE_KEY` | Queries only the official Kaggle REST API v1 endpoints. |

Never commit tokens to `profile.yaml`. Local `.env` files are excluded via [`.gitignore`](.gitignore).

## Embedding the updated SVGs in your profile `README.md`

After your workflow finishes its first run, switch the branch selector on GitHub from `main` to `output`. At the root of the `output` branch, GitHub renders `output/README.md`, which displays all 36 generated SVGs for the *Subject* alongside ready-to-copy Markdown blocks already containing your repository path.

In Modes A, B, and C (where the workflow runs inside your `dataprofile` fork), the raw image URL points to `https://raw.githubusercontent.com/<your-username>/dataprofile/output/generated/<school>/<panel>.svg`. In Mode D (where the workflow runs inside your `<your-username>/<your-username>` profile repository), replace `/dataprofile/` in the URL with `/<your-username>/` (`https://raw.githubusercontent.com/<your-username>/<your-username>/output/generated/<school>/<panel>.svg`).

Copy one of the layouts below into the `README.md` of your `<your-username>/<your-username>` profile repository (replacing `<your-username>` and switching `/few/` to `/tufte/`, `/cairo/`, or `/holmes/` as preferred):

```markdown
<!-- Layout 1: Single full-portfolio composite (animated header + elastic bio + 2x2 grid) -->
<p align="center">
  <img src="https://raw.githubusercontent.com/<your-username>/dataprofile/output/generated/few/full.svg" width="100%" alt="DataProfile Full Portfolio" />
</p>

<!-- Layout 2: Modular panels (header.svg, bio-full.svg, and 2x2 analytical grid) -->
<p align="center">
  <img src="https://raw.githubusercontent.com/<your-username>/dataprofile/output/generated/few/header.svg" width="100%" alt="Animated Portfolio Header" />
</p>
<p align="center">
  <img src="https://raw.githubusercontent.com/<your-username>/dataprofile/output/generated/few/bio-full.svg" width="100%" alt="Career Bio" />
</p>
<p align="center">
  <img src="https://raw.githubusercontent.com/<your-username>/dataprofile/output/generated/few/languages.svg" width="48%" alt="Languages & Engineering" />
  <img src="https://raw.githubusercontent.com/<your-username>/dataprofile/output/generated/few/stack.svg" width="48%" alt="Data & ML Stack" />
</p>
<p align="center">
  <img src="https://raw.githubusercontent.com/<your-username>/dataprofile/output/generated/few/kaggle.svg" width="48%" alt="Kaggle Standing" />
  <img src="https://raw.githubusercontent.com/<your-username>/dataprofile/output/generated/few/impact.svg" width="48%" alt="Community Impact" />
</p>

<!-- Layout 3: Mix and match different schools in the same grid -->
<p align="center">
  <img src="https://raw.githubusercontent.com/<your-username>/dataprofile/output/generated/tufte/languages.svg" width="48%" alt="Languages (Tufte)" />
  <img src="https://raw.githubusercontent.com/<your-username>/dataprofile/output/generated/cairo/stack.svg" width="48%" alt="Stack (Cairo)" />
</p>
<p align="center">
  <img src="https://raw.githubusercontent.com/<your-username>/dataprofile/output/generated/few/kaggle.svg" width="48%" alt="Kaggle (Few)" />
  <img src="https://raw.githubusercontent.com/<your-username>/dataprofile/output/generated/holmes/impact.svg" width="48%" alt="Impact (Holmes)" />
</p>
<p align="center">
  <sub>Generated with <a href="https://github.com/pablodiegoo/dataprofile">DataProfile</a></sub>
</p>
```

## Local CLI reference and repository layout

Because the engine uses only the Go 1.23+ standard library, you can run collection, rendering, validation, synchronization, and the automated test suite locally without installing third-party Go packages:

```bash
# Generate all 36 SVG panels + generated/README.md using live APIs, profile.yaml, and .env overrides
# Flags: --config=profile.yaml --env=.env --output=generated --locale=en|pt-BR --demo --no-env
go run ./cmd/dataprofile generate --config=profile.yaml

# Sync local .env overrides and secrets to GitHub Actions and trigger workflow regeneration
# Flags: --env=.env --repo=owner/repo --workflow=generate.yml --trigger=true --dry-run
go run ./cmd/dataprofile sync

# Generate deterministic Showcase SVGs in offline demo mode (ignoring .env and environment variables)
go run ./cmd/dataprofile generate --demo --no-env --output=showcase

# Render all panels in Brazilian Portuguese (pt-BR)
go run ./cmd/dataprofile generate --demo --locale=pt-BR

# Validate XML well-formedness, dual-surface WCAG AA contrast, and absence of school meta-labels
# Flags: --config=profile.yaml --env=.env --output=generated --no-env
go run ./cmd/dataprofile validate

# Run the automated unit and integration test suite
go test -v ./...
```

```text
dataprofile/
├── cmd/dataprofile/main.go                  # CLI entrypoint (generate, validate, and sync subcommands)
├── internal/
│   ├── domain/model.go                      # Canonical domain types (School, Panel, Palette, PortfolioData)
│   ├── config/                              # Declarative profile.yaml and .env [item 1; item 2] override loader
│   ├── normalizer/                          # Jupyter .ipynb base64 deflator and 3-pillar Stack scanner
│   ├── sources/{github,kaggle,huggingface}/ # Concurrent upstream Source collectors
│   └── render/                              # Parameterized SVG engine (4 Schools x 9 Panels, en and pt-BR)
├── prototypes/                              # Interactive HTML/SVG viewer (prototypes/index.html)
├── action.yml                               # Composite GitHub Marketplace Action definition
├── profile.yaml                             # Public declarative Profile Manifest
├── .env.example                             # Template for local and CI environment overrides
├── GLOSSARY.md                              # Domain vocabulary definitions
└── docs/adr/                                # Architectural Decision Records (ADR-0001 to ADR-0005)
```

## Contributing, security, and license

Please read [`CONTRIBUTING.md`](CONTRIBUTING.md), [`GLOSSARY.md`](GLOSSARY.md), and [`docs/adr/`](docs/adr/) before opening a pull request. Community interactions follow our [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md), and credential or vulnerability concerns should be reported privately according to [`SECURITY.md`](SECURITY.md).

DataProfile is licensed under the [Apache License, Version 2.0](LICENSE). Under Section 4(d) of the license, any redistribution or derivative work must preserve the attribution notice in [`NOTICE`](NOTICE) referencing the original repository (`https://github.com/pablodiegoo/dataprofile`) and author (Pablo Diego). The serverless GitHub Actions SVG workflow pattern was inspired by prior open-source projects [`jstrieb/github-stats`](https://github.com/jstrieb/github-stats) and [`rahul-jha98/github-stats-transparent`](https://github.com/rahul-jha98/github-stats-transparent); DataProfile is an independent implementation written from scratch in Go.
