package staticfs

import (
	"io/fs"
	"strings"
	"testing"
)

// TestCompiledCSSContainsEmptyNoteToneUtilities guards the EmptyNote Tone axis
// end-to-end. TestStaticFilesHashMatchesCommitted only catches a *stale* hash;
// it cannot catch a deliberate removal of a tone utility whose hash was then
// regenerated — the tone would silently stop applying in production.
//
// Each EmptyNote tone is a fixed base-content opacity in the safelist ladder
// (assets/app.css: `text-base-content/{5,10,…,65,…}`), compiled to two rules:
//
//	.text-base-content\/NN{color:var(--color-base-content)}                              // fallback
//	.text-base-content\/NN{color:color-mix(in oklab, var(--color-base-content) NN%, transparent)} // opacity
//
// The second rule is what applies the NN% opacity in modern browsers; without
// it the tone renders at full base-content colour.
//
// Each assertion pins the exact full rule (selector + property + value), not a
// loose substring and not a count. Both alternatives are brittle here:
//   - `color-mix(in oklab, var(--color-base-content) NN%, transparent)` as a
//     bare substring is also emitted by unrelated DaisyUI component rules
//     (disabled .input/.select/.textarea, .menu-title, the aria-selected
//     step), so for the /50 step a surviving group-hover variant — or even a
//     DaisyUI component rule — would keep a substring assertion green.
//   - A count of color-mix occurrences per step would break on a DaisyUI
//     version bump that adds/removes an unrelated component rule, even though
//     the EmptyNote tone is unchanged.
func TestCompiledCSSContainsEmptyNoteToneUtilities(t *testing.T) {
	css, err := fs.ReadFile(files, "static/css/app.css")
	if err != nil {
		t.Fatalf("read compiled css: %v", err)
	}
	src := string(css)

	// The three steps EmptyNoteTone maps to (components/ui/empty-note.templ):
	// Default → 50, Quiet → 40, Strong → 65.
	for _, step := range []string{"40", "50", "65"} {
		// Fallback rule: full base-content colour for browsers without
		// color-mix. Removing it leaves the utility undefined.
		base := `.text-base-content\/` + step + `{color:var(--color-base-content)}`
		if !strings.Contains(src, base) {
			t.Errorf("compiled css is missing the EmptyNote tone fallback rule %q — run `task build:ui`", base)
		}
		// Opacity rule: the per-step color-mix override that actually applies
		// the NN% opacity. Pinning the selector anchors it to the tone utility
		// itself, so a group-hover variant (`.group-hover\:text-base-content\/50`)
		// cannot satisfy it.
		mix := `.text-base-content\/` + step + `{color:color-mix(in oklab, var(--color-base-content) ` + step + `%, transparent)}`
		if !strings.Contains(src, mix) {
			t.Errorf("compiled css is missing the EmptyNote tone opacity rule %q — run `task build:ui` or the tone renders at full opacity", mix)
		}
	}
}

// TestCompiledCSSContainsHouseEdgeUtilities guards the card-outline contract
// from #79 end-to-end. components/ui/card_test.go pins that the components
// *emit* edge-line; this pins that the shipped stylesheet actually *defines* it
// and still raises --border for `prefers-contrast: more`.
//
// Both halves are needed: `edge-line` is hand-written component CSS, so a
// build that dropped the declaration (or a consumer bundle compiled without
// components/css/custom.css) would leave every card with no outline at all,
// while the class-string tests would stay green.
func TestCompiledCSSContainsHouseEdgeUtilities(t *testing.T) {
	css, err := fs.ReadFile(files, "static/css/app.css")
	if err != nil {
		t.Fatalf("read compiled css: %v", err)
	}
	src := string(css)

	// The width must be read from --border (with the 1px fallback). A literal
	// length here re-introduces the fixed-width defect.
	for _, rule := range []string{
		`.edge-line{border:var(--border,1px) solid var(--color-line)}`,
		`.edge-line-strong{border:var(--border,1px) solid var(--color-line-strong)}`,
	} {
		if !strings.Contains(src, rule) {
			t.Errorf("compiled css is missing the house edge rule %q — run `task build:ui`", rule)
		}
	}

	// Without this override the token stays 1px and the rule above is inert for
	// high-contrast users.
	const highContrast = `@media (prefers-contrast:more){:root{--border:2px}}`
	if !strings.Contains(src, highContrast) {
		t.Errorf("compiled css is missing %q — the house edge could never thicken under prefers-contrast: more", highContrast)
	}
}
