package routes

import (
	"github.com/berkaycubuk/billiard_software_api/controllers"
	"github.com/berkaycubuk/billiard_software_api/middlewares"
	"github.com/gin-gonic/gin"
)

func SubscriptionRoutes(r *gin.Engine) {
	v1 := r.Group("/v1/subscription")
	{
		v1.Use(middlewares.AuthMiddleware())
		v1.GET("/all", controllers.GetSubscriptions())
		v1.POST("/buy", controllers.BuySubscription())
		v1.POST("/update", controllers.UpdateSubscription())
		v1.POST("/pause", controllers.PauseSubscription())
		v1.POST("/unpause", controllers.UnpauseSubscription())
		v1.POST("/new", controllers.NewSubscription())
		v1.POST("/add-user-subscription", controllers.AddUserSub())
		v1.POST("/pause-user-subscription", controllers.PauseUserSub())
		v1.POST("/delete", controllers.DeleteSubscription())
		v1.POST("/delete-user-subscription", controllers.DeleteUserSubscription())
		v1.GET("/:id", controllers.GetSubscription())
	}
}
