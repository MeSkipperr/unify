package services

import (
	"fmt"
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

// StopAllServices sets the status of all services to "STOPPED"
func StopAllServices() error {
	result :=database.DB.Model(&models.Service{}).Where("status != ?", "STOPPED").Update("status", "STOPPED")
	if result.Error != nil {
		return fmt.Errorf("failed to stop services: %w", result.Error)
	}

	fmt.Printf("Updated %d service(s) to STOPPED status\n", result.RowsAffected)
	return nil
}
