package models

import "time"

type DeviceType string

const (
	CCTV DeviceType = "cctv"
	IPTV DeviceType = "iptv"
	AP   DeviceType = "access_point"
)

type DeviceStruct struct {
	DeviceID     string     `json:"device_id"`
	IPAddress    string     `json:"ip_address"`
	IsConnect    bool       `json:"is_connect"`
	ErrorCount   int        `json:"error_count"`
	Name         string     `json:"name"`
	RoomNumber   string     `json:"room_number"`
	Description  string     `json:"description"`
	Type         DeviceType `json:"type"`
	Notification bool       `json:"notification"`
	LastChange   time.Time  `json:"last_change"`
	MacAddress   string     `json:"mac_address"`
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
