package depthcharts

import (
	"context"
	"io"

	"github.com/tyler180/nfl-data-go/internal/datasets"
	"github.com/tyler180/nfl-data-go/internal/download"
	"github.com/tyler180/nfl-data-go/internal/source"
)

// NOTE: nflverse path includes "data/..." in the repo.
var src = datasets.Source{Repo: "nflverse/nflverse-data", Base: "data/depth_charts/depth_charts"}

// All seasons (combined)
func Load() ([]DepthChart, error) {
	// datasets.LoadFromSourceAs now takes a context.Context
	return datasets.LoadFromSourceAs[DepthChart](context.TODO(), src, 0, FromMap)
}

// Per-season (e.g., depth_charts_2024.*)
func LoadSeason(season int) ([]DepthChart, error) {
	return datasets.LoadFromSourceAs[DepthChart](context.TODO(), src, season, FromMap)
}

// Raw base asset (all seasons)
// If you prefer a ctx-param version, see LoadRawWithContext below.
func LoadRaw() ([]byte, string, error) {
	return LoadRawWithContext(context.TODO())
}

func LoadRawWithContext(ctx context.Context) ([]byte, string, error) {
	// Build raw.githubusercontent URL
	url := source.RawGitHubURL(src.Repo, src.Base)

	// Create a downloader (add options if you want cache/UA/timeouts)
	dl := download.New()
	rc, _, err := dl.Fetch(ctx, url)
	if err != nil {
		return nil, "", err
	}
	defer rc.Close()

	b, err := io.ReadAll(rc)
	if err != nil {
		return nil, "", err
	}
	return b, url, nil
}
