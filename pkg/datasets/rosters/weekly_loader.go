package rosters

import (
	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

// LoadWeekly downloads the week-level rosters dataset and returns typed rows.
// Source/tag: "weekly_rosters" in nflverse-data releases.
// Asset base name is the same as the tag (e.g., weekly_rosters.parquet/csv).
func LoadWeekly() ([]Roster, error) {
	b, usedURL, err := downloadpkg.Get().Download("nflverse-data", "weekly_rosters/weekly_rosters", nil, nil)
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

// LoadWeeklyRaw exposes the underlying bytes and provenance URL for the
// week-level rosters dataset.
func LoadWeeklyRaw() ([]byte, string, error) {
	return downloadpkg.Get().Download("nflverse-data", "weekly_rosters/weekly_rosters", nil, nil)
}
