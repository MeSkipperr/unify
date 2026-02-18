package sse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"unify-backend/config"
)

type Client struct {
	channel chan string
}

type SSEManager struct {
	mu      sync.RWMutex
	clients map[string]map[*Client]bool // channelName -> clients
}

func NewSSEManager() *SSEManager {
	return &SSEManager{
		clients: make(map[string]map[*Client]bool),
	}
}

// Register client ke channel tertentu
func (m *SSEManager) Subscribe(w http.ResponseWriter, r *http.Request, channelName string) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	origin := r.Header.Get("Origin")

	cfg := config.LoadConfig() 

	allowedMap := make(map[string]bool)
	for _, o := range cfg.AllowedOrigins {
		allowedMap[o] = true
	}

	if allowedMap[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	client := &Client{
		channel: make(chan string),
	}

	// Register client
	m.mu.Lock()
	if m.clients[channelName] == nil {
		m.clients[channelName] = make(map[*Client]bool)
	}
	m.clients[channelName][client] = true
	m.mu.Unlock()

	ctx := r.Context()

	for {
		select {
		case <-ctx.Done():
			m.removeClient(channelName, client)
			return
		case msg := <-client.channel:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}
}

func (m *SSEManager) removeClient(channelName string, client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.clients[channelName], client)
	close(client.channel)
}

// Broadcast ke channel tertentu
func (m *SSEManager) Broadcast(channel string, payload interface{}) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	for client := range m.clients[channel] {
		select {
		case client.channel <- string(data):
		default:
		}
	}
}
