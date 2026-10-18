package ui

import "github.com/a-h/templ"

// commandPaletteButtonAttrs builds the onclick wiring for a CommandPaletteButton:
// open the palette through its runtime so the focus trap and focus restore run
// on the button path too. The fallback keeps working without the runtime.
// Caller-supplied attrs override the generated onclick (take full control).
func commandPaletteButtonAttrs(id string, extra templ.Attributes) templ.Attributes {
	attrs := templ.Attributes{
		"onclick": "if (window.gdCommandPalette) { window.gdCommandPalette.open('" + id + "'); } else { var cb = document.getElementById('" + id + "-toggle'); if (cb) cb.checked = true; var i = document.getElementById('" + id + "-input'); if (i) i.focus(); }",
	}
	for k, v := range extra {
		attrs[k] = v
	}
	return attrs
}
