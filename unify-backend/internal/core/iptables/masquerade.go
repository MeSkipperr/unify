package iptables

import (
	"fmt"
	"unify-backend/models"
)

func BuildMasqueradeArgs(s models.SessionPortForward) []string {
	return []string{
		"-t", "nat",
		"-A", "POSTROUTING",
		"-p", s.Protocol,
		"-d", s.DestIP,
		"--dport", fmt.Sprint(s.DestPort),
		"-j", "MASQUERADE",
	}
}
