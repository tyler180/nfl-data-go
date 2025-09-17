package load

import (
	"github.com/tyler180/nfl-data-go/internal/download"
)

// Package nflreadgo: load_players.go
//
// This mirrors Python's load_players() using the Go downloader in this package.
// It downloads the players dataset from nflverse and returns parsed rows.
//
// Reference:
//  • Source repo path: nflverse-data / players/players (CSV or Parquet)
//  • Data dictionary: https://nflreadr.nflverse.com/articles/dictionary_players.html
//  • R/Python analogue: nflreadr::load_players()

// LoadPlayers downloads the canonical NFL players dataset and returns a slice
// of row maps (column name → value). The underlying downloader automatically
// picks Parquet (preferred) and falls back to CSV if needed, then ParseAuto
// normalizes the rows to []map[string]any.
func LoadPlayers() ([]map[string]any, error) {
	b, usedURL, err := download.GetDownloader().Download("nflverse-data", "players/players", nil, nil)
	if err != nil {
		return nil, err
	}
	rows, err := download.ParseAuto(b, usedURL)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// LoadPlayersRaw returns the raw bytes plus the URL used (helpful for caching
// or persisting the source). Most callers can use LoadPlayers() directly.
func LoadPlayersRaw() ([]byte, string, error) {
	return download.GetDownloader().Download("nflverse-data", "players/players", nil, nil)
}

// Example usage:
//
//  func main() {
//      rows, err := nflreadgo.LoadPlayers()
//      if err != nil { log.Fatal(err) }
//      fmt.Println("players:", len(rows))
//      // Access a field safely
//      if len(rows) > 0 {
//          name, _ := rows[0]["player_name"].(string)
//          fmt.Println("first player:", name)
//      }
//  }
