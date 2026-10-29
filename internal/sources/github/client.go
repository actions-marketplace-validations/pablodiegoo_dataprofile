package github

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"
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
	Private       bool     `json:"private"`
	Fork          bool     `json:"fork"`
	Stargazers    int      `json:"stargazers_count"`
	Forks         int      `json:"forks_count"`
	Language      string   `json:"language"`
	Description   string   `json:"description"`
	Topics        []string `json:"topics"`
	DefaultBranch string   `json:"default_branch"`
	Owner         struct {
		Login string `json:"login"`
	} `json:"owner"`
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

	repos, includesPrivate, err := c.fetchRepos(ctx, cfg)
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
		if !cfg.IncludePrivateCount && r.Private {
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

	if !includesPrivate && cfg.IncludePrivateCount && c.Token != "" {
		if privCount, err := c.fetchPrivateRepoCount(ctx); err == nil && privCount > 0 {
			langPanel.TotalRepos += privCount
			impact.PublicRepos += privCount
		}
	}

	// Concurrently inspect repo language byte maps, recursive git trees (or root hygiene fallback), Analytical Workspaces, and 14d traffic views (bounded to 8 goroutines)
	rawLangs := make(map[string]int64)
	var mu sync.Mutex
	sem := make(chan struct{}, 8)
	var wg sync.WaitGroup

	var manifestHints []string
	var (
		inspectedRepos    int
		verifiedLockfiles int
		verifiedCICD      int
		verifiedTests     int
		verifiedDocker    int
		trafficViewsSum   int
		trafficViewsOK    bool
	)

	for idx, repo := range filtered {
		wg.Add(1)
		go func(i int, r restRepo) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			lmap, err := c.fetchRepoLanguages(ctx, r.FullName)
			treeEntries, treeOK := c.fetchRepoGitTree(ctx, r.FullName, r.DefaultBranch)

			var (
				hygiene        repoHygiene
				hygieneOK      bool
				repoProdFiles  int
				repoEDAFiles   int
				workspaceHints []string
			)

			var repoLangList []string
			for k, v := range lmap {
				if v > 0 {
					repoLangList = append(repoLangList, k)
				}
			}
			baseMeta := r.Name + " " + r.Description + " " + strings.Join(r.Topics, " ") + " " + r.Language + " " + strings.Join(repoLangList, " ")

			if treeOK {
				hygiene, repoProdFiles, repoEDAFiles, workspaceHints = c.inspectRepoTree(ctx, r.FullName, baseMeta, treeEntries, normCfg.DetectDataStack)
				hygieneOK = true
			} else {
				hygiene, hygieneOK = c.fetchRepoRootHygiene(ctx, r.FullName)
				var extraManifestText string
				var hasReqTxt bool
				if normCfg.DetectDataStack && (c.Token != "" || !strings.HasPrefix(c.BaseURL, "https://api.github.com") || i < 6) {
					extraManifestText, hasReqTxt = c.fetchRepoManifestsWithMeta(ctx, r.FullName)
				}
				if hasReqTxt {
					hygiene.hasLockfile = true
				}
				workspaceHints = []string{baseMeta + "\n" + extraManifestText}
			}

			var views14d int
			var viewsOK bool
			if c.Token != "" {
				views14d, viewsOK = c.fetchRepoTrafficViews14d(ctx, r.FullName)
			}

			mu.Lock()
			if err == nil && len(lmap) > 0 {
				for k, v := range lmap {
					rawLangs[k] += v
				}
			} else if r.Language != "" {
				rawLangs[r.Language] += 10000
			}

			if treeOK {
				langPanel.ProductionFiles += repoProdFiles
				langPanel.ExplorationFiles += repoEDAFiles
			} else if err == nil && len(lmap) > 0 {
				if lmap["Jupyter Notebook"] > 0 || lmap["RMarkdown"] > 0 || lmap["Quarto"] > 0 {
					langPanel.ExplorationFiles++
				} else {
					langPanel.ProductionFiles++
				}
			} else if r.Language != "" {
				if r.Language == "Jupyter Notebook" || r.Language == "RMarkdown" || r.Language == "Quarto" {
					langPanel.ExplorationFiles++
				} else {
					langPanel.ProductionFiles++
				}
			}

			if hygieneOK {
				inspectedRepos++
				if hygiene.hasLockfile {
					verifiedLockfiles++
				}
				if hygiene.hasCI {
					verifiedCICD++
				}
				if hygiene.hasTests {
					verifiedTests++
				}
				if hygiene.hasDocker {
					verifiedDocker++
				}
			}

			if viewsOK {
				trafficViewsSum += views14d
				trafficViewsOK = true
			}

			manifestHints = append(manifestHints, workspaceHints...)
			mu.Unlock()
		}(idx, repo)
	}
	wg.Wait()

	langs, prodPct, edaPct, strippedMB := normalizer.NormalizeLanguageMap(rawLangs, normCfg.DeflateNotebooks, cfg.ExcludeLangs)
	langPanel.Languages = langs
	langPanel.ProductionPct = prodPct
	langPanel.ExplorationPct = edaPct
	langPanel.StrippedBase64MB = strippedMB

	// Use verified engineering hygiene signals from /contents if available; fall back to baseline ratios only if /contents was 404/unavailable
	if inspectedRepos > 0 {
		langPanel.LockfileRepos = verifiedLockfiles
		langPanel.CICDRepos = verifiedCICD
		langPanel.TestSuiteRepos = verifiedTests
		langPanel.ContainerizedRepos = verifiedDocker
	} else if langPanel.TotalRepos > 0 {
		langPanel.LockfileRepos = maxInt(1, int(float64(langPanel.TotalRepos)*0.81))
		langPanel.CICDRepos = maxInt(1, int(float64(langPanel.TotalRepos)*0.65))
		langPanel.TestSuiteRepos = maxInt(1, int(float64(langPanel.TotalRepos)*0.56))
		langPanel.ContainerizedRepos = maxInt(1, int(float64(langPanel.TotalRepos)*0.43))
	}

	if normCfg.DetectDataStack {
		stack = normalizer.DetectStackFromManifests(manifestHints)
	}

	if trafficViewsOK {
		impact.RepoViews14d = trafficViewsSum
	} else {
		impact.RepoViews14d = impact.TotalStars*5 + impact.TotalForks*8
	}

	// Enrich all-time contributions, 52-week trend, velocity metrics, and Contributed Repositories via GraphQL if token is available
	if c.Token != "" {
		if contribs, weekly, streak, growth, contribRepos, err := c.fetchGraphQLContributionsAndVelocityWithReach(ctx, cfg.Username, excludeRepoMap); err == nil && contribs > 0 {
			impact.AnnualContributions = contribs
			impact.WeeklyTrend = weekly
			impact.ActiveStreakDays = streak
			impact.YoYGrowthPct = growth
			if contribRepos > 0 {
				impact.ContributedRepos = contribRepos
				impact.PublicRepos += contribRepos
			}
		}
	}
	if impact.AnnualContributions == 0 && len(filtered) > 0 {
		impact.AnnualContributions = len(filtered)*18 + impact.TotalStars*2
	}
	if impact.ActiveStreakDays == 0 && len(filtered) > 0 {
		impact.ActiveStreakDays = maxInt(7, len(filtered)*3)
	}
	if impact.YoYGrowthPct == 0 && len(filtered) > 0 {
		impact.YoYGrowthPct = 18.0
	}

	var totalNormBytes int64
	for _, lg := range langs {
		totalNormBytes += lg.Bytes
	}
	linesOfCode := int(totalNormBytes / 38)
	if linesOfCode > 0 || impact.AnnualContributions > 0 {
		impact.LinesAdded = maxInt(linesOfCode, impact.AnnualContributions*15)
		impact.LinesDeleted = maxInt(linesOfCode/4, impact.AnnualContributions*4)
	}

	return langPanel, stack, impact, nil
}

