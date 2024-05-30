package table

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/middlewares"
	"github.com/berkaycubuk/billiard_software_api/models"
	"github.com/berkaycubuk/billiard_software_api/pkg/gamemodule"
	"github.com/berkaycubuk/billiard_software_api/pkg/gameuser"
	ordermodule "github.com/berkaycubuk/billiard_software_api/pkg/order"
	"github.com/berkaycubuk/billiard_software_api/utils"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type GameUserCustom struct {
	ID     uint            `json:"id"`
	UserID *uint           `json:"user_id"`
	Name   string          `json:"name"`
	Price  decimal.Decimal `json:"price"`
	Time   float64         `json:"time"`
	Status uint8           `json:"status"`
}

type GetTablesResponse struct {
	ID        uint             `json:"id"`
	Name      string           `json:"name"`
	Status    uint8            `json:"status"`
	Game      *models.Game     `json:"game"`
	GameUsers []GameUserCustom `json:"game_users"`
}

func getTables() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var tables []models.Table
		err := database.DB.Find(&tables).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
			})
			return
		}

		var tablesResponse []GetTablesResponse
		for _, v := range tables {
			var tableGame models.Game
			err := database.DB.Where("table_id = ? AND ended_at IS NULL", v.ID).Preload("GameUsers", "ended_at IS NULL").First(&tableGame).Error
			if err != nil {
				tablesResponse = append(tablesResponse, GetTablesResponse{
					ID:     v.ID,
					Name:   v.Name,
					Status: v.Status,
				})
			} else {
				var gameUsersCustom []GameUserCustom
				for _, j := range tableGame.GameUsers {
					price, _ := decimal.NewFromString(gameuser.CalcUserTotalPrice(tableGame.ID, j.ID).StringFixed(2))
					gameUsersCustom = append(gameUsersCustom, GameUserCustom{
						ID:     j.ID,
						UserID: j.UserID,
						Name:   j.Name,
						Price:  price,
						Time:   gameuser.CalcUserTotalTime(tableGame.ID, j.ID),
						Status: j.Status,
					})
				}
				tablesResponse = append(tablesResponse, GetTablesResponse{
					ID:        v.ID,
					Name:      v.Name,
					Status:    v.Status,
					Game:      &tableGame,
					GameUsers: gameUsersCustom,
				})
			}

		}

		ctx.JSON(http.StatusOK, gin.H{
			"tables":  tablesResponse,
			"success": true,
		})
	}
}

// need to check is user have permission
func deleteTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request DeleteTableRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		err := delete(request.ID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func updateTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request UpdateTableRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		var table models.Table
		err := database.DB.Where("id = ?", request.ID).First(&table).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.table_not_found",
				"success": false,
			})
		}

		err = database.DB.Model(&table).Update("name", request.Name).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"success": false,
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func newTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request NewTableRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		err := create(request.Name)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"success": false,
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

type GetTableRequest struct {
	ID string `uri:"id" binding:"required"`
}

type GetTableResponse struct {
	ID        uint             `json:"id"`
	Name      string           `json:"name"`
	Status    uint8            `json:"status"`
	CreatedAt time.Time        `json:"created_at"`
	Game      *models.Game     `json:"game"`
	GameUsers []GameUserCustom `json:"game_users"`
}

func getTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request GetTableRequest
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.id_required",
				"success": false,
			})
			return
		}

		var table models.Table
		err := database.DB.Where("id = ?", request.ID).First(&table).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.table_not_found",
				"success": false,
			})
			return
		}

		var game models.Game
		err = database.DB.Where("table_id = ? AND ended_at IS NULL", table.ID).Preload("GameUsers", "ended_at IS NULL").Preload("GameHistories.GameUser").First(&game).Error
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"table": GetTableResponse{
					ID:        table.ID,
					Name:      table.Name,
					Status:    table.Status,
					CreatedAt: table.CreatedAt,
				},
				"success": true,
			})
			return
		}

		var gameUsersCustom []GameUserCustom
		for _, v := range game.GameUsers {
			price, _ := decimal.NewFromString(gameuser.CalcUserTotalPrice(game.ID, v.ID).StringFixed(2))
			gameUsersCustom = append(gameUsersCustom, GameUserCustom{
				ID:     v.ID,
				UserID: v.UserID,
				Name:   v.Name,
				Price:  price,
				Time:   gameuser.CalcUserTotalTime(game.ID, v.ID),
				Status: v.Status,
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"table": GetTableResponse{
				ID:        table.ID,
				Name:      table.Name,
				Status:    table.Status,
				CreatedAt: table.CreatedAt,
				Game:      &game,
				GameUsers: gameUsersCustom,
			},
			"success": true,
		})
	}
}

