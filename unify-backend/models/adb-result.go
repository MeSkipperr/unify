package models

import (
	"time"
	"unify-backend/internal/adb"

	"github.com/google/uuid"
)

type AdbResult struct {
	ID           uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Status       adb.AdbStatus `gorm:"not null"`
	StartTime    time.Time     `gorm:"not null"`
	FinishTime   time.Time     `gorm:"not null"`
	IPAddress    string        `gorm:"size:45;not null"`
	Port         int           `gorm:"not null"`
	NameDevice   string        `gorm:"size:100;uniqueIndex;not null"`
	Result       string        `gorm:"type:text"`
	TypeServices string        `gorm:"type:text;not null"`

}
