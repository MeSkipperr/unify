package iptables

import (
	"os/exec"
	"time"
	"unify-backend/models"
)

func ApplyRule(s *models.SessionPortForward) error {
	if err := EnsureChain(s.Chain); err != nil {
		return err
	}

	args := BuildRuleArgs(*s)
	if err := exec.Command("iptables", args...).Run(); err != nil {
		return err
	}

	now := time.Now()
	s.LastAppliedAt = &now
	return nil
}
