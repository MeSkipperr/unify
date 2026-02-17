package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Service struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	ServiceName string `gorm:"type:varchar(100);uniqueIndex;not null" json:"serviceName"`
	DisplayName string `gorm:"type:varchar(150);not null" json:"displayName"`
	Description string `gorm:"type:text" json:"description,omitempty"`

	Status  string `gorm:"type:varchar(20);not null;index" json:"status"`
	Version string `gorm:"type:varchar(20)" json:"version,omitempty"`
	Type    string `gorm:"type:varchar(50)" json:"type,omitempty"`

	Config datatypes.JSON `gorm:"type:jsonb" json:"config,omitempty"`

	UpdatedAt time.Time `json:"updatedAt"`
}
