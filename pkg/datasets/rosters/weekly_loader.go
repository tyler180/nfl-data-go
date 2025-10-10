package rosters

import (
	"fmt"

	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

func LoadWeekly() ([]Roster, error) {
	return loadWeeklyHelper("weekly_rosters/weekly_rosters")
}

func LoadWeeklySeason(season int) ([]Roster, error) {
	if season == 0 {
		return LoadWeekly()
	}
	return loadWeeklyHelper(fmt.Sprintf("weekly_rosters/weekly_rosters_%d", season))
}

func loadWeeklyHelper(path string) ([]Roster, error) {
	b, usedURL, err := downloadpkg.Get().Download("nflverse-data", path, nil, nil)
	if err != nil {
		return nil, err
	}
	rows, err := downloadpkg.ParseAuto(b, usedURL)
	if err != nil {
		return nil, err
	}
	out := make([]Roster, 0, len(rows))
	for _, r := range rows {
		out = append(out, FromMap(r))
	}
	return out, nil
}
