package ui

import (
	"context"

	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/shared"
)

// ThemeRuntimeTag renders ThemeRuntimeJS() inside a <script> tag that carries
// the render-context CSP nonce.
//
// ThemeRuntimeJS itself is a bare JS body (no tags) because every consumer
// wraps it; those wrappers must use this helper instead of string concatenation
// so a strict `script-src 'nonce-…'` policy keeps working.
func ThemeRuntimeTag(ctx context.Context) string {
	return shared.ScriptTagNonce(templ.GetNonce(ctx), ThemeRuntimeJS())
}

// nonceAttrs returns the render-context CSP nonce as an attribute set, or nil
// when there is none. Spreading it keeps inline <script> output byte-identical
// for hosts that do not use CSP, instead of emitting nonce="".
func nonceAttrs(ctx context.Context) templ.Attributes { return shared.NonceAttrs(ctx) }
