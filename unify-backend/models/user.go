package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleUser       UserRole = "USER"
	RoleAdmin      UserRole = "ADMIN"
	RoleSuperAdmin UserRole = "SUPERADMIN"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	// Identitas (USER & SUPERADMIN)
	FirstName *string `gorm:"type:varchar(100)"`
	LastName  *string `gorm:"type:varchar(100)"`

	// Email (USER & SUPERADMIN)
	Email *string `gorm:"type:varchar(150);uniqueIndex"`

	// Credential (ADMIN & SUPERADMIN)
	Username *string `gorm:"type:varchar(100);uniqueIndex"`
	Password *string `gorm:"type:text"`

	// Role
	Role UserRole `gorm:"type:varchar(20);not null"`

	// Status
	IsActive bool `gorm:"default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
