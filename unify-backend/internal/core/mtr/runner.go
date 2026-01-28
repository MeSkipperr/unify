package mtr

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"unify-backend/internal/core/dns"
)

func Run(cfg Config) (*MtrResultJson, error) {
	if cfg.DestHost == "" {
		return nil, errors.New("DestHost is required")
	}

	applyDefaults(&cfg)

	if (cfg.Protocol == ProtocolTCP || cfg.Protocol == ProtocolUDP) && cfg.Port == nil {
		return nil, errors.New("Port is required for tcp/udp")
	}

	args := buildArgs(cfg)

	cmd := exec.Command("mtr", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("mtr failed: %v (%s)", err, stderr.String())
	}

	var raw MtrResultJson
	if err := json.Unmarshal(stdout.Bytes(), &raw); err != nil {
		return nil, err
	}

	raw.Report.Result.TotalHops = len(raw.Report.HopResult)
	if len(raw.Report.HopResult) > 0 {
		raw.Report.Result.Reachable =
			raw.Report.HopResult[len(raw.Report.HopResult)-1].Host == cfg.DestHost

			raw.Report.Result.AvgRTT = raw.Report.HopResult[len(raw.Report.HopResult)-1].Avg
		
			for i := 0; i < raw.Report.Result.TotalHops; i++ {
				raw.Report.HopResult[i].Dns = dns.ReverseDNS(raw.Report.HopResult[i].Host)[0]
			}
	}

	return &raw, nil
}
