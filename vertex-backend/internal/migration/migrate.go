package migration

import (
	"vertex-backend/internal/entity"

	"gorm.io/gorm"
)

// Run performs the database migration for the specified entities.
func Run(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.Devices{},
		&entity.DeviceService{},
		&entity.SelectOption{},
		&entity.ADBJob{},
		&entity.ADBJobTarget{},
		&entity.ADBJobResult{},
		&entity.SpeedTestJob{},
		&entity.SpeedTestResult{},
		&entity.MTRSession{},
		&entity.MTRResult{},
		&entity.MTRHop{},
		&entity.User{},
		&entity.Log{},
		&entity.Notification{},
	)
}