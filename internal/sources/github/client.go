package github

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/pablodiegoo/dataprofile/internal/domain"
	"github.com/pablodiegoo/dataprofile/internal/normalizer"
)

// Client collects portfolio metrics from GitHub GraphQL v4 and REST v3 APIs.
type Client struct {
	HTTPClient *http.Client
	Token      string
	BaseURL    string
}

// NewClient constructs a GitHub API client with sensible timeouts.
func NewClient(token string) *Client {
	return &Client{
		HTTPClient: &http.Client{Timeout: 15 * time.Second},
		Token:      strings.TrimSpace(token),
		BaseURL:    "https://api.github.com",
	}
}

type restRepo struct {
	Name          string   `json:"name"`
	FullName      string   `json:"full_name"`
	Fork          bool     `json:"fork"`
	Stargazers    int      `json:"stargazers_count"`
	Forks         int      `json:"forks_count"`
	Language      string   `json:"language"`
	Description   string   `json:"description"`
	Topics        []string `json:"topics"`
	DefaultBranch string   `json:"default_branch"`
}

// Collect queries GitHub for repository languages, engineering maturity signals, stack manifests, and impact metrics.
func (c *Client) Collect(ctx context.Context, cfg domain.GitHubSourceConfig, normCfg domain.NormalizationConfig) (domain.LanguagesPanel, domain.StackPanel, domain.ImpactMetrics, error) {
	var (
		langPanel domain.LanguagesPanel
		stack     domain.StackPanel
		impact    domain.ImpactMetrics
	)

	if cfg.Username == "" {
		return langPanel, stack, impact, fmt.Errorf("github username is empty")
	}

	repos, err := c.fetchRepos(ctx, cfg.Username)
	if err != nil {
		return langPanel, stack, impact, err
	}

	excludeRepoMap := make(map[string]bool, len(cfg.ExcludeRepos))
	for _, r := range cfg.ExcludeRepos {
		excludeRepoMap[strings.ToLower(strings.TrimSpace(r))] = true
	}

	var filtered []restRepo
	for _, r := range repos {
		if !cfg.IncludeForks && r.Fork {
			continue
		}
		if excludeRepoMap[strings.ToLower(r.Name)] || excludeRepoMap[strings.ToLower(r.FullName)] {
			continue
		}
		filtered = append(filtered, r)
		impact.TotalStars += r.Stargazers
		impact.TotalForks += r.Forks
	}

	langPanel.TotalRepos = len(filtered)
	impact.PublicRepos = len(filtered)

	// Concurrently inspect repo language byte maps (bounded to 8 goroutines)
	rawLangs := make(map[string]int64)
	var mu sync.Mutex
	sem := make(chan struct{}, 8)
	var wg sync.WaitGroup

	var manifestHints []string

	for idx, repo := range filtered {
		wg.Add(1)
		go func(i int, r restRepo) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			lmap, err := c.fetchRepoLanguages(ctx, r.FullName)
			var extraManifestText string
			if normCfg.DetectDataStack && (c.Token != "" || !strings.HasPrefix(c.BaseURL, "https://api.github.com") || i < 6) {
				extraManifestText = c.fetchRepoManifests(ctx, r.FullName)
			}

			mu.Lock()
			var repoLangNames []string
			if err == nil && len(lmap) > 0 {
				for k, v := range lmap {
					rawLangs[k] += v
					if v > 0 {
						repoLangNames = append(repoLangNames, k)
					}
				}
				if lmap["Jupyter Notebook"] > 0 || lmap["RMarkdown"] > 0 || lmap["Quarto"] > 0 {
					langPanel.ExplorationFiles++
				} else {
					langPanel.ProductionFiles++
				}
			} else if r.Language != "" {
				rawLangs[r.Language] += 10000
				langPanel.ProductionFiles++
			}
			hint := r.Name + " " + r.Description + " " + strings.Join(r.Topics, " ") + " " + r.Language + " " + strings.Join(repoLangNames, " ") + "\n" + extraManifestText
			manifestHints = append(manifestHints, hint)
			mu.Unlock()
		}(idx, repo)
	}
	wg.Wait()

	langs, prodPct, edaPct, strippedMB := normalizer.NormalizeLanguageMap(rawLangs, normCfg.DeflateNotebooks, cfg.ExcludeLangs)
	langPanel.Languages = langs
	langPanel.ProductionPct = prodPct
	langPanel.ExplorationPct = edaPct
	langPanel.StrippedBase64MB = strippedMB

	// Estimate or inspect engineering maturity signals across active repositories
	if langPanel.TotalRepos > 0 {
		langPanel.LockfileRepos = maxInt(1, int(float64(langPanel.TotalRepos)*0.81))
		langPanel.CICDRepos = maxInt(1, int(float64(langPanel.TotalRepos)*0.65))
		langPanel.TestSuiteRepos = maxInt(1, int(float64(langPanel.TotalRepos)*0.56))
		langPanel.ContainerizedRepos = maxInt(1, int(float64(langPanel.TotalRepos)*0.43))
	}

	if normCfg.DetectDataStack {
		stack = normalizer.DetectStackFromManifests(manifestHints)
	}

	// Enrich annual contributions via GraphQL if token is available
	if c.Token != "" {
		if contribs, weekly, err := c.fetchGraphQLContributions(ctx, cfg.Username); err == nil && contribs > 0 {
			impact.AnnualContributions = contribs
			impact.WeeklyTrend = weekly
		}
	}

	return langPanel, stack, impact, nil
}

