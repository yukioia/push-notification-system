package sse

import (
	"fmt"
	"net/http"
	"sync"
)

type Client struct {
	ClientID string
	Channel  chan string
}

type Manager struct {
	mu      sync.RWMutex
	clients map[string]*Client
}

func NewManager() *Manager {
	return &Manager{clients: make(map[string]*Client)}
}

func (m *Manager) Add(clientID string) *Client {
	m.mu.Lock()
	defer m.mu.Unlock()
	c := &Client{ClientID: clientID, Channel: make(chan string, 16)}
	m.clients[clientID] = c
	return c
}

func (m *Manager) Remove(clientID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if c, ok := m.clients[clientID]; ok {
		close(c.Channel)
		delete(m.clients, clientID)
	}
}

func (m *Manager) SendTo(clientID string, msg string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.clients[clientID]
	if !ok {
		return fmt.Errorf("client not connected")
	}
	select {
	case c.Channel <- msg:
		return nil
	default:
		// channel full -> drop
		return fmt.Errorf("client channel full")
	}
}

func (m *Manager) BroadcastToClients(clientIDs []string, msg string) {
	for _, id := range clientIDs {
		_ = m.SendTo(id, msg)
	}
}

// SSE handler helper
func SSEHandler(w http.ResponseWriter, r *http.Request, client *Client, mngr *Manager) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// keep connection open until client disconnects
	notify := r.Context().Done()
	for {
		select {
		case <-notify:
			mngr.Remove(client.ClientID)
			return
		case msg, ok := <-client.Channel:
			if !ok {
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}
}
