package shop

import (
	"fmt"
	"net/http"

	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/middlewares"
	"github.com/berkaycubuk/billiard_software_api/models"
	orderModule "github.com/berkaycubuk/billiard_software_api/pkg/order"
	"github.com/berkaycubuk/billiard_software_api/utils"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type CartItem struct {
	ID	uint	`json:"id"`
	Count	int	`json:"count"`
}

type AddToGuestAccountRequest struct {
	GameUserName	string	`json:"game_user_name"`
	Items	[]CartItem	`json:"items" binding:"required" validate:"required"`
}

func addToGuestAccount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request AddToGuestAccountRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		var gameuser models.GameUser
		err := database.DB.Where("name = ?", request.GameUserName).First(&gameuser).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.user_not_found",
				"success": false,
			})
			return
		}

		products := []struct {
			Product 	models.Product
			Count		int
		}{}
		totalPrice := decimal.New(0, -2)
		for _, v := range request.Items {
			var product models.Product
			err := database.DB.Where("id = ?", v.ID).First(&product).Error
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			totalPrice = totalPrice.Add(product.Price.Mul(decimal.New(int64(v.Count), 0)))

			products = append(products, struct{Product models.Product; Count int}{
				Product: product,
				Count: v.Count,
			})
		}

		if totalPrice.Cmp(decimal.New(0, -2)) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.items_not_found",
				"success": false,
			})
			return
		}

		order, err := orderModule.Create(nil, totalPrice, models.ORDER_STATUS_PAY_LATER)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.unable_to_create_order",
				"success": false,
			})
			return
		}

		for _, v := range products {
			orderModule.CreateOrderItem(order.ID, models.PRODUCT_TYPE_SHOP, &v.Product.ID, v.Product.Name, v.Count, v.Product.Price)
		}

		orderModule.CreateOrderDetails(order.ID, gameuser.Name, "", "", "")

		database.DB.Create(&models.GameUserOrder{
			OrderID: order.ID,
			GameUserID: gameuser.ID,
		})

		ctx.JSON(http.StatusOK, gin.H{
			"order_id": order.ID,
			"success": true,
		})
	}
}

type AddToAccountRequest struct {
	Items	[]CartItem	`json:"items" binding:"required" validate:"required"`
}

func addToAccount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request AddToAccountRequest
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

		products := []struct {
			Product 	models.Product
			Count		int
		}{}
		totalPrice := decimal.New(0, -2)
		for _, v := range request.Items {
			var product models.Product
			err := database.DB.Where("id = ?", v.ID).First(&product).Error
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			totalPrice = totalPrice.Add(product.Price.Mul(decimal.New(int64(v.Count), 0)))

			products = append(products, struct{Product models.Product; Count int}{
				Product: product,
				Count: v.Count,
			})
		}

		if totalPrice.Cmp(decimal.New(0, -2)) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.items_not_found",
				"success": false,
			})
			return
		}

		order, err := orderModule.Create(&user.ID, totalPrice, models.ORDER_STATUS_PAY_LATER)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.unable_to_create_order",
				"success": false,
			})
			return
		}

		for _, v := range products {
			orderModule.CreateOrderItem(order.ID, models.PRODUCT_TYPE_SHOP, &v.Product.ID, v.Product.Name, v.Count, v.Product.Price)
		}

		orderModule.CreateOrderDetails(order.ID, user.Name, user.Surname, user.Phone, user.Email)

		// if game user active add order to it
		gameUser, err := findGameUser(user.ID)
		if gameUser != nil {
			database.DB.Create(&models.GameUserOrder{
				OrderID: order.ID,
				GameUserID: gameUser.ID,
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"order_id": order.ID,
			"success": true,
		})
	}
}

type BuyRequest struct {
	Items	[]CartItem	`json:"items" binding:"required" validate:"required"`
}

func buy() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request BuyRequest
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

		products := []struct {
			Product 	models.Product
			Count		int
		}{}
		totalPrice := decimal.New(0, -2)
		for _, v := range request.Items {
			var product models.Product
			err := database.DB.Where("id = ?", v.ID).First(&product).Error
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			totalPrice = totalPrice.Add(product.Price.Mul(decimal.New(int64(v.Count), 0)))

			products = append(products, struct{Product models.Product; Count int}{
				Product: product,
				Count: v.Count,
			})
		}

		if totalPrice.Cmp(decimal.New(0, -2)) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.items_not_found",
				"success": false,
			})
			return
		}

		order, err := orderModule.Create(&user.ID, totalPrice, models.ORDER_STATUS_WAITING)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.unable_to_create_order",
				"success": false,
			})
			return
		}

		for _, v := range products {
			orderModule.CreateOrderItem(order.ID, models.PRODUCT_TYPE_SHOP, &v.Product.ID, v.Product.Name, v.Count, v.Product.Price)
		}

		orderModule.CreateOrderDetails(order.ID, user.Name, user.Surname, user.Phone, user.Email)

		// if game user active add order to it
		gameUser, err := findGameUser(user.ID)
		if gameUser != nil {
			database.DB.Create(&models.GameUserOrder{
				OrderID: order.ID,
				GameUserID: gameUser.ID,
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"order_id": order.ID,
			"success": true,
		})
	}
}

func findGameUser(userID uint) (*models.GameUser, error) {
	var gameUser models.GameUser
	err := database.DB.Where("user_id = ? AND ended_at IS NULL", userID).First(&gameUser).Error
	if err != nil {
		return nil, err
	}

	return &gameUser, nil
}

func Routes(r *gin.RouterGroup) {
	routes := r.Group("/shop")
	{
		routes.Use(middlewares.AuthMiddleware())
		routes.POST("/buy", buy())
		routes.POST("/add-to-account", addToAccount())
		routes.POST("/add-to-guest-account", addToGuestAccount())
	}
}
