package galleryruntime

import "testing"

// TestCustomThemeCSSDeterministic locks the emitted rule: map iteration order
// must not leak into the stylesheet (stable CSS diffs, stable snapshots).
func TestCustomThemeCSSDeterministic(t *testing.T) {
	theme := CustomTheme{
		Name:   "custom",
		Scheme: "dark",
		Vars: map[string]string{
			"--color-primary":  "oklch(60% 0.2 150)",
			"--color-base-100": "oklch(15% 0 0)",
			"--radius-box":     "0.5rem",
		},
	}
	want := `[data-theme="custom"]{color-scheme:dark;--color-base-100:oklch(15% 0 0);--color-primary:oklch(60% 0.2 150);--radius-box:0.5rem;}`
	for i := 0; i < 50; i++ {
		if got := theme.CSS(""); got != want {
			t.Fatalf("CSS() not deterministic:\n got: %s\nwant: %s", got, want)
		}
	}
	if got, wantSel := theme.CSS(".scope"), `.scope{color-scheme:dark;--color-base-100:oklch(15% 0 0);--color-primary:oklch(60% 0.2 150);--radius-box:0.5rem;}`; got != wantSel {
		t.Errorf("explicit selector:\n got: %s\nwant: %s", got, wantSel)
	}
	// Typography tokens survive verbatim (commas and quotes included): the
	// studio writes font stacks through the same custom-theme seam.
	if got := (CustomTheme{Vars: map[string]string{"--font-body": `"Segoe UI", system-ui, sans-serif`}}).CSS(""); got != `[data-theme="custom"]{color-scheme:light;--font-body:"Segoe UI", system-ui, sans-serif;}` {
		t.Errorf("font stack mangled: %s", got)
	}

	// Empty scheme defaults to light rather than emitting an empty value.
	if got := (CustomTheme{Vars: map[string]string{"--a": "b"}}).CSS(""); got != `[data-theme="custom"]{color-scheme:light;--a:b;}` {
		t.Errorf("scheme default: %s", got)
	}
}
