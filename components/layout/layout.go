// Package layout provides the top-level shell Templ components (Page, AppShell, Sidebar, Navbar).
package layout

import "github.com/a-h/templ"

// PageProps configures the full HTML document shell.
type PageProps struct {
	Title           string
	ThemeAttr       string
	// PinTheme omits ui.ThemeScript() from <head> so ThemeAttr is
	// authoritative. Use it for apps that ship a fixed brand theme and want the
	// persisted gallery palette/mode selection to never override it.
	PinTheme        bool
	Alpine          bool            // include Alpine.js <script> tag
	Morph           bool            // include idiomorph global (HTMX morphing is core in v4, no extension needed)
	Prefetch        bool            // include link prefetch on hover
	ViewTransitions bool            // enable View Transitions API (browser-native crossfade)
	ProgressBar     bool            // include auto-show progress bar during slow HTMX requests
	SSE             bool            // include htmx SSE extension
	WS              bool            // include htmx WebSocket extension
	BodyEnd         templ.Component // optional component rendered at the end of <body> (e.g. dev tools, feedback overlay)
}
