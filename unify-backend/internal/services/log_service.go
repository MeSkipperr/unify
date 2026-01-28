package services

import (
	"log"
	"time"
	"unify-backend/internal/database"
	"unify-backend/internal/repository"

	// "unify-backend/internal/services"
	"unify-backend/models"
)

/* =========================
	CREATE LOG PARAM
========================= */

type CreateLogParams struct {
	Level       string
	ServiceName string
	Message     string
	Timestamp   time.Time
}

/* =========================
   	SEARCH LOG PARAM
========================= */

type SearchLogParams struct {
	ServiceName string
	Level       string
	FromDate    *time.Time
	ToDate      *time.Time
	Limit       int
}

type LogService struct {
	repo *repository.LogRepository
}

func NewLogService(repo *repository.LogRepository) *LogService {
	return &LogService{repo: repo}
}

/* =========================
   CREATE LOG
========================= */

func (s *LogService) CreateLog(params CreateLogParams) error {
	log := models.Log{
		Level:       params.Level,
		ServiceName: params.ServiceName,
		Message:     params.Message,
	}

	return s.repo.CreateLog(&log)
}

/* =========================
    SEARCH LOG
========================= */

func (s *LogService) SearchLogs(params SearchLogParams) ([]models.Log, error) {
	return s.repo.FindLogs(repository.SearchLogParams{})
}

func CreateAppLog(params CreateLogParams) error {
	log.Print(
		// "[", params.Timestamp.Format("2006-01-02 15:04:05"), "] ",
		"[", params.Level, "] ",
		"[", params.ServiceName, "] ",
		params.Message,
	)
	logRepo := repository.NewLogRepository(database.DB)
	logService := NewLogService(logRepo)

	return logService.CreateLog(params)
}

func LogInfo(serviceName, msg string) {
	CreateAppLog(CreateLogParams{
		Level:       "INFO",
		ServiceName: serviceName,
		Message:     msg,
	})
}

func LogError(serviceName, msg string) {
	CreateAppLog(CreateLogParams{
		Level:       "ERROR",
		ServiceName: serviceName,
		Message:     msg,
	})
}
func LogWarning(serviceName, msg string) {
	CreateAppLog(CreateLogParams{
		Level:       "WARN",
		ServiceName: serviceName,
		Message:     msg,
	})
}
