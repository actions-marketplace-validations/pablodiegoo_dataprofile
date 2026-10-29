package config

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/pablodiegoo/dataprofile/internal/domain"
)

// Secrets holds authentication tokens loaded strictly from environment variables or a gitignored .env file.
type Secrets struct {
	GitHubToken    string
	KaggleUsername string
	KaggleKey      string
	HFToken        string
}

// SyncTarget represents a GitHub Actions Secret or Repository Variable synced by `dataprofile sync`.
type SyncTarget struct {
	Name     string
	Value    string
	IsSecret bool
}

// LoadDotenvMap parses a .env file into a key-value map.
func LoadDotenvMap(dotenvPath string) map[string]string {
	envMap := make(map[string]string)
	if dotenvPath == "" {
		return envMap
	}
	f, err := os.Open(dotenvPath)
	if err != nil {
		return envMap
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(stripEnvComment(scanner.Text()))
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := cleanYAMLScalar(parts[1])
			envMap[k] = v
		}
	}
	return envMap
}

// LoadSecrets reads optional .env file if present and resolves environment variables.
// Precedence: explicit environment variables > .env file values.
func LoadSecrets(dotenvPath string) Secrets {
	envMap := LoadDotenvMap(dotenvPath)

	get := func(key string) string {
		if v := strings.TrimSpace(os.Getenv(key)); v != "" {
			return v
		}
		return envMap[key]
	}

	ghToken := get("ACCESS_TOKEN")
	if ghToken == "" {
		ghToken = get("GITHUB_TOKEN")
	}

	kgKey := get("KAGGLE_KEY")
	if kgKey == "" {
		kgKey = get("KAGGLE_API_TOKEN")
	}

	hfTok := get("HF_TOKEN")
	if hfTok == "" {
		hfTok = get("HUGGINGFACE_TOKEN")
	}

	return Secrets{
		GitHubToken:    ghToken,
		KaggleUsername: get("KAGGLE_USERNAME"),
		KaggleKey:      kgKey,
		HFToken:        hfTok,
	}
}

// ParseEnvList parses a bracketed semicolon-delimited list `[item 1; item 2]`
// from an environment variable or .env file. Literal semicolons can be escaped as `\;`.
func ParseEnvList(raw string) ([]string, bool) {
	s := cleanYAMLScalar(strings.TrimSpace(raw))
	if !strings.HasPrefix(s, "[") || !strings.HasSuffix(s, "]") {
		return nil, false
	}
	inner := strings.TrimSpace(s[1 : len(s)-1])
	if inner == "" {
		return []string{}, true
	}

	var (
		out []string
		cur strings.Builder
	)
	for i := 0; i < len(inner); i++ {
		if inner[i] == '\\' && i+1 < len(inner) && inner[i+1] == ';' {
			cur.WriteByte(';')
			i++
			continue
		}
		if inner[i] == ';' {
			if item := cleanYAMLScalar(cur.String()); item != "" {
				out = append(out, item)
			}
			cur.Reset()
			continue
		}
		cur.WriteByte(inner[i])
	}
	if item := cleanYAMLScalar(cur.String()); item != "" {
		out = append(out, item)
	}
	return out, true
}

