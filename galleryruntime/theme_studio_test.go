package galleryruntime

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// TestThemeStudioRendersControls locks the studio's contract: it must render
// inside its own scoped preview theme (so editing never leaks into the shell),
// expose the apply/export affordances, and include the Alpine component.
func TestThemeStudioRendersControls(t *testing.T) {
	var buf bytes.Buffer
	if err := ThemeStudio().Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	html := buf.String()

	for _, want := range []string{
		`x-data="themeStudio()"`,
		`data-theme="gd-preview"`,
		"Apply to gallery",
		"Contrast",
		"Import CSS",
		"@plugin \"daisyui/theme\"",
		"function themeStudio()",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("studio output missing %q", want)
		}
	}

	// The preview theme appears twice: the scoped board element and the CSS
	// selector the studio generates. Nothing else may carry it, so the page
	// itself never switches the gallery theme.
	if got := strings.Count(html, `data-theme="gd-preview"`); got != 2 {
		t.Errorf("preview theme occurrences = %d, want 2 (board + selector)", got)
	}
	if strings.Contains(html, `x-data="themeStudio()" data-theme=`) {
		t.Error("studio root must not set data-theme")
	}
}

// scriptWraps reports whether needle sits inside a <script> element rather
// than leaking into the document as bare text.
func scriptWraps(html, needle string) bool {
	i := strings.Index(html, needle)
	if i < 0 {
		return false
	}
	return strings.LastIndex(html[:i], "<script") > strings.LastIndex(html[:i], "</script>")
}

// TestRenderSnippetPageInjectsCustomTheme guards the preview/iframe document:
// the theme runtime must be executable (inside a <script>) and the custom
// theme seam must be present, otherwise previews render unthemed in custom
// mode. Regression test for a runtime emitted as bare text between tags.
func TestRenderSnippetPageInjectsCustomTheme(t *testing.T) {
	html := renderSnippetPage("http://localhost", []string{"/static"}, "<p>preview</p>", false)

	if !scriptWraps(html, "window.__goDaisyTheme=(function()") {
		t.Error("ThemeRuntimeJS must be inside a <script> element")
	}
	for _, want := range []string{
		"gallery-custom-theme",  // style rule injected from localStorage
		"gallery-preview-theme", // pre-paint mode key
		"data-theme', 'custom'", // custom activation branch
	} {
		if !strings.Contains(html, want) {
			t.Errorf("preview document missing %q", want)
		}
	}
}
