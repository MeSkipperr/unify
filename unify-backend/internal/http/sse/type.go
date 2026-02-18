package sse

import (
	"time"
	"unify-backend/internal/core/mtr"

	"github.com/google/uuid"
)

const (
	SSEChannelNotif     = "notification"
	SSEChannelLogs      = "logs"
	SSEChannelServices  = "services"
	SSEChannelMTR       = "mtr"
	SSEChannelDashboard = "dashboard"
)

type MtrEvent struct {
	Time    time.Time         `json:"time"`
	ID      uuid.UUID         `json:"id"`
	Message mtr.MtrResultJson `json:"message"`
}

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
