package config

import (
	"bufio"
	"fmt"
	"os"
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

// LoadSecrets reads optional .env file if present and resolves environment variables.
// Precedence: explicit environment variables > .env file values.
func LoadSecrets(dotenvPath string) Secrets {
	envMap := make(map[string]string)
	if dotenvPath != "" {
		if f, err := os.Open(dotenvPath); err == nil {
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := strings.TrimSpace(stripYAMLComment(scanner.Text()))
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
		}
	}

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

	var (
		section    string
		subsection string
		listKey    string
		schoolsSet bool
		panelsSet  bool
		langsSet   bool
	)

	lines := strings.Split(string(data), "\n")
	for _, rawLine := range lines {
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
			continue
		}

		if indent == 2 {
			if items, ok := parseInlineYAMLArray(rawVal); ok && section == "output" {
				listKey = ""
				switch key {
				case "schools":
					manifest.Output.Schools = nil
					schoolsSet = true
					for _, it := range items {
						manifest.Output.Schools = append(manifest.Output.Schools, domain.School(it))
					}
				case "panels":
					manifest.Output.Panels = nil
					panelsSet = true
					for _, it := range items {
						manifest.Output.Panels = append(manifest.Output.Panels, domain.Panel(it))
					}
				}
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

func stripYAMLComment(s string) string {
	inSingle := false
	inDouble := false
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
			}
		case '#':
			if !inSingle && !inDouble {
				return s[:i]
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

func parseInlineYAMLArray(s string) ([]string, bool) {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "[") || !strings.HasSuffix(s, "]") {
		return nil, false
	}
	inner := strings.TrimSpace(s[1 : len(s)-1])
	if inner == "" {
		return nil, true
	}
	parts := strings.Split(inner, ",")
	var out []string
	for _, p := range parts {
		v := cleanYAMLScalar(p)
		if v != "" {
			out = append(out, v)
		}
	}
	return out, true
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
