package handler

import (
	"strconv"
	"unify-backend/internal/services"

	"github.com/gin-gonic/gin"
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
		"data":       data,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPage":  totalPage,
	})
}