type LeaveAsGuestRequest struct {
	ID   uint   `json:"id" binding:"required" validate:"required"`
	Name string `json:"name" binding:"required" validate:"required"`
}

func leaveAsGuest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request LeaveAsGuestRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		_, err := utils.CurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.user_not_found",
				"success": false,
			})
			return
		}

		// check is there game
		var game models.Game
		err = database.DB.Where("table_id = ? AND ended_at IS NULL", request.ID).Find(&game).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.table_game_not_found",
				"success": false,
			})
			return
		}

		var gameUser models.GameUser
		err = database.DB.Where("name = ? AND game_id = ? AND ended_at IS NULL", request.Name, game.ID).Find(&gameUser).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.game_user_not_found",
				"success": false,
			})
			return
		}

		// do not leave immediately
		// after the payment confirmation, leave.
		//database.DB.Model(&gameUser).Where("id = ?", gameUser.ID).Update("Status", models.GAME_USER_PAUSED)

		//syncUserChunksForPause(gameUser.ID, game.ID)

		//gamemodule.CreateHistory(gameUser.ID, game.ID, models.GAME_HISTORY_PAUSE)

		// create order
		price := gameuser.CalcUserTotalPrice(game.ID, gameUser.ID)
		order, err := ordermodule.Create(nil, price, models.ORDER_STATUS_WAITING)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "errors.unable_to_create_order",
				"success": false,
			})
			return
		}

		_, _ = ordermodule.CreateOrderDetails(order.ID, request.Name, "", "", "")

		_, _ = ordermodule.CreateOrderItem(order.ID, models.PRODUCT_TYPE_GAME, &game.ID, "game", 1, order.Price)

		// fetch pay later orders
		var gameUserOrders []models.GameUserOrder
		err = database.DB.Where("game_user_id = ?", gameUser.ID).Find(&gameUserOrders).Error
		if err == nil {
			for _, v := range gameUserOrders {
				ordermodule.TransferPayLater(v.OrderID, order.ID)
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"order":   order,
			"success": true,
		})
	}
}

type LeaveTableRequest struct {
	ID uint `json:"id" binding:"required" validate:"required"`
}

func leaveTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request JoinTableRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		user, err := utils.CurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.user_not_found",
				"success": false,
			})
			return
		}

		// check is there game
		var game models.Game
		err = database.DB.Where("table_id = ? AND ended_at IS NULL", request.ID).Find(&game).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.table_game_not_found",
				"success": false,
			})
			return
		}

		var gameUser models.GameUser
		err = database.DB.Where("user_id = ? AND game_id = ? AND ended_at IS NULL", user.ID, game.ID).Find(&gameUser).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.game_user_not_found",
				"success": false,
			})
			return
		}

		// create order
		price := gameuser.CalcUserTotalPrice(game.ID, gameUser.ID)
		order, err := ordermodule.Create(&user.ID, price, models.ORDER_STATUS_WAITING)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "errors.unable_to_create_order",
				"success": false,
			})
			return
		}

		_, _ = ordermodule.CreateOrderDetails(order.ID, user.Name, user.Surname, user.Phone, user.Email)

		_, _ = ordermodule.CreateOrderItem(order.ID, models.PRODUCT_TYPE_GAME, &game.ID, "game", 1, order.Price)

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

		// plus check pay laters created without table
		var payLaterOrders []models.Order
		err = database.DB.Where("user_id = ? and status = ?", user.ID, models.ORDER_STATUS_PAY_LATER).Find(&payLaterOrders).Error
		if err == nil {
			for _, v := range payLaterOrders {
				// check is it already added to game order
				isAdded := false
				for _, j := range gameUserOrders {
					if j.OrderID == v.ID {
						isAdded = true
					}
				}

				if !isAdded {
					ordermodule.TransferPayLater(v.ID, order.ID)
				}
			}
		}

		// create game_user_order
		newGameUserOrder := models.GameUserOrder{
			GameUserID: gameUser.ID,
			OrderID:    order.ID,
		}
		database.DB.Create(&newGameUserOrder)

		ctx.JSON(http.StatusOK, gin.H{
			"order":   order,
			"success": true,
		})
	}
}

