package payment

import (
	"fmt"
	"net/http"
	"time"

	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/middlewares"
	"github.com/berkaycubuk/billiard_software_api/models"
	"github.com/berkaycubuk/billiard_software_api/pkg/gameuser"
	orderModule "github.com/berkaycubuk/billiard_software_api/pkg/order"
	"github.com/berkaycubuk/billiard_software_api/pkg/vipps"
	"github.com/berkaycubuk/billiard_software_api/utils"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type InitializeRequest struct {
	OrderID uint  `json:"order_id" binding:"required" validate:"required"`
	Method  uint8 `json:"method" binding:"required" validate:"required"`
}

type StatusRequest struct {
	OrderID uint  `json:"order_id" binding:"required" validate:"required"`
	Method  uint8 `json:"method" binding:"required" validate:"required"`
}

type OverviewRequest struct {
	ID uint `uri:"id" binding:"required"`
}

type CaptureRequest struct {
	Reference string `uri:"reference" binding:"required"`
}

func initialize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request InitializeRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		// check order exists
		var order models.Order
		err := database.DB.Where("id = ?", request.OrderID).Preload("Detail").Find(&order).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.order_not_found",
				"success": false,
			})
			return
		}

		if order.Status == models.ORDER_STATUS_PAID || order.Status == models.ORDER_STATUS_APPROVED {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.order_completed",
				"success": false,
			})
			return
		}

		// update order price for game
		var orderItem models.OrderItem
		err = database.DB.Where("order_id = ? AND product_type = ?", order.ID, models.PRODUCT_TYPE_GAME).First(&orderItem).Error
		if err == nil { // found it
			// check is there game
			var game models.Game
			err = database.DB.Where("id = ? AND ended_at IS NULL", orderItem.ProductID).Find(&game).Error
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "errors.table_game_not_found",
					"success": false,
				})
				return
			}

			// is it user or guest
			if order.UserID != nil { // user
				var gameUser models.GameUser
				err = database.DB.Where("user_id = ? AND game_id = ? AND ended_at IS NULL", order.UserID, game.ID).Find(&gameUser).Error
				if err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"message": "errors.game_user_not_found",
						"success": false,
					})
					return
				}

				newOrderPrice := decimal.Zero
				foundGameOrder := false

				// check game user orders
				var gameuserOrders []models.GameUserOrder
				err = database.DB.Where("game_user_id = ?", gameUser.ID).Find(&gameuserOrders).Error
				if err == nil {
					foundGameOrder = true
					for _, v := range gameuserOrders {
						var gameOrder models.Order
						err = database.DB.Where("id = ? and status = ?", v.OrderID, models.ORDER_STATUS_PAY_LATER).First(&gameOrder).Error
						if err == nil {
							newOrderPrice = newOrderPrice.Add(gameOrder.Price)
						}
					}
				}

				// check orders added to account without game
				var accountOrders []models.Order
				err = database.DB.Where("user_id = ? and status = ?", order.UserID, models.ORDER_STATUS_PAY_LATER).Find(&accountOrders).Error
				if err == nil {
					if foundGameOrder {
						for _, v := range accountOrders {
							found := false
							for _, j := range gameuserOrders {
								if j.OrderID == v.ID {
									found = true
									break
								}
							}

							if !found {
								newOrderPrice = newOrderPrice.Add(v.Price)
							}
						}
					} else {
						for _, v := range accountOrders {
							newOrderPrice = newOrderPrice.Add(v.Price)
						}
					}
				}

				gamePrice := gameuser.CalcUserTotalPrice(game.ID, gameUser.ID)
				newOrderPrice = newOrderPrice.Add(gamePrice)

				database.DB.Model(&order).Update("price", newOrderPrice)
			} else { // guest
				var gameUser models.GameUser
				err = database.DB.Where("name = ? AND game_id = ? AND ended_at IS NULL", order.Detail.UserName, game.ID).Find(&gameUser).Error
				if err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"message": "errors.game_user_not_found",
						"success": false,
					})
					return
				}

				newOrderPrice := decimal.Zero

				// check game user orders
				var gameuserOrders []models.GameUserOrder
				err = database.DB.Where("game_user_id = ?", gameUser.ID).Find(&gameuserOrders).Error
				if err == nil {
					for _, v := range gameuserOrders {
						var gameOrder models.Order
						err = database.DB.Where("id = ? and status = ?", v.OrderID, models.ORDER_STATUS_PAY_LATER).First(&gameOrder).Error
						if err == nil {
							newOrderPrice = newOrderPrice.Add(gameOrder.Price)
						}
					}
				}

				gamePrice := gameuser.CalcUserTotalPrice(game.ID, gameUser.ID)
				newOrderPrice = newOrderPrice.Add(gamePrice)

				database.DB.Model(&order).Update("price", newOrderPrice)
			}
		}

		if request.Method == models.ORDER_METHOD_PHYSICAL {
			orderModule.CreateOrderHistory(order.ID, models.ORDER_ACTION_INIT, request.Method)
			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
			})
			return
		} else if request.Method == models.ORDER_METHOD_VIPPS {
			res, err := vipps.Initiate(order.ID, order.Price, order.Reference, "--description--")
			if err != nil {
				fmt.Println(err.Error())

				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "errors.payment_error",
					"success": false,
				})
				return
			}
			orderModule.CreateOrderHistory(order.ID, models.ORDER_ACTION_INIT, request.Method)

			ctx.JSON(http.StatusOK, gin.H{
				"token":                 res.Token,
				"checkout_frontend_url": res.CheckoutFrontendUrl,
				"polling_url":           res.PollingUrl,
				"success":               true,
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors.payment_method_not_found",
			"success": false,
		})
	}
}

