package playerstats

import "strconv"

// PlayerStat models a single row from the nflverse weekly player stats dataset
// ("player_stats" release). This combines offense, defense, and kicking stats
// as described in the nflreadr data dictionary. JSON tags mirror dataset
// column names.
//
// Identity & game context
//   - player_id (gsis)
//   - player_name / player_display_name
//   - position / position_group
//   - headshot_url
//   - season / week / season_type
//   - team / opponent_team
//
// Notes on types:
//   - Counts & yards are generally ints.
//   - Rates/percentages/ratios and EPA are float64.
//   - Some fields are lists encoded as comma-separated strings.
type PlayerStat struct {
	// identity & context
	PlayerID      string `json:"player_id"`
	PlayerName    string `json:"player_name"`
	PlayerDisplay string `json:"player_display_name"`
	Position      string `json:"position"`
	PositionGroup string `json:"position_group"`
	HeadshotURL   string `json:"headshot_url"`
	Season        int    `json:"season"`
	Week          int    `json:"week"`
	SeasonType    string `json:"season_type"` // REG or POST
	Team          string `json:"team"`
	OpponentTeam  string `json:"opponent_team"`

	// passing
	Completions            int     `json:"completions"`
	Attempts               int     `json:"attempts"`
	PassingYards           int     `json:"passing_yards"`
	PassingTDs             int     `json:"passing_tds"`
	PassingInterceptions   int     `json:"passing_interceptions"`
	SacksSuffered          int     `json:"sacks_suffered"`
	SackYardsLost          int     `json:"sack_yards_lost"`
	SackFumbles            int     `json:"sack_fumbles"`
	SackFumblesLost        int     `json:"sack_fumbles_lost"`
	PassingAirYards        int     `json:"passing_air_yards"`
	PassingYardsAfterCatch int     `json:"passing_yards_after_catch"`
	PassingFirstDowns      int     `json:"passing_first_downs"`
	PassingEPA             float64 `json:"passing_epa"`
	PassingCPOE            float64 `json:"passing_cpoe"`
	Passing2PtConversions  int     `json:"passing_2pt_conversions"`
	PACR                   float64 `json:"pacr"`

	// rushing
	Carries               int     `json:"carries"`
	RushingYards          int     `json:"rushing_yards"`
	RushingTDs            int     `json:"rushing_tds"`
	RushingFumbles        int     `json:"rushing_fumbles"`
	RushingFumblesLost    int     `json:"rushing_fumbles_lost"`
	RushingFirstDowns     int     `json:"rushing_first_downs"`
	RushingEPA            float64 `json:"rushing_epa"`
	Rushing2PtConversions int     `json:"rushing_2pt_conversions"`

	// receiving
	Receptions               int     `json:"receptions"`
	Targets                  int     `json:"targets"`
	ReceivingYards           int     `json:"receiving_yards"`
	ReceivingTDs             int     `json:"receiving_tds"`
	ReceivingFumbles         int     `json:"receiving_fumbles"`
	ReceivingFumblesLost     int     `json:"receiving_fumbles_lost"`
	ReceivingAirYards        int     `json:"receiving_air_yards"`
	ReceivingYardsAfterCatch int     `json:"receiving_yards_after_catch"`
	ReceivingFirstDowns      int     `json:"receiving_first_downs"`
	ReceivingEPA             float64 `json:"receiving_epa"`
	Receiving2PtConversions  int     `json:"receiving_2pt_conversions"`
	RACR                     float64 `json:"racr"`
	TargetShare              float64 `json:"target_share"`
	AirYardsShare            float64 `json:"air_yards_share"`
	WOPR                     float64 `json:"wopr"`

	// misc specials & defense
	SpecialTeamsTDs int `json:"special_teams_tds"`

	DefTacklesSolo         int `json:"def_tackles_solo"`
	DefTacklesWithAssist   int `json:"def_tackles_with_assist"`
	DefTackleAssists       int `json:"def_tackle_assists"`
	DefTacklesForLoss      int `json:"def_tackles_for_loss"`
	DefTacklesForLossYards int `json:"def_tackles_for_loss_yards"`
	DefFumblesForced       int `json:"def_fumbles_forced"`
	DefSacks               int `json:"def_sacks"`
	DefSackYards           int `json:"def_sack_yards"`
	DefQBHits              int `json:"def_qb_hits"`
	DefInterceptions       int `json:"def_interceptions"`
	DefInterceptionYards   int `json:"def_interception_yards"`
	DefPassDefended        int `json:"def_pass_defended"`
	DefTDs                 int `json:"def_tds"`
	DefFumbles             int `json:"def_fumbles"`
	DefSafeties            int `json:"def_safeties"`
	MiscYards              int `json:"misc_yards"`
	FumbleRecoveryOwn      int `json:"fumble_recovery_own"`
	FumbleRecoveryYardsOwn int `json:"fumble_recovery_yards_own"`
	FumbleRecoveryOpp      int `json:"fumble_recovery_opp"`
	FumbleRecoveryYardsOpp int `json:"fumble_recovery_yards_opp"`
	FumbleRecoveryTDs      int `json:"fumble_recovery_tds"`
	Penalties              int `json:"penalties"`
	PenaltyYards           int `json:"penalty_yards"`
	PuntReturns            int `json:"punt_returns"`
	PuntReturnYards        int `json:"punt_return_yards"`
	KickoffReturns         int `json:"kickoff_returns"`
	KickoffReturnYards     int `json:"kickoff_return_yards"`

	// kicking
	FGMade            int     `json:"fg_made"`
	FGAtt             int     `json:"fg_att"`
	FGMissed          int     `json:"fg_missed"`
	FGBlocked         int     `json:"fg_blocked"`
	FGLong            int     `json:"fg_long"`
	FGPct             float64 `json:"fg_pct"`
	FGMade0_19        int     `json:"fg_made_0_19"`
	FGMade20_29       int     `json:"fg_made_20_29"`
	FGMade30_39       int     `json:"fg_made_30_39"`
	FGMade40_49       int     `json:"fg_made_40_49"`
	FGMade50_59       int     `json:"fg_made_50_59"`
	FGMade60Plus      int     `json:"fg_made_60_"`
	FGMissed0_19      int     `json:"fg_missed_0_19"`
	FGMissed20_29     int     `json:"fg_missed_20_29"`
	FGMissed30_39     int     `json:"fg_missed_30_39"`
	FGMissed40_49     int     `json:"fg_missed_40_49"`
	FGMissed50_59     int     `json:"fg_missed_50_59"`
	FGMissed60Plus    int     `json:"fg_missed_60_"`
	FGMadeList        string  `json:"fg_made_list"`
	FGMissedList      string  `json:"fg_missed_list"`
	FGBlockedList     string  `json:"fg_blocked_list"`
	FGMadeDistance    int     `json:"fg_made_distance"`
	FGMissedDistance  int     `json:"fg_missed_distance"`
	FGBlockedDistance int     `json:"fg_blocked_distance"`
	PATMade           int     `json:"pat_made"`
	PATAtt            int     `json:"pat_att"`
	PATMissed         int     `json:"pat_missed"`
	PATBlocked        int     `json:"pat_blocked"`
	PATPct            float64 `json:"pat_pct"`
	GWFGMade          int     `json:"gwfg_made"`
	GWFGAtt           int     `json:"gwfg_att"`
	GWFGMissed        int     `json:"gwfg_missed"`
	GWFGBlocked       int     `json:"gwfg_blocked"`
	GWFGRawDistance   int     `json:"gwfg_distance"`

	// fantasy
	FantasyPoints    float64 `json:"fantasy_points"`
	FantasyPointsPPR float64 `json:"fantasy_points_ppr"`
}

