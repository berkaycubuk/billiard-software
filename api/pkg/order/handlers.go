package order

import (
	"net/http"
	"time"

	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/middlewares"
	"github.com/berkaycubuk/billiard_software_api/models"
	"github.com/berkaycubuk/billiard_software_api/utils"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type GetOrderRequest struct {
	ID uint `uri:"id" binding:"required"`
}

type DeleteOrderRequest struct {
	ID uint `uri:"id" binding:"required"`
}

func getAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		countChain := database.DB
		chain := database.DB

		start := ctx.Query("start")
		end := ctx.Query("end")
		zeroOrders := ctx.Query("zeroOrders")
		method := ctx.Query("method")

		var ordersTemp []models.Order
		var orders []models.Order

		if method != "" {
			chain = chain.Joins("JOIN order_histories ON orders.id = order_histories.order_id").
				Where("order_histories.method = ?", method)
			countChain = countChain.Joins("JOIN order_histories ON orders.id = order_histories.order_id").
				Where("order_histories.method = ?", method)

			if start != "" {
				chain = chain.Where("orders.created_at >= ?", start)
				countChain = countChain.Where("orders.created_at >= ?", start)
			}

			if end != "" {
				chain = chain.Where("orders.created_at <= ?", end+"T23:59:59")
				countChain = countChain.Where("orders.created_at <= ?", end+"T23:59:59")
			}

			if !(zeroOrders != "" && zeroOrders == "true") {
				chain = chain.Where("orders.price != ?", decimal.Zero)
				countChain = countChain.Where("orders.price != ?", decimal.Zero)
			}
			countChain.Where("orders.status != ?", models.ORDER_STATUS_TRANSFERRED).Select("id").Find(&ordersTemp)
			chain.Where("orders.status != ?", models.ORDER_STATUS_TRANSFERRED).Scopes(utils.Paginate(ctx)).Order("id desc").Preload("Detail").Find(&orders)
		} else {
			if start != "" {
				chain = chain.Where("created_at >= ?", start)
				countChain = countChain.Where("created_at >= ?", start)
			}

			if end != "" {
				chain = chain.Where("created_at <= ?", end+"T23:59:59")
				countChain = countChain.Where("created_at <= ?", end+"T23:59:59")
			}

			if !(zeroOrders != "" && zeroOrders == "true") {
				chain = chain.Where("price != ?", decimal.Zero)
				countChain = countChain.Where("price != ?", decimal.Zero)
			}
			countChain.Where("status != ?", models.ORDER_STATUS_TRANSFERRED).Select("id").Find(&ordersTemp)
			chain.Where("status != ?", models.ORDER_STATUS_TRANSFERRED).Scopes(utils.Paginate(ctx)).Order("id desc").Preload("Detail").Find(&orders)
		}

		count := len(ordersTemp)

		pageSize := count / 10
		if pageSize*10 < count {
			pageSize = pageSize + 1
		}

		ctx.Set("response", gin.H{
			"orders": orders,
			"pagination": gin.H{
				"size": pageSize,
			},
		})
	}
}

func activeOrders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var orders []models.Order

		database.DB.Where("status = ?", models.ORDER_STATUS_WAITING).Order("id desc").Preload("Detail").Find(&orders)

		ctx.Set("response", gin.H{
			"orders": orders,
		})
	}
}

func getMyOldOrders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := utils.CurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.user_not_found",
				"success": false,
			})
			return
		}

		countChain := database.DB
		chain := database.DB

		start := ctx.Query("start")
		if start != "" {
			chain = chain.Where("created_at >= ?", start)
			countChain = countChain.Where("created_at >= ?", start)
		}

		end := ctx.Query("end")
		if end != "" {
			chain = chain.Where("created_at <= ?", end+"T23:59:59")
			countChain = countChain.Where("created_at <= ?", end+"T23:59:59")
		}

		var ordersTemp []models.Order
		countChain.Where("user_id = ? AND status IN (?)", user.ID, []int{models.ORDER_STATUS_PAID, models.ORDER_STATUS_APPROVED}).Select("id").Find(&ordersTemp)
		count := len(ordersTemp)

		pageSize := count / 10
		if pageSize*10 < count {
			pageSize = pageSize + 1
		}

		var orders []models.Order

		chain.Where("user_id = ? AND status IN (?)", user.ID, []int{models.ORDER_STATUS_PAID, models.ORDER_STATUS_APPROVED}).Scopes(utils.Paginate(ctx)).Order("id desc").Preload("Detail").Find(&orders)

		ctx.Set("response", gin.H{
			"orders": orders,
			"pagination": gin.H{
				"size": pageSize,
			},
		})
	}
}

