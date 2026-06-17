package entity

import (
	"github.com/google/uuid"
)

// SelectOption stores reusable key-value labels assigned to a device for categorization.
type SelectOption struct {
	ID       int       `json:"id" gorm:"primaryKey;autoIncrement"`
	DeviceID uuid.UUID `gorm:"type:uuid;not null" json:"deviceId"`
	Label    string    `json:"label" gorm:"type:varchar(100);not null;index"`
	Value    string    `json:"value" gorm:"type:varchar(100);not null"`
}
