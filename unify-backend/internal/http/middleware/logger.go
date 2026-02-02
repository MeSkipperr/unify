package middleware

import (
	"fmt"
	"time"
	"unify-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// proses request
		c.Next()

		// setelah request selesai
		latency := time.Since(start)
		status := c.Writer.Status()

		// log info
		services.LogInfo(service, fmt.Sprintf("%s | %d | %s | %v", c.Request.Method, status, path, latency))
	}
}
