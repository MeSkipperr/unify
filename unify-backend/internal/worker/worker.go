package worker

import (
	"fmt"
	"sync"
	"unify-backend/internal/services"

	"github.com/robfig/cron/v3"
)

type Worker struct {
	Name    string
	Cron    string
	Task    func()
	status  Status
	cron    *cron.Cron
	mu      sync.Mutex
	RunOnce bool
}

func NewWorker(name, cronExpr string, task func()) *Worker {
	return &Worker{
		Name:   name,
		Cron:   cronExpr,
		Task:   task,
		status: StatusStopped,
	}
}

func (w *Worker) Start() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status == StatusStarted {
		return nil
	}

	err := services.CreateAppLog(services.CreateLogParams{
		Level:       "INFO",
		ServiceName: w.Name,
		Message:     "Services started",
	})
	if err != nil {
		return err
	}

	w.status = StatusStarted

	if w.RunOnce {
		go func() {
			w.Task()

			w.mu.Lock()
			w.status = StatusStopped
			w.mu.Unlock()
		}()
		return nil
	}

	if w.Cron == "" {
		w.status = StatusStopped
		return fmt.Errorf("cron expression is empty for worker %s", w.Name)
	}

	c := cron.New(cron.WithSeconds())
	_, err = c.AddFunc(w.Cron, w.Task)
	if err != nil {
		w.status = StatusStopped
		return err
	}

	w.cron = c
	w.cron.Start()
	return nil
}

func (w *Worker) Stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status == StatusStopped {
		return fmt.Errorf("worker %s already stopped", w.Name)
	}

	err := services.CreateAppLog(services.CreateLogParams{
		Level:       "INFO",
		ServiceName: w.Name,
		Message:     "Services stopped",
	})
	if err != nil {
		return err
	}

	if w.cron != nil {
		w.cron.Stop()
		w.cron = nil
	}

	w.status = StatusStopped
	return nil
}

func (w *Worker) Status() Status {
	return w.status
}

func (w *Worker) WithRunOnce(v bool) *Worker {
	w.RunOnce = v
	return w
}

type WorkerFactory func(*Manager) (*Worker, error)

func RegisterWorkersContinue(
    manager *Manager,
    factories []WorkerFactory,
) []error {
    var errs []error
    for _, factory := range factories {
        w, err := factory(manager)
        if err != nil {
            errs = append(errs, err)
            continue
        }
        if w != nil {
            manager.Register(w)
        }
    }
    return errs
}
