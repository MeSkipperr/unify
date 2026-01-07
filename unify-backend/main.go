// package main

// import (
// 	"log"
// 	"net/http"

// 	api "unify-backend/internal/http"
// 	"unify-backend/internal/services"
// 	"unify-backend/internal/worker"
// 	"unify-backend/internal/ws"
// )

// func main() {
// 	manager := worker.NewManager()

// 	projectHub := ws.NewHub()
// 	manager.SetProjectHub(projectHub)

// 	w, err := services.MonitoringNetwork(manager)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	manager.Register(w)

// 	mux := http.NewServeMux()
// 	apiHandler := api.NewHandler(manager)
// 	mux.Handle("/", apiHandler)

// 	// WebSocket endpoint
// 	mux.Handle("/ws/project", ws.ServeWS(projectHub))

// 	// 7️⃣ Start server
// 	log.Println("server running on :8080")
// 	if err := http.ListenAndServe(":8080", mux); err != nil {
// 		log.Fatal(err)
// 	}
// }

package main

import (
	"fmt"
	"log"
	"net/http"
	"unify-backend/internal/database"
	"unify-backend/internal/repository"
	"unify-backend/internal/services"
	"unify-backend/models"

	_ "github.com/lib/pq"
)

func main() {
	database.Connect()
	database.DB.AutoMigrate(
		&models.Service{}, // WAJIB PERTAMA
		&models.Log{},     // BARU LOG
		// &models.Devices{},     // BARU LOG
	)
	seedData()

	logRepo := repository.NewLogRepository(database.DB)

	logService := services.NewLogService(logRepo)

	err := logService.CreateLog(services.CreateLogParams{
		Level:       "INFO",
		ServiceName: "DeviceMonitor",
		Message:     "Service started successfully",
	})
	selectAllDevices()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Log created successfully")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Backend Go running on :8080")
	})

	log.Println("Backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func seedData() {
	device := models.Devices{
		Name: "cctv-04",
		IPAddress: "192.168.1.12",
		Type:      "cctv",
		// Status:    "active",
	}

	database.DB.FirstOrCreate(&device, models.Devices{Name: "cctv-04"})
	fmt.Println("✅ Seed data ready")
}

func selectAllDevices() {
	var devices []models.Devices

	result := database.DB.Find(&devices)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	fmt.Println("📦 Devices:")
	for _, d := range devices {
		fmt.Printf(
			"- ID=%d | DeviceID=%s | IP=%s | Type=%s 	\n",
			d.ID,
			d.Name,
			d.IPAddress,
			d.Type,
			// d.Status,
		)
	}
}