// FromMap converts a generic row (map[string]any) to a typed PlayerStat.
// Unknown/malformed fields are left at zero values.
func FromMap(row map[string]any) PlayerStat {
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

	ps := PlayerStat{
		PlayerID:      getS("player_id"),
		PlayerName:    getS("player_name"),
		PlayerDisplay: getS("player_display_name"),
		Position:      getS("position"),
		PositionGroup: getS("position_group"),
		HeadshotURL:   getS("headshot_url"),
		Season:        getI("season"),
		Week:          getI("week"),
		SeasonType:    getS("season_type"),
		Team:          getS("team"),
		OpponentTeam:  getS("opponent_team"),

		Completions:            getI("completions"),
		Attempts:               getI("attempts"),
		PassingYards:           getI("passing_yards"),
		PassingTDs:             getI("passing_tds"),
		PassingInterceptions:   getI("passing_interceptions"),
		SacksSuffered:          getI("sacks_suffered"),
		SackYardsLost:          getI("sack_yards_lost"),
		SackFumbles:            getI("sack_fumbles"),
		SackFumblesLost:        getI("sack_fumbles_lost"),
		PassingAirYards:        getI("passing_air_yards"),
		PassingYardsAfterCatch: getI("passing_yards_after_catch"),
		PassingFirstDowns:      getI("passing_first_downs"),
		PassingEPA:             getF("passing_epa"),
		PassingCPOE:            getF("passing_cpoe"),
		Passing2PtConversions:  getI("passing_2pt_conversions"),
		PACR:                   getF("pacr"),

		Carries:               getI("carries"),
		RushingYards:          getI("rushing_yards"),
		RushingTDs:            getI("rushing_tds"),
		RushingFumbles:        getI("rushing_fumbles"),
		RushingFumblesLost:    getI("rushing_fumbles_lost"),
		RushingFirstDowns:     getI("rushing_first_downs"),
		RushingEPA:            getF("rushing_epa"),
		Rushing2PtConversions: getI("rushing_2pt_conversions"),

		Receptions:               getI("receptions"),
		Targets:                  getI("targets"),
		ReceivingYards:           getI("receiving_yards"),
		ReceivingTDs:             getI("receiving_tds"),
		ReceivingFumbles:         getI("receiving_fumbles"),
		ReceivingFumblesLost:     getI("receiving_fumbles_lost"),
		ReceivingAirYards:        getI("receiving_air_yards"),
		ReceivingYardsAfterCatch: getI("receiving_yards_after_catch"),
		ReceivingFirstDowns:      getI("receiving_first_downs"),
		ReceivingEPA:             getF("receiving_epa"),
		Receiving2PtConversions:  getI("receiving_2pt_conversions"),
		RACR:                     getF("racr"),
		TargetShare:              getF("target_share"),
		AirYardsShare:            getF("air_yards_share"),
		WOPR:                     getF("wopr"),

		SpecialTeamsTDs: getI("special_teams_tds"),

		DefTacklesSolo:         getI("def_tackles_solo"),
		DefTacklesWithAssist:   getI("def_tackles_with_assist"),
		DefTackleAssists:       getI("def_tackle_assists"),
		DefTacklesForLoss:      getI("def_tackles_for_loss"),
		DefTacklesForLossYards: getI("def_tackles_for_loss_yards"),
		DefFumblesForced:       getI("def_fumbles_forced"),
		DefSacks:               getI("def_sacks"),
		DefSackYards:           getI("def_sack_yards"),
		DefQBHits:              getI("def_qb_hits"),
		DefInterceptions:       getI("def_interceptions"),
		DefInterceptionYards:   getI("def_interception_yards"),
		DefPassDefended:        getI("def_pass_defended"),
		DefTDs:                 getI("def_tds"),
		DefFumbles:             getI("def_fumbles"),
		DefSafeties:            getI("def_safeties"),
		MiscYards:              getI("misc_yards"),
		FumbleRecoveryOwn:      getI("fumble_recovery_own"),
		FumbleRecoveryYardsOwn: getI("fumble_recovery_yards_own"),
		FumbleRecoveryOpp:      getI("fumble_recovery_opp"),
		FumbleRecoveryYardsOpp: getI("fumble_recovery_yards_opp"),
		FumbleRecoveryTDs:      getI("fumble_recovery_tds"),
		Penalties:              getI("penalties"),
		PenaltyYards:           getI("penalty_yards"),
		PuntReturns:            getI("punt_returns"),
		PuntReturnYards:        getI("punt_return_yards"),
		KickoffReturns:         getI("kickoff_returns"),
		KickoffReturnYards:     getI("kickoff_return_yards"),

		FGMade:            getI("fg_made"),
		FGAtt:             getI("fg_att"),
		FGMissed:          getI("fg_missed"),
		FGBlocked:         getI("fg_blocked"),
		FGLong:            getI("fg_long"),
		FGPct:             getF("fg_pct"),
		FGMade0_19:        getI("fg_made_0_19"),
		FGMade20_29:       getI("fg_made_20_29"),
		FGMade30_39:       getI("fg_made_30_39"),
		FGMade40_49:       getI("fg_made_40_49"),
		FGMade50_59:       getI("fg_made_50_59"),
		FGMade60Plus:      getI("fg_made_60_"),
		FGMissed0_19:      getI("fg_missed_0_19"),
		FGMissed20_29:     getI("fg_missed_20_29"),
		FGMissed30_39:     getI("fg_missed_30_39"),
		FGMissed40_49:     getI("fg_missed_40_49"),
		FGMissed50_59:     getI("fg_missed_50_59"),
		FGMissed60Plus:    getI("fg_missed_60_"),
		FGMadeList:        getS("fg_made_list"),
		FGMissedList:      getS("fg_missed_list"),
		FGBlockedList:     getS("fg_blocked_list"),
		FGMadeDistance:    getI("fg_made_distance"),
		FGMissedDistance:  getI("fg_missed_distance"),
		FGBlockedDistance: getI("fg_blocked_distance"),
		PATMade:           getI("pat_made"),
		PATAtt:            getI("pat_att"),
		PATMissed:         getI("pat_missed"),
		PATBlocked:        getI("pat_blocked"),
		PATPct:            getF("pat_pct"),
		GWFGMade:          getI("gwfg_made"),
		GWFGAtt:           getI("gwfg_att"),
		GWFGMissed:        getI("gwfg_missed"),
		GWFGBlocked:       getI("gwfg_blocked"),
		GWFGRawDistance:   getI("gwfg_distance"),

		FantasyPoints:    getF("fantasy_points"),
		FantasyPointsPPR: getF("fantasy_points_ppr"),
	}
	return ps
}

