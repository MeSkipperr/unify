package notification

import (
	"unify-backend/internal/database"
	"unify-backend/internal/http/sse"
	"unify-backend/internal/worker"
	"unify-backend/models"
)

func SSENotification(data models.Notification) error {
	if err := database.DB.Create(&data).Error; err != nil {
		return err
	}
	sseManager := worker.ManagerGlobal.GetSSE()
	if sseManager == nil {
		return nil
	}
	sseManager.Broadcast(sse.SSEChannelNotif, data)

	return nil
}
