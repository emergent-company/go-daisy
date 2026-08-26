// Package stream provides Go helpers for generating Turbo-Stream-like HTML
// fragments and real-time server-sent updates. Use with HTMX SSE or WebSocket
// extensions for live page updates without polling.
//
// Usage:
//
//	stream.Append("#messages", messageRow)
//	stream.Replace("#counter", counterBadge)
//	stream.Remove("#alert-42")
//	stream.RefreshMorph()
//
// Each function returns HTML that HTMX processes as an out-of-band swap.
// When sent over SSE or WebSocket, the client automatically applies the mutation.
package stream

import (
	"context"
	"fmt"
	"html"
	"io"
	"strings"

	"github.com/a-h/templ"
)

// Action represents a DOM mutation type.
type Action string

const (
	ActionAppend  Action = "append"
	ActionPrepend Action = "prepend"
	ActionReplace Action = "replace"
	ActionUpdate  Action = "update"
	ActionRemove  Action = "remove"
	ActionBefore  Action = "before"
	ActionAfter   Action = "after"
)

// Builder accumulates stream actions into a single HTML response.
type Builder struct {
	w io.Writer
}

// NewBuilder creates a stream builder writing to w.
func NewBuilder(w io.Writer) *Builder {
	return &Builder{w: w}
}

// Write writes a stream action directly to the builder.
func (b *Builder) Write(action Action, target string, comp templ.Component) error {
	return b.writeStream(action, target, "", comp)
}

// WriteMorph writes a stream action with DOM morphing.
func (b *Builder) WriteMorph(action Action, target string, comp templ.Component) error {
	return b.writeStream(action, target, "morph", comp)
}

func (b *Builder) writeStream(action Action, target, method string, comp templ.Component) error {
	var sb strings.Builder
	sb.WriteString(`<div hx-swap-oob="`)
	if method != "" {
		sb.WriteString(string(action))
		sb.WriteString(":")
		sb.WriteString(method)
		sb.WriteString(":")
	} else {
		sb.WriteString(string(action))
	}
	sb.WriteString(target)
	sb.WriteString(`">`)
	if comp != nil && action != ActionRemove {
		if err := comp.Render(context.Background(), &sb); err != nil {
			return fmt.Errorf("stream: render component: %w", err)
		}
	}
	sb.WriteString(`</div>`)
	_, err := b.w.Write([]byte(sb.String()))
	return err
}

// Remove removes the element with the given target selector.
func (b *Builder) Remove(target string) error {
	return b.Write(ActionRemove, target, nil)
}

// Refresh writes a refresh stream action (triggers full page reload with morph).
func (b *Builder) Refresh() error {
	return b.writeStream("refresh", "", "", nil)
}

// String returns the accumulated stream HTML as a string. Panics on error.
func (b *Builder) String() string {
	if sb, ok := b.w.(*strings.Builder); ok {
		return sb.String()
	}
	return ""
}

// ── Convenience functions (no builder) ─────────────────────────────────────

// Append returns HTML that appends a component to the target element.
func Append(target string, comp templ.Component) (string, error) {
	return actionString(ActionAppend, target, "", comp)
}

// Prepend returns HTML that prepends a component to the target element.
func Prepend(target string, comp templ.Component) (string, error) {
	return actionString(ActionPrepend, target, "", comp)
}

// Replace returns HTML that replaces the target element with the component.
func Replace(target string, comp templ.Component) (string, error) {
	return actionString(ActionReplace, target, "", comp)
}

// Update returns HTML that replaces the innerHTML of the target element.
func Update(target string, comp templ.Component) (string, error) {
	return actionString(ActionUpdate, target, "", comp)
}

// Before returns HTML that inserts the component before the target element.
func Before(target string, comp templ.Component) (string, error) {
	return actionString(ActionBefore, target, "", comp)
}

// After returns HTML that inserts the component after the target element.
func After(target string, comp templ.Component) (string, error) {
	return actionString(ActionAfter, target, "", comp)
}

// RemoveTarget returns HTML that removes the target element from the DOM.
func RemoveTarget(target string) string {
	return fmt.Sprintf(`<div hx-swap-oob="delete:%s"></div>`, html.EscapeString(target))
}

// RefreshMorph returns HTML that triggers a full page refresh with DOM morphing.
func RefreshMorph() string {
	return `<div hx-swap-oob="refresh"></div>`
}

// Morph returns HTML that morphs the target element with the component.
func Morph(target string, comp templ.Component) (string, error) {
	return actionString("morph", target, "", comp)
}

func actionString(action Action, target, method string, comp templ.Component) (string, error) {
	var sb strings.Builder
	b := NewBuilder(&sb)
	if err := b.writeStream(action, target, method, comp); err != nil {
		return "", err
	}
	return sb.String(), nil
}

// TagSSE returns a script tag loading the HTMX SSE extension.
func TagSSE() templ.Component {
	return scriptTag("/static/js/htmx-sse.js")
}

// TagWS returns a script tag loading the HTMX WebSocket extension.
func TagWS() templ.Component {
	return scriptTag("/static/js/htmx-ws.js")
}

func scriptTag(src string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := fmt.Fprintf(w, `<script src="%s"></script>`, html.EscapeString(src))
		return err
	})
}
