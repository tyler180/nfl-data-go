package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tyler180/nfl-data-go/pkg/nflreadgo"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	rows, err := nflreadgo.LoadSnapCounts(
		ctx,
		2024, // selector: single season (supports many shapes; see below)
		nflreadgo.WithTimeout(45*time.Second),
		// nflreadgo.WithCache(nflreadgo.CacheFS, "~/.cache/nflreadgo", 24*time.Hour),
		// nflreadgo.WithUserAgent("nflreadgo/0.1 (+github.com/tyler180/nfl-data-go)"),
	)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(rows) && i < 5; i++ {
		r := rows[i]
		fmt.Printf("%d wk%-2d %s %s player=%s snaps=%d/%d (%.1f%%)\n",
			r.Season, r.Week, r.Team, r.GameID, r.PlayerID, r.PlayerSnaps, r.OffenseSnaps, r.SnapPct)
	}
}