// ToMap converts a PlayerStat back into a generic row map with dataset keys.
func (ps PlayerStat) ToMap() map[string]any {
	return map[string]any{
		"player_id":                   ps.PlayerID,
		"player_name":                 ps.PlayerName,
		"player_display_name":         ps.PlayerDisplay,
		"position":                    ps.Position,
		"position_group":              ps.PositionGroup,
		"headshot_url":                ps.HeadshotURL,
		"season":                      ps.Season,
		"week":                        ps.Week,
		"season_type":                 ps.SeasonType,
		"team":                        ps.Team,
		"opponent_team":               ps.OpponentTeam,
		"completions":                 ps.Completions,
		"attempts":                    ps.Attempts,
		"passing_yards":               ps.PassingYards,
		"passing_tds":                 ps.PassingTDs,
		"passing_interceptions":       ps.PassingInterceptions,
		"sacks_suffered":              ps.SacksSuffered,
		"sack_yards_lost":             ps.SackYardsLost,
		"sack_fumbles":                ps.SackFumbles,
		"sack_fumbles_lost":           ps.SackFumblesLost,
		"passing_air_yards":           ps.PassingAirYards,
		"passing_yac":                 ps.PassingYardsAfterCatch,
		"passing_yards_after_catch":   ps.PassingYardsAfterCatch,
		"passing_first_downs":         ps.PassingFirstDowns,
		"passing_epa":                 ps.PassingEPA,
		"passing_cpoe":                ps.PassingCPOE,
		"passing_2pt_conversions":     ps.Passing2PtConversions,
		"pacr":                        ps.PACR,
		"carries":                     ps.Carries,
		"rushing_yards":               ps.RushingYards,
		"rushing_tds":                 ps.RushingTDs,
		"rushing_fumbles":             ps.RushingFumbles,
		"rushing_fumbles_lost":        ps.RushingFumblesLost,
		"rushing_first_downs":         ps.RushingFirstDowns,
		"rushing_epa":                 ps.RushingEPA,
		"rushing_2pt_conversions":     ps.Rushing2PtConversions,
		"receptions":                  ps.Receptions,
		"targets":                     ps.Targets,
		"receiving_yards":             ps.ReceivingYards,
		"receiving_tds":               ps.ReceivingTDs,
		"receiving_fumbles":           ps.ReceivingFumbles,
		"receiving_fumbles_lost":      ps.ReceivingFumblesLost,
		"receiving_air_yards":         ps.ReceivingAirYards,
		"receiving_yards_after_catch": ps.ReceivingYardsAfterCatch,
		"receiving_first_downs":       ps.ReceivingFirstDowns,
		"receiving_epa":               ps.ReceivingEPA,
		"receiving_2pt_conversions":   ps.Receiving2PtConversions,
		"racr":                        ps.RACR,
		"target_share":                ps.TargetShare,
		"air_yards_share":             ps.AirYardsShare,
		"wopr":                        ps.WOPR,
		"special_teams_tds":           ps.SpecialTeamsTDs,
		"def_tackles_solo":            ps.DefTacklesSolo,
		"def_tackles_with_assist":     ps.DefTacklesWithAssist,
		"def_tackle_assists":          ps.DefTackleAssists,
		"def_tackles_for_loss":        ps.DefTacklesForLoss,
		"def_tackles_for_loss_yards":  ps.DefTacklesForLossYards,
		"def_fumbles_forced":          ps.DefFumblesForced,
		"def_sacks":                   ps.DefSacks,
		"def_sack_yards":              ps.DefSackYards,
		"def_qb_hits":                 ps.DefQBHits,
		"def_interceptions":           ps.DefInterceptions,
		"def_interception_yards":      ps.DefInterceptionYards,
		"def_pass_defended":           ps.DefPassDefended,
		"def_tds":                     ps.DefTDs,
		"def_fumbles":                 ps.DefFumbles,
		"def_safeties":                ps.DefSafeties,
		"misc_yards":                  ps.MiscYards,
		"fumble_recovery_own":         ps.FumbleRecoveryOwn,
		"fumble_recovery_yards_own":   ps.FumbleRecoveryYardsOwn,
		"fumble_recovery_opp":         ps.FumbleRecoveryOpp,
		"fumble_recovery_yards_opp":   ps.FumbleRecoveryYardsOpp,
		"fumble_recovery_tds":         ps.FumbleRecoveryTDs,
		"penalties":                   ps.Penalties,
		"penalty_yards":               ps.PenaltyYards,
		"punt_returns":                ps.PuntReturns,
		"punt_return_yards":           ps.PuntReturnYards,
		"kickoff_returns":             ps.KickoffReturns,
		"kickoff_return_yards":        ps.KickoffReturnYards,
		"fg_made":                     ps.FGMade,
		"fg_att":                      ps.FGAtt,
		"fg_missed":                   ps.FGMissed,
		"fg_blocked":                  ps.FGBlocked,
		"fg_long":                     ps.FGLong,
		"fg_pct":                      ps.FGPct,
		"fg_made_0_19":                ps.FGMade0_19,
		"fg_made_20_29":               ps.FGMade20_29,
		"fg_made_30_39":               ps.FGMade30_39,
		"fg_made_40_49":               ps.FGMade40_49,
		"fg_made_50_59":               ps.FGMade50_59,
		"fg_made_60_":                 ps.FGMade60Plus,
		"fg_missed_0_19":              ps.FGMissed0_19,
		"fg_missed_20_29":             ps.FGMissed20_29,
		"fg_missed_30_39":             ps.FGMissed30_39,
		"fg_missed_40_49":             ps.FGMissed40_49,
		"fg_missed_50_59":             ps.FGMissed50_59,
		"fg_missed_60_":               ps.FGMissed60Plus,
		"fg_made_list":                ps.FGMadeList,
		"fg_missed_list":              ps.FGMissedList,
		"fg_blocked_list":             ps.FGBlockedList,
		"fg_made_distance":            ps.FGMadeDistance,
		"fg_missed_distance":          ps.FGMissedDistance,
		"fg_blocked_distance":         ps.FGBlockedDistance,
		"pat_made":                    ps.PATMade,
		"pat_att":                     ps.PATAtt,
		"pat_missed":                  ps.PATMissed,
		"pat_blocked":                 ps.PATBlocked,
		"pat_pct":                     ps.PATPct,
		"gwfg_made":                   ps.GWFGMade,
		"gwfg_att":                    ps.GWFGAtt,
		"gwfg_missed":                 ps.GWFGMissed,
		"gwfg_blocked":                ps.GWFGBlocked,
		"gwfg_distance":               ps.GWFGRawDistance,
		"fantasy_points":              ps.FantasyPoints,
		"fantasy_points_ppr":          ps.FantasyPointsPPR,
	}
}
