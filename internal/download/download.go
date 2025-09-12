// Package nflreadgo provides data downloading functionality similar to the
// Python NflverseDownloader in downloader.py. It supports repository-aware URL
// building, format preference with automatic CSV<->Parquet fallback, a simple
// on-disk cache, custom headers, and basic CSV parsing helpers.
//
// Usage example:
//
//	cfg := DefaultConfig()
//	cfg.Verbose = true
//	dl := NewDownloader(cfg)
//	data, usedURL, err := dl.Download("nflverse-data", "week1/team_stats", nil, nil)
//	if err != nil { log.Fatal(err) }
//	rows, _ := CSVToMaps(data) // if CSV; Parquet bytes are returned as-is
//	fmt.Println("downloaded from", usedURL, "rows:", len(rows))
//
// Notes:
//   - Parquet bytes are returned intact; parsing to rows requires a library.
//     See ParseParquet placeholder for a suggested approach.
//   - The cache is a very simple file store keyed by the request URL.
//   - Format preference defaults to cfg.PreferFormat but can be overridden per-call.
package nflreadgo

import (
	"bytes"
	"crypto/sha1"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	parquet "github.com/parquet-go/parquet-go"
)

// Format represents a data serialization format.
type Format int

const (
	FormatParquet Format = iota
	FormatCSV
)

func (f Format) String() string {
	switch f {
	case FormatParquet:
		return "parquet"
	case FormatCSV:
		return "csv"
	default:
		return "unknown"
	}
}

// ParseFormat parses a string into a Format.
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

// Base repository URLs mirroring the Python implementation.
var baseURLs = map[string]string{
	"nflverse-data":  "https://github.com/nflverse/nflverse-data/releases/download/",
	"nfldata":        "https://github.com/nflverse/nfldata/raw/master/data/",
	"espnscraper":    "https://github.com/nflverse/espnscrapeR-data/raw/master/data/",
	"dynastyprocess": "https://github.com/dynastyprocess/data/raw/master/files/",
	"ffopportunity":  "https://github.com/ffverse/ffopportunity/releases/download/",
}

// Config controls downloader behavior.
type Config struct {
	UserAgent    string
	PreferFormat Format
	Verbose      bool
	Timeout      time.Duration
	CacheDir     string
	CacheEnabled bool
}

// DefaultConfig returns a usable default configuration.
func DefaultConfig() *Config {
	cacheDir := filepath.Join(os.TempDir(), "nflreadgo-cache")
	return &Config{
		UserAgent:    "nflreadgo/1.0 (+https://github.com/your/repo)",
		PreferFormat: FormatParquet,
		Verbose:      false,
		Timeout:      30 * time.Second,
		CacheDir:     cacheDir,
		CacheEnabled: true,
	}
}

// CacheManager is a tiny file-based cache keyed by the URL.
type CacheManager struct {
	enabled bool
	dir     string
}

func NewCacheManager(cfg *Config) *CacheManager {
	_ = os.MkdirAll(cfg.CacheDir, 0o755)
	return &CacheManager{enabled: cfg.CacheEnabled, dir: cfg.CacheDir}
}

func (c *CacheManager) key(u string) string {
	h := sha1.Sum([]byte(u))
	return hex.EncodeToString(h[:])
}

func (c *CacheManager) pathFor(u string) string {
	return filepath.Join(c.dir, c.key(u))
}

func (c *CacheManager) Get(u string) ([]byte, bool) {
	if !c.enabled {
		return nil, false
	}
	p := c.pathFor(u)
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, false
	}
	return b, true
}

func (c *CacheManager) Set(u string, data []byte) error {
	if !c.enabled {
		return nil
	}
	p := c.pathFor(u)
	return os.WriteFile(p, data, 0o644)
}

// Downloader fetches data from nflverse repositories.
type Downloader struct {
	client *http.Client
	cfg    *Config
	cache  *CacheManager
}

// NewDownloader constructs a Downloader with the given config.
func NewDownloader(cfg *Config) *Downloader {
	return &Downloader{
		client: &http.Client{Timeout: cfg.Timeout},
		cfg:    cfg,
		cache:  NewCacheManager(cfg),
	}
}

// headers returns the HTTP headers for requests.
func (d *Downloader) headers() http.Header {
	h := http.Header{}
	h.Set("User-Agent", d.cfg.UserAgent)
	h.Set("Accept", "application/octet-stream, text/csv, */*")
	return h
}

// buildURL builds the full URL for a data file, appending an extension when omitted.
func (d *Downloader) buildURL(repository, p string, format Format) (string, error) {
	base, ok := baseURLs[repository]
	if !ok {
		return "", fmt.Errorf("unknown repository: %s", repository)
	}

	if !strings.HasSuffix(p, ".parquet") && !strings.HasSuffix(p, ".csv") {
		ext := ".parquet"
		if format == FormatCSV {
			ext = ".csv"
		}
		p = p + ext
	}

	// Use url.JoinPath when possible to avoid duplicate slashes.
	joined, err := url.JoinPath(base, p)
	if err != nil {
		// Fallback: simple concat
		joined = strings.TrimRight(base, "/") + "/" + strings.TrimLeft(p, "/")
	}
	return joined, nil
}

