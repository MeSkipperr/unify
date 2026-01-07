package models

import (
	"github.com/google/uuid"

	"time"
)

type DeviceType string

const (
	CCTV   DeviceType = "cctv"
	IPTV   DeviceType = "iptv"
	AP     DeviceType = "access_point"
	SWIPTV DeviceType = "sw_iptv"
)

type Devices struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	IPAddress    string     `gorm:"size:45;not null"`
	IsConnect    bool       `gorm:"default:false"`
	ErrorCount   int        `gorm:"default:0"`
	Name         string     `gorm:"size:100;uniqueIndex;not null"`
	RoomNumber   string     `gorm:"size:50"`
	Description  string     `gorm:"type:text"`
	Type         DeviceType `gorm:"type:varchar(30);not null"`
	Notification bool       `gorm:"default:false"`
	MacAddress   string     `gorm:"size:50"`
	status_updated_at time.Time
}

var devicePorts = map[DeviceType][]int{
	CCTV: {554, 8000, 8899, 37777, 34567},
	IPTV: {80, 443, 22, 161, 8291},
	AP:   {80, 443, 5555, 8008, 8009},
}

func DevicePorts(t DeviceType) []int {
	return append([]int(nil), devicePorts[t]...)
}

var commonPorts = []int{80, 443, 22, 554}

func CommonPorts() []int {
	return append([]int(nil), commonPorts...)
}
