package ui

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// renderCommandPalette renders CommandPalette to its HTML string.
func renderCommandPalette(t *testing.T, props CommandPaletteProps) string {
	t.Helper()
	var buf bytes.Buffer
	if err := CommandPalette(props).Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestCommandPaletteDialogIsModal pins the a11y contract: the dialog must
// announce itself as modal and carry an accessible name, so screen readers do
// not treat the page behind the palette as reachable.
func TestCommandPaletteDialogIsModal(t *testing.T) {
	got := renderCommandPalette(t, CommandPaletteProps{ID: "spotlight"})
	for _, want := range []string{
		`role="dialog"`,
		`aria-modal="true"`,
		`aria-label="Command palette"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("CommandPalette must render %s:\n%s", want, got)
		}
	}
}

// TestCommandPaletteAriaLabelOverride lets apps localise the dialog name.
func TestCommandPaletteAriaLabelOverride(t *testing.T) {
	got := renderCommandPalette(t, CommandPaletteProps{ID: "spotlight", AriaLabel: "Springen"})
	if !strings.Contains(got, `aria-label="Springen"`) {
		t.Errorf("AriaLabel must override the default dialog name:\n%s", got)
	}
	if strings.Contains(got, `aria-label="Command palette"`) {
		t.Errorf("default aria-label must not be rendered alongside the override:\n%s", got)
	}
}

// TestCommandPaletteLoadsKeyboardRuntime pins the staticfs wiring: the focus
// trap lives in go-daisy-a11y.js and the palette logic in
// go-daisy-command-palette.js, each included exactly once.
func TestCommandPaletteLoadsKeyboardRuntime(t *testing.T) {
	got := renderCommandPalette(t, CommandPaletteProps{ID: "spotlight"})
	for _, file := range []string{"go-daisy-a11y.js", "go-daisy-command-palette.js"} {
		if n := strings.Count(got, file); n != 1 {
			t.Errorf("CommandPalette must include %s exactly once, got %d:\n%s", file, n, got)
		}
	}
	if !strings.Contains(got, `data-command-palette="true"`) {
		t.Errorf("toggle must stay discoverable by the runtime:\n%s", got)
	}
}

// TestCommandPaletteDoesNotAutofocus pins that opening focus comes from
// gdFocus.trap. autofocus on a server-rendered Open palette would make the
// hidden input the element that release() "restores" to after close.
func TestCommandPaletteDoesNotAutofocus(t *testing.T) {
	got := renderCommandPalette(t, CommandPaletteProps{ID: "spotlight", Open: true})
	if strings.Contains(got, "autofocus") {
		t.Errorf("CommandPalette must not render autofocus; the trap owns initial focus:\n%s", got)
	}
}
