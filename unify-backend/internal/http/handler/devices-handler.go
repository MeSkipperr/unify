package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"unify-backend/internal/database"
	"unify-backend/internal/services"
	"unify-backend/models"
	"unify-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetDevices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	status := c.QueryArray("status[]")
	notification := c.QueryArray("notification[]")
	types := c.QueryArray("type[]")
	sort := c.Query("sort")
	search := c.Query("search")

	data, total, totalPage, err := services.GetDevicesPaginated(
		page,
		pageSize,
		status,
		notification,
		types,
		sort,
		search,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data":      data,
		"total":     total,
		"page":      page,
		"pageSize":  pageSize,
		"totalPage": totalPage,
	})
}

func CreateDevice() gin.HandlerFunc {
	type CreateDeviceRequest struct {
		NormalizedPayload struct {
			Name        string `json:"name" binding:"required"`
			IPAddress   string `json:"ipAddress" binding:"required"`
			MacAddress  string `json:"macAddress" binding:"required"`
			RoomNumber  string `json:"roomNumber"`
			Description string `json:"description" binding:"required"`
			Type        string `json:"type" binding:"required"`
		} `json:"normalizedPayload" binding:"required"`
	}

	return func(c *gin.Context) {
		var req CreateDeviceRequest
		fmt.Println("Method:", c.Request.Method)
		fmt.Println("Path:", c.Request.URL.Path)
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			fmt.Println("Raw body:", string(bodyBytes))
			// Reset body supaya bisa digunakan lagi oleh ShouldBindJSON
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request body",
				"error":   err.Error(),
			})
			fmt.Println(err)
			return
		}
		payload := req.NormalizedPayload

		ip, err := utils.NormalizeIPv4(payload.IPAddress)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		mac, err := utils.NormalizeMac(payload.MacAddress)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// Convert string to DeviceType enum
		deviceType := models.DeviceType(payload.Type)
		switch deviceType {
		case models.AP, models.IPTV, models.CCTV, models.SW:
			// valid
		default:
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid device type"})
			return
		}

		device := models.Devices{
			IPAddress:    ip,
			Name:         strings.TrimSpace(payload.Name),
			RoomNumber:   strings.ToUpper(strings.TrimSpace(payload.RoomNumber)),
			Description:  strings.TrimSpace(payload.Description),
			Type:         deviceType,
			MacAddress:   mac,
			IsConnect:    false,
			ErrorCount:   0,
			Notification: false,
		}

		if err := database.DB.Create(&device).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to create device",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Device created successfully",
			"data":    device,
		})
	}
}

func DeleteDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ambil id dari param URL
		idParam := c.Param("id")
		deviceID, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid device ID",
				"error":   err.Error(),
			})
			return
		}

		// cari device
		var device models.Devices
		if err := database.DB.First(&device, "id = ?", deviceID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Device not found",
			})
			return
		}

		// hapus device
		if err := database.DB.Delete(&device).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to delete device",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Device deleted successfully",
		})
	}
}
