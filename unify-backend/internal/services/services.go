package services

import (
	"unify-backend/internal/database"
	"unify-backend/models"
)

func GetByServiceName(serviceName string) (*models.Service, error) {
	var service models.Service

	err := database.DB.
		Where("service_name = ?", serviceName).
		First(&service).Error

	if err != nil {
		return nil, err
	}

	return &service, nil
}
