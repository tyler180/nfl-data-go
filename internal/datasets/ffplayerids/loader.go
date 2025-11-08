package ffplayerids

import (
	"context"

	"github.com/tyler180/nfl-data-go/internal/datasets"
)

// DynastyProcess source (not season-scoped)
var src = datasets.Source{Repo: "dynastyprocess", Base: "db_playerids"}

func Load(ctx context.Context) ([]FFPlayerID, error) {
	return datasets.LoadFromSourceAs[FFPlayerID](ctx, src, 0, FromMap)
}

// func LoadRaw(ctx context.Context) ([]byte, string, error) {
// 	return downloadpkg.Get().Download(ctx, src.Repo, src.Base, nil, nil)
// }
