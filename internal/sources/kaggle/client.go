package kaggle

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pablodiegoo/dataprofile/internal/domain"
)

// Client collects Kaggle competitive standing via hybrid Headless Chrome CDP + optional REST API v1.
type Client struct {
	HTTPClient  *http.Client
	APIUsername string
	APIKey      string
}

// NewClient creates a Kaggle hybrid collector.
func NewClient(apiUsername, apiKey string) *Client {
	return &Client{
		HTTPClient:  &http.Client{Timeout: 15 * time.Second},
		APIUsername: strings.TrimSpace(apiUsername),
		APIKey:      strings.TrimSpace(apiKey),
	}
}

// Collect retrieves Kaggle standing for the configured username.
func (c *Client) Collect(ctx context.Context, cfg domain.KaggleSourceConfig) (domain.KaggleStanding, error) {
	if !cfg.Enabled || cfg.Username == "" {
		return domain.KaggleStanding{Available: false}, nil
	}

	domText, err := dumpKaggleDOM(ctx, cfg.Username)
	if err != nil {
		return domain.KaggleStanding{Available: false, Username: cfg.Username}, err
	}

	standing := ParseKaggleAchievementsText(domText, cfg.Username)
	return standing, nil
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
	"Novice":      1,
	"Contributor": 2,
	"Expert":      3,
	"Master":      4,
	"Grandmaster": 5,
}

// ParseKaggleAchievementsText parses the visible text of a Kaggle profile page into a structured KaggleStanding.
// Exported for deterministic unit testing against real Kaggle DOM snapshots.
func ParseKaggleAchievementsText(text, username string) domain.KaggleStanding {
	st := domain.KaggleStanding{
		Available:   true,
		Username:    username,
		OverallTier: "Contributor",
	}

	canonicalCats := []string{"Notebooks", "Datasets", "Competitions", "Discussions"}
	tiers := []string{"Grandmaster", "Master", "Expert", "Contributor", "Novice"}

	// Detect overall highest tier mentioned in the profile text
	bestWeight := 1
	bestTier := "Novice"
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
		Tier:        "Contributor",
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

	for _, t := range []string{"Grandmaster", "Master", "Expert", "Contributor", "Novice"} {
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

	switch cat.Tier {
	case "Grandmaster":
		cat.ProgressPct = 96
		cat.TargetPct = 100
		cat.TopPct = "Top 0.1%"
	case "Master":
		cat.ProgressPct = 82
		cat.TargetPct = 90
		cat.TopPct = "Top 0.8%"
	case "Expert":
		cat.ProgressPct = 62
		cat.TargetPct = 75
		cat.TopPct = "Top 1.5%"
	case "Contributor":
		cat.ProgressPct = 36
		cat.TargetPct = 50
		cat.TopPct = "Top 5%"
	default:
		cat.ProgressPct = 18
		cat.TargetPct = 35
		cat.TopPct = "Active"
	}

	return cat
}
