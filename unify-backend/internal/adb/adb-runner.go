package adb

import (
	"context"
	"os/exec"
	"runtime"
	"strings"
	"time"
)



func run(command string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/C", command)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", command)
	}

	out, err := cmd.CombinedOutput()

	// jika timeout
	if ctx.Err() == context.DeadlineExceeded {
		return "adb command timeout", ctx.Err()
	}

	return string(out), err
}

type AdbRunRequest struct {
	Config   *ADBConfig
	Template string
	Data   map[string]string
}

func AdbRun(opts AdbRunRequest) (AdbStatus, string) {
	connectOutput, err := run(RenderTemplate(opts.Config.CommandTemplate["connect"], opts.Data))
	if err != nil || strings.Contains(strings.ToLower(connectOutput), "failed") {
		return StatusFailed, connectOutput
	}

	adbOutput, err := run(RenderTemplate(opts.Template, opts.Data))
	if strings.Contains(strings.ToLower(adbOutput), "failed") {
		return StatusFailed, adbOutput
	} else if strings.Contains(strings.ToLower(adbOutput), "unauthorized") {
		return StatusUnauthorized, adbOutput
	}
	return StatusSuccess, adbOutput

}
