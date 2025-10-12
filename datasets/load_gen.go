package datasets

import (
	"fmt"

	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

// LoadRaw returns the raw bytes and provenance URL for a dataset key.
func LoadRaw(key Key) ([]byte, string, error) {
	path, ok := pathByKey[key]
	if !ok {
		return nil, "", fmt.Errorf("unknown dataset: %s", key)
	}
	return downloadpkg.Get().Download("nflverse-data", path, nil, nil)
}

// LoadRows returns generic []map[string]any using the downloader's auto-parser.
func LoadRows(key Key) ([]map[string]any, error) {
	b, usedURL, err := LoadRaw(key)
	if err != nil {
		return nil, err
	}
	return downloadpkg.ParseAuto(b, usedURL)
}

// LoadAs provides a typed, generic loader given a mapper function.
func LoadAs[T any](key Key, mapper func(map[string]any) T) ([]T, error) {
	rows, err := LoadRows(key)
	if err != nil {
		return nil, err
	}
	out := make([]T, 0, len(rows))
	for _, r := range rows {
		out = append(out, mapper(r))
	}
	return out, nil
}
