package ws

import "time"

type Message struct {
	Time    time.Time   `json:"time"`
	ID      int      `json:"id"`
	Message interface{} `json:"message"`
}
