package middleware

import (
	"os"
	"unify-backend/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(401, gin.H{"error": "unauthorized"})
		c.Abort()
		return
	}

	secret := os.Getenv("JWT_SECRET")

	claims, err := utils.VerifyJWT(token, secret)
	if err != nil || claims.Type != "access" {
		c.JSON(401, gin.H{"error": "invalid token"})
		c.Abort()
		return
	}

	newToken, _ := utils.GenerateAccessToken(
		claims.Sub,
		claims.Name,
		secret,
	)

	if newToken != token {
		c.SetCookie(
			"token",
			newToken,
			int(utils.AccessTokenTTL.Seconds()),
			"/",
			"",
			false,
			true,
		)
	}

	c.Next()
}