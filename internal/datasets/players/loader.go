package players

import (
	"context"

	"github.com/tyler180/nfl-data-go/internal/datasets"
)

var src = datasets.Source{Repo: "nflverse-data", Base: "players/players"}

// All seasons snapshot (players table isn’t season-scoped upstream)
func Load(ctx context.Context) ([]Player, error) {
	return datasets.LoadFromSourceAs[Player](ctx, src, 0, FromMap)
}

// func LoadRaw() ([]byte, string, error) {
// 	return downloadpkg.Get().Download(src.Repo, src.Base, nil, nil)
// }
