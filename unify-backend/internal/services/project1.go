package services

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"unify-backend/internal/worker"
)

type ProjectCronConfig struct {
	Service string `json:"service"`
	Cron    string `json:"cron"`
}


func loadProject1Cron() (string, error) {
	data, err := os.ReadFile("project1.json")
	if err != nil {
		return "", err
	}

	var cfg ProjectCronConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", err
	}

	if cfg.Cron == "" {
		return "", errors.New("cron is empty")
	}

	return cfg.Cron, nil
}


func Project1Worker() (*worker.Worker, error) {
	cron, err := loadProject1Cron()
	if err != nil {
		return nil, err
	}

	return worker.NewWorker(
		"project1",
		cron,
		func() {
			log.Println("project1 task running")
		},
	), nil
}
