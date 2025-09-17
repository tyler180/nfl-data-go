package players

// Player models the nflverse players dataset as a typed struct.
// Only common, stable fields are included here; you can add more as needed.
// Tags support JSON round-trips; CSV headers are documented in comments.
//
// Source reference: nflverse "players" dataset.
// (Fields intentionally use snake_case JSON to mirror dataset columns.)
type Player struct {
	// PlayerID         string `json:"pfr_id"` // unique id (e.g., pfr, gsis, or merged id)
	GSISID           string `json:"gsis_id"`
	PFRID            string `json:"pfr_id"`
	ESPNID           string `json:"espn_id"`
	PFFID            string `json:"pff_id"`
	ESDBID           string `json:"esb_id"`       // FootballDB unique id
	FullName         string `json:"display_name"` // full name
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	ShortName        string `json:"short_name"`    // e.g., "T.Brady"
	FootballName     string `json:"football_name"` // e.g., "T.Brady"
	Position         string `json:"position"`
	PositionGroup    string `json:"position_group"`
	PFFPosition      string `json:"pff_position"`
	PFFStatus        string `json:"pff_status"`
	NGSStatus        string `json:"ngs_status"`
	DraftTeam        string `json:"draft_team"`
	NGSPosition      string `json:"ngs_position"`
	NGSPositionGroup string `json:"ngs_position_group"`
	LatestTeam       string `json:"latest_team"`
	Status           string `json:"status"`
	Height           int    `json:"height"`     // inches
	Weight           int    `json:"weight"`     // pounds
	BirthDate        string `json:"birth_date"` // YYYY-MM-DD
	College          string `json:"college_name"`
	DraftYear        int    `json:"draft_year"`
	DraftRound       int    `json:"draft_round"`
	DraftPick        int    `json:"draft_pick"`
	YearsExp         int    `json:"years_of_experience"`
}

// FromMap performs a best-effort mapping from a generic row (map[string]any)
// into a Player. Missing or malformed fields are left at zero values.
func FromMap(row map[string]any) Player {
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
					// try simple Atoi
					if n, err := atoiSafe(t); err == nil {
						return n
					}
				}
			}
		}
		return 0
	}

	p := Player{
		// PlayerID:      getS("player_id", "id", "nfl_id"),
		GSISID:           getS("gsis_id", "gsisid"),
		PFRID:            getS("pfr_id", "pfr"),
		ESPNID:           getS("espn_id"),
		PFFID:            getS("pff_id"),
		FullName:         getS("player_name", "name", "name_full"),
		FirstName:        getS("first_name", "name_first", "firstname"),
		LastName:         getS("last_name", "name_last", "lastname"),
		Position:         getS("position"),
		PositionGroup:    getS("position_group"),
		LatestTeam:       getS("team", "recent_team"),
		Status:           getS("status"),
		Height:           getI("height", "height_in"),
		Weight:           getI("weight", "weight_lb"),
		BirthDate:        getS("birth_date", "birthdate"),
		College:          getS("college", "college_name"),
		DraftYear:        getI("draft_year"),
		DraftRound:       getI("draft_round"),
		DraftPick:        getI("draft_pick"),
		YearsExp:         getI("years_exp", "years_of_experience"),
		ShortName:        getS("name_short", "short_name"),
		FootballName:     getS("name_football", "football_name"),
		PFFPosition:      getS("pff_position"),
		PFFStatus:        getS("pff_status"),
		NGSStatus:        getS("ngs_status"),
		DraftTeam:        getS("draft_team"),
		NGSPosition:      getS("ngs_position"),
		NGSPositionGroup: getS("ngs_position_group"),
		ESDBID:           getS("football_db_id", "esb_id"),
	}
	return p
}

// ToMap converts a Player back into a generic row map matching dataset keys.
func (p Player) ToMap() map[string]any {
	return map[string]any{
		// "player_id":      p.PlayerID,
		"gsis_id":             p.GSISID,
		"pfr_id":              p.PFRID,
		"espn_id":             p.ESPNID,
		"pff_id":              p.PFFID,
		"name_full":           p.FullName,
		"name_first":          p.FirstName,
		"name_last":           p.LastName,
		"position":            p.Position,
		"position_group":      p.PositionGroup,
		"team":                p.LatestTeam,
		"status":              p.Status,
		"height":              p.Height,
		"weight":              p.Weight,
		"birth_date":          p.BirthDate,
		"college":             p.College,
		"draft_year":          p.DraftYear,
		"draft_round":         p.DraftRound,
		"draft_pick":          p.DraftPick,
		"years_of_experience": p.YearsExp,
		"name_short":          p.ShortName,
		"name_football":       p.FootballName,
		"pff_position":        p.PFFPosition,
		"pff_status":          p.PFFStatus,
		"ngs_status":          p.NGSStatus,
		"draft_team":          p.DraftTeam,
		"ngs_position":        p.NGSPosition,
		"ngs_position_group":  p.NGSPositionGroup,
		"esb_id":              p.ESDBID,
	}
}

// atoiSafe trims and converts a decimal string to int.
func atoiSafe(s string) (int, error) {
	// micro helper without pulling strconv everywhere here
	var n int
	neg := false
	for i, r := range s {
		if i == 0 && (r == '+' || r == '-') {
			neg = r == '-'
			continue
		}
		if r < '0' || r > '9' {
			return 0, fmtErr(s)
		}
		n = n*10 + int(r-'0')
	}
	if neg {
		n = -n
	}
	return n, nil
}

func fmtErr(s string) error { return &parseIntError{s: s} }

type parseIntError struct{ s string }

func (e *parseIntError) Error() string { return "invalid int: " + e.s }
