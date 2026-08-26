package main

import (
	"log"
	"net/http"

	"github.com/emergent-company/go-daisy/cmd/nexus/internal/handler"
	"github.com/emergent-company/go-daisy/staticfs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/static/*", echo.WrapHandler(staticfs.Handler("/static/")))
	e.Static("/images", "/root/nexus-html/html/images")
	e.Static("/js", "/root/nexus-html/html/js")

	h := handler.New()

	landing := e.Group("")
	h.RegisterLandingRoutes(landing)

	auth := e.Group("/auth")
	h.RegisterAuthRoutes(auth)

	admin := e.Group("")
	h.RegisterDashboardRoutes(admin)
	h.RegisterEcommerceRoutes(admin)
	h.RegisterGenAIRoutes(admin)
	h.RegisterAppRoutes(admin)
	h.RegisterPagesRoutes(admin)

	e.GET("/*", func(c echo.Context) error {
		return c.HTML(http.StatusNotFound, `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"/><meta name="viewport" content="width=device-width, initial-scale=1.0"/><link rel="stylesheet" href="/static/css/app.css"/></head><body class="bg-base-200 min-h-screen"><div class="flex flex-col items-center justify-center min-h-screen"><p class="text-4xl font-bold">404</p><a href="/landing" class="btn btn-primary mt-4">Go to Home</a></div></body></html>`)
	})

	log.Printf("Nexus test app starting on :11001")
	e.Logger.Fatal(e.Start(":11001"))
}
