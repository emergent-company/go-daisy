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
