package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
	"unify-backend/cmd"
	"unify-backend/config"
	"unify-backend/internal/database"
	api "unify-backend/internal/http"
	"unify-backend/internal/http/sse"
	"unify-backend/internal/notification"
	"unify-backend/internal/queue"
	"unify-backend/internal/worker"
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
	StartAutoDeviceNotification()
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

func StartAutoDeviceNotification() {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			<-ticker.C

			isOnline := rand.Intn(2) == 0
			now := time.Now()

			var level models.NotificationLevel
			var subject string
			var detail string

			if isOnline {
				level = models.NoticationStatusInfo
				subject = fmt.Sprintf("[INFO] %s - UP", "DPSCY")
				detail = fmt.Sprintf(
					"Device %s is reachable at %s.",
					"DPSCY",
					now.Format("15:04:05"),
				)
			} else {
				level = models.NoticationStatusAlert
				subject = fmt.Sprintf("[ALERT] %s - DOWN", "DPSCY")
				detail = fmt.Sprintf(
					"Device %s is unreachable at %s.",
					"DPSCY",
					now.Format("15:04:05"),
				)
			}

			notificationPayload := models.Notification{
				Level:  level,
				Title:  subject,
				Detail: detail,
				URL:    fmt.Sprintf("/devices?search=%s","DPSCY"),
			}

			notification.SSENotification(notificationPayload)
		}
	}()
}
