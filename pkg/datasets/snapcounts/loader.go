package snapcounts

import (
	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

// Load downloads the canonical nflverse snap counts dataset and returns
// a typed slice. Parquet is preferred with CSV fallback via ParseAuto.
func Load() ([]SnapCount, error) {
	b, usedURL, err := downloadpkg.Get().Download("nflverse-data", "snap_counts/snap_counts_2024", nil, nil)
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

// LoadRaw exposes the underlying bytes and URL used for snap counts.
func LoadRaw() ([]byte, string, error) {
	return downloadpkg.Get().Download("nflverse-data", "snap_counts/snap_counts", nil, nil)
}
