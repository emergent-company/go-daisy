package gallery

import (
	"context"
	"strings"
	"testing"
)

// AllComponents builds the gallery's static seed. Several entries render a
// component into a strings.Builder to embed its markup in a raw-HTML variant,
// and a mistake there panics at request time rather than failing a build — the
// gallery simply stops serving. Exercise every entry so that class of bug fails
// here instead.
func TestAllComponentsBuilds(t *testing.T) {
	comps := AllComponents()
	if len(comps) == 0 {
		t.Fatal("AllComponents returned an empty seed")
	}
	for i, c := range comps {
		if c.Name == "" {
			t.Errorf("component %d (%q) has no Name", i, c.Slug)
		}
		if c.Slug == "" {
			t.Errorf("component %d (%q) has no Slug", i, c.Name)
		}
	}
}

// TestSearchModalVariantsRender guards the specific defect that broke the
// gallery: the two SearchModal stories dropped their .Render() call, so the
// builder stayed empty and the non-nil templ.Component was passed to panic().
func TestSearchModalVariantsRender(t *testing.T) {
	var found bool
	for _, c := range AllComponents() {
		if c.Slug != "search-modal" {
			continue
		}
		found = true
		if len(c.Variants) == 0 {
			t.Fatal("search-modal has no variants")
		}
		for _, v := range c.Variants {
			if v.Templ == nil {
				t.Errorf("variant %q has no component", v.Name)
				continue
			}
			var buf strings.Builder
			if err := v.Templ.Render(context.Background(), &buf); err != nil {
				t.Fatalf("variant %q failed to render: %v", v.Name, err)
			}
			if !strings.Contains(buf.String(), `id="search-modal"`) {
				t.Errorf("variant %q rendered without the dialog (got %d bytes)", v.Name, buf.Len())
			}
		}
	}
	if !found {
		t.Fatal("no search-modal component in the seed")
	}
}
