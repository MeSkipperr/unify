package main

import (
	"log"
	"net/http"
	"unify-backend/cmd"
	"unify-backend/config"
	"unify-backend/internal/database"
	api "unify-backend/internal/http"
	"unify-backend/internal/http/sse"
	"unify-backend/internal/queue"
	"unify-backend/internal/worker"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database.Connect()
	database.Migrate()

	manager := worker.NewManager()
	sseManager := sse.NewSSEManager()
	manager.SetSSE(sseManager)

	queue.InitADBQueue(16)
	worker.StartADBWorkerPool(manager, 1, queue.GetADBQueue())

	errs := worker.RegisterWorkersContinue(manager, []worker.WorkerFactory{
		cmd.MonitoringNetwork,
		cmd.RemoveDataYoutubeADB,
		cmd.GetUptimeADB,
		cmd.GetSpeedtestNetwork,
		cmd.RunPortForwardSession,
		cmd.RunMTRSession,
	})
	for _, err := range errs {
		log.Println("worker error:", err)
	}

	apiHandler := api.NewHandler(manager)

	server := &http.Server{
		Addr:    config.ServerPort, 
		Handler: apiHandler,       
	}

	log.Println("Server running at: ", config.ServerPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
