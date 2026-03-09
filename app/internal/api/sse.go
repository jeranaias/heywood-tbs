package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"heywood-tbs/internal/middleware"
)

// SSEEvent represents a server-sent event payload.
type SSEEvent struct {
	Type string      `json:"type"` // "notification", "task", "at-risk-alert"
	Data interface{} `json:"data"`
}

// SSEBroker manages SSE client connections grouped by role.
type SSEBroker struct {
	clients map[string]map[chan SSEEvent]bool // role -> set of channels
	mu      sync.RWMutex
}

// NewSSEBroker creates an SSE broker.
func NewSSEBroker() *SSEBroker {
	return &SSEBroker{
		clients: make(map[string]map[chan SSEEvent]bool),
	}
}

// Register adds a client channel for a given role.
func (b *SSEBroker) Register(role string, ch chan SSEEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.clients[role] == nil {
		b.clients[role] = make(map[chan SSEEvent]bool)
	}
	b.clients[role][ch] = true
}

// Unregister removes a client channel.
func (b *SSEBroker) Unregister(role string, ch chan SSEEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.clients[role] != nil {
		delete(b.clients[role], ch)
		if len(b.clients[role]) == 0 {
			delete(b.clients, role)
		}
	}
}

// Broadcast sends an event to all clients for the given role.
// Use role="" to broadcast to all connected clients.
func (b *SSEBroker) Broadcast(role string, event SSEEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if role == "" {
		// Broadcast to all roles
		for _, chs := range b.clients {
			for ch := range chs {
				select {
				case ch <- event:
				default:
					// Drop if client can't keep up
				}
			}
		}
		return
	}

	if chs, ok := b.clients[role]; ok {
		for ch := range chs {
			select {
			case ch <- event:
			default:
			}
		}
	}
}

// ClientCount returns the total number of connected SSE clients.
func (b *SSEBroker) ClientCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	count := 0
	for _, chs := range b.clients {
		count += len(chs)
	}
	return count
}

// handleSSE is the HTTP handler for the SSE stream endpoint.
func (h *Handler) handleSSE(w http.ResponseWriter, r *http.Request) {
	if h.sseBroker == nil {
		writeError(w, 503, "SSE not available")
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, 500, "streaming not supported")
		return
	}

	role := middleware.GetRole(r.Context())

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	ch := make(chan SSEEvent, 16)
	h.sseBroker.Register(role, ch)
	defer h.sseBroker.Unregister(role, ch)

	// Send initial connection event
	fmt.Fprintf(w, "data: {\"type\":\"connected\",\"data\":{\"role\":\"%s\"}}\n\n", role)
	flusher.Flush()

	slog.Info("SSE client connected", "role", role)

	for {
		select {
		case <-r.Context().Done():
			slog.Info("SSE client disconnected", "role", role)
			return
		case event := <-ch:
			data, err := json.Marshal(event)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}
