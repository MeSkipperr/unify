package handler

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unify-backend/internal/database"
	"unify-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func cleanArray(arr []string) []string {
	result := make([]string, 0)
	for _, v := range arr {
		v = strings.TrimSpace(v)
		v = strings.Trim(v, `"`)
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}

type PortForwardWithIndex struct {
	models.SessionPortForward
	Index int `json:"index"`
}

func GetPortForward() gin.HandlerFunc {
	return func(c *gin.Context) {

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))

		if page < 1 {
			page = 1
		}
		if pageSize < 1 {
			pageSize = 50
		}

		protocol := cleanArray(c.QueryArray("protocol[]"))
		status := cleanArray(c.QueryArray("status[]"))
		sort := c.Query("sort")
		search := strings.TrimSpace(c.Query("search"))

		var portForward []models.SessionPortForward
		var total int64

		// =====================
		// BASE QUERY
		// =====================
		db := database.DB.Model(&models.SessionPortForward{})

		// =====================
		// FILTER
		// =====================
		if len(status) > 0 {
			db = db.Where("status IN ?", status)
		}

		if len(protocol) > 0 {
			db = db.Where("protocol IN ?", protocol)
		}

		// =====================
		// SEARCH
		// =====================
		if search != "" {
			like := "%" + search + "%"
			db = db.Where(`
				listen_ip ILIKE ? OR 
				dest_ip ILIKE ? OR 
				protocol ILIKE ? OR 
				status ILIKE ?
			`, like, like, like, like)
		}

		// =====================
		// COUNT (WAJIB pakai db yang sudah difilter)
		// =====================
		if err := db.Count(&total).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

		if page > totalPage && totalPage > 0 {
			page = totalPage
		}

		offset := (page - 1) * pageSize

		// =====================
		// SORTING
		// =====================
		if sort != "" {
			allowedSort := map[string]string{
				"startTime":  "created_at",
				"finishTime": "expires_at",
			}

			orders := strings.Split(sort, ",")
			for _, order := range orders {
				part := strings.Split(order, ":")
				if len(part) == 2 {
					if column, ok := allowedSort[part[0]]; ok {
						dir := strings.ToUpper(part[1])
						if dir != "DESCENDING" {
							dir = "ASC"
						} else {
							dir = "DESC"
						}
						db = db.Order(column + " " + dir)
					}
				}
			}
		}

		// =====================
		// QUERY DATA
		// =====================
		if total == 0 {
			c.JSON(http.StatusOK, gin.H{
				"data":      []models.SessionPortForward{},
				"total":     0,
				"page":      page,
				"pageSize":  pageSize,
				"totalPage": totalPage,
			})
			return
		}

		if err := db.
			Limit(pageSize).
			Offset(offset).
			Find(&portForward).Error; err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		startIndex := offset + 1

		var response []PortForwardWithIndex
		for i, item := range portForward {
			response = append(response, PortForwardWithIndex{
				Index:              startIndex + i,
				SessionPortForward: item,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"data":      response,
			"total":     total,
			"page":      page,
			"pageSize":  pageSize,
			"totalPage": totalPage,
		})
	}
}

func CreatePortForward() gin.HandlerFunc {
	return func(c *gin.Context) {

		var payload struct {
			ListenIP    string    `json:"listenIp" binding:"required"`
			DestIP      string    `json:"destIp" binding:"required"`
			DestPort    int       `json:"destPort" binding:"required"`
			Protocol    string    `json:"protocol" binding:"required"`
			ExpiresAt   time.Time `json:"expiresAt" binding:"required"`
			RuleComment string    `json:"ruleComment" binding:"required"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			fmt.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		db := database.DB

		const minPort = 20000
		const maxPort = 29999

		// 🔎 Ambil semua listen_port active dalam range
		var usedPorts []int
		if err := db.
			Model(&models.SessionPortForward{}).
			Where("status = ? AND listen_port BETWEEN ? AND ?", "active", minPort, maxPort).
			Pluck("listen_port", &usedPorts).Error; err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Buat map untuk cek cepat
		usedMap := make(map[int]bool)
		for _, p := range usedPorts {
			usedMap[p] = true
		}

		// 🔎 Cari port kosong
		var selectedPort int
		for port := minPort; port <= maxPort; port++ {
			if !usedMap[port] {
				selectedPort = port
				break
			}
		}

		// Jika tidak ada port tersedia
		if selectedPort == 0 {
			c.JSON(http.StatusConflict, gin.H{
				"message": "No available ports in range 20000-29999.",
			})
			return
		}

		// ✅ Buat data baru dengan status pending
		newSession := models.SessionPortForward{
			ListenIP:    payload.ListenIP,
			ListenPort:  selectedPort,
			DestIP:      payload.DestIP,
			DestPort:    payload.DestPort,
			Protocol:    payload.Protocol,
			ExpiresAt:   payload.ExpiresAt,
			Status:      "pending",
			RuleComment: payload.RuleComment,
		}

		if err := db.Create(&newSession).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Port forward request created successfully.",
			"data":    newSession,
		})
	}
}

func DeactivatedPortForward() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "ID is required",
			})
			return
		}

		db := database.DB

		// 🔎 Cari data berdasarkan ID
		var session models.SessionPortForward
		if err := db.First(&session, "id = ?", id).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Port forward not found",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Jika sudah inactive, tidak perlu update
		if session.Status == models.SessionStatusDeactivated {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Port forward is already deactivate",
			})
			return
		}

		// ✅ Update status jadi inactive
		if err := db.Model(&session).
			Update("status", models.SessionStatusDeactivated).Error; err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Port forward has been successfully deactivated.",
		})
	}
}
