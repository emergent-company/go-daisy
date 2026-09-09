// Package head provides base <head> dependencies for go-daisy consumer apps.
package head

// DepsProps controls which CSS/JS assets Dependencies() emits.
type DepsProps struct {
	// Core (always emitted): app.css, htmx.js
	Alpine   bool // include alpine.js
	Morph    bool // include morph.js (idiomorph global; htmx morphing is core in v4)
	Stimulus bool // include stimulus.js + controllers
	SSE      bool // include hx-sse.js extension
	WS       bool // include hx-ws.js extension
}
