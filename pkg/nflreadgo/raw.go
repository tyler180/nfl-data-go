package nflreadgo

import (
	"context"
	"io"
	"net/http"

	"github.com/tyler180/nfl-data-go/internal/download"
	"github.com/tyler180/nfl-data-go/internal/source"
)

// LoadSnapCountsRaw returns a byte slice for each resolved asset (usually one per season,
// or a single combined file if NFLREADGO_SNAP_URL is set). It also returns a best-effort
// MIME type for each blob (derived via http.DetectContentType).
func LoadSnapCountsRaw(ctx context.Context, sel any, opts ...Option) (blobs [][]byte, mimes []string, err error) {
	cfg := buildConfig(opts)
	dl := download.New(
		download.WithHTTPClient(cfg.HTTPClient()),
		download.WithCache(cfg.CacheBackend()),
		download.WithUserAgent(cfg.UserAgent),
	)

	seasons := expandSeasons(sel)
	if len(seasons) == 0 {
		seasons = []int{GetCurrentSeason()}
	}

	urls := source.NFLVerseSnapCountURLs(seasons)
	blobs = make([][]byte, 0, len(urls))
	mimes = make([]string, 0, len(urls))

	for _, u := range urls {
		rc, _, e := dl.Fetch(ctx, u)
		if e != nil {
			return nil, nil, e
		}
		b, e := io.ReadAll(rc)
		rc.Close()
		if e != nil {
			return nil, nil, e
		}
		mt := "application/octet-stream"
		if len(b) > 0 {
			peek := b
			if len(peek) > 512 {
				peek = peek[:512]
			}
			mt = http.DetectContentType(peek)
		}
		blobs = append(blobs, b)
		mimes = append(mimes, mt)
	}
	return blobs, mimes, nil
}
