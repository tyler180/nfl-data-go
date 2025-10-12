package rosters

import "strconv"

// Roster models a single row in the nflverse season-level rosters dataset.
// JSON tags mirror dataset column names from the nflreadr data dictionary.
// https://nflreadr.nflverse.com/articles/dictionary_rosters.html
type Roster struct {
	Season                int    `json:"season"`
	Team                  string `json:"team"`
	Position              string `json:"position"`
	DepthChartPosition    string `json:"depth_chart_position"`
	JerseyNumber          int    `json:"jersey_number"`
	Status                string `json:"status"`
	FullName              string `json:"full_name"`
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	BirthDate             string `json:"birth_date"` // YYYY-MM-DD
	HeightIn              int    `json:"height"`     // inches
	WeightLb              int    `json:"weight"`     // pounds
	College               string `json:"college"`
	HighSchool            string `json:"high_school"`
	GSISID                string `json:"gsis_id"`
	ESPNID                int    `json:"espn_id"`
	SportradarID          string `json:"sportradar_id"`
	YahooID               int    `json:"yahoo_id"`
	RotowireID            int    `json:"rotowire_id"`
	PFFID                 int    `json:"pff_id"`
	PFRID                 string `json:"pfr_id"`
	FantasyDataID         int    `json:"fantasy_data_id"`
	SleeperID             string `json:"sleeper_id"`
	YearsExp              int    `json:"years_exp"`
	HeadshotURL           string `json:"headshot_url"`
	NGSPosition           string `json:"ngs_position"`
	Week                  int    `json:"week"`
	GameType              string `json:"game_type"`
	StatusDescriptionAbbr string `json:"status_description_abbr"`
	FootballName          string `json:"football_name"`
	ESBID                 string `json:"esb_id"`
	GSISItID              int    `json:"gsis_it_id"`
	SmartID               string `json:"smart_id"`
	EntryYear             int    `json:"entry_year"`
	RookieYear            int    `json:"rookie_year"`
	DraftClub             string `json:"draft_club"`
	DraftNumber           int    `json:"draft_number"`
}

// FromMap converts a generic row map into a typed Roster.
func FromMap(row map[string]any) Roster {
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

	return Roster{
		Season:                getI("season"),
		Team:                  getS("team"),
		Position:              getS("position"),
		DepthChartPosition:    getS("depth_chart_position"),
		JerseyNumber:          getI("jersey_number"),
		Status:                getS("status"),
		FullName:              getS("full_name"),
		FirstName:             getS("first_name"),
		LastName:              getS("last_name"),
		BirthDate:             getS("birth_date"),
		HeightIn:              getI("height"),
		WeightLb:              getI("weight"),
		College:               getS("college"),
		HighSchool:            getS("high_school"),
		GSISID:                getS("gsis_id"),
		ESPNID:                getI("espn_id"),
		SportradarID:          getS("sportradar_id"),
		YahooID:               getI("yahoo_id"),
		RotowireID:            getI("rotowire_id"),
		PFFID:                 getI("pff_id"),
		PFRID:                 getS("pfr_id"),
		FantasyDataID:         getI("fantasy_data_id"),
		SleeperID:             getS("sleeper_id"),
		YearsExp:              getI("years_exp"),
		HeadshotURL:           getS("headshot_url"),
		NGSPosition:           getS("ngs_position"),
		Week:                  getI("week"),
		GameType:              getS("game_type"),
		StatusDescriptionAbbr: getS("status_description_abbr"),
		FootballName:          getS("football_name"),
		ESBID:                 getS("esb_id"),
		GSISItID:              getI("gsis_it_id"),
		SmartID:               getS("smart_id"),
		EntryYear:             getI("entry_year"),
		RookieYear:            getI("rookie_year"),
		DraftClub:             getS("draft_club"),
		DraftNumber:           getI("draft_number"),
	}
}

// ToMap converts a Roster back to dataset-style keys.
func (r Roster) ToMap() map[string]any {
	return map[string]any{
		"season":                  r.Season,
		"team":                    r.Team,
		"position":                r.Position,
		"depth_chart_position":    r.DepthChartPosition,
		"jersey_number":           r.JerseyNumber,
		"status":                  r.Status,
		"full_name":               r.FullName,
		"first_name":              r.FirstName,
		"last_name":               r.LastName,
		"birth_date":              r.BirthDate,
		"height":                  r.HeightIn,
		"weight":                  r.WeightLb,
		"college":                 r.College,
		"high_school":             r.HighSchool,
		"gsis_id":                 r.GSISID,
		"espn_id":                 r.ESPNID,
		"sportradar_id":           r.SportradarID,
		"yahoo_id":                r.YahooID,
		"rotowire_id":             r.RotowireID,
		"pff_id":                  r.PFFID,
		"pfr_id":                  r.PFRID,
		"fantasy_data_id":         r.FantasyDataID,
		"sleeper_id":              r.SleeperID,
		"years_exp":               r.YearsExp,
		"headshot_url":            r.HeadshotURL,
		"ngs_position":            r.NGSPosition,
		"week":                    r.Week,
		"game_type":               r.GameType,
		"status_description_abbr": r.StatusDescriptionAbbr,
		"football_name":           r.FootballName,
		"esb_id":                  r.ESBID,
		"gsis_it_id":              r.GSISItID,
		"smart_id":                r.SmartID,
		"entry_year":              r.EntryYear,
		"rookie_year":             r.RookieYear,
		"draft_club":              r.DraftClub,
		"draft_number":            r.DraftNumber,
	}
}
