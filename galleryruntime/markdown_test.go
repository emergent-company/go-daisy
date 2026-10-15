package galleryruntime

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

// buttonTempl returns a tiny templ component used as fixture content for the
// Templ-variant component in the markdown tests.
func buttonTempl() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<button class="btn">Click</button>`)
		return err
	})
}

func TestLLMsAndComponentMarkdown(t *testing.T) {
	components := []GalleryComponent{
		{
			Slug:        "button",
			Name:        "Button",
			Category:    CategoryComponents,
			Subcategory: "Buttons",
			Description: "A clickable button.",
			Variants: []GalleryStory{
				{Name: "Default", Description: "Primary button", Templ: buttonTempl()},
				{Name: "Loading", Description: "Shows a spinner", Templ: buttonTempl()},
			},
			Tokens: []DesignToken{
				{Label: "Padding scale", CSSVar: "--btn-p", Type: TokenTypeRange, Default: "0.75rem"},
				{Label: "Size", Group: "Component", Type: TokenTypeSelect, QueryParam: "size", Options: []TokenOption{{Value: "btn-lg", Label: "LG"}, {Value: "btn-sm", Label: "SM"}}},
			},
		},
		{
			Slug:        "badge",
			Name:        "Badge",
			Category:    CategoryComponents,
			Subcategory: "Badges",
			Description: "A status badge.",
			HTML:        `<span class="badge">Badge</span>`,
		},
	}

	e := echo.New()
	newGalleryHandler("Test Gallery", nil, components, nil, false).register(e)

	// /llms.txt lists every component with its markdown link.
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/llms.txt", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("llms.txt status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != mimeLLMsText {
		t.Fatalf("llms.txt Content-Type = %q, want %q", ct, mimeLLMsText)
	}
	llms := rec.Body.String()
	if !strings.Contains(llms, "# Test Gallery") {
		t.Errorf("llms.txt missing H1 title: %q", llms)
	}
	for _, comp := range components {
		want := "[" + comp.Name + "](/gallery/docs/" + comp.Slug + ".md)"
		if !strings.Contains(llms, want) {
			t.Errorf("llms.txt missing entry %q", want)
		}
	}

	// /gallery/docs/<slug>.md renders the component reference.
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/gallery/docs/button.md", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("markdown status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != mimeMarkdownText {
		t.Fatalf("markdown Content-Type = %q, want %q", ct, mimeMarkdownText)
	}
	md := rec.Body.String()
	for _, want := range []string{
		"# Button",
		"Category: Components",
		"Subcategory: Buttons",
		"## Variants",
		"**Default**",
		"**Loading**",
		"## Design tokens",
		"--btn-p",
		"Padding scale",
		"?size",
		"btn-lg",
		"## Example",
		"Click",
		"[Open in gallery](/gallery/button)",
	} {
		if !strings.Contains(md, want) {
			t.Errorf("button markdown missing %q", want)
		}
	}

	// HTML-only component still renders its example from the raw HTML field.
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/gallery/docs/badge.md", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("badge markdown status = %d, want %d", rec.Code, http.StatusOK)
	}
	if md := rec.Body.String(); !strings.Contains(md, `<span class="badge">Badge</span>`) {
		t.Errorf("badge markdown missing raw HTML example: %q", md)
	}

	// Unknown slug returns 404.
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/gallery/docs/nope.md", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("unknown slug status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
