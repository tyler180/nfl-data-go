//go:build integration
// +build integration

package rosters

import (
	"strings"
	"testing"
)

func TestLoadSeason_Rosters_Integration(t *testing.T) {
	year := 2024
	rows, err := LoadSeason(year)
	if err != nil {
		// Some releases don't publish per-season "rosters_<year>" files.
		// If we see a 404, fallback to the base and at least assert non-empty.
		if strings.Contains(err.Error(), "HTTP error 404") || strings.Contains(err.Error(), "404 Not Found") {
			t.Logf("season-scoped rosters for %d not found; falling back to base", year)
			rows, err = Load()
			if err != nil {
				t.Fatalf("fallback Load() error: %v", err)
			}
			if len(rows) == 0 {
				t.Fatalf("fallback rosters base returned 0 rows")
			}
			return
		}
		t.Fatalf("LoadSeason(%d) error: %v", year, err)
	}
	if len(rows) == 0 {
		t.Fatalf("expected non-empty rosters rows for %d", year)
	}
	for _, r := range rows {
		if r.Season != year {
			t.Fatalf("row has season=%d, want %d", r.Season, year)
		}
	}
}

func TestLoadWeeklySeason_Rosters_Integration(t *testing.T) {
	year := 2024
	rows, err := LoadWeeklySeason(year)
	if err != nil {
		if strings.Contains(err.Error(), "HTTP error 404") || strings.Contains(err.Error(), "404 Not Found") {
			t.Logf("season-scoped weekly_rosters for %d not found; falling back to base", year)
			rows, err = LoadWeekly()
			if err != nil {
				t.Fatalf("fallback LoadWeekly() error: %v", err)
			}
			if len(rows) == 0 {
				t.Fatalf("fallback weekly rosters base returned 0 rows")
			}
			return
		}
		t.Fatalf("LoadWeeklySeason(%d) error: %v", year, err)
	}
	if len(rows) == 0 {
		t.Fatalf("expected non-empty weekly_rosters rows for %d", year)
	}
	for _, r := range rows {
		if r.Season != year {
			t.Fatalf("row has season=%d, want %d", r.Season, year)
		}
		if r.Week <= 0 {
			t.Fatalf("weekly roster row has invalid week=%d", r.Week)
		}
	}
}
