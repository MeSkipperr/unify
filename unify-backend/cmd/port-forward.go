package cmd

import (
	"encoding/json"
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
		s.Status = models.SessionStatusActive
		return nil

	case models.SessionStatusActive:
		return nil

	case models.SessionStatusExpired, models.SessionStatusDisabled:
		if err := iptables.DeleteRule(*s); err != nil {
			return err
		}
		return nil

	default:
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
			continue
		}

		for _, s := range sessions {
			s.Status = models.SessionStatusExpired

			if err := SyncSessionPortForward(&s); err != nil {
				log.Println("[expire-worker] failed sync:", err)
				s.Status = models.SessionStatusError
			}

			if err := db.Save(&s).Error; err != nil {
				log.Println("[expire-worker] failed save:", err)
			}
		}
	}
}


type portForwardConfig struct {
	SyncInterval   int `json:"sync_interval"`
	ExpireInterval int `json:"expire_interval"`
}

func RunPortForwardServices(manager *worker.Manager) (*worker.Worker, error) {
	db := database.DB

	config := portForwardConfig{
		SyncInterval:   10,
		ExpireInterval: 30,
	}

	service, err := services.GetByServiceName(ServiceMonitoringNetwork)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(service.Config, &config); err != nil {
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
	ServiceMonitoringNetwork,
	"",
	func() {

		// 1. CLEAN
		if err := iptables.CleanupAllPortForwardRules(db); err != nil {
			log.Println("cleanup failed:", err)
		}

		// 2. REBUILD
		if err := iptables.RebuildActivePortForward(db); err != nil {
			log.Println("rebuild failed:", err)
		}

		// 3. START WORKERS
		go startSyncWorker(db, time.Duration(config.SyncInterval)*time.Second)
		go startExpireWorker(db, time.Duration(config.ExpireInterval)*time.Second)
	},
)


	// worker manager yang mengontrol lifecycle
	w.RunOnce = false
	return w, nil
}
