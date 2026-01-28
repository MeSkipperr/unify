package models

import (
	"time"

	"github.com/google/uuid"
)


type SpeedtestResult struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	TestedAt time.Time `gorm:"not null;index"`

	// Client / Interface
	NetworkName  string `gorm:"type:varchar(50);index"`
	InterfaceName string `gorm:"type:varchar(50)"`
	InternalIP   string `gorm:"type:varchar(45);index"`
	ExternalIP   string `gorm:"type:varchar(45)"`
	MACAddress   string `gorm:"type:varchar(20)"`
	ISPName      string `gorm:"type:varchar(100)"`

	// Server
	ServerID       int    `gorm:"index"`
	ServerName     string `gorm:"type:varchar(100)"`
	ServerCountry  string `gorm:"type:varchar(50)"`
	ServerLocation string `gorm:"type:varchar(100)"`

	// Result
	DownloadMbps  float64 `gorm:"type:decimal(8,2);comment:Download speed in Mbps"`
	UploadMbps    float64 `gorm:"type:decimal(8,2);comment:Upload speed in Mbps"`
	PingMs        float64 `gorm:"type:decimal(6,2);comment:Ping latency in ms"`
	ResultURL     string `gorm:"type:text"`

	CreatedAt time.Time
}
