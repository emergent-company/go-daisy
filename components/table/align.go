package table

import "strings"

// TableAlign controls the text alignment of a table header or body cell.
//
// The zero value emits no class, so a cell inherits the alignment of the table.
// The constant values are the logical vocabulary ("start" / "end"); they resolve
// to physical `text-left` / `text-right` utilities because the library has no
// RTL story yet, so a future switch to logical properties stays local to
// alignClasses.
type TableAlign string

const (
	AlignDefault TableAlign = ""      // zero value: no alignment class
	AlignStart   TableAlign = "start" // text-left
	AlignEnd     TableAlign = "end"   // text-right
)

// alignClasses returns the Tailwind alignment utilities for a cell. Numeric
// implies end alignment and adds tabular-nums so digits keep their column.
// Headers and body cells share this helper, so a right-aligned column cannot
// drift between the <th> and its <td>s.
func alignClasses(align TableAlign, numeric bool) string {
	if numeric {
		return "text-right tabular-nums"
	}
	switch align {
	case AlignStart:
		return "text-left"
	case AlignEnd:
		return "text-right"
	}
	return ""
}

// cellClasses joins the non-empty class fragments with single spaces.
// It never emits leading, trailing or doubled spaces, so an empty set of
// fragments renders as an empty class attribute, matching the pre-props
// components' output.
func cellClasses(parts ...string) string {
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return strings.Join(out, " ")
}

// endAligned reports whether a sortable header's column is right-aligned,
// which packs the link and direction icon against the cell's right edge.
func (p SortableTableHeaderProps) endAligned() bool {
	return p.Numeric || p.Align == AlignEnd
}
