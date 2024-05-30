package routes

import (
	"github.com/berkaycubuk/billiard_software_api/controllers"
	"github.com/berkaycubuk/billiard_software_api/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	v1 := r.Group("/v1/user")
	{
		v1.Use(middlewares.AuthMiddleware())
		v1.GET("/all", controllers.GetUsers())
		v1.GET("/all-without-pagination", controllers.GetUsersWithoutPagination())
		v1.GET("/guests", controllers.GetGuests())
		v1.GET("/profile", controllers.GetProfile())
		v1.GET("/permissions", controllers.GetPermissions())
		v1.GET("/my-subscriptions", controllers.GetMySubscriptions())
		v1.GET("/my-active-subscription", controllers.GetMyActiveSubscription())
		v1.POST("/update", controllers.UpdateUser())
		v1.POST("/create", controllers.CreateUser())
		v1.POST("/delete/:id", controllers.DeleteUser())
		v1.GET("/:id", controllers.GetUser())
	}
}
