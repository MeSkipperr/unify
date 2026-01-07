package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"unify-backend/cmd"
	"unify-backend/internal/worker"
	"unify-backend/internal/ws"
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

func NewHandler(m *worker.Manager) http.Handler {
	h := &Handler{
		manager: m,
		wsHub:   ws.NewHub(),
	}

	mux := http.NewServeMux()

	// existing http
	mux.HandleFunc("/services/", h.router)

	// NEW: websocket
	mux.Handle("/ws/services", ws.ServeWS(h.wsHub))

	return mux
}


func (h *Handler) router(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/status") {
		h.status(w, r)
		return
	}
	http.NotFound(w, r)
}

func (h *Handler) status(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	service := parts[2]

	switch r.Method {

	case http.MethodGet:
		status, ok := h.manager.Status(service)
		if !ok {
			http.Error(w, "service not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(StatusResponse{
			Service: service,
			Status:  status,
			At:      time.Now(),
		})

	case http.MethodPut:
		var payload StatusPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		switch payload.Status {
		case worker.StatusStarted, worker.StatusStopped:
			if err := h.manager.SetStatus(service, payload.Status); err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

		case worker.StatusRestart:
			if service == "monitoring-network" {
				if err := RestartMonitoringNetwork(h.manager); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				payload.Status = worker.StatusStarted
			} else {
				http.Error(w, "restart not supported for this service", http.StatusBadRequest)
				return
			}

		default:
			http.Error(w, "invalid status value", http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(StatusResponse{
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

	old, _ := manager.Get("monitoring-network")
	old.Stop()

	manager.Replace("monitoring-network", w)
	w.Start()
	manager.SetStatus("monitoring-network", worker.StatusStarted)

	return nil
}
