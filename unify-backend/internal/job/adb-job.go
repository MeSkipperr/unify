package job

import "github.com/google/uuid"

type ADBJob struct {
	ID        uuid.UUID `json:"id"`
	IPAddress string `json:"ipAddress"`
	Port      int    `json:"port"`
	Command   string `json:"command"`
	Name      string `json:"name"`
}

// Package   string `json:"package"`
