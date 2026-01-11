package cmd

import (
	"errors"
	"fmt"
	"log"
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

func selectDevices(types []models.DeviceType) ([]models.Devices, error) {
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
			"is_connect":         isConnect,
			"error_count":        errorCount,
			"status_updated_at":  time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("device not found")
	}

	return nil
}


func processConnection(dev models.Devices, times int) {
	isConnect := dev.IsConnect
	res := checkDeviceConnectivity(dev)

	if res.IsConnect {
		dev.ErrorCount = utils.Clamp(dev.ErrorCount-1, 0, times)
	} else {
		dev.ErrorCount = utils.Clamp(dev.ErrorCount+1, 0, times)
	}

	if dev.ErrorCount == 0 || dev.ErrorCount == times {
		if dev.IsConnect && dev.ErrorCount == times {
			//send error
			isConnect = false
			sendNotification(dev, false)
			services.CreateAppLog(services.CreateLogParams{
				Level:       "INFO",
				ServiceName: "monitoring-network",
				Message:     "Device down: " + dev.Name,
			})

		} else if !dev.IsConnect && dev.ErrorCount == 0 {
			// send recovery
			isConnect = true
			sendNotification(dev, true)
			services.CreateAppLog(services.CreateLogParams{
				Level:       "INFO",
				ServiceName: "monitoring-network",
				Message:     "Device recovered: " + dev.Name,
			})
		}
	}
	//update db
	err := updateDeviceStatus(dev.ID, isConnect, dev.ErrorCount)
	if err != nil {
		services.CreateAppLog(services.CreateLogParams{
			Level:       "ERROR",
			ServiceName: "monitoring-network",
			Message:     "Failed to update device status: " + err.Error(),
		})
	}
}

func MonitoringNetwork(manager *worker.Manager) (*worker.Worker, error) {
	err := services.CreateAppLog(services.CreateLogParams{
		Level:       "INFO",
		ServiceName: "monitoring-network",
		Message:     "Monitoring Network Service Started",
	})

	if err != nil {
		log.Fatal(err)
	}

	return worker.NewWorker(
		"monitoring-network",
		"* * * * * *",
		func() {
			types := []models.DeviceType{
				models.CCTV,
				models.IPTV,
			}

			devices, err := selectDevices(types)
			if err != nil {
				log.Println("Error selecting devices:", err)
				return
			}

			for _, dev := range devices {
				processConnection(dev, 2)
			}
		},
	), nil
}
