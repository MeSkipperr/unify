package services

import (
	"math"
	"strconv"
	"strings"
	"unify-backend/internal/database"
	"unify-backend/models"
)

type DeviceWithIndex struct {
	models.Devices
	Index int `json:"index"`
}

func GetDevicesPaginated(
	page int,
	pageSize int,
	status []string,
	notification []string,
	types []string,
	sort string,
	search string,
) ([]DeviceWithIndex, int64, int, error) {

	var devices []models.Devices
	result := make([]DeviceWithIndex, 0)
	var total int64

	db := database.DB.Model(&models.Devices{})

	// =====================
	// CLEAN INPUT
	// =====================
	status = cleanArray(status)
	notification = cleanArray(notification)
	types = cleanArray(types)
	search = strings.TrimSpace(search)

	// =====================
	// FILTER (OR via IN)
	// =====================
	if len(status) > 0 {
		db = db.Where("is_connect IN ?", toBoolSlice(status))
	}

	if len(notification) > 0 {
		db = db.Where("notification IN ?", toBoolSlice(notification))
	}

	if len(types) > 0 {
		db = db.Where("type IN ?", types)
	}

	// =====================
	// SEARCH
	// =====================
	if search != "" {
		like := "%" + search + "%"
		db = db.Where(
			"ip_address ILIKE ? OR name ILIKE ? OR room_number ILIKE ? OR type ILIKE ? OR mac_address ILIKE ? OR description ILIKE ?",
			like, like, like, like, like,like,
		)
	}

	// =====================
	// COUNT TOTAL
	// =====================
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	// clamp page
	if page > totalPage && totalPage > 0 {
		page = totalPage
	}

	offset := (page - 1) * pageSize

	// =====================
	// SORTING
	// =====================
	if sort != "" {
		allowedSort := map[string]string{
			"roomNumber": "room_number",
			"lastUpdate": "last_update",
			"index":      "id",
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
		return result, 0, totalPage, nil
	}

	if err := db.
		Limit(pageSize).
		Offset(offset).
		Find(&devices).Error; err != nil {
		return nil, 0, 0, err
	}

	// =====================
	// GLOBAL INDEX
	// =====================
	for i, device := range devices {
		result = append(result, DeviceWithIndex{
			Devices: device,
			Index:   offset + i + 1,
		})
	}

	return result, total, totalPage, nil
}

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

func toBoolSlice(arr []string) []bool {
	result := make([]bool, 0)
	for _, v := range arr {
		if b, err := strconv.ParseBool(v); err == nil {
			result = append(result, b)
		}
	}
	return result
}
