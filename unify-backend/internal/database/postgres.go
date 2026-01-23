package database

import (
	"fmt"
	"log"
	"os"
	"unify-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s ",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
		// os.Getenv("POSTGRES_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	DB = db
}

func Migrate() {
	err := DB.AutoMigrate(&models.AdbResult{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}
	err = DB.AutoMigrate(&models.Service{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}
	err = DB.AutoMigrate(&models.Devices{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}
	err = DB.AutoMigrate(&models.SpeedtestResult{})
	if err != nil {
		log.Fatal("failed to migrate SpeedtestResult:", err)
	}
	err = DB.AutoMigrate(&models.SessionPortForward{})
	if err != nil {
		log.Fatal("failed to migrate SessionPortForward:", err)
	}
}
