package adb

import (
	"os/exec"
	"runtime"
	"strings"
)



func run(command string) (string, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	out, err := cmd.CombinedOutput()
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
