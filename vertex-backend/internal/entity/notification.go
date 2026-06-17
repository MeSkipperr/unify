package entity

import (
	"time"

	"github.com/google/uuid"
)

// NotificationStatus defines the urgency or category of a notification.
type NotificationStatus string

const (
	NotificationStatusInfo    NotificationStatus = "info"
	NotificationStatusWarning NotificationStatus = "warning"
	NotificationStatusError   NotificationStatus = "error"
	NotificationStatusSuccess NotificationStatus = "success"
)

// Notification stores user notifications with a direct link to related information.
type Notification struct {
	ID        uuid.UUID            `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title     string               `gorm:"size:255;not null" json:"title"`
	Message   string               `gorm:"type:text;not null" json:"message"`
	Status    NotificationStatus   `gorm:"type:varchar(20);not null;default:'info'" json:"status"`
	URL       string               `gorm:"size:512" json:"url"`
	IsRead    bool                 `gorm:"default:false" json:"isRead"`
	CreatedAt time.Time            `gorm:"autoCreateTime;index" json:"createdAt"`
}
