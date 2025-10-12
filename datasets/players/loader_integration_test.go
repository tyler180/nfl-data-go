//go:build integration
// +build integration

package players

import "testing"

func TestLoad_Players_Integration(t *testing.T) {
	rows, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if len(rows) == 0 {
		t.Fatalf("expected non-empty players rows")
	}
}
