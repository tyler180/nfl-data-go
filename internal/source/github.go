package source

import "fmt"

// RawGitHubURL builds a raw.githubusercontent.com URL from repo + path.
// Accepts either "owner/repo" (preferred) or just "repo" (owner defaults to "nflverse").
// 'path' should include any subdirs + filename + extension (e.g., "data/snap_counts/snap_counts_2024.csv").
func RawGitHubURL(repo, path string) string {
	owner := "nflverse"
	name := repo
	if slash := indexByte(repo, '/'); slash >= 0 {
		owner, name = repo[:slash], repo[slash+1:]
	}
	// Branch is 'master' across nflverse-data; adjust to 'main' if needed per repo.
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/master/%s", owner, name, path)
}

func indexByte(s string, b byte) int {
	for i := range s {
		if s[i] == b {
			return i
		}
	}
	return -1
}
