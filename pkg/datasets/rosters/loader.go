package rosters

import (
	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

// Load downloads the season-level rosters dataset and returns typed rows.
func Load() ([]Roster, error) {
	b, usedURL, err := downloadpkg.Get().Download("nflverse-data", "rosters/rosters", nil, nil)
	if err != nil {
		return nil, err
	}
	rows, err := downloadpkg.ParseAuto(b, usedURL)
	if err != nil {
		return nil, err
	}
	out := make([]Roster, 0, len(rows))
	for _, r := range rows {
		out = append(out, FromMap(r))
	}
	return out, nil
}

// LoadRaw exposes the underlying bytes and provenance URL.
func LoadRaw() ([]byte, string, error) {
	return downloadpkg.Get().Download("nflverse-data", "rosters/rosters", nil, nil)
}
