package kaggle

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pablodiegoo/dataprofile/internal/domain"
)

// Client collects Kaggle competitive standing via hybrid Headless Chrome CDP + optional REST API v1.
type Client struct {
	BaseURL     string
	HTTPClient  *http.Client
	DOMDumper   func(ctx context.Context, username string) (string, error)
	APIUsername string
	APIKey      string
}

// NewClient creates a Kaggle hybrid collector.
func NewClient(apiUsername, apiKey string) *Client {
	return &Client{
		BaseURL:     "https://www.kaggle.com",
		HTTPClient:  &http.Client{Timeout: 15 * time.Second},
		DOMDumper:   dumpKaggleDOM,
		APIUsername: strings.TrimSpace(apiUsername),
		APIKey:      strings.TrimSpace(apiKey),
	}
}

// Collect retrieves Kaggle standing for the configured username.
func (c *Client) Collect(ctx context.Context, cfg domain.KaggleSourceConfig) (domain.KaggleStanding, error) {
	if !cfg.Enabled || cfg.Username == "" {
		return domain.KaggleStanding{Available: false}, nil
	}

	mode := strings.ToLower(strings.TrimSpace(cfg.Mode))
	if mode == "" {
		mode = domain.KaggleModeHybrid
	}

	dumper := c.DOMDumper
	if dumper == nil {
		dumper = dumpKaggleDOM
	}

	switch mode {
	case domain.KaggleModeBrowser:
		domText, err := dumper(ctx, cfg.Username)
		if err != nil {
			return domain.KaggleStanding{Available: false, Username: cfg.Username}, err
		}
		return ParseKaggleAchievementsText(domText, cfg.Username), nil

	case domain.KaggleModeAPI:
		if c.APIKey == "" || (!strings.HasPrefix(c.APIKey, "KGAT_") && c.APIUsername == "") {
			return domain.KaggleStanding{Available: false, Username: cfg.Username}, fmt.Errorf("kaggle REST API v1 mode requires KAGGLE_USERNAME and KAGGLE_KEY")
		}
		return c.fetchAPIStanding(ctx, cfg.Username)

	default: // domain.KaggleModeHybrid
		domText, domErr := dumper(ctx, cfg.Username)
		var domStanding domain.KaggleStanding
		if domErr == nil {
			domStanding = ParseKaggleAchievementsText(domText, cfg.Username)
		}

		if c.APIKey != "" && (strings.HasPrefix(c.APIKey, "KGAT_") || c.APIUsername != "") {
			apiStanding, apiErr := c.fetchAPIStanding(ctx, cfg.Username)
			if apiErr == nil {
				if domErr != nil {
					return apiStanding, nil
				}
				return enrichStanding(domStanding, apiStanding), nil
			}
			if domErr != nil {
				return domain.KaggleStanding{Available: false, Username: cfg.Username}, domErr
			}
		}

		if domErr != nil {
			return domain.KaggleStanding{Available: false, Username: cfg.Username}, domErr
		}
		return domStanding, nil
	}
}

type kaggleAPIArtifact struct {
	Ref        string `json:"ref"`
	Title      string `json:"title"`
	Medal      string `json:"medal"`
	Gold       int    `json:"gold"`
	Silver     int    `json:"silver"`
	Bronze     int    `json:"bronze"`
	TotalVotes int    `json:"totalVotes"`
}

