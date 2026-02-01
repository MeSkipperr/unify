package services

import (
	"errors"
	"time"

	"unify-backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo *repository.UserRepository
	Secret   string
}

func NewAuthService(repo *repository.UserRepository, secret string) *AuthService {
	return &AuthService{
		UserRepo: repo,
		Secret:   secret,
	}
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.UserRepo.FindByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	if user.Password == nil {
		return "", errors.New("password not set")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
