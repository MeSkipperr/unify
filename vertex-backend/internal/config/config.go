package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresName     string
	PostgresSSLMode  string
}

// Load reads the configuration from environment variables and returns a Config struct.
func Load() *Config {

	err := godotenv.Load()

	if err != nil {
		log.Println(".env not found")
	}

	return &Config{
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresName:     os.Getenv("POSTGRES_DB"),
		PostgresSSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}
}