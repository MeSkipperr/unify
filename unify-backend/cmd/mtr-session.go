package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"unify-backend/internal/core/mtr"
	"unify-backend/internal/database"
	"unify-backend/internal/mailer"
	"unify-backend/internal/notification"
	"unify-backend/internal/services"
	"unify-backend/internal/template/email"
	"unify-backend/internal/worker"
	"unify-backend/internal/ws"
	"unify-backend/models"

	"gorm.io/gorm"
)

type MTRSessionConfig struct {
	Cron               string `json:"cron"`
	RangeReachableLoss int    `json:"range_reachable_loss"`
}

type mtrSessionParms struct {
	db      *gorm.DB
	session models.MTRSession
	out     mtr.MtrResultJson
	config  MTRSessionConfig
	manager *worker.Manager
}

type reachableStatus string

const (
	UP      reachableStatus = "UP"
	DOWN    reachableStatus = "DOWN"
	PARTIAL reachableStatus = "PARTIAL"
)

func saveMtrResult(data mtrSessionParms) error {
	// Save the MTR result to the database
	db := data.db
	session := data.session
	out := data.out.Report

	mtrResult := models.MTRResult{
		SessionID: session.ID.String(),

		SourceIP:      out.HopResult[0].Host,
		DestinationIP: session.DestinationIP,
		Protocol:      session.Protocol,
		Port:          session.Port,
		Test:          session.Test,

		TotalHops: out.Result.TotalHops,
		Reachable: out.Result.Reachable,
		AvgRTT:    out.Result.AvgRTT,
	}

	result := db.Create(&mtrResult)
	if result.Error != nil {
		return result.Error
	}

	for _, hop := range out.HopResult {

		mtrHop := models.MTRHop{
			ResultID: mtrResult.ID.String(),
			Hop:      hop.Count,
			Host:     hop.Host,
			DNS:      hop.Dns,
			Loss:     hop.Loss,
			Sent:     hop.Snt,
			Last:     hop.Last,
			Avg:      hop.Avg,
			Best:     hop.Best,
			Worst:    hop.Worst,
			StdDev:   hop.Worst,
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

func checkReachable(logs []bool) reachableStatus {
	allUp := true
	allDown := true

	for _, v := range logs {
		if v {
			allDown = false
		} else {
			allUp = false
		}
	}

	switch {
	case allUp:
		return UP
	case allDown:
		return DOWN
	default:
		return PARTIAL
	}
}

func updateStatusReachableSession(session models.MTRSession, isReachable bool) error {
	result := database.DB.
		Model(&models.MTRSession{}).
		Where("id = ?", session.ID).
		Update("is_reachable", isReachable)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("Session Not Found")
	}

	return nil
}

func sendConnectionAlertNotification(data mtrSessionParms, isReachable bool) {
	subject := fmt.Sprintf("[ALERT] Destination %s is UNREACHABLE ", data.session.DestinationIP)
	if isReachable {
		subject = fmt.Sprintf("[ALERT] Destination %s is reachable ", data.session.DestinationIP)
	}

	notification.UserNotificationChannel(mailer.EmailData{
		Subject: subject,
		BodyTemplate: email.TraceRouteEmail(email.TraceRouteEmailParms{
			Session: data.session,
			Result:  data.out,
		}, isReachable),
		FileAttachment: []string{},
	})
}

func checkReachableStatus(data mtrSessionParms) {
	db := data.db
	session := data.session
	config := data.config

	var reachableLogs []bool

	err := db.Model(&models.MTRResult{}).
		Select("reachable").
		Where("session_id = ? AND send_notification = ?", session.ID, true).
		Order("created_at DESC").
		Limit(config.RangeReachableLoss).
		Pluck("reachable", &reachableLogs).Error

	if err != nil {
		services.LogError(ServiceMTRSession, "Failed to fetch reachable logs: "+err.Error())
		return
	}

	reachableStatus := checkReachable(reachableLogs)

	if reachableStatus == DOWN && session.IsReachable {
		// Host sekarang DOWN, sebelumnya reachable → warning
		services.LogWarning(ServiceMTRSession, fmt.Sprintf(
			"MTR Session Alert: Destination %s (Session ID: %s) is currently UNREACHABLE. Last successful check was at %s.",
			session.DestinationIP,
			session.ID,
			session.LastRunAt.Format("02/01/2006 15:04:05"),
		))
		updateStatusReachableSession(session, false)
		sendConnectionAlertNotification(data, false)
	} else if reachableStatus == UP && !session.IsReachable {
		// Host sekarang UP, sebelumnya unreachable → info
		services.LogInfo(ServiceMTRSession, fmt.Sprintf(
			"MTR Session Notice: Destination %s (Session ID: %s) is now reachable. Recovery detected at %s.",
			session.DestinationIP,
			session.ID,
			time.Now().Format("02/01/2006 15:04:05"),
		))
		updateStatusReachableSession(session, true)
		sendConnectionAlertNotification(data, true)
	}
}

func sendPacketUseWebsocket(data mtrSessionParms) {
	msg := ws.Message{
		Time:    time.Now(),
		ID:      data.session.ID.String(),
		Message: data.out,
	}

	data.manager.BroadcastProject(msg)
}

func startSyncSessionMTRWorker(db *gorm.DB, config MTRSessionConfig, manager *worker.Manager) {
	var sessions []models.MTRSession
	result := db.Find(&sessions).Where("status = ?", "active")
	if result.Error != nil {
		services.LogError(ServiceMTRSession, "Failed to fetch MTR sessions: "+result.Error.Error())
		return
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
			params := mtrSessionParms{
				db:      db,
				session: s,
				config:  config,
				out:     *out,
				manager: manager,
			}

			err = saveMtrResult(params)
			if err != nil {
				services.LogError(ServiceMTRSession, "Failed to save MTR result for session ID "+s.ID.String()+": "+err.Error())
				return
			}
			
			err = updateMtrSessionLastRun(params)
			if err != nil {
				services.LogError(ServiceMTRSession, "Failed to update last run for session ID "+s.ID.String()+": "+err.Error())
				return
			}
			checkReachableStatus(params)
			sendPacketUseWebsocket(params)
		}(session)
	}
}

func RunMTRSession(manager *worker.Manager) (*worker.Worker, error) {
	db := database.DB

	var config MTRSessionConfig

	service, err := services.GetByServiceName(ServiceMTRSession)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			services.LogInfo(ServiceMTRSession, "service mtr-session not found, worker disabled")
			return nil, nil
		}
		return nil, err
	}

	err = json.Unmarshal(service.Config, &config)
	if err != nil {
		return nil, err
	}

	w := worker.NewWorker(
		ServiceMTRSession,
		config.Cron,
		func() {
			services.LogInfo(ServiceMTRSession, "Starting MTR Session Session Worker")

			startSyncSessionMTRWorker(db, config, manager)
		},
	)

	return w, nil
}