// DownloadOptions controls caching per-call.
type DownloadOptions struct {
	Force bool // if true, bypass cache
}

// downloadOnce downloads the content of a URL, honoring cache and headers.
func (d *Downloader) downloadOnce(u string, opts *DownloadOptions) ([]byte, error) {
	if opts == nil {
		opts = &DownloadOptions{}
	}
	if !opts.Force {
		if b, ok := d.cache.Get(u); ok {
			if d.cfg.Verbose {
				fmt.Println("cache hit:", u)
			}
			return b, nil
		}
	}

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header = d.headers()

	if d.cfg.Verbose {
		fmt.Println("downloading:", u)
	}
	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download %s: %w", u, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP error %d for %s", resp.StatusCode, u)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading %s: %w", u, err)
	}
	_ = d.cache.Set(u, b)
	return b, nil
}

// Download fetches a file from a repo and returns its bytes and the final URL used.
// It tries the preferred format first, then automatically falls back to the
// alternative format on network/HTTP errors.
func (d *Downloader) Download(repository, path string, formatPreference *Format, opts *DownloadOptions) ([]byte, string, error) {
	format := d.cfg.PreferFormat
	if formatPreference != nil {
		format = *formatPreference
	}

	primaryURL, err := d.buildURL(repository, path, format)
	if err != nil {
		return nil, "", err
	}

	b, err := d.downloadOnce(primaryURL, opts)
	if err == nil {
		return b, primaryURL, nil
	}

	// Decide the alternative format.
	alt := FormatCSV
	if format == FormatCSV {
		alt = FormatParquet
	}
	if d.cfg.Verbose {
		fmt.Printf("failed to download %s (pref %s), trying %s...\n", primaryURL, format, alt)
	}
	altURL, err2 := d.buildURL(repository, path, alt)
	if err2 != nil {
		return nil, "", err // original error is more relevant
	}
	b2, err2 := d.downloadOnce(altURL, opts)
	if err2 != nil {
		return nil, "", err // bubble original error for clarity
	}
	return b2, altURL, nil
}

// DetectFormatFromURL infers Format from the file extension.
func DetectFormatFromURL(u string) (Format, error) {
	if strings.HasSuffix(strings.ToLower(u), ".parquet") {
		return FormatParquet, nil
	}
	if strings.HasSuffix(strings.ToLower(u), ".csv") {
		return FormatCSV, nil
	}
	return 0, errors.New("unable to detect format from URL")
}

// CSVToMaps parses CSV bytes into a slice of rows keyed by column name.
func CSVToMaps(data []byte) ([]map[string]string, error) {
	r := csv.NewReader(strings.NewReader(string(data)))
	r.ReuseRecord = false
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	headers := rows[0]
	var out []map[string]string
	for _, row := range rows[1:] {
		m := make(map[string]string, len(headers))
		for i := range headers {
			var v string
			if i < len(row) {
				v = row[i]
			}
			m[headers[i]] = v
		}
		out = append(out, m)
	}
	return out, nil
}

// ParseParquet decodes Parquet bytes into a slice of generic row maps using
// github.com/parquet-go/parquet-go. This keeps the downloader dependency-light
// while leaning on a widely adopted, actively maintained library.
//
// Notes:
//   - Logical types (DECIMAL, TIMESTAMP, etc.) are mapped to Go primitives
//     where possible by parquet-go. If you need custom coercions, transform the
//     returned maps after parsing.
//   - For very large files, prefer streaming APIs instead of loading all rows.
func ParseParquet(b []byte) ([]map[string]any, error) {
	size := int64(len(b))
	rows, err := parquet.Read[any](bytes.NewReader(b), size)
	if err != nil {
		return nil, fmt.Errorf("parquet read failed: %w", err)
	}
	out := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		switch v := r.(type) {
		case map[string]any:
			out = append(out, v)
		default:
			// Fallback: JSON roundtrip to map
			var m map[string]any
			jb, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("marshal parquet row: %w", err)
			}
			if err := json.Unmarshal(jb, &m); err != nil {
				return nil, fmt.Errorf("unmarshal parquet row: %w", err)
			}
			out = append(out, m)
		}
	}
	return out, nil
}

// ---- Global singleton (optional), mirroring get_downloader() in Python ----

var defaultDownloader = NewDownloader(DefaultConfig())

// GetDownloader returns the package-level Downloader instance.
func GetDownloader() *Downloader { return defaultDownloader }
