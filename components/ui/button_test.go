package ui

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// renderButton renders Button to its HTML string.
func renderButton(t *testing.T, props ButtonProps) string {
	t.Helper()
	var buf bytes.Buffer
	if err := Button(props).Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestButtonZeroTypeRendersPlainButton pins the default button type. The zero
// value of ButtonType is "", which HTML treats as submit — so an unset Type must
// render type="button", never type="" (implicit submit). This guards against the
// live defect where any ui.Button without an explicit Type inside a <form>
// silently became an implicit submit button.
func TestButtonZeroTypeRendersPlainButton(t *testing.T) {
	got := renderButton(t, ButtonProps{})

	if !strings.Contains(got, `type="button"`) {
		t.Errorf("Button with zero Type must render type=\"button\":\n%s", got)
	}
	for _, unwanted := range []string{`type=""`, `type="submit"`} {
		if strings.Contains(got, unwanted) {
			t.Errorf("Button renders %s — it would submit an enclosing form:\n%s", unwanted, got)
		}
	}
}

// TestButtonExplicitSubmitPreserved pins that an explicit ButtonTypeSubmit still
// renders type="submit", and an explicit ButtonTypeButton renders type="button".
func TestButtonExplicitSubmitPreserved(t *testing.T) {
	if got := renderButton(t, ButtonProps{Type: ButtonTypeSubmit}); !strings.Contains(got, `type="submit"`) {
		t.Errorf("Button with ButtonTypeSubmit must render type=\"submit\":\n%s", got)
	}
	if got := renderButton(t, ButtonProps{Type: ButtonTypeButton}); !strings.Contains(got, `type="button"`) {
		t.Errorf("Button with ButtonTypeButton must render type=\"button\":\n%s", got)
	}
}
