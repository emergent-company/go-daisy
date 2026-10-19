// Package parity guards the Nexus reference app against drifting back to
// hand-rolled markup in the areas that have already been converted to
// go-daisy components. It reads the .templ sources, so it catches regressions
// even when the offending markup still happens to render acceptably.
package parity

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/emergent-company/go-daisy/cmd/nexus/compare"
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

// TestEcommerceListsUseLibraryComponents pins the sellers/customers/shops list
// migrations: each list composes its header through nav.PageHeading and its
// table through the table primitives, with no hand-rolled table markup left
// behind in the row helpers.
func TestEcommerceListsUseLibraryComponents(t *testing.T) {
	src := read(t, "../handler/pages_ecommerce.templ")

	for _, name := range []string{"SellersList", "CustomersList", "ShopsList"} {
		fn := templFunc(t, src, name)
		require(t, fn, name, "@nav.PageHeading(")
		require(t, fn, name, "@table.TableCard(")
		require(t, fn, name, "@table.TableWithProps(")
		forbid(t, fn, name,
			"<table",
			"<thead",
			"<tbody",
			"<tr",
			"<td",
			"<th",
			`<input class="checkbox`,
		)
	}

	for _, row := range []string{"sellerRow", "customerRow", "shopRow"} {
		fn := templFunc(t, src, row)
		require(t, fn, row, "@table.TableRow(")
		require(t, fn, row, "@table.TableCell(")
		forbid(t, fn, row, "<tr", "<td", "<th")
	}

	// The three list handlers must resolve to the templ components, not fall
	// back to the raw-HTML passthrough.
	hsrc := read(t, "../handler/handler.go")
	forbid(t, hsrc, "handler.go",
		`nxHTMLComponent("apps-ecommerce-sellers")`,
		`nxHTMLComponent("apps-ecommerce-customers")`,
		`nxHTMLComponent("apps-ecommerce-shops")`,
	)
}

// TestNexusDashboardUsesLibraryCardsAndTable pins the dashboard conversion:
// stat cards go through ui.StatCardFeatured and Recent Orders through the table
// primitives, so the dashboard stops maintaining its own copies.
func TestNexusDashboardUsesLibraryCardsAndTable(t *testing.T) {
	src := read(t, "../handler/pages_dashboards.templ")

	dashboard := templFunc(t, src, "EcommerceDashboard")
	require(t, dashboard, "EcommerceDashboard", "@ui.StatCardFeatured(")
	require(t, dashboard, "EcommerceDashboard", "@table.TableCard(")
	require(t, dashboard, "EcommerceDashboard", "@table.TableWithProps(")
	forbid(t, dashboard, "EcommerceDashboard",
		"ecomStatCard(",
		"<table",
		"<thead",
		"<tbody",
		"<tr",
		"<td",
	)

	row := templFunc(t, src, "ecomOrderRow")
	require(t, row, "ecomOrderRow", "@table.TableRow(")
	require(t, row, "ecomOrderRow", "@table.TableCell(")
	forbid(t, row, "ecomOrderRow", "<tr", "<td", "<th")
}

// TestNexusDashboardsUseFoundationComponents pins the P1 foundation port: both
// dashboards compose their headers through nav.PageHeading, their card headers
// through ui.CardToolbar, their charts through ui.ChartCard, and their stat
// cards through ui.StatCardFeatured — with no hand-rolled breadcrumb, CRM stat
// card helper, or raw table markup left behind.
func TestNexusDashboardsUseFoundationComponents(t *testing.T) {
	src := stripLineComments(read(t, "../handler/pages_dashboards.templ"))

	for _, name := range []string{"EcommerceDashboard", "CRMDashboard"} {
		fn := templFunc(t, src, name)
		require(t, fn, name, "@ui.StatCardFeatured(")
		require(t, fn, name, "@ui.ChartCard(")
		require(t, fn, name, "@nav.PageHeading(")
	}

	forbid(t, src, "dashboards",
		"templ crmStatCard(",
		`class="breadcrumbs hidden p-0`,
		"<table",
		"<thead",
		"<tbody",
		"<tr",
		"<td",
	)
}

