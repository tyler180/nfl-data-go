package snapcounts

type SnapCount struct {
	GameID            string  `json:"game_id"`
	PFRGID            string  `json:"pfr_game_id"`
	Season            int     `json:"season"`
	Gametype          string  `json:"game_type"` // e.g., REG, PRE, POST
	Week              int     `json:"week"`
	Player            string  `json:"player"` // full name
	PlayerID          string  `json:"pfr_player_id"`
	Position          string  `json:"position"`
	Team              string  `json:"team"`
	Opponent          string  `json:"opponent"`
	OffenseSnaps      int     `json:"offensive_snaps"`
	OffensePct        float64 `json:"offense_pct"`
	DefenseSnaps      int     `json:"defensive_snaps"`
	DefensePct        float64 `json:"defense_pct"`
	SpecialTeamsSnaps int     `json:"st_snaps"`
	SpecialTeamsPct   float64 `json:"st_pct"`
	PlayerSnaps       int     `json:"player_snaps"`
	SnapPct           float64 `json:"snap_pct"`
}

// FromMap constructs a SnapCount from a generic map (e.g., from CSV or Parquet row).
// func FromMap(row map[string]any) SnapCount {
// 	getS := func(keys ...string) string {
// 		for _, k := range keys {
// 			if v, ok := row[k]; ok && v != nil {
// 				switch t := v.(type) {
// 				case string:
// 					if t != "" {
// 						return t
// 					}
// 				case []byte:
// 					if len(t) > 0 {
// 						return string(t)
// 					}
// 				}
// 			}
// 		}
// 		return ""
// 	}
// 	getI := func(keys ...string) int {
// 		for _, k := range keys {
// 			if v, ok := row[k]; ok && v != nil {
// 				switch t := v.(type) {
// 				case int:
// 					return t
// 				case int32:
// 					return int(t)
// 				case int64:
// 					return int(t)
// 				case float64:
// 					return int(t)
// 				}
// 			}
// 		}
// 		return 0
// 	}
// 	getF := func(keys ...string) float64 {
// 		for _, k := range keys {
// 			if v, ok := row[k]; ok && v != nil {
// 				switch t := v.(type) {
// 				case float64:
// 					return t
// 				case float32:
// 					return float64(t)
// 				case int:
// 					return float64(t)
// 				case int32:
// 					return float64(t)
// 				case int64:
// 					return float64(t)
// 				}
// 			}
// 		}
// 		return 0
// 	}

// 	sc := SnapCount{
// 		GameID:            getS("game_id"),
// 		PFRGID:            getS("pfr_game_id"),
// 		Season:            getI("season"),
// 		Gametype:          getS("game_type"),
// 		Week:              getI("week"),
// 		Player:            getS("player"),
// 		PlayerID:          getS("pfr_player_id"),
// 		Position:          getS("position"),
// 		Team:              getS("team"),
// 		Opponent:          getS("opponent"),
// 		OffenseSnaps:      getI("offensive_snaps"),
// 		OffensePct:        getF("offense_pct"),
// 		DefenseSnaps:      getI("defensive_snaps"),
// 		DefensePct:        getF("defense_pct"),
// 		SpecialTeamsSnaps: getI("st_snaps"),
// 		SpecialTeamsPct:   getF("st_pct"),
// 	}
// 	return sc
// }
