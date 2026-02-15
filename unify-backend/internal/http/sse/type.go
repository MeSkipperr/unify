package sse

import (
	"time"
	"unify-backend/internal/core/mtr"
)

const (
	SSEChannelNotif    = "notification"
	SSEChannelLogs     = "logs"
	SSEChannelServices = "services"
	SSEChannelMTR      = "mtr"
)

type NotificationEvent struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	URL    string `json:"url"`
}

type MtrEvent mtr.MtrResultJson

type ServicesEvent struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type LogsEvent struct {
	CreatedAt   time.Time `json:"createdAt"`
	Level       string    `json:"level"`
	Message     string    `json:"message"`
	ServiceName string    `json:"serviceName"`
}
