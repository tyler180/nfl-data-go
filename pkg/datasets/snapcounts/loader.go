package snapcounts

import (
	"fmt"

	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

func Load() ([]SnapCount, error) { return loadHelper("snap_counts/snap_counts") }

func LoadSeason(season int) ([]SnapCount, error) {
	if season == 0 {
		return Load()
	}
	return loadHelper(fmt.Sprintf("snap_counts/snap_counts_%d", season))
}

func loadHelper(path string) ([]SnapCount, error) {
	fmt.Printf("Loading snap counts from path: %s\n", path)
	b, usedURL, err := downloadpkg.Get().Download("nflverse-data", path, nil, nil)
	if err != nil {
		return nil, err
	}
	rows, err := downloadpkg.ParseAuto(b, usedURL)
	if err != nil {
		return nil, err
	}
	out := make([]SnapCount, 0, len(rows))
	for _, r := range rows {
		out = append(out, FromMap(r))
	}
	return out, nil
}
