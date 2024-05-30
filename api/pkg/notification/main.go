package notification

import (
	"net/http"
	"time"

	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/middlewares"
	"github.com/berkaycubuk/billiard_software_api/models"
	"github.com/berkaycubuk/billiard_software_api/utils"
	"github.com/gin-gonic/gin"
)

func Create(
				userID uint,
				medium uint8,
				typeCode uint8,
				message string,
) {
				notification := models.Notification{
								UserID: userID,
								Medium: medium,
							Type: typeCode,
								Message: message,
								CreatedAt: time.Now(),
				}

				database.DB.Create(&notification)
}

func Find(notificationID uint) (error, models.Notification) {
	var notification models.Notification
	err := database.DB.Where("id = ?", notificationID).First(&notification).Error
	if err != nil {
		return err, models.Notification{}
	}

	return nil, notification
}

func SetRead(notificationID uint) error {
	err, notification := Find(notificationID)
	if err != nil {
		return err
	}

	now := time.Now()
	notification.ReadAt = &now
	database.DB.Save(&notification)

	return nil
}

func GetActiveUserNotifications(userID uint) (error, []models.Notification) {
	var notifications []models.Notification
	err := database.DB.Where("user_id = ? AND read_at IS NULL", userID).
		Order("id desc").
		Find(&notifications).Error
	if err != nil {
		return err, nil
	}

	return nil, notifications
}

func getGameNotifications() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := utils.CurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "errors.auth_required",
				"success": false,
			})
			return
		}

		var notifications []models.Notification
		database.DB.Where("user_id = ? AND type = ?", user.ID, models.NOTIFICATION_TYPE_JOIN_GAME).Find(&notifications)

		ctx.JSON(http.StatusOK, gin.H{
			"notifications": notifications,
			"success": true,
		})
	}
}

func getMyNotifications() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := utils.CurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "errors.auth_required",
				"success": false,
			})
			return
		}

		_, notifications := GetActiveUserNotifications(user.ID)
		ctx.JSON(http.StatusOK, gin.H{
			"notifications": notifications,
			"success": true,
		})
	}
}

type MarkAsReadRequest struct {
	ID uint	`json:"id" binding:"required" validate:"required"`
}

func markAsRead() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request MarkAsReadRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		SetRead(request.ID)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func genNot() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		Create(1, models.NOTIFICATION_MEDIUM_APP, models.NOTIFICATION_TYPE_SIMPLE, "test")
	}
}

type CreateNotificationRequest struct {
	Message string	`json:"message" binding:"required" validate:"required"`
}

func createNotification() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request CreateNotificationRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		var users []models.User
		database.DB.Find(&users)

		for _, v := range users {
			Create(v.ID, models.NOTIFICATION_MEDIUM_APP, models.NOTIFICATION_TYPE_SIMPLE, request.Message)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func Routes(r *gin.RouterGroup) {
	routes := r.Group("/notification")
	{
		routes.Use(middlewares.AuthMiddleware())
		routes.GET("/active", getMyNotifications())
		routes.GET("/game", getGameNotifications())
		routes.POST("/gen", genNot())
		routes.POST("/create", createNotification())
		routes.POST("/mark", markAsRead())
	}
}
