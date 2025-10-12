package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"

	snaps "github.com/tyler180/nfl-data-go/datasets/snapcounts"
)

func main() {
	season := flag.Int("season", 2024, "season year (e.g., 2024)")
	week := flag.Int("week", 1, "week number (1-22)")
	seasonType := flag.String("season_type", "", "optional: REG|PRE|POST")
	pfrID := flag.String("pfr_id", "", "exact PFR player id (preferred)")
	name := flag.String("name", "", "player full name (case-insensitive)")
	team := flag.String("team", "", "optional team to disambiguate (e.g., KC)")
	flag.Parse()

	if *week <= 0 {
		log.Fatal("week must be > 0")
	}
	if *pfrID == "" && *name == "" {
		log.Fatal("provide either -pfr_id or -name (and optionally -team)")
	}

	// Load per-season file for performance.
	rows, err := snaps.LoadSeason(*season)
	if err != nil {
		log.Fatalf("snapcounts.LoadSeason(%d): %v", *season, err)
	}

	// Pick the best match.
	var hits []snaps.SnapCount
	for _, r := range rows {
		if r.Week != *week {
			continue
		}
		if *seasonType != "" && !strings.EqualFold(r.Gametype, *seasonType) {
			continue
		}

		match := false
		switch {
		case *pfrID != "":
			match = r.PlayerID == *pfrID
		case *name != "":
			if strings.EqualFold(strings.TrimSpace(r.Player), strings.TrimSpace(*name)) {
				match = true
				if *team != "" && !strings.EqualFold(r.Team, *team) {
					match = false
				}
			}
		}
		if match {
			hits = append(hits, r)
		}
	}

	if len(hits) == 0 {
		log.Printf("no match for season=%d week=%d season_type=%q pfr_id=%q name=%q team=%q",
			*season, *week, *seasonType, *pfrID, *name, *team)
		os.Exit(0)
	}

	// If multiple rows match (rare: multi-team weeks or dup names), keep them all.
	type out struct {
		Season      int     `json:"season"`
		Week        int     `json:"week"`
		GameType    string  `json:"game_type"`
		Team        string  `json:"team"`
		Player      string  `json:"player"`
		PFRPlayerID string  `json:"pfr_player_id"`
		Position    string  `json:"position"`
		DefSnaps    int     `json:"defensive_snaps"`
		DefSnapPct  float64 `json:"defense_pct"` // 0–100
		OffSnaps    int     `json:"offensive_snaps"`
		OffSnapPct  float64 `json:"offense_pct"` // 0–100
		StSnaps     int     `json:"st_snaps"`
		StSnapPct   float64 `json:"st_pct"` // 0–100
	}
	var result []out
	for _, r := range hits {
		result = append(result, out{
			Season:      r.Season,
			Week:        r.Week,
			GameType:    r.Gametype,
			Team:        r.Team,
			Player:      r.Player,
			PFRPlayerID: r.PlayerID,
			Position:    r.Position,
			DefSnaps:    r.DefenseSnaps,
			DefSnapPct:  r.DefensePct,
			OffSnaps:    r.OffenseSnaps,
			OffSnapPct:  r.OffensePct,
			StSnaps:     r.SpecialTeamsSnaps,
			StSnapPct:   r.SpecialTeamsPct,
		})
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		log.Fatal(err)
	}
}
