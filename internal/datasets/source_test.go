package datasets

import "testing"

func TestSeasonPath(t *testing.T) {
	if got := SeasonPath("injuries/injuries", 0); got != "injuries/injuries" {
		t.Fatalf("SeasonPath base: got %q", got)
	}
	if got := SeasonPath("injuries/injuries", 2024); got != "injuries/injuries_2024" {
		t.Fatalf("SeasonPath with year: got %q", got)
	}
}
