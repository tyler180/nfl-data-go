package players

import (
	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

// Load returns the typed players slice by downloading the canonical
// nflverse dataset (Parquet preferred with CSV fallback) and mapping rows.
func Load() ([]Player, error) {
	b, usedURL, err := downloadpkg.Get().Download("nflverse-data", "players/players", nil, nil)
	if err != nil {
		return nil, err
	}
	rows, err := downloadpkg.ParseAuto(b, usedURL)
	if err != nil {
		return nil, err
	}
	out := make([]Player, 0, len(rows))
	for _, r := range rows {
		out = append(out, FromMap(r))
	}
	return out, nil
}

// LoadRaw exposes the underlying bytes and URL used (useful for persisting
// provenance or re-parsing with custom logic).
func LoadRaw() ([]byte, string, error) {
	return downloadpkg.Get().Download("nflverse-data", "players/players", nil, nil)
}
