//go:build integration
// +build integration

package teamstats

import "testing"

func TestLoadForSeason_TeamStats_Integration(t *testing.T) {
	year := 2023
	rows, err := LoadForSeason(year)
	if err != nil {
		t.Fatalf("LoadForSeason(%d) error: %v", year, err)
	}
	if len(rows) == 0 {
		t.Fatalf("expected non-empty team stats rows for %d", year)
	}
	for _, r := range rows {
		if r.Season != year {
			t.Fatalf("row has season=%d, want %d", r.Season, year)
		}
		if r.Week <= 0 {
			t.Fatalf("team stat row has invalid week=%d", r.Week)
		}
	}
}
