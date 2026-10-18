package ui

import "strings"

// listRowClass returns the class list for a list row (or the anchor that
// carries the row when the row is a link). Variant and size modifiers live
// here so the plain and link branches cannot drift.
func listRowClass(p ListRowProps) string {
	classes := []string{"list-row"}
	switch p.Variant {
	case ListRowOutline:
		classes = append(classes, "rounded-box", "border", "border-base-300", "bg-base-100")
	case ListRowMuted:
		classes = append(classes, "bg-base-200")
	}
	switch p.Size {
	case ListRowSizeSM:
		classes = append(classes, "py-1")
	case ListRowSizeXS:
		classes = append(classes, "py-0.5")
	}
	return strings.Join(classes, " ")
}
