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

	claims, err := utils.VerifyJWT(token, os.Getenv("JWT_SECRET"))
	if err != nil || claims.Type != "access" {
		c.JSON(401, gin.H{"error": "invalid token"})
		c.Abort()
		return
	}

	c.Set("user_id", claims.Sub)
	c.Set("user_name", claims.Name)
	c.Next()
}
