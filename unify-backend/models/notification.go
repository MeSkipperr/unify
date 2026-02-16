package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationLevel string

const (
	NoticationStatusAlert   NotificationLevel = "notification-alert"
	NoticationStatusInfo    NotificationLevel = "notification-info"
	NoticationStatusWarning NotificationLevel = "notification-warning"
	NoticationStatusError   NotificationLevel = "notification-error"
)

type Notification struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	Level  NotificationLevel `gorm:"type:varchar(20);not null;index" json:"level"`
	Title  string            `gorm:"type:varchar(200);not null" json:"title"`
	Detail string            `gorm:"type:text" json:"detail"`
	URL    string            `gorm:"type:varchar(255)" json:"url"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}
