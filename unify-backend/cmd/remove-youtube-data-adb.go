package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"unify-backend/internal/adb"
	"unify-backend/internal/database"
	"unify-backend/internal/services"
	"unify-backend/internal/worker"
	"unify-backend/models"

	"gorm.io/gorm"
)

func processRemoveDataYoutube(dev models.AdbResult, adbConfig *adb.ADBConfig) (adb.AdbStatus, string) {
	data := map[string]string{
		"adbPath": adbConfig.ADBPath,
		"ip":      dev.IPAddress,
		"port":    fmt.Sprintf("%d", adbConfig.ADBPort),
		"package": adbConfig.Package["youtube"],
	}

	status, outputRemoveDataYoutube := adb.AdbRun(adb.AdbRunRequest{
		Config:   adbConfig,
		Template: adbConfig.CommandTemplate["clearData"],
		Data:     data,
	})
	return status, outputRemoveDataYoutube
}

type removeDataYoutubeConfig struct {
	Cron        string   `json:"cron"`
	DeviceTypes []string `json:"deviceType"`
}

func RemoveDataYoutubeADB(manager *worker.Manager) (*worker.Worker, error) {
	adbConfig, err := adb.LoadADBConfig("internal/adb/adb.linux.json")
	if err != nil {
		services.LogError(ServiceRemoveDataYoutubeADB, "Failed to load ADB config: "+err.Error())
		return nil, err
	}

	var config removeDataYoutubeConfig

	service, err := services.GetByServiceName(ServiceRemoveDataYoutubeADB)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			services.LogInfo(ServiceRemoveDataYoutubeADB, "service remove-youtube-data-adb not found, worker disabled")
			return nil, nil
		}
		return nil, err
	}

	err = json.Unmarshal(service.Config, &config)
	if err != nil {
		return nil, err
	}

	return worker.NewWorker(
		ServiceRemoveDataYoutubeADB,
		config.Cron,
		func() {
			services.LogInfo(ServiceRemoveDataYoutubeADB, "Starting Remove Data Youtube ADB Service")

			types := make([]models.DeviceType, 0, len(config.DeviceTypes))

			for _, t := range config.DeviceTypes {
				types = append(types, models.DeviceType(t))
			}

			devices, err := selectDevicesByTypes(types)
			if err != nil {
				services.LogError(ServiceRemoveDataYoutubeADB, "Error selecting devices: "+err.Error())
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
					StartTime:  time.Now().UTC(),
					FinishTime: time.Now().UTC(),
					IPAddress:  dev.IPAddress,
					Port:       adbConfig.ADBPort,
					NameDevice: dev.Name,
				})
			}

			for i := 0; i < adbConfig.VerificationSteps; i++ {
				err = adb.RestartADBServer(adbConfig.ADBPath)
				if err != nil {
					services.LogError(ServiceRemoveDataYoutubeADB, "Error restarting ADB server: "+err.Error())
					return
				}
				for i := range adbResult {
					if adbResult[i].Status == adb.StatusSuccess || adbResult[i].Status == adb.StatusNotConnected {
						continue
					}
					startedAt := time.Now().UTC()
					resStatus, value := processRemoveDataYoutube(adbResult[i], adbConfig)
					services.LogInfo(ServiceRemoveDataYoutubeADB, "Device "+adbResult[i].NameDevice+" - "+fmt.Sprintf("%v", resStatus)+" - "+value)

					adbResult[i].Status = resStatus
					adbResult[i].FinishTime = time.Now().UTC()
					adbResult[i].StartTime = startedAt
					adbResult[i].Result = value
					adbResult[i].TypeServices = ServiceRemoveDataYoutubeADB
				}
			}

			for _, res := range adbResult {
				err := database.DB.Create(&res).Error
				if err != nil {
					services.LogError(ServiceRemoveDataYoutubeADB, "Failed to save ADB result for device "+res.NameDevice+": "+err.Error())
				}
			}

			services.LogInfo(ServiceRemoveDataYoutubeADB, "Completed Get Uptime ADB Service")
		},
	), nil

}
