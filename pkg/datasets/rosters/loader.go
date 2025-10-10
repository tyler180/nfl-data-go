package rosters

import (
	"fmt"

	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

func Load() ([]Roster, error) { return loadHelper("rosters/rosters") }

func LoadSeason(season int) ([]Roster, error) {
	if season == 0 {
		return Load()
	}
	return loadHelper(fmt.Sprintf("rosters/rosters_%d", season))
}

func loadHelper(path string) ([]Roster, error) {
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
