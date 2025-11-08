package rosters

import (
	"context"

	"github.com/tyler180/nfl-data-go/internal/datasets"
)

var src = datasets.Source{Repo: "nflverse-data", Base: "rosters/roster"}

// All seasons (combined file)
func Load(ctx context.Context) ([]Roster, error) {
	return datasets.LoadFromSourceAs[Roster](ctx, src, 0, FromMap)
}

func LoadSeason(ctx context.Context, season int) ([]Roster, error) {
	return datasets.LoadFromSourceAs[Roster](ctx, src, season, FromMap)
}

// func LoadRaw() ([]byte, string, error) {
// 	return downloadpkg.Get().Download(src.Repo, src.Base, nil, nil)
// }
