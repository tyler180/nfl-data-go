package teamstats

import (
	"fmt"

	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

// WEEK-level (default)
func Load() ([]TeamStat, error) { return loadHelper("stats_team/stats_team_week") }

// WEEK-level for a single season (e.g., stats_team_week_2024)
func LoadForSeason(season int) ([]TeamStat, error) {
	if season == 0 {
		return Load()
	}
	return loadHelper(fmt.Sprintf("stats_team/stats_team_week_%d", season))
}

// Optional season summary helpers
func LoadSeasonReg() ([]TeamStat, error)     { return loadHelper("stats_team/stats_team_reg") }
func LoadSeasonPost() ([]TeamStat, error)    { return loadHelper("stats_team/stats_team_post") }
func LoadSeasonRegPost() ([]TeamStat, error) { return loadHelper("stats_team/stats_team_reg_post") }

func LoadSeasonRegForSeason(season int) ([]TeamStat, error) {
	return loadHelper(fmt.Sprintf("stats_team/stats_team_reg_%d", season))
}
func LoadSeasonPostForSeason(season int) ([]TeamStat, error) {
	return loadHelper(fmt.Sprintf("stats_team/stats_team_post_%d", season))
}
func LoadSeasonRegPostForSeason(season int) ([]TeamStat, error) {
	return loadHelper(fmt.Sprintf("stats_team/stats_team_reg_post_%d", season))
}

func LoadRaw() ([]byte, string, error) {
	return downloadpkg.Get().Download("nflverse-data", "stats_team/stats_team_week", nil, nil)
}

func loadHelper(path string) ([]TeamStat, error) {
	b, usedURL, err := downloadpkg.Get().Download("nflverse-data", path, nil, nil)
	if err != nil {
		return nil, err
	}
	rows, err := downloadpkg.ParseAuto(b, usedURL)
	if err != nil {
		return nil, err
	}
	out := make([]TeamStat, 0, len(rows))
	for _, r := range rows {
		out = append(out, FromMap(r))
	}
	return out, nil
}
