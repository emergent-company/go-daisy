package ui

import (
	"bytes"
	"html"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

// markdownRenderer renders GFM. Raw HTML is intentionally not enabled, so
// agent-authored markup is escaped rather than injected into the page.
var markdownRenderer = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
)

// renderMarkdown converts Markdown to safe HTML. Empty input returns "".
// On error it returns an escaped <pre> fallback.
func renderMarkdown(src string) string {
	if strings.TrimSpace(src) == "" {
		return ""
	}
	var buf bytes.Buffer
	if err := markdownRenderer.Convert([]byte(src), &buf); err != nil {
		return "<pre>" + html.EscapeString(src) + "</pre>"
	}
	return buf.String()
}
