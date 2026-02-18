package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unify-backend/internal/database"
	"unify-backend/internal/notification"
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

type CreateDeviceRequest struct {
	NormalizedPayload struct {
		Name          string `json:"name" binding:"required"`
		IPAddress     string `json:"ipAddress" binding:"required"`
		MacAddress    string `json:"macAddress" binding:"required"`
		RoomNumber    string `json:"roomNumber"`
		Description   string `json:"description" binding:"required"`
		Type          string `json:"type" binding:"required"`
		DeviceProduct string `json:"deviceProduct" binding:"required"`
	} `json:"normalizedPayload" binding:"required"`
}

func CreateDevice() gin.HandlerFunc {

	return func(c *gin.Context) {
		var req CreateDeviceRequest
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request body",
				"error":   err.Error(),
			})
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
			IPAddress:       ip,
			Name:            strings.TrimSpace(payload.Name),
			RoomNumber:      strings.ToUpper(strings.TrimSpace(payload.RoomNumber)),
			Description:     strings.TrimSpace(payload.Description),
			DeviceProduct:   payload.DeviceProduct,
			Type:            deviceType,
			MacAddress:      mac,
			IsConnect:       false,
			ErrorCount:      0,
			Notification:    true,
			StatusUpdatedAt: time.Now().UTC(),
		}

		if err := database.DB.Create(&device).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to create device",
				"error":   err.Error(),
			})
			return
		}

		notif := buildDeviceNotification(device, models.DeviceCreated)
		notification.SSENotification(notif)

		c.JSON(http.StatusCreated, gin.H{
			"message": "Device created successfully",
			"data":    device,
		})
	}
}

func ChangeDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil deviceId dari URL
		idParam := c.Param("id")
		deviceId, err := uuid.Parse(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid device ID",
				"error":   err.Error(),
			})
			return
		}
		// Binding request body
		var req CreateDeviceRequest

		if err := c.ShouldBindJSON(&req); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request body",
				"error":   err.Error(),
			})
			return
		}

		payload := req.NormalizedPayload

		// Normalisasi IP dan MAC
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

		// Cari device berdasarkan ID
		var device models.Devices
		if err := database.DB.First(&device, "id = ?", deviceId).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Device not found"})
			return
		}

		// Update field
		device.Name = strings.TrimSpace(payload.Name)
		device.IPAddress = ip
		device.MacAddress = mac
		device.RoomNumber = strings.ToUpper(strings.TrimSpace(payload.RoomNumber))
		device.Description = strings.TrimSpace(payload.Description)
		device.Type = deviceType
		device.DeviceProduct = payload.DeviceProduct

		if err := database.DB.Save(&device).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to update device",
				"error":   err.Error(),
			})
			return
		}

		notif := buildDeviceNotification(device, models.DeviceUpdated)
		notification.SSENotification(notif)

		c.JSON(http.StatusOK, gin.H{
			"message": "Device updated successfully",
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

		notif := buildDeviceNotification(device, models.DeviceDeleted)
		notification.SSENotification(notif)

		c.JSON(http.StatusOK, gin.H{
			"message": "Device deleted successfully",
		})
	}
}

func ChangeNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil deviceId dari URL
		idParam := c.Param("id")
		deviceId, err := uuid.Parse(idParam)
		if err != nil {

			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid device ID",
				"error":   err.Error(),
			})
			return
		}

		// Binding request body
		var req struct {
			Notification bool `json:"notification"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request body",
				"error":   err.Error(),
			})
			return
		}

		// Cari device berdasarkan ID
		var device models.Devices
		if err := database.DB.First(&device, "id = ?", deviceId).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Device not found"})
			return
		}

		// Update notification
		device.Notification = req.Notification

		if err := database.DB.Save(&device).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to update notification",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Notification updated successfully",
			"data":    device,
		})
	}
}

func buildDeviceNotification(device models.Devices, action models.DeviceAction) models.Notification {
	var level models.NotificationLevel
	var title string
	var detail string

	switch action {
	case models.DeviceCreated:
		level = models.NoticationStatusInfo
		title = fmt.Sprintf("[INFO] Device %s Created", device.Name)
		detail = fmt.Sprintf("Device %s was successfully created and is now registered in the system.", device.Name)

	case models.DeviceUpdated:
		level = models.NoticationStatusInfo
		title = fmt.Sprintf("[INFO] Device %s Updated", device.Name)
		detail = fmt.Sprintf("Device %s configuration was successfully updated.", device.Name)

	case models.DeviceDeleted:
		level = models.NoticationStatusAlert
		title = fmt.Sprintf("[ALERT] Device %s Deleted", device.Name)
		detail = fmt.Sprintf("Device %s was removed from the system.", device.Name)

	default:
		level = models.NoticationStatusInfo
		title = fmt.Sprintf("[INFO] Device %s Modified", device.Name)
		detail = fmt.Sprintf("Device %s has been modified.", device.Name)
	}

	return models.Notification{
		Level:  level,
		Title:  title,
		Detail: detail,
		URL:    fmt.Sprintf("/devices?search=%s", device.Name),
	}
}

func GetDeviceSummary() gin.HandlerFunc {
	return func(c *gin.Context) {

		typeParam := c.Query("type")

		// 🔹 Jika ada query type → return 1 type saja
		if typeParam != "" {

			var result struct {
				Total   int64
				Online  int64
				Offline int64
			}

			err := database.DB.
				Model(&models.Devices{}).
				Select(`
					COUNT(*) as total,
					SUM(CASE WHEN is_connect = true THEN 1 ELSE 0 END) as online,
					SUM(CASE WHEN is_connect = false THEN 1 ELSE 0 END) as offline
				`).
				Where("type = ?", typeParam).
				Scan(&result).Error

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"type":    typeParam,
				"total":   result.Total,
				"online":  result.Online,
				"offline": result.Offline,
			})
			return
		}

		// 🔹 Jika tidak ada query → return semua type
		var results []struct {
			Type    models.DeviceType `json:"type"`
			Total   int64             `json:"total"`
			Online  int64             `json:"online"`
			Offline int64             `json:"offline"`
		}

		err := database.DB.
			Model(&models.Devices{}).
			Select(`
				type,
				COUNT(*) as total,
				SUM(CASE WHEN is_connect = true THEN 1 ELSE 0 END) as online,
				SUM(CASE WHEN is_connect = false THEN 1 ELSE 0 END) as offline
			`).
			Group("type").
			Scan(&results).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": results,
		})
	}
}
