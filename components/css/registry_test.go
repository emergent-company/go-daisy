package css

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRegistryLoads verifies the CSS registry contains the core component
// packages and that every listed module set is usable.
func TestRegistryLoads(t *testing.T) {
	for _, pkg := range []string{"ui", "nav", "form", "table", "layout", "alpine", "modal"} {
		if _, ok := Registry[pkg]; !ok {
			t.Errorf("registry missing package %q", pkg)
		}
	}
	if len(Registry["ui"].DaisyUIModules) == 0 {
		t.Error("ui package module list is empty")
	}
	if len(Registry["form"].DaisyUIModules) == 0 {
		t.Error("form package module list is empty")
	}
}

// TestCustomCSSExists verifies the co-located custom CSS file is shipped
// alongside the registry.
func TestCustomCSSExists(t *testing.T) {
	// registry.go lives in components/css/, next to custom.css.
	here, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(here, filepath.Base(CustomCSS))
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("custom.css not found at %s: %v", path, err)
	}
}

// TestMemoryThemeCSS pins the brand theme contract: the file ships next to the
// registry, overrides daisyUI's `dark` theme (so consumers that prune themes
// keep it), and carries both bridges go-daisy markup depends on. Dropping the
// radius alias silently reverts every `rounded-lg` in the library to Tailwind's
// literal scale, and dropping the muted scale breaks `text-muted-faint`, both
// with no compile error.
func TestMemoryThemeCSS(t *testing.T) {
	here, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(here, filepath.Base(MemoryThemeCSS))
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("theme-memory.css not found at %s: %v", path, err)
	}
	src := string(data)

	for _, want := range []struct{ name, needle string }{
		{"daisyUI theme plugin", `@plugin "daisyui/theme"`},
		{"overrides the dark theme", `name: "dark"`},
		{"is the default theme", `default: true`},
		{"declares a dark color scheme", `color-scheme: dark`},
		{"brand primary", `--color-primary:`},
		{"radius bridge", `--radius-lg: var(--radius-box)`},
		{"muted-subtle step", `--color-muted-subtle:`},
		{"muted-faint step", `--color-muted-faint:`},
		{"muted step", `--color-muted:`},
		{"muted-strong step", `--color-muted-strong:`},
		{"muted-bright step", `--color-muted-bright:`},
		{"per-element resolution", "@theme inline {"},
		{"direct var() consumers", "\n:root {"},
	} {
		if !strings.Contains(src, want.needle) {
			t.Errorf("theme-memory.css missing %s (%q)", want.name, want.needle)
		}
	}
}

// TestMemoryThemeMutedScaleIsInlinedAndRooted pins the two-declaration rule: the
// muted scale is declared `@theme inline` (so utilities resolve per element
// under a nested theme) AND in a `:root` block (so a consumer's own
// `var(--color-muted-*)` rules still resolve). Only one of the two silently
// breaks a different consumer, so both must list the same five steps.
//
// The blocks are read as brace-balanced units rather than by slicing at the
// first `@theme inline`, because the radius bridge is another `@theme inline`
// block earlier in the file — a positional slice would silently swallow a
// regression that turned the muted block into a plain `@theme`.
func TestMemoryThemeMutedScaleIsInlinedAndRooted(t *testing.T) {
	here, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	src, err := os.ReadFile(filepath.Join(here, filepath.Base(MemoryThemeCSS)))
	if err != nil {
		t.Fatalf("theme-memory.css: %v", err)
	}
	text := string(src)

	steps := []string{"muted-subtle", "muted-faint", "muted", "muted-strong", "muted-bright"}
	hasAll := func(block string) bool {
		for _, step := range steps {
			if !strings.Contains(block, "--color-"+step+":") {
				return false
			}
		}
		return true
	}

	if !anyBlock(text, "@theme inline {", hasAll) {
		t.Errorf("no `@theme inline` block declares all five muted steps; utilities would freeze at the root theme")
	}
	if !anyBlock(text, ":root {", hasAll) {
		t.Errorf("no `:root` block declares all five muted steps; consumer `var(--color-muted-*)` rules would resolve to nothing")
	}
}

// anyBlock reports whether any brace-balanced block introduced by header
// satisfies want.
func anyBlock(text, header string, want func(string) bool) bool {
	for i := 0; i < len(text); {
		j := strings.Index(text[i:], header)
		if j < 0 {
			return false
		}
		start := i + j
		block, next := blockAfter(text, start)
		if block == "" {
			return false
		}
		if want(block) {
			return true
		}
		i = next
	}
	return false
}

// blockAfter returns the brace-balanced body that begins at the first `{` at or
// after start, plus the index just past it.
func blockAfter(text string, start int) (string, int) {
	open := strings.Index(text[start:], "{")
	if open < 0 {
		return "", len(text)
	}
	open += start
	depth := 0
	for i := open; i < len(text); i++ {
		switch text[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return text[open : i+1], i + 1
			}
		}
	}
	return text[open:], len(text)
}

// TestMemoryThemeCSSIsNotAPackageAsset pins the placement decision: the brand
// theme belongs to no component package, so no package may claim it. If it ever
// appears in Registry it would be pulled into every consumer that asks for
// `ui`, and a consumer running its own brand would silently lose it.
func TestMemoryThemeCSSIsNotAPackageAsset(t *testing.T) {
	for pkg, assets := range Registry {
		for _, f := range assets.CSSFiles {
			if f == MemoryThemeCSS {
				t.Errorf("package %q lists the brand theme in its CSS assets", pkg)
			}
		}
	}
}
