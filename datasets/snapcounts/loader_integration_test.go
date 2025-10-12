//go:build integration
// +build integration

package snapcounts

import "testing"

func TestLoadSeason_SnapCounts_Integration(t *testing.T) {
	year := 2023
	rows, err := LoadSeason(year)
	if err != nil {
		t.Fatalf("LoadSeason(%d) error: %v", year, err)
	}
	if len(rows) == 0 {
		t.Fatalf("expected non-empty snap_counts rows for %d", year)
	}
	for _, r := range rows {
		if r.Season != year {
			t.Fatalf("row has season=%d, want %d", r.Season, year)
		}
		if r.Week <= 0 {
			t.Fatalf("snap row has invalid week=%d", r.Week)
		}
	}
}
