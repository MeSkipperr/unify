package mtr

import (
	"errors"
	"fmt"
	"os/exec"
)

func Run(cfg Config) ([]byte, error) {
	if cfg.DestHost == "" {
		return nil, errors.New("DestHost is required")
	}

	applyDefaults(&cfg)

	// validation
	if (cfg.Protocol == ProtocolTCP || cfg.Protocol == ProtocolUDP) && cfg.Port == nil {
		return nil, errors.New("Port is required for tcp/udp")
	}

	args := buildArgs(cfg)

	cmd := exec.Command("mtr", args...)
	fmt.Println("EXEC:", cmd.String())

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return output, nil
}
