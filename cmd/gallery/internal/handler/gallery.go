package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/emergent-company/go-daisy/cmd/gallery/internal/pages"
	"github.com/emergent-company/go-daisy/render"
)

// Handler holds gallery route handlers.
type Handler struct{}

func New() *Handler { return &Handler{} }

func (h *Handler) Register(e *echo.Echo) {
	e.GET("/components", h.Home)
	e.GET("/components/buttons", h.Buttons)
	e.GET("/components/cards", h.Cards)
	e.GET("/components/forms", h.Forms)
	e.GET("/components/tables", h.Tables)
	e.GET("/components/modals", h.Modals)
	e.GET("/components/feedback", h.Feedback)
	e.GET("/components/nav", h.Nav)
	e.GET("/components/layout", h.Layout)
}

func (h *Handler) Home(c echo.Context) error {
	render.RenderPage(c.Response().Writer, c.Request(), pages.GalleryPage("Components", pages.HomeContent()))
	return nil
}

func (h *Handler) Buttons(c echo.Context) error {
	render.RenderAuto(c.Response().Writer, c.Request(),
		pages.GalleryPage("Buttons", pages.ButtonsContent()),
		pages.ButtonsContent())
	return nil
}

func (h *Handler) Cards(c echo.Context) error {
	render.RenderAuto(c.Response().Writer, c.Request(),
		pages.GalleryPage("Cards", pages.CardsContent()),
		pages.CardsContent())
	return nil
}

func (h *Handler) Forms(c echo.Context) error {
	render.RenderAuto(c.Response().Writer, c.Request(),
		pages.GalleryPage("Forms", pages.FormsContent()),
		pages.FormsContent())
	return nil
}

func (h *Handler) Tables(c echo.Context) error {
	render.RenderAuto(c.Response().Writer, c.Request(),
		pages.GalleryPage("Tables", pages.TablesContent()),
		pages.TablesContent())
	return nil
}

func (h *Handler) Modals(c echo.Context) error {
	render.RenderAuto(c.Response().Writer, c.Request(),
		pages.GalleryPage("Modals", pages.ModalsContent()),
		pages.ModalsContent())
	return nil
}

func (h *Handler) Feedback(c echo.Context) error {
	render.RenderAuto(c.Response().Writer, c.Request(),
		pages.GalleryPage("Feedback", pages.FeedbackContent()),
		pages.FeedbackContent())
	return nil
}

func (h *Handler) Nav(c echo.Context) error {
	render.RenderAuto(c.Response().Writer, c.Request(),
		pages.GalleryPage("Nav", pages.NavContent()),
		pages.NavContent())
	return nil
}

func (h *Handler) Layout(c echo.Context) error {
	render.RenderPage(c.Response().Writer, c.Request(), pages.LayoutDemo())
	return nil
}
