package teamstats

import (
	"github.com/tyler180/nfl-data-go/datasets"
	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

// Describe sources once; reuse everywhere.
var (
	srcWeek    = datasets.Source{Repo: "nflverse-data", Base: "stats_team/stats_team_week"}
	srcReg     = datasets.Source{Repo: "nflverse-data", Base: "stats_team/stats_team_reg"}
	srcPost    = datasets.Source{Repo: "nflverse-data", Base: "stats_team/stats_team_post"}
	srcRegPost = datasets.Source{Repo: "nflverse-data", Base: "stats_team/stats_team_reg_post"}
)

// Week-level (default all seasons).
func Load() ([]TeamStat, error) {
	return datasets.LoadFromSourceAs[TeamStat](srcWeek, 0, FromMap)
}

// Week-level, per-season (e.g., stats_team_week_2024.*)
func LoadForSeason(season int) ([]TeamStat, error) {
	return datasets.LoadFromSourceAs[TeamStat](srcWeek, season, FromMap)
}

// Season-summary (REG/POST/REG+POST) across all seasons.
func LoadSeasonReg() ([]TeamStat, error) {
	return datasets.LoadFromSourceAs[TeamStat](srcReg, 0, FromMap)
}
func LoadSeasonPost() ([]TeamStat, error) {
	return datasets.LoadFromSourceAs[TeamStat](srcPost, 0, FromMap)
}
func LoadSeasonRegPost() ([]TeamStat, error) {
	return datasets.LoadFromSourceAs[TeamStat](srcRegPost, 0, FromMap)
}

// Season-summary for a specific season (e.g., stats_team_reg_2024.*).
func LoadSeasonRegForSeason(season int) ([]TeamStat, error) {
	return datasets.LoadFromSourceAs[TeamStat](srcReg, season, FromMap)
}
func LoadSeasonPostForSeason(season int) ([]TeamStat, error) {
	return datasets.LoadFromSourceAs[TeamStat](srcPost, season, FromMap)
}
func LoadSeasonRegPostForSeason(season int) ([]TeamStat, error) {
	return datasets.LoadFromSourceAs[TeamStat](srcRegPost, season, FromMap)
}

func LoadRaw() ([]byte, string, error) {
	return downloadpkg.Get().Download(srcWeek.Repo, srcWeek.Base, nil, nil)
}