func getMyWaitingOrders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := utils.CurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.user_not_found",
				"success": false,
			})
			return
		}

		var orders []models.Order
		err = database.DB.Where("user_id = ? AND status IN (?)", user.ID, []uint8{models.ORDER_STATUS_PAY_LATER}).Order("id DESC").Preload("Items").Find(&orders).Error
		if err != nil {
			ctx.Set("response", nil)
			return
		}

		ctx.Set("response", orders)
	}
}

func get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request GetOrderRequest
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.id_required",
				"success": false,
			})
			return
		}

		order := Get(request.ID)
		if order.Price.IsZero() {
			if order.Status == models.ORDER_STATUS_PAID || order.Status == models.ORDER_STATUS_APPROVED {
				ctx.JSON(http.StatusOK, gin.H{
					"order":   order,
					"success": true,
				})
				return
			}

			database.DB.Model(&order).Update("status", models.ORDER_STATUS_PAID)

			var items []models.OrderItem
			err := database.DB.Where("order_id = ?", order.ID).Find(&items).Error
			if err == nil {
				for _, v := range items {
					HandleOrderItemAfterPayment(order, &v)
				}
			}

			CreateOrderHistory(order.ID, models.ORDER_ACTION_PAY, models.ORDER_METHOD_SYSTEM)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"order":   order,
			"success": true,
		})
	}
}

func approve() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request DeleteOrderRequest
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.id_required",
				"success": false,
			})
			return
		}

		err := Approve(request.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
			})
		}
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func del() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request DeleteOrderRequest
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.id_required",
				"success": false,
			})
			return
		}

		Delete(request.ID)
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func cancel() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request DeleteOrderRequest
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.id_required",
				"success": false,
			})
			return
		}

		var order models.Order
		err := database.DB.Where("id = ?", request.ID).First(&order).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.order_not_found",
				"success": false,
			})
			return
		}

		// check is order have an game item
		// var orderItem models.OrderItem
		// err = database.DB.Where("order_id = ? AND product_type = ?", order.ID, models.PRODUCT_TYPE_GAME).First(&orderItem).Error
		// if err == nil { // that means this order contains payment for a game
		// 	// unpause paused user
		// }

		now := time.Now()

		err = database.DB.Model(&order).Updates(models.Order{
			Status:    models.ORDER_STATUS_CANCELED,
			DeletedAt: &now,
		}).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.unable_to_update_order",
				"success": false,
			})
			return
		}

		CreateOrderHistory(order.ID, models.ORDER_ACTION_CANCEL, models.ORDER_METHOD_SYSTEM)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

type ApplyDiscountRequest struct {
	ID	uint	`json:"id" binding:"required" validate:"required"`
	Amount	decimal.Decimal	`json:"amount" binding:"required" validate:"required"`
}

func applyDiscount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request ApplyDiscountRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		var order models.Order
		err := database.DB.Where("id = ?", request.ID).First(&order).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.order_not_found",
				"success": false,
			})
			return
		}

		orderDiscount := models.OrderDiscount{
			OrderID: request.ID,
			Type: models.DISCOUNT_TYPE_PRICE,
			DiscountPrice: request.Amount,
		}

		err = database.DB.Create(&orderDiscount).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"success": false,
			})
			return
		}

		// update order price
		database.DB.Model(&order).Update("price", order.Price.Sub(request.Amount))

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

type GetUserOrdersRequest struct {
	ID	uint	`uri:"id" binding:"required"`
}

func getUserOrders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request GetUserOrdersRequest
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.id_required",
				"success": false,
			})
			return
		}

		var user models.User
		err := database.DB.Where("id = ?", request.ID).First(&user).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.user_not_found",
				"success": false,
			})
			return
		}

		var orders []models.Order
		err = database.DB.Where("user_id = ?", user.ID).Order("id DESC").Find(&orders).Error
		if err != nil {
			ctx.Set("response", nil)
			return
		}

		ctx.Set("response", gin.H{
			"orders": orders,
		})
	}
}

func Routes(r *gin.RouterGroup) {
	r.Use(middlewares.AuthMiddleware())
	r.GET("/order/all", getAll())
	r.GET("/order/active-orders", activeOrders())
	r.GET("/order/user-orders/:id", getUserOrders())
	r.GET("/order/my-old-orders", getMyOldOrders())
	r.GET("/order/my-waiting-orders", getMyWaitingOrders())
	r.POST("/order/apply-discount", applyDiscount())
	r.POST("/order/cancel/:id", cancel())
	r.POST("/order/delete/:id", del())
	r.POST("/order/approve/:id", approve())
	r.GET("/order/:id", get())
}
