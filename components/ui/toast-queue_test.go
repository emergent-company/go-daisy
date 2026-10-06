package ui

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func renderToastQueue(t *testing.T, props ToastQueueProps) string {
	t.Helper()
	var buf bytes.Buffer
	if err := ToastQueueWithProps(props).Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestToastQueuePersistentToastContract locks the client contract that makes a
// duration<=0 toast persistent with a close button, while auto-dismissing
// toasts keep click-anywhere dismissal and the countdown bar. These are
// runtime-gated in Alpine, so the assertions are on the emitted directives,
// not on rendered class strings.
func TestToastQueuePersistentToastContract(t *testing.T) {
	got := renderToastQueue(t, ToastQueueProps{PauseOnHover: true, Countdown: true})

	for _, want := range []string{
		`aria-label="Dismiss"`,        // explicit close affordance
		`@click.stop="dismiss(t.id)"`, // close button dismisses persistent toast
		`t.duration <= 0`,             // close button only for persistent toasts
		`t.duration > 0`,              // bar + click-dismiss only for auto toasts
		`role="alert"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("toast queue missing %q\n%s", want, got)
		}
	}

	// Persistent toasts must not arm a pointer cursor / whole-toast click.
	if strings.Contains(got, `class="alert shadow-lg cursor-pointer`) {
		t.Errorf("static cursor-pointer must be runtime-gated (persistent toasts are not click-dismissable)\n%s", got)
	}
}
