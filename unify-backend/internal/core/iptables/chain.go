package iptables

import "os/exec"

func EnsureChain(chain string) error {
	// create chain if not exist
	exec.Command("iptables", "-t", "nat", "-N", chain).Run()

	// link PREROUTING -> chain (ignore error if exists)
	cmd := exec.Command(
		"iptables", "-t", "nat",
		"-C", "PREROUTING",
		"-j", chain,
	)
	if cmd.Run() != nil {
		return exec.Command(
			"iptables", "-t", "nat",
			"-A", "PREROUTING",
			"-j", chain,
		).Run()
	}
	return nil
}
