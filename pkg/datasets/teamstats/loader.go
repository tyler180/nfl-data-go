package teamstats

import (
	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

// Load downloads the WEEK-level team stats dataset and returns typed rows.
// Mirrors nflreadr::load_team_stats(summary_level = "week").
func Load() ([]TeamStat, error) {
	return loadHelper("stats_team/stats_team_week")
}

// LoadRaw returns the raw bytes and provenance URL for the WEEK-level dataset.
func LoadRaw() ([]byte, string, error) {
	return downloadpkg.Get().Download("nflverse-data", "stats_team/stats_team_week", nil, nil)
}

// Season summary helpers (optional). These correspond to REG, POST, and REG+POST.
func LoadSeasonReg() ([]TeamStat, error)     { return loadHelper("stats_team/stats_team_reg") }
func LoadSeasonPost() ([]TeamStat, error)    { return loadHelper("stats_team/stats_team_post") }
func LoadSeasonRegPost() ([]TeamStat, error) { return loadHelper("stats_team/stats_team_reg_post") }

// loadHelper fetches, auto-parses (Parquet/CSV), and maps rows -> TeamStat.
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
