# Security Policy

## Supported versions

DataProfile releases follow semantic versioning (`v1.x.y`) with a moving `v1` major tag for the GitHub Marketplace Action. Security fixes are applied to the latest `v1` release branch and `main`.

| Version | Supported |
| :--- | :--- |
| `v1.x` (`main`) | Yes |
| `< v1.0.0` | No |

## Credential and secret hygiene

DataProfile is designed so that `profile.yaml` (`Profile Manifest`) contains only public usernames and display preferences. Secret tokens must never be written to `profile.yaml` or committed to git history.

In GitHub Actions workflows, credentials should only be supplied through encrypted **GitHub Actions Secrets** (`Settings -> Secrets and variables -> Actions`):

1. `GITHUB_TOKEN`: Provided automatically by GitHub Actions with `permissions: contents: write` scoped to the repository.
2. `ACCESS_TOKEN` (optional): A classic Personal Access Token (`read:user`, `repo`) or fine-grained token used only when the `Subject` wants private repository counts or 14-day traffic view totals.
3. `KAGGLE_USERNAME` and `KAGGLE_KEY` (optional): Kaggle API credentials (`KGAT_...` Bearer token or classic key) used for REST API v1 enrichment.

For local CLI runs, place credentials in a local `.env` file based on [`.env.example`](.env.example). The [`.gitignore`](.gitignore) file blocks `.env`, `.env.*`, and `kaggle.json` from being tracked, and the automated test `TestGitHubActionWorkflowAndSecretHygiene` in [`cmd/dataprofile/main_test.go`](cmd/dataprofile/main_test.go) verifies these exclusion rules on every commit.

## Reporting a vulnerability

If you discover a security issue (such as credential leakage in logs, XML/SVG injection via untrusted upstream API fields, or command injection in the composite GitHub Action), please do not open a public GitHub issue.

Instead, report it privately using [GitHub Private Vulnerability Reporting](https://github.com/pablodiegoo/dataprofile/security/advisories/new) on this repository. Include the affected file or workflow step, a minimal reproduction input, and the potential impact.

You will receive an acknowledgment within 72 hours. Once the fix is verified and released on `main` and the `v1` tag, we will publish the security advisory and credit your finding (unless you prefer to remain anonymous).
