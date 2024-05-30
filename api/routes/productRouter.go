package routes

import (
	"github.com/berkaycubuk/billiard_software_api/controllers"
	"github.com/gin-gonic/gin"
)

func ProductRoutes(r *gin.Engine) {
	v1 := r.Group("/v1/product")
	{
		v1.GET("/all", controllers.GetProducts())
		v1.POST("/new", controllers.CreateProduct())
		v1.POST("/delete", controllers.DeleteProduct())
		v1.POST("/update", controllers.UpdateProduct())
		v1.POST("/update-order", controllers.UpdateProductOrder())
		v1.GET("/:id", controllers.GetProduct())
	}
}
