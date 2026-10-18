package css

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// toneTokens are the semantic text tones defined in tokens.css.
var toneTokens = []string{"body", "label", "muted", "subtle", "faint"}

// rawTonePattern matches opacity utilities that should have been replaced by a
// tone token. Values below 40 (/20, /30) stay raw: they are reserved for
// decorative icons.
var rawTonePattern = regexp.MustCompile(`text-base-content/(40|45|50|55|60|70|80)\b`)

// librarySources lists the directories whose sources must use tone tokens.
var librarySources = []string{"components", "pages", "galleryruntime"}

// TestTokensCSSExistsAndIsInline locks the token layer contract: the tone
// ladder must be declared with `@theme inline` so the color-mix expressions are
// emitted into the utilities and resolve against the active theme (required by
// scoped theme previews / the gallery theme studio).
func TestTokensCSSExistsAndIsInline(t *testing.T) {
	root := filepath.Join("..", "..")
	path := filepath.Join(root, "components", "css", filepath.Base(TokensCSS))
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("tokens.css not found at %s: %v", path, err)
	}
	src := string(data)

	if !strings.Contains(src, "@theme inline") {
		t.Error("tokens.css must use `@theme inline`; plain @theme freezes values at :root and breaks scoped themes")
	}
	if strings.Contains(src, "@theme {") {
		t.Error("tokens.css contains a plain `@theme` block; use `@theme inline`")
	}
	for _, tok := range toneTokens {
		if !strings.Contains(src, "--color-"+tok+":") {
			t.Errorf("tokens.css missing --color-%s", tok)
		}
		if !strings.Contains(src, "var(--color-base-content)") {
			t.Error("tone tokens must derive from var(--color-base-content)")
			break
		}
	}
	if !strings.Contains(src, "@utility eyebrow") {
		t.Error("tokens.css missing @utility eyebrow")
	}
}

// TestRegistryShipsTokens verifies every component package advertises the token
// layer, so consumers compiling their own stylesheet import it.
func TestRegistryShipsTokens(t *testing.T) {
	for pkg, assets := range Registry {
		found := false
		for _, f := range assets.CSSFiles {
			if f == TokensCSS {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("registry entry %q does not list %s", pkg, TokensCSS)
		}
	}
}

// TestNoRawToneUtilities prevents the tone drift from growing back: library
// sources must use text-body/text-label/text-muted/text-subtle/text-faint
// instead of hand-written text-base-content/<opacity> (>= 40).
func TestNoRawToneUtilities(t *testing.T) {
	root := filepath.Join("..", "..")
	checked := 0
	for _, dir := range librarySources {
		base := filepath.Join(root, dir)
		err := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return err
			}
			name := info.Name()
			if strings.HasSuffix(name, "_templ.go") || strings.HasSuffix(name, "_test.go") {
				return nil
			}
			if !strings.HasSuffix(name, ".templ") && !strings.HasSuffix(name, ".go") {
				return nil
			}
			data, readErr := os.ReadFile(path)
			if readErr != nil {
				return readErr
			}
			checked++
			for i, line := range strings.Split(string(data), "\n") {
				if rawTonePattern.MatchString(line) {
					rel, _ := filepath.Rel(root, path)
					t.Errorf("%s:%d uses a raw tone opacity; use a tone token (text-body…text-faint)", rel, i+1)
				}
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", base, err)
		}
	}
	if checked == 0 {
		t.Fatal("guard test scanned no files; path resolution is wrong")
	}
}

// TestTypographyTokens locks the typography role contract: body and heading
// families plus the weights/tracking/leading that make the roles legible.
// They live in a plain `:root` base-layer block (not @theme) and must drive
// real element rules, otherwise the theme studio would emit variables that
// nothing consumes.
func TestTypographyTokens(t *testing.T) {
	root := filepath.Join("..", "..")
	path := filepath.Join(root, "components", "css", filepath.Base(TokensCSS))
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("tokens.css not found at %s: %v", path, err)
	}
	src := string(data)

	for _, tok := range []string{
		"--font-body:", "--font-heading:",
		"--type-weight-body:", "--type-weight-heading:",
		"--type-tracking-heading:", "--type-leading-body:",
	} {
		if !strings.Contains(src, tok) {
			t.Errorf("tokens.css missing typography token %s", tok)
		}
	}

	// Element rules: a token nobody applies is dead configuration.
	for _, want := range []string{
		"html,\n  :where([data-theme])",
		"font-family: var(--font-body)",
		"font-family: var(--font-heading)",
		"font-weight: var(--type-weight-heading)",
		"@utility font-heading",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("tokens.css missing typography rule %q", want)
		}
	}

	// The rules must sit in the base layer: unlayered declarations would beat
	// `font-bold`/`leading-*`/`tracking-*` utilities.
	if !strings.Contains(src, "@layer base {") {
		t.Error("typography role rules must be declared in @layer base so utilities still win")
	}
}
