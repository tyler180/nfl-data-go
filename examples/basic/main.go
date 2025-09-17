package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	configpkg "github.com/tyler180/nfl-data-go/internal/config"
	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
	datasets "github.com/tyler180/nfl-data-go/pkg/datasets"
	ffpkg "github.com/tyler180/nfl-data-go/pkg/datasets/ffplayerids"
	inj "github.com/tyler180/nfl-data-go/pkg/datasets/injuries"
	playerpkg "github.com/tyler180/nfl-data-go/pkg/datasets/players"
	pstatpkg "github.com/tyler180/nfl-data-go/pkg/datasets/playerstats"
	snappkg "github.com/tyler180/nfl-data-go/pkg/datasets/snapcounts"
)

func main() {
	var (
		dataset = flag.String("dataset", "players", "dataset to load: players|snapcounts|playerstats")
		limit   = flag.Int("limit", 3, "how many rows to print")
		format  = flag.String("format", "", "prefer format: parquet|csv (optional)")
		verbose = flag.Bool("v", true, "verbose HTTP/caching logs")
	)
	flag.Parse()

	// Configure the library at runtime
	opts := []configpkg.ConfigOption{configpkg.WithVerbose(*verbose)}
	switch *format {
	case "csv":
		opts = append(opts, configpkg.WithPreferFormat(downloadpkg.FormatCSV))
	case "parquet":
		opts = append(opts, configpkg.WithPreferFormat(downloadpkg.FormatParquet))
	}
	configpkg.UpdateConfig(opts...)

	switch *dataset {
	case "players":
		rows, err := datasets.LoadAs(datasets.Players, playerpkg.FromMap)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("players: %d rows\n", len(rows))
		printJSONRows(rowsToAny(rows, *limit))
	case "snapcounts":
		rows, err := datasets.LoadAs(datasets.SnapCounts, snappkg.FromMap)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("snapcounts: %d rows\n", len(rows))
		printJSONRows(rowsToAny(rows, *limit))
	case "playerstats":
		rows, err := datasets.LoadAs(datasets.PlayerStats, pstatpkg.FromMap)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("playerstats: %d rows\n", len(rows))
		printJSONRows(rowsToAny(rows, *limit))
	case "ff_playerids":
		rows, err := ffpkg.Load()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ff_playerids: %d rows\n", len(rows))
		printJSONRows(rowsToAny(rows, *limit))
	case "injuries":
		rows, err := inj.Load()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("injuries: %d rows\n", len(rows))
		printJSONRows(rowsToAny(rows, *limit))
	default:
		log.Fatalf("unknown dataset: %s (use players|snapcounts|playerstats)", *dataset)
	}
}

func rowsToAny[T any](rows []T, limit int) []any {
	if limit < 0 || limit > len(rows) {
		limit = len(rows)
	}
	out := make([]any, 0, limit)
	for i := 0; i < limit; i++ {
		out = append(out, rows[i])
	}
	return out
}

func printJSONRows(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	if err := enc.Encode(v); err != nil {
		log.Fatal(err)
	}
}
