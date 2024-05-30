package routes

import (
	"github.com/berkaycubuk/billiard_software_api/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	v1 := r.Group("/v1/auth")
	{
		v1.POST("/login", controllers.Login())
		v1.POST("/register", controllers.Register())
		v1.POST("/update-password", controllers.UpdatePassword())
		v1.POST("/forgot-password", controllers.ForgotPassword())
		v1.POST("/password-reset-complete", controllers.PasswordResetComplete())
		v1.POST("/verify", controllers.Verify())
	}
}
