package middlewares

import (
	"net/http"

	"github.com/berkaycubuk/billiard_software_api/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := utils.ValidateToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "error.auth_required",
				"success": false,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
