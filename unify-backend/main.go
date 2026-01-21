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
	"log"
	"net/http"
	"unify-backend/cmd"
	"unify-backend/internal/database"
	api "unify-backend/internal/http"
	"unify-backend/internal/worker"
	"unify-backend/internal/ws"

	_ "github.com/lib/pq"
)

func main() {
	database.Connect()
	database.Migrate()

	manager := worker.NewManager()

	projectHub := ws.NewHub()
	manager.SetProjectHub(projectHub)

	errs := worker.RegisterWorkersContinue(manager, []worker.WorkerFactory{
		cmd.MonitoringNetwork,
		cmd.RemoveDataYoutubeADB,
		cmd.GetUptimeADB,
	})

	for _, err := range errs {
		log.Println("worker error:", err)
	}

	mux := http.NewServeMux()
	apiHandler := api.NewHandler(manager)
	mux.Handle("/", apiHandler)

	mux.Handle("/ws/project", ws.ServeWS(projectHub))

	log.Println("server running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
