# Changelog

All notable changes to DataProfile are documented in this file. This project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Environment override engine (`internal/config`) supporting `.env` and GitHub Actions Variables/Secrets layered over `profile.yaml`, with `[item 1; item 2]` semicolon-delimited bracket list syntax (`ParseEnvList` and `ApplyEnvOverrides`) and `--no-env` Showcase isolation (`ADR-0005`).
- `dataprofile sync` CLI subcommand (`cmd/dataprofile`) with `--dry-run`, `--repo`, `--workflow`, and `--trigger` flags to synchronize local `.env` curriculum variables (`gh variable set`) and credentials (`gh secret set`) with GitHub Actions and trigger orphan `output` branch regeneration without commits on `main`.
- Dual-tree orphan `output` branch publishing (`output/showcase/` rendered with `--demo --no-env` for the project documentation Showcase and `output/generated/` rendered with live Subject overrides), removing tracked `generated/` SVGs from `main`.
- Animated portfolio `header` Panel (`900px` SMIL headline rotation), elastic `bio` (`440px`) and `bio-full` (`900px`) Panels with deterministic Go word-wrapping, and composite `full` Panel (`900px` combining `header`, `bio-full`, and the 2×2 analytical grid) across all four Schools (`ADR-0004`).
- Three-tier privacy-preserving collection (`Owned Repository`, `Contributed Repository`, `Passive Fork`), recursive monorepo `Analytical Workspace` detection, multi-year GitHub contribution velocity, and authenticated private Kaggle and Hugging Face Hub aggregation (`ADR-0003`).

## [v1.0.0] - 2026-10-28

### Added

- Core static Go 1.23 engine (`cmd/dataprofile`) with `generate` and `validate` subcommands and zero external runtime dependencies.
- Jupyter `.ipynb` base64 deflator (`internal/normalizer`) that strips embedded PNG/HTML cell outputs before computing language proportions and Exploration (EDA) versus Production code splits.
- Multi-manifest dependency scanner (`pyproject.toml`, `requirements.txt`, `environment.yml`, `Pipfile`, `DESCRIPTION`) across three pillars: Data & Analytics, Machine Learning & AI, and Visualization & Data Apps.
- Concurrent collectors for three upstream Sources: GitHub (REST + GraphQL 52-week contribution velocity), Kaggle (hybrid headless Chrome CDP + optional REST API v1 with `KGAT_` Bearer token support), and Hugging Face Hub (models, datasets, spaces, 30-day downloads, and likes).
- Parameterized SVG rendering engine (`internal/render`) producing 20 panels (4 Visualization Schools: `tufte`, `few`, `cairo`, `holmes` across 5 Panels: `languages`, `stack`, `kaggle`, `impact`, `summary`) in English (`en`) and Brazilian Portuguese (`pt-BR`).
- Universal Dual-Surface Palette (`fill="none"`) verified for WCAG AA contrast on both GitHub Light (`#ffffff`) and GitHub Dark (`#0d1117`), with Recruiter-First graceful degradation when Kaggle gold medals or Hugging Face artifacts are zero.
- Composite GitHub Action (`action.yml`) supporting both standalone Marketplace usage (`uses: pablodiegoo/dataprofile@v1`) and fork workflows, plus automated bilingual showcase `README.md` generation on the orphan `output` branch.
- Apache License 2.0 (`LICENSE`), mandatory attribution notice (`NOTICE`), embedded XML attribution comments in generated SVGs, `CITATION.cff`, and full GitHub community governance files.