func (c *Client) fetchAPIStanding(ctx context.Context, username string) (domain.KaggleStanding, error) {
	kernels, err := c.fetchAPIArtifacts(ctx, "/api/v1/kernels/list", username)
	if err != nil {
		return domain.KaggleStanding{Available: false, Username: username}, err
	}
	datasets, err := c.fetchAPIArtifacts(ctx, "/api/v1/datasets/list", username)
	if err != nil {
		return domain.KaggleStanding{Available: false, Username: username}, err
	}

	// When authenticated as the Subject, also aggregate private kernels and datasets via ?mine=true
	if strings.EqualFold(c.APIUsername, username) || (c.APIUsername == "" && strings.HasPrefix(c.APIKey, "KGAT_")) {
		if privKernels, err := c.fetchMineAPIArtifacts(ctx, "/api/v1/kernels/list"); err == nil {
			kernels = mergeKaggleArtifacts(kernels, privKernels)
		}
		if privDatasets, err := c.fetchMineAPIArtifacts(ctx, "/api/v1/datasets/list"); err == nil {
			datasets = mergeKaggleArtifacts(datasets, privDatasets)
		}
	}

	nbGold, nbSilver, nbBronze := tallyArtifactMedals(kernels)
	dsGold, dsSilver, dsBronze := tallyArtifactMedals(datasets)

	cats := []domain.KaggleCategory{
		buildCategoryFromMedals("Code", nbGold, nbSilver, nbBronze, len(kernels)),
		buildCategoryFromMedals("Datasets", dsGold, dsSilver, dsBronze, len(datasets)),
		buildCategoryFromMedals("Competitions", 0, 0, 0, 0),
	}

	st := domain.KaggleStanding{
		Available:   true,
		Username:    username,
		OverallTier: domain.KaggleTierNovice,
		Categories:  cats,
	}
	bestWeight := 1
	for _, cat := range cats {
		st.TotalGold += cat.Gold
		st.TotalSilver += cat.Silver
		st.TotalBronze += cat.Bronze
		if w := tierHierarchy[cat.Tier]; w > bestWeight {
			bestWeight = w
			st.OverallTier = cat.Tier
		}
	}
	return st, nil
}

func (c *Client) fetchAPIArtifacts(ctx context.Context, endpoint, username string) ([]kaggleAPIArtifact, error) {
	return c.fetchAPIArtifactsWithQuery(ctx, endpoint, "user="+url.QueryEscape(username))
}

func (c *Client) fetchMineAPIArtifacts(ctx context.Context, endpoint string) ([]kaggleAPIArtifact, error) {
	return c.fetchAPIArtifactsWithQuery(ctx, endpoint, "mine=true")
}

func (c *Client) fetchAPIArtifactsWithQuery(ctx context.Context, endpoint, rawQuery string) ([]kaggleAPIArtifact, error) {
	base := strings.TrimRight(c.BaseURL, "/")
	if base == "" {
		base = "https://www.kaggle.com"
	}
	reqURL := fmt.Sprintf("%s%s?%s", base, endpoint, rawQuery)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	if c.APIKey != "" {
		if strings.HasPrefix(c.APIKey, "KGAT_") || c.APIUsername == "" {
			req.Header.Set("Authorization", "Bearer "+c.APIKey)
		} else {
			req.SetBasicAuth(c.APIUsername, c.APIKey)
		}
	}
	req.Header.Set("Accept", "application/json")

	hc := c.HTTPClient
	if hc == nil {
		hc = &http.Client{Timeout: 15 * time.Second}
	}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kaggle API %s?%s returned status %d", endpoint, rawQuery, resp.StatusCode)
	}

	var items []kaggleAPIArtifact
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}

func mergeKaggleArtifacts(base, extra []kaggleAPIArtifact) []kaggleAPIArtifact {
	combined := make([]kaggleAPIArtifact, 0, len(base)+len(extra))
	combined = append(combined, base...)
	combined = append(combined, extra...)

	seen := make(map[string]bool, len(combined))
	out := make([]kaggleAPIArtifact, 0, len(combined))
	for _, it := range combined {
		key := strings.ToLower(strings.TrimSpace(it.Ref))
		if key == "" {
			key = strings.ToLower(strings.TrimSpace(it.Title))
		}
		if key != "" {
			if seen[key] {
				continue
			}
			seen[key] = true
		}
		out = append(out, it)
	}
	return out
}

