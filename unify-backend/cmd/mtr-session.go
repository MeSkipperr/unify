package cmd

import (
	// "errors"
	"fmt"
	"log"
	"time"
	"unify-backend/internal/core/dns"
	"unify-backend/internal/core/mtr"
	"unify-backend/internal/database"
	"unify-backend/internal/services"
	"unify-backend/internal/worker"
	"unify-backend/internal/ws"
	"unify-backend/models"

	"gorm.io/gorm"
)

type MTRSessionConfig struct {
	Interval           int `json:"interval"`
	RangeReachableLoss int `json:"range_reachable_loss"`
}

type mtrSessionParms struct {
	db      *gorm.DB
	session models.MTRSession
	out     mtr.Result
	config  MTRSessionConfig
	manager *worker.Manager
}

func saveMtrResult(data mtrSessionParms) error {
	// Save the MTR result to the database
	db := data.db
	session := data.session
	out := data.out

	mtrResult := models.MTRResult{
		SessionID: session.ID.String(),

		SourceIP:      out.Hops[0].IP,
		DestinationIP: session.DestinationIP,
		Protocol:      session.Protocol,
		Port:          session.Port,
		Test:          session.Test,

		TotalHops: out.TotalHops,
		Reachable: out.Reachable,
		MaxLoss:   out.MaxLoss,
		AvgRTT:    out.AvgRTT,
	}

	result := db.Create(&mtrResult)
	if result.Error != nil {
		return result.Error
	}

	for _, hop := range out.Hops {

		mtrHop := models.MTRHop{
			ResultID:  mtrResult.ID.String(),
			Hop:       hop.Hop,
			Host:      hop.IP,
			Loss:      hop.Loss,
			Avg:       hop.AvgRTT,
			IpAddress: hop.IP,
			DNS:       dns.ReverseDNS(hop.IP)[0],
		}
		result = db.Create(&mtrHop)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func updateMtrSessionLastRun(data mtrSessionParms) error {
	var session models.MTRSession
	result := data.db.First(&session, "id = ?", data.session.ID)
	if result.Error != nil {
		return result.Error
	}

	now := time.Now()

	session.LastRunAt = &now
	return data.db.Save(&session).Error
}

func sendLossConnectionAlert(data mtrSessionParms) {
	// Implement alerting logic here (e.g., send email or notification)
	db := data.db
	session := data.session
	config := data.config

	var reachableLogs []bool

	err := db.
		Model(&models.MTRResult{}).
		Select("reachable").
		Where("session_id = ?", session.ID).
		Order("created_at DESC").
		Limit(config.RangeReachableLoss).
		Pluck("reachable", &reachableLogs).Error

	if err != nil {
		log.Println(err)
	}

	allUnreachable := false

	for _, v := range reachableLogs {
		if !v {
			allUnreachable = false
			break
		}
	}

	if allUnreachable {
		alertMsg := fmt.Sprintf("MTR Session Alert: Destination %s is unreachable for the last %d checks.", session.DestinationIP, config.RangeReachableLoss)
		services.LogWarning(ServiceMTRSession, alertMsg)
	}
}

func sendPacketUseWebsocket(data mtrSessionParms) {
	// Implement WebSocket sending logic here
	fmt.Print("SEND WEB SOCKET")

	msg := ws.Message{
		Time:    time.Now(),
		ID:      data.session.ID.String(),
		Message: data.out,
	}

	data.manager.BroadcastProject(msg)
}

func startSyncSessionMTRWorker(db *gorm.DB, config MTRSessionConfig, manager *worker.Manager) {
	ticker := time.NewTicker(time.Duration(config.Interval) * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		var sessions []models.MTRSession
		result := db.Find(&sessions).Where("status = ?", "active")
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

				params := mtrSessionParms{
					db:      db,
					session: s,
					config:  config,
					out:     *out,
					manager: manager,
				}

				saveMtrResult(params)
				updateMtrSessionLastRun(params)
				sendLossConnectionAlert(params)
				sendPacketUseWebsocket(params)

				fmt.Println(out)
			}(session)
		}
	}
}

func RunMTRSession(manager *worker.Manager) (*worker.Worker, error) {
	db := database.DB

	config := MTRSessionConfig{
		Interval:           10,
		RangeReachableLoss: 10,
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

			go startSyncSessionMTRWorker(db, config, manager)
		},
	)

	// worker manager yang mengontrol lifecycle
	w.RunOnce = true
	return w, nil
}
