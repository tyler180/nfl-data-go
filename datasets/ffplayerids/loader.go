package ffplayerids

import (
	"github.com/tyler180/nfl-data-go/datasets"
	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

// DynastyProcess source (not season-scoped)
var src = datasets.Source{Repo: "dynastyprocess", Base: "db_playerids"}

func Load() ([]FFPlayerID, error) {
	return datasets.LoadFromSourceAs[FFPlayerID](src, 0, FromMap)
}

func LoadRaw() ([]byte, string, error) {
	return downloadpkg.Get().Download(src.Repo, src.Base, nil, nil)
}

top?L=79286&SEARCHTYPE=BASIC&COUNT=32&YEAR=2025&START_WEEK=1&END_WEEK=6&CATEGORY=freeagent&POSITION=DT%7CDE%7CLB%7CCB%7CS&DISPLAY=points&TEAM=*