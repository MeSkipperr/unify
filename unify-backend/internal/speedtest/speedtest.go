package speedtest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

var (
	ErrConfiguration = errors.New("speedtest configuration error")
	ErrNoResult      = errors.New("speedtest finished without result")
)

/*
Run executes speedtest by Ookla using source IP and server ID.

- sourceIP  : IP source (ex: 172.18.1.34)
- serverID  : speedtest server ID (ex: 13623)

Return:
- *SpeedtestResult if success
- error if fatal error detected
*/
func Run(sourceIP string, serverID string) (*SpeedtestResult, error) {
	args := []string{
		"--accept-license",
		"--accept-gdpr",
		"--format=json",
	}

	if sourceIP != "" {
		args = append(args, "--ip="+sourceIP)
	}

	if serverID != "" {
		args = append(args, "--server-id="+serverID)
	}

	cmd := exec.Command("speedtest", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var (
		resultJSON []byte
		fatalErr   error
	)

	scan := func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			// Skip non-json
			if !strings.HasPrefix(line, "{") {
				continue
			}

			// Detect type
			var meta struct {
				Type    string `json:"type"`
				Message string `json:"message"`
			}

			if err := json.Unmarshal([]byte(line), &meta); err != nil {
				continue
			}

			switch meta.Type {
			case "log":
				// Non-fatal bind error
				if strings.Contains(meta.Message, "bind(") {
					continue
				}

				// Fatal configuration error
				if strings.Contains(meta.Message, "Configuration") {
					fatalErr = ErrConfiguration
				}

			case "result":
				resultJSON = []byte(line)
			}
		}
	}

	go scan(bufio.NewScanner(stdout))
	go scan(bufio.NewScanner(stderr))

	if err := cmd.Wait(); err != nil && fatalErr == nil {
		return nil, err
	}

	if fatalErr != nil {
		return nil, fatalErr
	}

	if len(resultJSON) == 0 {
		return nil, ErrNoResult
	}

	var result SpeedtestResult
	decoder := json.NewDecoder(bytes.NewReader(resultJSON))
	decoder.UseNumber()

	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
