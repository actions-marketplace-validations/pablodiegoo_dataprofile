package kaggle

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
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
		if c.APIUsername == "" || c.APIKey == "" {
			return domain.KaggleStanding{Available: false, Username: cfg.Username}, fmt.Errorf("kaggle REST API v1 mode requires KAGGLE_USERNAME and KAGGLE_KEY")
		}
		return c.fetchAPIStanding(ctx, cfg.Username)

	default: // domain.KaggleModeHybrid
		domText, domErr := dumper(ctx, cfg.Username)
		var domStanding domain.KaggleStanding
		if domErr == nil {
			domStanding = ParseKaggleAchievementsText(domText, cfg.Username)
		}

		if c.APIUsername != "" && c.APIKey != "" {
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

	nbGold, nbSilver, nbBronze := tallyArtifactMedals(kernels)
	dsGold, dsSilver, dsBronze := tallyArtifactMedals(datasets)

	cats := []domain.KaggleCategory{
		buildCategoryFromMedals("Notebooks", nbGold, nbSilver, nbBronze, len(kernels)),
		buildCategoryFromMedals("Datasets", dsGold, dsSilver, dsBronze, len(datasets)),
		buildCategoryFromMedals("Competitions", 0, 0, 0, 0),
		buildCategoryFromMedals("Discussions", 0, 0, 0, 0),
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
	base := strings.TrimRight(c.BaseURL, "/")
	if base == "" {
		base = "https://www.kaggle.com"
	}
	reqURL := fmt.Sprintf("%s%s?user=%s", base, endpoint, url.QueryEscape(username))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	if c.APIUsername != "" && c.APIKey != "" {
		req.SetBasicAuth(c.APIUsername, c.APIKey)
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
		return nil, fmt.Errorf("kaggle API %s returned status %d", endpoint, resp.StatusCode)
	}

	var items []kaggleAPIArtifact
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
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
		}
	}
	return gold, silver, bronze
}

func buildCategoryFromMedals(name string, gold, silver, bronze, artifactCount int) domain.KaggleCategory {
	tier := domain.KaggleTierNovice
	total := gold + silver + bronze
	switch {
	case gold >= 5 && total >= 15:
		tier = domain.KaggleTierGrandmaster
	case gold >= 1 && total >= 10:
		tier = domain.KaggleTierMaster
	case total >= 3 || gold >= 1 || silver >= 2:
		tier = domain.KaggleTierExpert
	case artifactCount > 0 || total > 0:
		tier = domain.KaggleTierContributor
	}
	cat := domain.KaggleCategory{
		Name:   name,
		Tier:   tier,
		Gold:   gold,
		Silver: silver,
		Bronze: bronze,
	}
	applyTierProgress(&cat)
	return cat
}

func enrichStanding(base, api domain.KaggleStanding) domain.KaggleStanding {
	out := base
	apiByName := make(map[string]domain.KaggleCategory, len(api.Categories))
	for _, c := range api.Categories {
		apiByName[c.Name] = c
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
		if ac, ok := apiByName[cat.Name]; ok {
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
				applyTierProgress(cat)
			}
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

	cctx, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()

	url := fmt.Sprintf("https://www.kaggle.com/%s", username)
	cmd := exec.CommandContext(cctx, bin,
		"--headless=new",
		"--disable-gpu",
		"--no-sandbox",
		"--disable-dev-shm-usage",
		"--virtual-time-budget=5000",
		"--dump-dom",
		url,
	)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("running headless chrome on %s: %w", url, err)
	}
	return stripHTMLTags(string(out)), nil
}

var htmlTagRe = regexp.MustCompile(`<[^>]+>`)
var multiSpaceRe = regexp.MustCompile(`\s+`)

func stripHTMLTags(html string) string {
	text := htmlTagRe.ReplaceAllString(html, " ")
	return strings.TrimSpace(multiSpaceRe.ReplaceAllString(text, " "))
}

var tierHierarchy = map[string]int{
	domain.KaggleTierNovice:      1,
	domain.KaggleTierContributor: 2,
	domain.KaggleTierExpert:      3,
	domain.KaggleTierMaster:      4,
	domain.KaggleTierGrandmaster: 5,
}

