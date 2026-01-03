package worker

import (
	"errors"
	"sync"
	"unify-backend/internal/ws"
)

type Manager struct {
	mu      sync.RWMutex
	workers map[string]*Worker
	status  map[string]Status
	projectHub   *ws.Hub
}

func NewManager() *Manager {
	return &Manager{
		workers: make(map[string]*Worker),
		status:  make(map[string]Status),
	}
}

func (m *Manager) Register(w *Worker) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.workers[w.Name] = w
}

func (m *Manager) Status(name string) (Status, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	w, ok := m.workers[name]
	if !ok {
		return "", false
	}
	return w.Status(), true
}

func (m *Manager) SetStatus(name string, status Status) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	w, ok := m.workers[name]
	if !ok {
		return errors.New("service not found")
	}

	switch status {
	case StatusStarted:
		return w.Start()
	case StatusStopped:
		w.Stop()
	}
	return nil

}
func (m *Manager) Replace(name string, w *Worker) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.workers[name] = w
}
func (m *Manager) Get(name string) (*Worker, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	w, ok := m.workers[name]
	return w, ok
}

func (m *Manager) Restart(name string) error {
	m.mu.Lock()
	w, ok := m.workers[name]
	if !ok {
		m.mu.Unlock()
		return errors.New("service not found")
	}
	m.mu.Unlock()

	w.Stop()
	w.Start()

	m.mu.Lock()
	m.status[name] = StatusStarted
	m.mu.Unlock()

	return nil
}


func (m *Manager) BroadcastProject(msg ws.Message) {
	if m.projectHub == nil {
		return
	}
	m.projectHub.Broadcast(msg)
}


func (m *Manager) SetProjectHub(h *ws.Hub) {
	m.projectHub = h
}