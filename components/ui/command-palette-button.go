package ui

import "github.com/a-h/templ"

// commandPaletteButtonAttrs builds the onclick wiring for a CommandPaletteButton:
// check the palette's toggle input and focus its text input. Caller-supplied
// attrs override the generated onclick (take full control).
func commandPaletteButtonAttrs(id string, extra templ.Attributes) templ.Attributes {
	attrs := templ.Attributes{
		"onclick": "var cb = document.getElementById('" + id + "-toggle'); if (cb) cb.checked = true; var i = document.getElementById('" + id + "-input'); if (i) i.focus();",
	}
	for k, v := range extra {
		attrs[k] = v
	}
	return attrs
}