// ApplyEnvOverrides layers process environment variables and optional .env file values
// on top of a loaded ProfileManifest. Precedence: process env > .env file > profile.yaml.
func ApplyEnvOverrides(manifest *domain.ProfileManifest, dotenvPath string) {
	if manifest == nil {
		return
	}
	envMap := LoadDotenvMap(dotenvPath)

	get := func(keys ...string) (string, bool) {
		for _, k := range keys {
			if v := strings.TrimSpace(os.Getenv(k)); v != "" {
				return cleanYAMLScalar(v), true
			}
		}
		for _, k := range keys {
			if v := strings.TrimSpace(envMap[k]); v != "" {
				return v, true
			}
		}
		return "", false
	}

	if v, ok := get("SUBJECT_NAME"); ok {
		manifest.Subject.Name = v
	}
	if v, ok := get("SUBJECT_ROLE"); ok {
		manifest.Subject.Role = v
	}
	if v, ok := get("SUBJECT_LOCALE"); ok {
		manifest.Subject.Locale = v
	}
	if v, ok := get("SUBJECT_HEADLINES"); ok {
		if list, isList := ParseEnvList(v); isList {
			manifest.Subject.Headlines = list
		} else {
			manifest.Subject.Headlines = splitMultilineScalar(v)
		}
	}
	if v, ok := get("SUBJECT_BIO"); ok {
		if list, isList := ParseEnvList(v); isList {
			manifest.Subject.Bio = list
		} else {
			manifest.Subject.Bio = splitMultilineScalar(v)
		}
	}

	if v, ok := get("GITHUB_USERNAME", "GH_USERNAME"); ok {
		manifest.Sources.GitHub.Username = v
	}
	if v, ok := get("GITHUB_INCLUDE_FORKS", "GH_INCLUDE_FORKS"); ok {
		manifest.Sources.GitHub.IncludeForks = parseBool(v, manifest.Sources.GitHub.IncludeForks)
	}
	if v, ok := get("GITHUB_INCLUDE_PRIVATE_COUNT", "GITHUB_INCLUDE_PRIVATE", "GH_INCLUDE_PRIVATE_COUNT", "GH_INCLUDE_PRIVATE"); ok {
		manifest.Sources.GitHub.IncludePrivateCount = parseBool(v, manifest.Sources.GitHub.IncludePrivateCount)
	}
	if v, ok := get("GITHUB_EXCLUDE_REPOS", "GH_EXCLUDE_REPOS"); ok {
		if list, isList := ParseEnvList(v); isList {
			manifest.Sources.GitHub.ExcludeRepos = list
		}
	}
	if v, ok := get("GITHUB_EXCLUDE_LANGS", "GH_EXCLUDE_LANGS"); ok {
		if list, isList := ParseEnvList(v); isList {
			manifest.Sources.GitHub.ExcludeLangs = list
		}
	}

	if v, ok := get("KAGGLE_USERNAME"); ok {
		manifest.Sources.Kaggle.Username = v
	}
	if v, ok := get("KAGGLE_ENABLED"); ok {
		manifest.Sources.Kaggle.Enabled = parseBool(v, manifest.Sources.Kaggle.Enabled)
	}
	if v, ok := get("KAGGLE_MODE"); ok {
		manifest.Sources.Kaggle.Mode = v
	}

	if v, ok := get("HF_USERNAME", "HUGGINGFACE_USERNAME"); ok {
		manifest.Sources.HuggingFace.Username = v
	}
	if v, ok := get("HF_ENABLED", "HUGGINGFACE_ENABLED"); ok {
		manifest.Sources.HuggingFace.Enabled = parseBool(v, manifest.Sources.HuggingFace.Enabled)
	}

	if v, ok := get("DEFLATE_NOTEBOOKS"); ok {
		manifest.Normalization.DeflateNotebooks = parseBool(v, manifest.Normalization.DeflateNotebooks)
	}
	if v, ok := get("DETECT_DATA_STACK"); ok {
		manifest.Normalization.DetectDataStack = parseBool(v, manifest.Normalization.DetectDataStack)
	}

	if v, ok := get("OUTPUT_DIRECTORY"); ok {
		manifest.Output.Directory = v
	}
	if v, ok := get("OUTPUT_SCHOOLS"); ok {
		if list, isList := ParseEnvList(v); isList && len(list) > 0 {
			manifest.Output.Schools = nil
			for _, it := range list {
				manifest.Output.Schools = append(manifest.Output.Schools, domain.School(it))
			}
		}
	}
	if v, ok := get("OUTPUT_PANELS"); ok {
		if list, isList := ParseEnvList(v); isList && len(list) > 0 {
			manifest.Output.Panels = nil
			for _, it := range list {
				manifest.Output.Panels = append(manifest.Output.Panels, domain.Panel(it))
			}
		}
	}
}

