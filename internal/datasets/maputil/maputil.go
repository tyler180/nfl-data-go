package maputil

import (
	"strconv"
	"strings"
)

func Get(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != nil {
			s := strings.TrimSpace(toString(v))
			if s != "" {
				return s
			}
		}
	}
	return ""
}

func Upper(m map[string]any, keys ...string) string {
	return strings.ToUpper(Get(m, keys...))
}

func Int(m map[string]any, keys ...string) int {
	s := Get(m, keys...)
	if s == "" {
		return 0
	}
	// fast-ish parse without panics
	sign, n := 1, 0
	for i, r := range s {
		if i == 0 && r == '-' {
			sign = -1
			continue
		}
		if r < '0' || r > '9' {
			break
		}
		n = n*10 + int(r-'0')
	}
	return sign * n
}

func Float(m map[string]any, keys ...string) float64 {
	s := Get(m, keys...)
	if s == "" {
		return 0
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	// minimal fallback
	var sign float64 = 1
	i := 0
	if len(s) > 0 && s[0] == '-' {
		sign = -1
		i++
	}
	var intPart int64
	var fracPart int64
	var fracDiv float64 = 1
	seenDot := false
	for ; i < len(s); i++ {
		c := s[i]
		if c == '.' && !seenDot {
			seenDot = true
			continue
		}
		if c < '0' || c > '9' {
			break
		}
		if !seenDot {
			intPart = intPart*10 + int64(c-'0')
		} else {
			fracPart = fracPart*10 + int64(c-'0')
			fracDiv *= 10
		}
	}
	return sign * (float64(intPart) + float64(fracPart)/fracDiv)
}

func toString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	default:
		return strings.TrimSpace(strings.ReplaceAll(strings.TrimSpace(fmtAny(v)), "\u0000", ""))
	}
}

func fmtAny(v any) string {
	// tiny, alloc-light fmt.Sprint alternative for common scalar types
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
		return strconv.FormatFloat(t, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(t), 'f', -1, 32)
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}
