package datasets

import (
	"fmt"
	"strings"

	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

// Source describes where a dataset lives and the base path inside that repo.
// Examples:
//
//	Repo: "nflverse-data",      Base: "injuries/injuries"
//	Repo: "DynastyProcess/data", Base: "files/db_playerids"
type Source struct {
	Repo string
	Base string
}

// SeasonPath returns base or base_YYYY when season > 0.
func SeasonPath(base string, season int) string {
	if season > 0 {
		return fmt.Sprintf("%s_%d", base, season)
	}
	return base
}

// LoadFromSourceAs downloads (Repo, Base[_season]) and maps rows using mapper.
// CSV vs Parquet is auto-detected by the downloader's ParseAuto.
// If a season-specific asset 404s, this automatically falls back to the base asset.
func LoadFromSourceAs[T any](src Source, season int, mapper func(map[string]any) T) ([]T, error) {
	// Try season-scoped first when requested.
	if season > 0 {
		p := SeasonPath(src.Base, season)
		b, usedURL, err := downloadpkg.Get().Download(src.Repo, p, nil, nil)
		if err == nil {
			rows, err := downloadpkg.ParseAuto(b, usedURL)
			if err != nil {
				return nil, err
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
	b, usedURL, err := downloadpkg.Get().Download(src.Repo, src.Base, nil, nil)
	if err != nil {
		return nil, err
	}
	rows, err := downloadpkg.ParseAuto(b, usedURL)
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
	return strings.Contains(s, "HTTP error 404") || strings.Contains(s, "404 Not Found") || strings.Contains(s, "not found")
}

// LoadFromPathAs is a convenience for direct (repo, path) loads without a season param.
func LoadFromPathAs[T any](repo, path string, mapper func(map[string]any) T) ([]T, error) {
	b, usedURL, err := downloadpkg.Get().Download(repo, path, nil, nil)
	if err != nil {
		return nil, err
	}
	rows, err := downloadpkg.ParseAuto(b, usedURL)
	if err != nil {
		return nil, err
	}
	out := make([]T, 0, len(rows))
	for _, r := range rows {
		out = append(out, mapper(r))
	}
	return out, nil
}
