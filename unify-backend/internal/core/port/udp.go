package port

import (
	"net"
	"strconv"
	"time"
)

func checkUDP(p Params) Result {
	if p.Timeout <= 0 {
		p.Timeout = 3 * time.Second
	}

	addr := net.JoinHostPort(p.Target, strconv.Itoa(p.Port))

	conn, err := net.DialTimeout("udp", addr, p.Timeout)
	if err != nil {
		return Result{
			Target:   p.Target,
			Port:     p.Port,
			Protocol: UDP,
			Open:     false,
			Error:    err.Error(),
		}
	}
	defer conn.Close()

	// Send empty packet
	_, err = conn.Write([]byte{})
	if err != nil {
		return Result{
			Target:   p.Target,
			Port:     p.Port,
			Protocol: UDP,
			Open:     false,
			Error:    err.Error(),
		}
	}

	// Try to read (optional)
	_ = conn.SetReadDeadline(time.Now().Add(p.Timeout))
	buf := make([]byte, 1)
	_, err = conn.Read(buf)

	// UDP behavior:
	// - no error → maybe open
	// - timeout → open/filtered
	// - ICMP unreachable → closed
	if err == nil {
		return Result{
			Target:   p.Target,
			Port:     p.Port,
			Protocol: UDP,
			Open:     true,
		}
	}

	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return Result{
			Target:   p.Target,
			Port:     p.Port,
			Protocol: UDP,
			Open:     true, // open or filtered
		}
	}

	return Result{
		Target:   p.Target,
		Port:     p.Port,
		Protocol: UDP,
		Open:     false,
		Error:    err.Error(),
	}
}
