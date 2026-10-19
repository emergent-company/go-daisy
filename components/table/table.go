package table

import (
	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/components/ui"
	"github.com/emergent-company/go-daisy/shared"
)

func ternary(cond bool, a, b string) string { return shared.Ternary(cond, a, b) }

// ternaryInt returns a when cond is true, otherwise b.
func ternaryInt(cond bool, a, b int) int {
	if cond {
		return a
	}
	return b
}

// ternaryVariant returns the error button variant when danger is set, else the zero value.
func ternaryVariant(danger bool) ui.ButtonVariant {
	if danger {
		return ui.ButtonError
	}
	return ""
}

// dataTablePageSize returns the page size with default fallback.
func dataTablePageSize(n int) int {
	if n == 0 {
		return 10
	}
	return n
}

// rowAttrs merges the hover class with any caller-supplied "class" in attrs so
// TableRow never emits a duplicate class attribute (which the browser would
// otherwise resolve to the first one, silently dropping the caller's classes).
func rowAttrs(hover bool, attrs templ.Attributes) templ.Attributes {
	out := templ.Attributes{}
	for k, v := range attrs {
		out[k] = v
	}
	cls := ""
	if hover {
		cls = "hover"
	}
	if extra, ok := out["class"].(string); ok && extra != "" {
		if cls != "" {
			cls += " "
		}
		cls += extra
	}
	if cls != "" {
		out["class"] = cls
	} else {
		delete(out, "class")
	}
	return out
}
