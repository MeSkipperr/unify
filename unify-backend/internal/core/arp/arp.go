package arp

import (
	"bytes"
	"os/exec"
	"runtime"
	"strings"
	"unify-backend/internal/core/ping"
)

type Params struct {
	IP        string
	Interface string
	Warmup    bool
}

type Result struct {
	IP      string `json:"ip"`
	MAC     string `json:"mac,omitempty"`
	Exists  bool   `json:"exists"`
	RawLine string `json:"raw,omitempty"`
	Error   string `json:"error,omitempty"`
}

func Check(p Params) Result {
	if p.IP == "" {
		return Result{
			Error: "ip is required",
		}
	}

	// 🔥 Warmup ARP cache (ping once)
	if p.Warmup {
		ping.Ping(ping.Params{
			Target: p.IP,
			Times:  1,
		})

	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		if p.Interface != "" {
			// ip neigh show dev eth0 192.168.1.1
			cmd = exec.Command(
				"ip", "neigh", "show", "dev", p.Interface, p.IP,
			)
		} else {
			cmd = exec.Command("ip", "neigh", "show", p.IP)
		}

	case "windows":
		// interface selection not supported by arp -a
		cmd = exec.Command("arp", "-a", p.IP)

	default:
		return Result{
			IP:    p.IP,
			Error: "unsupported OS",
		}
	}

	out, err := cmd.Output()
	if err != nil {
		return Result{
			IP:    p.IP,
			Error: err.Error(),
		}
	}

	lines := bytes.Split(out, []byte("\n"))

	for _, line := range lines {
		text := strings.TrimSpace(string(line))
		if text == "" {
			continue
		}

		if strings.Contains(text, p.IP) {
			return Result{
				IP:      p.IP,
				MAC:     extractMAC(text),
				Exists:  true,
				RawLine: text,
			}
		}
	}

	return Result{
		IP:     p.IP,
		Exists: false,
	}
}
func extractMAC(s string) string {
	parts := strings.Fields(s)
	for _, p := range parts {
		if strings.Count(p, ":") == 5 || strings.Count(p, "-") == 5 {
			return p
		}
	}
	return ""
}
