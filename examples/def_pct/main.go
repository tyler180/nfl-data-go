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
	season := flag.Int("season", 2024, "season year")
	week := flag.Int("week", 1, "week (1-22)")
	team := flag.String("team", "", "optional team filter (e.g., KC)")
	flag.Parse()

	rows, err := snaps.LoadSeason(*season) // season-aware loader
	if err != nil {
		log.Fatal(err)
	}

	// keep only this week, defenders (had >0 defensive snaps), and optional team
	type rec struct {
		Team, Player, Pos, ID string
		DefSnaps              int
		DefPct                float64 // 0–100 as provided
	}
	var out []rec
	for _, r := range rows {
		if r.Week != *week {
			continue
		}
		if *team != "" && !strings.EqualFold(r.Team, *team) {
			continue
		}
		if r.DefenseSnaps <= 0 {
			continue
		}
		out = append(out, rec{
			Team: r.Team, Player: r.Player, Pos: r.Position, ID: r.PlayerID,
			DefSnaps: r.DefenseSnaps, DefPct: r.DefensePct,
		})
	}

	// sort by team then descending % (or just by % if you prefer)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Team == out[j].Team {
			return out[i].DefPct > out[j].DefPct
		}
		return out[i].Team < out[j].Team
	})

	// write CSV to stdout
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()
	_ = w.Write([]string{"team", "player", "pos", "pfr_player_id", "def_snaps", "def_snap_pct"})
	for _, r := range out {
		_ = w.Write([]string{
			r.Team, r.Player, r.Pos, r.ID,
			fmt.Sprint(r.DefSnaps),
			fmt.Sprintf("%.1f", r.DefPct), // already 0–100
		})
	}
}
