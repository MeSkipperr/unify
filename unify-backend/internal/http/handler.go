package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"unify-backend/cmd"
	"unify-backend/internal/http/handler"
	"unify-backend/internal/http/middleware"
	"unify-backend/internal/services"
	"unify-backend/internal/worker"
	"unify-backend/internal/ws"

	"github.com/gin-contrib/cors"
)

type Handler struct {
	manager *worker.Manager
	wsHub   *ws.Hub
}

type StatusPayload struct {
	Status worker.Status `json:"status"`
}

type StatusResponse struct {
	Service string        `json:"service"`
	Status  worker.Status `json:"status"`
	At      time.Time     `json:"at"`
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewHandler(m *worker.Manager) *gin.Engine {
	h := &Handler{
		manager: m,
		wsHub:   ws.NewHub(),
	}

	router := gin.Default()
	router.RemoveExtraSlash = true
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE","PATCH"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	router.Use(middleware.LoggerMiddleware("backend-route"))

	auth := router.Group("/auth")
	{
		auth.POST("/login", services.LoginHandler)
		auth.POST("/refresh", services.RefreshTokenHandler)
		auth.POST("/me", services.Me)
	}

	api := router.Group("/api", middleware.AuthMiddleware)

	{
		api.POST("/users", func(c *gin.Context) {
			var newUser User
			if err := c.ShouldBindJSON(&newUser); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload"})
				return
			}
			fmt.Println("Payload struct:", newUser.Username)
			fmt.Println("Payload struct:", newUser.Password)

			c.JSON(http.StatusOK, gin.H{
				"message": "user received",
				"user":    newUser,
			})
		})
		api.GET("/devices", handler.GetDevices)
		api.PATCH("/devices/:id/notification", handler.ChangeNotification())
		api.POST("/devices", handler.CreateDevice())
		api.DELETE("/devices/:id", handler.DeleteDevice())
		api.PUT("/devices/:id", handler.ChangeDevice())


		api.GET("/services",handler.GetServices())
		api.GET("/services/:serviceName",handler.GetServiceByName())
		api.GET("/services/adb",handler.GetAdbResults)

		api.GET("/services/speedtest", handler.GetSpeedtestByInternalIPAndServer())
	}

	// Existing HTTP endpoints
	router.GET("/services/:service/status", h.getStatus)
	router.PUT("/services/:service/status", h.updateStatus)

	// Websocket
	router.GET("/ws/services", func(c *gin.Context) {
		ws.ServeWS(h.wsHub).ServeHTTP(c.Writer, c.Request)
	})

	return router
}

// GET /services/:service/status
func (h *Handler) getStatus(c *gin.Context) {
	service := c.Param("service")
	status, ok := h.manager.Status(service)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}

	c.JSON(http.StatusOK, StatusResponse{
		Service: service,
		Status:  status,
		At:      time.Now(),
	})
}

// PUT /services/:service/status
func (h *Handler) updateStatus(c *gin.Context) {
	service := c.Param("service")
	var payload StatusPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	switch payload.Status {
	case worker.StatusStarted, worker.StatusStopped:
		if err := h.manager.SetStatus(service, payload.Status); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
	case worker.StatusRestart:
		if service == cmd.ServiceMonitoringNetwork {
			if err := RestartMonitoringNetwork(h.manager); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			payload.Status = worker.StatusStarted
		} else {
			if err := h.manager.Restart(service); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "restart not supported for this service"})
				return
			}
			if err := h.manager.SetStatus(service, worker.StatusStarted); err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			payload.Status = worker.StatusRestart
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status value"})
		return
	}

	c.JSON(http.StatusOK, StatusResponse{
		Service: service,
		Status:  payload.Status,
		At:      time.Now(),
	})
}

func RestartMonitoringNetwork(manager *worker.Manager) error {
	w, err := cmd.MonitoringNetwork(manager)
	if err != nil {
		return err
	}

	old, _ := manager.Get(cmd.ServiceMonitoringNetwork)
	old.Stop()

	manager.Replace(cmd.ServiceMonitoringNetwork, w)
	w.Start()
	manager.SetStatus(cmd.ServiceMonitoringNetwork, worker.StatusStarted)

	return nil
}
