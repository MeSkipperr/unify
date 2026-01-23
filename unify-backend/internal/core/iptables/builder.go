package iptables

import (
	"fmt"
	"unify-backend/models"
)

func BuildRuleArgs(s models.SessionPortForward) []string {
	args := []string{
		"-t", "nat",
		"-A", s.Chain,
		"-d", s.ListenIP,
		"-p", s.Protocol,
		"--dport", fmt.Sprint(s.ListenPort),
	}

	if s.Interface != "" {
		args = append(args, "-i", s.Interface)
	}

	if s.AllowSourceIP != "" {
		args = append(args, "-s", s.AllowSourceIP)
	}

	args = append(args,
		"-j", "DNAT",
		"--to-destination",
		fmt.Sprintf("%s:%d", s.DestIP, s.DestPort),
	)

	if s.RuleComment != "" {
		args = append(args, "-m", "comment", "--comment", s.RuleComment)
	}

	return args
}
