package download

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Format int

const (
	FormatParquet Format = iota
	FormatCSV
)

type Client struct {
	http      *http.Client
	cache     Cache // interface in cache.go
	userAgent string
}

type Option func(*Client)

func WithHTTPClient(h *http.Client) Option { return func(c *Client) { c.http = h } }
func WithCache(cache Cache) Option         { return func(c *Client) { c.cache = cache } }
func WithUserAgent(ua string) Option       { return func(c *Client) { c.userAgent = ua } }

func New(opts ...Option) *Client {
	c := &Client{
		http:      &http.Client{Timeout: 30 * time.Second},
		userAgent: "nflreadgo/0.1 (+github.com/tyler180/nfl-data-go)",
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Fetch returns a readable body and closes it when the ctx is done.
// If a cache is configured, it will attempt conditional GETs with ETag/Last-Modified.
func (c *Client) Fetch(ctx context.Context, url string) (io.ReadCloser, Metadata, error) {
	if c.http == nil {
		return nil, Metadata{}, errors.New("nil http client")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, Metadata{}, err
	}
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	// If cache has validators, set If-None-Match / If-Modified-Since
	var meta Metadata
	if c.cache != nil {
		if tag, t, ok := c.cache.Validators(url); ok {
			if tag != "" {
				req.Header.Set("If-None-Match", tag)
			}
			if !t.IsZero() {
				req.Header.Set("If-Modified-Since", t.UTC().Format(http.TimeFormat))
			}
		}
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, Metadata{}, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		meta = ParseRespMeta(resp) // pulls ETag/Last-Modified, Size, etc.
		if c.cache != nil {
			// stream into cache while returning a tee'd reader
			rc, _ := c.cache.StoreStream(url, meta, resp.Body)
			return rc, meta, nil
		}
		return resp.Body, meta, nil

	case http.StatusNotModified:
		if c.cache == nil {
			_ = resp.Body.Close()
			return nil, Metadata{}, errors.New("304 but no cache configured")
		}
		_ = resp.Body.Close()
		rc, meta, err := c.cache.Open(url)
		return rc, meta, err

	default:
		defer resp.Body.Close()
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return nil, Metadata{}, &HTTPError{Code: resp.StatusCode, Body: string(b)}
	}
}

func ParseFormat(s string) (Format, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	switch s {
	case "parquet":
		return FormatParquet, nil
	case "csv":
		return FormatCSV, nil
	default:
		return 0, fmt.Errorf("unknown format %q", s)
	}
}
