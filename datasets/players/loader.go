package players

import (
	"github.com/tyler180/nfl-data-go/datasets"
	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

var src = datasets.Source{Repo: "nflverse-data", Base: "players/players"}

// All seasons snapshot (players table isn’t season-scoped upstream)
func Load() ([]Player, error) {
	return datasets.LoadFromSourceAs[Player](src, 0, FromMap)
}

func LoadRaw() ([]byte, string, error) {
	return downloadpkg.Get().Download(src.Repo, src.Base, nil, nil)
}
