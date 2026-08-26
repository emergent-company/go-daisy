// Example: go-daisy form with HTMX submit + inline validation.
//
// Build: templ generate && go run ./examples/forms
// Open:  http://localhost:11002
package main

import (
	"net/http"
	"strings"

	"github.com/emergent-company/go-daisy/render"
	"github.com/emergent-company/go-daisy/staticfs"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.HideBanner = true

	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", staticfs.Handler("/"))))

	e.GET("/", func(c echo.Context) error {
		render.RenderPage(c.Response().Writer, c.Request(), FormPage("", ""))
		return nil
	})

	e.POST("/submit", func(c echo.Context) error {
		name := strings.TrimSpace(c.Request().FormValue("name"))
		email := strings.TrimSpace(c.Request().FormValue("email"))
		var errMsg string
		if name == "" {
			errMsg = "Name is required"
		} else if email == "" {
			errMsg = "Email is required"
		}
		render.RenderPartial(c.Response().Writer, c.Request(), FormPage(name, errMsg))
		return nil
	})

	e.Logger.Fatal(e.Start(":11002"))
}
