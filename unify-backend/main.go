package main

import (
	"log"
	"net/http"

	api "unify-backend/internal/http"
	"unify-backend/internal/services"
	"unify-backend/internal/worker"
	"unify-backend/internal/ws"
)

func main() {	
	manager := worker.NewManager()

	projectHub := ws.NewHub()
	manager.SetProjectHub(projectHub)

	w, err := services.MonitoringNetwork(manager)
	if err != nil {
		log.Fatal(err)
	}

	manager.Register(w)

	mux := http.NewServeMux()
	apiHandler := api.NewHandler(manager)
	mux.Handle("/", apiHandler)

	// WebSocket endpoint
	mux.Handle("/ws/project", ws.ServeWS(projectHub))

	// 7️⃣ Start server
	log.Println("server running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
