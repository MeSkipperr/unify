package worker

import (
	"sync"

	"github.com/robfig/cron/v3"
)

type Worker struct {
	Name   string
	Cron   string
	Task   func()
	status Status

	cron *cron.Cron
	mu   sync.Mutex
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

	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc(w.Cron, w.Task)
	if err != nil {
		return err
	}

	w.cron = c
	w.cron.Start()
	w.status = StatusStarted
	return nil
}

func (w *Worker) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.cron != nil {
		w.cron.Stop()
	}
	w.status = StatusStopped
}

func (w *Worker) Status() Status {
	return w.status
}
