package form

import (
	"strings"
	"testing"
)

// TestRichTextEditorLabelAssociation locks the Quill label path: a <label for>
// cannot label a Quill-managed div, so the label text also travels as data-label
// and the init script copies it onto Quill's contenteditable root as aria-label.
func TestRichTextEditorLabelAssociation(t *testing.T) {
	got := renderForm(t, RichTextEditor(RichTextEditorProps{ID: "rte", Label: "Body"}))

	for _, want := range []string{
		`for="rte"`,
		`id="rte"`,
		`data-quill`,
		`data-label="Body"`,
		`q.root.setAttribute('aria-label', lbl)`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("rich text editor missing %q in:\n%s", want, got)
		}
	}

	// Without a Label there is nothing to propagate.
	bare := renderForm(t, RichTextEditor(RichTextEditorProps{ID: "rte2"}))
	if strings.Contains(bare, "data-label=") {
		t.Errorf("unlabelled editor must not emit data-label:\n%s", bare)
	}
}
