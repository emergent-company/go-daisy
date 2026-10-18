package galleryruntime

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ComponentSummary is the JSON representation of a single component returned by
// the /api/components endpoint. It deliberately contains no rendered HTML and
// no function values, so it marshals deterministically and is safe for
// automation (deployment health checks, the panel, MCP) to consume.
type ComponentSummary struct {
	Slug         string   `json:"slug"`
	Name         string   `json:"name"`
	Category     Category `json:"category"`
	Subcategory  string   `json:"subcategory"`
	Description  string   `json:"description"`
	HidePreview  bool     `json:"hide_preview"`
	Variants     []string `json:"variants"`
	VariantCount int      `json:"variant_count"`
	TokenCount   int      `json:"token_count"`
}

// componentSummary derives the API summary for a single component from its
// registry entry. Variant names and count come from EffectiveVariants, so a
// component with no explicit Variants reports a single synthesised "Default"
// story (variant_count 1). TokenCount is the total across the component and
// its variants.
func componentSummary(c GalleryComponent) ComponentSummary {
	variants := c.EffectiveVariants()
	names := make([]string, 0, len(variants))
	for _, v := range variants {
		names = append(names, v.Name)
	}
	return ComponentSummary{
		Slug:         c.Slug,
		Name:         c.Name,
		Category:     c.Category,
		Subcategory:  c.Subcategory,
		Description:  c.Description,
		HidePreview:  c.HidePreview,
		Variants:     names,
		VariantCount: len(variants),
		TokenCount:   tokenCount(c),
	}
}

// tokenCount returns the total number of design tokens across a component and
// its effective variants.
func tokenCount(c GalleryComponent) int {
	n := len(c.Tokens)
	for _, v := range c.EffectiveVariants() {
		n += len(v.Tokens)
	}
	return n
}

// handleHealthz reports liveness with no dependencies and no auth.
func (h *galleryHandler) handleHealthz(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// handleVersion reports the build version and commit.
func (h *galleryHandler) handleVersion(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"version": Version, "commit": Commit})
}

// handleComponents lists every registered component as a JSON summary, in
// registry order.
func (h *galleryHandler) handleComponents(c echo.Context) error {
	summaries := make([]ComponentSummary, 0, len(h.components))
	for _, comp := range h.components {
		summaries = append(summaries, componentSummary(comp))
	}
	return c.JSON(http.StatusOK, map[string]any{
		"count":      len(summaries),
		"components": summaries,
	})
}
