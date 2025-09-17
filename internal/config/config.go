// Package nflreadgo: config.go
//
// This file provides a Go equivalent of the Python `config.py` used by
// nflreadpy. It defines:
//   - CacheMode and preferred data format settings
//   - A package-level Config with sane defaults
//   - Environment variable overrides compatible with nflreadpy
//   - NFLREADPY_CACHE              (memory|filesystem|off)
//   - NFLREADPY_CACHE_DIR          (path)
//   - NFLREADPY_CACHE_DURATION     (seconds)
//   - NFLREADPY_PREFER             (parquet|csv)
//   - NFLREADPY_VERBOSE            (true|false)
//   - NFLREADPY_TIMEOUT            (seconds)
//   - NFLREADPY_USER_AGENT         (string)
//   - Functions to get/update/reset the config and to apply it to the
//     default downloader and cache.
//
// Notes
//   - `.env` support: if a `.env` file is present in the working directory,
//     we parse simple KEY=VALUE lines and use them as fallbacks when OS env
//     vars are absent (comments and blank lines are ignored).
//   - The defaults mirror the Python module: memory cache, 24h TTL, Parquet
//     preferred, verbose on, 30s timeout, and an nflverse-flavored UA.
package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	cachepkg "github.com/tyler180/nfl-data-go/internal/cache"
	downloadpkg "github.com/tyler180/nfl-data-go/internal/download"
)

// Version is the library version used in the default User-Agent.
const Version = "v0.1.0"

// CacheMode controls where downloaded data is cached.
type CacheMode string

const (
	CacheModeMemory     CacheMode = "memory"
	CacheModeFilesystem CacheMode = "filesystem"
	CacheModeOff        CacheMode = "off"
)

// DataFormat mirrors the Python enum (but we map to our internal Format).
type DataFormat string

const (
	DataFormatParquet DataFormat = "parquet"
	DataFormatCSV     DataFormat = "csv"
)

// AppConfig contains user-configurable settings. It is intentionally
// independent from Downloader.Config to keep a stable public surface.
type AppConfig struct {
	CacheMode     CacheMode
	CacheDir      string
	CacheDuration time.Duration // TTL for cache entries

	Prefer    downloadpkg.Format // preferred download format
	Verbose   bool
	Timeout   time.Duration // HTTP timeout
	UserAgent string
}

// defaultCacheDir attempts to mirror platformdirs.user_cache_dir("nflreadpy").
func defaultCacheDir() string {
	if d, err := os.UserCacheDir(); err == nil && d != "" {
		return filepath.Join(d, "nflreadpy") // keep python-compatible folder name
	}
	return filepath.Join(os.TempDir(), "nflreadpy")
}

// DefaultAppConfig returns library defaults analogous to nflreadpy.
func DefaultAppConfig() *AppConfig {
	return &AppConfig{
		CacheMode:     CacheModeMemory,
		CacheDir:      defaultCacheDir(),
		CacheDuration: 24 * time.Hour,
		Prefer:        downloadpkg.FormatParquet,
		Verbose:       true,
		Timeout:       30 * time.Second,
		UserAgent:     fmt.Sprintf("nflverse/nflreadgo %s (%s)", Version, runtime.Version()),
	}
}

var (
	cfgMu       sync.RWMutex
	globalCfg   *AppConfig
	dotenvOnce  sync.Once
	dotenvPairs map[string]string
)

// GetConfig returns a snapshot pointer to the current config.
func GetConfig() *AppConfig {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	c := *globalCfg
	return &c
}

// ResetConfig resets the global config to defaults, then applies .env and
// environment variable overrides, and wires the downloader+cache to match.
func ResetConfig() {
	cfgMu.Lock()
	defer cfgMu.Unlock()
	globalCfg = DefaultAppConfig()
	loadDotEnvIfPresent()
	applyEnvOverrides(globalCfg)
	applyToSubsystems(globalCfg)
}

// UpdateConfig applies one or more functional options, then re-wires the
// subsystems (downloader + cache) to reflect the new settings.
func UpdateConfig(opts ...ConfigOption) {
	cfgMu.Lock()
	defer cfgMu.Unlock()
	for _, opt := range opts {
		if opt != nil {
			opt(globalCfg)
		}
	}
	applyToSubsystems(globalCfg)
}

// ConfigOption is a functional option for UpdateConfig.
type ConfigOption func(*AppConfig)

