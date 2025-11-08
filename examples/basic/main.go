package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	configpkg "github.com/tyler180/nfl-data-go/internal/config"
	"github.com/tyler180/nfl-data-go/internal/datasets"
	dchartpkg "github.com/tyler180/nfl-data-go/internal/datasets/depthcharts"
	ffpidpkg "github.com/tyler180/nfl-data-go/internal/datasets/ffplayerids"
	injpkg "github.com/tyler180/nfl-data-go/internal/datasets/injuries"
	playerpkg "github.com/tyler180/nfl-data-go/internal/datasets/players"
	pstatpkg "github.com/tyler180/nfl-data-go/internal/datasets/playerstats"
	rosterpkg "github.com/tyler180/nfl-data-go/internal/datasets/rosters"
	snappkg "github.com/tyler180/nfl-data-go/internal/datasets/snapcounts"
	tstatpkg "github.com/tyler180/nfl-data-go/internal/datasets/teamstats"
	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

func main() {
	var (
		dataset    = flag.String("dataset", "players", "dataset: players|snapcounts|playerstats|rosters|rosters_weekly|teamstats|depth_charts|injuries|ff_playerids")
		limit      = flag.Int("limit", 3, "how many rows to print")
		format     = flag.String("format", "", "prefer format: parquet|csv (optional)")
		verbose    = flag.Bool("v", true, "verbose HTTP/caching logs")
		season     = flag.Int("season", 0, "download a specific season file when available (e.g., 2023). 0 = all seasons (if available)")
		week       = flag.Int("week", 0, "filter to a specific week (1-22). 0 = no filter")
		seasonType = flag.String("season_type", "", "filter by season type: REG|POST (optional)")
	)
	flag.Parse()

	ctx := context.Background()
	// Configure the library at runtime
	opts := []configpkg.ConfigOption{configpkg.WithVerbose(*verbose)}
	switch strings.ToLower(*format) {
	case "csv":
		opts = append(opts, configpkg.WithPreferFormat(downloadpkg.FormatCSV))
	case "parquet":
		opts = append(opts, configpkg.WithPreferFormat(downloadpkg.FormatParquet))
	}
	configpkg.UpdateConfig(opts...)

	switch *dataset {
	case "players":
		rows, err := datasets.LoadAs[playerpkg.Player](ctx, datasets.Players, playerpkg.FromMap)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("players: %d rows\n", len(rows))
		printJSONRows(rowsToAny(rows, *limit))

	case "snapcounts":
		var (
			rows []snappkg.SnapCount
			err  error
		)
		if *season != 0 {
			rows, err = snappkg.LoadSeason(ctx, *season)
		} else {
			rows, err = snappkg.Load(ctx)
		}
		if err != nil {
			log.Fatal(err)
		}
		rows = filter(rows, func(r snappkg.SnapCount) bool {
			if *week != 0 && r.Week != *week {
				return false
			}
			return true
		})
		fmt.Printf("snapcounts: %d rows (after filters)\n", len(rows))
		printJSONRows(rowsToAny(rows, *limit))

	case "playerstats":
		var (
			rows []pstatpkg.PlayerStat
			err  error
		)
		if *season != 0 {
			rows, err = pstatpkg.LoadForSeason(ctx, *season) // week-level per-season file
		} else {
			rows, err = pstatpkg.Load(ctx) // all seasons (week-level)
		}
		if err != nil {
			log.Fatal(err)
		}
		rows = filter(rows, func(r pstatpkg.PlayerStat) bool {
			if *week != 0 && r.Week != *week {
				return false
			}
			if *seasonType != "" && !strings.EqualFold(r.SeasonType, *seasonType) {
				return false
			}
			return true
		})
		fmt.Printf("playerstats: %d rows (after filters)\n", len(rows))
		printJSONRows(rowsToAny(rows, *limit))

	case "teamstats":
		var (
			rows []tstatpkg.TeamStat
			err  error
		)
		if *season != 0 {
			rows, err = tstatpkg.LoadForSeason(ctx, *season) // week-level per-season file
		} else {
			rows, err = tstatpkg.Load(ctx) // all seasons (week-level)
		}
		if err != nil {
			log.Fatal(err)
		}
		rows = filter(rows, func(r tstatpkg.TeamStat) bool {
			if *week != 0 && r.Week != *week {
				return false
			}
			if *seasonType != "" && !strings.EqualFold(r.SeasonType, *seasonType) {
				return false
			}
			return true
		})
		fmt.Printf("teamstats (week): %d rows (after filters)\n", len(rows))
		printJSONRows(rowsToAny(rows, *limit))

	case "rosters":
		var (
			rows []rosterpkg.Roster
			err  error
		)
		if *season != 0 {
			rows, err = rosterpkg.LoadSeason(ctx, *season)
		} else {
			rows, err = rosterpkg.Load(ctx)
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("rosters: %d rows\n", len(rows))
		printJSONRows(rowsToAny(rows, *limit))

	case "rosters_weekly":
		var (
			rows []rosterpkg.Roster
			err  error
		)
		if *season != 0 {
			rows, err = rosterpkg.LoadWeeklySeason(ctx, *season)
		} else {
			rows, err = rosterpkg.LoadWeekly(ctx)
		}
		if err != nil {
			log.Fatal(err)
		}
		rows = filter(rows, func(r rosterpkg.Roster) bool {
			if *week != 0 && r.Week != *week {
				return false
			}
			return true
		})
		fmt.Printf("rosters_weekly: %d rows (after filters)\n", len(rows))
		printJSONRows(rowsToAny(rows, *limit))

	case "depth_charts":
		var (
			rows []dchartpkg.DepthChart
			err  error
		)
		if *season != 0 {
			rows, err = dchartpkg.LoadSeason(*season)
		} else {
			rows, err = dchartpkg.Load()
		}
		if err != nil {
			log.Fatal(err)
		}
		rows = filter(rows, func(r dchartpkg.DepthChart) bool {
			if *week != 0 && r.Week != *week {
				return false
			}
			return true
		})
		fmt.Printf("depth_charts: %d rows (after filters)\n", len(rows))
		printJSONRows(rowsToAny(rows, *limit))

	case "injuries":
		var (
			rows []injpkg.Injury
			err  error
		)
		if *season != 0 {
			rows, err = injpkg.LoadSeason(ctx, *season)
		} else {
			rows, err = injpkg.Load(ctx)
		}
		if err != nil {
			log.Fatal(err)
		}
		rows = filter(rows, func(r injpkg.Injury) bool {
			if *week != 0 && r.Week != *week {
				return false
			}
			if *seasonType != "" && !strings.EqualFold(r.SeasonType, *seasonType) {
				return false
			}
			return true
		})
		fmt.Printf("injuries: %d rows (after filters)\n", len(rows))
		printJSONRows(rowsToAny(rows, *limit))

	case "ff_playerids":
		rows, err := ffpidpkg.Load(ctx)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ff_playerids: %d rows\n", len(rows))
		printJSONRows(rowsToAny(rows, *limit))

	case "players_components":

	default:
		log.Fatalf("unknown dataset: %s (use players|snapcounts|playerstats|rosters|rosters_weekly|teamstats|depth_charts|injuries|ff_playerids)", *dataset)
	}
}

// filter keeps rows for which keep(x) == true.
func filter[T any](rows []T, keep func(T) bool) []T {
	out := rows[:0]
	for _, r := range rows {
		if keep(r) {
			out = append(out, r)
		}
	}
	return out
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
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		log.Fatal(err)
	}
}

