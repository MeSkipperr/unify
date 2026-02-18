package config

import (
	"os"
	"strings"
)

const ServerPort = ":8080"


type Config struct {
	AllowedOrigins []string
}

func LoadConfig() *Config {
	origins := os.Getenv("CORS_ALLOWED_ORIGINS")

	// CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:80

	var allowedOrigins []string
	if origins != "" {
		allowedOrigins = strings.Split(origins, ",")
	}

	return &Config{
		AllowedOrigins: allowedOrigins,
	}
}
