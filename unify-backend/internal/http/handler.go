package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"unify-backend/cmd"
	"unify-backend/internal/http/handler"
	"unify-backend/internal/http/middleware"
	"unify-backend/internal/http/sse"
	"unify-backend/internal/services"
	"unify-backend/internal/worker"
	"unify-backend/internal/ws"

	"github.com/gin-contrib/cors"
)

type Handler struct {
	manager *worker.Manager
	wsHub   *ws.Hub
	sse     *sse.SSEManager
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
	sseManager := m.GetSSE()
	if sseManager == nil {
		sseManager = sse.NewSSEManager()
	}
	h := &Handler{
		manager: m,
		wsHub:   ws.NewHub(),
		sse:     sseManager,
	}

	router := gin.Default()
	router.RemoveExtraSlash = true
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5500",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "PATCH",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"X-Timezone", 
		},
		AllowCredentials: true,
	}))

	router.Use(middleware.LoggerMiddleware("backend-route"))

	auth := router.Group("/auth")
	{
		auth.POST("/login", services.LoginHandler)
		auth.POST("/refresh", services.RefreshTokenHandler)
		auth.POST("/me", services.Me)
	}

	events := router.Group("/events")
	{
		events.GET("/services", func(c *gin.Context) {
			h.sse.Subscribe(c.Writer, c.Request, sse.SSEChannelServices)
		})
		events.GET("/notification", func(c *gin.Context) {
			h.sse.Subscribe(c.Writer, c.Request, sse.SSEChannelNotif)
		})
	}

	api := router.Group("/api", middleware.AuthMiddleware)
	{
		// =========================
		// DEVICES
		// =========================
		devices := api.Group("/devices")
		{
			devices.GET("", handler.GetDevices)
			devices.POST("", handler.CreateDevice())
			devices.PUT("/:id", handler.ChangeDevice())
			devices.DELETE("/:id", handler.DeleteDevice())
			devices.PATCH("/:id/notification", handler.ChangeNotification())
		}

		// =========================
		// SERVICES
		// =========================
		services := api.Group("/services")
		{
			services.GET("", handler.GetServices())
			services.GET("/:serviceName", handler.GetServiceByName())

			// ----- ADB -----
			adb := services.Group("/adb")
			{
				adb.GET("", handler.GetAdbResults)
				adb.POST("", handler.CreateADBJob())
			}

			// ----- SPEEDTEST -----
			speedtest := services.Group("/speedtest")
			{
				speedtest.GET("", handler.GetSpeedtestByInternalIPAndServer())
			}

			// ----- MTR -----
			mtr := services.Group("/mtr-sessions")
			{
				mtr.GET("/active", handler.GetActiveMTRSessions())
				mtr.POST("", handler.CreateMTRSession())
				mtr.PUT("", handler.DisableMTRSession())
				mtr.GET("/result/:id", handler.GetMTRResult())
			}

			// ----- PORT FORWARD -----
			portForward := services.Group("/port-forward")
			{
				portForward.GET("", handler.GetPortForward())
				portForward.POST("", handler.CreatePortForward())
				portForward.PATCH("/:id/deactivate", handler.DeactivatedPortForward())
			}
		}
	}

	// Existing HTTP endpoints
	router.GET("/services/:service/status", h.getStatus)
	router.PUT("/services/:service/status", h.updateStatus)

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
		At:      time.Now().UTC(),
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
		At:      time.Now().UTC(),
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