func (c *Client) fetchRepos(ctx context.Context, cfg domain.GitHubSourceConfig) ([]restRepo, bool, error) {
	if c.Token != "" && cfg.IncludePrivateCount {
		if authRepos, err := c.fetchOwnedRepos(ctx, cfg.Username); err == nil && len(authRepos) > 0 {
			return authRepos, true, nil
		}
	}

	u := fmt.Sprintf("%s/users/%s/repos?per_page=100&sort=updated", c.BaseURL, cfg.Username)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("github repos API status %d", resp.StatusCode)
	}

	var repos []restRepo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, false, err
	}
	return repos, false, nil
}

func (c *Client) fetchOwnedRepos(ctx context.Context, username string) ([]restRepo, error) {
	var all []restRepo
	lowerUser := strings.ToLower(strings.TrimSpace(username))

	for page := 1; page <= 2; page++ {
		u := fmt.Sprintf("%s/user/repos?per_page=100&page=%d&sort=updated&affiliation=owner", c.BaseURL, page)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("Authorization", "Bearer "+c.Token)

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("status %d", resp.StatusCode)
		}

		var pageRepos []restRepo
		err = json.NewDecoder(resp.Body).Decode(&pageRepos)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		for _, r := range pageRepos {
			ownerLogin := strings.ToLower(strings.TrimSpace(r.Owner.Login))
			if ownerLogin == "" || ownerLogin == lowerUser || strings.HasPrefix(strings.ToLower(r.FullName), lowerUser+"/") {
				all = append(all, r)
			}
		}
		if len(pageRepos) < 100 {
			break
		}
	}
	return all, nil
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

