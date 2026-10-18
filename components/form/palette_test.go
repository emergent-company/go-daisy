package form

import (
	"strings"
	"testing"
)

// TestPaletteXDataStartsEmpty is a regression test for a live bug on master:
// Palette rendered x-data="{open:false}" (alpine.Toggle) while its markup binds
// x-model="picked" and every swatch handler assigns `picked`. Alpine evaluates
// x-model during init, so the missing property produced
// "ReferenceError: picked is not defined" on every palette preview — the
// picker's hidden input never got a value and clicks silently failed.
//
// It also locks the initial selection: no selection means picked is "" (the
// reset handler assigns the same value), not an absent key. The x-data object
// is asserted verbatim rather than by substring because `picked` appears in the
// click handlers and x-model regardless — only the x-data attribute proves the
// property is declared.
func TestPaletteXDataStartsEmpty(t *testing.T) {
	html := renderForm(t, Palette(PaletteProps{ID: "pal", Hues: []string{"#ff0000"}}))

	if !strings.Contains(html, `x-data="{&#34;picked&#34;:&#34;&#34;}"`) {
		t.Errorf("Palette x-data is not the empty-selection state:\n%s", html)
	}
	if !strings.Contains(html, `x-model="picked"`) {
		t.Errorf("Palette hidden input lost its x-model binding:\n%s", html)
	}
}

// TestPaletteNameIsLiteralAttribute is the second half of the same bug: the
// hidden input carried x-bind:name={props.ID}, and Alpine evaluates an x-bind
// value as a JavaScript expression. With ID "demo-palette" Alpine compiled
// `demo-palette` (subtraction) and threw "ReferenceError: demo is not defined"
// during initialization, so the field never submitted under its real name.
// The name is static markup — it must not be bound.
func TestPaletteNameIsLiteralAttribute(t *testing.T) {
	html := renderForm(t, Palette(PaletteProps{ID: "demo-palette", Hues: []string{"#ff0000"}}))

	if strings.Contains(html, "x-bind:name") {
		t.Errorf("Palette binds its hidden input name as an Alpine expression:\n%s", html)
	}
	if !strings.Contains(html, `name="demo-palette"`) {
		t.Errorf("Palette hidden input is missing the literal name attribute:\n%s", html)
	}
}
