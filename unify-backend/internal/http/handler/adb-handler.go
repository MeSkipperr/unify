package handler

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"unify-backend/internal/database"
	"unify-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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
	data []models.AdbResult,
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

	// =====================
	// SORT
	// =====================
	allowedSort := map[string]string{
		"start_time":  "start_time",
		"finish_time": "finish_time",
	}

	if sort != "" {
		sorts := strings.Split(sort, ",")
		for _, s := range sorts {
			parts := strings.Split(s, ":")
			if len(parts) != 2 {
				continue
			}

			field := parts[0]
			order := strings.ToUpper(parts[1])

			if col, ok := allowedSort[field]; ok {
				if order != "ASC" && order != "DESC" {
					order = "ASC"
				}
				query = query.Order(col + " " + order)
			}
		}
	} else {
		// default
		query = query.Order("start_time DESC")
	}

	if date != "" {
		query = query.Where(
			"DATE(start_time) = ?", date,
		)
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

	totalPage = int(math.Ceil(float64(total) / float64(pageSize)))
	return
}
