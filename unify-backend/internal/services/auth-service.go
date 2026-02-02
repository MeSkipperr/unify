package services

import (
	"fmt"
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

	refreshToken, err := utils.GenerateRefreshToken(
		user.ID.String(),
		os.Getenv("JWT_SECRET"),
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate refresh token"})
		return
	}
	// Access Token
	c.SetCookie(
		"token",
		accessToken,
		int(utils.AccessTokenTTL.Seconds()),
		"/",
		"",
		true,
		true,
	)

	// Refresh Token
	c.SetCookie(
		"refresh_token",
		refreshToken,
		int(utils.RefreshTokenTTL.Seconds()),
		"/",
		"",
		true,
		true,
	)

	c.JSON(200, gin.H{"message": "login success"})
}

func Me(c *gin.Context) {
	secret := os.Getenv("JWT_SECRET")
	accessToken, err := c.Cookie("token")
	if err == nil {
		claims, err := utils.VerifyJWT(accessToken, secret)
		if err == nil && claims.Type == "access" {
			c.JSON(200, gin.H{
				"user_id": claims.Sub,
				"name":    claims.Name,
			})
			return
		}
	}

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(401, gin.H{"message": "unauthorized"})
		return
	}

	claims, err := utils.VerifyJWT(refreshToken, secret)
	if err != nil || claims.Type != "refresh" {
		c.JSON(401, gin.H{"message": "session expired"})
		return
	}

	newAccess, _ := utils.GenerateAccessToken(claims.Sub, "", secret)

	c.SetCookie(
		"access_token",
		newAccess,
		int(utils.AccessTokenTTL.Seconds()),
		"/",
		"",
		true,
		true,
	)

	c.JSON(200, gin.H{
		"user_id": claims.Sub,
	})
}

func RefreshTokenHandler(c *gin.Context) {
	token, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(401, gin.H{"error": "refresh token missing"})
		return
	}

	fmt.Println("token", token)
	claims, err := utils.VerifyJWT(token, os.Getenv("JWT_SECRET"))
	fmt.Println("claims,", claims)
	if err != nil || claims.Type != "refresh" {
		c.JSON(401, gin.H{"error": "invalid refresh token"})
		return
	}

	newAccessToken, _ := utils.GenerateAccessToken(
		claims.Sub,
		claims.Name,
		os.Getenv("JWT_SECRET"),
	)

	c.SetCookie("token", newAccessToken, 900, "/", "", true, true)
	c.JSON(200, gin.H{"message": "token refreshed"})
}
