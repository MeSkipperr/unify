package main

import (
	"vertex-backend/internal/config"
	"vertex-backend/internal/database"
	"vertex-backend/internal/migration"
)

func main() {

	cfg := config.Load()

	db, err := database.Connect(cfg)

	if err != nil {
		panic(err)
	}

	err = migration.Run(db)

	if err != nil {
		panic(err)
	}
}