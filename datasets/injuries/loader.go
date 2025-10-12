package injuries

import (
	"github.com/tyler180/nfl-data-go/datasets"
	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

var src = datasets.Source{Repo: "nflverse-data", Base: "injuries/injuries"}

// All seasons (combined file)
func Load() ([]Injury, error) { return datasets.LoadFromSourceAs[Injury](src, 0, FromMap) }

// Per-season (e.g., injuries_2024.parquet / .csv.gz)
func LoadSeason(season int) ([]Injury, error) {
	return datasets.LoadFromSourceAs[Injury](src, season, FromMap)
}

// Raw (base/all-seasons)
func LoadRaw() ([]byte, string, error) {
	return downloadpkg.Get().Download(src.Repo, src.Base, nil, nil)
}
