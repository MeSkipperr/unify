package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
	"unify-backend/internal/core/iptables"
	"unify-backend/internal/database"
	"unify-backend/internal/notification"
	"unify-backend/internal/services"
	"unify-backend/internal/worker"
	"unify-backend/models"

	"gorm.io/gorm"
)

func handleExpiry(s *models.SessionPortForward, now time.Time) {
	if s.Status == models.SessionStatusDeactivated {
		s.Status = models.SessionStatusInactive
		return
	}
	if !s.ExpiresAt.Before(now) {
		return
	}
	switch s.Status {
	case models.SessionStatusActive:
		s.Status = models.SessionStatusExpired

	case models.SessionStatusPending:
		s.Status = models.SessionStatusInactive
	}
}

func handleSessionState(db *gorm.DB, s *models.SessionPortForward) {
	switch s.Status {

	case models.SessionStatusPending:
		handlePending(db, s)

	case models.SessionStatusExpired,
		models.SessionStatusInactive:
		handleExpired(db, s)

	case models.SessionStatusDeactivated:
		db.Save(s)

	case models.SessionStatusActive:
		return
	}
}

func sendRuleApplied(level models.NotificationLevel, data *models.SessionPortForward) {
	sseManager := worker.ManagerGlobal.GetSSE()
	if sseManager == nil || data == nil {
		return
	}

	listenAddr := fmt.Sprintf("%s:%d", data.ListenIP, data.ListenPort)
	destAddr := fmt.Sprintf("%s:%d", data.DestIP, data.DestPort)

	var title string
	var detail string

	switch level {
	case models.NoticationStatusInfo:
		title = "Port Forward Rule Applied"
		detail = fmt.Sprintf(
			"Port forwarding rule successfully applied: %s → %s (%s).",
			listenAddr,
			destAddr,
			data.Protocol,
		)

	case models.NoticationStatusError:
		title = "Port Forward Rule Failed"
		detail = fmt.Sprintf(
			"Failed to apply port forwarding rule: %s → %s (%s).",
			listenAddr,
			destAddr,
			data.Protocol,
		)

	default:
		title = "Port Forward Rule Update"
		detail = fmt.Sprintf(
			"Port forwarding rule update: %s → %s (%s).",
			listenAddr,
			destAddr,
			data.Protocol,
		)
	}

	notificationPayload := models.Notification{
		Level:     level,
		Title:     title,
		Detail:    detail,
		URL:       listenAddr,
		CreatedAt: data.CreatedAt,
	}

	notification.SSENotification(notificationPayload)
}

func handlePending(db *gorm.DB, s *models.SessionPortForward) {
	if err := iptables.ApplyRule(s); err != nil {
		s.Status = models.SessionStatusError
		sendRuleApplied(models.NoticationStatusError, s)
		services.LogError(
			ServicePortForward,
			fmt.Sprintf("failed to apply rule for session %s: %v", s.ID, err),
		)
		db.Save(s)
		return
	}

	s.Status = models.SessionStatusActive
	db.Save(s)

	sendRuleApplied(models.NoticationStatusInfo, s)

	services.LogInfo(
		ServicePortForward,
		fmt.Sprintf("applied port forward rule for session %s", s.ID),
	)
}

func handleExpired(db *gorm.DB, s *models.SessionPortForward) {
	if err := iptables.DeleteRule(*s); err != nil {
		s.Status = models.SessionStatusError
		services.LogError(
			ServicePortForward,
			fmt.Sprintf("failed to delete rule for session %s: %v", s.ID, err),
		)
		db.Save(s)
		return
	}

	db.Save(s)

	services.LogInfo(
		ServicePortForward,
		fmt.Sprintf("deleted expired rule for session %s", s.ID),
	)
}

func startSyncSessionPortForwardWorker(db *gorm.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		var sessions []models.SessionPortForward
		now := time.Now().UTC()

		if err := db.Where(
			"status IN ?",
			[]models.SessionStatus{
				models.SessionStatusPending,
				models.SessionStatusActive,
				models.SessionStatusDeactivated,
			},
		).Find(&sessions).Error; err != nil {
			continue
		}

		for _, s := range sessions {
			handleExpiry(&s, now)
			handleSessionState(db, &s)
		}
	}
}

type portForwardConfig struct {
	SyncInterval int `json:"sync_interval"`
}

func RunPortForwardSession(manager *worker.Manager) (*worker.Worker, error) {
	db := database.DB

	var config portForwardConfig

	service, err := services.GetByServiceName(ServicePortForward)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			services.LogInfo(ServicePortForward, "service port-forward not found, worker disabled")
			return nil, nil
		}
		return nil, err
	}

	err = json.Unmarshal(service.Config, &config)
	if err != nil {
		return nil, err
	}

	chain := os.Getenv("IPTABLES_NAT_CHAIN_PORT_FORWARD")
	if chain == "" {
		return nil, errors.New("IPTABLES_NAT_CHAIN_PORT_FORWARD is not set")
	}

	if err := iptables.EnsureChain(chain); err != nil {
		return nil, err
	}

	w := worker.NewWorker(
		ServicePortForward,
		"",
		func() {
			services.LogInfo(ServicePortForward, "Starting Port Forward Session Worker")

			// 1. CLEAN
			if err := iptables.CleanupAllPortForwardRules(db, chain); err != nil {
				services.LogError(ServicePortForward, "Error Clean Up ip tables : "+err.Error())
			}

			// 2. REBUILD
			if err := iptables.RebuildActivePortForward(db); err != nil {
				services.LogError(ServicePortForward, "Error Rebuild ip tables : "+err.Error())
			}

			// 3. START WORKERS
			go startSyncSessionPortForwardWorker(db, time.Duration(config.SyncInterval)*time.Second)
		},
	)

	// worker manager yang mengontrol lifecycle
	w.RunOnce = true
	return w, nil
}
