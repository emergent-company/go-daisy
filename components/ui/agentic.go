package ui

import (
	"context"
	"strings"

	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/devmode"
)

// agenticAttrs merges the component's dev-mode annotation with caller attrs.
// Caller attrs win, so consumers can override defaults on the root element.
func agenticAttrs(ctx context.Context, component string, extra templ.Attributes) templ.Attributes {
	out := templ.Attributes{}
	for k, v := range devmode.Attrs(ctx, component) {
		out[k] = v
	}
	for k, v := range extra {
		out[k] = v
	}
	return out
}

// humanizeToolName turns a raw tool id ("schema-list-installed") into a quiet
// one-line label ("Schema list installed"), matching the emergent.memory chat UI.
func humanizeToolName(name string) string {
	t := strings.TrimSpace(name)
	if t == "" {
		return "Tool"
	}
	t = strings.NewReplacer("_", " ", ".", " ", "-", " ").Replace(t)
	t = strings.Join(strings.Fields(t), " ")
	t = strings.ToLower(t)
	if t == "" {
		return "Tool"
	}
	return strings.ToUpper(t[:1]) + t[1:]
}
