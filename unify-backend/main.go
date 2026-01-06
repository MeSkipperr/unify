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
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	db := connectDB()
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Backend Go running on :8080")
	})

	log.Println("Backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func connectDB() *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("DB connection failed:", err)
	}

	log.Println("PostgreSQL connected")
	return db
}
