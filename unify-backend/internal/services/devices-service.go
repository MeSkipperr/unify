package services

import (
	"unify-backend/internal/database"
	"unify-backend/models"
)

func GetDevicesPaginated(page int, pageSize int) ([]models.Devices, int64, error) {
	var devices []models.Devices
	var total int64

	// Hitung total data
	if err := database.DB.Model(&models.Devices{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Hitung offset berdasarkan page dan pageSize
	offset := (page - 1) * pageSize

	// Query dengan pagination
	result := database.DB.
		Limit(pageSize).
		Offset(offset).
		Find(&devices)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return devices, total, nil
}
