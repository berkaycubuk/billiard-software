package controllers

import (
	"net/http"
	"time"

	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/models"
	orderModule "github.com/berkaycubuk/billiard_software_api/pkg/order"
	"github.com/berkaycubuk/billiard_software_api/pkg/subscription"
	"github.com/berkaycubuk/billiard_software_api/utils"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type GetSubscriptionRequest struct {
	ID uint `uri:"id" binding:"required"`
}

type UpdateSubscriptionRequest struct {
	ID    uint    `json:"id" binding:"required" validate:"required"`
	Name  string  `json:"name" binding:"required" validate:"required"`
	Price float64 `json:"price" binding:"required" validate:"required"`
	Hours int     `json:"hours" binding:"required" validate:"required"`
	Role uint	  `json:"role"`
	Hidden bool	  `json:"hidden"`
}

type NewSubscriptionRequest struct {
	Name  string  `json:"name" binding:"required" validate:"required"`
	Price float64 `json:"price" binding:"required" validate:"required"`
	Hours int     `json:"hours" binding:"required" validate:"required"`
	Role uint	  `json:"role"`
	Hidden bool	  `json:"hidden"`
}

type DeleteSubscriptionRequest struct {
	ID uint `json:"id" binding:"required" validate:"required"`
}

type PauseSubscriptionRequest struct {
	ID uint `json:"id" binding:"required" validate:"required"`
}

func GetSubscription() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request GetSubscriptionRequest
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.id_required",
				"success": false,
			})
			return
		}

		var subscription models.Subscription
		err := database.DB.Where("id = ?", request.ID).First(&subscription).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"success": false,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"subscription": subscription,
			"success":      true,
		})
	}
}

func GetSubscriptions() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var subscriptions []models.Subscription
		err := database.DB.Find(&subscriptions).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"success": false,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"subscriptions": subscriptions,
			"success":       true,
		})
	}
}

type PauseUserSubscriptionRequest struct {
	UserID uint `json:"user_id"`
	SubID uint `json:"sub_id"`
}

func PauseUserSub() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request PauseUserSubscriptionRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		var user models.User
		err := database.DB.Where("id = ?", request.UserID).First(&user).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error.user_not_found",
				"success": false,
			})
			return
		}

		// check subscription exists
		var subscription models.UserSubscription
		err = database.DB.Where("id = ?", request.SubID).First(&subscription).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.subscription_not_found",
				"success": false,
			})
			return
		}

		// check it's state
		if subscription.Status == models.USER_SUBSCRIPTION_STATUS_ACTIVE {
			database.DB.Create(&models.UserSubscriptionChunk{
				UserSubscriptionID: subscription.ID,
				Action:             models.USER_SUBSCRIPTION_ACTION_PAUSED,
				ActionBy:           &user.ID,
				CreatedAt:          time.Now(),
			})

			database.DB.Model(&subscription).Update("status", models.USER_SUBSCRIPTION_STATUS_PAUSED)

			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
			})
			return
		} else if subscription.Status == models.USER_SUBSCRIPTION_STATUS_PAUSED {
			database.DB.Create(&models.UserSubscriptionChunk{
				UserSubscriptionID: subscription.ID,
				Action:             models.USER_SUBSCRIPTION_ACTION_ACTIVATED,
				ActionBy:           &user.ID,
				CreatedAt:          time.Now(),
			})

			database.DB.Model(&subscription).Update("status", models.USER_SUBSCRIPTION_STATUS_ACTIVE)

			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors.nothing_happened",
			"success": false,
		})
	}
}

type AddUserSubscriptionRequest struct {
	UserID uint `json:"user_id"`
	SubID uint `json:"sub_id"`
}

