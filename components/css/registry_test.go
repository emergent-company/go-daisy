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
func TestMemoryThemeMutedScaleIsInlinedAndRooted(t *testing.T) {
	here, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	src, err := os.ReadFile(filepath.Join(here, filepath.Base(MemoryThemeCSS)))
	if err != nil {
		t.Fatalf("theme-memory.css: %v", err)
	}

	steps := []string{"muted-subtle", "muted-faint", "muted", "muted-strong", "muted-bright"}
	text := string(src)
	inlineAt := strings.Index(text, "@theme inline {")
	if inlineAt < 0 {
		t.Fatal("no @theme inline block")
	}
	rootAt := strings.Index(text, "\n:root {")
	if rootAt < 0 {
		t.Fatal("no :root block")
	}
	inlineBlock := text[inlineAt:rootAt]
	rootBlock := text[rootAt:]

	for _, step := range steps {
		decl := "--color-" + step + ":"
		if !strings.Contains(inlineBlock, decl) {
			t.Errorf("@theme inline block missing %s", decl)
		}
		if !strings.Contains(rootBlock, decl) {
			t.Errorf(":root block missing %s", decl)
		}
	}
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
