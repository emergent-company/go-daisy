package galleryruntime

import (
	"github.com/emergent-company/go-daisy/render"
	"github.com/labstack/echo/v4"
)

// handleThemeStudio renders the DaisyUI theme generator.
//
// The studio is intentionally client-side: it reads the active theme's custom
// properties, edits them live on a scoped [data-theme="gd-preview"] board, and
// persists the result under CustomThemeStorageKey. Applying a theme only writes
// localStorage — nothing is written to the repository.
func (h *galleryHandler) handleThemeStudio(c echo.Context) error {
	all := h.components
	categories := BuildCategoryGroups(all)

	content := ThemeStudio()
	render.RenderAuto(c.Response().Writer, c.Request(),
		GalleryPage(h.title, "theme-studio", categories, h.logo, content),
		GalleryPageContent(h.title, "theme-studio", categories, h.logo, content),
	)
	return nil
}
