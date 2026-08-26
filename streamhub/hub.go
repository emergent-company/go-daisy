// Package streamhub provides a Go-native pub/sub SSE hub for real-time
// page updates. Register connections, broadcast to channels, handle
// disconnect/reconnect. Works with any http.Handler or Echo framework.
//
// Usage:
//
//	hub := streamhub.New()
//	e.GET("/events", func(c echo.Context) error {
//	    return hub.ServeHTTP(c.Response(), c.Request())
//	})
//	hub.Broadcast("orders", stream.Append("#list", itemRow))
package streamhub

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// Hub manages SSE connections and broadcasts.
type Hub struct {
	mu          sync.RWMutex
	conns       map[string]*sseConn
	logger      *slog.Logger
	heartbeat   time.Duration
	maxBodySize int64
	nextID      int64
}

// sseConn represents a single SSE connection.
type sseConn struct {
	id       string
	w        io.Writer
	flusher  http.Flusher
	channels []string
	done     chan struct{}
	closed   bool
}

// New creates a new SSE hub with defaults.
func New() *Hub {
	return &Hub{
		conns:       make(map[string]*sseConn),
		logger:      slog.Default(),
		heartbeat:   30 * time.Second,
		maxBodySize: 64 * 1024,
	}
}

// ServeHTTP handles an SSE connection upgrade. The caller must set headers
// before calling this:
//
//	w.Header().Set("Content-Type", "text/event-stream")
//	w.Header().Set("Cache-Control", "no-cache")
//	w.Header().Set("Connection", "keep-alive")
//	w.WriteHeader(http.StatusOK)
func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	h.mu.Lock()
	h.nextID++
	id := fmt.Sprintf("sse-%d", h.nextID)
	conn := &sseConn{
		id:      id,
		w:       w,
		flusher: flusher,
		done:    make(chan struct{}),
	}
	h.conns[id] = conn
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.conns, id)
		h.mu.Unlock()
	}()

	// Send initial connection event.
	fmt.Fprintf(w, "data: {\"type\":\"connected\",\"id\":\"%s\"}\n\n", id)
	flusher.Flush()

	// Heartbeat.
	ticker := time.NewTicker(h.heartbeat)
	defer ticker.Stop()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.mu.RLock()
			closed := conn.closed
			h.mu.RUnlock()
			if closed {
				return
			}
			_, err := fmt.Fprintf(w, ": heartbeat\n\n")
			if err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

// Broadcast sends data to all connected clients.
func (h *Hub) Broadcast(data string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, conn := range h.conns {
		h.writeToConn(conn, data)
	}
}

// BroadcastChannel sends data to clients subscribed to a specific channel.
func (h *Hub) BroadcastChannel(channel string, data string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, conn := range h.conns {
		for _, ch := range conn.channels {
			if ch == channel {
				h.writeToConn(conn, data)
				break
			}
		}
	}
}

// Subscribe adds a client to a channel by connection ID.
func (h *Hub) Subscribe(connID string, channels ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conn, ok := h.conns[connID]; ok {
		conn.channels = append(conn.channels, channels...)
	}
}

// Count returns the number of connected clients.
func (h *Hub) Count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.conns)
}

// BroadcastRefresh sends a refresh stream action to all connected clients.
// Browsers will re-fetch the current page and morph the DOM (idiomorph required).
// This is the simplest real-time pattern — one call, no targeted selectors needed.
func (h *Hub) BroadcastRefresh() {
	h.Broadcast(`<div hx-swap-oob="refresh"></div>`)
}

// BroadcastChannelRefresh sends a refresh to clients on a specific channel.
func (h *Hub) BroadcastChannelRefresh(channel string) {
	h.BroadcastChannel(channel, `<div hx-swap-oob="refresh"></div>`)
}

func (h *Hub) writeToConn(conn *sseConn, data string) {
	msg := fmt.Sprintf("data: %s\n\n", data)
	_, err := io.WriteString(conn.w, msg)
	if err != nil {
		conn.closed = true
		return
	}
	conn.flusher.Flush()
}
