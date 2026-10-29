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
A standalone SVG visual module representing one analytical dimension of the Subject's portfolio (language distribution, data stack, Kaggle standing, community impact, or executive summary).
_Avoid_: Widget, Card, Badge, Image

**School**:
A data-visualization paradigm (Tufte, Few, Cairo, or Holmes) that governs chart geometry, information density, annotation depth, and data-ink ratio across Panels.
_Avoid_: Theme, Skin, Template, Style

**Palette**:
A chromatic and contrast scheme (such as the dual-surface transparent palette) applied within a School so that text and marks remain legible across light and dark backgrounds.
_Avoid_: Theme, Mode, Skin

**Profile Manifest**:
The public, declarative YAML specification (`profile.yaml`) that identifies the Subject, active Sources, and selected Schools without storing secret credentials.
_Avoid_: Config file, Settings, Environment file
