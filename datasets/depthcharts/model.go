package depthcharts

import "strconv"

// DepthChart models a single row from the nflverse depth charts dataset.
// JSON tags mirror dataset column names. The fields below cover the most
// commonly used identifiers and chart attributes; extend as needed.
// Data dictionary: https://nflreadr.nflverse.com/articles/dictionary_depth_charts.html

type DepthChart struct {
	Season             int    `json:"season"`
	Week               int    `json:"week"`
	Team               string `json:"team"`
	Position           string `json:"position"`
	Depth              int    `json:"depth"`                // chart order (1 = starter)
	DepthChartPosition string `json:"depth_chart_position"` // e.g., WR1, WR2, etc.

	PlayerID     string `json:"player_id"`
	FullName     string `json:"full_name"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	JerseyNumber int    `json:"jersey_number"`
	Status       string `json:"status"`

	GSISID     string `json:"gsis_id"`
	PFRID      string `json:"pfr_id"`
	ESPNID     int    `json:"espn_id"`
	PFFID      int    `json:"pff_id"`
	RotowireID int    `json:"rotowire_id"`
	YahooID    int    `json:"yahoo_id"`

	HeadshotURL string `json:"headshot_url"`
}

// FromMap converts a generic row into a typed DepthChart.
func FromMap(row map[string]any) DepthChart {
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

	return DepthChart{
		Season:             getI("season"),
		Week:               getI("week"),
		Team:               getS("team"),
		Position:           getS("position"),
		Depth:              getI("depth"),
		DepthChartPosition: getS("depth_chart_position", "chart_position"),

		PlayerID:     getS("player_id"),
		FullName:     getS("full_name", "player_name"),
		FirstName:    getS("first_name"),
		LastName:     getS("last_name"),
		JerseyNumber: getI("jersey_number", "jersey"),
		Status:       getS("status"),

		GSISID:     getS("gsis_id"),
		PFRID:      getS("pfr_id"),
		ESPNID:     getI("espn_id"),
		PFFID:      getI("pff_id"),
		RotowireID: getI("rotowire_id"),
		YahooID:    getI("yahoo_id"),

		HeadshotURL: getS("headshot_url"),
	}
}

// ToMap converts a DepthChart back to dataset-style keys.
func (d DepthChart) ToMap() map[string]any {
	return map[string]any{
		"season":               d.Season,
		"week":                 d.Week,
		"team":                 d.Team,
		"position":             d.Position,
		"depth":                d.Depth,
		"depth_chart_position": d.DepthChartPosition,
		"player_id":            d.PlayerID,
		"full_name":            d.FullName,
		"first_name":           d.FirstName,
		"last_name":            d.LastName,
		"jersey_number":        d.JerseyNumber,
		"status":               d.Status,
		"gsis_id":              d.GSISID,
		"pfr_id":               d.PFRID,
		"espn_id":              d.ESPNID,
		"pff_id":               d.PFFID,
		"rotowire_id":          d.RotowireID,
		"yahoo_id":             d.YahooID,
		"headshot_url":         d.HeadshotURL,
	}
}
