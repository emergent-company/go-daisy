package form

import (
	"strings"
	"testing"
)

// TestSelectShellTriggerTextBinding is a regression test for a live bug on
// master: with no ValueExpr the component emitted x-text="" (an empty Alpine
// expression) and the gallery seeded it with "picked", a property no ancestor
// x-data declares. Either way Alpine threw at init and the trigger text was
// undefined behaviour. The placeholder must survive as static text when no
// binding is requested.
func TestSelectShellTriggerTextBinding(t *testing.T) {
	cases := []struct {
		name      string
		valueExpr string
		want      []string
	}{
		{
			name:      "no_binding_keeps_placeholder",
			valueExpr: "",
			want:      []string{`<span class="block truncate flex-1">Select...</span>`},
		},
		{
			name:      "binding_renders_expression",
			valueExpr: "picked",
			want:      []string{`<span class="block truncate flex-1" x-text="picked">Select...</span>`},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			html := renderForm(t, SelectShell(SelectShellProps{
				ID: "shell", Placeholder: "Select...", ValueExpr: tc.valueExpr,
			}))
			for _, want := range tc.want {
				if !strings.Contains(html, want) {
					t.Errorf("missing %q in:\n%s", want, html)
				}
			}
			if strings.Contains(html, `x-text=""`) {
				t.Errorf("SelectShell emitted an empty x-text expression:\n%s", html)
			}
		})
	}
}

// TestSelectShellHonoursSelectClose locks the contract the component's own doc
// comment states ("Dispatch select-close on pick"): children in the panel are
// expected to dispatch that event to close the shell. Master documented it but
// never listened for it, so an option click left the panel open — Escape and
// click-outside were the only ways to close it.
//
// The listener must be element-scoped. `$dispatch` bubbles, so a `.window`
// modifier would close every SelectShell on the page instead of the one whose
// option was clicked; options are DOM descendants of this root, so a plain
// listener already sees the event.
func TestSelectShellHonoursSelectClose(t *testing.T) {
	html := renderForm(t, SelectShell(SelectShellProps{ID: "shell", Placeholder: "Select..."}))

	if !strings.Contains(html, `@select-close="open = false"`) {
		t.Errorf("SelectShell does not listen for the documented select-close event:\n%s", html)
	}
	if strings.Contains(html, "@select-close.window") {
		t.Errorf("select-close listener is window-scoped and will close other shells:\n%s", html)
	}
}

// TestSelectShellPanelIsWired checks the trigger/listbox ARIA wiring: the panel
// needs an id for the trigger to point at.
func TestSelectShellPanelIsWired(t *testing.T) {
	html := renderForm(t, SelectShell(SelectShellProps{ID: "shell", Placeholder: "Select..."}))

	for _, want := range []string{`aria-controls="shell-panel"`, `id="shell-panel"`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %q in:\n%s", want, html)
		}
	}
}
