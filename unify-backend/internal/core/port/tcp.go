package port

import (
	"net"
	"strconv"
	"time"
)

func checkTCP(p Params) Result {
	if p.Timeout <= 0 {
		p.Timeout = 3 * time.Second
	}

	addr := net.JoinHostPort(p.Target, strconv.Itoa(p.Port))

	conn, err := net.DialTimeout("tcp", addr, p.Timeout)

	if err != nil {
		return Result{
			Target:   p.Target,
			Port:     p.Port,
			Protocol: TCP,
			Open:     false,
			Error:    err.Error(),
		}
	}
	defer conn.Close()

	return Result{
		Target:   p.Target,
		Port:     p.Port,
		Protocol: TCP,
		Open:     true,
	}
}