// package main

// import (
// 	"encoding/json"
// 	"flag"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strings"

// 	configpkg "github.com/tyler180/nfl-data-go/internal/config"
// 	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
// 	"github.com/tyler180/nfl-data-go/internal/datasets"
// 	dchartpkg "github.com/tyler180/nfl-data-go/internal/datasets/depthcharts"
// 	ffpidpkg "github.com/tyler180/nfl-data-go/internal/datasets/ffplayerids"
// 	injpkg "github.com/tyler180/nfl-data-go/internal/datasets/injuries"
// 	playerpkg "github.com/tyler180/nfl-data-go/internal/datasets/players"
// 	pstatpkg "github.com/tyler180/nfl-data-go/internal/datasets/playerstats"
// 	rosterpkg "github.com/tyler180/nfl-data-go/internal/datasets/rosters"
// 	snappkg "github.com/tyler180/nfl-data-go/internal/datasets/snapcounts"
// 	tstatpkg "github.com/tyler180/nfl-data-go/internal/datasets/teamstats"
// )

// func main() {
// 	var (
// 		dataset    = flag.String("dataset", "players", "dataset: players|snapcounts|playerstats|rosters|rosters_weekly|teamstats|depth_charts|injuries|ff_playerids")
// 		limit      = flag.Int("limit", 3, "how many rows to print")
// 		format     = flag.String("format", "", "prefer format: parquet|csv (optional)")
// 		verbose    = flag.Bool("v", true, "verbose HTTP/caching logs")
// 		season     = flag.Int("season", 0, "filter to a specific season (e.g., 2023). 0 = no filter")
// 		week       = flag.Int("week", 0, "filter to a specific week (1-22). 0 = no filter")
// 		seasonType = flag.String("season_type", "", "filter by season type: REG|POST (optional)")
// 	)
// 	flag.Parse()

