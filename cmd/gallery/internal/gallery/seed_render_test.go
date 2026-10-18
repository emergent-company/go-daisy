package gallery

import (
	"bytes"
	"context"
	"net/url"
	"strings"
	"testing"

	"golang.org/x/net/html"

	"github.com/emergent-company/go-daisy/galleryruntime"
)

// TestRegistryRenderSmoke renders every registered gallery story.
//
// This is the closest thing to per-component coverage we get for free: the
// registry is the single list of everything go-daisy ships, so a component that
// panics, returns nil, or emits unparseable markup fails here without a
// bespoke test per component. It is intentionally shallow — it proves the
// component renders, not that its markup is correct. Attribute-level contracts
// live next to the components (see components/ui/render_smoke_test.go).
func TestRegistryRenderSmoke(t *testing.T) {
	components := AllComponents()
	if len(components) == 0 {
		t.Fatal("registry is empty")
	}

	stories := 0
	for _, component := range components {
		component := component
		t.Run(component.Slug, func(t *testing.T) {
			if component.Name == "" {
				t.Error("empty Name")
			}
			if component.Description == "" {
				t.Error("empty Description")
			}
			if component.Category == "" {
				t.Error("empty Category")
			}

			variants := component.EffectiveVariants()
			if len(variants) == 0 {
				t.Fatal("no variants")
			}
			for _, variant := range variants {
				variant := variant
				if variant.SubExamples != nil {
					// Sub-examples render as independent iframes from the detail
					// page, not through this story's own renderer.
					return
				}
				if variant.RenderFunc == nil && variant.Templ == nil && variant.HTML == "" {
					t.Error("story has no RenderFunc, Templ, HTML, or SubExamples")
					return
				}
				stories++
				t.Run(variant.Name, func(t *testing.T) {
					rendered := renderStory(t, component, variant)
					if strings.TrimSpace(rendered) == "" {
						t.Fatal("rendered empty output")
					}
					assertParseable(t, rendered)
				})
			}
		})
	}
	t.Logf("rendered %d stories across %d components", stories, len(components))
	if stories < len(components) {
		t.Errorf("expected at least one story per component: %d stories for %d components", stories, len(components))
	}
}

// TestRegistrySlugsUnique catches two components fighting over one gallery URL
// (the detail route keys off the slug, so a duplicate silently hides one).
func TestRegistrySlugsUnique(t *testing.T) {
	seen := make(map[string]string, len(AllComponents()))
	for _, component := range AllComponents() {
		if previous, ok := seen[component.Slug]; ok {
			t.Errorf("duplicate slug %q: %s and %s", component.Slug, previous, component.Name)
			continue
		}
		seen[component.Slug] = component.Name
	}
}

// TestRegistrySlugsResolvable catches a slug-map entry pointing at a gallery
// page that does not exist, which makes the dev overlay's click-through 404.
func TestRegistrySlugsResolvable(t *testing.T) {
	known := make(map[string]bool, len(AllComponents()))
	for _, component := range AllComponents() {
		known[component.Slug] = true
	}
	for name, slug := range galleryruntime.ComponentSlugs {
		if !known[slug] {
			t.Errorf("ComponentSlugs[%q] = %q: no component registered with that slug", name, slug)
		}
	}
}

func renderStory(t *testing.T, component galleryruntime.GalleryComponent, variant galleryruntime.GalleryStory) (out string) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("%s/%s panicked: %v", component.Slug, variant.Name, r)
		}
	}()

	ctx := context.Background()
	var buf bytes.Buffer
	var err error

	switch {
	case variant.RenderFunc != nil:
		err = variant.RenderFunc(url.Values{}).Render(ctx, &buf)
	case variant.Templ != nil:
		err = variant.Templ.Render(ctx, &buf)
	default:
		// Raw HTML stories are served as-is by the gallery.
		buf.WriteString(variant.HTML)
	}
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// assertParseable parses the fragment as HTML. html.Parse accepts nearly
// anything, so this catches truncated or malformed markup (unclosed quotes,
// stray "<") rather than full validity.
func assertParseable(t *testing.T, markup string) {
	t.Helper()
	if _, err := html.Parse(strings.NewReader(markup)); err != nil {
		t.Fatalf("unparseable html: %v", err)
	}
}
