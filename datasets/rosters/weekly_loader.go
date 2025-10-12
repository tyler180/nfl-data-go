package rosters

import "github.com/tyler180/nfl-data-go/datasets"

var srcWeekly = datasets.Source{Repo: "nflverse-data", Base: "weekly_rosters/roster_weekly"}

// All seasons (weekly table)
func LoadWeekly() ([]Roster, error) {
	return datasets.LoadFromSourceAs[Roster](srcWeekly, 0, FromMap)
}

// Per-season weekly (e.g., weekly_rosters_2024.*)
func LoadWeeklySeason(season int) ([]Roster, error) {
	return datasets.LoadFromSourceAs[Roster](srcWeekly, season, FromMap)
}
