package iptables

import (
	"os/exec"
	"unify-backend/models"
)

func DeleteRule(s models.SessionPortForward) error {
	args := BuildRuleArgs(s)
	args[2] = "-D" // replace -A with -D
	return exec.Command("iptables", args...).Run()
}
