package main

import (
	"log"
	"net/http"

	api "unify-backend/internal/http"
	"unify-backend/internal/services"
	"unify-backend/internal/worker"
)

func main() {
	manager := worker.NewManager()
	w, err := services.Project1Worker()
	if err != nil {
		log.Fatal(err)
	}

	manager.Register(w)

	handler := api.NewHandler(manager)

	log.Println("server running on :8080")
	http.ListenAndServe(":8080", handler)
}
