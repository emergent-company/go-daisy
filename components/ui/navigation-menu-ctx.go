package ui

import "context"

type navigationMenuItemIDKey struct{}

// WithNavigationMenuItemID returns ctx annotated with the id of the enclosing
// NavigationMenuItem, so child parts (Trigger, Content) can derive their own
// aria-controls / content ids without an explicit parent reference.
func WithNavigationMenuItemID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, navigationMenuItemIDKey{}, id)
}

// navigationMenuItemIDFromCtx reads the enclosing NavigationMenuItem id, if any.
// The ok result is false when no NavigationMenuItem wraps the current render, in
// which case callers must fall back to an empty id and not emit an empty
// aria-controls / id attribute.
func navigationMenuItemIDFromCtx(ctx context.Context) (id string, ok bool) {
	id, ok = ctx.Value(navigationMenuItemIDKey{}).(string)
	return id, ok
}
