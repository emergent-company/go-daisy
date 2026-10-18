package galleryruntime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

// newAPITestServer builds an Echo instance with the production-contract routes
// (healthz, version, api/components) plus the rest of the gallery routes.
func newAPITestServer(components []GalleryComponent) *echo.Echo {
	e := echo.New()
	newGalleryHandler("Test Gallery", nil, components, nil, false).register(e)
	return e
}

func TestHealthz(t *testing.T) {
	e := newAPITestServer(nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json", ct)
	}
	var got map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal body %q: %v", rec.Body.String(), err)
	}
	if got["status"] != "ok" {
		t.Errorf("status = %q, want %q", got["status"], "ok")
	}
}

func TestVersion(t *testing.T) {
	e := newAPITestServer(nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/version", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var got struct {
		Version string `json:"version"`
		Commit  string `json:"commit"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal body %q: %v", rec.Body.String(), err)
	}
	if got.Version != Version {
		t.Errorf("version = %q, want %q", got.Version, Version)
	}
	if got.Commit != Commit {
		t.Errorf("commit = %q, want %q", got.Commit, Commit)
	}
}

type componentsResponse struct {
	Count      int                `json:"count"`
	Components []ComponentSummary `json:"components"`
}

func TestComponentsEmptyRegistry(t *testing.T) {
	e := newAPITestServer(nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/components", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var got componentsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal body %q: %v", rec.Body.String(), err)
	}
	if got.Count != 0 {
		t.Errorf("count = %d, want 0", got.Count)
	}
	if got.Components == nil || len(got.Components) != 0 {
		t.Errorf("components = %#v, want non-nil empty slice", got.Components)
	}
}

func TestComponentsSummary(t *testing.T) {
	components := []GalleryComponent{
		{
			Slug:        "button",
			Name:        "Button",
			Category:    CategoryComponents,
			Subcategory: "Buttons",
			Description: "A clickable button.",
			HidePreview: true,
			Tokens:      []DesignToken{{Label: "Root1"}, {Label: "Root2"}},
			Variants: []GalleryStory{
				{Name: "Default", Tokens: []DesignToken{{Label: "A"}}},
				{Name: "Loading", Tokens: []DesignToken{{Label: "B"}, {Label: "C"}}},
			},
		},
		{
			Slug:     "badge",
			Name:     "Badge",
			Category: CategoryComponents,
			HTML:     `<span class="badge">Badge</span>`,
		},
	}

	e := newAPITestServer(components)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/components", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var got componentsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal body %q: %v", rec.Body.String(), err)
	}

	if got.Count != 2 {
		t.Fatalf("count = %d, want 2", got.Count)
	}
	if len(got.Components) != 2 {
		t.Fatalf("len(components) = %d, want 2", len(got.Components))
	}

	button := got.Components[0]
	if button.Slug != "button" {
		t.Errorf("button.slug = %q, want %q", button.Slug, "button")
	}
	if button.Name != "Button" {
		t.Errorf("button.name = %q, want %q", button.Name, "Button")
	}
	if button.Category != CategoryComponents {
		t.Errorf("button.category = %q, want %q", button.Category, CategoryComponents)
	}
	if button.Subcategory != "Buttons" {
		t.Errorf("button.subcategory = %q, want %q", button.Subcategory, "Buttons")
	}
	if button.Description != "A clickable button." {
		t.Errorf("button.description = %q", button.Description)
	}
	if !button.HidePreview {
		t.Errorf("button.hide_preview = false, want true")
	}
	if !reflect.DeepEqual(button.Variants, []string{"Default", "Loading"}) {
		t.Errorf("button.variants = %v, want [Default Loading]", button.Variants)
	}
	if button.VariantCount != 2 {
		t.Errorf("button.variant_count = %d, want 2", button.VariantCount)
	}
	if button.TokenCount != 5 {
		t.Errorf("button.token_count = %d, want 5 (2 root + 1 + 2)", button.TokenCount)
	}

	badge := got.Components[1]
	if badge.HidePreview {
		t.Errorf("badge.hide_preview = true, want false")
	}
	if !reflect.DeepEqual(badge.Variants, []string{"Default"}) {
		t.Errorf("badge.variants = %v, want [Default] (synthesised)", badge.Variants)
	}
	if badge.VariantCount != 1 {
		t.Errorf("badge.variant_count = %d, want 1", badge.VariantCount)
	}
	if badge.TokenCount != 0 {
		t.Errorf("badge.token_count = %d, want 0", badge.TokenCount)
	}
}

func TestComponentSummaryVariantsAndTokens(t *testing.T) {
	// Explicit variants: names and counts derive from EffectiveVariants.
	c := GalleryComponent{
		Tokens: []DesignToken{{}, {}, {}},
		Variants: []GalleryStory{
			{Name: "Default", Tokens: []DesignToken{{}, {}}},
			{Name: "Error", Tokens: []DesignToken{{}}},
		},
	}
	s := componentSummary(c)
	if s.VariantCount != 2 {
		t.Errorf("variant_count = %d, want 2", s.VariantCount)
	}
	if !reflect.DeepEqual(s.Variants, []string{"Default", "Error"}) {
		t.Errorf("variants = %v, want [Default Error]", s.Variants)
	}
	if s.TokenCount != 6 {
		t.Errorf("token_count = %d, want 6 (3 + 2 + 1)", s.TokenCount)
	}

	// No explicit variants: a single "Default" story is synthesised.
	c2 := GalleryComponent{Tokens: []DesignToken{{}}}
	s2 := componentSummary(c2)
	if s2.VariantCount != 1 {
		t.Errorf("synthesised variant_count = %d, want 1", s2.VariantCount)
	}
	if !reflect.DeepEqual(s2.Variants, []string{"Default"}) {
		t.Errorf("synthesised variants = %v, want [Default]", s2.Variants)
	}
	if s2.TokenCount != 1 {
		t.Errorf("synthesised token_count = %d, want 1", s2.TokenCount)
	}
}
