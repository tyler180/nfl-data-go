package nflreadgo

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/tyler180/nfl-data-go/internal/download"
)

type CacheMode string

const (
	CacheOff CacheMode = "off"
	CacheFS  CacheMode = "filesystem"
	CacheMem CacheMode = "memory"
)

type Config struct {
	CacheMode CacheMode
	CacheDir  string
	CacheTTL  time.Duration

	Timeout   time.Duration
	UserAgent string
	Verbose   bool
}

type Option func(*Config)

func WithCache(mode CacheMode, dir string, ttl time.Duration) Option {
	return func(c *Config) { c.CacheMode, c.CacheDir, c.CacheTTL = mode, dir, ttl }
}
func WithTimeout(d time.Duration) Option { return func(c *Config) { c.Timeout = d } }
func WithUserAgent(ua string) Option     { return func(c *Config) { c.UserAgent = ua } }
func WithVerbose(v bool) Option          { return func(c *Config) { c.Verbose = v } }

func DefaultConfig() Config {
	return Config{
		CacheMode: CacheFS,
		CacheDir:  os.ExpandEnv("$HOME/.cache/nflreadgo"),
		CacheTTL:  24 * time.Hour,

		Timeout:   30 * time.Second,
		UserAgent: "nflreadgo/0.1 (+github.com/tyler180/nfl-data-go)",
	}
}

// buildConfig applies options and environment variables.
func buildConfig(opts []Option) Config {
	c := DefaultConfig()
	for _, o := range opts {
		o(&c)
	}

	// --- env overrides (optional, mirrors nflreadpy-style knobs) ---
	// NFLREADGO_CACHE = off|filesystem|memory
	if v := os.Getenv("NFLREADGO_CACHE"); v != "" {
		switch v {
		case "off":
			c.CacheMode = CacheOff
		case "filesystem":
			c.CacheMode = CacheFS
		case "memory":
			c.CacheMode = CacheMem
		}
	}
	if v := os.Getenv("NFLREADGO_CACHE_DIR"); v != "" {
		c.CacheDir = v
	}
	if v := os.Getenv("NFLREADGO_CACHE_TTL"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			c.CacheTTL = time.Duration(n) * time.Second
		}
	}
	if v := os.Getenv("NFLREADGO_TIMEOUT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			c.Timeout = time.Duration(n) * time.Second
		}
	}
	if v := os.Getenv("NFLREADGO_USER_AGENT"); v != "" {
		c.UserAgent = v
	}
	if v := os.Getenv("NFLREADGO_VERBOSE"); v != "" {
		c.Verbose = v == "1" || v == "true" || v == "TRUE"
	}
	return c
}

// HTTPClient returns a ready http.Client that respects Config.Timeout.
// Callers should not mutate its Transport.
func (c Config) HTTPClient() *http.Client {
	return &http.Client{
		Timeout: c.Timeout,
		// You can inject a custom Transport here (proxy, TLS tweaks, retries, etc.)
	}
}

// CacheBackend returns a download.Cache implementation wired to the config.
// Returns nil when caching is disabled.
func (c Config) CacheBackend() download.Cache {
	switch c.CacheMode {
	case CacheFS:
		if c.CacheDir == "" {
			return nil
		}
		return download.NewFSCache(c.CacheDir, c.CacheTTL)
	case CacheMem:
		return download.NewMemCache(c.CacheTTL)
	default:
		return nil
	}
}
