package mtr

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)
func Run(cfg Config) (*Result, error) {
	if cfg.DestHost == "" {
		return nil, errors.New("DestHost is required")
	}

	applyDefaults(&cfg)

	if (cfg.Protocol == ProtocolTCP || cfg.Protocol == ProtocolUDP) && cfg.Port == nil {
		return nil, errors.New("Port is required for tcp/udp")
	}

	args := buildArgs(cfg)

	cmd := exec.Command("mtr", args...)
	fmt.Println("EXEC:", cmd.String())

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("mtr failed: %v (%s)", err, stderr.String())
	}

	return ParseResult(stdout.Bytes())
}
