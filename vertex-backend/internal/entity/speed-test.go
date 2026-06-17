package entity

import (
	"time"

	"github.com/google/uuid"
)

// SpeedTestJob stores scheduled network speed test configurations.
// Relation: SpeedTestJob -> SpeedTestResult
type SpeedTestJob struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name            string     `gorm:"size:255;not null" json:"name"`
	NetworkName     string     `gorm:"size:255;not null" json:"networkName"`
	InterfaceName   string     `gorm:"size:255;not null" json:"interfaceName"`
	IPAddress       string     `gorm:"size:45;not null" json:"ipAddress"`
	ISP             string     `gorm:"size:255;not null" json:"isp"`
	ServerID        string     `gorm:"size:255;not null" json:"serverId"`
	ServerName      string     `gorm:"size:255;not null" json:"serverName"`
	ServerLocation  string     `gorm:"size:255;not null" json:"serverLocation"`
	IntervalSeconds int        `gorm:"not null;default:0" json:"intervalSeconds"`
	Enabled         bool       `gorm:"default:true" json:"enabled"`
	LastRunAt       *time.Time `json:"lastRunAt"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"createdAt"`

	Results []SpeedTestResult `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE" json:"results"`
}

// SpeedTestResult stores the output of a speed test execution (ping, download, upload).
type SpeedTestResult struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	JobID        uuid.UUID  `gorm:"type:uuid;not null" json:"jobId"`
	Timestamp    int64      `gorm:"not null" json:"timestamp"`
	PingMs       float64    `gorm:"not null;default:0" json:"pingMs"`
	DownloadMbps float64    `gorm:"not null;default:0" json:"downloadMbps"`
	UploadMbps   float64    `gorm:"not null;default:0" json:"uploadMbps"`
	StartedAt    time.Time  `json:"startedAt"`
	CompletedAt  *time.Time `json:"completedAt"`

	Job *SpeedTestJob `gorm:"foreignKey:JobID;references:ID" json:"job,omitempty"`
}
