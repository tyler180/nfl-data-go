package ffplayerids

import (
	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

// Load fetches the DynastyProcess player IDs table and returns typed rows.
// Source: https://github.com/DynastyProcess/data (files/db_playerids.csv)
func Load() ([]FFPlayerID, error) {
	return loadHelper("DynastyProcess/data", "files/db_playerids")
}

// LoadRaw returns the raw bytes and provenance URL.
func LoadRaw() ([]byte, string, error) {
	return downloadpkg.Get().Download("DynastyProcess/data", "files/db_playerids", nil, nil)
}

func loadHelper(repo, path string) ([]FFPlayerID, error) {
	b, usedURL, err := downloadpkg.Get().Download(repo, path, nil, nil)
	if err != nil {
		return nil, err
	}
	rows, err := downloadpkg.ParseAuto(b, usedURL)
	if err != nil {
		return nil, err
	}
	out := make([]FFPlayerID, 0, len(rows))
	for _, r := range rows {
		out = append(out, FromMap(r))
	}
	return out, nil
}
