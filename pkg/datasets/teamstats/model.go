package teamstats

import "strconv"

// TeamStat models a single row in the nflverse team summary stats dataset
// produced by nflfastR::calculate_stats(stat_type = "team").
// It intentionally mirrors the player stats naming where possible.
//
// NOTE: The dataset supports multiple summary levels (week, reg, post, reg+post).
// This struct works for all of them; for season summaries, Week will be 0.
//
// Data dictionary: https://nflreadr.nflverse.com/articles/dictionary_team_stats.html
type TeamStat struct {
	Season     int    `json:"season"`
	Week       int    `json:"week"`
	SeasonType string `json:"season_type"` // REG, POST, or REG+POST
	Team       string `json:"team"`

	// Passing (official-style box score fields)
	Completions            int `json:"completions"`
	Attempts               int `json:"attempts"`
	PassingYards           int `json:"passing_yards"`
	PassingTDs             int `json:"passing_tds"`
	PassingInterceptions   int `json:"passing_interceptions"`
	SacksSuffered          int `json:"sacks_suffered"`
	SackYardsLost          int `json:"sack_yards_lost"`
	PassingAirYards        int `json:"passing_air_yards"`
	PassingYardsAfterCatch int `json:"passing_yards_after_catch"`
	PassingFirstDowns      int `json:"passing_first_downs"`

	// Rushing
	Carries           int `json:"carries"`
	RushingYards      int `json:"rushing_yards"`
	RushingTDs        int `json:"rushing_tds"`
	RushingFirstDowns int `json:"rushing_first_downs"`

	// Receiving
	Targets                  int `json:"targets"`
	Receptions               int `json:"receptions"`
	ReceivingYards           int `json:"receiving_yards"`
	ReceivingTDs             int `json:"receiving_tds"`
	ReceivingFirstDowns      int `json:"receiving_first_downs"`
	ReceivingAirYards        int `json:"receiving_air_yards"`
	ReceivingYardsAfterCatch int `json:"receiving_yards_after_catch"`

	// Ball security
	Fumbles     int `json:"fumbles"`
	FumblesLost int `json:"fumbles_lost"`

	// Kicking / ST (selected common fields)
	FieldGoalsMade       int `json:"field_goals_made"`
	FieldGoalsAttempted  int `json:"field_goals_attempts"`
	ExtraPointsMade      int `json:"extra_points_made"`
	ExtraPointsAttempted int `json:"extra_points_attempts"`
	Punts                int `json:"punts"`
	PuntYards            int `json:"punt_yards"`
}

// FromMap converts a generic row map into a typed TeamStat.
func FromMap(row map[string]any) TeamStat {
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

	return TeamStat{
		Season:     getI("season"),
		Week:       getI("week"),
		SeasonType: getS("season_type"),
		Team:       getS("team"),

		Completions:            getI("completions"),
		Attempts:               getI("attempts"),
		PassingYards:           getI("passing_yards"),
		PassingTDs:             getI("passing_tds"),
		PassingInterceptions:   getI("passing_interceptions"),
		SacksSuffered:          getI("sacks_suffered"),
		SackYardsLost:          getI("sack_yards_lost"),
		PassingAirYards:        getI("passing_air_yards"),
		PassingYardsAfterCatch: getI("passing_yards_after_catch"),
		PassingFirstDowns:      getI("passing_first_downs"),

		Carries:           getI("carries"),
		RushingYards:      getI("rushing_yards"),
		RushingTDs:        getI("rushing_tds"),
		RushingFirstDowns: getI("rushing_first_downs"),

		Targets:                  getI("targets"),
		Receptions:               getI("receptions"),
		ReceivingYards:           getI("receiving_yards"),
		ReceivingTDs:             getI("receiving_tds"),
		ReceivingFirstDowns:      getI("receiving_first_downs"),
		ReceivingAirYards:        getI("receiving_air_yards"),
		ReceivingYardsAfterCatch: getI("receiving_yards_after_catch"),

		Fumbles:     getI("fumbles"),
		FumblesLost: getI("fumbles_lost"),

		FieldGoalsMade:       getI("field_goals_made"),
		FieldGoalsAttempted:  getI("field_goals_attempts"),
		ExtraPointsMade:      getI("extra_points_made"),
		ExtraPointsAttempted: getI("extra_points_attempts"),
		Punts:                getI("punts"),
		PuntYards:            getI("punt_yards"),
	}
}

// ToMap converts a TeamStat back to dataset-style keys.
func (t TeamStat) ToMap() map[string]any {
	return map[string]any{
		"season":      t.Season,
		"week":        t.Week,
		"season_type": t.SeasonType,
		"team":        t.Team,

		"completions":               t.Completions,
		"attempts":                  t.Attempts,
		"passing_yards":             t.PassingYards,
		"passing_tds":               t.PassingTDs,
		"passing_interceptions":     t.PassingInterceptions,
		"sacks_suffered":            t.SacksSuffered,
		"sack_yards_lost":           t.SackYardsLost,
		"passing_air_yards":         t.PassingAirYards,
		"passing_yards_after_catch": t.PassingYardsAfterCatch,
		"passing_first_downs":       t.PassingFirstDowns,

		"carries":             t.Carries,
		"rushing_yards":       t.RushingYards,
		"rushing_tds":         t.RushingTDs,
		"rushing_first_downs": t.RushingFirstDowns,

		"targets":                     t.Targets,
		"receptions":                  t.Receptions,
		"receiving_yards":             t.ReceivingYards,
		"receiving_tds":               t.ReceivingTDs,
		"receiving_first_downs":       t.ReceivingFirstDowns,
		"receiving_air_yards":         t.ReceivingAirYards,
		"receiving_yards_after_catch": t.ReceivingYardsAfterCatch,

		"fumbles":      t.Fumbles,
		"fumbles_lost": t.FumblesLost,

		"field_goals_made":      t.FieldGoalsMade,
		"field_goals_attempts":  t.FieldGoalsAttempted,
		"extra_points_made":     t.ExtraPointsMade,
		"extra_points_attempts": t.ExtraPointsAttempted,
		"punts":                 t.Punts,
		"punt_yards":            t.PuntYards,
	}
}
