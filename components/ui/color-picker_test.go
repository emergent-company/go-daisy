package ui

import (
	"strings"
	"testing"
)

func TestColorPickerRendersTextSwatchAndPresets(t *testing.T) {
	html := renderTempl(t, ColorPicker(ColorPickerProps{
		ID:          "cp",
		Name:        "color",
		Value:       "#4F46E5",
		Placeholder: "#4F46E5",
		AllowEmpty:  true,
		Presets:     []string{"#4F46E5", "#EF4444"},
	}))

	for _, want := range []string{
		`data-gd-color-picker`,
		`name="color"`,
		`value="#4F46E5"`,
		`data-gd-color-swatch`,
		`value="#4f46e5"`, // swatch normalised to lowercase hex
		`data-gd-color-value="#4F46E5"`,
		`aria-pressed="true"`,
		`aria-label="Clear color"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %q", want)
		}
	}
}

func TestColorPickerNoClearWithoutAllowEmpty(t *testing.T) {
	html := renderTempl(t, ColorPicker(ColorPickerProps{ID: "cp", Name: "c", Value: "#000000"}))
	if strings.Contains(html, `aria-label="Clear color"`) {
		t.Error("clear button emitted despite AllowEmpty=false")
	}
	if strings.Contains(html, `aria-label="Use color`) {
		t.Error("preset buttons emitted despite empty Presets")
	}
}

func TestColorPickerSwatchHex(t *testing.T) {
	cases := map[string]string{
		"#abc":         "#aabbcc",
		"#ABCDEF":      "#abcdef",
		"#11223344":    "#112233",
		"oklch(1 0 0)": "#000000",
		"blue":         "#000000",
		"":             "#000000",
		"  #0d9488  ":  "#0d9488",
	}
	for in, want := range cases {
		if got := colorPickerSwatchHex(in); got != want {
			t.Errorf("colorPickerSwatchHex(%q) = %q, want %q", in, got, want)
		}
	}
}
