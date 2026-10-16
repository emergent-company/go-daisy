package ui

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

// renderAlert renders Alert to its HTML string.
func renderAlert(t *testing.T, props AlertProps) string {
	t.Helper()
	var buf bytes.Buffer
	if err := Alert(props).Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// strBody returns a templ.Component that writes the given text as the body.
func strBody(text string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, text)
		return err
	})
}

// TestAlertMessageOnlyPath pins the message-only rendering so the Body slot is
// purely additive and existing callers stay byte-compatible.
func TestAlertMessageOnlyPath(t *testing.T) {
	got := renderAlert(t, AlertProps{Type: AlertInfo, Message: "Operation completed successfully."})
	if !strings.Contains(got, `<span>Operation completed successfully.</span>`) {
		t.Errorf("Message-only path must render the plain message <span>:\n%s", got)
	}
}

// TestAlertBodySlot pins the Body slot renders in place of the message span.
func TestAlertBodySlot(t *testing.T) {
	got := renderAlert(t, AlertProps{Type: AlertInfo, Body: strBody(`<button>Retry</button>`)})
	if !strings.Contains(got, `<button>Retry</button>`) {
		t.Errorf("Body slot must render the provided body:\n%s", got)
	}
	if strings.Contains(got, `Operation completed successfully.`) {
		t.Errorf("Body slot must replace the message span:\n%s", got)
	}
}
