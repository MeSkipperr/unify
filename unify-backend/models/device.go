package models

import (
	"github.com/google/uuid"

	"time"
)

type DeviceType string

const (
	CCTV   DeviceType = "cctv"
	IPTV   DeviceType = "iptv"
	AP     DeviceType = "access-point"
	SW     DeviceType = "sw"
	SWIPTV DeviceType = "sw_iptv"
)

type DeviceAction string

const (
	DeviceCreated DeviceAction = "created"
	DeviceUpdated DeviceAction = "updated"
	DeviceDeleted DeviceAction = "deleted"
)

type Devices struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	IPAddress       string     `gorm:"size:45;not null" json:"ipAddress"`
	IsConnect       bool       `gorm:"default:false" json:"isConnect"`
	ErrorCount      int        `gorm:"default:0" json:"errorCount"`
	Name            string     `gorm:"size:100;uniqueIndex;not null" json:"name"`
	RoomNumber      string     `gorm:"size:50" json:"roomNumber"`
	Description     string     `gorm:"type:text" json:"description"`
	Type            DeviceType `gorm:"type:varchar(30);not null" json:"type"`
	Notification    bool       `gorm:"default:false" json:"notification"`
	MacAddress      string     `gorm:"size:50" json:"macAddress"`
	DeviceProduct   string     `gorm:"type:varchar(255)" json:"deviceProduct"`
	StatusUpdatedAt time.Time  `json:"statusUpdatedAt"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"createdAt"`
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

