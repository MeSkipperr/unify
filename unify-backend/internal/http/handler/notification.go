package handler

import (
	"math"
	"net/http"
	"strconv"
	"unify-backend/internal/database"
	"unify-backend/models"

	"github.com/gin-gonic/gin"
)

func GetNotifications() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Ambil query params
        pageStr := c.DefaultQuery("page", "1")
        pageSizeStr := c.DefaultQuery("pageSize", "25")

        page, err := strconv.Atoi(pageStr)
        if err != nil || page < 1 {
            page = 1
        }

        pageSize, err := strconv.Atoi(pageSizeStr)
        if err != nil || pageSize < 1 {
            pageSize = 25
        }

        offset := (page - 1) * pageSize

        var notifications []models.Notification
        var total int64

        // Hitung total data untuk pagination info
        if err := database.DB.Model(&models.Notification{}).Count(&total).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": err.Error(),
            })
            return
        }

        // Ambil data dengan limit & offset, urut DESC (terbaru di atas)
        if err := database.DB.
            Order("created_at DESC").
            Limit(pageSize).
            Offset(offset).
            Find(&notifications).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": err.Error(),
            })
            return
        }

        // Response dengan info pagination
        c.JSON(http.StatusOK, gin.H{
            "data":       notifications,
            "page":       page,
            "pageSize":   pageSize,
            "total":      total,
            "totalPages": int(math.Ceil(float64(total) / float64(pageSize))),
        })
    }
}
