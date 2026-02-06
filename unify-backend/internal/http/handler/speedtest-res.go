package handler

import (
	"net/http"
	"strconv"
	"unify-backend/internal/database"
	"unify-backend/models"

	"github.com/gin-gonic/gin"
)

func GetSpeedtestByInternalIPAndServer() gin.HandlerFunc {
	return func(c *gin.Context) {
		internalIP := c.Query("internalIp")
		serverIdStr := c.Query("serverId")

		if internalIP == "" || serverIdStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "internalIp and serverId are required",
			})
			return
		}

		serverID, err := strconv.Atoi(serverIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "serverId must be a valid number",
			})
			return
		}

		var results []models.SpeedtestResult

		err = database.DB.
			Where("internal_ip = ? AND server_id = ?", internalIP, serverID).
			Order("tested_at ASC").
			Limit(30).
			Find(&results).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if len(results) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Speedtest results not found",
			})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}
