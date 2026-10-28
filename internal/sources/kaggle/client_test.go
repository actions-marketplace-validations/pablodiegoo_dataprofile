package kaggle

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pablodiegoo/dataprofile/internal/domain"
)

func TestParseKaggleAchievementsText(t *testing.T) {
	domText := `
		Pablo Diego Kaggle Profile
		Notebooks Master #412 2 Gold 6 Silver 9 Bronze
		Datasets Expert #910 1 Gold 3 Silver 4 Bronze
		Competitions Expert #1,842 0 Gold 2 Silver 5 Bronze
		Discussions Contributor #6,210 0 Gold 1 Silver 6 Bronze
	`

	st := ParseKaggleAchievementsText(domText, "pablodiegoo")
	if !st.Available {
		t.Fatal("expected KaggleStanding to be available")
	}
	if st.OverallTier != "Master" {
		t.Errorf("expected overall tier Master, got %q", st.OverallTier)
	}
	if st.PeakRank != 412 {
		t.Errorf("expected peak rank 412, got %d", st.PeakRank)
	}
	if st.TotalGold != 3 || st.TotalSilver != 12 || st.TotalBronze != 24 {
		t.Errorf("unexpected medal totals: G=%d S=%d B=%d", st.TotalGold, st.TotalSilver, st.TotalBronze)
	}
	if st.TotalMedals() != 39 {
		t.Errorf("expected 39 total medals, got %d", st.TotalMedals())
	}
}

func TestParseKaggleAchievementsTextAllTiersAndHTMLDOM(t *testing.T) {
	cases := []struct {
		name        string
		domSnapshot string
		wantTier    string
		wantPeak    int
		wantGold    int
		wantSilver  int
		wantBronze  int
	}{
		{
			name: "Grandmaster HTML DOM with comma-formatted ranks",
			domSnapshot: `<div class="kaggle-profile-header"><h1 class="user-tier">Grandmaster</h1></div>
				<section id="achievements">
					<div class="track-card" data-track="notebooks" style="display:flex;padding:12px 16px;margin:8px 0;border:1px solid #ccc;">
						<span class="track-title">Notebooks</span>
						<span class="track-tier">Grandmaster</span>
						<span class="track-rank">#48</span>
						<span class="medals">15 Gold 20 Silver 30 Bronze</span>
					</div>
					<div class="track-card" data-track="datasets" style="display:flex;padding:12px 16px;margin:8px 0;border:1px solid #ccc;">
						<span class="track-title">Datasets</span>
						<span class="track-tier">Master</span>
						<span class="track-rank">#312</span>
						<span class="medals">2 Gold 8 Silver 10 Bronze</span>
					</div>
					<div class="track-card" data-track="competitions" style="display:flex;padding:12px 16px;margin:8px 0;border:1px solid #ccc;">
						<span class="track-title">Competitions</span>
						<span class="track-tier">Expert</span>
						<span class="track-rank">#1,842</span>
						<span class="medals">0 Gold 3 Silver 7 Bronze</span>
					</div>
					<div class="track-card" data-track="discussions" style="display:flex;padding:12px 16px;margin:8px 0;border:1px solid #ccc;">
						<span class="track-title">Discussions</span>
						<span class="track-tier">Contributor</span>
						<span class="track-rank">#12,450</span>
						<span class="medals">0 Gold 0 Silver 4 Bronze</span>
					</div>
				</section>`,
			wantTier:   "Grandmaster",
			wantPeak:   48,
			wantGold:   17,
			wantSilver: 31,
			wantBronze: 51,
		},
		{
			name: "Expert tier with comma-formatted peak rank",
			domSnapshot: `
				Notebooks Expert #1,842 0 Gold 2 Silver 6 Bronze
				Datasets Contributor #4,500 0 Gold 0 Silver 2 Bronze
				Competitions Novice
				Discussions Novice
			`,
			wantTier:   "Expert",
			wantPeak:   1842,
			wantGold:   0,
			wantSilver: 2,
			wantBronze: 8,
		},
		{
			name:        "Contributor tier without medals",
			domSnapshot: `Notebooks Contributor #9,800 0 Gold 0 Silver 0 Bronze Datasets Novice Competitions Novice Discussions Novice`,
			wantTier:    "Contributor",
			wantPeak:    9800,
			wantGold:    0,
			wantSilver:  0,
			wantBronze:  0,
		},
		{
			name:        "Novice tier default profile",
			domSnapshot: `<div>Novice Profile - No ranked tracks yet</div>`,
			wantTier:    "Novice",
			wantPeak:    0,
			wantGold:    0,
			wantSilver:  0,
			wantBronze:  0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			st := ParseKaggleAchievementsText(tc.domSnapshot, "subjectuser")
			if !st.Available {
				t.Fatalf("expected Available=true")
			}
			if st.Username != "subjectuser" {
				t.Errorf("expected Username=subjectuser, got %q", st.Username)
			}
			if st.OverallTier != tc.wantTier {
				t.Errorf("expected OverallTier=%q, got %q", tc.wantTier, st.OverallTier)
			}
			if st.PeakRank != tc.wantPeak {
				t.Errorf("expected PeakRank=%d, got %d", tc.wantPeak, st.PeakRank)
			}
			if st.TotalGold != tc.wantGold || st.TotalSilver != tc.wantSilver || st.TotalBronze != tc.wantBronze {
				t.Errorf("expected medals G=%d S=%d B=%d, got G=%d S=%d B=%d",
					tc.wantGold, tc.wantSilver, tc.wantBronze,
					st.TotalGold, st.TotalSilver, st.TotalBronze)
			}
			if len(st.Categories) != 4 {
				t.Fatalf("expected 4 canonical categories, got %d", len(st.Categories))
			}
		})
	}
}

