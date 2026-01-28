package iptables

import (
	"fmt"
	"unify-backend/models"
)

func BuildRuleArgs(s models.SessionPortForward) []string {
	args := []string{
		"-t", "nat",
		"-A", s.Chain,
		"-p", s.Protocol,
		"--dport", fmt.Sprint(s.ListenPort),
	}

	if s.ListenIP != "" {
		args = append(args, "-d", s.ListenIP)
	}

	if s.Interface != "" {
		args = append(args, "-i", s.Interface)
	}

	if s.AllowSourceIP != "" {
		args = append(args, "-s", s.AllowSourceIP)
	}

	// 🔥 PENTING: pilih DNAT vs REDIRECT
	if s.ListenIP == s.DestIP || s.DestIP == "127.0.0.1" {
		args = append(args,
			"-j", "REDIRECT",
			"--to-port", fmt.Sprint(s.DestPort),
		)
	} else {
		args = append(args,
			"-j", "DNAT",
			"--to-destination",
			fmt.Sprintf("%s:%d", s.DestIP, s.DestPort),
		)
	}

	if s.RuleComment != "" {
		args = append(args, "-m", "comment", "--comment", s.RuleComment)
	}

	return args
}

