## Summary

<!-- State what this change does and why in 1-3 direct sentences. Reference the issue number (e.g. Closes #14) if applicable. -->

## Affected Domain Components

- [ ] `School` / `Panel` SVG Renderer (`internal/render`)
- [ ] Jupyter Deflator or Stack Scanner (`internal/normalizer`)
- [ ] Upstream `Source` Collector (`internal/sources/{github,kaggle,huggingface}`)
- [ ] `Profile Manifest` / `.env` Overrides / CLI / GitHub Action (`cmd/dataprofile`, `internal/config`, `action.yml`)
- [ ] Documentation / Governance (`README.md`, `GLOSSARY.md`, `docs/adr/`)

## Verification Checklist

- [ ] `go test -v ./...` passes locally with zero failures.
- [ ] `go run ./cmd/dataprofile generate --demo --no-env --output=showcase && go run ./cmd/dataprofile validate --no-env --output=showcase` passes (valid XML, dual-surface contrast on `#ffffff` and `#0d1117`, zero school meta-labels).
- [ ] If any analytical SVG panel (`languages`, `stack`, `kaggle`, `impact`) changed, `prototypes/schools/<school>/<panel>.svg` was updated in sync with the rendered `Showcase` output.
- [ ] Domain terms match [`GLOSSARY.md`](GLOSSARY.md) (`Subject`, `Source`, `Panel`, `School`, `Palette`, `Profile Manifest`, `Showcase`) and attribution in [`NOTICE`](NOTICE) is preserved.
