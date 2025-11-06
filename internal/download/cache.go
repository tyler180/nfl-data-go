package download

import (
	"io"
	"time"
)

type Metadata struct {
	ETag          string
	LastModified  time.Time
	ContentLength int64
}

type Cache interface {
	// Validators returns ETag/Last-Modified to use for conditional GET.
	Validators(url string) (etag string, lastMod time.Time, ok bool)

	// StoreStream writes the response stream to cache and returns a ReadCloser
	// that the caller reads from (so we only stream once).
	StoreStream(url string, meta Metadata, body io.ReadCloser) (io.ReadCloser, io.Closer)

	// Open returns a cached object body and its metadata.
	Open(url string) (io.ReadCloser, Metadata, error)
}