func tallyArtifactMedals(items []kaggleAPIArtifact) (gold, silver, bronze int) {
	for _, it := range items {
		gold += it.Gold
		silver += it.Silver
		bronze += it.Bronze
		switch strings.ToLower(strings.TrimSpace(it.Medal)) {
		case "gold":
			gold++
		case "silver":
			silver++
		case "bronze":
			bronze++
		default:
			// Kaggle REST API v1 (/kernels/list, /datasets/list) omits the "medal" field and returns raw totalVotes.
			// When explicit medal fields are absent, count artifacts with >= 5 votes conservatively as Bronze.
			if it.Gold == 0 && it.Silver == 0 && it.Bronze == 0 && it.TotalVotes >= 5 {
				bronze++
			}
		}
	}
	return gold, silver, bronze
}

func expertMedalRequirement(trackName string) int {
	switch canonicalTrackName(trackName) {
	case "Code":
		return 5
	case "Datasets":
		return 3
	case "Competitions":
		return 2
	default:
		return 3
	}
}

func canonicalTrackName(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "code", "notebooks", "kernels":
		return "Code"
	case "datasets", "dataset":
		return "Datasets"
	case "competitions", "competition":
		return "Competitions"
	default:
		return name
	}
}

func buildCategoryFromMedals(name string, gold, silver, bronze, artifactCount int) domain.KaggleCategory {
	name = canonicalTrackName(name)
	tier := domain.KaggleTierNovice
	total := gold + silver + bronze
	expertReq := expertMedalRequirement(name)
	switch {
	case gold >= 5 && total >= 15:
		tier = domain.KaggleTierGrandmaster
	case gold >= 1 && total >= 10:
		tier = domain.KaggleTierMaster
	case total >= expertReq || gold >= 1 || silver >= 2:
		tier = domain.KaggleTierExpert
	case artifactCount > 0 || total > 0:
		tier = domain.KaggleTierContributor
	}
	cat := domain.KaggleCategory{
		Name:          name,
		Tier:          tier,
		Gold:          gold,
		Silver:        silver,
		Bronze:        bronze,
		ArtifactCount: artifactCount,
	}
	applyTierProgress(&cat)
	return cat
}

func enrichStanding(base, api domain.KaggleStanding) domain.KaggleStanding {
	out := base
	apiByName := make(map[string]domain.KaggleCategory, len(api.Categories))
	for _, c := range api.Categories {
		apiByName[canonicalTrackName(c.Name)] = c
	}

	out.TotalGold = 0
	out.TotalSilver = 0
	out.TotalBronze = 0
	bestWeight := tierHierarchy[out.OverallTier]
	if bestWeight == 0 {
		bestWeight = 1
	}

	for i := range out.Categories {
		cat := &out.Categories[i]
		cat.Name = canonicalTrackName(cat.Name)
		if ac, ok := apiByName[cat.Name]; ok {
			if ac.ArtifactCount > cat.ArtifactCount {
				cat.ArtifactCount = ac.ArtifactCount
			}
			if ac.Gold > cat.Gold {
				cat.Gold = ac.Gold
			}
			if ac.Silver > cat.Silver {
				cat.Silver = ac.Silver
			}
			if ac.Bronze > cat.Bronze {
				cat.Bronze = ac.Bronze
			}
			if tierHierarchy[ac.Tier] > tierHierarchy[cat.Tier] {
				cat.Tier = ac.Tier
			}
			applyTierProgress(cat)
		}
		out.TotalGold += cat.Gold
		out.TotalSilver += cat.Silver
		out.TotalBronze += cat.Bronze
		if w := tierHierarchy[cat.Tier]; w > bestWeight {
			bestWeight = w
			out.OverallTier = cat.Tier
		}
	}
	return out
}

