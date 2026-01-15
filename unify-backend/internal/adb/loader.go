package adb

import (
	"encoding/json"
	"os"
)

type ADBConfig struct {
	ADBPath            string            `json:"adbPath"`
	ADBPort            int               `json:"adbPort"`
	Package            map[string]string `json:"package"`
	VerificationSteps  int               `json:"verificationSteps"`
	StatusMessage      map[string]string `json:"statusMessage"`
	CommandTemplate    map[string]string `json:"adbCommandTemplate"`
}

func LoadADBConfig(path string) (*ADBConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg ADBConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