func accept() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request InitializeRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		// check order exists
		var order models.Order
		err := database.DB.Where("id = ?", request.OrderID).Find(&order).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.order_not_found",
				"success": false,
			})
			return
		}

		if request.Method == models.ORDER_METHOD_PHYSICAL {
			//orderModule.UpdateOrderStatus(&order, models.ORDER_STATUS_PAID)
			orderModule.CreateOrderHistory(order.ID, models.ORDER_ACTION_INIT, request.Method)
			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
			})
			return
		} else if request.Method == models.ORDER_METHOD_VIPPS {
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors.payment_method_not_found",
			"success": false,
		})
	}
}

func status() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request StatusRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		// check order exists
		var order models.Order
		err := database.DB.Where("id = ?", request.OrderID).Find(&order).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.order_not_found",
				"success": false,
			})
			return
		}

		// check is it zero or not
		if order.Price.IsZero() {
			ctx.JSON(http.StatusOK, gin.H{
				"status":  models.ORDER_STATUS_PAID,
				"success": true,
			})
			return
		}

		var method uint8
		var history models.OrderHistory
		err = database.DB.Where("order_id = ? AND action = ?", order.ID, models.ORDER_ACTION_INIT).Last(&history).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.payment_not_initialized",
				"success": false,
			})
			return
		} else {
			method = history.Method
		}

		if method == models.ORDER_METHOD_PHYSICAL {
			orderModule.UpdateOrderStatus(&order, models.ORDER_STATUS_PAID, models.ORDER_ACTION_PAY, models.ORDER_METHOD_PHYSICAL)

			ctx.JSON(http.StatusOK, gin.H{
				"status":  order.Status,
				"success": true,
			})
			return
		} else if method == models.ORDER_METHOD_VIPPS {
			res, err := vipps.GetPayment(order.Reference)
			if err != nil {
				fmt.Println(err.Error())
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  order.Status,
					"success": false,
				})
				return
			}

			// update order status
			if res.State == vipps.PAYMENT_COMPLETED {
				orderModule.UpdateOrderStatus(&order, models.ORDER_STATUS_PAID, models.ORDER_ACTION_PAY, models.ORDER_METHOD_VIPPS)

				CapturePayment(&order)
			} else if res.State == vipps.PAYMENT_CANCELED {
				orderModule.UpdateOrderStatus(&order, models.ORDER_STATUS_CANCELED, models.ORDER_ACTION_CANCEL, models.ORDER_METHOD_VIPPS)
			} else if res.State == vipps.PAYMENT_WAITING {
			}

			ctx.JSON(http.StatusOK, gin.H{
				"status":  order.Status,
				"success": true,
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "errors.payment_method_not_found",
			"success": false,
		})
	}
}

func overview() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request OverviewRequest
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.id_required",
				"success": false,
			})
			return
		}

		var order models.Order
		err := database.DB.Where("id = ?", request.ID).Preload("Items").First(&order).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.order_not_found",
				"success": false,
			})
			return
		}

		if order.Price.IsZero() {
			if order.Status == models.ORDER_STATUS_PAID || order.Status == models.ORDER_STATUS_APPROVED {
				return
			}

			database.DB.Model(&order).Update("status", models.ORDER_STATUS_PAID)

			var items []models.OrderItem
			err := database.DB.Where("order_id = ?", order.ID).Find(&items).Error
			if err == nil {
				for _, v := range items {
					orderModule.HandleOrderItemAfterPayment(&order, &v)
				}
			}

			orderModule.CreateOrderHistory(order.ID, models.ORDER_ACTION_PAY, models.ORDER_METHOD_SYSTEM)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"order":   order,
			"success": true,
		})
	}
}

func captureReq() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request CaptureRequest
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.reference_required",
				"success": false,
			})
			return
		}

		var order models.Order
		err := database.DB.Where("reference = ?", request.Reference).First(&order).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"success": false,
			})
			return
		}

		err = vipps.Capture(order.Reference, order.Price)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
				"success": false,
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Success",
			"success": true,
		})
	}
}

func Routes(r *gin.RouterGroup) {
	routes := r.Group("/payment")
	{
		routes.Use(middlewares.AuthMiddleware())
		routes.POST("/init", initialize())
		routes.POST("/accept", accept())
		routes.GET("/overview/:id", overview())
		routes.GET("/capture/:reference", captureReq())
		routes.POST("/status", status())
	}
}

func CapturePayment(order *models.Order) {
	err := vipps.Capture(order.Reference, order.Price)	
	if err == nil {
		database.DB.Model(&order).Update("status", models.ORDER_STATUS_CAPTURED)
		fmt.Println("Captured")
		fmt.Println(order.ID)
	} else {
			fmt.Println(err.Error())
	}
}

func CapturePayments() {
	// get completed vipps order histories
	var histories []models.OrderHistory
	err := database.DB.Where("method = ? and action = ?", models.ORDER_METHOD_VIPPS, models.ORDER_ACTION_PAY).Find(&histories).Error
	if err == nil {
		for _, v := range histories {
			// find related order
			var order models.Order
			err = database.DB.Where("id = ? and status = ?", v.OrderID, models.ORDER_STATUS_PAID).Find(&order).Error
			if err == nil {
				CapturePayment(&order)
				// sleep for a second
				time.Sleep(time.Second)
			}
		}
	}

	/*
		var orders []models.Order
		err := database.DB.Where("status = ?", models.ORDER_STATUS_PAID).Find(&orders).Error
		if err == nil {
				for _, v := range orders {
						CapturePayment(&v)
				}
		}
		*/
}
