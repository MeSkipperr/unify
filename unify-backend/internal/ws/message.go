package ws

import "time"

type Message struct {
	Time    time.Time   `json:"time"`
	ID      string      `json:"id"`
	Message interface{} `json:"message"`
}
