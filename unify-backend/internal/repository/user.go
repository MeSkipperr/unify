package repository

import (
	"unify-backend/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	result := r.DB.Where("username = ?", username).First(&user)
	return &user, result.Error
}