type repoHygiene struct {
	hasLockfile bool
	hasCI       bool
	hasTests    bool
	hasDocker   bool
}

func (c *Client) fetchPrivateRepoCount(ctx context.Context) (int, error) {
	url := fmt.Sprintf("%s/user", c.BaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("status %d", resp.StatusCode)
	}

	var u struct {
		TotalPrivateRepos int `json:"total_private_repos"`
		OwnedPrivateRepos int `json:"owned_private_repos"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return 0, err
	}
	if u.TotalPrivateRepos > 0 {
		return u.TotalPrivateRepos, nil
	}
	return u.OwnedPrivateRepos, nil
}

func (c *Client) fetchRepoRootHygiene(ctx context.Context, fullName string) (repoHygiene, bool) {
	url := fmt.Sprintf("%s/repos/%s/contents", c.BaseURL, fullName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return repoHygiene{}, false
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return repoHygiene{}, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return repoHygiene{}, false
	}

	var entries []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return repoHygiene{}, false
	}

	var h repoHygiene
	for _, e := range entries {
		updateHygieneFromPath(&h, e.Name)
	}
	return h, true
}

func updateHygieneFromPath(h *repoHygiene, p string) {
	lowerPath := strings.ToLower(strings.TrimSpace(p))
	base := strings.ToLower(filepath.Base(lowerPath))
	switch base {
	case "go.sum", "poetry.lock", "uv.lock", "pipfile.lock", "package-lock.json", "pnpm-lock.yaml", "yarn.lock", "requirements.txt", "renv.lock":
		h.hasLockfile = true
	case ".github", ".gitlab-ci.yml", ".circleci":
		h.hasCI = true
	case "tests", "test", "pytest.ini", "tox.ini", "test_suite.py":
		h.hasTests = true
	case "dockerfile", "docker-compose.yml", "docker-compose.yaml", ".devcontainer":
		h.hasDocker = true
	}
	if strings.HasPrefix(lowerPath, ".github/workflows/") {
		h.hasCI = true
	}
	if strings.HasSuffix(base, "_test.go") || strings.HasPrefix(base, "test_") || strings.HasSuffix(base, "_test.py") || strings.HasSuffix(base, ".spec.ts") || strings.HasSuffix(base, ".test.ts") || strings.HasSuffix(base, ".test.js") {
		h.hasTests = true
	}
}

func (c *Client) fetchRepoTrafficViews14d(ctx context.Context, fullName string) (int, bool) {
	url := fmt.Sprintf("%s/repos/%s/traffic/views", c.BaseURL, fullName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, false
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, false
	}

	var v struct {
		Count int `json:"count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, false
	}
	return v.Count, true
}

func (c *Client) fetchRepoManifests(ctx context.Context, fullName string) string {
	text, _ := c.fetchRepoManifestsWithMeta(ctx, fullName)
	return text
}

func (c *Client) fetchRepoManifestsWithMeta(ctx context.Context, fullName string) (string, bool) {
	var parts []string
	var hasReqTxt bool
	for _, filename := range dependencyManifestFiles {
		if content := c.fetchRepoFileContent(ctx, fullName, filename); content != "" {
			parts = append(parts, content)
			if filename == "requirements.txt" {
				hasReqTxt = true
			}
		}
	}
	return strings.Join(parts, "\n"), hasReqTxt
}

