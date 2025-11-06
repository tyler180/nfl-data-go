package download

import (
	"net/http"
	"strconv"
	"strings"
)

// ParseRespMeta extracts cache-related metadata from an HTTP response.
// - ETag:           from "ETag" header (stripped of weak/quotes as-is preserved)
// - Last-Modified:  parsed via http.ParseTime
// - ContentLength:  from resp.ContentLength or "Content-Length" header
func ParseRespMeta(resp *http.Response) Metadata {
	var m Metadata

	// ETag (keep as-is; callers can decide how to compare/use it)
	if et := resp.Header.Get("ETag"); et != "" {
		// Trim spaces; keep weak validators (W/...) intact
		m.ETag = strings.TrimSpace(et)
	}

	// Last-Modified
	if lm := resp.Header.Get("Last-Modified"); lm != "" {
		if t, err := http.ParseTime(lm); err == nil {
			m.LastModified = t
		}
	}

	// Content-Length
	if resp.ContentLength >= 0 {
		m.ContentLength = resp.ContentLength
	} else if cl := resp.Header.Get("Content-Length"); cl != "" {
		if n, err := strconv.ParseInt(cl, 10, 64); err == nil {
			m.ContentLength = n
		}
	}

	return m
}
