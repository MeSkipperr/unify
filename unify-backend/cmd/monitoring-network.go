package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"unify-backend/internal/core/arp"
	"unify-backend/internal/core/port"
	"unify-backend/internal/database"
	"unify-backend/internal/mailer"
	"unify-backend/internal/notification"
	"unify-backend/internal/services"
	"unify-backend/internal/template/email"
	"unify-backend/utils"

	"unify-backend/internal/worker"
	"unify-backend/models"

	"github.com/google/uuid"
)

type ConnectivityResult struct {
	MACAddress  string `json:"mac_address"`
	IsConnect   bool   `json:"is_connect"`
	ConnectPort int    `json:"connect_port"`
}

func selectDevicesByTypes(types []models.DeviceType) ([]models.Devices, error) {
	var devices []models.Devices

	result := database.DB.
		Where("type IN ?", types).
		Find(&devices)
	if result.Error != nil {
		return nil, result.Error
	}

	return devices, nil
}

func checkDeviceConnectivity(dev models.Devices) ConnectivityResult {
	result := ConnectivityResult{
		MACAddress:  dev.MacAddress,
		IsConnect:   false,
		ConnectPort: 0,
	}

	// STEP 1: Ping / ARP
	arpRes := arp.Check(arp.Params{
		IP:     dev.IPAddress,
		Warmup: true,
	})

	if arpRes.Exists {
		result.MACAddress = arpRes.MAC
		result.IsConnect = true
		result.ConnectPort = 0
		return result
	}

	// STEP 2: If cannot ping to device use testing port
	ports := models.CommonPorts()

	if typePorts := models.DevicePorts(dev.Type); len(typePorts) > 0 {
		ports = append(ports, typePorts...)
	}

	for _, p := range ports {
		portRes := port.Check(port.Params{
			Target:   dev.IPAddress,
			Port:     p,
			Protocol: port.TCP,
			Timeout:  2 * time.Second,
		})

		if portRes.Open {
			result.MACAddress = dev.MacAddress
			result.IsConnect = true
			result.ConnectPort = p
			return result
		}
	}

	// STEP 3: All Method cannot connect to device
	return result
}

func sendNotification(dev models.Devices, isConnect bool) {
	subject := fmt.Sprintf("[ALERT] %s - DOWN", dev.Name)
	if isConnect {
		subject = fmt.Sprintf("[ALERT] %s - UP", dev.Name)
	}

	notification.UserNotificationChannel(mailer.EmailData{
		Subject:        subject,
		BodyTemplate:   email.DeviceStatusEmail(dev, isConnect),
		FileAttachment: []string{},
	})
}

func updateDeviceStatus(
	deviceID uuid.UUID,
	isConnect bool,
	errorCount int,
) error {

	result := database.DB.
		Model(&models.Devices{}).
		Where("id = ?", deviceID).
		Updates(map[string]interface{}{
			"is_connect":        isConnect,
			"error_count":       errorCount,
			"status_updated_at": time.Now().UTC(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("device not found")
	}

	return nil
}

func processConnection(dev models.Devices, maxTimes int) {
	prevErrorCount := dev.ErrorCount
	prevIsConnect := dev.IsConnect

	result := checkDeviceConnectivity(dev)

	// update error count
	if result.IsConnect {
		dev.ErrorCount = utils.Clamp(dev.ErrorCount-1, 0, maxTimes)
	} else {
		dev.ErrorCount = utils.Clamp(dev.ErrorCount+1, 0, maxTimes)
	}

	// tentukan status baru
	newIsConnect := prevIsConnect

	// DOWN (dari connect → disconnect)
	if prevIsConnect && dev.ErrorCount == maxTimes {
		newIsConnect = false
		sendNotification(dev, false)
		services.LogInfo(ServiceMonitoringNetwork, "Device down: "+dev.Name)
	} else if !prevIsConnect && dev.ErrorCount == 0 {
		// RECOVER (dari disconnect → connect)
		newIsConnect = true
		sendNotification(dev, true)
		services.LogInfo(ServiceMonitoringNetwork, "Device recovered: "+dev.Name)
	}

	// ⛔ tidak ada perubahan → stop
	if prevErrorCount == dev.ErrorCount && prevIsConnect == newIsConnect {
		return
	}

	// update DB
	if err := updateDeviceStatus(dev.ID, newIsConnect, dev.ErrorCount); err != nil {
		services.LogError(ServiceMonitoringNetwork, "Failed to update device status: "+err.Error())
	}
}

type monitoringNetworkConfig struct {
	Delay         int      `json:"delay"`
	CheckingTimes int      `json:"checkingTimes"`
	DeviceTypes   []string `json:"deviceType"`
}

func MonitoringNetwork(manager *worker.Manager) (*worker.Worker, error) {

	var config monitoringNetworkConfig

	service, err := services.GetByServiceName(ServiceMonitoringNetwork)
	err = json.Unmarshal(service.Config, &config)

	if err != nil {
		return nil, err
	}

	w := worker.NewWorker(
		ServiceMonitoringNetwork,
		"",
		func() {
			types := make([]models.DeviceType, 0, len(config.DeviceTypes))

			for _, t := range config.DeviceTypes {
				types = append(types, models.DeviceType(t))
			}

			for {
				devices, err := selectDevicesByTypes(types)
				if err != nil {
					services.LogError(ServiceMonitoringNetwork, "Error selecting devices: "+err.Error())
					return
				}

				for _, dev := range devices {
					processConnection(dev, config.CheckingTimes)
				}

				time.Sleep(time.Duration(config.Delay))
			}
		},
	)

	w.RunOnce = true
	return w, nil
}
