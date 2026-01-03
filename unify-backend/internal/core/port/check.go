package port

import "time"

type Protocol string

const (
	TCP Protocol = "tcp"
	UDP Protocol = "udp"
)

type Params struct {
	Target   string
	Port     int
	Protocol Protocol
	Timeout  time.Duration
}

type Result struct {
	Target   string   `json:"target"`
	Port     int      `json:"port"`
	Protocol Protocol `json:"protocol"`
	Open     bool     `json:"open"`
	Error    string   `json:"error,omitempty"`
}

// Check selects TCP or UDP check
func Check(p Params) Result {
	switch p.Protocol {
	case UDP:
		return checkUDP(p)
	default:
		return checkTCP(p)
	}
}
