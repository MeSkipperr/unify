package entity

import (
	"github.com/google/uuid"

	"time"
)

// DeviceStatus represents the connectivity state of a device.
type DeviceStatus string

const (
	DeviceStatusOnline      DeviceStatus = "online"      // Device is reachable and responding normally.
	DeviceStatusOffline     DeviceStatus = "offline"     // Device is not responding, considered disconnected.
	DeviceStatusUnreachable DeviceStatus = "unreachable" // Device exists but cannot be reached (routing/firewall).
	DeviceStatusStealth     DeviceStatus = "stealth"     // Device may be online but ignores discovery requests.
)


// Devices stores network device information and connectivity status.
type Devices struct {
	ID            uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	IPAddress     string       `gorm:"size:45;not null" json:"ipAddress"`
	ErrorCount    int          `gorm:"default:0" json:"errorCount"`
	Status        DeviceStatus `gorm:"type:varchar(20);not null" json:"status"`
	Name          string       `gorm:"size:256;uniqueIndex;not null" json:"name"`
	RoomNumber    string       `gorm:"size:50" json:"roomNumber"`
	Description   string       `gorm:"type:text" json:"description"`
	Notification  bool         `gorm:"default:false" json:"notification"`
	MacAddress    string       `gorm:"size:50" json:"macAddress"`
	DeviceProduct string       `gorm:"type:varchar(255)" json:"deviceProduct"`
	LastSeen      time.Time    `json:"lastSeen"`
	LastError     time.Time    `json:"lastError"`
	CreatedAt     time.Time    `gorm:"autoCreateTime;index" json:"createdAt"`

	DeviceType []SelectOption  `gorm:"foreignKey:DeviceID;references:ID;constraint:OnDelete:CASCADE" json:"deviceType"`
	Services   []DeviceService `gorm:"foreignKey:DeviceID;references:ID;constraint:OnDelete:CASCADE" json:"services"`
}

// DeviceService stores network service endpoints exposed by a device (e.g., HTTP, SSH).
type DeviceService struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	DeviceID    uuid.UUID `gorm:"type:uuid;not null" json:"deviceId"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Port        int       `gorm:"not null" json:"port"`

	Device *Devices `gorm:"foreignKey:DeviceID;references:ID" json:"device,omitempty"`
}
