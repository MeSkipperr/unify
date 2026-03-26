package services

import (
	"os"
	"unify-backend/internal/database"
	"unify-backend/models"
	"unify-backend/utils"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string
	Password string
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	var user models.User

	if err := database.DB.
		Where("username = ?", req.Username).
		First(&user).Error; err != nil {

		c.JSON(401, gin.H{"error": "invalid username or password"})
		return
	}

	if err := utils.CheckPassword(*user.Password, req.Password); err != nil {
		c.JSON(401, gin.H{"error": "invalid username or password"})
		return
	}
	accessToken, err := utils.GenerateAccessToken(
		user.ID.String(),
		*user.Username,
		os.Getenv("JWT_SECRET"),
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}
	// Access Token
	c.SetCookie(
		"token",
		accessToken,
		int(utils.AccessTokenTTL.Seconds()),
		"/",
		"",
		false,
		true,
	)


	c.JSON(200, gin.H{"message": "login success"})
}

func Me(c *gin.Context) {
	secret := os.Getenv("JWT_SECRET")

	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(401, gin.H{"message": "unauthorized"})
		return
	}

	claims, err := utils.VerifyJWT(token, secret)
	if err != nil {
		c.JSON(401, gin.H{"message": "invalid token"})
		return
	}

	c.JSON(200, gin.H{
		"user_id": claims.Sub,
		"name":    claims.Name,
	})
}

func LogoutHandler(c *gin.Context) {
	// Hapus Access Token
	c.SetCookie(
		"token",
		"",
		-1, 
		"/",
		"",
		false, 
		true, 
	)

	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)
	c.JSON(200, gin.H{"message": "logged out successfully"})
}
