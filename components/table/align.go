package table

import (
	"strings"

	"github.com/a-h/templ"
)

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

// mergeAttrs merges a cell's default attributes with the caller's so every
// attribute is emitted exactly once (a duplicate `class` or `scope` would be
// resolved by the browser to the first one, silently dropping the other).
//
// `class` is always combined: the defaults' classes come first, any
// caller-supplied class is appended — the same rule TableRow's rowAttrs uses.
// Every other key, including `scope`, lets the caller override the default.
func mergeAttrs(defaults, caller templ.Attributes) templ.Attributes {
	out := templ.Attributes{}
	for k, v := range defaults {
		out[k] = v
	}
	for k, v := range caller {
		if k == "class" {
			extra, ok := v.(string)
			if !ok || strings.TrimSpace(extra) == "" {
				continue
			}
			base, _ := out["class"].(string)
			out["class"] = cellClasses(base, extra)
			continue
		}
		out[k] = v
	}
	return out
}

// headCellAttrs builds the <th> attributes for TableHeadCellWithProps.
func headCellAttrs(props TableHeadCellProps) templ.Attributes {
	return mergeAttrs(templ.Attributes{
		"scope": "col",
		"class": cellClasses("text-sm font-medium text-muted", alignClasses(props.Align, props.Numeric), props.Class),
	}, props.Attrs)
}

// bodyCellAttrs builds the <td> attributes for TableCellWithProps.
func bodyCellAttrs(props TableCellProps) templ.Attributes {
	return mergeAttrs(templ.Attributes{
		"class": cellClasses(alignClasses(props.Align, props.Numeric), props.Class),
	}, props.Attrs)
}

// sortableHeaderAttrs builds the <th> attributes for SortableTableHeader. The
// aria-sort default is only a default: callers that manage sorting themselves
// may pass their own.
func sortableHeaderAttrs(props SortableTableHeaderProps) templ.Attributes {
	attrs := mergeAttrs(templ.Attributes{
		"scope": "col",
		"class": cellClasses(alignClasses(props.Align, props.Numeric), props.Class),
	}, props.Attrs)
	if aria := props.Sort.AriaSort(props.Column); aria != "" {
		if _, overridden := props.Attrs["aria-sort"]; !overridden {
			attrs["aria-sort"] = aria
		}
	}
	return attrs
}
