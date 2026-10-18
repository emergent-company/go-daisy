package galleryruntime

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

// TestFeedbackOverlayScript pins the emergent.feedback install snippet. The
// service renamed its bundle from feedback-overlay.js to emergent-feedback.js;
// the old path now answers 401, which the browser blocks as an opaque response,
// so the overlay silently never mounts. The label is the repo label a report
// lands under; gallery reports belong to the plain "feedback" label.
func TestFeedbackOverlayScript(t *testing.T) {
	var buf bytes.Buffer
	if err := feedbackOverlayScript().Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	html := buf.String()

	for _, want := range []string{
		`src="https://feedback.emergent-company.ai/emergent-feedback.js"`,
		`data-api="https://feedback.emergent-company.ai"`,
		`data-repo="emergent-company/go-daisy"`,
		`data-label="feedback"`,
		`async`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("feedback overlay script missing %q", want)
		}
	}

	if strings.Contains(html, "feedback-overlay.js") {
		t.Error("feedback overlay script still uses the retired feedback-overlay.js path")
	}
}

// TestGalleryPageMountsFeedbackOverlay pins the wiring as well as the snippet:
// deleting the @feedbackOverlayScript() call would leave the isolated component
// test green while the feedback tool silently disappears from every page.
func TestGalleryPageMountsFeedbackOverlay(t *testing.T) {
	var buf bytes.Buffer
	if err := GalleryPage("Test", "", nil, nil, templ.NopComponent).Render(context.Background(), &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	if !strings.Contains(buf.String(), "emergent-feedback.js") {
		t.Error("GalleryPage does not mount the feedback overlay script")
	}
}
