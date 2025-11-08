package parse

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

// Auto parses bytes into []map[string]any by sniffing URL/bytes.
// CSV is supported; Parquet returns a clear "not implemented" for now.
func Auto(b []byte, usedURL string) ([]map[string]any, error) {
	ext := strings.ToLower(filepath.Ext(usedURL))
	if ext == ".csv" || looksLikeCSV(b) {
		return parseCSVMaps(bytes.NewReader(b))
	}
	if ext == ".parquet" {
		return nil, errors.New("parquet parsing not implemented yet")
	}
	// Fallback by content-type sniffing
	mt := http.DetectContentType(peek512(b))
	if strings.Contains(mt, "text/plain") || strings.Contains(mt, "text/csv") {
		return parseCSVMaps(bytes.NewReader(b))
	}
	return nil, errors.New("unknown content type; cannot parse")
}

func parseCSVMaps(r io.Reader) ([]map[string]any, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	hdr, err := cr.Read()
	if err != nil {
		return nil, err
	}
	norm := make([]string, len(hdr))
	dupCount := map[string]int{}
	for i, h := range hdr {
		base := normalize(h)
		dupCount[base]++
		if dupCount[base] > 1 {
			base = base + "_" + itoa(dupCount[base]) // disambiguate dup headers
		}
		norm[i] = base
	}

	var out []map[string]any
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		m := make(map[string]any, len(norm))
		for i := range norm {
			if i < len(rec) {
				m[norm[i]] = strings.TrimSpace(rec[i])
			}
		}
		out = append(out, m)
	}
	return out, nil
}

// func normalize(s string) string {
// 	s = strings.TrimSpace(strings.ToLower(s))
// 	s = strings.ReplaceAll(s, " ", "_")
// 	s = strings.ReplaceAll(s, "-", "_")
// 	return s
// }

func itoa(i int) string { return strconvItoa(i) }

// small local itoa to avoid extra imports
func strconvItoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	var b [20]byte
	pos := len(b)
	for i > 0 {
		pos--
		b[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		b[pos] = '-'
	}
	return string(b[pos:])
}

func looksLikeCSV(b []byte) bool {
	s := string(peek512(b))
	// very light heuristic
	return strings.Contains(s, ",") && strings.Contains(s, "\n")
}

func peek512(b []byte) []byte {
	if len(b) > 512 {
		return b[:512]
	}
	return b
}
