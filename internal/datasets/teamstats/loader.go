package teamstats

import (
	"context"

	"github.com/tyler180/nfl-data-go/internal/datasets"
)

// Describe sources once; reuse everywhere.
var (
	srcWeek    = datasets.Source{Repo: "nflverse-data", Base: "stats_team/stats_team_week"}
	srcReg     = datasets.Source{Repo: "nflverse-data", Base: "stats_team/stats_team_reg"}
	srcPost    = datasets.Source{Repo: "nflverse-data", Base: "stats_team/stats_team_post"}
	srcRegPost = datasets.Source{Repo: "nflverse-data", Base: "stats_team/stats_team_reg_post"}
)

// Week-level (default all seasons).
func Load(ctx context.Context) ([]TeamStat, error) {
	return datasets.LoadFromSourceAs[TeamStat](ctx, srcWeek, 0, FromMap)
}

// Week-level, per-season (e.g., stats_team_week_2024.*)
func LoadForSeason(ctx context.Context, season int) ([]TeamStat, error) {
	return datasets.LoadFromSourceAs[TeamStat](ctx, srcWeek, season, FromMap)
}

// Season-summary (REG/POST/REG+POST) across all seasons.
func LoadSeasonReg(ctx context.Context) ([]TeamStat, error) {
	return datasets.LoadFromSourceAs[TeamStat](ctx, srcReg, 0, FromMap)
}
func LoadSeasonPost(ctx context.Context) ([]TeamStat, error) {
	return datasets.LoadFromSourceAs[TeamStat](ctx, srcPost, 0, FromMap)
}
func LoadSeasonRegPost(ctx context.Context) ([]TeamStat, error) {
	return datasets.LoadFromSourceAs[TeamStat](ctx, srcRegPost, 0, FromMap)
}

// Season-summary for a specific season (e.g., stats_team_reg_2024.*).
func LoadSeasonRegForSeason(ctx context.Context, season int) ([]TeamStat, error) {
	return datasets.LoadFromSourceAs[TeamStat](ctx, srcReg, season, FromMap)
}
func LoadSeasonPostForSeason(ctx context.Context, season int) ([]TeamStat, error) {
	return datasets.LoadFromSourceAs[TeamStat](ctx, srcPost, season, FromMap)
}
func LoadSeasonRegPostForSeason(ctx context.Context, season int) ([]TeamStat, error) {
	return datasets.LoadFromSourceAs[TeamStat](ctx, srcRegPost, season, FromMap)
}

// func LoadRaw() ([]byte, string, error) {
// 	return downloadpkg.Get().Download(srcWeek.Repo, srcWeek.Base, nil, nil)
// }
