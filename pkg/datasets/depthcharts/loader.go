package depthcharts

import (
	"fmt"

	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

func Load() ([]DepthChart, error) { return loadHelper("depth_charts/depth_charts") }

func LoadSeason(season int) ([]DepthChart, error) {
	if season == 0 {
		return Load()
	}
	return loadHelper(fmt.Sprintf("depth_charts/depth_charts_%d", season))
}

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
