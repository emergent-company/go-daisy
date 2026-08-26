// Example: go-daisy landing page with Hero, Features, and CTA.
//
// Build: templ generate && go run ./examples/landing
// Open:  http://localhost:11001
package main

import (
	"net/http"

	"github.com/emergent-company/go-daisy/render"
	"github.com/emergent-company/go-daisy/staticfs"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.HideBanner = true

	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", staticfs.Handler("/"))))

	e.GET("/", func(c echo.Context) error {
		render.RenderPage(c.Response().Writer, c.Request(), LandingPage())
		return nil
	})

	e.Logger.Fatal(e.Start(":11001"))
}
