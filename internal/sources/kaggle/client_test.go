package kaggle

import "testing"

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
