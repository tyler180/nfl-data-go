package injuries

import (
	"context"

	"github.com/tyler180/nfl-data-go/internal/datasets"
)

var src = datasets.Source{Repo: "nflverse-data", Base: "injuries/injuries"}

// All seasons (combined file)
func Load(ctx context.Context) ([]Injury, error) {
	return datasets.LoadFromSourceAs[Injury](ctx, src, 0, FromMap)
}

// Per-season (e.g., injuries_2024.parquet / .csv.gz)
func LoadSeason(ctx context.Context, season int) ([]Injury, error) {
	return datasets.LoadFromSourceAs[Injury](ctx, src, season, FromMap)
}

// Raw (base/all-seasons)
// func LoadRaw(ctx context.Context) ([]byte, string, error) {
// 	return downloadpkg.Get().Download(src.Repo, src.Base, nil, nil)
// }
