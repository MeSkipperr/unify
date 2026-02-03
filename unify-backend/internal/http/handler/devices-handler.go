package handler

import (
	"strconv"
	"unify-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func  GetDevices(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))

    devices, total, err := services.GetDevicesPaginated(page, pageSize)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{
        "data":  devices,
        "total": total,
        "page":  page,
        "pageSize": pageSize,
    })
}
