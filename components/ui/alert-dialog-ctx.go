package ui

import "context"

type alertDialogIDKey struct{}

// WithAlertDialogID returns ctx annotated with the id of the enclosing
// AlertDialog, so child parts (Title, Description) can derive their own
// aria-labelledby/aria-describedby ids without an explicit parent reference.
func WithAlertDialogID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, alertDialogIDKey{}, id)
}

// alertDialogIDFromCtx reads the enclosing AlertDialog id, if any. The ok
// result is false when no AlertDialog wraps the current render, in which case
// callers must fall back to an empty id and not emit an empty id attribute.
func alertDialogIDFromCtx(ctx context.Context) (id string, ok bool) {
	id, ok = ctx.Value(alertDialogIDKey{}).(string)
	return id, ok
}
