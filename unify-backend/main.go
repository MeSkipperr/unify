package main

import (
	"log"
	"net/http"
	"time"
	"unify-backend/cmd"
	"unify-backend/config"
	"unify-backend/internal/database"
	api "unify-backend/internal/http"
	"unify-backend/internal/worker"
	"unify-backend/internal/ws"
	"unify-backend/models"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database.Connect()
	database.Migrate()

	manager := worker.NewManager()

	projectHub := ws.NewHub()
	manager.SetProjectHub(projectHub)

	errs := worker.RegisterWorkersContinue(manager, []worker.WorkerFactory{
		cmd.MonitoringNetwork,
		cmd.RemoveDataYoutubeADB,
		cmd.GetUptimeADB,
		cmd.GetSpeedtestNetwork,
		cmd.RunPortForwardSession,
	})

	for _, err := range errs {
		log.Println("worker error:", err)
	}
	session := models.SessionPortForward{
		ListenIP:   "172.19.186.63",
		ListenPort: 3000,
		DestIP:     "172.19.186.63",
		DestPort:   8000,
		Protocol:   "tcp",
		Status:     models.SessionStatusPending,

		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(5 * time.Minute), 
	}

	if err := database.DB.Create(&session).Error; err != nil {
		log.Println("failed create session:", err)
	}

	mux := http.NewServeMux()
	apiHandler := api.NewHandler(manager)
	mux.Handle("/", apiHandler)

	mux.Handle("/ws/project", ws.ServeWS(projectHub))

	log.Println("server running on port", config.ServerPort)
	if err := http.ListenAndServe(config.ServerPort, mux); err != nil {
		log.Fatal(err)
	}
}
