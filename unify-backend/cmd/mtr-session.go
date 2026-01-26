package cmd

import (
	// "errors"
	"fmt"
	"log"
	"time"
	"unify-backend/internal/core/mtr"
	"unify-backend/internal/database"
	"unify-backend/internal/services"
	"unify-backend/internal/worker"
	"unify-backend/models"

	"gorm.io/gorm"
)

// func saveMtrResult(db *gorm.DB, sessionID string, out mtr.Result) error {
// 	// Save the MTR result to the database
// 	return nil
// }

// func updateMtrSessionLastRun(db *gorm.DB, sessionID string, lastRun time.Time) error {
// 	var session models.MTRSession
// 	result := db.First(&session, "id = ?", sessionID)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	// session.LastRunAt = lastRun
// 	return db.Save(&session).Error
// }

// func sendLossConnectionAlert(session models.MTRSession, packetLoss float64) {
// 	// Implement alerting logic here (e.g., send email or notification)
// }

// func sendPacketUseWebsocket(session models.MTRSession, out mtr.Result) {
// 	// Implement WebSocket sending logic here
// }

func startSyncSessionMTRWorker(db *gorm.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		var sessions []models.MTRSession
		result := db.Find(&sessions).Where("enabled = ?", true)
		if result.Error != nil {
			services.LogError(ServiceMTRSession, "Failed to fetch MTR sessions: "+result.Error.Error())
			continue
		}

		for _, session := range sessions {
			go func(s models.MTRSession) {
				// Here you would add the logic to perform MTR and store results
				services.LogInfo(ServiceMTRSession, "Processing MTR session ID: "+s.ID.String())
				out, err := mtr.Run(mtr.Config{
					DestHost: s.DestinationIP,
					SourceIP: s.SourceIP,
					Protocol: mtr.Protocol(s.Protocol),
					Port:     s.Port,
					Count:    s.Test,
					JSON:     true,
					UseDNS:   false,
				})

				if err != nil {
					services.LogError(ServiceMTRSession, "MTR run failed for session ID "+s.ID.String()+": "+err.Error())
					return
				}
				log.Println("[MTR RESULT]")

				fmt.Println(string(out))
			}(session)
		}
	}
}

type MTRSessionConfig struct {
	Interval int `json:"interval"`
}

func RunMTRSession(manager *worker.Manager) (*worker.Worker, error) {
	db := database.DB

	config := MTRSessionConfig{
		Interval: 10,
	}

	// service, err := services.GetByServiceName(ServiceMTRSession)
	// if err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		log.Println("service mtr-session not found, worker disabled")
	// 		return nil, nil
	// 	}
	// 	return nil, err
	// }

	// err = json.Unmarshal(service.Config, &config)
	// if err != nil {
	// 	return nil, err
	// }

	w := worker.NewWorker(
		ServiceMTRSession,
		"",
		func() {
			services.LogInfo(ServiceMTRSession, "Starting MTR Session Session Worker")

			go startSyncSessionMTRWorker(db, time.Duration(config.Interval)*time.Second)
		},
	)

	// worker manager yang mengontrol lifecycle
	w.RunOnce = true
	return w, nil
}
