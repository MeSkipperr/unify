package handler

import (
	"math"
	"strconv"
	"time"
	"unify-backend/internal/database"
	"unify-backend/models"
	"unify-backend/utils"

	"github.com/gin-gonic/gin"
)

func GetActiveMTRSessions() gin.HandlerFunc {
	return func(c *gin.Context) {
		var sessions []models.MTRSession

		if err := database.DB.
			Where("status = ?", "active").
			Order("created_at DESC").
			Find(&sessions).Error; err != nil {

			c.JSON(500, gin.H{
				"message": "failed to fetch active sessions",
			})
			return
		}

		c.JSON(200, gin.H{
			"data": sessions,
		})
	}
}

func GetMTRResult() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if id == "" {
			c.JSON(400, gin.H{
				"message": "session id is required",
			})
			return
		}

		// --- Query param ---
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}

		pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "100"))
		if err != nil || pageSize <= 0 {
			pageSize = 100
		}

		if pageSize > 500 {
			pageSize = 500
		}

		offset := (page - 1) * pageSize

		var results []models.MTRResult
		var total int64

		// --- Hitung total ---
		err = database.DB.
			Model(&models.MTRResult{}).
			Where("session_id = ?", id).
			Count(&total).Error

		if err != nil {
			c.JSON(500, gin.H{
				"message": "failed to count mtr results",
				"error":   err.Error(),
			})
			return
		}

		// --- Subquery: Pagination dalam kondisi DESC ---
		subQuery := database.DB.
			Model(&models.MTRResult{}).
			Where("session_id = ?", id).
			Order("created_at DESC").
			Limit(pageSize).
			Offset(offset)

		// --- Outer query: urutkan ulang ASC ---
		err = database.DB.
			Table("(?) as t", subQuery).
			Order("created_at ASC").
			Find(&results).Error

		if err != nil {
			c.JSON(500, gin.H{
				"message": "failed to fetch mtr results",
				"error":   err.Error(),
			})
			return
		}

		totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

		c.JSON(200, gin.H{
			"data":       results,
			"total":      total,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": totalPages,
		})
	}
}

func DisableMTRSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		now := time.Now().UTC()

		result := database.DB.Model(&models.MTRSession{}).
			Where("id = ? AND status = ?", id, "active").
			Updates(map[string]interface{}{
				"status":      "disable",
				"last_run_at": &now,
			})

		if result.Error != nil {
			c.JSON(500, gin.H{
				"message": "failed to update session",
			})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(404, gin.H{
				"message": "active session not found",
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "session disabled successfully",
		})
	}
}

type CreateMTRSessionRequest struct {
	SourceIP      string `json:"source_ip"`
	DestinationIP string `json:"destination_ip" binding:"required"`

	Protocol string `json:"protocol"`
	Port     *int   `json:"port"`

	Test int `json:"test"`

	Note string `json:"note"`

	SendNotification bool `json:"send_notification"`
}

func CreateMTRSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateMTRSessionRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"message": "invalid request body",
			})
			return
		}

		session := models.MTRSession{
			Status:           "active",
			IsReachable:      false,
			SourceIP:         req.SourceIP,
			DestinationIP:    req.DestinationIP,
			Protocol:         utils.DefaultString(req.Protocol, "icmp"),
			Port:             req.Port,
			Test:             utils.DefaultInt(req.Test, 10),
			Note:             req.Note,
			SendNotification: req.SendNotification,
		}

		if err := database.DB.Create(&session).Error; err != nil {
			c.JSON(500, gin.H{
				"message": "failed to create session",
			})
			return
		}

		c.JSON(201, gin.H{
			"data": session,
		})
	}
}