type UnpauseAsGuestRequest struct {
	ID   uint   `json:"id" binding:"required" validate:"required"`
	Name string `json:"name" binding:"required" validate:"required"`
}

func unpauseAsGuest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request UnpauseAsGuestRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		// check is there game already
		var game models.Game
		err := database.DB.Where("table_id = ? AND ended_at IS NULL", request.ID).First(&game).Error
		if err != nil { // I'll guess it's not found
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.game_not_found",
				"success": false,
			})
			return
		}

		var gameUser models.GameUser
		err = database.DB.Where("name = ? AND game_id = ? AND ended_at IS NULL", request.Name, game.ID).First(&gameUser).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.game_user_not_found",
				"success": false,
			})
			return
		}

		database.DB.Model(&gameUser).Where("id = ?", gameUser.ID).Update("Status", models.GAME_USER_PLAYING)

		syncUserChunksForUnpause(gameUser.ID, game.ID)

		gamemodule.CreateHistory(gameUser.ID, game.ID, models.GAME_HISTORY_UNPAUSE)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

type UnpauseTableRequest struct {
	ID uint `json:"id" binding:"required" validate:"required"`
}

func unpauseTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request JoinTableRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		user, err := utils.CurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.user_not_found",
				"success": false,
			})
			return
		}

		// check is there game already
		var game models.Game
		err = database.DB.Where("table_id = ? AND ended_at IS NULL", request.ID).First(&game).Error
		if err != nil { // I'll guess it's not found
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.game_not_found",
				"success": false,
			})
			return
		}

		var gameUser models.GameUser
		err = database.DB.Where("user_id = ? AND game_id = ? AND ended_at IS NULL", user.ID, game.ID).First(&gameUser).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.game_user_not_found",
				"success": false,
			})
			return
		}

		database.DB.Model(&gameUser).Where("id = ?", gameUser.ID).Update("Status", models.GAME_USER_PLAYING)

		syncUserChunksForUnpause(gameUser.ID, game.ID)

		gamemodule.CreateHistory(gameUser.ID, game.ID, models.GAME_HISTORY_UNPAUSE)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

type PauseAsGuestRequest struct {
	ID   uint   `json:"id" binding:"required" validate:"required"`
	Name string `json:"name" binding:"required" validate:"required"`
}

func pauseAsGuest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request PauseAsGuestRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		// check is there game already
		var game models.Game
		err := database.DB.Where("table_id = ? AND ended_at IS NULL", request.ID).First(&game).Error
		if err != nil { // I'll guess it's not found
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.game_not_found",
				"success": false,
			})
			return
		}

		var gameUser models.GameUser
		err = database.DB.Where("name = ? AND game_id = ? AND ended_at IS NULL", request.Name, game.ID).First(&gameUser).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.game_user_not_found",
				"success": false,
			})
			return
		}

		database.DB.Model(&gameUser).Where("id = ?", gameUser.ID).Update("Status", models.GAME_USER_PAUSED)

		syncUserChunksForPause(gameUser.ID, game.ID)

		gamemodule.CreateHistory(gameUser.ID, game.ID, models.GAME_HISTORY_PAUSE)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

type PauseTableRequest struct {
	ID uint `json:"id" binding:"required" validate:"required"`
}

func pauseTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request JoinTableRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		user, err := utils.CurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.user_not_found",
				"success": false,
			})
			return
		}

		// check is there game already
		var game models.Game
		err = database.DB.Where("table_id = ? AND ended_at IS NULL", request.ID).First(&game).Error
		if err != nil { // I'll guess it's not found
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.game_not_found",
				"success": false,
			})
			return
		}

		var gameUser models.GameUser
		err = database.DB.Where("user_id = ? AND game_id = ? AND ended_at IS NULL", user.ID, game.ID).First(&gameUser).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.game_user_not_found",
				"success": false,
			})
			return
		}

		database.DB.Model(&gameUser).Where("id = ?", gameUser.ID).Update("Status", models.GAME_USER_PAUSED)

		syncUserChunksForPause(gameUser.ID, game.ID)

		gamemodule.CreateHistory(gameUser.ID, game.ID, models.GAME_HISTORY_PAUSE)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

type JoinAsGuestRequest struct {
	ID   uint   `json:"id" binding:"required" validate:"required"`
	Name string `json:"name" binding:"required" validate:"required"`
}

func joinAsGuest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request JoinAsGuestRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		username := fmt.Sprintf("Guest-%v", strings.ReplaceAll(request.Name, " ", ""))

		_, err := utils.CurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.user_not_found",
				"success": false,
			})
			return
		}

		// check is there game already
		var game models.Game
		err = database.DB.Where("table_id = ? AND ended_at IS NULL", request.ID).Preload("GameUsers").First(&game).Error
		if err != nil { // I'll guess it's not found
			game = models.Game{
				TableID:   request.ID,
				StartedAt: time.Now(),
			}
			err = database.DB.Create(&game).Error
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
					"success": false,
				})
				return
			}
		}

		var gameUsers []models.GameUser
		err = database.DB.Where("game_id = ? AND ended_at IS NULL", game.ID).Find(&gameUsers).Error
		if err == nil {
			// check player count
			if len(gameUsers) > 3 {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "errors.table_full",
					"success": false,
				})
				return
			}
		}

		// check name
		var gameUser models.GameUser
		err = database.DB.Where("name = ? AND ended_at IS NULL", username).First(&gameUser).Error
		if err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.player_name_taken",
				"success": false,
			})
			return
		}

		err = database.DB.Where("game_id = ? AND name = ? AND ended_at IS NULL", game.ID, username).First(&gameUser).Error
		if err != nil {
			gameUser = models.GameUser{
				UserID:    nil,
				GameID:    game.ID,
				Name:      username,
				Status:    models.GAME_USER_PLAYING,
				StartedAt: time.Now(),
			}
			err = database.DB.Create(&gameUser).Error
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "errors.unable_to_join_game",
					"success": false,
				})
				return
			}
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.player_name_taken",
				"success": false,
			})
			return
		}

		syncUserChunksForJoin(gameUser.ID, game.ID)

		gamemodule.CreateHistory(gameUser.ID, game.ID, models.GAME_HISTORY_JOIN)

		ctx.JSON(http.StatusOK, gin.H{
			"game_user": gameUser,
			"success":   true,
		})
	}
}

type JoinTableRequest struct {
	ID uint `json:"id" binding:"required" validate:"required"`
}

func joinTable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request JoinTableRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		user, err := utils.CurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.user_not_found",
				"success": false,
			})
			return
		}

		// check user is not playing other game
		var myGameUser models.GameUser
		err = database.DB.Where("user_id = ? AND ended_at IS NULL", user.ID).First(&myGameUser).Error
		if err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.user_already_playing",
				"success": false,
			})
			return
		}

		// check is there game already
		var game models.Game
		err = database.DB.Where("table_id = ? AND ended_at IS NULL", request.ID).First(&game).Error
		if err != nil { // I'll guess it's not found
			game = models.Game{
				TableID:   request.ID,
				StartedAt: time.Now(),
			}
			err = database.DB.Create(&game).Error
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
					"success": false,
				})
				return
			}
		}

		var gameUsers []models.GameUser
		err = database.DB.Where("game_id = ? AND ended_at IS NULL", game.ID).Find(&gameUsers).Error
		if err == nil {
			// check player count
			if len(gameUsers) > 3 {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "errors.table_full",
					"success": false,
				})
				return
			}
		}

		var gameUser models.GameUser
		err = database.DB.Where("user_id = ? AND game_id = ? AND ended_at IS NULL", user.ID, game.ID).First(&gameUser).Error
		if err != nil {
			gameUser = models.GameUser{
				UserID:    &user.ID,
				GameID:    game.ID,
				Name:      user.Name,
				Status:    models.GAME_USER_PLAYING,
				StartedAt: time.Now(),
			}
			err = database.DB.Create(&gameUser).Error
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "errors.unable_to_join_game",
					"success": false,
				})
				return
			}
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.user_already_joined_game",
				"success": false,
			})
			return
		}

		syncUserChunksForJoin(gameUser.ID, game.ID)

		gamemodule.CreateHistory(gameUser.ID, game.ID, models.GAME_HISTORY_JOIN)

		ctx.JSON(http.StatusOK, gin.H{
			"game_user": gameUser,
			"success":   true,
		})
	}
}

