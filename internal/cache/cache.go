package cache

// Package cache: cache.go
//
// A small, pragmatic cache layer inspired by the Python cache module.
// It supports in‑memory and filesystem caching with TTL-based expiration,
// easy clearing, and an optional default (package-level) cache instance.
//
// Notes
//  • Keys are arbitrary strings (e.g., URLs, repo/path, dataset names).
//  • Disk entries are stored under CacheDir with filenames derived from
//    SHA‑1(key). Expiration for disk is computed from file mtime + TTL.
//  • Pattern clearing uses filepath.Match against the *key* (not the hash).
//  • This cache is independent of Downloader’s built‑in cache, but you can
//    point both at the same CacheDir if you like. If you want Downloader to
//    rely solely on this cache, wire Cache.Get/Set into your download path.

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type CacheMode int

const (
	CacheOff CacheMode = iota
	CacheMemory
	CacheDisk
	CacheBoth
)

// Cache provides memory+disk caching with TTL.
type Cache struct {
	mu     sync.RWMutex
	mode   CacheMode
	ttl    time.Duration
	dir    string
	maxMem int
	mem    map[string]memEntry
}

type memEntry struct {
	data      []byte
	expiresAt time.Time
}

// NewCache constructs a Cache. If mode uses disk, dir will be created.
func NewCache(mode CacheMode, ttl time.Duration, dir string, maxMemEntries int) *Cache {
	c := &Cache{mode: mode, ttl: ttl, dir: dir, maxMem: maxMemEntries, mem: make(map[string]memEntry)}
	if (mode == CacheDisk || mode == CacheBoth) && dir != "" {
		_ = os.MkdirAll(dir, 0o755)
	}
	return c
}

func (c *Cache) diskPathFor(key string) string {
	h := sha1.Sum([]byte(key))
	return filepath.Join(c.dir, hex.EncodeToString(h[:]))
}

// Get returns cached bytes if present and not expired.
func (c *Cache) Get(key string) ([]byte, bool) {
	if c == nil || c.mode == CacheOff {
		return nil, false
	}
	// Memory first
	if c.mode == CacheMemory || c.mode == CacheBoth {
		c.mu.RLock()
		ent, ok := c.mem[key]
		c.mu.RUnlock()
		if ok {
			if time.Now().Before(ent.expiresAt) {
				return ent.data, true
			}
			// expired
			c.mu.Lock()
			delete(c.mem, key)
			c.mu.Unlock()
		}
	}
	// Disk
	if (c.mode == CacheDisk || c.mode == CacheBoth) && c.dir != "" {
		p := c.diskPathFor(key)
		fi, err := os.Stat(p)
		if err == nil {
			if time.Since(fi.ModTime()) <= c.ttl {
				b, err := os.ReadFile(p)
				if err == nil {
					// hydrate memory if enabled
					if c.mode == CacheBoth {
						c.mu.Lock()
						c.mem[key] = memEntry{data: b, expiresAt: time.Now().Add(c.ttl)}
						c.mu.Unlock()
					}
					return b, true
				}
			}
			// expired on disk
			_ = os.Remove(p)
		}
	}
	return nil, false
}

// Set stores bytes in the configured cache layers.
func (c *Cache) Set(key string, data []byte) error {
	if c == nil || c.mode == CacheOff {
		return nil
	}
	now := time.Now()
	if c.mode == CacheMemory || c.mode == CacheBoth {
		c.mu.Lock()
		// Optional tiny eviction once over capacity (FIFO-ish)
		if c.maxMem > 0 && len(c.mem) >= c.maxMem {
			for k := range c.mem {
				delete(c.mem, k)
				break
			}
		}
		c.mem[key] = memEntry{data: append([]byte(nil), data...), expiresAt: now.Add(c.ttl)}
		c.mu.Unlock()
	}
	if (c.mode == CacheDisk || c.mode == CacheBoth) && c.dir != "" {
		p := c.diskPathFor(key)
		if err := os.WriteFile(p, data, 0o644); err != nil {
			return err
		}
		_ = os.Chtimes(p, now, now)
	}
	return nil
}

