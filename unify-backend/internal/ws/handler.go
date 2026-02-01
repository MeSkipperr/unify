package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
)


// Upgrader exported supaya bisa dipakai di package lain
var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // untuk local / internal network
	},
}


func ServeWS(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := Upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		hub.Register(conn)

		go func() {
			defer hub.Unregister(conn)
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					break
				}
			}
		}()
	}
}
