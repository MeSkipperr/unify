package models

import (
	"time"

	"github.com/google/uuid"
)

type Log struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt   time.Time
	Level       string  `gorm:"type:varchar(10);not null;index"`
	Message     string  `gorm:"type:text;not null"`
	ServiceName string  `gorm:"type:varchar(100);not null;index	"`
}
