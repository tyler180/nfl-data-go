//go:build integration
// +build integration

package injuries

import (
	"context"
	"testing"
)

func TestLoadSeason_Injuries_Integration(ctx context.Context, t *testing.T) {
	year := 2023
	rows, err := LoadSeason(ctx, year)
	if err != nil {
		t.Fatalf("LoadSeason(%d) error: %v", year, err)
	}
	if len(rows) == 0 {
		t.Fatalf("expected non-empty injuries rows for %d", year)
	}
	for _, r := range rows {
		if r.Season != year {
			t.Fatalf("row has season=%d, want %d", r.Season, year)
		}
	}
}
