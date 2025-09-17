package ffplayerids

import "strconv"

// FFPlayerID models a single row from DynastyProcess' fantasy player IDs table.
// Fields follow the nflreadr data dictionary for ff_playerids.
// https://nflreadr.nflverse.com/articles/dictionary_ff_playerids.html

type FFPlayerID struct {
	// Core IDs
	MFLID         string `json:"mfl_id"`
	SportradarID  string `json:"sportradar_id"`
	FantasyProsID string `json:"fantasypros_id"`
	GSISID        string `json:"gsis_id"`
	PFFID         string `json:"pff_id"`
	SleeperID     string `json:"sleeper_id"`
	NFLID         string `json:"nfl_id"`
	ESPNID        string `json:"espn_id"`
	YahooID       string `json:"yahoo_id"`
	FleaflickerID string `json:"fleaflicker_id"`
	CBSID         string `json:"cbs_id"`
	RotowireID    string `json:"rotowire_id"`
	RotoworldID   string `json:"rotoworld_id"`
	KTCID         string `json:"ktc_id"`
	PFRID         string `json:"pfr_id"`
	CFBRefID      string `json:"cfbref_id"`
	StatsID       string `json:"stats_id"`
	StatsGlobalID string `json:"stats_global_id"`
	FantasyDataID string `json:"fantasy_data_id"`
	SwishID       string `json:"swish_id"`

	// Descriptive
	Name            string  `json:"name"`
	MergeName       string  `json:"merge_name"`
	Position        string  `json:"position"`
	Team            string  `json:"team"`
	Birthdate       string  `json:"birthdate"` // YYYY-MM-DD
	Age             float64 `json:"age"`
	DraftYear       int     `json:"draft_year"`
	DraftRound      int     `json:"draft_round"`
	DraftPick       int     `json:"draft_pick"`
	DraftOverallStr string  `json:"draft_ovr"`
	TwitterUsername string  `json:"twitter_username"`
	HeightIn        int     `json:"height"`
	WeightLb        int     `json:"weight"`
	College         string  `json:"college"`
	DBSeason        int     `json:"db_season"`
}

// FromMap converts a generic row map into a typed FFPlayerID.
func FromMap(row map[string]any) FFPlayerID {
	getS := func(keys ...string) string {
		for _, k := range keys {
			if v, ok := row[k]; ok && v != nil {
				switch t := v.(type) {
				case string:
					if t != "" {
						return t
					}
				case []byte:
					if len(t) > 0 {
						return string(t)
					}
				default:
					// accept numbers as strings where IDs are numeric
					s := toString(v)
					if s != "" {
						return s
					}
				}
			}
		}
		return ""
	}
	getI := func(keys ...string) int {
		for _, k := range keys {
			if v, ok := row[k]; ok && v != nil {
				switch t := v.(type) {
				case int:
					return t
				case int32:
					return int(t)
				case int64:
					return int(t)
				case float64:
					return int(t)
				case float32:
					return int(t)
				case string:
					if n, err := strconv.Atoi(t); err == nil {
						return n
					}
				}
			}
		}
		return 0
	}
	getF := func(keys ...string) float64 {
		for _, k := range keys {
			if v, ok := row[k]; ok && v != nil {
				switch t := v.(type) {
				case float64:
					return t
				case float32:
					return float64(t)
				case int:
					return float64(t)
				case int64:
					return float64(t)
				case string:
					if f, err := strconv.ParseFloat(t, 64); err == nil {
						return f
					}
				}
			}
		}
		return 0
	}

	return FFPlayerID{
		MFLID:         getS("mfl_id"),
		SportradarID:  getS("sportradar_id", "sportsdata_id"),
		FantasyProsID: getS("fantasypros_id"),
		GSISID:        getS("gsis_id"),
		PFFID:         getS("pff_id"),
		SleeperID:     getS("sleeper_id"),
		NFLID:         getS("nfl_id"),
		ESPNID:        getS("espn_id"),
		YahooID:       getS("yahoo_id"),
		FleaflickerID: getS("fleaflicker_id"),
		CBSID:         getS("cbs_id"),
		RotowireID:    getS("rotowire_id"),
		RotoworldID:   getS("rotoworld_id"),
		KTCID:         getS("ktc_id"),
		PFRID:         getS("pfr_id"),
		CFBRefID:      getS("cfbref_id"),
		StatsID:       getS("stats_id"),
		StatsGlobalID: getS("stats_global_id"),
		FantasyDataID: getS("fantasy_data_id"),
		SwishID:       getS("swish_id"),

		Name:            getS("name"),
		MergeName:       getS("merge_name"),
		Position:        getS("position"),
		Team:            getS("team"),
		Birthdate:       getS("birthdate"),
		Age:             getF("age"),
		DraftYear:       getI("draft_year"),
		DraftRound:      getI("draft_round"),
		DraftPick:       getI("draft_pick"),
		DraftOverallStr: getS("draft_ovr"),
		TwitterUsername: getS("twitter_username"),
		HeightIn:        getI("height"),
		WeightLb:        getI("weight"),
		College:         getS("college"),
		DBSeason:        getI("db_season"),
	}
}

func toString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	case float64:
		return strconv.FormatInt(int64(t), 10)
	default:
		return ""
	}
}

// ToMap converts FFPlayerID back to dataset-style keys.
func (p FFPlayerID) ToMap() map[string]any {
	return map[string]any{
		"mfl_id":          p.MFLID,
		"sportradar_id":   p.SportradarID,
		"fantasypros_id":  p.FantasyProsID,
		"gsis_id":         p.GSISID,
		"pff_id":          p.PFFID,
		"sleeper_id":      p.SleeperID,
		"nfl_id":          p.NFLID,
		"espn_id":         p.ESPNID,
		"yahoo_id":        p.YahooID,
		"fleaflicker_id":  p.FleaflickerID,
		"cbs_id":          p.CBSID,
		"rotowire_id":     p.RotowireID,
		"rotoworld_id":    p.RotoworldID,
		"ktc_id":          p.KTCID,
		"pfr_id":          p.PFRID,
		"cfbref_id":       p.CFBRefID,
		"stats_id":        p.StatsID,
		"stats_global_id": p.StatsGlobalID,
		"fantasy_data_id": p.FantasyDataID,
		"swish_id":        p.SwishID,

		"name":             p.Name,
		"merge_name":       p.MergeName,
		"position":         p.Position,
		"team":             p.Team,
		"birthdate":        p.Birthdate,
		"age":              p.Age,
		"draft_year":       p.DraftYear,
		"draft_round":      p.DraftRound,
		"draft_pick":       p.DraftPick,
		"draft_ovr":        p.DraftOverallStr,
		"twitter_username": p.TwitterUsername,
		"height":           p.HeightIn,
		"weight":           p.WeightLb,
		"college":          p.College,
		"db_season":        p.DBSeason,
	}
}
