package source

import (
	"fmt"
	"os"
	"sort"
)

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