// PlanGitHubSync converts a parsed .env map into sorted GitHub Actions Secret and Variable targets.
// Note: GitHub Actions forbids user-defined variables or secrets prefixed with `GITHUB_`, so
// `GITHUB_USERNAME` etc. are automatically mapped to `GH_USERNAME` etc. for GitHub Actions Variables,
// and `GITHUB_TOKEN` is omitted (as GitHub Actions provides `secrets.GITHUB_TOKEN` natively).
func PlanGitHubSync(envMap map[string]string) []SyncTarget {
	secretKeys := map[string]string{
		"ACCESS_TOKEN":      "ACCESS_TOKEN",
		"KAGGLE_KEY":        "KAGGLE_KEY",
		"KAGGLE_API_TOKEN":  "KAGGLE_KEY",
		"HF_TOKEN":          "HF_TOKEN",
		"HUGGINGFACE_TOKEN": "HF_TOKEN",
	}

	varKeys := map[string]string{
		"SUBJECT_NAME":                 "SUBJECT_NAME",
		"SUBJECT_ROLE":                 "SUBJECT_ROLE",
		"SUBJECT_LOCALE":               "SUBJECT_LOCALE",
		"SUBJECT_HEADLINES":            "SUBJECT_HEADLINES",
		"SUBJECT_BIO":                  "SUBJECT_BIO",
		"GITHUB_USERNAME":              "GH_USERNAME",
		"GH_USERNAME":                  "GH_USERNAME",
		"GITHUB_INCLUDE_FORKS":         "GH_INCLUDE_FORKS",
		"GH_INCLUDE_FORKS":             "GH_INCLUDE_FORKS",
		"GITHUB_INCLUDE_PRIVATE_COUNT": "GH_INCLUDE_PRIVATE_COUNT",
		"GITHUB_INCLUDE_PRIVATE":       "GH_INCLUDE_PRIVATE_COUNT",
		"GH_INCLUDE_PRIVATE_COUNT":     "GH_INCLUDE_PRIVATE_COUNT",
		"GH_INCLUDE_PRIVATE":           "GH_INCLUDE_PRIVATE_COUNT",
		"GITHUB_EXCLUDE_REPOS":         "GH_EXCLUDE_REPOS",
		"GH_EXCLUDE_REPOS":             "GH_EXCLUDE_REPOS",
		"GITHUB_EXCLUDE_LANGS":         "GH_EXCLUDE_LANGS",
		"GH_EXCLUDE_LANGS":             "GH_EXCLUDE_LANGS",
		"KAGGLE_USERNAME":              "KAGGLE_USERNAME",
		"KAGGLE_ENABLED":               "KAGGLE_ENABLED",
		"KAGGLE_MODE":                  "KAGGLE_MODE",
		"HF_USERNAME":                  "HF_USERNAME",
		"HUGGINGFACE_USERNAME":         "HF_USERNAME",
		"HF_ENABLED":                   "HF_ENABLED",
		"HUGGINGFACE_ENABLED":          "HF_ENABLED",
		"DEFLATE_NOTEBOOKS":            "DEFLATE_NOTEBOOKS",
		"DETECT_DATA_STACK":            "DETECT_DATA_STACK",
		"OUTPUT_DIRECTORY":             "OUTPUT_DIRECTORY",
		"OUTPUT_SCHOOLS":               "OUTPUT_SCHOOLS",
		"OUTPUT_PANELS":                "OUTPUT_PANELS",
	}

	resolved := make(map[string]SyncTarget)
	for rawKey, rawVal := range envMap {
		val := strings.TrimSpace(rawVal)
		if val == "" {
			continue
		}
		if targetName, ok := secretKeys[rawKey]; ok {
			resolved["secret:"+targetName] = SyncTarget{
				Name:     targetName,
				Value:    val,
				IsSecret: true,
			}
			continue
		}
		if targetName, ok := varKeys[rawKey]; ok {
			resolved["var:"+targetName] = SyncTarget{
				Name:     targetName,
				Value:    val,
				IsSecret: false,
			}
		}
	}

	out := make([]SyncTarget, 0, len(resolved))
	for _, t := range resolved {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].IsSecret != out[j].IsSecret {
			return !out[i].IsSecret // variables first, then secrets
		}
		return out[i].Name < out[j].Name
	})
	return out
}

