//go:build integration
// +build integration

package ffplayerids

import (
	"context"
	"strings"
	"testing"
)

func TestLoad_FFPlayerIDs_Integration(ctx context.Context, t *testing.T) {
	rows, err := Load(ctx)
	if err != nil {
		// If the downloader hasn't been configured for DynastyProcess/data yet,
		// skip the test with a clear hint rather than failing the whole suite.
		if strings.Contains(err.Error(), "unknown repository") {
			t.Skipf("skipping: %v (configure internal/download to support DynastyProcess/data)", err)
		}
		t.Fatalf("Load() error: %v", err)
	}
	if len(rows) == 0 {
		t.Fatalf("expected non-empty ffplayerids rows")
	}
}
