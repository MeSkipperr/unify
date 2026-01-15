package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"unify-backend/internal/adb"
	"unify-backend/internal/services"
	"unify-backend/internal/worker"
	"unify-backend/models"
)



func processGetUptime(dev models.Devices, adbConfig *adb.ADBConfig) (adb.AdbStatus, error) {
	//Check if device is connect to system
	if !dev.IsConnect {
		return adb.StatusNotConnected, fmt.Errorf("device %s is not connected", dev.IPAddress)
	}

	// Implementation for processing ADB devices goes here
	target := map[string]string{
		"adbPath": adbConfig.ADBPath,
		"ip":      dev.IPAddress,
		"port":    fmt.Sprintf("%d", adbConfig.ADBPort),
		"package": adbConfig.Package["youtube"],
	}

	clearCmd := adb.RenderTemplate(adbConfig.CommandTemplate["clearData"], target);
	fmt.Println("RUN:", clearCmd)

	return adb.StatusFailed, nil
}

func GetUptimeADB(manager *worker.Manager) (*worker.Worker, error) {
	const serviceName = "get-uptime-adb"
	adbConfig, err := adb.LoadADBConfig("internal/adb/adb.linux.json")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	config := struct {
		Cron        string   `json:"cron"`
		DeviceTypes []string `json:"deviceTypes"`
	}{
		Cron:        "*/10 * * * * *",
		DeviceTypes: []string{""},
	}

	service, err := services.GetByServiceName(serviceName)
	err = json.Unmarshal(service.Config, &config)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return worker.NewWorker(
		serviceName,
		config.Cron,
		func() {
			types := make([]models.DeviceType, 0, len(config.DeviceTypes))

			for _, t := range config.DeviceTypes {
				types = append(types, models.DeviceType(t))
			}

			devices, err := selectDevicesByTypes(types)
			if err != nil {
				log.Println("Error selecting devices:", err)
				return
			}

			err = adb.RestartADBServer(adbConfig.ADBPath)
			if err != nil {
				log.Println("Error restarting ADB server:", err)
				return 
			}

			for i := 0; i < adbConfig.VerificationSteps; i++ {

				for _, dev := range devices {
					resStatus, err := processGetUptime(dev, adbConfig)
					if err != nil {
						log.Println("Error processing device uptime:", err)
					}
				}
			}
		},
	), nil

}
