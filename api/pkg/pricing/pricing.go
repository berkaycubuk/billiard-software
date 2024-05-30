package pricing

import (
	"fmt"
	"log"
	"net/http"

	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/middlewares"
	"github.com/berkaycubuk/billiard_software_api/models"
	"github.com/berkaycubuk/billiard_software_api/utils"
	"github.com/berkaycubuk/gofusion"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func PerMinutePricingForUser(gameUserID uint, playerCount uint8) (*decimal.Decimal, error) {
	var gameUser models.GameUser
	err := database.DB.Where("id = ?", gameUserID).First(&gameUser).Error
	if err != nil { // not found I guess
		log.Println("----- PerMinutePricingForUser() | where game_user id = ? ----")
		log.Println(err.Error())
		log.Println("-----")
		return nil, err
	}

	// new

	// get role
	userRoleID := uint(0)
	if gameUser.UserID != nil {
		var userRole models.UserRole
		err = database.DB.Where("user_id = ?", gameUser.UserID).First(&userRole).Error
		if err == nil {
			userRoleID = userRole.RoleID
		} else {
			log.Println("ERR: userRoleID section")
		}
	}

	userSubscriptionID := uint(0)
	if gameUser.UserID != nil {
		var userSubscription models.UserSubscription
		err = database.DB.Where("user_id = ? AND status = ?", gameUser.UserID, models.USER_SUBSCRIPTION_STATUS_ACTIVE).First(&userSubscription).Error
		if err == nil {
			if userSubscription.SubscriptionID == 1 {
				userSubscriptionID = 2
			} else {
				userSubscriptionID = userSubscription.SubscriptionID
			}
		}
	}

	var pricing models.Pricing
	if userRoleID == uint(0) {
		err = database.DB.Where("role_id IS NULL AND subscription_id IS NULL AND player_count = ?", playerCount).First(&pricing).Error
		if err != nil {
			log.Println("ERR: userRoleID == 0 section")
		}
	} else if userSubscriptionID == uint(0) {
		err = database.DB.Where("role_id = ? AND subscription_id IS NULL AND player_count = ?", userRoleID, playerCount).First(&pricing).Error
		if err != nil {
			err = database.DB.Where("role_id IS NULL AND player_count = ?", playerCount).First(&pricing).Error
		}
	} else {
		err = database.DB.Where("role_id = ? AND subscription_id = ? AND player_count = ?", userRoleID, userSubscriptionID, playerCount).First(&pricing).Error
		if err != nil {
			err = database.DB.Where("role_id = ? AND player_count = ?", userRoleID, playerCount).First(&pricing).Error
			if err != nil {
				err = database.DB.Where("role_id IS NULL AND player_count = ?", playerCount).First(&pricing).Error
			}
		}
	}

	if err != nil {
		log.Println("----- PerMinutePricingForUser() | final err check ----")
		log.Println(err.Error())
		log.Println(gameUserID, playerCount)
		log.Println("-----")
		return nil, err
	}

	return &pricing.PerMinute, nil
}

// return as role groups
func contains(value uint, array []uint) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}

	return false
}

func allPricing() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var pricings []models.Pricing
		err := database.DB.Preload("Subscription").Preload("Role").Order("id DESC").Find(&pricings).Error
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"pricings": nil,
				"success":  true,
			})
			return
		}

		// var roleIDs []uint
		// for _, v := range pricings {
		// 	if v.RoleID != nil {
		// 		if !contains(*v.RoleID, roleIDs) {
		// 			roleIDs = append(roleIDs, *v.RoleID)
		// 		}
		// 	}
		// }

		// var roles []models.Role
		// err = database.DB.Where("id IN ?", roleIDs).Find(&roles).Error
		// if err != nil {
		// 	ctx.JSON(http.StatusOK, gin.H{
		// 		"roles":   nil,
		// 		"success": true,
		// 	})
		// 	return
		// }

		ctx.JSON(http.StatusOK, gin.H{
			"pricings": pricings,
			"success":  true,
		})
	}
}

type GetPricingRequest struct {
	ID uint `uri:"id" bindin:"required"`
}

type SubscriptionGroup struct {
	SubscriptionID   uint    `json:"subscription_id"`
	SubscriptionName string  `json:"subscription_name"`
	Blocks           []Block `json:"blocks"`
}

type Block struct {
	ID          uint            `json:"id"`
	PlayerCount int             `json:"player_count"`
	PerMinute   decimal.Decimal `json:"per_minute"`
}

