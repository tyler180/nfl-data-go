package playerstats

import (
	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

func Load() ([]PlayerStat, error) {
	b, usedURL, err := downloadpkg.Get().Download("nflverse-data", "player_stats/player_stats_2024", nil, nil)
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

func LoadRaw() ([]byte, string, error) {
	return downloadpkg.Get().Download("nflverse-data", "player_stats/player_stats_2024", nil, nil)
}
