package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/emergent-company/go-daisy/cmd/gallery/internal/handler"
	"github.com/emergent-company/go-daisy/staticfs"
)

func main() {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Static files
	e.GET("/static/*", echo.WrapHandler(staticfs.Handler("/static/")))

	// Redirect root to /components
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/components")
	})

	h := handler.New()
	h.Register(e)

	log.Println("gallery listening on :4100")
	log.Fatal(e.Start(":4100"))
}
