// examples/def_snap_pct/main.go
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	snaps "github.com/tyler180/nfl-data-go/pkg/datasets/snapcounts"
)

func main() {
	season := flag.Int("season", 2024, "season year (e.g., 2024)")
	week := flag.Int("week", 0, "week number (1-22). 0 = show all")
	team := flag.String("team", "", "optional team filter (e.g., KC)")
	seasonType := flag.String("season_type", "", "optional: REG|PRE|POST")
	requireDefSnaps := flag.Bool("require_def_snaps", true, "drop rows with 0 defensive snaps")
	out := flag.String("out", "", "optional: write CSV to this path instead of stdout")
	flag.Parse()

	rows, err := snaps.LoadSeason(*season)
	if err != nil {
		log.Fatalf("snapcounts.LoadSeason(%d): %v", *season, err)
	}
	log.Printf("loaded %d snap-count rows for season %d", len(rows), *season)

	// Show what (week, game_type) exist so your filters make sense.
	foundWeeks := map[string]struct{}{}
	for _, r := range rows {
		key := fmt.Sprintf("week=%d game_type=%s", r.Week, r.Gametype)
		foundWeeks[key] = struct{}{}
	}
	log.Printf("available (week, game_type) pairs:")
	for k := range foundWeeks {
		log.Printf("  %s", k)
	}

	// Filter
	type rec struct {
		Season, Week      int
		GameType          string
		Team, Player, Pos string
		PFRID             string
		DefSnaps          int
		DefPct            float64 // if your model has int, cast with float64(...)
	}
	var outRows []rec
	for _, r := range rows {
		if *week != 0 && r.Week != *week {
			continue
		}
		if *team != "" && !strings.EqualFold(r.Team, *team) {
			continue
		}
		if *seasonType != "" && !strings.EqualFold(r.Gametype, *seasonType) {
			continue
		}
		if *requireDefSnaps && r.DefenseSnaps <= 0 {
			continue
		}
		outRows = append(outRows, rec{
			Season:   r.Season,
			Week:     r.Week,
			GameType: r.Gametype,
			Team:     r.Team,
			Player:   r.Player,
			Pos:      r.Position,
			PFRID:    r.PlayerID,
			DefSnaps: r.DefenseSnaps,
			DefPct:   r.DefensePct, // 0–100 as provided in your model
		})
	}
	log.Printf("after filters: %d rows", len(outRows))

	// Sort by team, week, then descending pct
	sort.Slice(outRows, func(i, j int) bool {
		if outRows[i].Team == outRows[j].Team {
			if outRows[i].Week == outRows[j].Week {
				return outRows[i].DefPct > outRows[j].DefPct
			}
			return outRows[i].Week < outRows[j].Week
		}
		return outRows[i].Team < outRows[j].Team
	})

	// Output CSV
	var w *csv.Writer
	if *out == "" {
		w = csv.NewWriter(os.Stdout)
	} else {
		f, err := os.Create(*out)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		w = csv.NewWriter(f)
	}
	defer w.Flush()

	_ = w.Write([]string{"season", "week", "game_type", "team", "player", "pos", "pfr_player_id", "def_snaps", "def_snap_pct"})
	for _, r := range outRows {
		_ = w.Write([]string{
			fmt.Sprint(r.Season),
			fmt.Sprint(r.Week),
			r.GameType,
			r.Team,
			r.Player,
			r.Pos,
			r.PFRID,
			fmt.Sprint(r.DefSnaps),
			fmt.Sprintf("%.1f", r.DefPct),
		})
	}
}