func getPricings() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request GetPricingRequest
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "ID required.",
				"success": false,
			})
			return
		}

		var pricings []models.Pricing
		err := database.DB.Where("role_id = ?", request.ID).Find(&pricings).Error
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"pricings": nil,
				"success":  true,
			})
			return
		}

		var subscriptionIDs []uint
		for _, v := range pricings {
			if v.SubscriptionID != nil {
				if !contains(*v.SubscriptionID, subscriptionIDs) {
					subscriptionIDs = append(subscriptionIDs, *v.SubscriptionID)
				}
			}
		}

		var subscriptionGroups []*SubscriptionGroup
		for _, v := range subscriptionIDs {
			// get subscription name
			var subscription models.Subscription
			err := database.DB.Where("id = ?", v).First(&subscription).Error

			// get blocks
			var temp []models.Pricing
			database.DB.Where("role_id = ? AND subscription_id = ?", request.ID, v).Find(&temp)

			var blocks []Block
			for _, i := range temp {
				blocks = append(blocks, Block{
					ID:          i.ID,
					PlayerCount: int(i.PlayerCount),
					PerMinute:   i.PerMinute,
				})
			}

			if err == nil {
				subscriptionGroups = append(subscriptionGroups, &SubscriptionGroup{
					SubscriptionID:   v,
					SubscriptionName: subscription.Name,
					Blocks:           blocks,
				})
			}
		}

		// add nil subscription_id pricings

		// get blocks
		var temp []models.Pricing
		database.DB.Where("role_id = ? AND subscription_id IS NULL", request.ID).Find(&temp)

		var blocks []Block
		for _, i := range temp {
			blocks = append(blocks, Block{
				ID:          i.ID,
				PlayerCount: int(i.PlayerCount),
				PerMinute:   i.PerMinute,
			})
		}

		subscriptionGroups = append(subscriptionGroups, &SubscriptionGroup{
			SubscriptionID:   0,
			SubscriptionName: "",
			Blocks:           blocks,
		})

		fmt.Println(subscriptionGroups)

		ctx.JSON(http.StatusOK, gin.H{
			"pricings": subscriptionGroups,
			"success":  true,
		})
	}
}

type CreatePricingRequest struct {
	RoleID         *uint            `json:"role_id" binding:"-" validate:"-"`
	SubscriptionID *uint            `json:"subscription_id" binding:"-" validate:"-"`
	PlayerCount    int             `json:"player_count" binding:"required" validate:"required"`
	PerMinute      decimal.Decimal `json:"per_minute" binding:"required" validate:"required"`
}

func createPricing() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request CreatePricingRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		pricing := &models.Pricing{
			RoleID:         request.RoleID,
			SubscriptionID: request.SubscriptionID,
			PlayerCount:    int8(request.PlayerCount),
			PerMinute:      request.PerMinute,
		}
		database.DB.Create(&pricing)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

type UpdatePricingRequest struct {
	ID             uint            `json:"id" binding:"required" validate:"required"`
	RoleID         uint            `json:"role_id" binding:"-" validate:"-"`
	SubscriptionID uint            `json:"subscription_id" binding:"-" validate:"-"`
	PlayerCount    int             `json:"player_count" binding:"required" validate:"required"`
	PerMinute      decimal.Decimal `json:"per_minute" binding:"required" validate:"required"`
}

func updatePricing() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request UpdatePricingRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		var pricing models.Pricing
		err := database.DB.Where("id = ?", request.ID).First(&pricing).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.pricing_not_found",
				"success": false,
			})
			return
		}

		database.DB.Model(&pricing).Updates(models.Pricing{
			RoleID:         &request.RoleID,
			SubscriptionID: &request.SubscriptionID,
			PlayerCount:    int8(request.PlayerCount),
			PerMinute:      request.PerMinute,
		})

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

type DeletePricingRequest struct {
	ID uint `json:"id" binding:"required" validate:"required"`
}

func deletePricing() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request DeletePricingRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		var pricing models.Pricing
		err := database.DB.Where("id = ?", request.ID).First(&pricing).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.pricing_not_found",
				"success": false,
			})
			return
		}

		database.DB.Delete(&pricing)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func Routes(r *gin.RouterGroup) {
	r.Use(middlewares.AuthMiddleware())
	r.GET("/pricing/all", allPricing())
	r.POST("/pricing/create", createPricing())
	r.POST("/pricing/delete", deletePricing())
	r.POST("/pricing/update", updatePricing())
	r.GET("/pricing/:id", gofusion.ValidateURIParam("id"), getPricings())
}
