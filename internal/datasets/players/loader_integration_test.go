//go:build integration
// +build integration

package players

import (
	"context"
	"testing"
)

func TestLoad_Players_Integration(ctx context.Context, t *testing.T) {
	rows, err := Load(ctx)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if len(rows) == 0 {
		t.Fatalf("expected non-empty players rows")
	}
}
