package staticfs

import (
	"io/fs"
	"strings"
	"testing"
)

// TestCompiledCSSContainsEmptyNoteToneUtilities guards the EmptyNote Tone axis
// end-to-end. TestStaticFilesHashMatchesCommitted only catches a *stale* hash;
// it cannot catch a deliberate removal of a tone utility whose hash was then
// regenerated — the tone would silently stop applying in production. Each
// EmptyNote tone is a fixed base-content opacity in the safelist ladder
// (assets/app.css: `text-base-content/{5,10,…,65,…}`), compiled to a base rule
// (`.text-base-content\/NN{color:var(--color-base-content)}`) plus an
// `@supports color-mix` override that actually applies the NN% opacity. Pin
// both for every EmptyNote tone so a removal fails loudly rather than shipping
// a note at full opacity.
func TestCompiledCSSContainsEmptyNoteToneUtilities(t *testing.T) {
	css, err := fs.ReadFile(files, "static/css/app.css")
	if err != nil {
		t.Fatalf("read compiled css: %v", err)
	}
	src := string(css)

	// The three steps EmptyNoteTone maps to (components/ui/empty-note.templ):
	// Default → 50, Quiet → 40, Strong → 65.
	for _, step := range []string{"40", "50", "65"} {
		selector := `.text-base-content\/` + step
		if !strings.Contains(src, selector) {
			t.Errorf("compiled css is missing the EmptyNote tone selector %q — run `task build:ui` or the tone breaks in production", selector)
		}
		// The color-mix rule is what applies the per-step opacity in modern
		// browsers; without it the selector renders at full base-content color.
		mix := `color-mix(in oklab, var(--color-base-content) ` + step + `%, transparent)`
		if !strings.Contains(src, mix) {
			t.Errorf("compiled css is missing the EmptyNote tone opacity rule %q — run `task build:ui` or the tone applies at full opacity", mix)
		}
	}
}