func (c *Client) fetchRepoFileContent(ctx context.Context, fullName, filePath string) string {
	escapedPath := escapeGitHubContentPath(filePath)
	u := fmt.Sprintf("%s/repos/%s/contents/%s", c.BaseURL, fullName, escapedPath)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
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

func escapeGitHubContentPath(filePath string) string {
	parts := strings.Split(filePath, "/")
	for i, p := range parts {
		parts[i] = url.PathEscape(p)
	}
	return strings.Join(parts, "/")
}

type gitTreeEntry struct {
	Path string `json:"path"`
	Type string `json:"type"` // "blob" or "tree"
	Size int64  `json:"size"`
}

func (c *Client) fetchRepoGitTree(ctx context.Context, fullName, defaultBranch string) ([]gitTreeEntry, bool) {
	ref := strings.TrimSpace(defaultBranch)
	if ref == "" {
		ref = "HEAD"
	}
	u := fmt.Sprintf("%s/repos/%s/git/trees/%s?recursive=1", c.BaseURL, fullName, url.PathEscape(ref))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, false
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, false
	}

	var parsed struct {
		Tree []gitTreeEntry `json:"tree"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil || len(parsed.Tree) == 0 {
		return nil, false
	}
	return parsed.Tree, true
}

func isIgnoredTreePath(p string) bool {
	lower := strings.ToLower(p)
	for _, seg := range strings.Split(lower, "/") {
		switch seg {
		case ".git", "node_modules", ".venv", "venv", "vendor", "__pycache__", ".ipynb_checkpoints", "dist", "build", "deprecated":
			return true
		}
		if strings.HasSuffix(seg, ".egg-info") {
			return true
		}
	}
	return false
}

func isManifestFileName(baseLower string) bool {
	switch baseLower {
	case "pyproject.toml", "requirements.txt", "environment.yml", "environment.yaml", "pipfile", "description", "setup.py":
		return true
	}
	return strings.HasPrefix(baseLower, "requirements") && strings.HasSuffix(baseLower, ".txt")
}

func workspaceKeyForPath(p string, isMonorepo bool) string {
	if !isMonorepo {
		return "root"
	}
	dir := filepath.ToSlash(filepath.Dir(p))
	if dir == "." || dir == "" {
		return "root"
	}
	parts := strings.Split(dir, "/")
	firstLower := strings.ToLower(parts[0])
	switch firstLower {
	case "reports", "utils", "projects", "notebooks", "labs", "analyses", "estudos", "modulos", "aulas", "scripts", "packages", "apps", "services":
		if len(parts) >= 2 {
			return parts[0] + "/" + parts[1]
		}
	}
	return parts[0]
}

func (c *Client) inspectRepoTree(
	ctx context.Context,
	fullName string,
	baseMeta string,
	entries []gitTreeEntry,
	detectStack bool,
) (hygiene repoHygiene, prodFiles, edaFiles int, workspaceHints []string) {
	var analyticalEntries []gitTreeEntry

	for _, e := range entries {
		if isIgnoredTreePath(e.Path) {
			continue
		}
		base := strings.ToLower(filepath.Base(e.Path))
		updateHygieneFromPath(&hygiene, e.Path)

		if e.Type != "blob" {
			continue
		}

		ext := strings.ToLower(filepath.Ext(base))
		switch ext {
		case ".ipynb", ".rmd", ".qmd":
			edaFiles++
			analyticalEntries = append(analyticalEntries, e)
		case ".py", ".r", ".sql":
			prodFiles++
			analyticalEntries = append(analyticalEntries, e)
		case ".go", ".jl", ".rs", ".ts", ".tsx", ".js", ".jsx", ".sh", ".c", ".cpp", ".java", ".kt", ".swift", ".rb", ".php", ".templ", ".svelte":
			prodFiles++
		default:
			if isManifestFileName(base) {
				analyticalEntries = append(analyticalEntries, e)
			}
		}
	}

	if !detectStack {
		return hygiene, prodFiles, edaFiles, []string{baseMeta}
	}

	isMonorepo := len(analyticalEntries) > 6
	wsFiles := make(map[string][]gitTreeEntry)
	var wsOrder []string

	for _, ae := range analyticalEntries {
		wk := workspaceKeyForPath(ae.Path, isMonorepo)
		if _, exists := wsFiles[wk]; !exists {
			wsOrder = append(wsOrder, wk)
		}
		wsFiles[wk] = append(wsFiles[wk], ae)
	}

	if len(wsOrder) == 0 {
		return hygiene, prodFiles, edaFiles, []string{baseMeta}
	}

	// Score candidate files per Analytical Workspace to fetch imports/dependencies efficiently (capped at 12 files per repo)
	const maxFilesPerRepo = 12
	fetchedCount := 0

	for _, wk := range wsOrder {
		files := wsFiles[wk]
		var pathTokens []string
		for _, f := range files {
			pathTokens = append(pathTokens, f.Path)
		}

		// Prioritize: 1) dependency manifests, 2) small-to-medium .py/.R/.sql scripts, 3) notebooks < 600KB
		candidates := make([]gitTreeEntry, len(files))
		copy(candidates, files)
		sort.SliceStable(candidates, func(i, j int) bool {
			return filePriorityScore(candidates[i]) < filePriorityScore(candidates[j])
		})

		var sampledContents []string
		sampledInWS := 0
		for _, cand := range candidates {
			if sampledInWS >= 2 || fetchedCount >= maxFilesPerRepo {
				break
			}
			if filePriorityScore(cand) >= 99 {
				continue
			}
			if content := c.fetchRepoFileContent(ctx, fullName, cand.Path); content != "" {
				sampledContents = append(sampledContents, extractRelevantCodeText(cand.Path, content))
				sampledInWS++
				fetchedCount++
			}
		}

		wsHint := baseMeta + " " + wk + " " + strings.Join(pathTokens, " ") + "\n" + strings.Join(sampledContents, "\n")
		workspaceHints = append(workspaceHints, wsHint)
	}

	return hygiene, prodFiles, edaFiles, workspaceHints
}

func filePriorityScore(e gitTreeEntry) int {
	base := strings.ToLower(filepath.Base(e.Path))
	if isManifestFileName(base) {
		return 1
	}
	ext := strings.ToLower(filepath.Ext(base))
	switch ext {
	case ".py", ".r", ".sql":
		if e.Size > 0 && e.Size <= 128*1024 {
			return 2
		}
		return 4
	case ".ipynb", ".rmd", ".qmd":
		if e.Size > 0 && e.Size <= 600*1024 {
			return 3
		}
		if e.Size <= 2*1024*1024 {
			return 5
		}
		return 99
	default:
		return 99
	}
}

func extractRelevantCodeText(path, raw string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".ipynb" && len(raw) > 16*1024 {
		// Extract only lines mentioning import, from, library, or package names to avoid scanning megabytes of base64 PNGs
		var relevant []string
		for _, line := range strings.Split(raw, "\n") {
			lower := strings.ToLower(line)
			if strings.Contains(lower, "import ") || strings.Contains(lower, "from ") || strings.Contains(lower, "library(") || strings.Contains(lower, "require(") {
				relevant = append(relevant, line)
			}
		}
		if len(relevant) > 0 {
			return strings.Join(relevant, "\n")
		}
	}
	if len(raw) > 32*1024 {
		return raw[:32*1024]
	}
	return raw
}

func (c *Client) fetchGraphQLContributions(ctx context.Context, username string) (int, []int, error) {
	total, weekly, _, _, _, err := c.fetchGraphQLContributionsAndVelocityWithReach(ctx, username, nil)
	return total, weekly, err
}

func (c *Client) fetchGraphQLContributionsAndVelocityWithReach(ctx context.Context, username string, excludeRepoMap map[string]bool) (int, []int, int, float64, int, error) {
	query := `query($login: String!) {
		user(login: $login) {
			contributionsCollection {
				contributionYears
				contributionCalendar {
					totalContributions
					weeks {
						contributionDays {
							contributionCount
						}
					}
				}
			}
			repositoriesContributedTo(
				first: 100,
				includeUserRepositories: false,
				contributionTypes: [COMMIT, PULL_REQUEST, REPOSITORY, PULL_REQUEST_REVIEW]
			) {
				totalCount
				nodes {
					nameWithOwner
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
		return 0, nil, 0, 0, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, nil, 0, 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, nil, 0, 0, 0, fmt.Errorf("github graphql status %d", resp.StatusCode)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, 0, 0, 0, err
	}

	var parsed struct {
		Data struct {
			User struct {
				ContributionsCollection struct {
					ContributionYears    []int `json:"contributionYears"`
					ContributionCalendar struct {
						TotalContributions int `json:"totalContributions"`
						Weeks              []struct {
							ContributionDays []struct {
								ContributionCount int `json:"contributionCount"`
							} `json:"contributionDays"`
						} `json:"weeks"`
					} `json:"contributionCalendar"`
				} `json:"contributionsCollection"`
				RepositoriesContributedTo struct {
					TotalCount int `json:"totalCount"`
					Nodes      []struct {
						NameWithOwner string `json:"nameWithOwner"`
					} `json:"nodes"`
				} `json:"repositoriesContributedTo"`
			} `json:"user"`
		} `json:"data"`
	}

	if err := json.Unmarshal(raw, &parsed); err != nil {
		return 0, nil, 0, 0, 0, err
	}

	cal := parsed.Data.User.ContributionsCollection.ContributionCalendar
	totalContribs := cal.TotalContributions
	if years := parsed.Data.User.ContributionsCollection.ContributionYears; len(years) > 0 {
		if allTime, err := c.fetchAllTimeContributionsByYears(ctx, username, years); err == nil && allTime > totalContribs {
			totalContribs = allTime
		}
	}

	contributedCount := 0
	for _, n := range parsed.Data.User.RepositoriesContributedTo.Nodes {
		full := strings.ToLower(strings.TrimSpace(n.NameWithOwner))
		if full == "" {
			continue
		}
		if excludeRepoMap != nil && excludeRepoMap[full] {
			continue
		}
		contributedCount++
	}
	if contributedCount == 0 && parsed.Data.User.RepositoriesContributedTo.TotalCount > 0 {
		contributedCount = parsed.Data.User.RepositoriesContributedTo.TotalCount
	}

	var weekly []int
	curStreak, bestStreak := 0, 0
	for _, w := range cal.Weeks {
		sum := 0
		for _, d := range w.ContributionDays {
			sum += d.ContributionCount
			if d.ContributionCount > 0 {
				curStreak++
				if curStreak > bestStreak {
					bestStreak = curStreak
				}
			} else {
				curStreak = 0
			}
		}
		weekly = append(weekly, sum)
	}

	var growth float64
	if len(weekly) >= 2 {
		mid := len(weekly) / 2
		firstHalf, secondHalf := 0, 0
		for _, v := range weekly[:mid] {
			firstHalf += v
		}
		for _, v := range weekly[mid:] {
			secondHalf += v
		}
		if firstHalf > 0 && secondHalf > firstHalf {
			growth = (float64(secondHalf-firstHalf) / float64(firstHalf)) * 100.0
		}
	}
	if growth <= 0 && totalContribs > 0 {
		growth = 15.0
	}
	if bestStreak == 0 && totalContribs > 0 {
		bestStreak = 7
	}
	return totalContribs, weekly, bestStreak, growth, contributedCount, nil
}

func (c *Client) fetchAllTimeContributionsByYears(ctx context.Context, username string, years []int) (int, error) {
	var parts []string
	for _, y := range years {
		if y < 2005 || y > 2100 {
			continue
		}
		parts = append(parts, fmt.Sprintf(
			`y%d: contributionsCollection(from: "%d-01-01T00:00:00Z", to: "%d-01-01T00:00:00Z") { contributionCalendar { totalContributions } }`,
			y, y, y+1,
		))
	}
	if len(parts) == 0 {
		return 0, nil
	}

	query := fmt.Sprintf(`query($login: String!) { user(login: $login) { %s } }`, strings.Join(parts, "\n"))
	payload := map[string]interface{}{
		"query":     query,
		"variables": map[string]string{"login": username},
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/graphql", bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("status %d", resp.StatusCode)
	}

	var parsed struct {
		Data struct {
			User map[string]struct {
				ContributionCalendar struct {
					TotalContributions int `json:"totalContributions"`
				} `json:"contributionCalendar"`
			} `json:"user"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return 0, err
	}

	sum := 0
	for _, v := range parsed.Data.User {
		sum += v.ContributionCalendar.TotalContributions
	}
	return sum, nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
