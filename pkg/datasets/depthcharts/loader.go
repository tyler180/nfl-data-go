package depthcharts

import (
	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

// Load downloads the depth charts dataset and returns typed rows.
// Source/tag: "depth_charts" in nflverse-data.
func Load() ([]DepthChart, error) {
	return loadHelper("depth_charts/depth_charts")
}

// LoadRaw exposes the underlying bytes and provenance URL.
func LoadRaw() ([]byte, string, error) {
	return downloadpkg.Get().Download("nflverse-data", "depth_charts/depth_charts", nil, nil)
}

func loadHelper(path string) ([]DepthChart, error) {
	b, usedURL, err := downloadpkg.Get().Download("nflverse-data", path, nil, nil)
	if err != nil {
		return nil, err
	}
	rows, err := downloadpkg.ParseAuto(b, usedURL)
	if err != nil {
		return nil, err
	}
	out := make([]DepthChart, 0, len(rows))
	for _, r := range rows {
		out = append(out, FromMap(r))
	}
	return out, nil
}
