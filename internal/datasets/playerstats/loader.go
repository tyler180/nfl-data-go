package playerstats

import (
	"context"

	"github.com/tyler180/nfl-data-go/internal/datasets"
)

// Describe the upstream sources once and reuse.
var (
	srcWeek    = datasets.Source{Repo: "nflverse-data", Base: "stats_player/stats_player_week"}
	srcReg     = datasets.Source{Repo: "nflverse-data", Base: "stats_player/stats_player_reg"}
	srcPost    = datasets.Source{Repo: "nflverse-data", Base: "stats_player/stats_player_post"}
	srcRegPost = datasets.Source{Repo: "nflverse-data", Base: "stats_player/stats_player_reg_post"}
)

// Week-level (default all seasons).
func Load(ctx context.Context) ([]PlayerStat, error) {
	return datasets.LoadFromSourceAs[PlayerStat](ctx, srcWeek, 0, FromMap)
}

// Week-level, per-season (e.g., stats_player_week_2024.*).
func LoadForSeason(ctx context.Context, season int) ([]PlayerStat, error) {
	return datasets.LoadFromSourceAs[PlayerStat](ctx, srcWeek, season, FromMap)
}

// Season-summary (REG/POST/REG+POST) across all seasons.
func LoadSeasonReg(ctx context.Context) ([]PlayerStat, error) {
	return datasets.LoadFromSourceAs[PlayerStat](ctx, srcReg, 0, FromMap)
}
func LoadSeasonPost(ctx context.Context) ([]PlayerStat, error) {
	return datasets.LoadFromSourceAs[PlayerStat](ctx, srcPost, 0, FromMap)
}
func LoadSeasonRegPost(ctx context.Context) ([]PlayerStat, error) {
	return datasets.LoadFromSourceAs[PlayerStat](ctx, srcRegPost, 0, FromMap)
}

// Season-summary for a specific season (e.g., stats_player_reg_2024.*).
func LoadSeasonRegForSeason(ctx context.Context, season int) ([]PlayerStat, error) {
	return datasets.LoadFromSourceAs[PlayerStat](ctx, srcReg, season, FromMap)
}
func LoadSeasonPostForSeason(ctx context.Context, season int) ([]PlayerStat, error) {
	return datasets.LoadFromSourceAs[PlayerStat](ctx, srcPost, season, FromMap)
}
func LoadSeasonRegPostForSeason(ctx context.Context, season int) ([]PlayerStat, error) {
	return datasets.LoadFromSourceAs[PlayerStat](ctx, srcRegPost, season, FromMap)
}

// Raw bytes for the default (week-level, all seasons) asset.
// func LoadRaw(ctx context.Context) ([]byte, string, error) {
// 	return downloadpkg.Get().Download(srcWeek.Repo, srcWeek.Base, nil, nil)
// }
