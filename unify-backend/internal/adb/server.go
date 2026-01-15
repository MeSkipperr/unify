package adb

import (
	"fmt"
	"time"
	"unify-backend/internal/system/command"
)

func Kill(ADBPath string) error {
	cmd := fmt.Sprintf("%s kill-server", ADBPath)
	_, err := command.Run(cmd)
	return err
}

func Start(ADBPath string) error {
	cmd := fmt.Sprintf("%s start-server", ADBPath)
	_, err := command.Run(cmd)
	return err
}

func RestartADBServer(ADBPath string) error {
	if err := Kill(ADBPath); err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	
	if err := Start(ADBPath); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	return nil
}
