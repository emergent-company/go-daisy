package galleryruntime

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	"github.com/emergent-company/go-daisy/devmode"
	"github.com/emergent-company/go-daisy/render"
)

// galleryHandler holds gallery route handlers and dependencies.
type galleryHandler struct {
	title          string
	logo           templ.Component // optional branded logo; nil falls back to title text
	components     []GalleryComponent
	staticPrefixes []string // all CSS prefixes: default "/static/" + any ExtraStaticPrefixes
	devMode        bool     // when true, component boundary annotations are injected
}

// newGalleryHandler creates a new gallery handler.
func newGalleryHandler(title string, logo templ.Component, components []GalleryComponent, staticPrefixes []string, devModeEnabled bool) *galleryHandler {
	return &galleryHandler{
		title:          title,
		logo:           logo,
		components:     components,
		staticPrefixes: staticPrefixes,
		devMode:        devModeEnabled,
	}
}

// register mounts all gallery routes on the Echo instance.
func (h *galleryHandler) register(e *echo.Echo) {
	e.GET("/gallery", h.handleIndex)
	e.GET("/gallery/render/:slug", h.handleRender)
	e.GET("/gallery/render/:slug/examples", h.handleRenderSubExample)
	e.GET("/gallery/render/:slug/:variant", h.handleRenderVariant)
	e.GET("/gallery/group/:category/:subcategory", h.handleSubGroup)
	e.GET("/gallery/group/:category", h.handleGroup)
	e.GET("/gallery/:slug", h.handleDetail)
}

// handleIndex renders the gallery landing page.
func (h *galleryHandler) handleIndex(c echo.Context) error {
	all := h.components
	categories := BuildCategoryGroups(all)
	content := GalleryIndex()
	render.RenderAuto(c.Response().Writer, c.Request(),
		GalleryPage(h.title, "", categories, h.logo, content),
		GalleryPageContent(h.title, "", categories, h.logo, content),
	)
	return nil
}

// handleDetail renders the detail page for a single component.
func (h *galleryHandler) handleDetail(c echo.Context) error {
	slug := c.Param("slug")
	comp, ok := ComponentBySlug(h.components, slug)
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "component not found")
	}

	all := h.components
	categories := BuildCategoryGroups(all)

	content := ComponentDetail(comp)
	render.RenderAuto(c.Response().Writer, c.Request(),
		GalleryPage(h.title, slug, categories, h.logo, content),
		GalleryPageContent(h.title, slug, categories, h.logo, content),
	)
	return nil
}

// handleGroup renders the category group summary page.
func (h *galleryHandler) handleGroup(c echo.Context) error {
	slug := c.Param("category")
	cat, ok := CategoryBySlug(slug)
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "category not found")
	}

	all := h.components
	categories := BuildCategoryGroups(all)
	var subcategories []SubcategoryGroup
	for _, cg := range categories {
		if cg.Name == cat {
			subcategories = cg.Subcategories
			break
		}
	}
	if subcategories == nil {
		subcategories = []SubcategoryGroup{}
	}
	content := GalleryGroup(cat, subcategories)

	render.RenderAuto(c.Response().Writer, c.Request(),
		GalleryPage(h.title, "", categories, h.logo, content),
		GalleryPageContent(h.title, "", categories, h.logo, content),
	)
	return nil
}

// handleSubGroup renders the subcategory summary page.
func (h *galleryHandler) handleSubGroup(c echo.Context) error {
	catSlug := c.Param("category")
	subSlug := c.Param("subcategory")
	cat, ok := CategoryBySlug(catSlug)
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "category not found")
	}

	all := h.components
	categories := BuildCategoryGroups(all)
	var found *SubcategoryGroup
	for _, cg := range categories {
		if cg.Name != cat {
			continue
		}
		for i := range cg.Subcategories {
			if Slugify(cg.Subcategories[i].Name) == subSlug {
				found = &cg.Subcategories[i]
				break
			}
		}
		if found != nil {
			break
		}
	}
	if found == nil {
		return echo.NewHTTPError(http.StatusNotFound, "subcategory not found")
	}

	content := GallerySubGroup(cat, *found)
	render.RenderAuto(c.Response().Writer, c.Request(),
		GalleryPage(h.title, "", categories, h.logo, content),
		GalleryPageContent(h.title, "", categories, h.logo, content),
	)
	return nil
}

// handleRender renders a component as a standalone HTML page for use as an iframe src.
// Supports both Templ-based components and HTML-snippet components.
func (h *galleryHandler) handleRender(c echo.Context) error {
	slug := c.Param("slug")
	comp, ok := ComponentBySlug(h.components, slug)
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "component not found")
	}

	c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")

	baseURL := h.baseURL(c)

	if comp.Templ != nil {
		return h.renderTemplPage(c, baseURL, comp.Templ)
	}

	if comp.HTML != "" {
		html := renderSnippetPage(baseURL, h.staticPrefixes, comp.HTML, false)
		_, err := c.Response().Writer.Write([]byte(html))
		return err
	}

	// Fall back to first variant's RenderFunc/Templ/HTML
	if variants := comp.EffectiveVariants(); len(variants) > 0 {
		v := variants[0]
		if v.RenderFunc != nil {
			return h.renderTemplPage(c, baseURL, v.RenderFunc(c.Request().URL.Query()))
		}
		if v.Templ != nil {
			return h.renderTemplPage(c, baseURL, v.Templ)
		}
		if v.HTML != "" {
			html := renderSnippetPage(baseURL, h.staticPrefixes, v.HTML, false)
			_, err := c.Response().Writer.Write([]byte(html))
			return err
		}
	}

	return echo.NewHTTPError(http.StatusNotFound, "component has no renderable content")
}

