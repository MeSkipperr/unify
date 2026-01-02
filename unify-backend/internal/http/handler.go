package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"unify-backend/internal/services"
	"unify-backend/internal/worker"
)

type Handler struct {
	manager *worker.Manager
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
	h := &Handler{manager: m}
	mux := http.NewServeMux()
	mux.HandleFunc("/services/", h.router)
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
			if service == "project1" {
				if err := RestartProject1(h.manager); err != nil {
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

func RestartProject1(manager *worker.Manager) error {
	w, err := services.Project1Worker() // baca config terbaru
	if err != nil {
		return err
	}

	old, _ := manager.Get("project1")
	old.Stop()

	manager.Replace("project1", w)
	w.Start()
	manager.SetStatus("project1", worker.StatusStarted)

	return nil
}