// DefaultManifest returns sensible defaults when fields are omitted from profile.yaml.
func DefaultManifest() domain.ProfileManifest {
	return domain.ProfileManifest{
		Subject: domain.SubjectConfig{
			Name:   "Data Practitioner",
			Role:   "Data Scientist & Analytics Engineer",
			Locale: "en",
		},
		Sources: domain.SourcesConfig{
			GitHub: domain.GitHubSourceConfig{
				IncludeForks:        false,
				IncludePrivateCount: true,
				ExcludeLangs: []string{
					"HTML",
					"CSS",
					"SCSS",
					"TeX",
					"G-code",
					"Roff",
					"Liquid",
					"Blade",
					"Rich Text Format",
					"Batchfile",
				},
			},
			Kaggle: domain.KaggleSourceConfig{
				Enabled: true,
				Mode:    "hybrid",
			},
			HuggingFace: domain.HuggingFaceSourceConfig{
				Enabled: true,
			},
		},
		Normalization: domain.NormalizationConfig{
			DeflateNotebooks: true,
			DetectDataStack:  true,
		},
		Output: domain.OutputConfig{
			Directory: "generated",
			Schools:   domain.AllSchools(),
			Panels:    domain.AllPanels(),
		},
	}
}

// LoadManifest parses a declarative profile.yaml manifest file.
func LoadManifest(path string) (domain.ProfileManifest, error) {
	manifest := DefaultManifest()
	data, err := os.ReadFile(path)
	if err != nil {
		return manifest, fmt.Errorf("reading manifest %s: %w", path, err)
	}

	type subjectContinuationMode int
	const (
		contNone subjectContinuationMode = iota
		contBlockScalar
		contMultilineScalar
	)

	var (
		section    string
		subsection string
		listKey    string
		contMode   subjectContinuationMode
		schoolsSet bool
		panelsSet  bool
		langsSet   bool
	)

	appendSubjectLines := func(field string, linesToAppend ...string) {
		if field == "bio" {
			manifest.Subject.Bio = append(manifest.Subject.Bio, linesToAppend...)
		} else if field == "headlines" {
			manifest.Subject.Headlines = append(manifest.Subject.Headlines, linesToAppend...)
		}
	}

	lines := strings.Split(string(data), "\n")
	for _, rawLine := range lines {
		rawIndent := len(rawLine) - len(strings.TrimLeft(rawLine, " "))

		// Handle YAML block scalar lines (| or >) or physical multiline scalar continuation lines
		// under subject.bio or subject.headlines before comment stripping,
		// since '#' inside a YAML block scalar or multiline string is literal text (e.g., "Kaggle Master (#412)").
		if section == "subject" && rawIndent >= 4 && contMode != contNone {
			trimmedRaw := strings.TrimSpace(rawLine)
			if contMode == contBlockScalar {
				if textLine := cleanYAMLScalar(trimmedRaw); textLine != "" {
					appendSubjectLines(listKey, textLine)
				}
			} else {
				if textLine := cleanPartialQuotedScalar(trimmedRaw); textLine != "" {
					appendSubjectLines(listKey, splitMultilineScalar(textLine)...)
				}
			}
			continue
		}

		// Strip inline comments outside quotes
		line := stripYAMLComment(rawLine)
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		indent := len(line) - len(strings.TrimLeft(line, " "))

		// Handle list items (- "value")
		if strings.HasPrefix(trimmed, "-") {
			val := cleanYAMLScalar(strings.TrimSpace(strings.TrimPrefix(trimmed, "-")))
			if val == "" {
				continue
			}
			switch {
			case section == "subject" && listKey == "headlines":
				manifest.Subject.Headlines = append(manifest.Subject.Headlines, val)
			case section == "subject" && listKey == "bio":
				manifest.Subject.Bio = append(manifest.Subject.Bio, val)
			case section == "sources" && subsection == "github" && listKey == "exclude_repos":
				manifest.Sources.GitHub.ExcludeRepos = append(manifest.Sources.GitHub.ExcludeRepos, val)
			case section == "sources" && subsection == "github" && listKey == "exclude_langs":
				if !langsSet {
					manifest.Sources.GitHub.ExcludeLangs = nil
					langsSet = true
				}
				manifest.Sources.GitHub.ExcludeLangs = append(manifest.Sources.GitHub.ExcludeLangs, val)
			case section == "output" && listKey == "schools":
				if !schoolsSet {
					manifest.Output.Schools = nil
					schoolsSet = true
				}
				manifest.Output.Schools = append(manifest.Output.Schools, domain.School(val))
			case section == "output" && listKey == "panels":
				if !panelsSet {
					manifest.Output.Panels = nil
					panelsSet = true
				}
				manifest.Output.Panels = append(manifest.Output.Panels, domain.Panel(val))
			}
			continue
		}

		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		rawVal := strings.TrimSpace(parts[1])
		val := cleanYAMLScalar(rawVal)

		if indent == 0 {
			section = key
			subsection = ""
			listKey = ""
			contMode = contNone
			continue
		}

		if indent == 2 {
			contMode = contNone
			if items, ok := parseInlineYAMLArray(rawVal); ok {
				listKey = ""
				switch {
				case section == "output" && key == "schools":
					manifest.Output.Schools = nil
					schoolsSet = true
					for _, it := range items {
						manifest.Output.Schools = append(manifest.Output.Schools, domain.School(it))
					}
				case section == "output" && key == "panels":
					manifest.Output.Panels = nil
					panelsSet = true
					for _, it := range items {
						manifest.Output.Panels = append(manifest.Output.Panels, domain.Panel(it))
					}
				case section == "subject" && key == "headlines":
					manifest.Subject.Headlines = items
				case section == "subject" && key == "bio":
					manifest.Subject.Bio = items
				}
				continue
			}
			if isYAMLBlockScalarIndicator(val) && section == "subject" && (key == "bio" || key == "headlines") {
				listKey = key
				contMode = contBlockScalar
				continue
			}
			if val == "" {
				if section == "sources" {
					subsection = key
					listKey = ""
				} else {
					listKey = key
				}
				continue
			}
			listKey = ""
			switch section {
			case "subject":
				switch key {
				case "name":
					manifest.Subject.Name = val
				case "role":
					manifest.Subject.Role = val
				case "locale":
					manifest.Subject.Locale = val
				case "headlines", "bio":
					listKey = key
					contMode = contMultilineScalar
					appendSubjectLines(key, splitMultilineScalar(cleanPartialQuotedScalar(val))...)
				}
			case "normalization":
				switch key {
				case "deflate_notebooks":
					manifest.Normalization.DeflateNotebooks = parseBool(val, true)
				case "detect_data_stack":
					manifest.Normalization.DetectDataStack = parseBool(val, true)
				}
			case "output":
				if key == "directory" {
					manifest.Output.Directory = val
				}
			}
			continue
		}

		if indent >= 4 && section == "sources" {
			if items, ok := parseInlineYAMLArray(rawVal); ok && subsection == "github" {
				listKey = ""
				switch key {
				case "exclude_repos":
					manifest.Sources.GitHub.ExcludeRepos = items
				case "exclude_langs":
					manifest.Sources.GitHub.ExcludeLangs = items
					langsSet = true
				}
				continue
			}
			if val == "" {
				listKey = key
				continue
			}
			switch subsection {
			case "github":
				switch key {
				case "username":
					manifest.Sources.GitHub.Username = val
				case "include_forks":
					manifest.Sources.GitHub.IncludeForks = parseBool(val, false)
				case "include_private_count", "include_private":
					manifest.Sources.GitHub.IncludePrivateCount = parseBool(val, true)
				}
			case "kaggle":
				switch key {
				case "username":
					manifest.Sources.Kaggle.Username = val
				case "enabled":
					manifest.Sources.Kaggle.Enabled = parseBool(val, true)
				case "mode":
					manifest.Sources.Kaggle.Mode = val
				}
			case "huggingface":
				switch key {
				case "username":
					manifest.Sources.HuggingFace.Username = val
				case "enabled":
					manifest.Sources.HuggingFace.Enabled = parseBool(val, true)
				}
			}
		}
	}

	return manifest, nil
}