// handleRenderSubExample renders an individual sub-example (from GallerySubExample)
// as a standalone iframe page. The story index and sub-example index are passed
// via query params: ?s=<storyIdx>&e=<subExampleIdx>.
func (h *galleryHandler) handleRenderSubExample(c echo.Context) error {
	slug := c.Param("slug")
	comp, ok := ComponentBySlug(h.components, slug)
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "component not found")
	}

	si, siErr := strconv.Atoi(c.QueryParam("s"))
	ei, eiErr := strconv.Atoi(c.QueryParam("e"))

	variants := comp.EffectiveVariants()

	// When query params are missing or invalid, fall back to the first
	// non-Interactive story that has SubExamples.
	needsFallback := siErr != nil || eiErr != nil || si < 0 || si >= len(variants)
	if needsFallback {
		found := false
		for i, v := range variants {
			if v.Name != "Interactive" && len(v.SubExamples) > 0 {
				si, ei = i, 0
				found = true
				break
			}
		}
		if !found {
			return echo.NewHTTPError(http.StatusNotFound, "no sub-examples available for this component")
		}
	}

	story := variants[si]
	if ei < 0 || ei >= len(story.SubExamples) {
		// Fall back to first sub-example if index is out of range.
		if len(story.SubExamples) > 0 {
			ei = 0
		} else {
			return echo.NewHTTPError(http.StatusNotFound, "no sub-examples available for this story")
		}
	}
	sub := story.SubExamples[ei]

	c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	return h.renderTemplPage(c, h.baseURL(c), sub.RenderFunc(c.Request().URL.Query()))
}

// handleRenderVariant renders a specific named variant/story of a component.
func (h *galleryHandler) handleRenderVariant(c echo.Context) error {
	slug := c.Param("slug")
	variantSlug := c.Param("variant")
	comp, ok := ComponentBySlug(h.components, slug)
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "component not found")
	}

	story := StoryByName(comp, variantSlug)

	c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")

	baseURL := h.baseURL(c)

	if story.RenderFunc != nil {
		return h.renderTemplPage(c, baseURL, story.RenderFunc(c.Request().URL.Query()))
	}

	if story.Templ != nil {
		return h.renderTemplPage(c, baseURL, story.Templ)
	}

	if story.HTML != "" {
		html := renderSnippetPage(baseURL, h.staticPrefixes, story.HTML, false)
		_, err := c.Response().Writer.Write([]byte(html))
		return err
	}

	return echo.NewHTTPError(http.StatusNotFound, "variant has no renderable content")
}

// baseURL returns the scheme+host for the current request.
func (h *galleryHandler) baseURL(c echo.Context) string {
	scheme := "http"
	if c.Request().TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + c.Request().Host
}

// renderTemplPage renders a templ.Component wrapped in a full HTML shell with
// all project CSS injected, so partial/fragment components display correctly.
func (h *galleryHandler) renderTemplPage(c echo.Context, baseURL string, comp templ.Component) error {
	ctx := c.Request().Context()
	if h.devMode {
		ctx = devmode.WithDevMode(ctx)
	}
	var buf bytes.Buffer
	if err := comp.Render(ctx, &buf); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("render error: %v", err))
	}
	html := renderSnippetPage(baseURL, h.staticPrefixes, buf.String(), h.devMode)
	_, err := c.Response().Writer.Write([]byte(html))
	return err
}

// renderSnippetPage wraps an HTML snippet in a complete standalone HTML document
// with all CSS links injected. When devMode is true, the hover overlay script
// for component boundary visualisation is injected into the document.
func renderSnippetPage(baseURL string, staticPrefixes []string, snippet string, devMode bool) string {
	var cssLinks strings.Builder
	for _, prefix := range staticPrefixes {
		// Ensure prefix ends with /
		p := strings.TrimRight(prefix, "/") + "/"
		fmt.Fprintf(&cssLinks, `  <link href="%s%scss/app.css" rel="stylesheet" type="text/css"/>`, baseURL, p)
		cssLinks.WriteString("\n")
	}

	devScript := ""
	if devMode {
		devScript = devOverlayScript
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en" data-theme="light">
<head>
  <meta charset="UTF-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
%s  <script>
    try {
      var t = localStorage.getItem('gallery-preview-theme');
      if (t) document.documentElement.setAttribute('data-theme', t);
    } catch(e) {}
  </script>
  <style>
    html { margin: 0; padding: 0; background: transparent; }
    body { margin: 0; padding: 16px; background: transparent; }
  </style>
</head>
<body>
%s%s
</body>
</html>`, cssLinks.String(), snippet, devScript)
}


