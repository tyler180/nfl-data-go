package nflreadgo

import "github.com/tyler180/nfl-data-go/internal/schema"

// filterBySelection narrows []SnapCount according to 'sel'.
// Supported selectors (same shapes accepted by your loaders):
//   - Seasons{...} or []int (treated as seasons): filter by Season only
//   - []SeasonWeeks: filter by Season AND listed Weeks (empty Weeks = all)
//   - int (single season)
//   - Weeks (applies to current season)
//   - bool(true) => no filtering (“all”)
func filterBySelection(rows []schema.SnapCount, sel any) []schema.SnapCount {
	if len(rows) == 0 || sel == nil {
		return rows
	}

	switch v := sel.(type) {
	case Seasons:
		seasonSet := make(map[int]struct{}, len(v))
		for _, s := range v {
			seasonSet[s] = struct{}{}
		}
		out := rows[:0]
		for _, r := range rows {
			if _, ok := seasonSet[r.Season]; ok {
				out = append(out, r)
			}
		}
		return out

	case []int:
		seasonSet := make(map[int]struct{}, len(v))
		for _, s := range v {
			seasonSet[s] = struct{}{}
		}
		out := rows[:0]
		for _, r := range rows {
			if _, ok := seasonSet[r.Season]; ok {
				out = append(out, r)
			}
		}
		return out

	case int:
		season := v
		out := rows[:0]
		for _, r := range rows {
			if r.Season == season {
				out = append(out, r)
			}
		}
		return out

	case Weeks:
		// Interpret bare Weeks as "weeks from current season".
		cur := GetCurrentSeason()
		weekSet := make(map[int]struct{}, len(v))
		for _, w := range v {
			weekSet[w] = struct{}{}
		}
		out := rows[:0]
		for _, r := range rows {
			if r.Season == cur {
				if len(weekSet) == 0 {
					out = append(out, r)
					continue
				}
				if _, ok := weekSet[r.Week]; ok {
					out = append(out, r)
				}
			}
		}
		return out

	case []SeasonWeeks:
		type weekset map[int]struct{}
		bySeason := map[int]weekset{}
		for _, sw := range v {
			ws := bySeason[sw.Season]
			if ws == nil {
				ws = weekset{}
				bySeason[sw.Season] = ws
			}
			for _, w := range sw.Weeks {
				ws[w] = struct{}{}
			}
			// empty Weeks means "all weeks" → store nil to signal no week filtering for that season
			if len(sw.Weeks) == 0 {
				bySeason[sw.Season] = nil
			}
		}

		out := rows[:0]
		for _, r := range rows {
			ws, ok := bySeason[r.Season]
			if !ok {
				continue
			}
			if ws == nil { // all weeks for that season
				out = append(out, r)
				continue
			}
			if _, ok := ws[r.Week]; ok {
				out = append(out, r)
			}
		}
		return out

	case bool:
		// true => all; false => no-op; keep as-is
		return rows

	default:
		// Unknown selector → pass-through
		return rows
	}
}
