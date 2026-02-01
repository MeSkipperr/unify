package main

import (
	"log"
	"os"
	"unify-backend/cmd"
	"unify-backend/config"
	"unify-backend/internal/database"
	api "unify-backend/internal/http"
	"unify-backend/internal/http/middleware"
	"unify-backend/internal/repository"
	"unify-backend/internal/services"
	"unify-backend/internal/worker"
	"unify-backend/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect DB & migrate
	database.Connect()
	database.Migrate()

	// Setup Gin
	router := gin.Default()
	manager := worker.NewManager()
	handler := api.NewGinHandler(manager) // gunakan constructor
	handler.RegisterRoutes(router)        // daftarkan routes
	// Setup workers
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

	// Setup repository & service untuk auth
	userRepo := repository.NewUserRepository(database.DB)
	authService := services.NewAuthService(userRepo, os.Getenv("JWT_SECRET"))
	authHandler := api.NewAuthHandler(authService)

	// ===== Public Routes =====
	public := router.Group("/api")
	{
		public.POST("/login", authHandler.Login)
	}

	// ===== Protected Routes =====
	protected := router.Group("/api")
	protected.Use(middleware.GinJWTMiddleware(os.Getenv("JWT_SECRET")))
	{
		protected.GET("/user/profile", func(c *gin.Context) {
			// Contoh ambil user ID dari JWT claims
			claims := c.MustGet("claims").(map[string]interface{})
			c.JSON(200, api.SuccessResponse(claims, "Authenticated"))
		})
	}

	// ===== WebSocket Routes =====
	mtrSocket := ws.NewHub()
	manager.SetMTRhub(mtrSocket)

	router.GET("/ws/mtr", gin.WrapH(ws.ServeWS(mtrSocket)))

	// Run server
	log.Println("server running on port", config.ServerPort)
	if err := router.Run(config.ServerPort); err != nil {
		log.Fatal(err)
	}
}
