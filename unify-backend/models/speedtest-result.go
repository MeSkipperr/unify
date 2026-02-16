package models

import (
	"time"

	"github.com/google/uuid"
)


type SpeedtestResult struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	TestedAt time.Time `gorm:"not null;index" json:"testedAt"`

	// Client / Interface
	NetworkName   string `gorm:"type:varchar(50);index" json:"networkName"`
	InterfaceName string `gorm:"type:varchar(50)" json:"interfaceName"`
	InternalIP    string `gorm:"type:varchar(45);index" json:"internalIP"`
	ExternalIP    string `gorm:"type:varchar(45)" json:"externalIP"`
	MACAddress    string `gorm:"type:varchar(20)" json:"macAddress"`
	ISPName       string `gorm:"type:varchar(100)" json:"ispName"`

	// Server
	ServerID       int    `gorm:"index" json:"serverId"`
	ServerName     string `gorm:"type:varchar(100)" json:"serverName"`
	ServerCountry  string `gorm:"type:varchar(50)" json:"serverCountry"`
	ServerLocation string `gorm:"type:varchar(100)" json:"serverLocation"`

	// Result
	DownloadMbps float64 `gorm:"type:decimal(8,2);comment:Download speed in Mbps" json:"downloadMbps"`
	UploadMbps   float64 `gorm:"type:decimal(8,2);comment:Upload speed in Mbps" json:"uploadMbps"`
	PingMs       float64 `gorm:"type:decimal(6,2);comment:Ping latency in ms" json:"pingMs"`
	ResultURL    string  `gorm:"type:text" json:"resultUrl"`

	CreatedAt time.Time `json:"createdAt"`
}
