package entity

import (
	"time"
)

// LogLevel defines the severity level of a log entry.
type LogLevel string

const (
	LogLevelInfo  LogLevel = "INFO"
	LogLevelError LogLevel = "ERROR"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelDebug LogLevel = "DEBUG"
)

// Log stores application log entries with severity levels.
type Log struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Level     LogLevel  `gorm:"type:varchar(10);not null;index" json:"level"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	Source    string    `gorm:"size:255" json:"source"`
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"createdAt"`
}
