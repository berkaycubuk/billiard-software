package dashboard

import (
	"fmt"
	"net/http"
	"time"

	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/middlewares"
	"github.com/berkaycubuk/billiard_software_api/models"
	"github.com/berkaycubuk/billiard_software_api/pkg/gamemodule"
	"github.com/berkaycubuk/billiard_software_api/pkg/gameuser"
	ordermodule "github.com/berkaycubuk/billiard_software_api/pkg/order"
	"github.com/berkaycubuk/billiard_software_api/utils"
	"github.com/gin-gonic/gin"
)

func syncUserChunksForUnpause(gameUserID uint, gameID uint) error {
	// keep the current time
	now := time.Now()

	// fill price and ended_at values for current chunks
	var chunks []models.GameUserChunk
	database.DB.Where("game_id = ? AND ended_at IS NULL", gameID).Find(&chunks)
	for _, v := range chunks {
		price, err := gameuser.CalcUserChunkPrice(&v, now)
		if err != nil {
			return err
		}
		database.DB.Where("id = ?", v.ID).
			Updates(models.GameUserChunk{
				Price:   price,
				EndedAt: &now,
			})
	}

	// create new chunks
	var users []models.GameUser
	database.DB.Where("game_id = ? AND ended_at IS NULL", gameID).Find(&users)

	for _, v := range users {
		if v.ID == gameUserID {
			database.DB.Create(&models.GameUserChunk{
				GameID:     gameID,
				GameUserID: v.ID,
				Status:     models.GAME_USER_PLAYING,
				Players:    uint8(len(users)),
				StartedAt:  now,
			})
		} else {
			database.DB.Create(&models.GameUserChunk{
				GameID:     gameID,
				GameUserID: v.ID,
				Status:     v.Status,
				Players:    uint8(len(users)),
				StartedAt:  now,
			})
		}
	}

	return nil
}

func syncUserChunksForPause(gameUserID uint, gameID uint) error {
	// keep the current time
	now := time.Now()

	// fill price and ended_at values for current chunks
	var chunks []models.GameUserChunk
	database.DB.Where("game_id = ? AND ended_at IS NULL", gameID).Find(&chunks)
	for _, v := range chunks {
		price, err := gameuser.CalcUserChunkPrice(&v, now)
		if err != nil {
			return err
		}
		database.DB.Where("id = ?", v.ID).
			Updates(models.GameUserChunk{
				Price:   price,
				EndedAt: &now,
			})
	}

	// create new chunks
	var users []models.GameUser
	database.DB.Where("game_id = ? AND ended_at IS NULL", gameID).Find(&users)

	for _, v := range users {
		if v.ID == gameUserID {
			database.DB.Create(&models.GameUserChunk{
				GameID:     gameID,
				GameUserID: v.ID,
				Status:     models.GAME_USER_PAUSED,
				Players:    uint8(len(users)),
				StartedAt:  now,
			})
		} else {
			database.DB.Create(&models.GameUserChunk{
				GameID:     gameID,
				GameUserID: v.ID,
				Status:     v.Status,
				Players:    uint8(len(users)),
				StartedAt:  now,
			})
		}
	}

	return nil
}

func activeGames() gin.HandlerFunc {
				return func(ctx *gin.Context) {
								var games []models.Game
								err := database.DB.Where("ended_at IS NULL").Preload("GameUsers").Find(&games).Error
								if err != nil {
												ctx.JSON(http.StatusOK, gin.H{
																"games": nil,
												})
												return
								}

								ctx.JSON(http.StatusOK, gin.H{
												"games": games,
								})
				}
}

func kickUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request LeaveTableRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		var gameUser models.GameUser
				err := database.DB.Where("id = ?", request.ID).Find(&gameUser).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.game_user_not_found",
				"success": false,
			})
			return
		}

		database.DB.Model(&gameUser).Where("id = ?", gameUser.ID).Update("Status", models.GAME_USER_PAUSED)

		syncUserChunksForPause(gameUser.ID, gameUser.GameID)

		gamemodule.CreateHistory(gameUser.ID, gameUser.GameID, models.GAME_HISTORY_PAUSE)
		price := gameuser.CalcUserTotalPrice(gameUser.GameID, gameUser.ID)

		if gameUser.UserID != nil {
			order, err := ordermodule.Create(gameUser.UserID, price, models.ORDER_STATUS_WAITING)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "errors.unable_to_create_order",
					"success": false,
				})
				return
			}

			var user models.User
			err = database.DB.Where("id = ?", gameUser.UserID).First(&user).Error
			if err != nil {
						ctx.JSON(http.StatusInternalServerError, gin.H{
							"message": "errors.user_not_found",
							"success": false,
						})
			}

			_, _ = ordermodule.CreateOrderDetails(order.ID, user.Name, user.Surname, user.Phone, user.Email)
			_, _ = ordermodule.CreateOrderItem(order.ID, models.PRODUCT_TYPE_GAME, &gameUser.GameID, "game", 1, order.Price)

			database.DB.Model(&order).Update("status", models.ORDER_STATUS_PAID)

			var items []models.OrderItem
			err = database.DB.Where("order_id = ?", order.ID).Find(&items).Error
			if err == nil {
				for _, v := range items {
					ordermodule.HandleOrderItemAfterPayment(order, &v)
				}
			}

			ordermodule.CreateOrderHistory(order.ID, models.ORDER_ACTION_PAY, models.ORDER_METHOD_SYSTEM)

			ctx.JSON(http.StatusOK, gin.H{
				"order":   order,
				"success": true,
			})
					return
		} else {
			order, err := ordermodule.Create(nil, price, models.ORDER_STATUS_WAITING)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "errors.unable_to_create_order",
					"success": false,
				})
				return
			}
			_, _ = ordermodule.CreateOrderDetails(order.ID, gameUser.Name, "", "", "")
			_, _ = ordermodule.CreateOrderItem(order.ID, models.PRODUCT_TYPE_GAME, &gameUser.GameID, "game", 1, order.Price)

			database.DB.Model(&order).Update("status", models.ORDER_STATUS_PAID)

			var items []models.OrderItem
			err = database.DB.Where("order_id = ?", order.ID).Find(&items).Error
			if err == nil {
				for _, v := range items {
					ordermodule.HandleOrderItemAfterPayment(order, &v)
				}
			}

			ordermodule.CreateOrderHistory(order.ID, models.ORDER_ACTION_PAY, models.ORDER_METHOD_SYSTEM)

			ctx.JSON(http.StatusOK, gin.H{
				"order":   order,
				"success": true,
			})
					return
		}
	}
}

type PauseUnpauseGameRequest struct {
	ID uint `json:"id" binding:"required" validate:"required"`
}

func pauseUnpauseGame() gin.HandlerFunc {
				return func(ctx *gin.Context) {
						var request PauseUnpauseGameRequest
						if !utils.ValidateRequest(ctx, &request) {
							return
						}

							var gameUser models.GameUser
								err := database.DB.Where("id = ?", request.ID).First(&gameUser).Error
								if err != nil {
												ctx.JSON(http.StatusBadRequest, gin.H{
																"message": "Game user not found",
																"success": false,
												})
												return
								}

								if gameUser.Status == models.GAME_USER_PLAYING {
														database.DB.Model(&gameUser).Where("id = ?", gameUser.ID).Update("Status", models.GAME_USER_PAUSED)
												syncUserChunksForPause(gameUser.ID, gameUser.GameID)
												gamemodule.CreateHistory(gameUser.ID, gameUser.GameID, models.GAME_HISTORY_PAUSE)
								} else {
														database.DB.Model(&gameUser).Where("id = ?", gameUser.ID).Update("Status", models.GAME_USER_PLAYING)
												syncUserChunksForUnpause(gameUser.ID, gameUser.GameID)
												gamemodule.CreateHistory(gameUser.ID, gameUser.GameID, models.GAME_HISTORY_UNPAUSE)
								}

										ctx.JSON(http.StatusOK, gin.H{
											"success": true,
										})
				}
}

type LeaveTableRequest struct {
	ID uint `json:"id" binding:"required" validate:"required"`
}

func leaveTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request LeaveTableRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		var gameUser models.GameUser
				err := database.DB.Where("id = ?", request.ID).Find(&gameUser).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.game_user_not_found",
				"success": false,
			})
			return
		}

		price := gameuser.CalcUserTotalPrice(gameUser.GameID, gameUser.ID)

				// determine guest or user
				if gameUser.UserID != nil {
						order, err := ordermodule.Create(gameUser.UserID, price, models.ORDER_STATUS_WAITING)
						if err != nil {
							ctx.JSON(http.StatusInternalServerError, gin.H{
								"message": "errors.unable_to_create_order",
								"success": false,
							})
							return
						}

								var user models.User
								err = database.DB.Where("id = ?", gameUser.UserID).First(&user).Error
								if err != nil {
											ctx.JSON(http.StatusInternalServerError, gin.H{
												"message": "errors.user_not_found",
												"success": false,
											})
								}

						_, _ = ordermodule.CreateOrderDetails(order.ID, user.Name, user.Surname, user.Phone, user.Email)

						_, _ = ordermodule.CreateOrderItem(order.ID, models.PRODUCT_TYPE_GAME, &gameUser.GameID, "game", 1, order.Price)

						// fetch pay later orders
						var gameUserOrders []models.GameUserOrder
						err = database.DB.Where("game_user_id = ?", gameUser.ID).Find(&gameUserOrders).Error
						if err != nil {
							fmt.Println(err.Error())
						} else {
							for _, v := range gameUserOrders {
								ordermodule.TransferPayLater(v.OrderID, order.ID)
							}
						}

						// create game_user_order
						newGameUserOrder := models.GameUserOrder{
							GameUserID: gameUser.ID,
							OrderID:    order.ID,
						}
						database.DB.Create(&newGameUserOrder)

								if order.Price.IsZero() {

											database.DB.Model(&order).Update("status", models.ORDER_STATUS_PAID)

											var items []models.OrderItem
											err := database.DB.Where("order_id = ?", order.ID).Find(&items).Error
											if err == nil {
												for _, v := range items {
													ordermodule.HandleOrderItemAfterPayment(order, &v)
												}
											}

											ordermodule.CreateOrderHistory(order.ID, models.ORDER_ACTION_PAY, models.ORDER_METHOD_SYSTEM)
								}

						ctx.JSON(http.StatusOK, gin.H{
							"order":   order,
							"success": true,
						})
								return
				} else {
						order, err := ordermodule.Create(nil, price, models.ORDER_STATUS_WAITING)
						if err != nil {
							ctx.JSON(http.StatusInternalServerError, gin.H{
								"message": "errors.unable_to_create_order",
								"success": false,
							})
							return
						}

						_, _ = ordermodule.CreateOrderDetails(order.ID, gameUser.Name, "", "", "")

						_, _ = ordermodule.CreateOrderItem(order.ID, models.PRODUCT_TYPE_GAME, &gameUser.GameID, "game", 1, order.Price)

						// fetch pay later orders
						var gameUserOrders []models.GameUserOrder
						err = database.DB.Where("game_user_id = ?", gameUser.ID).Find(&gameUserOrders).Error
						if err != nil {
							fmt.Println(err.Error())
						} else {
							for _, v := range gameUserOrders {
								ordermodule.TransferPayLater(v.OrderID, order.ID)
							}
						}

								if order.Price.IsZero() {

											database.DB.Model(&order).Update("status", models.ORDER_STATUS_PAID)

											var items []models.OrderItem
											err := database.DB.Where("order_id = ?", order.ID).Find(&items).Error
											if err == nil {
												for _, v := range items {
													ordermodule.HandleOrderItemAfterPayment(order, &v)
												}
											}

											ordermodule.CreateOrderHistory(order.ID, models.ORDER_ACTION_PAY, models.ORDER_METHOD_SYSTEM)
								}

						ctx.JSON(http.StatusOK, gin.H{
							"order":   order,
							"success": true,
						})
								return
				}
	}
}

func Routes(r *gin.RouterGroup) {
	r.Use(middlewares.AuthMiddleware())
	r.GET("/dashboard/active-games", activeGames())
	r.POST("/dashboard/pause-unpause", pauseUnpauseGame())
	r.POST("/dashboard/kick-user", kickUser())
	r.POST("/dashboard/leave-table", leaveTable())
}
