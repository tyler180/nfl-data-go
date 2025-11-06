package download

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"
)

type fsCache struct {
	dir string
	ttl time.Duration
}

type sidecar struct {
	ETag         string    `json:"etag,omitempty"`
	LastModified time.Time `json:"last_modified,omitempty"`
	Size         int64     `json:"size,omitempty"`
	SavedAt      time.Time `json:"saved_at"`
}

// NewFSCache returns a filesystem-backed Cache. Files are stored as two
// siblings: <hash>.data and <hash>.json containing ETag/Last-Modified/TTL info.
func NewFSCache(dir string, ttl time.Duration) Cache {
	_ = os.MkdirAll(dir, 0o755)
	return &fsCache{dir: dir, ttl: ttl}
}

func (c *fsCache) Validators(url string) (string, time.Time, bool) {
	meta, ok := c.readMeta(url)
	if !ok {
		return "", time.Time{}, false
	}
	// TTL: if expired, fetch fresh (skip validators to force 200)
	if c.ttl > 0 && time.Since(meta.SavedAt) > c.ttl {
		return "", time.Time{}, false
	}
	return meta.ETag, meta.LastModified, true
}

func (c *fsCache) StoreStream(url string, m Metadata, body io.ReadCloser) (io.ReadCloser, io.Closer) {
	defer body.Close()
	b, err := io.ReadAll(body)
	if err != nil {
		return io.NopCloser(bytes.NewReader(nil)), io.NopCloser(nil)
	}
	_ = os.MkdirAll(c.dir, 0o755)

	base := c.base(url)
	_ = os.WriteFile(base+".data", b, 0o644)

	sc := sidecar{
		ETag:         m.ETag,
		LastModified: m.LastModified,
		Size:         int64(len(b)),
		SavedAt:      time.Now().UTC(),
	}
	if j, err := json.Marshal(sc); err == nil {
		_ = os.WriteFile(base+".json", j, 0o644)
	}
	return io.NopCloser(bytes.NewReader(b)), io.NopCloser(nil)
}

func (c *fsCache) Open(url string) (io.ReadCloser, Metadata, error) {
	base := c.base(url)
	b, err := os.ReadFile(base + ".data")
	if err != nil {
		return nil, Metadata{}, err
	}
	meta := Metadata{}
	if j, err := os.ReadFile(base + ".json"); err == nil {
		var sc sidecar
		if json.Unmarshal(j, &sc) == nil {
			meta.ETag = sc.ETag
			meta.LastModified = sc.LastModified
			meta.ContentLength = sc.Size
		}
	}
	return io.NopCloser(bytes.NewReader(b)), meta, nil
}

func (c *fsCache) base(url string) string {
	h := sha1.Sum([]byte(url))
	return filepath.Join(c.dir, hex.EncodeToString(h[:]))
}

func (c *fsCache) readMeta(url string) (sidecar, bool) {
	j, err := os.ReadFile(c.base(url) + ".json")
	if err != nil {
		return sidecar{}, false
	}
	var sc sidecar
	if err := json.Unmarshal(j, &sc); err != nil {
		return sidecar{}, false
	}
	// TTL gate (also enforced in Validators)
	if c.ttl > 0 && time.Since(sc.SavedAt) > c.ttl {
		return sidecar{}, false
	}
	return sc, true
}
