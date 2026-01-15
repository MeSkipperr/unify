package models

import (
	"time"
	"unify-backend/internal/adb"

	"github.com/google/uuid"
)

type DeviceState struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Status  adb.AdbStatus
	startTime  time.Time
	finishTime time.Time
	CreatedAt  time.Time
	IPAddress  string `gorm:"size:45;not null"`
	Port       int    `gorm:"not null"`
	NameDevice string `gorm:"size:100;uniqueIndex;not null"`
	
}
