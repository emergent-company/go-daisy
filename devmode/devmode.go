// Package devmode provides component boundary annotation for gallery dev tooling.
//
// In dev mode, [ComponentBoundary] wraps any [templ.Component] output in a
// display:contents <div> annotated with data-component and data-props attributes.
// This makes the component hierarchy visible in DevTools and enables the gallery's
// hover overlay, component tree panel, and annotated source view.
//
// Usage in the gallery server:
//
//	ctx = devmode.WithDevMode(ctx)   // inject once per request
//	comp = devmode.ComponentBoundary("Button", props, ui.Button(props))
//
// In production (when [IsDevMode] returns false), [ComponentBoundary] is a
// zero-overhead passthrough — no wrapper element is emitted.
package devmode

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

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
