package download

import (
	"bytes"
	"io"
	"sync"
	"time"
)

type memCache struct {
	ttl  time.Duration
	mu   sync.RWMutex
	data map[string]struct {
		b    []byte
		meta Metadata
		ts   time.Time
	}
}

// NewMemCache returns an in-memory Cache with a simple TTL.
func NewMemCache(ttl time.Duration) Cache {
	return &memCache{
		ttl: ttl,
		data: make(map[string]struct {
			b    []byte
			meta Metadata
			ts   time.Time
		}),
	}
}

func (c *memCache) Validators(url string) (string, time.Time, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.data[url]
	if !ok {
		return "", time.Time{}, false
	}
	if c.ttl > 0 && time.Since(e.ts) > c.ttl {
		return "", time.Time{}, false
	}
	return e.meta.ETag, e.meta.LastModified, true
}

func (c *memCache) StoreStream(url string, meta Metadata, body io.ReadCloser) (io.ReadCloser, io.Closer) {
	defer body.Close()
	b, err := io.ReadAll(body)
	if err != nil {
		return io.NopCloser(bytes.NewReader(nil)), io.NopCloser(nil)
	}
	c.mu.Lock()
	c.data[url] = struct {
		b    []byte
		meta Metadata
		ts   time.Time
	}{append([]byte(nil), b...), meta, time.Now()}
	c.mu.Unlock()
	return io.NopCloser(bytes.NewReader(b)), io.NopCloser(nil)
}

func (c *memCache) Open(url string) (io.ReadCloser, Metadata, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.data[url]
	if !ok {
		return nil, Metadata{}, io.EOF
	}
	if c.ttl > 0 && time.Since(e.ts) > c.ttl {
		return nil, Metadata{}, io.EOF
	}
	return io.NopCloser(bytes.NewReader(e.b)), e.meta, nil
}