func dumpKaggleDOM(ctx context.Context, username string) (string, error) {
	browsers := []string{"google-chrome", "google-chrome-stable", "chromium", "chromium-browser"}
	var bin string
	for _, candidate := range browsers {
		if p, err := exec.LookPath(candidate); err == nil {
			bin = p
			break
		}
	}
	if bin == "" {
		return "", fmt.Errorf("headless chrome binary not found in PATH")
	}

	cctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	mainURL := fmt.Sprintf("https://www.kaggle.com/%s", username)
	mainHTML, err := dumpSingleKagglePage(cctx, bin, mainURL)
	if err != nil {
		return "", fmt.Errorf("running headless chrome on %s: %w", mainURL, err)
	}

	cleanedMain := preprocessKaggleHTML(mainHTML)
	var b strings.Builder
	b.WriteString(cleanedMain)

	// Check which active Progression Tracks are not already detailed in "Kaggle Achievements" on the main profile page.
	// For tracks at the Kaggler (pre-Expert) tier, Kaggle omits the summary card on /<username> and renders
	// the track tier and individual medal icons on /<username>/code, /<username>/datasets, and /<username>/competitions.
	achSlice := extractAchievementsSection(cleanedMain)
	subTabs := []struct {
		track   string
		aliases []string
		slug    string
	}{
		{track: "Code", aliases: []string{"Code", "Notebooks"}, slug: "code"},
		{track: "Datasets", aliases: []string{"Datasets"}, slug: "datasets"},
		{track: "Competitions", aliases: []string{"Competitions"}, slug: "competitions"},
	}

	type subResult struct {
		track string
		text  string
	}
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results []subResult
		sem     = make(chan struct{}, 2)
	)

	for _, st := range subTabs {
		if achSlice != "" && findTrackTierWindow(achSlice, st.aliases) != "" {
			continue
		}
		wg.Add(1)
		go func(track, slug string) {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-cctx.Done():
				return
			}
			tabURL := fmt.Sprintf("https://www.kaggle.com/%s/%s", username, slug)
			if tabHTML, tabErr := dumpSingleKagglePage(cctx, bin, tabURL); tabErr == nil && tabHTML != "" {
				cleanedTab := preprocessKaggleHTML(tabHTML)
				mu.Lock()
				results = append(results, subResult{track: track, text: cleanedTab})
				mu.Unlock()
			}
		}(st.track, st.slug)
	}
	wg.Wait()

	for _, r := range results {
		b.WriteString(" [[SUBTAB:")
		b.WriteString(r.track)
		b.WriteString("]] ")
		b.WriteString(r.text)
	}

	return b.String(), nil
}

