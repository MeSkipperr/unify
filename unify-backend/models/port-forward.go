package models

import (
	"time"

	"github.com/google/uuid"
)

type SessionStatus string

const (
	SessionStatusPending     SessionStatus = "pending"     // waiting to be applied
	SessionStatusActive      SessionStatus = "active"      // currently applied
	SessionStatusInactive    SessionStatus = "inactive"    // disabled
	SessionStatusExpired     SessionStatus = "expired"     // expired and removed
	SessionStatusDeactivated SessionStatus = "deactivated" // manually disabled
	SessionStatusError       SessionStatus = "error"       // error occurred during application

	// State transitions:
	//  pending -> active
	//  active -> if expires_at < now() -> expired
	//  active -> if manually disabled  -> deactivated -> inactive
)

type SessionPortForward struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	CreatedAt time.Time
	ExpiresAt time.Time
	Status    SessionStatus `gorm:"type:varchar(20);not null;default:'pending';index"`

	ListenIP   string `gorm:"type:varchar(45);not null;index"`
	ListenPort int    `gorm:"not null;index"`

	DestIP   string `gorm:"type:varchar(45);not null"`
	DestPort int    `gorm:"not null"`

	Protocol string `gorm:"type:varchar(10);not null;default:'tcp'"`

	Chain         string `gorm:"type:varchar(50);not null;default:'PORT_FORWARD'"`
	Interface     string `gorm:"type:varchar(30)"`
	AllowSourceIP string `gorm:"type:varchar(45);default:'0.0.0.0/0'"`
	RuleComment   string
	LastAppliedAt *time.Time
}
