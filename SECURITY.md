# Security Policy

## Supported versions

DataProfile releases follow semantic versioning (`v1.x.y`) with a moving `v1` major tag for the GitHub Marketplace Action. Security fixes are applied to the latest `v1` release branch and `main`.

| Version | Supported |
| :--- | :--- |
| `v1.x` (`main`) | Yes |
| `< v1.0.0` | No |

## Credential and secret hygiene

DataProfile is designed so that `profile.yaml` (`Profile Manifest`) on `main` holds only a neutral reference `Showcase` configuration and zero secret credentials. Secret tokens must never be written to `profile.yaml` or committed to git history.

In GitHub Actions workflows, credentials are supplied through encrypted **GitHub Actions Secrets**, while non-secret curriculum overrides (`SUBJECT_*`, `GH_*`, `KAGGLE_ENABLED`, `HF_USERNAME`, `OUTPUT_*`) are stored in **GitHub Actions Variables** (`Settings -> Secrets and variables -> Actions`):

1. `GITHUB_TOKEN`: Provided automatically by GitHub Actions with `permissions: contents: write` scoped to the repository.
2. `ACCESS_TOKEN` (optional Secret): A classic Personal Access Token (`read:user`, `repo`) or fine-grained token used when the `Subject` wants private `Owned Repository` aggregation or 14-day traffic view totals.
3. `KAGGLE_USERNAME` (Variable or Secret) and `KAGGLE_KEY` (optional Secret): Kaggle API credentials (`KGAT_...` Bearer token or classic key) used for REST API v1 enrichment and private artifact counts.
4. `HF_TOKEN` (optional Secret): Hugging Face Hub token used to include private models, datasets, and spaces in aggregate counts.

For local CLI runs and curriculum management, place credentials and `[item 1; item 2]` manifest overrides in a local `.env` file based on [`.env.example`](.env.example). When running `go run ./cmd/dataprofile sync`, `PlanGitHubSync` routes tokens (`ACCESS_TOKEN`, `KAGGLE_KEY`, `HF_TOKEN`) exclusively to `gh secret set`, maps public configuration fields to `gh variable set`, and masks secret values as `***` during `--dry-run`.

The [`.gitignore`](.gitignore) file blocks `.env`, `.env.*`, `kaggle.json`, `/generated/`, and `/showcase/` from being tracked on `main`, and the automated test `TestGitHubActionWorkflowAndSecretHygiene` in [`cmd/dataprofile/main_test.go`](cmd/dataprofile/main_test.go) verifies these exclusion rules on every commit.

## Reporting a vulnerability

If you discover a security issue (such as credential leakage in logs, XML/SVG injection via untrusted upstream API fields, or command injection in the composite GitHub Action), please do not open a public GitHub issue.

Instead, report it privately using [GitHub Private Vulnerability Reporting](https://github.com/pablodiegoo/dataprofile/security/advisories/new) on this repository. Include the affected file or workflow step, a minimal reproduction input, and the potential impact.

You will receive an acknowledgment within 72 hours. Once the fix is verified and released on `main` and the `v1` tag, we will publish the security advisory and credit your finding (unless you prefer to remain anonymous).