// ParseKaggleAchievementsText parses the visible text of a Kaggle profile page into a structured KaggleStanding.
// Exported for deterministic unit testing against real Kaggle DOM snapshots.
func ParseKaggleAchievementsText(text, username string) domain.KaggleStanding {
	text = stripHTMLTags(text)
	st := domain.KaggleStanding{
		Available:   true,
		Username:    username,
		OverallTier: domain.KaggleTierContributor,
	}

	canonicalCats := []string{"Notebooks", "Datasets", "Competitions", "Discussions"}
	tiers := []string{
		domain.KaggleTierGrandmaster,
		domain.KaggleTierMaster,
		domain.KaggleTierExpert,
		domain.KaggleTierContributor,
		domain.KaggleTierNovice,
	}

	// Detect overall highest tier mentioned in the profile text
	bestWeight := 1
	bestTier := domain.KaggleTierNovice
	for _, t := range tiers {
		if strings.Contains(text, t) {
			if w := tierHierarchy[t]; w > bestWeight {
				bestWeight = w
				bestTier = t
			}
		}
	}
	st.OverallTier = bestTier

	// Parse per-category blocks if present: e.g. "Notebooks Master #412 2 Gold 6 Silver 9 Bronze"
	for _, catName := range canonicalCats {
		cat := parseCategoryBlock(text, catName, bestTier)
		st.TotalGold += cat.Gold
		st.TotalSilver += cat.Silver
		st.TotalBronze += cat.Bronze
		if cat.Rank > 0 && (st.PeakRank == 0 || cat.Rank < st.PeakRank) {
			st.PeakRank = cat.Rank
		}
		st.Categories = append(st.Categories, cat)
	}

	return st
}

var (
	rankRe   = regexp.MustCompile(`#\s*([0-9][0-9,]*)`)
	goldRe   = regexp.MustCompile(`([0-9]+)\s*(?:Gold|G\b)`)
	silverRe = regexp.MustCompile(`([0-9]+)\s*(?:Silver|S\b)`)
	bronzeRe = regexp.MustCompile(`([0-9]+)\s*(?:Bronze|B\b)`)
)

func parseCategoryBlock(fullText, catName, fallbackTier string) domain.KaggleCategory {
	cat := domain.KaggleCategory{
		Name:        catName,
		Tier:        domain.KaggleTierContributor,
		ProgressPct: 35,
		TargetPct:   55,
	}

	idx := strings.Index(strings.ToLower(fullText), strings.ToLower(catName))
	if idx == -1 {
		return cat
	}

	window := fullText[idx:]
	if len(window) > 180 {
		window = window[:180]
	}
	lowerWindow := strings.ToLower(window)
	for _, otherCat := range []string{"Notebooks", "Datasets", "Competitions", "Discussions"} {
		if strings.EqualFold(otherCat, catName) {
			continue
		}
		if nextIdx := strings.Index(lowerWindow[len(catName):], strings.ToLower(otherCat)); nextIdx != -1 {
			cutoff := len(catName) + nextIdx
			if cutoff < len(window) {
				window = window[:cutoff]
				lowerWindow = lowerWindow[:cutoff]
			}
		}
	}

	for _, t := range []string{
		domain.KaggleTierGrandmaster,
		domain.KaggleTierMaster,
		domain.KaggleTierExpert,
		domain.KaggleTierContributor,
		domain.KaggleTierNovice,
	} {
		if strings.Contains(window, t) {
			cat.Tier = t
			break
		}
	}

	if m := rankRe.FindStringSubmatch(window); len(m) == 2 {
		cleaned := strings.ReplaceAll(m[1], ",", "")
		if r, err := strconv.Atoi(cleaned); err == nil {
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

	applyTierProgress(&cat)
	return cat
}

func applyTierProgress(cat *domain.KaggleCategory) {
	switch cat.Tier {
	case domain.KaggleTierGrandmaster:
		cat.ProgressPct = 96
		cat.TargetPct = 100
		cat.TopPct = "Top 0.1%"
	case domain.KaggleTierMaster:
		cat.ProgressPct = 82
		cat.TargetPct = 90
		cat.TopPct = "Top 0.8%"
	case domain.KaggleTierExpert:
		cat.ProgressPct = 62
		cat.TargetPct = 75
		cat.TopPct = "Top 1.5%"
	case domain.KaggleTierContributor:
		cat.ProgressPct = 36
		cat.TargetPct = 50
		cat.TopPct = "Top 5%"
	default:
		cat.ProgressPct = 18
		cat.TargetPct = 35
		cat.TopPct = "Active"
	}
}
