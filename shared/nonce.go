package shared

import (
	"context"
	"html"

	"github.com/a-h/templ"
)

// NonceAttrs returns a templ.Attributes set carrying the CSP nonce from the
// render context (set with templ.WithNonce) for inline <script> or <style>
// elements.
//
// Returns nil when the context has no nonce, so a component rendered outside a
// CSP host emits exactly the same HTML as before. Spread it on the element.
func NonceAttrs(ctx context.Context) templ.Attributes {
	if n := templ.GetNonce(ctx); n != "" {
		return templ.Attributes{"nonce": n}
	}
	return nil
}

// ScriptTag renders an inline script body wrapped in a <script> tag that carries
// the context nonce. Prefer this over hand-writing a <script> tag in Go string
// concatenation (e.g. theme runtime bootstraps).
func ScriptTag(ctx context.Context, js string) string {
	if n := templ.GetNonce(ctx); n != "" {
		return `<script nonce="` + html.EscapeString(n) + `">` + js + "</script>"
	}
	return "<script>" + js + "</script>"
}
