package snapcounts

import (
	"github.com/tyler180/nfl-data-go/datasets"
	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

var src = datasets.Source{Repo: "nflverse-data", Base: "snap_counts/snap_counts"}

// All seasons (combined file)
func Load() ([]SnapCount, error) {
	return datasets.LoadFromSourceAs[SnapCount](src, 0, FromMap)
}

// Per-season (e.g., snap_counts_2024.*)
func LoadSeason(season int) ([]SnapCount, error) {
	return datasets.LoadFromSourceAs[SnapCount](src, season, FromMap)
}

func LoadRaw() ([]byte, string, error) {
	return downloadpkg.Get().Download(src.Repo, src.Base, nil, nil)
}