type TransferRequest struct {
	ID         uint `json:"id" binding:"required" validate:"required"`
	NewTableID uint `json:"new_table_id" binding:"required" validate:"required"`
}

func transfer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request TransferRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		user, err := utils.CurrentUser(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.user_not_found",
				"success": false,
			})
			return
		}

		// check is there game on current table
		var oldGame models.Game
		err = database.DB.Where("table_id = ? AND ended_at IS NULL", request.ID).First(&oldGame).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.game_not_found",
				"success": false,
			})
			return
		}

		// check current game user
		var gameUser models.GameUser
		err = database.DB.Where("user_id = ? AND game_id = ? AND ended_at IS NULL", user.ID, oldGame.ID).First(&gameUser).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.game_user_not_found",
				"success": false,
			})
			return
		}

		// check is there game already on new table
		var newGame models.Game
		err = database.DB.Where("table_id = ? AND ended_at IS NULL", request.NewTableID).Preload("GameUsers").First(&newGame).Error
		if err != nil { // I'll guess it's not found
			newGame = models.Game{
				TableID:   request.NewTableID,
				StartedAt: time.Now(),
			}
			err = database.DB.Create(&newGame).Error
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
					"success": false,
				})
				return
			}
		}

		// check user count
		var newGameUsers []models.GameUser
		err = database.DB.Where("game_id = ? and ended_at IS NULL", newGame.ID).Find(&newGameUsers).Error
		if err == nil {
				if len(newGameUsers) > 3 {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"message": "errors.table_full",
						"success": false,
					})
					return
				}
		}

		// syncUserChunksForLeave(gameUser.ID, oldGame.ID)

		// now transfer
		err = database.DB.Model(&gameUser).Where("id = ?", gameUser.ID).Update("game_id", newGame.ID).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "errors.unable_to_transfer_user",
				"success": false,
			})
			return
		}

		// create new chunk for it
		syncUserChunksForTransfer(gameUser.ID, newGame.ID, oldGame.ID)

		gamemodule.CreateHistory(gameUser.ID, oldGame.ID, models.GAME_HISTORY_TRANSFER)
		gamemodule.CreateHistory(gameUser.ID, newGame.ID, models.GAME_HISTORY_JOIN)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

type TransferGuestRequest struct {
	ID         uint   `json:"id" binding:"required" validate:"required"`
	NewTableID uint   `json:"new_table_id" binding:"required" validate:"required"`
	Name       string `json:"name" binding:"required" validate:"required"`
}

