package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"unify-backend/internal/adb"
	"unify-backend/internal/database"
	"unify-backend/internal/services"
	"unify-backend/internal/worker"
	"unify-backend/models"
)

func processGetUptime(dev models.AdbResult, adbConfig *adb.ADBConfig) (adb.AdbStatus, string) {
	data := map[string]string{
		"adbPath": adbConfig.ADBPath,
		"ip":      dev.IPAddress,
		"port":    fmt.Sprintf("%d", adbConfig.ADBPort),
	}

	status, uptimeOutput := adb.AdbRun(adb.AdbRunRequest{
		Config:   adbConfig,
		Template: adbConfig.CommandTemplate["getUptime"],
		Data:     data,
	})
	if status == adb.StatusSuccess {
		parts := strings.Split(uptimeOutput, " ")

		if len(parts) > 0 {
			uptimeSeconds, err := strconv.ParseFloat(parts[0], 64)
			if err != nil {
				return adb.StatusFailedUptime, "Failed to parse uptime output"
			} else {
				uptimeDays := uptimeSeconds / (60 * 60 * 24)
				return status, fmt.Sprintf("Success - Uptime %.2f Days", uptimeDays)
			}
		} else {
			return adb.StatusFailedUptime, "Unexpected uptime output format"
		}
	}
	return status, uptimeOutput
}

type getUptimeConfig struct {
	Cron        string   `json:"cron"`
	DeviceTypes []string `json:"deviceType"`
}

func GetUptimeADB(manager *worker.Manager) (*worker.Worker, error) {
	adbConfig, err := adb.LoadADBConfig("internal/adb/adb.linux.json")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	config := getUptimeConfig{
		Cron:        "0 0 10 * * *",
		DeviceTypes: []string{""},
	}

	service, err := services.GetByServiceName(ServiceGetUptimeADB)
	err = json.Unmarshal(service.Config, &config)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return worker.NewWorker(
		ServiceGetUptimeADB,
		config.Cron,
		func() {
			services.LogInfo(ServiceGetUptimeADB, "Starting Get Uptime ADB Service")

			types := make([]models.DeviceType, 0, len(config.DeviceTypes))

			for _, t := range config.DeviceTypes {
				types = append(types, models.DeviceType(t))
			}

			devices, err := selectDevicesByTypes(types)
			if err != nil {
				log.Println("Error selecting devices:", err)
				return
			}

			adbResult := []models.AdbResult{}
			for _, dev := range devices {
				status := adb.StatusNotStarted

				if !dev.IsConnect {
					status = adb.StatusNotConnected
				}

				adbResult = append(adbResult, models.AdbResult{
					Status:     status,
					StartTime:  time.Now(),
					FinishTime: time.Now(),
					IPAddress:  dev.IPAddress,
					Port:       adbConfig.ADBPort,
					NameDevice: dev.Name,
				})
			}

			for i := 0; i < adbConfig.VerificationSteps; i++ {
				err = adb.RestartADBServer(adbConfig.ADBPath)
				if err != nil {
					log.Println("Error restarting ADB server:", err)
					services.LogError(ServiceGetUptimeADB, "Error restarting ADB server: "+err.Error())
					return
				}
				for i := range adbResult {
					if adbResult[i].Status == adb.StatusSuccess || adbResult[i].Status == adb.StatusNotConnected {
						continue
					}
					startedAt := time.Now()
					resStatus, value := processGetUptime(adbResult[i], adbConfig)
					services.LogInfo(ServiceGetUptimeADB, "Device "+adbResult[i].NameDevice+" - "+fmt.Sprintf("%v", resStatus)+" - "+value)

					adbResult[i].Status = resStatus
					adbResult[i].FinishTime = time.Now()
					adbResult[i].StartTime = startedAt
					adbResult[i].Result = value
				}
			}

			for _, res := range adbResult {
				err := database.DB.Create(&res).Error
				if err != nil {
					services.LogError(ServiceGetUptimeADB, "Failed to save ADB result for device "+res.NameDevice+": "+err.Error())
				}
			}

			services.LogInfo(ServiceGetUptimeADB, "Completed Get Uptime ADB Service")
		},
	), nil

}
