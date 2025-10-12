package injuries

import "strconv"

// Injury models a single row from the nflverse injuries dataset.
// JSON tags mirror dataset column names exactly.
// Data dictionary: https://nflreadr.nflverse.com/articles/dictionary_injuries.html

type Injury struct {
	Season                  int    `json:"season"`
	SeasonType              string `json:"season_type"` // REG or POST
	Team                    string `json:"team"`
	Week                    int    `json:"week"`
	GSISID                  string `json:"gsis_id"`
	Position                string `json:"position"`
	FullName                string `json:"full_name"`
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name"`
	ReportPrimaryInjury     string `json:"report_primary_injury"`
	ReportSecondaryInjury   string `json:"report_secondary_injury"`
	ReportStatus            string `json:"report_status"`
	PracticePrimaryInjury   string `json:"practice_primary_injury"`
	PracticeSecondaryInjury string `json:"practice_secondary_injury"`
	PracticeStatus          string `json:"practice_status"`
	DateModified            string `json:"date_modified"` // ISO8601 timestamp
}

// FromMap converts a generic row into a typed Injury.
func FromMap(row map[string]any) Injury {
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
					// numbers to string (gsis_id sometimes numeric)
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

	return Injury{
		Season:                  getI("season"),
		SeasonType:              getS("season_type"),
		Team:                    getS("team"),
		Week:                    getI("week"),
		GSISID:                  getS("gsis_id"),
		Position:                getS("position"),
		FullName:                getS("full_name"),
		FirstName:               getS("first_name"),
		LastName:                getS("last_name"),
		ReportPrimaryInjury:     getS("report_primary_injury"),
		ReportSecondaryInjury:   getS("report_secondary_injury"),
		ReportStatus:            getS("report_status"),
		PracticePrimaryInjury:   getS("practice_primary_injury"),
		PracticeSecondaryInjury: getS("practice_secondary_injury"),
		PracticeStatus:          getS("practice_status"),
		DateModified:            getS("date_modified"),
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

// ToMap converts a typed Injury back to dataset-style keys.
func (x Injury) ToMap() map[string]any {
	return map[string]any{
		"season":                    x.Season,
		"season_type":               x.SeasonType,
		"team":                      x.Team,
		"week":                      x.Week,
		"gsis_id":                   x.GSISID,
		"position":                  x.Position,
		"full_name":                 x.FullName,
		"first_name":                x.FirstName,
		"last_name":                 x.LastName,
		"report_primary_injury":     x.ReportPrimaryInjury,
		"report_secondary_injury":   x.ReportSecondaryInjury,
		"report_status":             x.ReportStatus,
		"practice_primary_injury":   x.PracticePrimaryInjury,
		"practice_secondary_injury": x.PracticeSecondaryInjury,
		"practice_status":           x.PracticeStatus,
		"date_modified":             x.DateModified,
	}
}
