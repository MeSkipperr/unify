package repository

import (
	"time"
	"unify-backend/models"

	"gorm.io/gorm"
)

type LogRepository struct {
	db *gorm.DB
}

func NewLogRepository(db *gorm.DB) *LogRepository {
	return &LogRepository{db: db}
}

/* =========================
   CREATE
========================= */

func (r *LogRepository) CreateLog(log *models.Log) error {
	return r.db.Create(log).Error
}

/* =========================
   SEARCH (KEY BASED)
========================= */

type SearchLogParams struct {
	ServiceName string
	Level       string
	FromDate    *time.Time
	ToDate      *time.Time
	Limit       int
}

func (r *LogRepository) FindLogs(params SearchLogParams) ([]models.Log, error) {
	var logs []models.Log

	query := r.db.Model(&models.Log{})

	if params.ServiceName != "" {
		query = query.Where("service_name = ?", params.ServiceName)
	}

	if params.Level != "" {
		query = query.Where("level = ?", params.Level)
	}

	if params.FromDate != nil {
		query = query.Where("created_at >= ?", *params.FromDate)
	}

	if params.ToDate != nil {
		query = query.Where("created_at <= ?", *params.ToDate)
	}

	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}

	err := query.Order("created_at DESC").Find(&logs).Error
	return logs, err
}