func WithCacheMode(m CacheMode) ConfigOption { return func(c *AppConfig) { c.CacheMode = m } }
func WithCacheDir(dir string) ConfigOption {
	return func(c *AppConfig) {
		if dir != "" {
			c.CacheDir = dir
		}
	}
}
func WithCacheDuration(ttl time.Duration) ConfigOption {
	return func(c *AppConfig) {
		if ttl >= 0 {
			c.CacheDuration = ttl
		}
	}
}
func WithPreferFormat(f downloadpkg.Format) ConfigOption { return func(c *AppConfig) { c.Prefer = f } }
func WithVerbose(v bool) ConfigOption                    { return func(c *AppConfig) { c.Verbose = v } }
func WithTimeout(d time.Duration) ConfigOption {
	return func(c *AppConfig) {
		if d > 0 {
			c.Timeout = d
		}
	}
}
func WithUserAgent(ua string) ConfigOption {
	return func(c *AppConfig) {
		if ua != "" {
			c.UserAgent = ua
		}
	}
}

// applyToSubsystems wires the downloader and the package-level cache
// to reflect the current global configuration.
func applyToSubsystems(c *AppConfig) {
	// Configure the package-level cache
	switch c.CacheMode {
	case CacheModeMemory:
		_ = cachepkg.SetCacheOptions(cachepkg.CacheMemory, c.CacheDuration, c.CacheDir, 256)
	case CacheModeFilesystem:
		_ = cachepkg.SetCacheOptions(cachepkg.CacheDisk, c.CacheDuration, c.CacheDir, 0)
	case CacheModeOff:
		_ = cachepkg.SetCacheOptions(cachepkg.CacheOff, c.CacheDuration, c.CacheDir, 0)
	default:
		_ = cachepkg.SetCacheOptions(cachepkg.CacheMemory, c.CacheDuration, c.CacheDir, 256)
	}
}

// --- Environment & .env helpers ---

func loadDotEnvIfPresent() {
	dotenvOnce.Do(func() {
		dotenvPairs = make(map[string]string)
		f, err := os.Open(".env")
		if err != nil {
			return
		}
		defer f.Close()
		s := bufio.NewScanner(f)
		for s.Scan() {
			line := strings.TrimSpace(s.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if i := strings.IndexByte(line, '='); i > 0 {
				k := strings.TrimSpace(line[:i])
				v := strings.TrimSpace(line[i+1:])
				// remove optional surrounding quotes
				v = strings.Trim(v, "\"'")
				dotenvPairs[k] = v
			}
		}
	})
}

func envOrDotenv(key string) (string, bool) {
	if v := os.Getenv(key); v != "" {
		return v, true
	}
	if dotenvPairs == nil {
		return "", false
	}
	v, ok := dotenvPairs[key]
	return v, ok
}

func applyEnvOverrides(c *AppConfig) {
	if v, ok := envOrDotenv("NFLREADPY_CACHE"); ok {
		s := strings.ToLower(strings.TrimSpace(v))
		switch s {
		case "memory":
			c.CacheMode = CacheModeMemory
		case "filesystem", "disk":
			c.CacheMode = CacheModeFilesystem
		case "off", "none", "disabled":
			c.CacheMode = CacheModeOff
		}
	}
	if v, ok := envOrDotenv("NFLREADPY_CACHE_DIR"); ok {
		if v != "" {
			c.CacheDir = v
		}
	}
	if v, ok := envOrDotenv("NFLREADPY_CACHE_DURATION"); ok {
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil && n >= 0 {
			c.CacheDuration = time.Duration(n) * time.Second
		}
	}
	if v, ok := envOrDotenv("NFLREADPY_PREFER"); ok {
		if f, err := downloadpkg.ParseFormat(v); err == nil {
			c.Prefer = f
		}
	}
	if v, ok := envOrDotenv("NFLREADPY_VERBOSE"); ok {
		if b, err := parseBool(v); err == nil {
			c.Verbose = b
		}
	}
	if v, ok := envOrDotenv("NFLREADPY_TIMEOUT"); ok {
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil && n > 0 {
			c.Timeout = time.Duration(n) * time.Second
		}
	}
	if v, ok := envOrDotenv("NFLREADPY_USER_AGENT"); ok {
		v = strings.TrimSpace(v)
		if v != "" {
			c.UserAgent = v
		}
	}
}

func parseBool(s string) (bool, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "1", "true", "t", "yes", "y", "on":
		return true, nil
	case "0", "false", "f", "no", "n", "off":
		return false, nil
	default:
		return false, errors.New("invalid bool")
	}
}

// init sets up the global config on package import.
func init() { ResetConfig() }
