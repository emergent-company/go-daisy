package ui

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

// renderTempl renders a component to its HTML string for assertions.
func renderTempl(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

func TestIconPickerRendersOptionsAndHiddenInput(t *testing.T) {
	html := renderTempl(t, IconPicker(IconPickerProps{
		ID:           "ip",
		Name:         "icon",
		Value:        "",
		DefaultLabel: "Default icon",
		DefaultClass: "lucide--box",
		Options: []IconPickerOption{
			{Name: "file-text", Label: "File Text", Class: "lucide--file-text"},
			{Name: "user", Label: "User", Class: "lucide--user"},
		},
	}))

	if !strings.Contains(html, "data-gd-icon-picker") {
		t.Error("missing picker root attribute")
	}
	if !strings.Contains(html, `type="hidden"`) || !strings.Contains(html, `name="icon"`) {
		t.Error("hidden input with name=icon not emitted")
	}
	if got := strings.Count(html, "data-gd-icon-option"); got != 3 {
		t.Errorf("option count = %d, want 3 (default + 2)", got)
	}
	for _, want := range []string{
		`data-gd-icon-value="file-text"`,
		`data-gd-icon-value="user"`,
		`lucide--file-text`,
		`lucide--user`,
		`data-gd-icon-default`,
		"Default icon",
		"Search icons",
		"No icons match",
		"Reset to default",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %q", want)
		}
	}
}

func TestIconPickerNoHiddenInputWithoutName(t *testing.T) {
	html := renderTempl(t, IconPicker(IconPickerProps{ID: "ip", Value: ""}))
	if strings.Contains(html, `type="hidden"`) {
		t.Error("hidden input emitted despite empty Name")
	}
}

func TestIconPickerSelectsMatchingOption(t *testing.T) {
	html := renderTempl(t, IconPicker(IconPickerProps{
		ID:           "ip",
		Name:         "icon",
		Value:        "file-text",
		DefaultLabel: "Default icon",
		Options: []IconPickerOption{
			{Name: "file-text", Label: "File Text", Class: "lucide--file-text"},
			{Name: "user", Label: "User", Class: "lucide--user"},
		},
	}))

	// The matching option is the single tab stop and reports selected.
	if !strings.Contains(html, `aria-selected="true"`) {
		t.Error("no selected option")
	}
	if !strings.Contains(html, `value="file-text"`) {
		t.Error("hidden input did not carry the selected value")
	}
	if !strings.Contains(html, "File Text") {
		t.Error("trigger label did not resolve to the option label")
	}
	// Default tile must not claim selection.
	if !strings.Contains(html, `aria-selected="false"`) {
		t.Error("expected at least one unselected option")
	}
}
