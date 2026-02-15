package handler

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"unify-backend/internal/adb"
	"unify-backend/internal/database"
	"unify-backend/internal/job"
	"unify-backend/internal/queue"
	"unify-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdbResultWithIndex struct {
	models.AdbResult
	Index int `json:"index" gorm:"-"`
}

func GetAdbResults(c *gin.Context) {

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	typeServices := c.QueryArray("typeServices[]")
	sort := c.Query("sort")
	search := c.Query("search")
	date := c.Query("date")

	data, total, totalPage, err := GetAdbResultsPaginated(
		database.DB,
		page,
		pageSize,
		typeServices,
		sort,
		search,
		date,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get adb results",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      data,
		"total":     total,
		"page":      page,
		"pageSize":  pageSize,
		"totalPage": totalPage,
	})
}

func GetAdbResultsPaginated(
	db *gorm.DB,
	page int,
	pageSize int,
	typeServices []string,
	sort string,
	search string,
	date string,
) (
	data []AdbResultWithIndex,
	total int64,
	totalPage int,
	err error,
) {
	offset := (page - 1) * pageSize

	query := db.Model(&models.AdbResult{})

	// =====================
	// FILTER: TypeServices
	// =====================
	if len(typeServices) > 0 {
		query = query.Where("type_services IN ?", typeServices)
	}

	// =====================
	// SEARCH
	// =====================
	if search != "" {
		like := "%" + search + "%"
		query = query.Where(`
			CAST(status AS TEXT) ILIKE ? OR
			ip_address ILIKE ? OR
			CAST(port AS TEXT) ILIKE ? OR
			name_device ILIKE ?
		`, like, like, like, like)
	}

	// =====================
	// COUNT
	// =====================
	if err = query.Count(&total).Error; err != nil {
		return
	}

	if date != "" {
		query = query.Where(
			"DATE(start_time) = ?", date,
		)
	}
	// =====================
	// SORT
	// =====================

	if sort != "" {
		allowedSort := map[string]string{
			"startTime":  "start_time",
			"finishTime": "finish_time",
		}

		orders := strings.Split(sort, ",")

		for _, order := range orders {
			part := strings.Split(order, ":")
			if len(part) != 2 {
				continue
			}

			columnKey := strings.TrimSpace(part[0])
			directionRaw := strings.ToUpper(strings.TrimSpace(part[1]))

			column, ok := allowedSort[columnKey]
			if !ok {
				continue
			}

			var dir string
			switch directionRaw {
			case "DESC", "DESCENDING":
				dir = "DESC"
			case "ASC", "ASCENDING":
				dir = "ASC"
			default:
				dir = "ASC"
			}

			db = query.Order(column + " " + dir)
		}
	}
	// =====================
	// FETCH DATA
	// =====================
	err = query.
		Limit(pageSize).
		Offset(offset).
		Find(&data).
		Error

	if err != nil {
		return
	}

	for i := range data {
		data[i].Index = offset + i + 1
	}

	totalPage = int(math.Ceil(float64(total) / float64(pageSize)))
	return
}

func CreateADBJob() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			IPAddress string `json:"ipAddress" binding:"required"`
			Port      int    `json:"port" binding:"required"`
			Command   string `json:"command" binding:"required"`
			Name      string `json:"name" binding:"required"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			fmt.Println("Error binding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		now := time.Now().UTC()

		data := models.AdbResult{
			Status:       adb.StatusNotStarted,
			FinishTime:   now,
			StartTime:    now,
			IPAddress:    payload.IPAddress,
			Port:         payload.Port,
			NameDevice:   payload.Name,
			Result:       "",
			TypeServices: "manual",
			Command:      payload.Command,
		}

		if err := database.DB.Create(&data).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		newJob := job.ADBJob{
			ID:        data.ID,
			Command:   payload.Command,
			IPAddress: payload.IPAddress,
			Port:      payload.Port,
			Name:      payload.Name,
		}

		queue.EnqueueADB(newJob)

		c.JSON(http.StatusCreated, gin.H{
			"message": "ADB Job created successfully.",
			"data":    newJob,
		})
	}
}
