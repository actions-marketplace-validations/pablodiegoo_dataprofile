## Summary

<!-- State what this change does and why in 1-3 direct sentences. Reference the issue number (e.g. Closes #12) if applicable. -->

## Affected Domain Components

- [ ] `School` / `Panel` SVG Renderer (`internal/render`)
- [ ] Jupyter Deflator or Stack Scanner (`internal/normalizer`)
- [ ] Upstream `Source` Collector (`internal/sources/{github,kaggle,huggingface}`)
- [ ] `Profile Manifest` / CLI / GitHub Action (`cmd/dataprofile`, `internal/config`, `action.yml`)
- [ ] Documentation / Governance (`README.md`, `GLOSSARY.md`, `docs/adr/`)

## Verification Checklist

- [ ] `go test -v ./...` passes locally with zero failures.
- [ ] `go run ./cmd/dataprofile generate --demo && go run ./cmd/dataprofile validate` passes (valid XML, dual-surface contrast on `#ffffff` and `#0d1117`, zero school meta-labels).
- [ ] If any analytical SVG panel changed, `prototypes/schools/<school>/<panel>.svg` was updated in sync with `generated/<school>/<panel>.svg`.
- [ ] Domain terms match [`GLOSSARY.md`](GLOSSARY.md) (`Subject`, `Source`, `Panel`, `School`, `Palette`, `Profile Manifest`) and attribution in [`NOTICE`](NOTICE) is preserved.
