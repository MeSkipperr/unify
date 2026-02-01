package http

import (
	"fmt"
	"net/http"
	"time"

	"unify-backend/cmd"
	"unify-backend/internal/worker"
	"unify-backend/internal/ws"

	"github.com/gin-gonic/gin"
)

type GinHandler struct {
	manager *worker.Manager
	wsHub   *ws.Hub
}

func NewGinHandler(m *worker.Manager) *GinHandler {
	return &GinHandler{
		manager: m,
		wsHub:   ws.NewHub(),
	}
}

// Register routes ke router Gin
func (h *GinHandler) RegisterRoutes(router *gin.Engine) {
	fmt.Print()
	// Service routes
	router.GET("/services/:service/status", h.status)
	router.PUT("/services/:service/status", h.status)

	// WebSocket route
	router.GET("/ws/services", func(c *gin.Context) {
		conn, err := ws.Upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		h.wsHub.Register(conn)

		go func() {
			defer h.wsHub.Unregister(conn)
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					break
				}
			}
		}()
	})
}

// ========================
// Handler logic tetap sama
// ========================

type StatusPayload struct {
	Status worker.Status `json:"status"`
}

type StatusResponse struct {
	Service string        `json:"service"`
	Status  worker.Status `json:"status"`
	At      time.Time     `json:"at"`
}

func (h *GinHandler) status(c *gin.Context) {
	service := c.Param("service")
	if service == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service"})
		return
	}

	switch c.Request.Method {
	case http.MethodGet:
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

	case http.MethodPut:
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
				err := h.manager.Restart(service)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "restart not supported for this service"})
					return
				}
				if err := h.manager.SetStatus(service, worker.StatusStarted); err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
					return
				}
			}

			payload.Status = worker.StatusRestart

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
