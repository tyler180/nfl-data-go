package nflreadgo

import (
	"context"
	"sort"
	"time"

	"github.com/tyler180/nfl-data-go/internal/download"
	"github.com/tyler180/nfl-data-go/internal/parse"
	"github.com/tyler180/nfl-data-go/internal/schema"
	"github.com/tyler180/nfl-data-go/internal/source"
	// "github.com/tyler180/nfl-data-go/internal/source"
)

// type CacheMode string

// const (
// 	CacheOff CacheMode = "off"
// 	CacheFS  CacheMode = "filesystem"
// 	CacheMem CacheMode = "memory"
// )

// type Config struct {
// 	CacheMode CacheMode
// 	CacheDir  string
// 	CacheTTL  time.Duration
// 	Timeout   time.Duration
// 	UserAgent string
// 	Verbose   bool
// }

// type Option func(*Config)

// func WithCache(mode CacheMode, dir string, ttl time.Duration) Option {
// 	return func(c *Config) {
// 		c.CacheMode, c.CacheDir, c.CacheTTL = mode, dir, ttl
// 	}
// }
// func WithTimeout(d time.Duration) Option { return func(c *Config) { c.Timeout = d } }
// func WithUserAgent(ua string) Option     { return func(c *Config) { c.UserAgent = ua } }
// func WithVerbose(v bool) Option          { return func(c *Config) { c.Verbose = v } }

// func DefaultConfig() Config {
// 	return Config{
// 		CacheMode: CacheFS,
// 		CacheDir:  "~/.cache/nflreadgo",
// 		CacheTTL:  24 * time.Hour,
// 		Timeout:   30 * time.Second,
// 		UserAgent: "nflreadgo/0.1 (+github.com/tyler180/nfl-data-go)",
// 	}
// }

// // helper to apply options and env (NFLREADGO_CACHE, ...).
// func buildConfig(opts []Option) Config {
// 	cfg := DefaultConfig()
// 	for _, o := range opts {
// 		o(&cfg)
// 	}
// 	return cfg
// }

// ----- Selectors / Params -----

type Seasons []int
type Weeks []int
type SeasonWeeks struct {
	Season int
	Weeks  []int // empty = all available
}

// ----- Return types (example) -----

type SnapCount struct {
	Season       int
	Week         int
	GameID       string
	PlayerID     string
	Team         string
	OffenseSnaps int
	PlayerSnaps  int
	SnapPct      float64 // 0..100
}

func LoadSnapCounts(ctx context.Context, sel any, opts ...Option) ([]schema.SnapCount, error) {
	cfg := buildConfig(opts)
	dl := download.New(
		download.WithUserAgent(cfg.UserAgent),
		download.WithHTTPClient(cfg.HTTPClient()),
		download.WithCache(cfg.CacheBackend()), // fs or memory or nil
	)

	selInt := expandSeasons(sel)
	if len(selInt) == 0 {
		return nil, nil // nothing to load
	}

	urls := source.NFLVerseSnapCountURLs(selInt) // returns []string
	var out []schema.SnapCount

	for _, u := range urls {
		rc, _, err := dl.Fetch(ctx, u)
		if err != nil {
			return nil, err
		}
		rows, err := parse.SnapCountsCSV(rc)
		rc.Close()
		if err != nil {
			return nil, err
		}
		out = append(out, rows...)
	}
	return filterBySelection(out, sel), nil
}

// ---- selection expansion ----

func expandSeasons(sel any) []int {
	switch v := sel.(type) {
	case Seasons:
		cp := append([]int(nil), v...)
		uniqueSort(&cp)
		return cp
	case []SeasonWeeks:
		m := map[int]struct{}{}
		for _, sw := range v {
			m[sw.Season] = struct{}{}
		}
		out := make([]int, 0, len(m))
		for yr := range m {
			out = append(out, yr)
		}
		sort.Ints(out)
		return out
	case bool:
		// true == "all" (follow nflreadpy pattern). For now, return a sane range.
		// You can refine this by reading available tags or index files.
		return []int{2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024, GetCurrentSeason()}
	case int:
		return []int{v}
	case []int:
		cp := append([]int(nil), v...)
		uniqueSort(&cp)
		return cp
	default:
		return nil
	}
}

func uniqueSort(xs *[]int) {
	if xs == nil || len(*xs) == 0 {
		return
	}
	sort.Ints(*xs)
	j := 0
	for i := 1; i < len(*xs); i++ {
		if (*xs)[i] != (*xs)[j] {
			j++
			(*xs)[j] = (*xs)[i]
		}
	}
	*xs = (*xs)[:j+1]
}

// func LoadRostersWeekly(ctx context.Context, seasons Seasons, opts ...Option) ([]RosterWeekly, error) { /* ... */
// 	return nil, nil
// }
// func LoadPlayers(ctx context.Context, opts ...Option) ([]schema.Player, error) { /* ... */ return nil, nil }

func GetCurrentSeason() int { /* mirror nflreadpy */ return time.Now().Year() /* refine by month */ }

// func GetCurrentWeek(ctx context.Context) (int, error) { /* derive from schedules */ return 1, nil }
