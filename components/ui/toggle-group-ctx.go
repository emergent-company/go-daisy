package ui

import (
	"context"
	"encoding/json"
)

// toggleGroupDefaults carries the inherited variant/size/disabled state a
// ToggleGroup shares with its child ToggleButtons. Stored on the render context
// so children can pick it up without an explicit parent reference.
type toggleGroupDefaults struct {
	variant  ToggleVariant
	size     ToggleSize
	disabled bool
}

type toggleGroupDefaultsKey struct{}

// WithToggleGroupDefaults returns ctx annotated with the group-level defaults a
// ToggleGroup should propagate to its child ToggleButtons.
func WithToggleGroupDefaults(ctx context.Context, variant ToggleVariant, size ToggleSize, disabled bool) context.Context {
	return context.WithValue(ctx, toggleGroupDefaultsKey{}, toggleGroupDefaults{variant: variant, size: size, disabled: disabled})
}

// toggleGroupDefaultsFromCtx reads the inherited group defaults, if any.
func toggleGroupDefaultsFromCtx(ctx context.Context) (variant ToggleVariant, size ToggleSize, disabled bool, ok bool) {
	d, ok := ctx.Value(toggleGroupDefaultsKey{}).(toggleGroupDefaults)
	if !ok {
		return "", "", false, false
	}
	return d.variant, d.size, d.disabled, true
}

// toggleGroupValueJSON marshals a value slice to its JSON representation for the
// root data-gd-toggle-value attribute. Nil and empty slices both marshal to "[]".
func toggleGroupValueJSON(values []string) string {
	if len(values) == 0 {
		return "[]"
	}
	b, err := json.Marshal(values)
	if err != nil {
		return "[]"
	}
	return string(b)
}
