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

// NonceAttr renders the ` nonce="…"` attribute fragment (leading space) for a
// nonce, HTML-escaped, or an empty string when the nonce is empty.
//
// Every hand-built script tag goes through this (directly or via
// ScriptTagNonce) so escaping and the omit-when-empty rule cannot drift between
// callers.
func NonceAttr(nonce string) string {
	if nonce == "" {
		return ""
	}
	return ` nonce="` + html.EscapeString(nonce) + `"`
}

// ScriptTagNonce renders an inline script body in a <script> tag carrying the
// given nonce. An empty nonce renders a plain tag, byte-identical to pre-CSP
// output.
func ScriptTagNonce(nonce, js string) string {
	return "<script" + NonceAttr(nonce) + ">" + js + "</script>"
}

// ScriptTag renders an inline script body wrapped in a <script> tag that carries
// the context nonce. Prefer this over hand-writing a <script> tag in Go string
// concatenation (e.g. theme runtime bootstraps).
func ScriptTag(ctx context.Context, js string) string {
	return ScriptTagNonce(templ.GetNonce(ctx), js)
}
