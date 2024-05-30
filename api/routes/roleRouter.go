package routes

import (
	"github.com/berkaycubuk/billiard_software_api/controllers"
	"github.com/gin-gonic/gin"
)

func RoleRoutes(r *gin.Engine) {
	v1 := r.Group("/v1/role")
	{
		v1.GET("/all", controllers.GetRoles())
		v1.POST("/new", controllers.CreateRole())
		v1.POST("/update", controllers.UpdateRole())
		v1.POST("/delete", controllers.DeleteRole())
		v1.GET("/:id", controllers.GetRole())
	}
}