func dumpSingleKagglePage(ctx context.Context, bin, pageURL string) (string, error) {
	pctx, cancel := context.WithTimeout(ctx, 14*time.Second)
	defer cancel()

	cmd := exec.CommandContext(pctx, bin,
		"--headless=new",
		"--disable-gpu",
		"--no-sandbox",
		"--disable-dev-shm-usage",
		"--virtual-time-budget=5000",
		"--dump-dom",
		pageURL,
	)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

var (
	scriptStyleTagRe   = regexp.MustCompile(`(?is)<(?:script|style)[^>]*>.*?</(?:script|style)>`)
	achievementMedalRe = regexp.MustCompile(`(?is)<img[^>]*?(?:alt|title)="(gold|silver|bronze)\s+medal"[^>]*>\s*<span[^>]*>\s*([0-9][0-9,]*)\s*</span>`)
	itemMedalRe        = regexp.MustCompile(`(?is)<img[^>]*?(?:alt|title)="(gold|silver|bronze)\s+medal"[^>]*>`)
	tierImgRe          = regexp.MustCompile(`(?is)<img[^>]*?(?:alt|title)="(Grandmaster|Master|Expert|Kaggler|Contributor|Novice)"[^>]*>`)
	htmlTagRe          = regexp.MustCompile(`<[^>]+>`)
	multiSpaceRe       = regexp.MustCompile(`\s+`)
)

func preprocessKaggleHTML(raw string) string {
	s := scriptStyleTagRe.ReplaceAllString(raw, " ")
	s = achievementMedalRe.ReplaceAllString(s, " ${2} ${1} ")
	s = itemMedalRe.ReplaceAllString(s, " __ITEM_MEDAL_${1}__ ")
	s = tierImgRe.ReplaceAllString(s, " ${1} ")
	s = htmlTagRe.ReplaceAllString(s, " ")
	return strings.TrimSpace(multiSpaceRe.ReplaceAllString(s, " "))
}

func stripHTMLTags(html string) string {
	return preprocessKaggleHTML(html)
}

var tierHierarchy = map[string]int{
	domain.KaggleTierNovice:      1,
	domain.KaggleTierContributor: 2,
	domain.KaggleTierKaggler:     2,
	domain.KaggleTierExpert:      3,
	domain.KaggleTierMaster:      4,
	domain.KaggleTierGrandmaster: 5,
}

var allTrackBoundaryNames = []string{"Code", "Notebooks", "Datasets", "Competitions", "Discussions", "Discussion", "Models", "Benchmarks", "Followers", "Pinned Work", "Awards"}

// ParseKaggleAchievementsText parses the visible text or HTML of a Kaggle profile page (and optional track sub-tabs)
// into a structured KaggleStanding across the 3 active Kaggle Progression Tracks (Code, Datasets, Competitions).
// Exported for deterministic unit testing against real Kaggle DOM snapshots.
func ParseKaggleAchievementsText(text, username string) domain.KaggleStanding {
	text = preprocessKaggleHTML(text)
	mainText, subTabs := splitSubTabSections(text)

	st := domain.KaggleStanding{
		Available:   true,
		Username:    username,
		OverallTier: domain.KaggleTierContributor,
	}

	tiers := []string{
		domain.KaggleTierGrandmaster,
		domain.KaggleTierMaster,
		domain.KaggleTierExpert,
		domain.KaggleTierKaggler,
		domain.KaggleTierContributor,
		domain.KaggleTierNovice,
	}

	// Detect overall highest tier mentioned in the profile text
	bestWeight := 0
	bestTier := domain.KaggleTierNovice
	for _, t := range tiers {
		if containsTierWord(text, t) {
			w := tierHierarchy[t]
			if w > bestWeight {
				bestWeight = w
				bestTier = t
			}
		}
	}
	st.OverallTier = bestTier

	codeArtifacts := extractTrackArtifactCount(text, codeCountRe)
	dsArtifacts := extractTrackArtifactCount(text, datasetsCountRe)
	compArtifacts := extractTrackArtifactCount(text, competitionsCountRe)

	tracks := []struct {
		name          string
		aliases       []string
		artifactCount int
	}{
		{name: "Code", aliases: []string{"Code", "Notebooks"}, artifactCount: codeArtifacts},
		{name: "Datasets", aliases: []string{"Datasets"}, artifactCount: dsArtifacts},
		{name: "Competitions", aliases: []string{"Competitions"}, artifactCount: compArtifacts},
	}

	achSlice := extractAchievementsSection(mainText)
	searchBase := mainText
	if achSlice != "" {
		searchBase = achSlice
	}

	for _, tr := range tracks {
		cat := parseTrackCategory(searchBase, mainText, subTabs[tr.name], tr.name, tr.aliases, tr.artifactCount, bestTier)
		st.TotalGold += cat.Gold
		st.TotalSilver += cat.Silver
		st.TotalBronze += cat.Bronze
		if cat.Rank > 0 && (st.PeakRank == 0 || cat.Rank < st.PeakRank) {
			st.PeakRank = cat.Rank
		}
		if w := tierHierarchy[cat.Tier]; w > bestWeight {
			bestWeight = w
			st.OverallTier = cat.Tier
		}
		st.Categories = append(st.Categories, cat)
	}

	return st
}

var (
	codeCountRe         = regexp.MustCompile(`(?i)\b(?:Public\s+Notebooks|Code|Notebooks)\s*\(\s*([0-9][0-9,]*)\s*\)`)
	datasetsCountRe     = regexp.MustCompile(`(?i)\b(?:Public\s+Datasets|Datasets)\s*\(\s*([0-9][0-9,]*)\s*\)`)
	competitionsCountRe = regexp.MustCompile(`(?i)\bCompetitions\s*\(\s*([0-9][0-9,]*)\s*\)`)

	rankOfTotalRe = regexp.MustCompile(`(?i)\bRank\s+([0-9][0-9,]*)\s+of\s+([0-9][0-9,]*)`)
	rankHashRe    = regexp.MustCompile(`#\s*([0-9][0-9,]*)`)
	goldRe        = regexp.MustCompile(`(?i)\b([0-9]+)\s*(?:Gold|G\b)`)
	silverRe      = regexp.MustCompile(`(?i)\b([0-9]+)\s*(?:Silver|S\b)`)
	bronzeRe      = regexp.MustCompile(`(?i)\b([0-9]+)\s*(?:Bronze|B\b)`)

	itemGoldMedalRe   = regexp.MustCompile(`(?i)__ITEM_MEDAL_gold__`)
	itemSilverMedalRe = regexp.MustCompile(`(?i)__ITEM_MEDAL_silver__`)
	itemBronzeMedalRe = regexp.MustCompile(`(?i)__ITEM_MEDAL_bronze__`)
)

func splitSubTabSections(text string) (string, map[string]string) {
	subTabs := make(map[string]string)
	marker := "[[SUBTAB:"
	firstIdx := strings.Index(text, marker)
	if firstIdx == -1 {
		return text, subTabs
	}
	mainText := strings.TrimSpace(text[:firstIdx])
	rest := text[firstIdx:]
	for {
		idx := strings.Index(rest, marker)
		if idx == -1 {
			break
		}
		after := rest[idx+len(marker):]
		closeIdx := strings.Index(after, "]]")
		if closeIdx == -1 {
			break
		}
		trackName := canonicalTrackName(after[:closeIdx])
		body := after[closeIdx+2:]
		nextIdx := strings.Index(body, marker)
		if nextIdx == -1 {
			subTabs[trackName] = strings.TrimSpace(body)
			break
		}
		subTabs[trackName] = strings.TrimSpace(body[:nextIdx])
		rest = body[nextIdx:]
	}
	return mainText, subTabs
}

func extractAchievementsSection(text string) string {
	lower := strings.ToLower(text)
	idx := strings.Index(lower, "kaggle achievements")
	if idx == -1 {
		return ""
	}
	slice := text[idx+len("kaggle achievements"):]
	lowerSlice := strings.ToLower(slice)
	for _, stop := range []string{"pinned work", "awards", "public activity", "[[subtab:"} {
		if sIdx := strings.Index(lowerSlice, stop); sIdx != -1 {
			slice = slice[:sIdx]
			lowerSlice = lowerSlice[:sIdx]
		}
	}
	return strings.TrimSpace(slice)
}

func extractTrackArtifactCount(text string, re *regexp.Regexp) int {
	matches := re.FindAllStringSubmatch(text, -1)
	best := 0
	for _, m := range matches {
		if len(m) == 2 {
			cleaned := strings.ReplaceAll(m[1], ",", "")
			if n, err := strconv.Atoi(cleaned); err == nil && n > best {
				best = n
			}
		}
	}
	return best
}

func containsTierWord(text, tier string) bool {
	lower := strings.ToLower(text)
	target := strings.ToLower(tier)
	if target == "grandmaster" {
		// Avoid matching "Legacy Grandmaster" as an active track tier if "Grandmaster" only appears after "Legacy"
		cleaned := strings.ReplaceAll(lower, "legacy grandmaster", "")
		return strings.Contains(cleaned, "grandmaster")
	}
	return strings.Contains(lower, target)
}

func findTrackTierWindow(text string, aliases []string) string {
	lowerText := strings.ToLower(text)
	tiers := []string{
		domain.KaggleTierGrandmaster,
		domain.KaggleTierMaster,
		domain.KaggleTierExpert,
		domain.KaggleTierKaggler,
		domain.KaggleTierContributor,
		domain.KaggleTierNovice,
	}

	for _, alias := range aliases {
		lowerAlias := strings.ToLower(alias)
		searchFrom := 0
		for searchFrom < len(lowerText) {
			relIdx := strings.Index(lowerText[searchFrom:], lowerAlias)
			if relIdx == -1 {
				break
			}
			idx := searchFrom + relIdx
			searchFrom = idx + len(lowerAlias)

			// Require word boundary before alias
			if idx > 0 && isWordChar(lowerText[idx-1]) {
				continue
			}
			endAlias := idx + len(lowerAlias)
			if endAlias < len(lowerText) && isWordChar(lowerText[endAlias]) {
				continue
			}

			window := text[idx:]
			if len(window) > 220 {
				window = window[:220]
			}
			lowerWindow := strings.ToLower(window)

			// Cut window at the next track boundary name
			for _, other := range allTrackBoundaryNames {
				isSelf := false
				for _, a := range aliases {
					if strings.EqualFold(other, a) {
						isSelf = true
						break
					}
				}
				if isSelf {
					continue
				}
				if nextIdx := indexWord(lowerWindow[len(alias):], strings.ToLower(other)); nextIdx != -1 {
					cutoff := len(alias) + nextIdx
					if cutoff < len(window) {
						window = window[:cutoff]
						lowerWindow = lowerWindow[:cutoff]
					}
				}
			}

			// Verify that this window actually describes a track standing (contains a tier name)
			for _, t := range tiers {
				if containsTierWord(window, t) {
					return window
				}
			}
		}
	}
	return ""
}

func isWordChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}