// 	// Configure the library at runtime
// 	opts := []configpkg.ConfigOption{configpkg.WithVerbose(*verbose)}
// 	if *format == "csv" {
// 		opts = append(opts, configpkg.WithPreferFormat(downloadpkg.FormatCSV))
// 	} else if *format == "parquet" {
// 		opts = append(opts, configpkg.WithPreferFormat(downloadpkg.FormatParquet))
// 	}
// 	configpkg.UpdateConfig(opts...)

// 	switch *dataset {
// 	case "players":
// 		rows, err := datasets.LoadAs[playerpkg.Player](datasets.Players, playerpkg.FromMap)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("players: %d rows\n", len(rows))
// 		printJSONRows(rowsToAny(filter(rows, func(_ playerpkg.Player) bool { return true }), *limit))

// 	case "snapcounts":
// 		rows, err := datasets.LoadAs[snappkg.SnapCount](datasets.SnapCounts, snappkg.FromMap)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		rows = filter(rows, func(r snappkg.SnapCount) bool {
// 			if *season != 0 && r.Season != *season {
// 				return false
// 			}
// 			if *week != 0 && r.Week != *week {
// 				return false
// 			}
// 			return true
// 		})
// 		fmt.Printf("snapcounts: %d rows (filtered)\n", len(rows))
// 		printJSONRows(rowsToAny(rows, *limit))

// 	case "playerstats":
// 		rows, err := datasets.LoadAs[pstatpkg.PlayerStat](datasets.PlayerStats, pstatpkg.FromMap)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		rows = filter(rows, func(r pstatpkg.PlayerStat) bool {
// 			if *season != 0 && r.Season != *season {
// 				return false
// 			}
// 			if *week != 0 && r.Week != *week {
// 				return false
// 			}
// 			if *seasonType != "" && !strings.EqualFold(r.SeasonType, *seasonType) {
// 				return false
// 			}
// 			return true
// 		})
// 		fmt.Printf("playerstats: %d rows (filtered)\n", len(rows))
// 		printJSONRows(rowsToAny(rows, *limit))

// 	case "rosters":
// 		rows, err := rosterpkg.Load()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		rows = filter(rows, func(r rosterpkg.Roster) bool {
// 			if *season != 0 && r.Season != *season {
// 				return false
// 			}
// 			return true
// 		})
// 		fmt.Printf("rosters: %d rows (filtered)\n", len(rows))
// 		printJSONRows(rowsToAny(rows, *limit))

// 	case "rosters_weekly":
// 		rows, err := rosterpkg.LoadWeekly()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		rows = filter(rows, func(r rosterpkg.Roster) bool {
// 			if *season != 0 && r.Season != *season {
// 				return false
// 			}
// 			if *week != 0 && r.Week != *week {
// 				return false
// 			}
// 			return true
// 		})
// 		fmt.Printf("rosters_weekly: %d rows (filtered)\n", len(rows))
// 		printJSONRows(rowsToAny(rows, *limit))

// 	case "teamstats":
// 		rows, err := tstatpkg.Load()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		rows = filter(rows, func(r tstatpkg.TeamStat) bool {
// 			if *season != 0 && r.Season != *season {
// 				return false
// 			}
// 			if *week != 0 && r.Week != *week {
// 				return false
// 			}
// 			if *seasonType != "" && !strings.EqualFold(r.SeasonType, *seasonType) {
// 				return false
// 			}
// 			return true
// 		})
// 		fmt.Printf("teamstats (week): %d rows (filtered)\n", len(rows))
// 		printJSONRows(rowsToAny(rows, *limit))

// 	case "depth_charts":
// 		rows, err := dchartpkg.Load()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		rows = filter(rows, func(r dchartpkg.DepthChart) bool {
// 			if *season != 0 && r.Season != *season {
// 				return false
// 			}
// 			if *week != 0 && r.Week != *week {
// 				return false
// 			}
// 			return true
// 		})
// 		fmt.Printf("depth_charts: %d rows (filtered)\n", len(rows))
// 		printJSONRows(rowsToAny(rows, *limit))

// 	case "injuries":
// 		rows, err := datasets.LoadAs[injpkg.Injury](datasets.Injuries, injpkg.FromMap)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		rows = filter(rows, func(r injpkg.Injury) bool {
// 			if *season != 0 && r.Season != *season {
// 				return false
// 			}
// 			if *week != 0 && r.Week != *week {
// 				return false
// 			}
// 			if *seasonType != "" && !strings.EqualFold(r.SeasonType, *seasonType) {
// 				return false
// 			}
// 			return true
// 		})
// 		fmt.Printf("injuries: %d rows (filtered)\n", len(rows))
// 		printJSONRows(rowsToAny(rows, *limit))

