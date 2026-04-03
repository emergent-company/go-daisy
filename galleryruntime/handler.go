package galleryruntime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/emergent-company/go-daisy/render"
)

// galleryHandler holds gallery route handlers and dependencies.
type galleryHandler struct {
	title      string
	components []GalleryComponent
	store      *Store
	github     *GitHubClient
}

// newGalleryHandler creates a new gallery handler.
func newGalleryHandler(title string, components []GalleryComponent, store *Store, gh *GitHubClient) *galleryHandler {
	return &galleryHandler{
		title:      title,
		components: components,
		store:      store,
		github:     gh,
	}
}

// register mounts all gallery routes on the Echo instance.
func (h *galleryHandler) register(e *echo.Echo) {
	e.GET("/gallery", h.handleIndex)
	e.GET("/gallery/render/:slug", h.handleRender)
	e.GET("/gallery/render/:slug/:variant", h.handleRenderVariant)
	e.GET("/gallery/:slug", h.handleDetail)

	// Feedback routes
	e.POST("/gallery/:slug/feedback", h.handleCreateFeedback)
	e.GET("/gallery/:slug/feedback", h.handleListFeedback)
	e.GET("/gallery/:slug/feedback/count", h.handleCountFeedback)
	e.DELETE("/gallery/:slug/feedback/:id", h.handleDeleteFeedback)
	e.POST("/gallery/:slug/feedback/export-issue", h.handleExportIssue)
}

// handleIndex renders the gallery landing page.
func (h *galleryHandler) handleIndex(c echo.Context) error {
	all := h.components
	categories := BuildCategoryGroups(all)
	content := GalleryIndex()
	render.RenderAuto(c.Response().Writer, c.Request(),
		GalleryPage(h.title, "", categories, content),
		GalleryPageContent(h.title, "", categories, content),
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

	var feedbackCount int64
	if h.store != nil {
		feedbackCount, _ = h.store.Count(c.Request().Context(), slug)
	}

	content := ComponentDetail(comp, feedbackCount, h.github != nil)
	render.RenderAuto(c.Response().Writer, c.Request(),
		GalleryPage(h.title, slug, categories, content),
		GalleryPageContent(h.title, slug, categories, content),
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

	if comp.Templ != nil {
		return comp.Templ.Render(c.Request().Context(), c.Response().Writer)
	}

	if comp.HTML != "" {
		scheme := "http"
		if c.Request().TLS != nil {
			scheme = "https"
		}
		baseURL := scheme + "://" + c.Request().Host
		html := renderSnippetPage(baseURL, comp.HTML)
		_, err := c.Response().Writer.Write([]byte(html))
		return err
	}

	return echo.NewHTTPError(http.StatusNotFound, "component has no renderable content")
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

	if story.Templ != nil {
		return story.Templ.Render(c.Request().Context(), c.Response().Writer)
	}

	if story.HTML != "" {
		scheme := "http"
		if c.Request().TLS != nil {
			scheme = "https"
		}
		baseURL := scheme + "://" + c.Request().Host
		html := renderSnippetPage(baseURL, story.HTML)
		_, err := c.Response().Writer.Write([]byte(html))
		return err
	}

	return echo.NewHTTPError(http.StatusNotFound, "variant has no renderable content")
}

// renderSnippetPage wraps an HTML snippet in a complete standalone HTML document.
func renderSnippetPage(baseURL, snippet string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en" data-theme="light">
<head>
  <meta charset="UTF-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
  <link href="%s/static/css/app.css" rel="stylesheet" type="text/css"/>
  <script>
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
%s
</body>
</html>`, baseURL, snippet)
}

// feedbackRequest is the JSON body for POST /gallery/:slug/feedback.
type feedbackRequest struct {
	Comment     string          `json:"comment"`
	ContextJSON json.RawMessage `json:"context_json"`
}

// feedbackResponse is the JSON body returned after creating feedback.
type feedbackResponse struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

// handleCreateFeedback handles POST /gallery/:slug/feedback.
func (h *galleryHandler) handleCreateFeedback(c echo.Context) error {
	if h.store == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "gallery database not available"})
	}

	slug := c.Param("slug")

	var req feedbackRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON body")
	}
	if req.Comment == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "comment is required"})
	}

	contextJSON := string(req.ContextJSON)
	if contextJSON == "" {
		contextJSON = "{}"
	}

	record, err := h.store.Create(c.Request().Context(), CreateParams{
		ComponentSlug: slug,
		Comment:       req.Comment,
		ContextJSON:   contextJSON,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to save feedback: %v", err))
	}

	c.Response().WriteHeader(http.StatusCreated)
	return c.JSON(http.StatusCreated, feedbackResponse{
		ID:        record.ID,
		CreatedAt: record.CreatedAt,
	})
}

// handleListFeedback handles GET /gallery/:slug/feedback.
func (h *galleryHandler) handleListFeedback(c echo.Context) error {
	var items []Feedback
	if h.store != nil {
		slug := c.Param("slug")
		var err error
		items, err = h.store.List(c.Request().Context(), slug)
		if err != nil {
			items = []Feedback{}
		}
	}
	if items == nil {
		items = []Feedback{}
	}

	c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	return FeedbackListPartial(items).Render(c.Request().Context(), c.Response().Writer)
}

// handleCountFeedback handles GET /gallery/:slug/feedback/count.
func (h *galleryHandler) handleCountFeedback(c echo.Context) error {
	count := int64(0)
	slug := c.Param("slug")
	if h.store != nil {
		var err error
		count, err = h.store.Count(c.Request().Context(), slug)
		if err != nil {
			count = 0
		}
	}

	c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := fmt.Fprintf(c.Response().Writer,
		`<span id="feedback-count-%s" class="badge badge-xs badge-primary">%d</span>`,
		slug, count)
	return err
}

// handleDeleteFeedback handles DELETE /gallery/:slug/feedback/:id.
func (h *galleryHandler) handleDeleteFeedback(c echo.Context) error {
	if h.store == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "gallery database not available")
	}

	slug := c.Param("slug")
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	if err := h.store.Delete(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "feedback not found")
	}

	count, _ := h.store.Count(c.Request().Context(), slug)
	c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = fmt.Fprintf(c.Response().Writer,
		`<span id="feedback-count-%s" class="badge badge-xs badge-primary">%d</span>`,
		slug, count)
	return err
}

// exportIssueResponse is the JSON body returned after creating a GitHub issue.
type exportIssueResponse struct {
	IssueURL string `json:"issue_url"`
}

// handleExportIssue handles POST /gallery/:slug/feedback/export-issue.
func (h *galleryHandler) handleExportIssue(c echo.Context) error {
	if h.github == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "GitHub integration not configured"})
	}
	if h.store == nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "gallery database not available"})
	}

	slug := c.Param("slug")
	comp, ok := ComponentBySlug(h.components, slug)
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "component not found")
	}

	items, err := h.store.ListOpen(c.Request().Context(), slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to list feedback: %v", err))
	}
	if len(items) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no open feedback items to export"})
	}

	scheme := "http"
	if c.Request().TLS != nil {
		scheme = "https"
	}
	baseURL := scheme + "://" + c.Request().Host

	title, body := BuildIssueContent(comp, items, baseURL)
	issueURL, err := h.github.CreateIssue(c.Request().Context(), title, body, []string{"gallery-feedback"})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to create GitHub issue: %v", err))
	}

	return c.JSON(http.StatusCreated, exportIssueResponse{IssueURL: issueURL})
}
