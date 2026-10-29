# Contributing to DataProfile

Contributions from data scientists, machine learning engineers, information designers, and Go developers are welcome. Whether you want to add a library to the dependency scanner, refine a chart grammar, or propose a new visualization school, this guide explains how the codebase is organized and how we validate changes before merging.

## Start with the domain language

Before writing code or opening an issue, read [`GLOSSARY.md`](GLOSSARY.md) and the Architectural Decision Records in [`docs/adr/`](docs/adr/). We use six canonical terms across all structs, flags, and pull requests: `Subject` (the data professional portrayed, never "User" or "Client"), `Source` (the upstream platform such as GitHub, Kaggle, or Hugging Face Hub), `Panel` (one of the five SVG modules), `School` (a visualization paradigm such as Tufte, Few, Cairo, or Holmes), `Palette` (the dual-surface color scheme), and `Profile Manifest` (`profile.yaml`). Keeping these terms consistent prevents drift between the documentation and the Go types in [`internal/domain/model.go`](internal/domain/model.go).

## Local development setup

DataProfile compiles as a static Go 1.23+ binary using only the Go standard library. That choice eliminates external package installation steps, so you can clone the repository and run the full test and validation suite immediately:

```bash
# Run all unit and integration tests across collectors, normalizers, and renderers
go test -v ./...

# Render all 20 SVG panels (4 Schools x 5 Panels) using deterministic offline data
go run ./cmd/dataprofile generate --demo

# Render in Brazilian Portuguese
go run ./cmd/dataprofile generate --demo --locale=pt-BR

# Verify XML well-formedness, dual-surface contrast, and absence of school meta-labels
go run ./cmd/dataprofile validate
```

When you modify any of the four analytical panels (`languages.svg`, `stack.svg`, `kaggle.svg`, or `impact.svg`) in [`internal/render/renderer.go`](internal/render/renderer.go), regenerate the demo assets and synchronize the interactive prototype copies in `prototypes/schools/<school>/`. The integration test `TestPrototypesShowcaseAndGeneratedSVGs` in [`cmd/dataprofile/main_test.go`](cmd/dataprofile/main_test.go) checks that both directories match byte-for-byte.

## Common contribution paths

### 1. Adding a data or ML library to the Stack scanner

The `stack.svg` panel inspects `pyproject.toml`, `requirements.txt`, `environment.yml`, `Pipfile`, and `DESCRIPTION` files to detect which tools a `Subject` actually uses in production and research. If a library is missing or needs an alias, edit `canonicalStackRules` in [`internal/normalizer/notebooks.go`](internal/normalizer/notebooks.go) and add a corresponding unit test in [`internal/normalizer/notebooks_test.go`](internal/normalizer/notebooks_test.go).

### 2. Refining or adding a Visualization School (`School`)

Each `School` translates a published theory of information design into parameterized SVG geometry. To add a new `School`, define its constant in [`internal/domain/model.go`](internal/domain/model.go) and implement its five panel renderers (`languages`, `stack`, `kaggle`, `impact`, and `summary`) in [`internal/render/renderer.go`](internal/render/renderer.go), supporting both English (`en`) and Brazilian Portuguese (`pt-BR`).

Every rendered SVG must satisfy four invariants enforced by `ValidateSVG`:

1. Valid, self-contained XML with a transparent root canvas (`fill="none"`).
2. Legible contrast on both GitHub Light (`#ffffff`) and GitHub Dark (`#0d1117`) themes, avoiding washed-out pale greys (`#e6edf3`, `#c9d1d9`, `#8b949e`, `#94a3b8`).
3. Zero visual meta-labels inside the SVG artwork mentioning the school's author (for example, never print "Edward Tufte" or "Stephen Few" inside the chart).
4. Graceful degradation when a metric is zero: if a `Subject` has zero Kaggle gold medals or zero Hugging Face models, the layout pivots to active tracks or engineering velocity rather than displaying `0 Gold` or `0 / 0 / 0`.

### 3. Improving a Source collector

Source collectors live under `internal/sources/{github,kaggle,huggingface}/`. Because GitHub Actions workflows run unattended on a daily schedule, collectors must respect context timeouts, degrade gracefully when optional credentials (`ACCESS_TOKEN` or `KAGGLE_KEY`) are absent, and include table-driven unit tests using local HTTP test servers rather than live network calls.

## Pull request checklist

1. Open an issue first if you plan a large architectural change, or reference an existing issue in your pull request description.
2. Run `go test -v ./...` and `go run ./cmd/dataprofile validate` locally and confirm both pass with zero warnings.
3. Keep prose in documentation direct and factual: state the technical reason and the concrete number, avoid marketing buzzwords, and follow the domain terms in [`GLOSSARY.md`](GLOSSARY.md).
4. By submitting a contribution, you agree that your work will be licensed under the project's [Apache License 2.0](LICENSE) and retain the project attribution in [`NOTICE`](NOTICE).
