// Package head provides base <head> dependencies for go-daisy consumer apps.
package head

// DepsProps controls which CSS/JS assets Dependencies() emits.
type DepsProps struct {
	// Core (always emitted): app.css, htmx.js
	Alpine  bool // include alpine.js
	Morph   bool // include morph.js
	Stimulus bool // include stimulus.js + controllers
	SSE     bool // include htmx SSE extension
	WS      bool // include htmx WebSocket extension
}
