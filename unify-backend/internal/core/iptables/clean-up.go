package iptables
import (	
	"log"
	"unify-backend/models"
	"gorm.io/gorm"
)

func CleanupAllPortForwardRules(db *gorm.DB) error {
	var sessions []models.SessionPortForward

	if err := db.Find(&sessions).Error; err != nil {
		return err
	}

	for _, s := range sessions {
		// paksa delete, abaikan status
		if err := DeleteRule(s); err != nil {
			log.Println("[startup-cleanup] delete failed:", err)
		}
	}

	return nil
}
