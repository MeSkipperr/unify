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
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	CreatedAt time.Time     `json:"createdAt"`
	ExpiresAt time.Time     `json:"expiresAt"`
	Status    SessionStatus `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`

	ListenIP   string `gorm:"type:varchar(45);not null;index" json:"listenIp"`
	ListenPort int    `gorm:"not null;index" json:"listenPort"`

	DestIP   string `gorm:"type:varchar(45);not null" json:"destIp"`
	DestPort int    `gorm:"not null" json:"destPort"`

	Protocol string `gorm:"type:varchar(10);not null;default:'tcp'" json:"protocol"`

	Chain         string     `gorm:"type:varchar(50);not null;default:'PORT_FORWARD'" json:"chain"`
	Interface     string     `gorm:"type:varchar(30)" json:"interface"`
	AllowSourceIP string     `gorm:"type:varchar(45);default:'0.0.0.0/0'" json:"allowSourceIp"`
	RuleComment   string     `json:"ruleComment"`
	LastAppliedAt *time.Time `json:"lastAppliedAt"`
}
