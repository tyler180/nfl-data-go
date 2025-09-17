package injuries

import (
	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

// Load downloads the injuries dataset and returns typed rows.
// Mirrors nflreadr::load_injuries(); pulls from the nflverse-data release.
func Load() ([]Injury, error) {
	return loadHelper("injuries/injuries")
}

// LoadRaw returns the raw bytes and provenance URL for the injuries dataset.
func LoadRaw() ([]byte, string, error) {
	return downloadpkg.Get().Download("nflverse-data", "injuries/injuries", nil, nil)
}

func loadHelper(path string) ([]Injury, error) {
	b, usedURL, err := downloadpkg.Get().Download("nflverse-data", path, nil, nil)
	if err != nil {
		return nil, err
	}
	rows, err := downloadpkg.ParseAuto(b, usedURL)
	if err != nil {
		return nil, err
	}
	out := make([]Injury, 0, len(rows))
	for _, r := range rows {
		out = append(out, FromMap(r))
	}
	return out, nil
}