func transferGuest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request TransferGuestRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		// check is there game on current table
		var oldGame models.Game
		err := database.DB.Where("table_id = ? AND ended_at IS NULL", request.ID).First(&oldGame).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.game_not_found",
				"success": false,
			})
			return
		}

		// check current game user
		var gameUser models.GameUser
		err = database.DB.Where("name = ? AND game_id = ? AND ended_at IS NULL", request.Name, oldGame.ID).First(&gameUser).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.game_user_not_found",
				"success": false,
			})
			return
		}

		// check is there game already on new table
		var newGame models.Game
		err = database.DB.Where("table_id = ? AND ended_at IS NULL", request.NewTableID).First(&newGame).Error
		if err != nil { // I'll guess it's not found
			newGame = models.Game{
				TableID:   request.NewTableID,
				StartedAt: time.Now(),
			}
			err = database.DB.Create(&newGame).Error
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
					"success": false,
				})
				return
			}
		}

		// check user count
		var newGameUsers []models.GameUser
		err = database.DB.Where("game_id = ? and ended_at IS NULL", newGame.ID).Find(&newGameUsers).Error
		if err == nil {
				if len(newGameUsers) > 3 {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"message": "errors.table_full",
						"success": false,
					})
					return
				}
		}

		// now transfer
		err = database.DB.Model(&gameUser).Where("id = ?", gameUser.ID).Update("game_id", newGame.ID).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "errors.unable_to_transfer_user",
				"success": false,
			})
			return
		}

		// create new chunk for it
		syncUserChunksForTransfer(gameUser.ID, newGame.ID, oldGame.ID)

		gamemodule.CreateHistory(gameUser.ID, oldGame.ID, models.GAME_HISTORY_TRANSFER)
		gamemodule.CreateHistory(gameUser.ID, newGame.ID, models.GAME_HISTORY_JOIN)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

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

	active_player_count := 0
	for _, v := range users {
		if v.Status == models.GAME_USER_PLAYING {
			active_player_count = active_player_count + 1
		}
	}

	for _, v := range users {
		if v.ID == gameUserID {
			database.DB.Create(&models.GameUserChunk{
				GameID:     gameID,
				GameUserID: v.ID,
				Status:     models.GAME_USER_PLAYING,
				Players:    uint8(active_player_count),
				StartedAt:  now,
			})
		} else {
			database.DB.Create(&models.GameUserChunk{
				GameID:     gameID,
				GameUserID: v.ID,
				Status:     v.Status,
				Players:    uint8(active_player_count),
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

	active_player_count := 0
	for _, v := range users {
		if v.Status == models.GAME_USER_PLAYING {
			active_player_count = active_player_count + 1
		}
	}

	for _, v := range users {
		if v.ID == gameUserID {
			database.DB.Create(&models.GameUserChunk{
				GameID:     gameID,
				GameUserID: v.ID,
				Status:     models.GAME_USER_PAUSED,
				Players:    uint8(active_player_count),
				StartedAt:  now,
			})
		} else {
			database.DB.Create(&models.GameUserChunk{
				GameID:     gameID,
				GameUserID: v.ID,
				Status:     v.Status,
				Players:    uint8(active_player_count),
				StartedAt:  now,
			})
		}
	}

	return nil
}

func syncUserChunksForTransfer(IDToExclude uint, gameID uint, oldGameID uint) error {
	// keep the current time
	now := time.Now()

	var savedPrice decimal.Decimal
	var playTimeFloat float64

	// fill price and ended_at values for current chunks
	var chunks []models.GameUserChunk
	database.DB.Where("game_id = ? AND ended_at IS NULL", oldGameID).Find(&chunks)
	for _, v := range chunks {
		if v.GameUserID == IDToExclude {
			totalPrice := gameuser.CalcUserTotalPrice(oldGameID, IDToExclude)
			savedPrice = totalPrice
			playTimeFloat = gameuser.CalcUserTotalTime(oldGameID, IDToExclude)
		}

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

	// create next chunks for old game
	var users []models.GameUser
	err := database.DB.Where("game_id = ? AND id != ? AND ended_at IS NULL", oldGameID, IDToExclude).Find(&users).Error
	if err != nil { // no one left so end old game
		database.DB.Model(&models.Game{}).Where("id = ?", oldGameID).Update("ended_at", now)
	}

	active_player_count := 0
	for _, v := range users {
		if v.Status == models.GAME_USER_PLAYING {
			active_player_count = active_player_count + 1
		}
	}

	if len(users) == 0 {
		database.DB.Model(&models.Game{}).Where("id = ?", oldGameID).Update("ended_at", now)
	} else {
		for _, v := range users {
			database.DB.Create(&models.GameUserChunk{
				GameID:     oldGameID,
				GameUserID: v.ID,
				Status:     v.Status,
				Players:    uint8(active_player_count),
				StartedAt:  now,
			})
		}
	}

	// finish current chunks for new game if exists
	err = database.DB.Where("game_id = ? AND ended_at IS NULL", gameID).Find(&chunks).Error
	if err == nil { // if chunks exists for new game, seal them
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
	}

	// convert float64 to time
	playTime := now.Add(time.Duration(-1 * playTimeFloat) * time.Minute)

	// fill for just the transfered player's price
	database.DB.Where("game_id = ? AND ended_at IS NULL", gameID).Find(&users)

	active_player_count = 0
	for _, v := range users {
		if v.Status == models.GAME_USER_PLAYING {
			active_player_count = active_player_count + 1
		}
	}

	for _, v := range users {
		priceForChunk := decimal.New(0, -2)

		if v.ID == IDToExclude {
			database.DB.Create(&models.GameUserChunk{
				GameID:     gameID,
				GameUserID: v.ID,
				Status:     models.GAME_USER_PLAYING,
				Price:      &savedPrice,
				Players:    uint8(active_player_count),
				//StartedAt:  now, // for client request will try to keep transfered user's play time
				StartedAt: playTime,
				EndedAt:    &now,
			})
		} else {
			database.DB.Create(&models.GameUserChunk{
				GameID:     gameID,
				GameUserID: v.ID,
				Status:     v.Status,
				Price:      &priceForChunk,
				Players:    uint8(active_player_count),
				StartedAt:  now,
				EndedAt:    &now,
			})
		}
	}

	// create new chunks for new game
	database.DB.Where("game_id = ? AND ended_at IS NULL", gameID).Find(&users)

	active_player_count = 0
	for _, v := range users {
		if v.Status == models.GAME_USER_PLAYING {
			active_player_count = active_player_count + 1
		}
	}

	for _, v := range users {
		database.DB.Create(&models.GameUserChunk{
			GameID:     gameID,
			GameUserID: v.ID,
			Status:     v.Status,
			Players:    uint8(active_player_count),
			StartedAt:  now,
		})
	}

	return nil
}

func syncUserChunksForJoin(IDToExclude uint, gameID uint) error {
	// keep the current time
	now := time.Now()

	// fill price and ended_at values for current chunks
	var chunks []models.GameUserChunk
	database.DB.Where("game_id = ? AND ended_at IS NULL", gameID).Find(&chunks)
	for _, v := range chunks {
		log.Println("----- chunk -----")
		log.Println(v)
		log.Println("-----")
		price, err := gameuser.CalcUserChunkPrice(&v, now)
		if err != nil {
			log.Println("----- gameuser.CalcUserChunkPrice() -----")
			log.Println(err.Error())
			log.Println("-----")
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

	active_player_count := 0
	for _, v := range users {
		if v.Status == models.GAME_USER_PLAYING {
			active_player_count = active_player_count + 1
		}
	}

	for _, v := range users {
		database.DB.Create(&models.GameUserChunk{
			GameID:     gameID,
			GameUserID: v.ID,
			Status:     v.Status,
			Players:    uint8(active_player_count),
			StartedAt:  now,
		})
	}

	return nil
}

func syncUserChunksForLeave(gameUserID uint, gameID uint) error {
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
	database.DB.Where("id != ? AND game_id = ? AND ended_at IS NULL", gameUserID, gameID).Find(&users)

	active_player_count := 0
	for _, v := range users {
		if v.Status == models.GAME_USER_PLAYING {
			active_player_count = active_player_count + 1
		}
	}

	for _, v := range users {
		database.DB.Create(&models.GameUserChunk{
			GameID:     gameID,
			GameUserID: v.ID,
			Status:     v.Status,
			Players:    uint8(active_player_count),
			StartedAt:  now,
		})
	}

	return nil
}

func Routes(r *gin.RouterGroup) {
	r.Use(middlewares.AuthMiddleware())
	r.GET("/table/all", getTables())
	r.POST("/table/new", newTable())
	r.POST("/table/update", updateTable())
	r.POST("/table/delete", deleteTable())
	r.GET("/table/:id", getTable())
	r.POST("/table/join", joinTable())
	r.POST("/table/join-as-guest", joinAsGuest())
	r.POST("/table/leave", leaveTable())
	r.POST("/table/leave-as-guest", leaveAsGuest())
	r.POST("/table/pause", pauseTable())
	r.POST("/table/pause-as-guest", pauseAsGuest())
	r.POST("/table/unpause", unpauseTable())
	r.POST("/table/unpause-as-guest", unpauseAsGuest())

	// transfer to other table
	r.POST("/table/transfer", transfer())
	r.POST("/table/transfer-as-guest", transferGuest())
}
