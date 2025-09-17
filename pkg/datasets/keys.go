package datasets

// Key identifies a supported dataset. Keep names stable and lowercase.
type Key string

const (
	Players         Key = "players"
	SnapCounts      Key = "snapcounts"
	PlayerStats     Key = "playerstats"
	Rosters         Key = "rosters"
	RostersWeekly   Key = "rosters_weekly"
	TeamStatsWeekly Key = "teamstats_week"
	DepthCharts     Key = "depth_charts"
	Injuries        Key = "injuries"
)

// pathByKey maps dataset keys to their nflverse-data repo paths.
var pathByKey = map[Key]string{
	Players:         "players/players",
	SnapCounts:      "snap_counts/snap_counts",
	PlayerStats:     "player_stats/player_stats",
	Rosters:         "rosters/rosters",
	RostersWeekly:   "weekly_rosters/weekly_rosters",
	TeamStatsWeekly: "stats_team/stats_team_week",
	DepthCharts:     "depth_charts/depth_charts",
	Injuries:        "injuries/injuries",
}
