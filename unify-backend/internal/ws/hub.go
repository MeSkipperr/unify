package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients map[*websocket.Conn]bool
	lock    sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (h *Hub) Register(conn *websocket.Conn) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.clients[conn] = true
}

func (h *Hub) Unregister(conn *websocket.Conn) {
	h.lock.Lock()
	defer h.lock.Unlock()
	delete(h.clients, conn)
	conn.Close()
}

func (h *Hub) Broadcast(data interface{}) {
	h.lock.Lock()
	defer h.lock.Unlock()

	for c := range h.clients {
		if err := c.WriteJSON(data); err != nil {
			c.Close()
			delete(h.clients, c)
		}
	}
}
