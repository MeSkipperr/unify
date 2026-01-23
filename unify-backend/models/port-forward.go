package models

import (
	"time"

	"github.com/google/uuid"
)

type SessionStatus string

const (
	SessionStatusPending  SessionStatus = "pending"
	SessionStatusActive   SessionStatus = "active"
	SessionStatusExpired  SessionStatus = "expired"
	SessionStatusDisabled SessionStatus = "disabled"
	SessionStatusError    SessionStatus = "error"
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
