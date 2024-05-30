package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Paginate(ctx *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(ctx.Query("page"))
		if page <= 0 {
			page = 1
		}
		limit := 10

		offset := (page - 1) * limit

		return db.Offset(offset).Limit(limit)
	}
}
