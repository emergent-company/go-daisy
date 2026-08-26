package streamhub

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// EchoHandler returns an Echo handler function for SSE connections.
// Usage:
//
//	hub := streamhub.New()
//	e.GET("/events", streamhub.EchoHandler(hub))
func EchoHandler(hub *Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
		w := c.Response()
		r := c.Request()

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		hub.ServeHTTP(w, r)
		return nil
	}
}