func updateQuoteState(ch byte, inSingle, inDouble *bool) {
	switch ch {
	case '\'':
		if !*inDouble {
			*inSingle = !*inSingle
		}
	case '"':
		if !*inSingle {
			*inDouble = !*inDouble
		}
	}
}

func stripYAMLComment(s string) string {
	var inSingle, inDouble bool
	for i := 0; i < len(s); i++ {
		updateQuoteState(s[i], &inSingle, &inDouble)
		if s[i] == '#' && !inSingle && !inDouble {
			return s[:i]
		}
	}
	return s
}

func stripEnvComment(s string) string {
	var (
		inSingle, inDouble bool
		bracketDepth       int
	)
	for i := 0; i < len(s); i++ {
		if bracketDepth == 0 {
			updateQuoteState(s[i], &inSingle, &inDouble)
		}
		if !inSingle && !inDouble {
			switch s[i] {
			case '[':
				bracketDepth++
			case ']':
				if bracketDepth > 0 {
					bracketDepth--
				}
			case '#':
				if bracketDepth == 0 && (i == 0 || s[i-1] == ' ' || s[i-1] == '\t') {
					return s[:i]
				}
			}
		}
	}
	return s
}

func cleanYAMLScalar(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func cleanPartialQuotedScalar(s string) string {
	s = cleanYAMLScalar(s)
	if len(s) >= 1 && (s[0] == '"' || s[0] == '\'') && s[len(s)-1] != s[0] {
		s = strings.TrimSpace(s[1:])
	} else if len(s) >= 1 && (s[len(s)-1] == '"' || s[len(s)-1] == '\'') && s[0] != s[len(s)-1] {
		s = strings.TrimSpace(s[:len(s)-1])
	}
	return s
}

func parseInlineYAMLArray(s string) ([]string, bool) {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "[") || !strings.HasSuffix(s, "]") {
		return nil, false
	}
	inner := strings.TrimSpace(s[1 : len(s)-1])
	if inner == "" {
		return nil, true
	}
	var (
		out                []string
		inSingle, inDouble bool
		start              int
	)
	for i := 0; i < len(inner); i++ {
		updateQuoteState(inner[i], &inSingle, &inDouble)
		if inner[i] == ',' && !inSingle && !inDouble {
			if v := cleanYAMLScalar(inner[start:i]); v != "" {
				out = append(out, v)
			}
			start = i + 1
		}
	}
	if v := cleanYAMLScalar(inner[start:]); v != "" {
		out = append(out, v)
	}
	return out, true
}

func isYAMLBlockScalarIndicator(s string) bool {
	switch strings.TrimSpace(s) {
	case "|", "|-", "|+", ">", ">-", ">+":
		return true
	default:
		return false
	}
}

func splitMultilineScalar(val string) []string {
	var out []string
	for _, p := range strings.Split(strings.ReplaceAll(val, `\n`, "\n"), "\n") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func parseBool(s string, fallback bool) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true", "yes", "1":
		return true
	case "false", "no", "0":
		return false
	default:
		return fallback
	}
}
