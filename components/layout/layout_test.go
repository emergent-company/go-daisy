package layout

import (
	"context"
	"strings"
	"testing"
)

// TestPageFullPinTheme guards the PinTheme contract: when set, PageFull must
// omit ui.ThemeScript() (the persisted gallery palette/mode selection) so the
// ThemeAttr value is authoritative; when unset, the script markers are present.
func TestPageFullPinTheme(t *testing.T) {
	render := func(props PageProps) string {
		t.Helper()
		var buf strings.Builder
		if err := PageFull(props).Render(context.Background(), &buf); err != nil {
			t.Fatalf("render failed: %v", err)
		}
		return buf.String()
	}

	pinned := render(PageProps{Title: "t", ThemeAttr: "nexus", PinTheme: true})
	if !strings.Contains(pinned, `data-theme="nexus"`) {
		t.Errorf("pinned page missing data-theme=\"nexus\": %s", pinned)
	}
	for _, marker := range []string{"gallery-ui-palette", "__goDaisyTheme"} {
		if strings.Contains(pinned, marker) {
			t.Errorf("pinned page still emits ThemeScript marker %q", marker)
		}
	}

	unpinned := render(PageProps{Title: "t", ThemeAttr: "nexus", PinTheme: false})
	for _, marker := range []string{"gallery-ui-palette", "__goDaisyTheme"} {
		if !strings.Contains(unpinned, marker) {
			t.Errorf("unpinned page missing ThemeScript marker %q", marker)
		}
	}
}
