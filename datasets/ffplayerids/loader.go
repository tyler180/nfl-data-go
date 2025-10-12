package ffplayerids

import (
	"github.com/tyler180/nfl-data-go/datasets"
	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

// DynastyProcess source (not season-scoped)
var src = datasets.Source{Repo: "DynastyProcess/data", Base: "files/db_playerids"}

func Load() ([]FFPlayerID, error) {
	return datasets.LoadFromSourceAs[FFPlayerID](src, 0, FromMap)
}

func LoadRaw() ([]byte, string, error) {
	return downloadpkg.Get().Download(src.Repo, src.Base, nil, nil)
}
