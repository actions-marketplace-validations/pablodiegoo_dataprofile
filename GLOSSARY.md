# DataProfile

Automated portfolio and statistical visualization generator built for data professionals (data analysts, data scientists, machine learning engineers, and quantitative researchers).

## Language

**Subject**:
The data professional whose public portfolio, code engineering metrics, and community work are portrayed.
_Avoid_: User, Account, Client, Developer

**Source**:
An external platform from which raw professional signals, artifacts, and metrics are collected (GitHub, Kaggle, or Hugging Face Hub).
_Avoid_: Integration, Connector, Plugin, Provider

**Panel**:
A standalone SVG visual module portraying the Subject's portfolio through an analytical dimension (language distribution, data stack, Kaggle standing, community impact), an identity or narrative block (animated header, career bio), or a composite synthesis (executive summary, full portfolio).
_Avoid_: Widget, Card, Badge, Image

**School**:
A data-visualization paradigm (Tufte, Few, Cairo, or Holmes) that governs chart geometry, information density, annotation depth, and data-ink ratio across Panels.
_Avoid_: Theme, Skin, Template, Style

**Palette**:
A chromatic and contrast scheme (such as the dual-surface transparent palette) applied within a School so that text and marks remain legible across light and dark backgrounds.
_Avoid_: Theme, Mode, Skin

**Profile Manifest**:
The declarative specification (`profile.yaml`, optionally overridden by local or CI environment variables) that identifies the Subject, active Sources, and selected Schools without committing secret credentials or personal curriculum churn to upstream source control.
_Avoid_: Config file, Settings

**Showcase**:
The deterministic reference set of Panels rendered from the default Profile Manifest to illustrate the four Schools in project documentation, kept separate from any live Subject portfolio.
_Avoid_: Sample images, Mockup, Preview

**Owned Repository**:
A public or private non-fork repository belonging to the Subject whose source code, notebooks, dependency manifests, and engineering hygiene signals are fully attributed to the Subject's portfolio.
_Avoid_: Public repo, User repo

**Passive Fork**:
A forked repository whose upstream codebase and community stars originate from third parties and are excluded from the Subject's language proportions, stack detection, and star counts.
_Avoid_: Forked repo

**Contributed Repository**:
An external or organization repository where the Subject has authored verified commits, pull requests, or code reviews, counted toward repository reach and contribution velocity without claiming upstream code bytes or stars.
_Avoid_: Ignored repo, Collaborator repo

**Notebook Deflation**:
The normalization process that strips embedded base64 chart images, execution outputs, and JSON schema overhead from interactive notebooks so language proportions reflect only executable source code.
_Avoid_: Bloat filter, File compression

**Analytical Workspace**:
A standalone Owned Repository or a self-contained project subdirectory within a monorepo that houses its own notebooks, source scripts, or dependency manifests, serving as the unit of ecosystem stack detection.
_Avoid_: Subfolder, Directory, Package
