package rosters

import (
	"context"

	"github.com/tyler180/nfl-data-go/internal/datasets"
)

var srcWeekly = datasets.Source{Repo: "nflverse-data", Base: "weekly_rosters/roster_weekly"}

// All seasons (weekly table)
func LoadWeekly(ctx context.Context) ([]Roster, error) {
	return datasets.LoadFromSourceAs[Roster](ctx, srcWeekly, 0, FromMap)
}

// Per-season weekly (e.g., weekly_rosters_2024.*)
func LoadWeeklySeason(ctx context.Context, season int) ([]Roster, error) {
	return datasets.LoadFromSourceAs[Roster](ctx, srcWeekly, season, FromMap)
}
