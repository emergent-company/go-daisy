package ui

import (
	"context"

	"github.com/a-h/templ"
)

// ThemeRuntimeTag renders ThemeRuntimeJS() inside a <script> tag that carries
// the render-context CSP nonce.
//
// ThemeRuntimeJS itself is a bare JS body (no tags) because every consumer
// wraps it; those wrappers must use this helper instead of string concatenation
// so a strict `script-src 'nonce-…'` policy keeps working.
func ThemeRuntimeTag(ctx context.Context) string {
	js := ThemeRuntimeJS()
	if n := templ.GetNonce(ctx); n != "" {
		return "<script nonce=\"" + n + "\">" + js + "</script>"
	}
	return "<script>" + js + "</script>"
}
