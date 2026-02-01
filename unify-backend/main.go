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

	"github.com/gin-gonic/gin"
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

	router := gin.Default()
	apiHandler := api.NewHandler(manager)

	mtrSocket := ws.NewHub()
	manager.SetMTRhub(mtrSocket)

	router.GET("/ws/mtr", func(c *gin.Context) {
		ws.ServeWS(mtrSocket).ServeHTTP(c.Writer, c.Request)
	})

	server := &http.Server{
		Addr:    config.ServerPort,
		Handler: apiHandler,
	}

	log.Println("Server running at: ", config.ServerPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
