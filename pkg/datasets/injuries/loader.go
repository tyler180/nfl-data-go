package injuries

import (
	"fmt"

	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

// All seasons (combined file)
func Load() ([]Injury, error) { return loadHelper("injuries/injuries") }

// Single season (e.g., injuries_2024.parquet)
func LoadSeason(season int) ([]Injury, error) {
	if season == 0 {
		return Load()
	}
	return loadHelper(fmt.Sprintf("injuries/injuries_%d", season))
}

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
