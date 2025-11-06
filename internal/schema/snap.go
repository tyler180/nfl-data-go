package schema

type SnapCount struct {
	Season       int
	Week         int
	GameID       string
	PlayerID     string
	Team         string
	OffenseSnaps int
	PlayerSnaps  int
	SnapPct      float64
}
