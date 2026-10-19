package handler

import (
	"strconv"
	"strings"
)

func initials(name string) string {
	if len(name) == 0 {
		return "?"
	}
	r := []rune(name)
	if len(r) >= 2 {
		return string(r[:2])
	}
	return string(r[0])
}

func capFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func productNameOnly(full string) string {
	if idx := strings.LastIndex(full, " #"); idx > 0 {
		return full[:idx]
	}
	return full
}

func productSKU(full string) string {
	if idx := strings.LastIndex(full, " #"); idx > 0 {
		return full[idx+1:]
	}
	return ""
}

// shopRatingClass maps a numeric rating to its star tone: success for >= 4.0,
// warning for >= 3.0 and error below.
func shopRatingClass(rating string) string {
	f, err := strconv.ParseFloat(rating, 64)
	if err != nil {
		return ""
	}
	switch {
	case f >= 4.0:
		return "text-success"
	case f >= 3.0:
		return "text-warning"
	default:
		return "text-error"
	}
}
