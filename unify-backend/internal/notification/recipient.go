package notification

import (
	"unify-backend/internal/database"
	"unify-backend/models"
)

func SelectUserByType(types []models.UserRole) ([]models.User, error) {
	var users []models.User

	result := database.DB.
		Where("role IN ?", types).
		Where("is_active = ?", true).
		Find(&users)
		
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}
