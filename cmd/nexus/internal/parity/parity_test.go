// Package parity guards the Nexus reference app against drifting back to
// hand-rolled markup in the areas that have already been converted to
// go-daisy components. It reads the .templ sources, so it catches regressions
// even when the offending markup still happens to render acceptably.
package parity

import (
	"os"
	"strings"
	"testing"
)

func read(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}

// templFunc returns the source of one top-level templ declaration, from its
// signature to the next top-level declaration.
func templFunc(t *testing.T, src, name string) string {
	t.Helper()
	marker := "templ " + name + "("
	start := strings.Index(src, marker)
	if start < 0 {
		t.Fatalf("templ %s not found", name)
	}
	rest := src[start+len(marker):]
	if end := strings.Index(rest, "\ntempl "); end >= 0 {
		rest = rest[:end]
	}
	return rest
}

func require(t *testing.T, src, what, needle string) {
	t.Helper()
	if !strings.Contains(src, needle) {
		t.Errorf("%s must contain %q", what, needle)
	}
}

func forbid(t *testing.T, src, what string, needles ...string) {
	t.Helper()
	for _, needle := range needles {
		if strings.Contains(src, needle) {
			t.Errorf("%s must not contain %q (use the library component instead)", what, needle)
		}
	}
}

// TestNexusShellComposesAdminShell pins the shell migration.
func TestNexusShellComposesAdminShell(t *testing.T) {
	src := read(t, "../shell/shell.templ")

	require(t, src, "Nexus shell", "@layout.AdminShell(")
	require(t, src, "Nexus shell", `layout.TopbarVariant("app"`)
	require(t, src, "Nexus shell", "BodyEnd: FeedbackOverlay()")

	// AdminShell owns the toggle checkbox, backdrop and the sidebar width,
	// topbar height and content padding rules.
	forbid(t, src, "Nexus shell",
		`aria-label="Navbar"`,
		`id="_layout-sidebar-toggle-trigger"`,
		`#_layout-sidebar {`,
		`#_layout-topbar {`,
		`#layout-content {`,
	)
}

// TestNexusChatUsesChatLayout pins the chat migration.
func TestNexusChatUsesChatLayout(t *testing.T) {
	src := read(t, "../handler/pages_chat.templ")

	require(t, src, "Nexus chat", "@pages.ChatLayout(")
	forbid(t, src, "Nexus chat",
		`class="chat chat-`,
		"data-simplebar",
		"chatContactNx",
		"chatContact(",
	)
}

// TestNexusCommerceTablesUseTableComponents pins the two converted lists.
func TestNexusCommerceTablesUseTableComponents(t *testing.T) {
	src := read(t, "../handler/pages_ecommerce.templ")

	for _, name := range []string{"ProductsList", "OrdersList"} {
		fn := templFunc(t, src, name)
		require(t, fn, name, "@table.TableCard(")
		require(t, fn, name, "@table.TableWithProps(")
		forbid(t, fn, name,
			`<table class="table`,
			"<table",
			"<thead",
			"<tbody",
			"<tr",
			"<td",
			`<input class="checkbox`,
		)
	}
}
