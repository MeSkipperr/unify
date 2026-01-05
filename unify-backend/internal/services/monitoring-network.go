package services

import (
	"fmt"
	"log"
	"time"
	"unify-backend/internal/core/arp"
	"unify-backend/internal/core/port"
	"unify-backend/utils"

	"unify-backend/internal/worker"
	"unify-backend/internal/ws"
	"unify-backend/models"
)

type ConnectivityResult struct {
	MACAddress  string `json:"mac_address"`
	IsConnect   bool   `json:"is_connect"`
	ConnectPort int    `json:"connect_port"`
}

func CheckDeviceConnectivity(dev models.DeviceStruct) ConnectivityResult {
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

func sendNotification(dev models.DeviceStruct, errorNotification bool) {

}

func ProcessConnection(dev models.DeviceStruct) {
	times := 2

	res := CheckDeviceConnectivity(dev)

	if res.IsConnect {
		dev.ErrorCount = utils.Clamp(dev.ErrorCount-1, 0, times)
	} else {
		dev.ErrorCount = utils.Clamp(dev.ErrorCount+1, 0, times)
	}

	if dev.ErrorCount == 0 || dev.ErrorCount == times {
		if dev.IsConnect && dev.ErrorCount == times {
			//send error
			sendNotification(dev, true)
		} else if !dev.IsConnect && dev.ErrorCount == 0 {
			// send recovery
			sendNotification(dev, false)
		}
	}
	//create log
	//update db
}

func MonitoringNetwork(manager *worker.Manager) (*worker.Worker, error) {

	return worker.NewWorker(
		"monitoring-network",
		"* * * * * *",
		func() {
			log.Println("project1 task running")
			i := 1

			//PING
			// res := ping.Ping(ping.Params{
			// 	Target: "1.1.1.1",
			// 	Times:  1,
			// })

			// fmt.Printf("%+v\n", res)

			// res := port.Check(port.Params{
			// 	Target:   "0.0.0.0",
			// 	Port:     3000,
			// 	Protocol: port.TCP,
			// })

			res := arp.Check(arp.Params{
				IP: "172.19.176.1",
				// Interface: "eth0",
				Warmup: true,
			})
			fmt.Printf("%+v\n", res)

			msg := ws.Message{
				Time: time.Now(),
				ID:   i,
				Message: map[string]interface{}{
					"status": "running",
					"cpu":    42,
				},
			}

			i++

			manager.BroadcastProject(msg)
		},
	), nil
}
