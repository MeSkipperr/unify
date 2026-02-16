package cmd

import (
	"encoding/json"
	"errors"
	"time"
	"unify-backend/internal/core/ip"
	"unify-backend/internal/database"
	"unify-backend/internal/http/sse"
	"unify-backend/internal/services"
	"unify-backend/internal/speedtest"
	"unify-backend/internal/worker"
	"unify-backend/models"
	"unify-backend/utils"

	"gorm.io/gorm"
)

type getSpeedtestNetworkConfig struct {
	Cron     string          `json:"cron"`
	Network  []networkConfig `json:"network"`
	ServerID []string        `json:"server_id"`
}

type networkConfig struct {
	Name      string `json:"name"`
	Interface string `json:"interface"`
	IPAddress string `json:"ip_address"`
}

func sendSSESpeedtest(manager *worker.Manager, data models.SpeedtestResult) {
	sseManager := manager.GetSSE()

	res := sse.ServicesEvent{
		Type: ServiceGetSpeedtestNetwork,
		Data: data,
	}

	if sseManager != nil {
		sseManager.Broadcast(sse.SSEChannelServices, res)
	}

}

func GetSpeedtestNetwork(manager *worker.Manager) (*worker.Worker, error) {

	var config getSpeedtestNetworkConfig

	service, err := services.GetByServiceName(ServiceGetSpeedtestNetwork)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			services.LogInfo(ServiceGetSpeedtestNetwork, "service get-speedtest-network not found, worker disabled")
			return nil, nil
		}
		return nil, err
	}

	err = json.Unmarshal(service.Config, &config)
	if err != nil {
		return nil, err
	}
	return worker.NewWorker(
		ServiceGetSpeedtestNetwork,
		config.Cron,
		func() {
			services.LogInfo(ServiceGetSpeedtestNetwork, "Starting Get Speedtest Network Service")
			for _, network := range config.Network {
				client, err := ip.CheckLocalIp(network.IPAddress)
				if err != nil {
					services.LogError(ServiceGetSpeedtestNetwork, "Failed to get local IP info: "+err.Error())
					continue
				}
				if !client.Exists {
					services.LogError(ServiceGetSpeedtestNetwork, "Local IP "+network.IPAddress+" not found on any interface")
					continue
				}
				for _, serverID := range config.ServerID {
					if serverID == "" {
						continue
					}
					services.LogInfo(ServiceGetSpeedtestNetwork, "Running speedtest for IP "+network.IPAddress+" and server ID "+serverID)
					result, err := speedtest.Run(network.IPAddress, serverID)
					if err != nil {
						services.LogError(ServiceGetSpeedtestNetwork, "Failed to run speedtest for IP "+network.IPAddress+" and server ID "+serverID+": "+err.Error())
						continue
					}
					resultRecord := models.SpeedtestResult{
						TestedAt: time.Now().UTC(),

						NetworkName:   network.Name,
						InterfaceName: result.Interface.Name,
						InternalIP:    result.Interface.InternalIP,
						ExternalIP:    result.Interface.ExternalIP,
						MACAddress:    result.Interface.MACAddr,
						ISPName:       result.ISP,

						ServerID:       result.Server.ID,
						ServerName:     result.Server.Name,
						ServerCountry:  result.Server.Country,
						ServerLocation: result.Server.Location,

						DownloadMbps: utils.BytesPerSecToMbps(result.Download.Bandwidth),
						UploadMbps:   utils.BytesPerSecToMbps(result.Upload.Bandwidth),
						PingMs:       result.Ping.Latency,
						ResultURL:    result.Result.URL,
					}

					// Save resultRecord to database
					if err := database.DB.Create(&resultRecord).Error; err != nil {
						services.LogError(ServiceGetSpeedtestNetwork, "Failed to save speedtest result for IP "+network.IPAddress+" and server ID "+serverID+": "+err.Error())
						continue
					}
					sendSSESpeedtest(manager, resultRecord)
				}
			}
			services.LogInfo(ServiceGetSpeedtestNetwork, "Completed Get Speedtest Network Service")
		},
	), nil

}
