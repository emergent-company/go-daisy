// Package devmode provides component boundary annotation for gallery dev tooling.
//
// In dev mode, [ComponentBoundary] wraps any [templ.Component] output in a
// display:contents <div> annotated with data-component and data-props attributes.
// This makes the component hierarchy visible in DevTools and enables the gallery's
// hover overlay, component tree panel, and annotated source view.
//
// For structural HTML elements that cannot legally contain a <div> wrapper
// (e.g. <thead>, <tbody>, <tr>, <td>, <th>), use [ElementBoundary] instead.
// It injects the data-component/data-props attributes directly onto the first
// opening tag emitted by the inner component rather than wrapping it.
//
// Usage in the gallery server:
//
//	ctx = devmode.WithDevMode(ctx)   // inject once per request
//	comp = devmode.ComponentBoundary("Button", props, ui.Button(props))
//	comp = devmode.ElementBoundary("TableRow", props, table.TableRow(id, hover))
//
// In production (when [IsDevMode] returns false), both functions are a
// zero-overhead passthrough — no wrapper element or extra attributes are emitted.
package devmode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/a-h/templ"
)

// contextKey is the unexported typed key used to store devmode state in context.
// Using a named type prevents collisions with any other package's context keys.
type contextKey struct{}

// WithDevMode returns a new context with dev mode enabled.
// Pass this context (or a context derived from it) to [templ.Component.Render]
// to activate component boundary annotations.
func WithDevMode(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey{}, true)
}

// IsDevMode reports whether dev mode is active in the given context.
func IsDevMode(ctx context.Context) bool {
	v, _ := ctx.Value(contextKey{}).(bool)
	return v
}

// ComponentBoundary wraps inner with a display:contents <div> annotated with
// data-component (the component name) and data-props (the props as JSON).
//
// The wrapper is only emitted when [IsDevMode] returns true for the render
// context — in all other cases the inner component is rendered unchanged.
//
// The props argument can be any JSON-serialisable value; nil is safe and
// results in data-props="null". For best results pass a map[string]any or
// the component's props struct.
func ComponentBoundary(name string, props any, inner templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if !IsDevMode(ctx) {
			return inner.Render(ctx, w)
		}

		propsJSON, err := json.Marshal(props)
		if err != nil {
			propsJSON = []byte("null")
		}

		// Emit the opening wrapper. display:contents makes the div invisible to
		// layout engines (flexbox, grid) while preserving its presence in the DOM
		// so DevTools and JavaScript can query [data-component] freely.
		if _, err := fmt.Fprintf(w,
			`<div data-component=%q data-props=%q style="display:contents">`,
			name, string(propsJSON),
		); err != nil {
			return err
		}

		if err := inner.Render(ctx, w); err != nil {
			return err
		}

		_, err = io.WriteString(w, `</div>`)
		return err
	})
}

// ElementBoundary annotates the first opening tag emitted by inner with
// data-component and data-props attributes, without adding any wrapper element.
//
// Use this for structural HTML elements that cannot legally contain a <div>
// wrapper, such as <thead>, <tbody>, <tr>, <td>, and <th>. The browser's
// HTML parser will strip a <div> placed directly inside a <table> or <tr>,
// making [ComponentBoundary] ineffective for these elements.
//
// Like [ComponentBoundary], this is a zero-overhead passthrough in production
// (when [IsDevMode] returns false).
func ElementBoundary(name string, props any, inner templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if !IsDevMode(ctx) {
			return inner.Render(ctx, w)
		}

		propsJSON, err := json.Marshal(props)
		if err != nil {
			propsJSON = []byte("null")
		}
		attrs := fmt.Sprintf(` data-component=%q data-props=%q`, name, string(propsJSON))

		// Render the inner component into a buffer so we can inject attributes
		// into the first opening tag before forwarding to the real writer.
		var buf bytes.Buffer
		if err := inner.Render(ctx, &buf); err != nil {
			return err
		}

		// Find the first '>' and insert the attributes just before it,
		// after any existing attributes on the tag.
		html := buf.String()
		idx := strings.Index(html, ">")
		if idx < 0 {
			// Shouldn't happen for well-formed templ output; fall through.
			_, err = io.WriteString(w, html)
			return err
		}
		// Handle self-closing tags: don't inject before '/>'
		insert := idx
		if insert > 0 && html[insert-1] == '/' {
			insert--
		}
		annotated := html[:insert] + attrs + html[insert:]
		_, err = io.WriteString(w, annotated)
		return err
	})
}
