package datasets

// import (
// 	"fmt"

// 	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
// 	playerpkg "github.com/tyler180/nfl-data-go/pkg/datasets/players"
// 	snappkg "github.com/tyler180/nfl-data-go/pkg/datasets/snapcounts"
// )

// // Load provides a single entry point that returns a concrete typed slice for
// // known datasets (players, snapcounts) and falls back to generic rows for
// // unknown keys. The return type is 'any' to allow different concrete slices.
// //
// // Example:
// //
// //	v, err := datasets.Load(datasets.Players)
// //	ps := v.([]players.Player)
// func Load(key Key) (any, error) {
// 	switch key {
// 	case Players:
// 		return playerpkg.Load()
// 	case SnapCounts:
// 		return snappkg.Load()
// 	default:
// 		rows, err := LoadRows(key)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return rows, nil // []map[string]any for untyped datasets
// 	}
// }

// // LoadRaw returns the raw bytes and provenance URL for a dataset key.
// func LoadRaw(key Key) ([]byte, string, error) {
// 	path, ok := pathByKey[key]
// 	if !ok {
// 		return nil, "", fmt.Errorf("unknown dataset: %s", key)
// 	}
// 	return downloadpkg.Get().Download("nflverse-data", path, nil, nil)
// }

// // LoadRows returns generic rows for a dataset key using the downloader's
// // CSV/Parquet auto-parser.
// func LoadRows(key Key) ([]map[string]any, error) {
// 	b, usedURL, err := LoadRaw(key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return downloadpkg.ParseAuto(b, usedURL)
// }

// // LoadAs provides a typed, generic loader given a mapper function.
// // Callers supply a dataset key and a row->T mapper.
// //
// // Example:
// //
// //	players, _ := datasets.LoadAs[players.Player](datasets.Players, players.FromMap)
// func LoadAs[T any](key Key, mapper func(map[string]any) T) ([]T, error) {
// 	rows, err := LoadRows(key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	out := make([]T, 0, len(rows))
// 	for _, r := range rows {
// 		out = append(out, mapper(r))
// 	}
// 	return out, nil
// }