func indexWord(s, word string) int {
	from := 0
	for from < len(s) {
		i := strings.Index(s[from:], word)
		if i == -1 {
			return -1
		}
		pos := from + i
		beforeOK := pos == 0 || !isWordChar(s[pos-1])
		afterPos := pos + len(word)
		afterOK := afterPos >= len(s) || !isWordChar(s[afterPos])
		if beforeOK && afterOK {
			return pos
		}
		from = pos + len(word)
	}
	return -1
}

func parseTrackCategory(searchBase, mainText, subTabText, canonicalName string, aliases []string, artifactCount int, fallbackTier string) domain.KaggleCategory {
	defaultTier := domain.KaggleTierContributor
	if fallbackTier == domain.KaggleTierNovice || fallbackTier == domain.KaggleTierKaggler {
		defaultTier = fallbackTier
	}

	cat := domain.KaggleCategory{
		Name:          canonicalName,
		Tier:          defaultTier,
		ArtifactCount: artifactCount,
	}

	window := findTrackTierWindow(searchBase, aliases)
	if window == "" && searchBase != mainText {
		window = findTrackTierWindow(mainText, aliases)
	}
	if window == "" && subTabText != "" {
		window = findTrackTierWindow(subTabText, aliases)
	}

	tiers := []string{
		domain.KaggleTierGrandmaster,
		domain.KaggleTierMaster,
		domain.KaggleTierExpert,
		domain.KaggleTierKaggler,
		domain.KaggleTierContributor,
		domain.KaggleTierNovice,
	}

	if window != "" {
		for _, t := range tiers {
			if containsTierWord(window, t) {
				cat.Tier = t
				break
			}
		}

		if m := rankOfTotalRe.FindStringSubmatch(window); len(m) == 3 {
			if r, err := strconv.Atoi(strings.ReplaceAll(m[1], ",", "")); err == nil {
				cat.Rank = r
			}
			if tot, err := strconv.Atoi(strings.ReplaceAll(m[2], ",", "")); err == nil {
				cat.TotalUsers = tot
			}
		} else if m := rankHashRe.FindStringSubmatch(window); len(m) == 2 {
			if r, err := strconv.Atoi(strings.ReplaceAll(m[1], ",", "")); err == nil {
				cat.Rank = r
			}
		}

		if m := goldRe.FindStringSubmatch(window); len(m) == 2 {
			cat.Gold, _ = strconv.Atoi(m[1])
		}
		if m := silverRe.FindStringSubmatch(window); len(m) == 2 {
			cat.Silver, _ = strconv.Atoi(m[1])
		}
		if m := bronzeRe.FindStringSubmatch(window); len(m) == 2 {
			cat.Bronze, _ = strconv.Atoi(m[1])
		}
	}

	// Enrich from track sub-tab (/<username>/code, /datasets, /competitions) when medals are rendered per item
	if subTabText != "" {
		if subWin := findTrackTierWindow(subTabText, aliases); subWin != "" {
			for _, t := range tiers {
				if containsTierWord(subWin, t) {
					if tierHierarchy[t] >= tierHierarchy[cat.Tier] || cat.Tier == domain.KaggleTierContributor {
						cat.Tier = t
					}
					break
				}
			}
		}
		if cat.Gold == 0 && cat.Silver == 0 && cat.Bronze == 0 {
			cat.Gold = len(itemGoldMedalRe.FindAllStringIndex(subTabText, -1))
			cat.Silver = len(itemSilverMedalRe.FindAllStringIndex(subTabText, -1))
			cat.Bronze = len(itemBronzeMedalRe.FindAllStringIndex(subTabText, -1))
		}
	}

	applyTierProgress(&cat)
	return cat
}

