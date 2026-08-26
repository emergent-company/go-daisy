// Package html provides server-side streaming helpers for progressive rendering.
// Send the page skeleton immediately, then stream slow content chunks as they
// become available. Works with HTMX for partial swaps and full page loads.
//
// Usage:
//
//	w := streamhtml.NewWriter(c.Response())
//	w.Write(header)       // sent immediately
//	w.Flush()
//	w.Write(slowContent)  // streamed later
//	w.Flush()
//	w.Close()
package html

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/a-h/templ"
)

// Writer wraps an http.ResponseWriter for chunked streaming.
type Writer struct {
	w       io.Writer
	flusher http.Flusher
	mu      sync.Mutex
	closed  bool
}

// NewWriter creates a streaming writer. Caller must set Content-Type and
// flush headers before first Write.
func NewWriter(w http.ResponseWriter) (*Writer, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("stream: response writer does not support flushing")
	}
	return &Writer{w: w, flusher: flusher}, nil
}

// Write writes bytes to the stream.
func (sw *Writer) Write(b []byte) (int, error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.w.Write(b)
}

// WriteString writes a string to the stream.
func (sw *Writer) WriteString(s string) (int, error) {
	return sw.Write([]byte(s))
}

// Render renders a Templ component to the stream.
func (sw *Writer) Render(comp templ.Component) error {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return comp.Render(context.Background(), sw.w)
}

// Flush sends buffered data to the client.
func (sw *Writer) Flush() {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.flusher.Flush()
}

// Close marks the writer as done.
func (sw *Writer) Close() {
	sw.mu.Lock()
	sw.closed = true
	sw.mu.Unlock()
}

// RenderFlush renders a component and flushes immediately.
func (sw *Writer) RenderFlush(comp templ.Component) error {
	if err := sw.Render(comp); err != nil {
		return err
	}
	sw.Flush()
	return nil
}