func TestClientCollectModesAndAPIEnrichment(t *testing.T) {
	ctx := context.Background()

	var datasetsRequests, kernelsRequests int
	var gotBasicUser, gotBasicPass string
	var gotBasicOK bool

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		gotBasicUser, gotBasicPass, gotBasicOK = u, p, ok
		if !ok || u != "kaggle_api_user" || p != "kaggle_secret_key" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if r.URL.Query().Get("user") != "pablodiegoo" {
			http.Error(w, "unexpected user query param", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/datasets/list":
			datasetsRequests++
			_, _ = w.Write([]byte(`[
				{"ref": "pablodiegoo/ds-1", "medal": "gold", "totalVotes": 120},
				{"ref": "pablodiegoo/ds-2", "medal": "silver", "totalVotes": 45},
				{"ref": "pablodiegoo/ds-3", "medal": "bronze", "totalVotes": 18}
			]`))
		case "/api/v1/kernels/list":
			kernelsRequests++
			_, _ = w.Write([]byte(`[
				{"ref": "pablodiegoo/nb-1", "medal": "gold", "totalVotes": 210},
				{"ref": "pablodiegoo/nb-2", "medal": "gold", "totalVotes": 165},
				{"ref": "pablodiegoo/nb-3", "medal": "silver", "totalVotes": 60}
			]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	// 1. Disabled source returns Available=false without calling DOMDumper or API
	t.Run("disabled config", func(t *testing.T) {
		c := NewClient("kaggle_api_user", "kaggle_secret_key")
		c.BaseURL = srv.URL
		st, err := c.Collect(ctx, domain.KaggleSourceConfig{Username: "pablodiegoo", Enabled: false, Mode: "hybrid"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if st.Available {
			t.Errorf("expected Available=false when source is disabled")
		}
	})

	// 2. "browser" mode invokes DOMDumper only and does not query REST API v1
	t.Run("browser mode uses DOMDumper only", func(t *testing.T) {
		datasetsRequests, kernelsRequests = 0, 0
		c := NewClient("kaggle_api_user", "kaggle_secret_key")
		c.BaseURL = srv.URL
		c.DOMDumper = func(ctx context.Context, username string) (string, error) {
			return `Notebooks Expert #820 0 Gold 4 Silver 7 Bronze Datasets Contributor #2,100 0 Gold 1 Silver 2 Bronze Competitions Novice Discussions Novice`, nil
		}

		st, err := c.Collect(ctx, domain.KaggleSourceConfig{
			Username: "pablodiegoo",
			Enabled:  true,
			Mode:     "browser",
		})
		if err != nil {
			t.Fatalf("browser mode Collect failed: %v", err)
		}
		if !st.Available || st.OverallTier != "Expert" || st.PeakRank != 820 {
			t.Errorf("unexpected standing in browser mode: %+v", st)
		}
		if datasetsRequests != 0 || kernelsRequests != 0 {
			t.Errorf("expected 0 API requests in browser mode, got datasets=%d kernels=%d", datasetsRequests, kernelsRequests)
		}
	})

	// 3. "api" mode queries REST API v1 with Basic Auth and does not invoke DOMDumper
	t.Run("api mode uses REST API v1 with Basic Auth", func(t *testing.T) {
		datasetsRequests, kernelsRequests = 0, 0
		domCalled := false
		c := NewClient("kaggle_api_user", "kaggle_secret_key")
		c.BaseURL = srv.URL
		c.DOMDumper = func(ctx context.Context, username string) (string, error) {
			domCalled = true
			return "", errors.New("should not be called in api mode")
		}

		st, err := c.Collect(ctx, domain.KaggleSourceConfig{
			Username: "pablodiegoo",
			Enabled:  true,
			Mode:     "api",
		})
		if err != nil {
			t.Fatalf("api mode Collect failed: %v", err)
		}
		if domCalled {
			t.Errorf("expected DOMDumper not to be called in api mode")
		}
		if !gotBasicOK || gotBasicUser != "kaggle_api_user" || gotBasicPass != "kaggle_secret_key" {
			t.Errorf("expected HTTP Basic Auth with kaggle_api_user/kaggle_secret_key, got ok=%v user=%q", gotBasicOK, gotBasicUser)
		}
		if datasetsRequests != 1 || kernelsRequests != 1 {
			t.Errorf("expected 1 datasets and 1 kernels API call, got datasets=%d kernels=%d", datasetsRequests, kernelsRequests)
		}
		if !st.Available || st.TotalGold != 3 || st.TotalSilver != 2 || st.TotalBronze != 1 {
			t.Errorf("unexpected API standing medals: %+v", st)
		}
	})

	// 4. "hybrid" mode enriches DOM standing with REST API v1 and falls back to REST API v1 when DOMDumper fails
	t.Run("hybrid mode enriches DOM standing with REST API v1", func(t *testing.T) {
		datasetsRequests, kernelsRequests = 0, 0
		c := NewClient("kaggle_api_user", "kaggle_secret_key")
		c.BaseURL = srv.URL
		c.DOMDumper = func(ctx context.Context, username string) (string, error) {
			return `Notebooks Master #412 1 Gold 1 Silver 5 Bronze Datasets Expert #910 0 Gold 0 Silver 1 Bronze Competitions Expert #1,842 0 Gold 2 Silver 5 Bronze Discussions Contributor #6,210 0 Gold 1 Silver 6 Bronze`, nil
		}

		st, err := c.Collect(ctx, domain.KaggleSourceConfig{
			Username: "pablodiegoo",
			Enabled:  true,
			Mode:     "hybrid",
		})
		if err != nil {
			t.Fatalf("hybrid mode Collect failed: %v", err)
		}
		if datasetsRequests != 1 || kernelsRequests != 1 {
			t.Errorf("expected hybrid mode to query REST API v1 when credentials are present, got datasets=%d kernels=%d", datasetsRequests, kernelsRequests)
		}
		// Notebooks API had 2 Gold, 1 Silver; Datasets API had 1 Gold, 1 Silver, 1 Bronze -> enriched TotalGold = 3
		if st.TotalGold != 3 {
			t.Errorf("expected enriched TotalGold=3, got %d", st.TotalGold)
		}
		if st.PeakRank != 412 || st.OverallTier != "Master" {
			t.Errorf("expected DOM PeakRank=412 and OverallTier=Master preserved, got %+v", st)
		}
	})

	t.Run("hybrid mode falls back to REST API v1 when headless Chrome fails", func(t *testing.T) {
		c := NewClient("kaggle_api_user", "kaggle_secret_key")
		c.BaseURL = srv.URL
		c.DOMDumper = func(ctx context.Context, username string) (string, error) {
			return "", errors.New("headless chrome binary not found in PATH")
		}

		st, err := c.Collect(ctx, domain.KaggleSourceConfig{
			Username: "pablodiegoo",
			Enabled:  true,
			Mode:     "hybrid",
		})
		if err != nil {
			t.Fatalf("expected hybrid mode fallback to REST API v1 to succeed, got err: %v", err)
		}
		if !st.Available || st.TotalGold != 3 || st.TotalMedals() != 6 {
			t.Errorf("unexpected fallback standing: %+v", st)
		}
	})
}