func applyTierProgress(cat *domain.KaggleCategory) {
	totalMedals := cat.Gold + cat.Silver + cat.Bronze
	expertReq := expertMedalRequirement(cat.Name)

	switch cat.Tier {
	case domain.KaggleTierGrandmaster:
		cat.ProgressPct = 96
		cat.TargetPct = 100
		cat.TopPct = formatPercentileOrDefault(cat.Rank, cat.TotalUsers, "Top 0.1%")
	case domain.KaggleTierMaster:
		cat.ProgressPct = 82
		cat.TargetPct = 90
		cat.TopPct = formatPercentileOrDefault(cat.Rank, cat.TotalUsers, "Top 0.8%")
	case domain.KaggleTierExpert:
		cat.ProgressPct = 62
		cat.TargetPct = 75
		cat.TopPct = formatPercentileOrDefault(cat.Rank, cat.TotalUsers, "Top 1.5%")
	case domain.KaggleTierKaggler, domain.KaggleTierContributor:
		cat.TargetPct = 50
		if totalMedals > 0 && totalMedals < expertReq {
			cat.ProgressPct = math.Round(20.0 + (float64(totalMedals)/float64(expertReq))*28.0)
			if cat.Rank > 0 {
				cat.TopPct = formatPercentileOrDefault(cat.Rank, cat.TotalUsers, "Top 5%")
			} else {
				cat.TopPct = fmt.Sprintf("%d/%d Expert", totalMedals, expertReq)
			}
		} else if cat.Rank > 0 {
			cat.ProgressPct = 32
			cat.TopPct = formatPercentileOrDefault(cat.Rank, cat.TotalUsers, "Top 5%")
		} else {
			if cat.ArtifactCount > 0 {
				cat.ProgressPct = 18
			} else {
				cat.ProgressPct = 12
			}
			cat.TopPct = "Active"
		}
	default:
		cat.ProgressPct = 18
		cat.TargetPct = 35
		cat.TopPct = "Active"
	}
}

func formatPercentileOrDefault(rank, totalUsers int, fallback string) string {
	if rank > 0 && totalUsers > 0 {
		pct := (float64(rank) / float64(totalUsers)) * 100.0
		if pct < 0.1 {
			return "Top 0.1%"
		}
		if pct < 10.0 {
			return fmt.Sprintf("Top %.1f%%", pct)
		}
		return fmt.Sprintf("Top %.0f%%", math.Ceil(pct))
	}
	return fallback
}
