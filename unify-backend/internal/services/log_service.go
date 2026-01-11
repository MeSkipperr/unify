package services

import (
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
	logRepo := repository.NewLogRepository(database.DB)
	logService := NewLogService(logRepo)

	return logService.CreateLog(params)
}
