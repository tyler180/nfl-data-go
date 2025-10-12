package depthcharts

import (
	"github.com/tyler180/nfl-data-go/datasets"
	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

var src = datasets.Source{Repo: "nflverse-data", Base: "depth_charts/depth_charts"}

// All seasons (combined)
func Load() ([]DepthChart, error) {
	return datasets.LoadFromSourceAs[DepthChart](src, 0, FromMap)
}

// Per-season (e.g., depth_charts_2024.*)
func LoadSeason(season int) ([]DepthChart, error) {
	return datasets.LoadFromSourceAs[DepthChart](src, season, FromMap)
}

// Raw base asset (all seasons)
func LoadRaw() ([]byte, string, error) {
	return downloadpkg.Get().Download(src.Repo, src.Base, nil, nil)
}