func AddUserSub() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request AddUserSubscriptionRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		// check user
		var user models.User
		err := database.DB.Where("id = ?", request.UserID).First(&user).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.user_not_found",
				"success": false,
			})
			return
		}

		// check sub
		var sub models.Subscription
		err = database.DB.Where("id = ?", request.SubID).First(&sub).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.subscription_not_found",
				"success": false,
			})
			return
		}

		err = subscription.ActivateSubscription(request.UserID, request.SubID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"success": false,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

type BuySubscriptionRequest struct {
	ID uint `json:"id"`
}

func BuySubscription() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request BuySubscriptionRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		user, err := utils.CurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "errors.auth_required",
				"success": false,
			})
			return
		}

		// find subscription
		var subscription models.Subscription
		err = database.DB.Where("id = ?", request.ID).First(&subscription).Error
		if err != nil { // not found
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.subscription_not_found",
				"success": false,
			})
			return
		}

		order, err := orderModule.Create(&user.ID, subscription.Price, models.ORDER_STATUS_WAITING)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.unable_to_create_order",
				"success": false,
			})
			return
		}

		orderModule.CreateOrderItem(order.ID, models.PRODUCT_TYPE_SUBSCRIPTION, &subscription.ID, subscription.Name, 1, subscription.Price)

		orderModule.CreateOrderDetails(order.ID, user.Name, user.Surname, user.Phone, user.Email)

		ctx.JSON(http.StatusOK, gin.H{
			"order_id": order.ID,
			"success":  true,
		})
	}
}

func DeleteUserSubscription() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request DeleteSubscriptionRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		// check subscription exists
		var subscription models.UserSubscription
		err := database.DB.Where("id = ?", request.ID).First(&subscription).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.subscription_not_found",
				"success": false,
			})
			return
		}

		err = database.DB.Delete(&subscription).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"success": false,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func DeleteSubscription() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request DeleteSubscriptionRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		// check subscription exists
		var subscription models.Subscription
		err := database.DB.Where("id = ?", request.ID).First(&subscription).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.subscription_not_found",
				"success": false,
			})
			return
		}

		err = database.DB.Delete(&subscription).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"success": false,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func UpdateSubscription() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request UpdateSubscriptionRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		// check subscription exists
		var subscription models.Subscription
		err := database.DB.Where("id = ?", request.ID).First(&subscription).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.subscription_not_found",
				"success": false,
			})
			return
		}

		err = database.DB.Model(&subscription).Updates(map[string]interface{}{
			"name":  request.Name,
			"price": decimal.NewFromFloat(request.Price),
			"hours": request.Hours,
			"role": request.Role,
			"hidden": request.Hidden,
		}).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"success": false,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func NewSubscription() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request NewSubscriptionRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		subscription := models.Subscription{
			Name:  request.Name,
			Price: decimal.NewFromFloat(request.Price),
			Hours: request.Hours,
			Role: request.Role,
			Hidden: request.Hidden,
		}
		err := database.DB.Create(&subscription).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"success": false,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func PauseSubscription() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request PauseSubscriptionRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		user, err := utils.CurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error.user_not_found",
				"success": false,
			})
			return
		}

		// check subscription exists
		var subscription models.UserSubscription
		err = database.DB.Where("id = ?", request.ID).First(&subscription).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.subscription_not_found",
				"success": false,
			})
			return
		}

		// check it's state
		if subscription.Status != models.USER_SUBSCRIPTION_STATUS_ACTIVE {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.subscription_already_paused",
				"success": false,
			})
			return
		}

		database.DB.Create(&models.UserSubscriptionChunk{
			UserSubscriptionID: subscription.ID,
			Action:             models.USER_SUBSCRIPTION_ACTION_PAUSED,
			ActionBy:           &user.ID,
			CreatedAt:          time.Now(),
		})

		database.DB.Model(&subscription).Update("status", models.USER_SUBSCRIPTION_STATUS_PAUSED)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func UnpauseSubscription() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request PauseSubscriptionRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		user, err := utils.CurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "error.user_not_found",
				"success": false,
			})
			return
		}

		// check subscription exists
		var subscription models.UserSubscription
		err = database.DB.Where("id = ?", request.ID).First(&subscription).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.subscription_not_found",
				"success": false,
			})
			return
		}

		// check it's state
		if subscription.Status != models.USER_SUBSCRIPTION_STATUS_PAUSED {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.subscription_already_paused",
				"success": false,
			})
			return
		}

		database.DB.Create(&models.UserSubscriptionChunk{
			UserSubscriptionID: subscription.ID,
			Action:             models.USER_SUBSCRIPTION_ACTION_ACTIVATED,
			ActionBy:           &user.ID,
			CreatedAt:          time.Now(),
		})

		database.DB.Model(&subscription).Update("status", models.USER_SUBSCRIPTION_STATUS_ACTIVE)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}
