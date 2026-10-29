# Contributing to DataProfile

Contributions from data scientists, machine learning engineers, information designers, and Go developers are welcome. Whether you want to add a library to the dependency scanner, refine a chart grammar, or propose a new visualization school, this guide explains how the codebase is organized and how we validate changes before merging.

## Start with the domain language

Before writing code or opening an issue, read [`GLOSSARY.md`](GLOSSARY.md) and the Architectural Decision Records in [`docs/adr/`](docs/adr/) (`ADR-0001` through `ADR-0005`). We use canonical terms across all structs, flags, and pull requests: `Subject` (the data professional portrayed, never "User" or "Client"), `Source` (the upstream platform such as GitHub, Kaggle, or Hugging Face Hub), `Panel` (one of the nine SVG modules: `languages`, `stack`, `kaggle`, `impact`, `summary`, `header`, `bio`, `bio-full`, and `full`), `School` (a visualization paradigm such as Tufte, Few, Cairo, or Holmes), `Palette` (the dual-surface color scheme), `Profile Manifest` (`profile.yaml` with optional `.env` / CI variable overrides), and `Showcase` (deterministic reference panels rendered with `--demo --no-env`). Keeping these terms consistent prevents drift between the documentation and the Go types in [`internal/domain/model.go`](internal/domain/model.go).

## Local development setup

DataProfile compiles as a static Go 1.23+ binary using only the Go standard library. Generated SVG directories (`/generated/` and `/showcase/`) are excluded from the `main` branch via [`.gitignore`](.gitignore) and published exclusively to the orphan `output` branch (`output/showcase/<school>/<panel>.svg` and `output/generated/<school>/<panel>.svg`) so local runs and personal curriculum updates never pollute `main`. You can clone the repository and run the full test and validation suite immediately:

```bash
# Run all unit and integration tests across collectors, normalizers, config overrides, and renderers
go test -v ./...

# Render all 36 Showcase SVG panels (4 Schools x 9 Panels) using deterministic offline data and pure profile.yaml
go run ./cmd/dataprofile generate --demo --no-env --output=showcase

# Validate XML well-formedness, dual-surface contrast, and absence of school meta-labels
go run ./cmd/dataprofile validate --no-env --output=showcase

# Render in Brazilian Portuguese
go run ./cmd/dataprofile generate --demo --no-env --locale=pt-BR --output=showcase

# Preview planned GitHub Actions Variable/Secret synchronization from local .env without executing
go run ./cmd/dataprofile sync --dry-run
```

When you modify any of the four analytical panels (`languages.svg`, `stack.svg`, `kaggle.svg`, or `impact.svg`) in [`internal/render/renderer.go`](internal/render/renderer.go), regenerate the reference assets (`--demo --no-env`) and synchronize the interactive prototype copies in `prototypes/schools/<school>/`. The integration test `TestPrototypesShowcaseAndGeneratedSVGs` in [`cmd/dataprofile/main_test.go`](cmd/dataprofile/main_test.go) checks that rendered Showcase panels and `prototypes/schools/<school>/` match byte-for-byte.

## Common contribution paths

### 1. Adding a data or ML library to the Stack scanner

The `stack.svg` panel inspects `pyproject.toml`, `requirements.txt`, `environment.yml`, `Pipfile`, and `DESCRIPTION` files across each `Analytical Workspace` to detect which tools a `Subject` actually uses in production and research. If a library is missing or needs an alias, edit `canonicalStackRules` in [`internal/normalizer/notebooks.go`](internal/normalizer/notebooks.go) and add a corresponding unit test in [`internal/normalizer/notebooks_test.go`](internal/normalizer/notebooks_test.go).

### 2. Refining or adding a Visualization School (`School`)

Each `School` translates a published theory of information design into parameterized SVG geometry. To add a new `School`, define its constant in [`internal/domain/model.go`](internal/domain/model.go) and implement its nine panel renderers (`languages`, `stack`, `kaggle`, `impact`, `summary`, `header`, `bio`, `bio-full`, and `full`) in [`internal/render/renderer.go`](internal/render/renderer.go), supporting both English (`en`) and Brazilian Portuguese (`pt-BR`).

Every rendered SVG must satisfy four invariants enforced by `ValidateSVG`:

1. Valid, self-contained XML with a transparent root canvas (`fill="none"`).
2. Legible contrast on both GitHub Light (`#ffffff`) and GitHub Dark (`#0d1117`) themes, avoiding washed-out pale greys (`#e6edf3`, `#c9d1d9`, `#8b949e`, `#94a3b8`).
3. Zero visual meta-labels inside the SVG artwork mentioning the school's author (for example, never print "Edward Tufte" or "Stephen Few" inside the chart).
4. Graceful degradation when a metric is zero: if a `Subject` has zero Kaggle gold medals or zero Hugging Face models, the layout pivots to active tracks or engineering velocity rather than displaying `0 Gold` or `0 / 0 / 0`.

### 3. Improving a Source collector

Source collectors live under `internal/sources/{github,kaggle,huggingface}/`. Because GitHub Actions workflows run unattended on a daily schedule, collectors must respect context timeouts, degrade gracefully when optional credentials (`ACCESS_TOKEN`, `KAGGLE_KEY`, or `HF_TOKEN`) are absent, and include table-driven unit tests using local HTTP test servers rather than live network calls.

## Pull request checklist

1. Open an issue first if you plan a large architectural change, or reference an existing issue in your pull request description.
2. Run `go test -v ./...` and `go run ./cmd/dataprofile generate --demo --no-env --output=showcase && go run ./cmd/dataprofile validate --no-env --output=showcase` locally and confirm all pass with zero warnings.
3. Keep prose in documentation direct and factual: state the technical reason and the concrete number, avoid marketing buzzwords, and follow the domain terms in [`GLOSSARY.md`](GLOSSARY.md).
4. By submitting a contribution, you agree that your work will be licensed under the project's [Apache License 2.0](LICENSE) and retain the project attribution in [`NOTICE`](NOTICE).