// 	case "ff_playerids":
// 		rows, err := ffpidpkg.Load()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		// ff_playerids has db_season, not season/week; we skip filtering here.
// 		fmt.Printf("ff_playerids: %d rows\n", len(rows))
// 		printJSONRows(rowsToAny(rows, *limit))

// 	default:
// 		log.Fatalf("unknown dataset: %s (use players|snapcounts|playerstats|rosters|rosters_weekly|teamstats|depth_charts|injuries|ff_playerids)", *dataset)
// 	}
// }

// // filter keeps rows for which keep(x) == true.
// func filter[T any](rows []T, keep func(T) bool) []T {
// 	out := rows[:0]
// 	for _, r := range rows {
// 		if keep(r) {
// 			out = append(out, r)
// 		}
// 	}
// 	return out
// }

// func rowsToAny[T any](rows []T, limit int) []any {
// 	if limit < 0 || limit > len(rows) {
// 		limit = len(rows)
// 	}
// 	out := make([]any, 0, limit)
// 	for i := 0; i < limit; i++ {
// 		out = append(out, rows[i])
// 	}
// 	return out
// }

// func printJSONRows(v any) {
// 	enc := json.NewEncoder(os.Stdout)
// 	enc.SetIndent("", "  ")
// 	if err := enc.Encode(v); err != nil {
// 		log.Fatal(err)
// 	}
// }

// package main

// import (
// 	"encoding/json"
// 	"flag"
// 	"fmt"
// 	"log"
// 	"os"

// 	configpkg "github.com/tyler180/nfl-data-go/internal/config"
// 	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
// 	datasets "github.com/tyler180/nfl-data-go/internal/datasets"
// 	ffpkg "github.com/tyler180/nfl-data-go/internal/datasets/ffplayerids"
// 	inj "github.com/tyler180/nfl-data-go/internal/datasets/injuries"
// 	playerpkg "github.com/tyler180/nfl-data-go/internal/datasets/players"
// 	pstatpkg "github.com/tyler180/nfl-data-go/internal/datasets/playerstats"
// 	snappkg "github.com/tyler180/nfl-data-go/internal/datasets/snapcounts"
// )

// func main() {
// 	var (
// 		dataset = flag.String("dataset", "players", "dataset to load: players|snapcounts|playerstats")
// 		limit   = flag.Int("limit", 3, "how many rows to print")
// 		format  = flag.String("format", "", "prefer format: parquet|csv (optional)")
// 		verbose = flag.Bool("v", true, "verbose HTTP/caching logs")
// 	)
// 	flag.Parse()

// 	// Configure the library at runtime
// 	opts := []configpkg.ConfigOption{configpkg.WithVerbose(*verbose)}
// 	switch *format {
// 	case "csv":
// 		opts = append(opts, configpkg.WithPreferFormat(downloadpkg.FormatCSV))
// 	case "parquet":
// 		opts = append(opts, configpkg.WithPreferFormat(downloadpkg.FormatParquet))
// 	}
// 	configpkg.UpdateConfig(opts...)

// 	switch *dataset {
// 	case "players":
// 		rows, err := datasets.LoadAs(datasets.Players, playerpkg.FromMap)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("players: %d rows\n", len(rows))
// 		printJSONRows(rowsToAny(rows, *limit))
// 	case "snapcounts":
// 		rows, err := datasets.LoadAs(datasets.SnapCounts, snappkg.FromMap)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("snapcounts: %d rows\n", len(rows))
// 		printJSONRows(rowsToAny(rows, *limit))
// 	case "playerstats":
// 		rows, err := datasets.LoadAs(datasets.PlayerStats, pstatpkg.FromMap)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("playerstats: %d rows\n", len(rows))
// 		printJSONRows(rowsToAny(rows, *limit))
// 	case "ff_playerids":
// 		rows, err := ffpkg.Load()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("ff_playerids: %d rows\n", len(rows))
// 		printJSONRows(rowsToAny(rows, *limit))
// 	case "injuries":
// 		rows, err := inj.Load()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("injuries: %d rows\n", len(rows))
// 		printJSONRows(rowsToAny(rows, *limit))
// 	default:
// 		log.Fatalf("unknown dataset: %s (use players|snapcounts|playerstats)", *dataset)
// 	}
// }

// func rowsToAny[T any](rows []T, limit int) []any {
// 	if limit < 0 || limit > len(rows) {
// 		limit = len(rows)
// 	}
// 	out := make([]any, 0, limit)
// 	for i := 0; i < limit; i++ {
// 		out = append(out, rows[i])
// 	}
// 	return out
// }

// func printJSONRows(v any) {
// 	enc := json.NewEncoder(os.Stdout)
// 	enc.SetIndent("", " ")
// 	if err := enc.Encode(v); err != nil {
// 		log.Fatal(err)
// 	}
// }
