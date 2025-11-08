package datasets

import (
	"context"
	"fmt"
	"io"

	"github.com/tyler180/nfl-data-go/internal/download"
	"github.com/tyler180/nfl-data-go/internal/parse"
	"github.com/tyler180/nfl-data-go/internal/source"
)

// LoadRaw returns the raw bytes and provenance URL for a dataset key.
func LoadRaw(ctx context.Context, key Key) ([]byte, string, error) {
	path, ok := pathByKey[key]
	if !ok {
		return nil, "", fmt.Errorf("unknown dataset: %s", key)
	}

	// Build raw.githubusercontent URL (owner defaults to nflverse if not provided)
	url := source.RawGitHubURL("nflverse/nflverse-data", path)

	// Create a downloader (you can add options: WithCache, WithUserAgent, etc.)
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

// LoadRows returns generic []map[string]any using the parser's auto-detection (CSV now; Parquet TODO).
func LoadRows(ctx context.Context, key Key) ([]map[string]any, error) {
	b, usedURL, err := LoadRaw(ctx, key)
	if err != nil {
		return nil, err
	}
	return parse.Auto(b, usedURL)
}

// LoadAs provides a typed, generic loader given a mapper function.
func LoadAs[T any](ctx context.Context, key Key, mapper func(map[string]any) T) ([]T, error) {
	rows, err := LoadRows(ctx, key)
	if err != nil {
		return nil, err
	}
	out := make([]T, 0, len(rows))
	for _, r := range rows {
		out = append(out, mapper(r))
	}
	return out, nil
}
