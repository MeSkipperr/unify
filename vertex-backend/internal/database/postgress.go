package database

import (
	"fmt"
	"vertex-backend/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect establishes a connection to the PostgreSQL database using GORM.
func Connect(cfg *config.Config) (*gorm.DB, error) {

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresName,
		cfg.PostgresSSLMode,
	)

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}