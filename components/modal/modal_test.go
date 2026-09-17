package modal

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

// renderModal renders a component to its HTML string for assertions.
func renderModal(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestModalVanillaBackdropTargetsOwnDialog guards the regression where the
// vanilla Modal backdrop form hardcoded onclick="document.getElementById('modal').remove()".
// With a non-default ModalProps.ID the handler must target its own <dialog> via
// closest('dialog'), not a hardcoded id (which either throws on a missing node or
// destroys an unrelated element that happens to have id=modal).
func TestModalVanillaBackdropTargetsOwnDialog(t *testing.T) {
	html := renderModal(t, Modal(ModalProps{ID: "custom-modal", Title: "Confirm"}))

	if !strings.Contains(html, `id="custom-modal"`) {
		t.Errorf("dialog missing custom id:\n%s", html)
	}
	if !strings.Contains(html, `this.closest('dialog').remove()`) {
		t.Errorf("backdrop handler must target its own dialog via closest('dialog'):\n%s", html)
	}
	if strings.Contains(html, `document.getElementById('modal').remove()`) {
		t.Errorf("backdrop handler must not hardcode id 'modal':\n%s", html)
	}
}

// TestModalAlpineBranchUnaffected locks the Alpine driver output: it keeps the
// non-destructive x-show/@click toggle and must not gain the destructive
// remove() handler.
func TestModalAlpineBranchUnaffected(t *testing.T) {
	html := renderModal(t, Modal(ModalProps{Title: "T", Driver: ModalDriverAlpine}))

	if !strings.Contains(html, `x-show="open"`) {
		t.Errorf("alpine branch missing x-show:\n%s", html)
	}
	if !strings.Contains(html, `@click="open = false"`) {
		t.Errorf("alpine backdrop must use @click=open=false:\n%s", html)
	}
	if strings.Contains(html, `this.closest('dialog').remove()`) {
		t.Errorf("alpine branch must not emit destructive remove():\n%s", html)
	}
}

// TestModalStimulusBranchUnaffected locks the Stimulus driver output: it keeps
// the data-controller wiring and a plain method=dialog backdrop form.
func TestModalStimulusBranchUnaffected(t *testing.T) {
	html := renderModal(t, Modal(ModalProps{Title: "T", Driver: ModalDriverStimulus}))

	if !strings.Contains(html, `data-controller="modal"`) {
		t.Errorf("stimulus branch missing controller:\n%s", html)
	}
	if !strings.Contains(html, `data-modal-target="dialog"`) {
		t.Errorf("stimulus branch missing dialog target:\n%s", html)
	}
	if strings.Contains(html, `this.closest('dialog').remove()`) {
		t.Errorf("stimulus branch must not emit destructive remove():\n%s", html)
	}
}

// TestModalDefaultOutput locks the vanilla default rendering (empty ID → id="modal")
// so the backdrop-handler fix cannot silently change the default markup.
func TestModalDefaultOutput(t *testing.T) {
	html := renderModal(t, Modal(ModalProps{Title: "Hello"}))

	for _, want := range []string{
		`id="modal"`,
		`class="modal modal-open"`,
		`role="dialog"`,
		`aria-modal="true"`,
		`aria-labelledby="modal-title-modal"`,
		`class="modal-box"`,
		`this.closest('dialog').remove()`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("default modal missing %q:\n%s", want, html)
		}
	}
}
