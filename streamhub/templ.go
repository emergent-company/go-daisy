package streamhub

import (
	"context"
	"strings"

	"github.com/a-h/templ"
)

// Send sends raw SSE data to a specific connection by ID.
func (h *Hub) Send(connID string, data string) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if conn, ok := h.conns[connID]; ok {
		h.writeToConn(conn, data)
	}
}

// SendComponent renders a Templ component and sends it to a specific connection.
func (h *Hub) SendComponent(connID string, comp templ.Component) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	conn, ok := h.conns[connID]
	if !ok {
		return
	}
	var sb strings.Builder
	if err := comp.Render(context.Background(), &sb); err != nil {
		return
	}
	h.writeToConn(conn, sb.String())
}

// BroadcastComponent renders a Templ component and broadcasts it to all clients.
func (h *Hub) BroadcastComponent(comp templ.Component) {
	var sb strings.Builder
	if err := comp.Render(context.Background(), &sb); err != nil {
		return
	}
	h.Broadcast(sb.String())
}

// BroadcastChannelComponent renders a Templ component and sends it to a channel.
func (h *Hub) BroadcastChannelComponent(channel string, comp templ.Component) {
	var sb strings.Builder
	if err := comp.Render(context.Background(), &sb); err != nil {
		return
	}
	h.BroadcastChannel(channel, sb.String())
}

// BroadcastAppend renders a component and appends it to a CSS selector on all clients.
// Target is the CSS selector (e.g. "#messages").
func (h *Hub) BroadcastAppend(target string, comp templ.Component) {
	var sb strings.Builder
	sb.WriteString(`<div hx-swap-oob="beforeend:`)
	sb.WriteString(target)
	sb.WriteString(`">`)
	_ = comp.Render(context.Background(), &sb)
	sb.WriteString(`</div>`)
	h.Broadcast(sb.String())
}

// BroadcastReplace renders a component and replaces a CSS selector on all clients.
func (h *Hub) BroadcastReplace(target string, comp templ.Component) {
	var sb strings.Builder
	sb.WriteString(`<div hx-swap-oob="outerHTML:`)
	sb.WriteString(target)
	sb.WriteString(`">`)
	_ = comp.Render(context.Background(), &sb)
	sb.WriteString(`</div>`)
	h.Broadcast(sb.String())
}

// BroadcastUpdate renders a component and updates innerHTML of a selector on all clients.
func (h *Hub) BroadcastUpdate(target string, comp templ.Component) {
	var sb strings.Builder
	sb.WriteString(`<div hx-swap-oob="innerHTML:`)
	sb.WriteString(target)
	sb.WriteString(`">`)
	_ = comp.Render(context.Background(), &sb)
	sb.WriteString(`</div>`)
	h.Broadcast(sb.String())
}

// BroadcastRemove removes a CSS selector on all clients.
func (h *Hub) BroadcastRemove(target string) {
	h.Broadcast(`<div hx-swap-oob="delete:` + target + `"></div>`)
}
