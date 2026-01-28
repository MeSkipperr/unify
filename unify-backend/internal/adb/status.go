package adb

type AdbStatus string

const (
	StatusNotStarted   AdbStatus = "NOT_STARTED"
	StatusFailed       AdbStatus = "FAILED"
	StatusNotConnected AdbStatus = "NOT_CONNECTED"
	StatusUnauthorized AdbStatus = "UNAUTHORIZED"
	StatusFailedClear  AdbStatus = "FAILED_CLEAR"
	StatusFailedUptime AdbStatus = "FAILED_UPTIME"
	StatusSuccess      AdbStatus = "SUCCESS"
)
