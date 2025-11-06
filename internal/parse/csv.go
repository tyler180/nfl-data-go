package parse

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/tyler180/nfl-data-go/internal/schema"
)

// SnapCountsCSV parses a CSV stream into []nflreadgo.SnapCount.
//
// Expected columns (case/underscore-insensitive; extra columns ignored):
//
//	season, week, game_id, player_id|gsis_id, team,
//	offense_snaps|team_snaps, player_snaps|snaps, snap_pct (optional)
//
// If snap_pct is missing, it is computed as 100 * player_snaps / offense_snaps (when both present).
func SnapCountsCSV(r io.Reader) ([]schema.SnapCount, error) {
	cr := csv.NewReader(r)
	cr.ReuseRecord = true
	cr.FieldsPerRecord = -1

	header, err := cr.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	idx := indexHeader(header)

	var out []schema.SnapCount
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read row: %w", err)
		}
		row := func(key string) string {
			if i, ok := idx[key]; ok && i >= 0 && i < len(rec) {
				return strings.TrimSpace(rec[i])
			}
			return ""
		}

		season := atoi(row("season"))
		week := atoi(row("week"))
		gameID := row("game_id")
		if gameID == "" {
			gameID = row("gameid")
		}
		playerID := firstNonEmpty(row("player_id"), row("gsis_id"), row("playerid"))
		team := strings.ToUpper(row("team"))

		offSnaps := atoi(firstNonEmpty(row("offense_snaps"), row("team_snaps")))
		playerSnaps := atoi(firstNonEmpty(row("player_snaps"), row("snaps")))
		var snapPct float64
		if v := row("snap_pct"); v != "" {
			snapPct = atof(v)
		} else if offSnaps > 0 && playerSnaps >= 0 {
			snapPct = 100.0 * float64(playerSnaps) / float64(offSnaps)
		}

		out = append(out, schema.SnapCount{
			Season:       season,
			Week:         week,
			GameID:       gameID,
			PlayerID:     playerID,
			Team:         team,
			OffenseSnaps: offSnaps,
			PlayerSnaps:  playerSnaps,
			SnapPct:      snapPct,
		})
	}
	return out, nil
}

// ---- helpers ----

func normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	return s
}

func indexHeader(hdr []string) map[string]int {
	// Map normalized header -> index, with a few synonyms added
	idx := make(map[string]int, len(hdr))
	for i, h := range hdr {
		idx[normalize(h)] = i
	}
	// Add synonym keys pointing at the same index if present
	setSyn := func(alias, base string) {
		if i, ok := idx[base]; ok {
			idx[alias] = i
		}
	}
	setSyn("player_id", "gsis_id")
	setSyn("gsis_id", "player_id")
	setSyn("offense_snaps", "team_snaps")
	setSyn("team_snaps", "offense_snaps")
	setSyn("player_snaps", "snaps")
	setSyn("snap_pct", "snap_percentage")
	setSyn("gameid", "game_id")
	setSyn("playerid", "player_id")
	return idx
}

func firstNonEmpty(vs ...string) string {
	for _, v := range vs {
		if v != "" {
			return v
		}
	}
	return ""
}

func atoi(s string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(s))
	return i
}

func atof(s string) float64 {
	f, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return f
}
