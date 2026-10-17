package pages

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

// renderPages renders a page component to its HTML string.
func renderPages(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestNotificationToggleOutputLocked locks the settings row toggle markup.
// The row owns its title and description, so the control must stay a bare
// ToggleInput carrying the primary colour: no form-control wrapper, no
// component-rendered label, and the accessible name comes from the row title.
func TestNotificationToggleOutputLocked(t *testing.T) {
	got := renderPages(t, notificationToggle("mentions", "Mentions", "Get notified when someone mentions you", true))
	want := `<div class="flex items-start justify-between gap-4 py-2"><div><p class="text-sm font-medium text-base-content">Mentions</p><p class="text-xs text-base-content/50">Get notified when someone mentions you</p></div><input type="checkbox" name="mentions" class="toggle toggle-primary" checked aria-label="Mentions"></div>`
	if got != want {
		t.Errorf("output changed:\n got: %q\nwant: %q", got, want)
	}
}

// TestSettingsTogglesKeepPrimaryEmphasis asserts every toggle on the settings
// page keeps primary emphasis and stays a bare input with no wrapper layer. The
// two appearance rows plus the six notification rows render eight switches, and
// each one must be an immediate sibling of its description block (no
// form-control wrapper, no component-owned label between them).
func TestSettingsTogglesKeepPrimaryEmphasis(t *testing.T) {
	out := renderPages(t, SettingsPage())

	if n := strings.Count(out, `class="toggle toggle-primary"`); n != 8 {
		t.Errorf("primary toggles = %d, want 8", n)
	}
	if n := strings.Count(out, `type="checkbox" name=`); n != 8 {
		t.Errorf("toggle inputs = %d, want 8", n)
	}
	if strings.Contains(out, `class="toggle"`) {
		t.Error("settings page renders a neutral toggle; primary emphasis was dropped")
	}
	if n := strings.Count(out, `</div><input type="checkbox" name=`); n != 8 {
		t.Errorf("bare row toggles = %d, want 8 (a wrapper layer was reintroduced)", n)
	}
}

// TestSettingsTogglesHaveAccessibleNames guards the a11y reason the migration
// passes Attrs: the visible row title is a <p>, not a <label>, so each switch
// must carry its own aria-label.
func TestSettingsTogglesHaveAccessibleNames(t *testing.T) {
	out := renderPages(t, SettingsPage())

	for _, name := range []string{"Compact view", "Reduced motion", "Mentions", "Marketing emails"} {
		if !strings.Contains(out, `aria-label="`+name+`"`) {
			t.Errorf("no aria-label %q on the rendered page", name)
		}
	}
}
