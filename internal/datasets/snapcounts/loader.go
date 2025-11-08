package snapcounts

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/tyler180/nfl-data-go/internal/datasets"
)

// NOTE: nflverse path includes "data/..." in the repo.
var src = datasets.Source{Repo: "nflverse/nflverse-data", Base: "data/snap_counts/snap_counts"}

// SnapCount is the typed row for snap counts.
// type SnapCount struct {
// 	Season       int
// 	Week         int
// 	GameID       string
// 	PlayerID     string
// 	Team         string
// 	OffenseSnaps int
// 	PlayerSnaps  int
// 	SnapPct      float64 // 0..100
// }

// FromMap maps a generic CSV row (normalized headers) into SnapCount.
// parse.Auto() normalizes headers to lowercase_with_underscores.
func FromMap(m map[string]any) SnapCount {
	get := func(key string) string {
		if v, ok := m[key]; ok && v != nil {
			return strings.TrimSpace(fmt.Sprint(v))
		}
		return ""
	}
	season := atoi(get("season"))
	week := atoi(get("week"))
	gameID := firstNonEmpty(get("game_id"), get("gameid"))
	playerID := firstNonEmpty(get("player_id"), get("gsis_id"), get("playerid"))
	team := strings.ToUpper(get("team"))

	offSnaps := atoi(firstNonEmpty(get("offense_snaps"), get("team_snaps")))
	playerSnaps := atoi(firstNonEmpty(get("player_snaps"), get("snaps")))

	var pct float64
	if v := get("snap_pct"); v != "" {
		pct = atof(v)
	} else if offSnaps > 0 && playerSnaps >= 0 {
		pct = 100.0 * float64(playerSnaps) / float64(offSnaps)
	}

	return SnapCount{
		Season:       season,
		Week:         week,
		GameID:       gameID,
		PlayerID:     playerID,
		Team:         team,
		OffenseSnaps: offSnaps,
		PlayerSnaps:  playerSnaps,
		SnapPct:      pct,
	}
}

// All seasons (combined file if provided; otherwise base per-repo behavior)
func Load(ctx context.Context) ([]SnapCount, error) {
	return datasets.LoadFromSourceAs[SnapCount](ctx, src, 0, FromMap)
}

func LoadSeason(ctx context.Context, season int) ([]SnapCount, error) {
	return datasets.LoadFromSourceAs[SnapCount](ctx, src, season, FromMap)
}

// ---- helpers ----

func firstNonEmpty(vs ...string) string {
	for _, v := range vs {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func atoi(s string) int {
	n := 0
	sign := 1
	for i, r := range s {
		if i == 0 && r == '-' {
			sign = -1
			continue
		}
		if r < '0' || r > '9' {
			break
		}
		n = n*10 + int(r-'0')
	}
	return sign * n
}

func atof(s string) float64 {
	// minimal, fast parser good enough for CSV numeric fields
	var sign float64 = 1
	i := 0
	if len(s) > 0 && s[0] == '-' {
		sign = -1
		i++
	}
	var intPart int64
	var fracPart int64
	var fracDiv float64 = 1
	seenDot := false
	for ; i < len(s); i++ {
		c := s[i]
		if c == '.' && !seenDot {
			seenDot = true
			continue
		}
		if c < '0' || c > '9' {
			break
		}
		if !seenDot {
			intPart = intPart*10 + int64(c-'0')
		} else {
			fracPart = fracPart*10 + int64(c-'0')
			fracDiv *= 10
		}
	}
	return sign * (float64(intPart) + float64(fracPart)/fracDiv)
}

const defaultSnapPerSeasonPattern = "https://raw.githubusercontent.com/nflverse/nflverse-data/master/data/snap_counts/snap_counts_%d.csv"

// NFLVerseSnapCountURLs returns URLs for the given seasons.
// If NFLREADGO_SNAP_URL is set, it returns exactly that single URL.
func NFLVerseSnapCountURLs(seasons []int) []string {
	if u := os.Getenv("NFLREADGO_SNAP_URL"); u != "" {
		return []string{u}
	}
	pattern := os.Getenv("NFLREADGO_SNAP_PATTERN")
	if pattern == "" {
		pattern = defaultSnapPerSeasonPattern
	}
	ss := append([]int(nil), seasons...)
	sort.Ints(ss)

	urls := make([]string, 0, len(ss))
	for _, yr := range ss {
		urls = append(urls, fmt.Sprintf(pattern, yr))
	}
	return urls
}