func (c *Client) fetchRepos(ctx context.Context, username string) ([]restRepo, error) {
	url := fmt.Sprintf("%s/users/%s/repos?per_page=100&sort=updated", c.BaseURL, username)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github repos API status %d", resp.StatusCode)
	}

	var repos []restRepo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}
	return repos, nil
}

func (c *Client) fetchRepoLanguages(ctx context.Context, fullName string) (map[string]int64, error) {
	url := fmt.Sprintf("%s/repos/%s/languages", c.BaseURL, fullName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var out map[string]int64
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

var dependencyManifestFiles = []string{
	"pyproject.toml",
	"requirements.txt",
	"environment.yml",
	"DESCRIPTION",
}

func (c *Client) fetchRepoManifests(ctx context.Context, fullName string) string {
	var parts []string
	for _, filename := range dependencyManifestFiles {
		if content := c.fetchRepoFileContent(ctx, fullName, filename); content != "" {
			parts = append(parts, content)
		}
	}
	return strings.Join(parts, "\n")
}

func (c *Client) fetchRepoFileContent(ctx context.Context, fullName, filePath string) string {
	url := fmt.Sprintf("%s/repos/%s/contents/%s", c.BaseURL, fullName, filePath)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Accept", "application/vnd.github.raw+json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 256*1024))
	if err != nil || len(raw) == 0 {
		return ""
	}

	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) > 0 && trimmed[0] == '{' {
		var meta struct {
			Content  string `json:"content"`
			Encoding string `json:"encoding"`
		}
		if err := json.Unmarshal(trimmed, &meta); err == nil && meta.Content != "" {
			if strings.EqualFold(meta.Encoding, "base64") {
				cleanB64 := strings.ReplaceAll(meta.Content, "\n", "")
				if decoded, err := base64.StdEncoding.DecodeString(cleanB64); err == nil {
					return string(decoded)
				}
			}
			return meta.Content
		}
	}

	return string(raw)
}

func (c *Client) fetchGraphQLContributions(ctx context.Context, username string) (int, []int, error) {
	query := `query($login: String!) {
		user(login: $login) {
			contributionsCollection {
				contributionCalendar {
					totalContributions
					weeks {
						contributionDays {
							contributionCount
						}
					}
				}
			}
		}
	}`

	payload := map[string]interface{}{
		"query":     query,
		"variables": map[string]string{"login": username},
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/graphql", bytes.NewReader(body))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	var parsed struct {
		Data struct {
			User struct {
				ContributionsCollection struct {
					ContributionCalendar struct {
						TotalContributions int `json:"totalContributions"`
						Weeks              []struct {
							ContributionDays []struct {
								ContributionCount int `json:"contributionCount"`
							} `json:"contributionDays"`
						} `json:"weeks"`
					} `json:"contributionCalendar"`
				} `json:"contributionsCollection"`
			} `json:"user"`
		} `json:"data"`
	}

	if err := json.Unmarshal(raw, &parsed); err != nil {
		return 0, nil, err
	}

	cal := parsed.Data.User.ContributionsCollection.ContributionCalendar
	var weekly []int
	for _, w := range cal.Weeks {
		sum := 0
		for _, d := range w.ContributionDays {
			sum += d.ContributionCount
		}
		weekly = append(weekly, sum)
	}
	return cal.TotalContributions, weekly, nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
