package models

import (
	"time"

	"github.com/google/uuid"
)

type MTRSession struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	Status      string `gorm:"type:varchar(20);not null;index"`
	IsReachable bool   `gorm:"not null; default:false"`

	CreatedAt time.Time

	LastRunAt *time.Time `gorm:"index"`

	SourceIP      string `gorm:"type:varchar(45)"`
	DestinationIP string `gorm:"type:varchar(45);not null;index"`

	Protocol string `gorm:"type:varchar(10);not null;default:'icmp';index"`
	Port     *int

	Test int `gorm:"not null;default:10"`

	Note string `gorm:"type:varchar(500);default:'';comment:'Additional notes'"`

	SendNotification bool `gorm:"not null;default:false"`
}

type MTRResult struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	SessionID string    `gorm:"type:uuid;index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	SourceIP      string `gorm:"type:varchar(45)"`
	DestinationIP string `gorm:"type:varchar(45)"`
	Protocol      string `gorm:"type:varchar(10)"`
	Port          *int
	Test          int `gorm:"not null;default:10"`

	TotalHops int
	Reachable bool
	AvgRTT    float64

	Hops []MTRHop `gorm:"foreignKey:ResultID;constraint:OnDelete:CASCADE"`
}

type MTRHop struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ResultID string    `gorm:"type:uuid;index"`

	Hop  int    `gorm:"index"`
	Host string `gorm:"type:varchar(255)"`
	DNS  string `gorm:"type:varchar(255);index"`

	Loss   float64
	Sent   int
	Last   float64
	Avg    float64
	Best   float64
	Worst  float64
	StdDev float64
}
