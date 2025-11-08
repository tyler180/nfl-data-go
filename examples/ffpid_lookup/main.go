// examples/ffpid_lookup/main.go
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	ffpid "github.com/tyler180/nfl-data-go/internal/datasets/ffplayerids"
)

// Heuristics for which struct fields we consider "IDs".
func isIDField(field reflect.StructField) bool {
	name := field.Name
	tag := field.Tag.Get("json")
	nameLower := strings.ToLower(name)
	tagLower := strings.ToLower(tag)

	// Common patterns
	if strings.HasSuffix(name, "ID") || strings.Contains(tagLower, "_id") {
		return true
	}

	// Some columns don't strictly end with "ID"; whitelist a few
	whitelist := []string{
		"pfr", "pfr_id", "pfr_player_id", "gsis", "gsis_id", "sportradar_id",
		"espn_id", "yahoo_id", "sleeper_id", "mfl_id", "cbs_id",
		"pff_id", "fleaflicker_id", "fantasypros_id", "ktc_id", "rotowire_id",
	}
	for _, w := range whitelist {
		if strings.EqualFold(nameLower, w) || strings.EqualFold(tagLower, w) {
			return true
		}
	}
	return false
}

// Useful descriptive (non-ID) fields to show in output if present.
func isMetaField(field reflect.StructField) bool {
	name := strings.ToLower(field.Name)
	tag := strings.ToLower(field.Tag.Get("json"))
	meta := []string{
		"name", "player", "player_name", "full_name",
		"position", "pos",
		"team", "team_abbr",
		"first_name", "last_name",
	}
	for _, m := range meta {
		if name == m || tag == m {
			return true
		}
	}
	return false
}

// Extract all ID fields (and some meta) from a struct value into maps.
func collectFields(v any) (ids map[string]string, meta map[string]string) {
	ids = map[string]string{}
	meta = map[string]string{}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	rt := rv.Type()
	if rt.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if !f.IsExported() {
			continue
		}
		val := rv.Field(i)
		// Render value -> string
		var s string
		switch val.Kind() {
		case reflect.String:
			s = strings.TrimSpace(val.String())
		case reflect.Int, reflect.Int32, reflect.Int64:
			s = fmt.Sprintf("%d", val.Int())
		case reflect.Uint, reflect.Uint32, reflect.Uint64:
			s = fmt.Sprintf("%d", val.Uint())
		default:
			continue // ignore non-scalar fields
		}
		if s == "" {
			continue
		}
		// Preferred key = json tag if present, else field name (lowercase)
		key := f.Tag.Get("json")
		if key == "" || key == "-" {
			key = strings.ToLower(f.Name)
		}

		if isIDField(f) {
			ids[key] = s
		} else if isMetaField(f) {
			meta[key] = s
		}
	}
	return
}

func main() {
	query := flag.String("id", "", "any player ID (Sleeper, ESPN, Yahoo, MFL, GSIS, PFR, PFF, Sportradar, CBS, Fleaflicker, FantasyPros, KTC, Rotowire, etc.)")
	flag.Parse()

	if strings.TrimSpace(*query) == "" {
		fmt.Fprintln(os.Stderr, "usage: ffpid_lookup -id <some_id_value>")
		os.Exit(2)
	}

	ctx := context.Background()
	rows, err := ffpid.Load(ctx)
	if err != nil {
		log.Fatalf("ffplayerids.Load: %v", err)
	}

	q := strings.TrimSpace(*query)
	qLower := strings.ToLower(q)

	type Out struct {
		Meta map[string]string            `json:"meta,omitempty"`
		IDs  map[string]string            `json:"ids"`
		Raw  map[string]map[string]string `json:"_debug,omitempty"` // optional, remove if not wanted
	}
	var results []Out

	for _, r := range rows {
		ids, meta := collectFields(r)
		// match against ANY id value (case-insensitive)
		found := false
		for _, v := range ids {
			if strings.EqualFold(v, qLower) || strings.EqualFold(v, q) {
				found = true
				break
			}
			// also tolerate stripped leading zeros
			if strings.TrimLeft(strings.ToLower(v), "0") == strings.TrimLeft(qLower, "0") && v != "" && q != "" {
				found = true
				break
			}
		}
		if !found {
			continue
		}

		results = append(results, Out{
			Meta: meta,
			IDs:  ids,
		})
	}

	if len(results) == 0 {
		// Give the user a hint of what fields exist by sampling the first row
		if len(rows) > 0 {
			ids, _ := collectFields(rows[0])
			have := make([]string, 0, len(ids))
			for k := range ids {
				have = append(have, k)
			}
			log.Printf("no matches for %q; known ID fields include: %s", q, strings.Join(have, ", "))
		} else {
			log.Printf("no playerid rows loaded")
		}
		return
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(results); err != nil {
		log.Fatal(err)
	}
}
