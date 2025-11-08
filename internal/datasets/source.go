package datasets

import (
	"context"
	"io"
	"strings"

	"github.com/tyler180/nfl-data-go/internal/download"
	"github.com/tyler180/nfl-data-go/internal/parse"
	"github.com/tyler180/nfl-data-go/internal/source"
)

// Source: where a dataset lives in a GitHub repo.
// Repo examples:
//
//	"nflverse/nflverse-data" (preferred explicit form)
//	"nflverse-data"          (owner defaults to "nflverse")
//
// Base examples:
//
//	"data/injuries/injuries"
//	"files/db_playerids"
type Source struct {
	Repo string
	Base string
}

// SeasonPath returns base or base_YYYY when season > 0.
func SeasonPath(base string, season int) string {
	if season > 0 {
		return base + "_" + itoa(season)
	}
	return base
}

// LoadFromSourceAs downloads (Repo, Base[_season]) and maps rows using mapper.
// CSV vs Parquet is auto-detected by parse.Auto.
// If a season-specific asset 404s, this automatically falls back to the base asset.
func LoadFromSourceAs[T any](ctx context.Context, src Source, season int, mapper func(map[string]any) T) ([]T, error) {
	dl := download.New() // default client (30s timeout, no cache unless you add options)

	// Try season-scoped first when requested.
	if season > 0 {
		p := SeasonPath(src.Base, season)
		url := source.RawGitHubURL(src.Repo, p)
		rc, _, err := dl.Fetch(ctx, url)
		if err == nil {
			defer rc.Close()
			b, rerr := io.ReadAll(rc)
			if rerr != nil {
				return nil, rerr
			}
			rows, rerr := parse.Auto(b, url)
			if rerr != nil {
				return nil, rerr
			}
			out := make([]T, 0, len(rows))
			for _, r := range rows {
				out = append(out, mapper(r))
			}
			return out, nil
		}
		// Fallback to base if the season file isn't published.
		if !shouldFallbackToBase(err) {
			return nil, err
		}
	}

	// Base (all seasons)
	url := source.RawGitHubURL(src.Repo, src.Base)
	rc, _, err := dl.Fetch(ctx, url)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	b, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	rows, err := parse.Auto(b, url)
	if err != nil {
		return nil, err
	}
	out := make([]T, 0, len(rows))
	for _, r := range rows {
		out = append(out, mapper(r))
	}
	return out, nil
}

// LoadFromPathAs is a convenience for direct (repo, path) loads without a season param.
func LoadFromPathAs[T any](ctx context.Context, repo, path string, mapper func(map[string]any) T) ([]T, error) {
	dl := download.New()
	url := source.RawGitHubURL(repo, path)
	rc, _, err := dl.Fetch(ctx, url)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	b, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	rows, err := parse.Auto(b, url)
	if err != nil {
		return nil, err
	}
	out := make([]T, 0, len(rows))
	for _, r := range rows {
		out = append(out, mapper(r))
	}
	return out, nil
}

// shouldFallbackToBase reports whether a season-scoped download error should retry the base asset.
// We conservatively match common 404 phrases so callers don't need to import internal error types.
func shouldFallbackToBase(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "404") || strings.Contains(s, "not found")
}

// tiny local itoa to avoid extra deps here
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	var b [20]byte
	pos := len(b)
	for i > 0 {
		pos--
		b[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		b[pos] = '-'
	}
	return string(b[pos:])
}
