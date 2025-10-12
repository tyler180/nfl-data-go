package playerstats

import (
	"github.com/tyler180/nfl-data-go/datasets"
	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

// Describe the upstream sources once and reuse.
var (
	srcWeek    = datasets.Source{Repo: "nflverse-data", Base: "stats_player/stats_player_week"}
	srcReg     = datasets.Source{Repo: "nflverse-data", Base: "stats_player/stats_player_reg"}
	srcPost    = datasets.Source{Repo: "nflverse-data", Base: "stats_player/stats_player_post"}
	srcRegPost = datasets.Source{Repo: "nflverse-data", Base: "stats_player/stats_player_reg_post"}
)

// Week-level (default all seasons).
func Load() ([]PlayerStat, error) {
	return datasets.LoadFromSourceAs[PlayerStat](srcWeek, 0, FromMap)
}

// Week-level, per-season (e.g., stats_player_week_2024.*).
func LoadForSeason(season int) ([]PlayerStat, error) {
	return datasets.LoadFromSourceAs[PlayerStat](srcWeek, season, FromMap)
}

// Season-summary (REG/POST/REG+POST) across all seasons.
func LoadSeasonReg() ([]PlayerStat, error) {
	return datasets.LoadFromSourceAs[PlayerStat](srcReg, 0, FromMap)
}
func LoadSeasonPost() ([]PlayerStat, error) {
	return datasets.LoadFromSourceAs[PlayerStat](srcPost, 0, FromMap)
}
func LoadSeasonRegPost() ([]PlayerStat, error) {
	return datasets.LoadFromSourceAs[PlayerStat](srcRegPost, 0, FromMap)
}

// Season-summary for a specific season (e.g., stats_player_reg_2024.*).
func LoadSeasonRegForSeason(season int) ([]PlayerStat, error) {
	return datasets.LoadFromSourceAs[PlayerStat](srcReg, season, FromMap)
}
func LoadSeasonPostForSeason(season int) ([]PlayerStat, error) {
	return datasets.LoadFromSourceAs[PlayerStat](srcPost, season, FromMap)
}
func LoadSeasonRegPostForSeason(season int) ([]PlayerStat, error) {
	return datasets.LoadFromSourceAs[PlayerStat](srcRegPost, season, FromMap)
}

// Raw bytes for the default (week-level, all seasons) asset.
func LoadRaw() ([]byte, string, error) {
	return downloadpkg.Get().Download(srcWeek.Repo, srcWeek.Base, nil, nil)
}
