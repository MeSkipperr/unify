package services

import (
	"errors"
	"net/http"
	"os"
	"time"

	"unify-backend/internal/database"
	"unify-backend/internal/repository"
	"unify-backend/models"

	"github.com/gin-gonic/gin"
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

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(c *gin.Context) {
	var payload LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", payload.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if user.Password == nil {
		return 
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(payload.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// buat JWT
	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate token"})
		return
	}

	

	// set HttpOnly cookie
	c.SetCookie(
		"token",
		signedToken,
		3600*24, // 1 day
		"/",
		"localhost", // domain sesuai environment
		false,       // secure true jika https
		true,        // HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   signedToken, // optional, bisa diambil client
	})
}

func ProfileHandler(c *gin.Context) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	c.JSON(http.StatusOK, gin.H{
		"id":       claims["id"],
		"username": claims["username"],
		"message":  "Authenticated",
	})
}
