package handler

import (
	"github.com/a-h/templ"
	"github.com/emergent-company/go-daisy/cmd/nexus/internal/shell"
	"github.com/emergent-company/go-daisy/components/nav"
	"github.com/emergent-company/go-daisy/render"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterLandingRoutes(g *echo.Group) {
	g.GET("", h.handleLanding)
	g.GET("/landing", h.handleLanding)
}

func (h *Handler) RegisterAuthRoutes(g *echo.Group) {
	g.GET("/login", h.handleLogin)
	g.GET("/register", h.handleRegister)
	g.GET("/forgot-password", h.handleForgotPassword)
	g.GET("/reset-password", h.handleResetPassword)
}

func (h *Handler) RegisterDashboardRoutes(g *echo.Group) {
	g.GET("/dashboards/ecommerce", h.handleEcommerceDashboard)
	g.GET("/dashboards/crm", h.handleCRMDashboard)
}

func (h *Handler) RegisterEcommerceRoutes(g *echo.Group) {
	g.GET("/apps/ecommerce/products", h.handleProductsList)
	g.GET("/apps/ecommerce/products/create", h.handleProductCreate)
	g.GET("/apps/ecommerce/products/:id", h.handleProductEdit)
	g.GET("/apps/ecommerce/orders", h.handleOrdersList)
	g.GET("/apps/ecommerce/orders/:id", h.handleOrderDetail)
	g.GET("/apps/ecommerce/sellers", h.handleSellersList)
	g.GET("/apps/ecommerce/sellers/create", h.handleSellerCreate)
	g.GET("/apps/ecommerce/sellers/:id", h.handleSellerEdit)
	g.GET("/apps/ecommerce/customers", h.handleCustomersList)
	g.GET("/apps/ecommerce/customers/create", h.handleCustomerCreate)
	g.GET("/apps/ecommerce/customers/:id", h.handleCustomerEdit)
	g.GET("/apps/ecommerce/shops", h.handleShopsList)
	g.GET("/apps/ecommerce/shops/create", h.handleShopCreate)
	g.GET("/apps/ecommerce/shops/:id", h.handleShopEdit)
}

func (h *Handler) RegisterGenAIRoutes(g *echo.Group) {
	g.GET("/apps/gen-ai/home", h.handleGenAIHome)
	g.GET("/apps/gen-ai/content", h.handleGenAIContent)
	g.GET("/apps/gen-ai/image", h.handleGenAIImage)
	g.GET("/apps/gen-ai/library", h.handleGenAILibrary)
}

func (h *Handler) RegisterAppRoutes(g *echo.Group) {
	g.GET("/apps/file-manager", h.handleFileManager)
	g.GET("/apps/chat", h.handleChat)
}

func (h *Handler) RegisterPagesRoutes(g *echo.Group) {
	g.GET("/pages/settings", h.handleSettings)
	g.GET("/pages/get-help", h.handleGetHelp)
	g.GET("/ui/components", h.handleUIComponents)
	g.GET("/ui/forms", h.handleUIForms)
	g.GET("/ui/charts", h.handleUICharts)
	g.GET("/ui/components/:name", h.handleUIComponentDetail)
	g.GET("/ui/forms/:name", h.handleUIFormDetail)
	g.GET("/ui/charts/:name", h.handleUIChartDetail)
}

func (h *Handler) handleUIComponents(c echo.Context) error {
	return adminRender(c, UIComponentsPage())
}

func (h *Handler) handleUIForms(c echo.Context) error {
	return adminRender(c, UIFormsPage())
}

func (h *Handler) handleUICharts(c echo.Context) error {
	return adminRender(c, UIChartsPage())
}

func (h *Handler) handleUIComponentDetail(c echo.Context) error {
	name := c.Param("name")
	// Pages ported to library components with a faithful card grid.
	switch name {
	case "accordion", "alert", "avatar", "badge", "breadcrumb", "button",
		"countdown", "drawer", "dropdown", "indicator", "loading", "menu",
		"modal", "pagination", "progress", "step", "tab", "table", "timeline",
		"toast", "tooltip":
		return adminRender(c, UIComponentDetailPage(name))
	}
	return adminRender(c, nxHTMLComponent("ui-components-"+name))
}

func (h *Handler) handleUIFormDetail(c echo.Context) error {
	name := c.Param("name")
	switch name {
	case "checkbox", "fieldset", "file-input", "input", "label", "radio",
		"range", "rating", "select", "textarea", "toggle", "validator":
		return adminRender(c, UIFormDetailPage(name))
	}
	return adminRender(c, nxHTMLComponent("ui-forms-"+name))
}

func (h *Handler) handleUIChartDetail(c echo.Context) error {
	name := c.Param("name")
	// Apex chart showcase pages ported to the library ui.Chart component.
	switch name {
	case "area", "bar", "column", "line", "pie":
		return adminRender(c, UIChartDetailPage(name))
	}
	return adminRender(c, nxHTMLComponent("ui-charts-apex-"+name))
}

func adminRender(c echo.Context, content templ.Component) error {
	w := c.Response().Writer
	r := c.Request()
	render.RenderAuto(w, r,
		shell.NexusShell(content, r.URL.Path),
		content,
	)
	return nil
}

func (h *Handler) handleLanding(c echo.Context) error {
	render.RenderAuto(c.Response().Writer, c.Request(),
		LandingPageFull(),
		LandingPage(),
	)
	return nil
}

func (h *Handler) handleLogin(c echo.Context) error {
	render.RenderAuto(c.Response().Writer, c.Request(),
		LoginPageFull(),
		LoginPage(),
	)
	return nil
}

func (h *Handler) handleRegister(c echo.Context) error {
	render.RenderAuto(c.Response().Writer, c.Request(),
		RegisterPageFull(),
		RegisterPage(),
	)
	return nil
}

func (h *Handler) handleForgotPassword(c echo.Context) error {
	render.RenderAuto(c.Response().Writer, c.Request(),
		ForgotPasswordPageFull(),
		ForgotPasswordPage(),
	)
	return nil
}

func (h *Handler) handleResetPassword(c echo.Context) error {
	render.RenderAuto(c.Response().Writer, c.Request(),
		ResetPasswordPageFull(),
		ResetPasswordPage(),
	)
	return nil
}

func (h *Handler) handleEcommerceDashboard(c echo.Context) error {
	return adminRender(c, EcommerceDashboard())
}

func (h *Handler) handleCRMDashboard(c echo.Context) error {
	return adminRender(c, CRMDashboard())
}

func (h *Handler) handleProductsList(c echo.Context) error {
	return adminRender(c, ProductsList())
}

func (h *Handler) handleProductCreate(c echo.Context) error {
	return adminRender(c, ProductForm(true))
}

func (h *Handler) handleProductEdit(c echo.Context) error {
	return adminRender(c, ProductForm(false))
}

func (h *Handler) handleOrdersList(c echo.Context) error {
	return adminRender(c, OrdersList())
}

func (h *Handler) handleOrderDetail(c echo.Context) error {
	return adminRender(c, OrderDetail())
}

func (h *Handler) handleSellersList(c echo.Context) error {
	return adminRender(c, SellersList())
}

func (h *Handler) handleSellerCreate(c echo.Context) error {
	return adminRender(c, SellerForm(true))
}

func (h *Handler) handleSellerEdit(c echo.Context) error {
	return adminRender(c, SellerForm(false))
}

func (h *Handler) handleCustomersList(c echo.Context) error {
	return adminRender(c, CustomersList())
}

func (h *Handler) handleCustomerCreate(c echo.Context) error {
	return adminRender(c, CustomerForm(true))
}

func (h *Handler) handleCustomerEdit(c echo.Context) error {
	return adminRender(c, CustomerForm(false))
}

func (h *Handler) handleShopsList(c echo.Context) error {
	return adminRender(c, ShopsList())
}

func (h *Handler) handleShopCreate(c echo.Context) error {
	return adminRender(c, ShopForm(true))
}

func (h *Handler) handleShopEdit(c echo.Context) error {
	return adminRender(c, ShopForm(false))
}

func (h *Handler) handleGenAIHome(c echo.Context) error {
	return adminRender(c, GenAIHome())
}

func (h *Handler) handleGenAIContent(c echo.Context) error {
	return adminRender(c, GenAIContent())
}

func (h *Handler) handleGenAIImage(c echo.Context) error {
	return adminRender(c, GenAIImage())
}

func (h *Handler) handleGenAILibrary(c echo.Context) error {
	return adminRender(c, GenAILibrary())
}

func (h *Handler) handleFileManager(c echo.Context) error {
	return adminRender(c, FileManager())
}

func (h *Handler) handleChat(c echo.Context) error {
	return adminRender(c, ChatPage())
}

func (h *Handler) handleSettings(c echo.Context) error {
	return adminRender(c, SettingsPage())
}

func (h *Handler) handleGetHelp(c echo.Context) error {
	return adminRender(c, GetHelpPage())
}

var _ = nav.Crumbs
