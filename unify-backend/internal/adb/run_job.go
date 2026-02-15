package adb

import (
	"fmt"
	"unify-backend/internal/job"
)

// DefaultConfigPath is the default path to the ADB config file.
const DefaultConfigPath = "internal/adb/adb.linux.json"

// RunJob runs an ADB job using the given config path. Returns status and output/error message.
func RunJob(j job.ADBJob, configPath string) (AdbStatus, string) {
	adbConfig, err := LoadADBConfig(configPath)
	if err != nil {
		return StatusFailed, err.Error()
	}
	data := map[string]string{
		"adbPath": adbConfig.ADBPath,
		"ip":      j.IPAddress,
		"port":    fmt.Sprintf("%d", j.Port),
		"package": adbConfig.Package["youtube"],
	}

	status, output := AdbRun(AdbRunRequest{
		Config:   adbConfig,
		Template: j.Command,
		Data:     data,
	})
	return status, output
}
