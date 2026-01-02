package worker

type Status string

const (
	StatusStarted Status = "RUNNING"
	StatusStopped Status = "STOPPED"
	StatusRestart Status = "RESTART"
)
