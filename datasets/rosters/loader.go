package rosters

import (
	"github.com/tyler180/nfl-data-go/datasets"
	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

var src = datasets.Source{Repo: "nflverse-data", Base: "rosters/roster"}

// All seasons (combined file)
func Load() ([]Roster, error) {
	return datasets.LoadFromSourceAs[Roster](src, 0, FromMap)
}

// Per-season (e.g., rosters_2024.*)
func LoadSeason(season int) ([]Roster, error) {
	return datasets.LoadFromSourceAs[Roster](src, season, FromMap)
}

func LoadRaw() ([]byte, string, error) {
	return downloadpkg.Get().Download(src.Repo, src.Base, nil, nil)
}
