package entity

import (
	"time"

	"github.com/google/uuid"
)

// UserRole defines the access level of a user.
type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

// User stores registered user accounts and authentication data.
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FirstName string    `gorm:"size:255;not null" json:"firstName"`
	LastName  string    `gorm:"size:255;not null" json:"lastName"`
	Email     string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Role      UserRole  `gorm:"type:varchar(20);not null;default:'user'" json:"role"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}
