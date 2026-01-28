package iptables

import (
	"log"
	"unify-backend/models"
	"gorm.io/gorm"
)

func RebuildActivePortForward(db *gorm.DB) error {
	var sessions []models.SessionPortForward

	if err := db.Where(
		"status = ?",
		models.SessionStatusActive,
	).Find(&sessions).Error; err != nil {
		return err
	}

	for _, s := range sessions {
		if err := ApplyRule(&s); err != nil {
			log.Println("[startup-rebuild] apply failed:", err)
		}
	}

	return nil
}
