package cmd

import (
	// "encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	"unify-backend/internal/core/iptables"
	"unify-backend/internal/database"
	"unify-backend/internal/services"
	"unify-backend/internal/worker"
	"unify-backend/models"

	"gorm.io/gorm"
)

func SyncSessionPortForward(s *models.SessionPortForward) error {
	switch s.Status {

	case models.SessionStatusPending:
		if err := iptables.ApplyRule(s); err != nil {
			s.Status = models.SessionStatusError
			return err
		}
		log.Println("applied port forward rule for session ID", s.ID)
		s.Status = models.SessionStatusActive
		return nil

	case models.SessionStatusActive:
		return nil

	case models.SessionStatusExpired, models.SessionStatusDisabled:
		services.LogInfo(ServicePortForward,fmt.Sprintf("deleting port forward rule for session ID %d", s.ID))
		if err := iptables.DeleteRule(*s); err != nil {
			return err
		}
		return nil

	default:
		services.LogError(ServicePortForward,fmt.Sprintf("invalid session status: %s", s.Status))
		return fmt.Errorf("invalid session status: %s", s.Status)
	}
}

func ExpireWorkerPortForward(db *gorm.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			var sessions []models.SessionPortForward

			db.Where(
				"expires_at < NOW() AND status = ?",
				models.SessionStatusActive,
			).Find(&sessions)

			for _, s := range sessions {
				if err := iptables.DeleteRule(s); err != nil {
					continue
				}

				s.Status = models.SessionStatusExpired
				db.Save(&s)
			}
		}
	}()
}

func startSyncWorker(db *gorm.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		var sessions []models.SessionPortForward

		db.Where(
			"status = ?",
			models.SessionStatusPending,
		).Find(&sessions)

		for _, s := range sessions {
			if err := SyncSessionPortForward(&s); err != nil {
				s.Status = models.SessionStatusError
			}
			db.Save(&s)
		}
	}
}

func startExpireWorker(db *gorm.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		var sessions []models.SessionPortForward

		err := db.Where(
			"status = ? AND expires_at <= NOW()",
			models.SessionStatusActive,
		).Find(&sessions).Error
		if err != nil {
			log.Println("[expire-worker] failed query:", err)
			services.LogError(ServicePortForward, "Error Query Expire Port Forward Sessions : "+err.Error())
			continue
		}

		for _, s := range sessions {
			s.Status = models.SessionStatusExpired

			if err := SyncSessionPortForward(&s); err != nil {
				log.Println("[expire-worker] failed sync:", err)
				services.LogError(ServicePortForward, "Error Expire Port Forward Session ID "+fmt.Sprint(s.ID)+" : "+err.Error())
				s.Status = models.SessionStatusError
			}

			if err := db.Save(&s).Error; err != nil {
				log.Println("[expire-worker] failed save:", err)
				services.LogError(ServicePortForward, "Error Save Expire Port Forward Session ID "+fmt.Sprint(s.ID)+" : "+err.Error())
			}
		}
	}
}

type portForwardConfig struct {
	SyncInterval   int `json:"sync_interval"`
	ExpireInterval int `json:"expire_interval"`
}

func RunPortForwardSession(manager *worker.Manager) (*worker.Worker, error) {
	db := database.DB

	config := portForwardConfig{
		SyncInterval:   10,
		ExpireInterval: 30,
	}

	// service, err := services.GetByServiceName(ServiceMonitoringNetwork)
	// if err != nil {
	// 	return nil, err
	// }

	// if err := json.Unmarshal(service.Config, &config); err != nil {
	// 	return nil, err
	// }

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
			if err := iptables.CleanupAllPortForwardRules(db); err != nil {
				log.Println("cleanup failed:", err)
				services.LogError(ServicePortForward, "Error Clean Up ip tables : "+err.Error())
			}

			// 2. REBUILD
			if err := iptables.RebuildActivePortForward(db); err != nil {
				log.Println("rebuild failed:", err)
				services.LogError(ServicePortForward, "Error Rebuild ip tables : "+err.Error())
			}

			// 3. START WORKERS
			go startSyncWorker(db, time.Duration(config.SyncInterval)*time.Second)
			go startExpireWorker(db, time.Duration(config.ExpireInterval)*time.Second)
		},
	)

	// worker manager yang mengontrol lifecycle
	w.RunOnce = true
	return w, nil
}
