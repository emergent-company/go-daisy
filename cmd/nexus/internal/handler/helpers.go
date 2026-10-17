package handler

import (
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
