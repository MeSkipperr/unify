package main

import (
	"log"
	"net/http"
	"unify-backend/cmd"
	"unify-backend/config"
	"unify-backend/internal/database"
	api "unify-backend/internal/http"
	"unify-backend/internal/worker"
	"unify-backend/internal/ws"

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


	mux := http.NewServeMux()
	apiHandler := api.NewHandler(manager)
	mux.Handle("/", apiHandler)

	mux.Handle("/ws/project", ws.ServeWS(projectHub))

	log.Println("server running on port", config.ServerPort)
	if err := http.ListenAndServe(config.ServerPort, mux); err != nil {
		log.Fatal(err)
	}
}
