package entity

import (
	"time"

	"github.com/google/uuid"
)

// ScheduleType defines how an ADB job is triggered.
type ScheduleType string

const (
	ScheduleTypeCron     ScheduleType = "cron"
	ScheduleTypeInterval ScheduleType = "interval"
)

// ADBJobStatus represents the execution state of an ADB job, target, or result.
type ADBJobStatus string

const (
	ADBJobStatusPending ADBJobStatus = "pending"
	ADBJobStatusRunning ADBJobStatus = "running"
	ADBJobStatusSuccess ADBJobStatus = "success"
	ADBJobStatusFailed  ADBJobStatus = "failed"
	ADBJobStatusTimeout ADBJobStatus = "timeout"
)

// ADBJob stores scheduled ADB command executions with schedule config.
// Relation: ADBJob -> ADBJobTarget -> ADBJobResult
type ADBJob struct {
	ID              uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name            string       `gorm:"size:255;not null" json:"name"`
	CommandKey      string       `gorm:"size:255;not null" json:"commandKey"`
	PackageName     string       `gorm:"size:255" json:"packageName"`
	ScheduleType    ScheduleType `gorm:"type:varchar(20);not null" json:"scheduleType"`
	CronExpression  string       `gorm:"size:100" json:"cronExpression"`
	IntervalSeconds int          `gorm:"default:0" json:"intervalSeconds"`
	Enabled         bool         `gorm:"default:true" json:"enabled"`
	LastRunAt       *time.Time   `json:"lastRunAt"`
	NextRunAt       *time.Time   `json:"nextRunAt"`
	CreatedAt       time.Time    `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time    `gorm:"autoUpdateTime" json:"updatedAt"`

	Targets []ADBJobTarget `gorm:"foreignKey:JobID;references:ID;constraint:OnDelete:CASCADE" json:"targets"`
}

// ADBJobTarget links an ADB job to a specific device for execution.
type ADBJobTarget struct {
	ID        uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	JobID     uuid.UUID    `gorm:"type:uuid;not null" json:"jobId"`
	DeviceID  uuid.UUID    `gorm:"type:uuid;not null" json:"deviceId"`
	Status    ADBJobStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	CreatedAt time.Time    `gorm:"autoCreateTime" json:"createdAt"`

	Job     *ADBJob        `gorm:"foreignKey:JobID;references:ID" json:"job,omitempty"`
	Device  *Devices       `gorm:"foreignKey:DeviceID;references:ID" json:"device,omitempty"`
	Results []ADBJobResult `gorm:"foreignKey:TargetID;references:ID;constraint:OnDelete:CASCADE" json:"results"`
}

// ADBJobResult stores the output of an ADB command execution on a target device.
type ADBJobResult struct {
	ID          uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TargetID    uuid.UUID    `gorm:"type:uuid;not null" json:"targetId"`
	Status      ADBJobStatus `gorm:"type:varchar(20);not null" json:"status"`
	Output      string       `gorm:"type:text" json:"output"`
	StartedAt   time.Time    `json:"startedAt"`
	CompletedAt *time.Time   `json:"completedAt"`

	Target *ADBJobTarget `gorm:"foreignKey:TargetID;references:ID" json:"target,omitempty"`
}

// ADBCommand stores a named ADB command template from configuration.
type ADBCommand struct {
	Key     string `json:"key"`
	Command string `json:"command"`
}

// ADBPackage stores an Android application package entry from configuration.
type ADBPackage struct {
	PackageName string `json:"packageName"`
	DisplayName string `json:"displayName"`
}

// ADBConfig is the top-level ADB subsystem configuration.
type ADBConfig struct {
	ADBPath           string                `json:"adbPath"`
	ADBPort           int                   `json:"adbPort"`
	VerificationSteps int                   `json:"verificationSteps"`
	ADBListCommand    []ADBCommand          `json:"adbListCommand"`
	Packages          map[string]ADBPackage `json:"packages"`
}
