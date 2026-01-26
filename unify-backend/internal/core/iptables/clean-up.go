package iptables

import (
	"log"
	"os/exec"
	"unify-backend/models"

	"gorm.io/gorm"
)

func CleanupAllPortForwardRulesLinux(chain string) error {
	return exec.Command("sudo", "iptables", "-t", "nat", "-F", chain).Run()
}

func CleanupAllPortForwardRulesDB(db *gorm.DB) error {
	var sessions []models.SessionPortForward

	if err := db.
		Where(`
		(status = ?)
		OR (status = ? AND expires_at < NOW())
		OR (status = ?)
	`,
			models.SessionStatusActive,
			models.SessionStatusPending,
			models.SessionStatusDeactivated,
		).
		Find(&sessions).Error; err != nil {
		return err
	}

	for _, s := range sessions {
		switch s.Status {
		case models.SessionStatusActive:
			s.Status = models.SessionStatusExpired

		default:
			// deactivated & pending
			s.Status = models.SessionStatusInactive
		}

		if err := db.Save(&s).Error; err != nil {
			return err
		}
	}

	return nil
}

func CleanupAllPortForwardRules(db *gorm.DB, chain string) error {
	err := CleanupAllPortForwardRulesLinux(chain)

	if err != nil {
		log.Println("failed to clean up iptables rules:", err)
		return err
	}

	err = CleanupAllPortForwardRulesDB(db)
	if err != nil {
		log.Println("failed to clean up database sessions:", err)
		return err
	}

	return nil
}
