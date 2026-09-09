// Package css catalogs the CSS each go-daisy component package depends on so
// consumers can prune daisyUI and include the co-located custom CSS when they
// compile their own stylesheet instead of loading the monolithic bundle.
//
// The daisyUI module names are the `include`/`exclude` names accepted by
// daisyUI 5's `@plugin "daisyui"` config (see daisyui.com/docs/config).
// The sets below are derived by static analysis of each package's components
// (class tokens in .templ/.go sources); a consumer still owns the final
// include/exclude decision and verifies it by rendering.
package css

import (
	"path"
	"sort"
)

// CustomCSS is the co-located custom component CSS that go-daisy ships for
// consumers (sidebar/layout/topbar/alpine/view-transition/icon rules). It is
// imported by go-daisy's own assets/app.css and can be @imported by a
// consumer's Tailwind build from the vendored module:
//
//	@import "../components/css/custom.css";
const CustomCSS = "css/custom.css"

// PackageAssets describes the CSS a go-daisy component package needs.
type PackageAssets struct {
	// DaisyUIModules lists the daisyUI modules the package's components
	// reference (daisyUI config include/exclude names).
	DaisyUIModules []string
	// CSSFiles lists co-located custom CSS files the package's components rely
	// on (paths relative to the components/ root).
	CSSFiles []string
}

// Registry maps each component package directory name to its CSS assets.
var Registry = map[string]PackageAssets{
	"alpine": {DaisyUIModules: nil, CSSFiles: []string{CustomCSS}},
	"form": {DaisyUIModules: []string{
		"badge", "button", "calendar", "card", "dropdown", "fieldset", "input",
		"join", "label", "loading", "menu", "progress", "textarea", "validator",
	}},
	"layout": {DaisyUIModules: []string{
		"avatar", "button", "collapse", "dropdown", "input", "menu", "progress",
		"select", "toast",
	}, CSSFiles: []string{CustomCSS}},
	"modal":      {DaisyUIModules: []string{"loading", "modal"}},
	"nav":        {DaisyUIModules: []string{"avatar", "badge", "breadcrumbs", "button", "dock", "drawer", "dropdown", "input", "mask", "menu", "modal", "select", "steps"}, CSSFiles: []string{CustomCSS}},
	"schemaform": {DaisyUIModules: []string{"badge"}},
	"table":      {DaisyUIModules: []string{"button", "card", "checkbox", "dropdown", "input", "join", "loading", "select", "table"}},
	"ui": {DaisyUIModules: []string{
		"alert", "avatar", "badge", "button", "card", "collapse", "countdown",
		"diff", "divider", "dropdown", "fieldset", "hero", "indicator", "input",
		"join", "kbd", "label", "link", "list", "loading", "mask", "menu",
		"modal", "progress", "rating", "skeleton", "stack", "stat", "steps",
		"swap", "tab", "timeline", "tooltip", "validator",
	}, CSSFiles: []string{CustomCSS}},
}

// ModulesFor returns the union of daisyUI modules used by the given component
// packages (directory names under components/, or full import paths whose last
// segment is a package directory name), sorted and deduplicated — the set a
// consumer's `@plugin "daisyui"` include/exclude must cover when it compiles
// those packages. Unknown package names are ignored.
func ModulesFor(pkgs ...string) []string {
	seen := map[string]bool{}
	var out []string
	for _, p := range pkgs {
		pa, ok := Registry[path.Base(p)]
		if !ok {
			continue
		}
		for _, m := range pa.DaisyUIModules {
			if !seen[m] {
				seen[m] = true
				out = append(out, m)
			}
		}
	}
	sort.Strings(out)
	return out
}

// CSSFilesFor returns the union of co-located custom CSS files (relative to
// the components/ root) the given packages rely on, sorted and deduplicated.
func CSSFilesFor(pkgs ...string) []string {
	seen := map[string]bool{}
	var out []string
	for _, p := range pkgs {
		pa, ok := Registry[path.Base(p)]
		if !ok {
			continue
		}
		for _, f := range pa.CSSFiles {
			if !seen[f] {
				seen[f] = true
				out = append(out, f)
			}
		}
	}
	sort.Strings(out)
	return out
}
