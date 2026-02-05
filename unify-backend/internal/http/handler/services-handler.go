package handler

import (
	"net/http"
	"strconv"
	"unify-backend/internal/database"
	"unify-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler untuk mendapatkan seluruh service dengan pagination
func GetServices() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))

		if page < 1 {
			page = 1
		}
		if pageSize < 1 {
			pageSize = 50
		}

		var services []models.Service
		var total int64

		// Hitung total data
		database.DB.Model(&models.Service{}).Count(&total)

		// Ambil data dengan limit & offset
		err := database.DB.
			Limit(pageSize).
			Offset((page - 1) * pageSize).
			Find(&services).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		totalPage := (total + int64(pageSize) - 1) / int64(pageSize) // hitung total page

		c.JSON(http.StatusOK, gin.H{
			"data":      services,
			"total":     total,
			"page":      page,
			"pageSize":  pageSize,
			"totalPage": totalPage,
		})
	}
}

// Handler untuk mendapatkan service berdasarkan ServiceName
func GetServiceByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceName := c.Param("serviceName")
		if serviceName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "serviceName is required"})
			return
		}

		var service models.Service
		err := database.DB.Preload("ServiceStatus").
			Where("service_name = ?", serviceName).
			First(&service).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, service)
	}
}