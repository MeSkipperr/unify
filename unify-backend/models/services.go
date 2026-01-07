package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ServiceStatus struct {
	Status string `gorm:"type:varchar(20);primaryKey"`
}

type Service struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	ServiceName string `gorm:"type:varchar(100);uniqueIndex;not null"`
	DisplayName string `gorm:"type:varchar(150);not null"`
	Description string `gorm:"type:text"`

	Status  string `gorm:"type:varchar(20);not null;index"`
	Version string `gorm:"type:varchar(20)"`
	Type    string `gorm:"type:varchar(50)"`

	Config datatypes.JSON `gorm:"type:jsonb"`

	UpdatedAt time.Time

	// Relation
	ServiceStatus ServiceStatus `gorm:"foreignKey:Status;references:Status"`
}
