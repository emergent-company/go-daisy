package ui

import (
	"bytes"
	"context"
	"regexp"
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

// toastMessageSpanClass returns the class list of the message <span> that Alpine
// fills via x-text. The regression guards below are exact about WHICH element
// carries the wrapping classes, since daisyUI's .alert is a single-row grid
// (grid-auto-flow: column) whose track otherwise refuses to shrink for a long
// message and pushes the close button out of the clipped box.
func toastMessageSpanClass(t *testing.T, html string) string {
	t.Helper()
	m := regexp.MustCompile(`<span class="([^"]*)" x-text="t\.message">`).FindStringSubmatch(html)
	if m == nil {
		t.Fatalf("message span not found in rendered toast queue\n%s", html)
	}
	return m[1]
}

// TestToastQueueLongMessageWrapsAndStaysScrollable locks the layout contract
// that fixes a long error toast being unreadable with an invisible close button:
//
//   - min-w-0 lets the message column shrink below its content width inside the
//     .alert grid instead of overflowing it.
//   - wrap-anywhere (overflow-wrap: anywhere) breaks unbroken tokens such as
//     SOURCES_OAUTH_GOOGLE_CLIENT_ID/SECRET and, unlike break-word, shrinks the
//     element's min-content size so the track actually collapses.
//   - max-h-[…] + overflow-y-auto cap the toast height at the viewport with
//     scrolling, so a very long message cannot grow past the screen.
func TestToastQueueLongMessageWrapsAndStaysScrollable(t *testing.T) {
	class := toastMessageSpanClass(t, renderToastQueue(t, ToastQueueProps{PauseOnHover: true, Countdown: true}))

	for _, want := range []string{
		"min-w-0",         // allow the grid track to shrink
		"wrap-anywhere",   // break long unbroken tokens / shrink min-content
		"max-h-[",         // sane vertical cap
		"overflow-y-auto", // scroll instead of growing past the viewport
	} {
		if !strings.Contains(class, want) {
			t.Errorf("message span class %q missing %q (long messages would be clipped or overflow the viewport)", class, want)
		}
	}
}

// TestToastQueueMessageWrappingIsNotCountdownGated ensures the wrapping fix
// applies to every queue configuration, not only the countdown one the Memory
// shell happens to use.
func TestToastQueueMessageWrappingIsNotCountdownGated(t *testing.T) {
	for _, props := range []ToastQueueProps{
		{},
		{Position: ToastQueueBottomCenter},
		{Countdown: true},
		{PauseOnHover: true, Countdown: true},
	} {
		class := toastMessageSpanClass(t, renderToastQueue(t, props))
		if !strings.Contains(class, "wrap-anywhere") || !strings.Contains(class, "min-w-0") {
			t.Errorf("props %+v: message span must wrap in every configuration, got class %q", props, class)
		}
	}
}

// TestToastQueueContainerWidth locks the reconsidered measure: the queue is
// wide enough for long messages but still capped to the viewport on small
// screens, for both anchor positions.
func TestToastQueueContainerWidth(t *testing.T) {
	for _, pos := range []string{ToastQueueTopRight, ToastQueueBottomCenter} {
		class := toastQueuePositionClass(pos)
		if !strings.Contains(class, "w-96") {
			t.Errorf("position %q: container should use the wider w-96 measure, got %q", pos, class)
		}
		if !strings.Contains(class, "max-w-[calc(100vw-2rem)]") {
			t.Errorf("position %q: container must keep its viewport cap, got %q", pos, class)
		}
	}
}