// Delete removes a specific key from cache.
func (c *Cache) Delete(key string) error {
	if c == nil {
		return nil
	}
	if c.mode == CacheMemory || c.mode == CacheBoth {
		c.mu.Lock()
		delete(c.mem, key)
		c.mu.Unlock()
	}
	if (c.mode == CacheDisk || c.mode == CacheBoth) && c.dir != "" {
		p := c.diskPathFor(key)
		_ = os.Remove(p)
	}
	return nil
}

// Clear removes cached entries. If pattern is empty or "*", clears all.
// Pattern uses filepath.Match semantics against the *key* (not filename).
// For disk entries, we don’t store keys on disk; therefore pattern clearing
// will remove all disk files when pattern != "" (best-effort) unless you pass
// exact keys via Delete().
func (c *Cache) Clear(pattern string) error {
	if c == nil {
		return nil
	}
	if pattern == "" || pattern == "*" {
		// memory
		c.mu.Lock()
		c.mem = make(map[string]memEntry)
		c.mu.Unlock()
		// disk
		if (c.mode == CacheDisk || c.mode == CacheBoth) && c.dir != "" {
			_ = filepath.WalkDir(c.dir, func(path string, d fs.DirEntry, err error) error {
				if err == nil && !d.IsDir() {
					_ = os.Remove(path)
				}
				return nil
			})
		}
		return nil
	}
	// key-based memory clear
	c.mu.Lock()
	for k := range c.mem {
		if ok, _ := filepath.Match(pattern, k); ok {
			delete(c.mem, k)
		}
	}
	c.mu.Unlock()
	// Disk: without a persisted key index we can’t selectively match.
	// Fall back to full-disk clear only if the pattern looks like a blanket.
	if strings.ContainsAny(pattern, "*") {
		if (c.mode == CacheDisk || c.mode == CacheBoth) && c.dir != "" {
			_ = filepath.WalkDir(c.dir, func(path string, d fs.DirEntry, err error) error {
				if err == nil && !d.IsDir() {
					_ = os.Remove(path)
				}
				return nil
			})
		}
	}
	return nil
}

// Cleanup removes expired entries from memory and disk.
func (c *Cache) Cleanup() {
	if c == nil {
		return
	}
	now := time.Now()
	if c.mode == CacheMemory || c.mode == CacheBoth {
		c.mu.Lock()
		for k, ent := range c.mem {
			if now.After(ent.expiresAt) {
				delete(c.mem, k)
			}
		}
		c.mu.Unlock()
	}
	if (c.mode == CacheDisk || c.mode == CacheBoth) && c.dir != "" {
		_ = filepath.WalkDir(c.dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			fi, err := os.Stat(path)
			if err == nil && time.Since(fi.ModTime()) > c.ttl {
				_ = os.Remove(path)
			}
			return nil
		})
	}
}

// ---- Package-level default cache (optional) ----

var (
	// Defaults chosen to be conservative.
	defaultCache = NewCache(CacheBoth, 24*time.Hour, filepath.Join(os.TempDir(), "nflreadgo-cache"), 256)
)

// GetCache returns the package-level cache.
func GetCache() *Cache { return defaultCache }

// SetCacheOptions allows changing default cache mode/ttl/dir at runtime.
func SetCacheOptions(mode CacheMode, ttl time.Duration, dir string, maxMemEntries int) error {
	if dir == "" && (mode == CacheDisk || mode == CacheBoth) {
		return errors.New("cache dir required for disk modes")
	}
	if dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create cache dir: %w", err)
		}
	}
	defaultCache = NewCache(mode, ttl, dir, maxMemEntries)
	return nil
}

// ClearCache clears the package-level cache. Empty pattern = clear all.
func ClearCache(pattern string) { _ = defaultCache.Clear(pattern) }
