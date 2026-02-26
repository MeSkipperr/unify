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
		// Use arp -a instead of ip neigh
		cmd = exec.Command("arp", "-a", p.IP)

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

		// check if line contains the target IP
		if strings.Contains(text, p.IP) {
			// Example ARP line:
			// ? (172.18.0.9) at 74:81:9a:f2:d0:04 [ether] on eno1
			mac := extractMAC(text)

			exists := mac != "" && !strings.Contains(text, "<incomplete>")

			return Result{
				IP:      p.IP,
				MAC:     mac,
				Exists:  exists,
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
		// check for MAC address format
		if strings.Count(p, ":") == 5 || strings.Count(p, "-") == 5 {
			return p
		}
	}
	return ""
}