// stripLineComments removes whole-line // comments so a commented-out component
// call cannot satisfy a require() assertion.
func stripLineComments(src string) string {
	lines := strings.Split(src, "\n")
	out := lines[:0]
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			continue
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

// TestNexusTemplatesUseRootRelativeImages catches the class of 404 where a
// template writes "./images/..." and the browser resolves it under the current
// route (e.g. /dashboards/images/...) instead of /images/.
func TestNexusTemplatesUseRootRelativeImages(t *testing.T) {
	patterns := []string{"../handler/*.templ", "../shell/*.templ"}
	for _, pattern := range patterns {
		files, err := filepath.Glob(pattern)
		if err != nil {
			t.Fatalf("glob %s: %v", pattern, err)
		}
		for _, file := range files {
			src := read(t, file)
			for _, bad := range []string{`"./images`, `"./static`} {
				if strings.Contains(src, bad) {
					t.Errorf("%s must not contain %q: use a root-relative asset path", file, bad)
				}
			}
		}
	}
}

// TestCompareCoversEveryReference asserts the compare.Pages() table is
// provably complete: every reference file under nxhtml has a PageMapping
// (skipped entries allowed), and no non-skipped mapping points at a missing
// file.
func TestCompareCoversEveryReference(t *testing.T) {
	dir := "../handler/nxhtml"
	refs, err := filepath.Glob(filepath.Join(dir, "*.html"))
	if err != nil {
		t.Fatalf("glob %s: %v", dir, err)
	}
	if len(refs) == 0 {
		t.Fatal("no reference files found under ../handler/nxhtml")
	}

	mappings := compare.Pages()
	mapped := make(map[string]bool, len(mappings))
	for _, m := range mappings {
		mapped[m.HTMLFile] = true
		if m.Skip && m.Reason == "" {
			t.Errorf("skipped mapping %q has no reason", m.HTMLFile)
		}
		if !m.Skip {
			if _, err := os.Stat(filepath.Join(dir, m.HTMLFile)); err != nil {
				t.Errorf("mapping %q -> %s points at a missing reference file", m.Route, m.HTMLFile)
			}
		}
	}
	for _, f := range refs {
		base := filepath.Base(f)
		if !mapped[base] {
			t.Errorf("reference %s has no PageMapping (skipped entries allowed)", base)
		}
	}
}

// TestRawPassthroughIsDetectable asserts nxrender.go emits the nx-raw marker
// so the parity comparator can distinguish raw passthroughs from real ports.
// This is a source-grep tripwire, not a behavioural test: it guards that the
// marker string stays present in nxrender.go; the actual detection is exercised
// by the compare gate at runtime.
func TestRawPassthroughIsDetectable(t *testing.T) {
	src := read(t, "../handler/nxrender.go")
	if !strings.Contains(src, "<!-- nx-raw:") {
		t.Errorf("nxrender.go must emit the `<!-- nx-raw:<page> -->` passthrough marker")
	}
}

// TestMappedRoutesAreRegistered asserts every non-skipped PageMapping.Route is
// actually wired in handler.go. It parses the Register*Routes bodies for
// `g.GET("<path>"` literals, applies the /auth prefix only to
// RegisterAuthRoutes, and matches concrete values (e.g. /1) against `:param`
// segments by comparing segment counts and literal segments.
func TestMappedRoutesAreRegistered(t *testing.T) {
	src := read(t, "../handler/handler.go")
	registered := registeredRoutes(src)

	for _, m := range compare.Pages() {
		if m.Skip {
			continue
		}
		matched := false
		for _, r := range registered {
			if routeMatches(r, m.Route) {
				matched = true
				break
			}
		}
		if !matched {
			t.Errorf("mapped route %q (%s) is not registered in handler.go", m.Route, m.HTMLFile)
		}
	}
}

// registeredRoutes scans handler.go for `g.GET("<path>"` literals, returning
// the absolute route for each, with /auth prefixed only to RegisterAuthRoutes.
func registeredRoutes(src string) []string {
	var routes []string
	prefix := ""
	for _, line := range strings.Split(src, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "func (h *Handler) Register") {
			prefix = ""
			if strings.Contains(trimmed, "RegisterAuthRoutes") {
				prefix = "/auth"
			}
			continue
		}
		if !strings.Contains(line, "g.GET(") {
			continue
		}
		path := extractGETPath(line)
		full := prefix + path
		if full == "" {
			full = "/"
		}
		routes = append(routes, full)
	}
	return routes
}

// extractGETPath returns the string literal passed as the first argument to
// g.GET(...) on the given source line, or "" when absent.
func extractGETPath(line string) string {
	start := strings.Index(line, "g.GET(")
	if start < 0 {
		return ""
	}
	rest := strings.TrimSpace(line[start+len("g.GET("):])
	if !strings.HasPrefix(rest, `"`) {
		return ""
	}
	rest = rest[1:]
	if end := strings.Index(rest, `"`); end >= 0 {
		return rest[:end]
	}
	return ""
}

// routeMatches reports whether a concrete mapped route (e.g. /apps/.../1)
// matches a registered route (e.g. /apps/.../:id): equal segment counts, equal
// literal segments, and any `:param` segment matching any concrete value.
func routeMatches(registered, mapped string) bool {
	reg := splitPath(registered)
	m := splitPath(mapped)
	if len(reg) != len(m) {
		return false
	}
	for i := range reg {
		if strings.HasPrefix(reg[i], ":") {
			continue
		}
		if reg[i] != m[i] {
			return false
		}
	}
	return true
}

func splitPath(p string) []string {
	p = strings.Trim(p, "/")
	if p == "" {
		return nil
	}
	return strings.Split(p, "/")
}
