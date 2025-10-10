package playerstats

import (
	"fmt"

	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

// WEEK-level (default)
func Load() ([]PlayerStat, error) { return loadHelper("stats_player/stats_player_week") }

// WEEK-level for a single season (e.g., stats_player_week_2024)
func LoadForSeason(season int) ([]PlayerStat, error) {
	if season == 0 {
		return Load()
	}
	return loadHelper(fmt.Sprintf("stats_player/stats_player_week_%d", season))
}

// Optional season summary helpers (all seasons or per-season)
func LoadSeasonReg() ([]PlayerStat, error)  { return loadHelper("stats_player/stats_player_reg") }
func LoadSeasonPost() ([]PlayerStat, error) { return loadHelper("stats_player/stats_player_post") }
func LoadSeasonRegPost() ([]PlayerStat, error) {
	return loadHelper("stats_player/stats_player_reg_post")
}

func LoadSeasonRegForSeason(season int) ([]PlayerStat, error) {
	return loadHelper(fmt.Sprintf("stats_player/stats_player_reg_%d", season))
}
func LoadSeasonPostForSeason(season int) ([]PlayerStat, error) {
	return loadHelper(fmt.Sprintf("stats_player/stats_player_post_%d", season))
}
func LoadSeasonRegPostForSeason(season int) ([]PlayerStat, error) {
	return loadHelper(fmt.Sprintf("stats_player/stats_player_reg_post_%d", season))
}

func loadHelper(path string) ([]PlayerStat, error) {
	b, usedURL, err := downloadpkg.Get().Download("nflverse-data", path, nil, nil)
	if err != nil {
		return nil, err
	}
	rows, err := downloadpkg.ParseAuto(b, usedURL)
	if err != nil {
		return nil, err
	}
	out := make([]PlayerStat, 0, len(rows))
	for _, r := range rows {
		out = append(out, FromMap(r))
	}
	return out, nil
}
