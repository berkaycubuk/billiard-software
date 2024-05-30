package reporting

import (
	"fmt"
	"net/http"
	"time"

	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/middlewares"
	"github.com/berkaycubuk/billiard_software_api/models"
	"github.com/berkaycubuk/billiard_software_api/utils"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

var notCountedOrderTypes = []int{models.ORDER_STATUS_CANCELED, models.ORDER_STATUS_DELETED, models.ORDER_STATUS_TRANSFERRED, int(models.ORDER_STATUS_WAITING)}

type GetGameRequest struct {
	ID	uint	`uri:"id" binding:"required"`
}

type AllTableSalesRequest struct {
	From time.Time `json:"from" binding:"required" validate:"required"`
	To   time.Time `json:"tow" binding:"required" validate:"required"`
}

type TableSalesRequest struct {
	ID   uint      `json:"id" binding:"required" validate:"required"`
}

func tableSales() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request TableSalesRequest
		if !utils.ValidateRequest(ctx, &request) {
			return
		}

		startParam := ctx.Query("start")
		endParam := ctx.Query("end")

		sales, games := TableGameSales(request.ID, startParam, endParam)

		ctx.JSON(http.StatusOK, gin.H{
			"sales":   sales,
			"games":   games,
			"success": true,
		})
	}
}

func allTableSales() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// var request AllTableSalesRequest
		// if !utils.ValidateRequest(ctx, &request) {
		// 	return
		// }

		startParam := ctx.Query("start")
		endParam := ctx.Query("end")

		sales := AllTableGameSales(startParam, endParam)

		ctx.JSON(http.StatusOK, gin.H{
			"sales":   sales,
			"success": true,
		})
	}
}

func allGames() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		countChain := database.DB
		chain := database.DB

		start := ctx.Query("start")
		if start != "" {
			chain = chain.Where("started_at >= ?", start)
			countChain = countChain.Where("started_at >= ?", start)
		}

		end := ctx.Query("end")
		if end != "" {
			chain = chain.Where("started_at <= ?", end+"T23:59:59")
			countChain = countChain.Where("started_at <= ?", end+"T23:59:59")
		}

		var gamesTemp []models.Game
		countChain.Select("id").Find(&gamesTemp)
		count := len(gamesTemp)

		pageSize := count / 10
		if pageSize*10 < count {
			pageSize = pageSize + 1
		}

		var games []models.Game
		chain.Scopes(utils.Paginate(ctx)).Order("id desc").Find(&games)

		ctx.Set("response", gin.H{
			"games":   games,
			"pagination": gin.H{
				"size": pageSize,
			},
		})
	}
}

func getGame() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request GetGameRequest
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.id_required",
				"success": false,
			})
			return
		}
		
		var game models.Game
		err := database.DB.Where("id = ?", request.ID).Preload("GameUsers").First(&game).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "errors.game_not_found",
				"success": false,
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"game":   game,
			"success": true,
		})
	}
}

func subscriptionSales() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// check query parameters
		from := ctx.Query("start")
		to := ctx.Query("end")

		var orderItems []models.OrderItem
		if from != "" && to != "" {
			err := database.DB.
				Where("created_at BETWEEN ? AND ?", from, to).
				Joins("JOIN order_items ON orders.id = order_items.order_id").
				Where("product_type = ?", models.PRODUCT_TYPE_SUBSCRIPTION).
				Preload("Items").
				Find(&orderItems).Error
			if err != nil {
				fmt.Println(err.Error())
				ctx.JSON(http.StatusOK, gin.H{
					"sales": decimal.New(0, -2),
					"success": true,
				})
			}
		} else if from != "" && to == "" {
			err := database.DB.
				Where("created_at >= ?", from).
				Joins("JOIN order_items ON orders.id = order_items.order_id").
				Where("product_type = ?", models.PRODUCT_TYPE_SUBSCRIPTION).
				Preload("Items").
				Find(&orderItems).Error
			if err != nil {
				fmt.Println(err.Error())
				ctx.JSON(http.StatusOK, gin.H{
					"sales": decimal.New(0, -2),
					"success": true,
				})
			}
		} else if from == "" && to != "" {
			err := database.DB.
				Where("created_at <= ?", to).
				Joins("JOIN order_items ON orders.id = order_items.order_id").
				Where("product_type = ?", models.PRODUCT_TYPE_SUBSCRIPTION).
				Preload("Items").
				Find(&orderItems).Error
			if err != nil {
				fmt.Println(err.Error())
				ctx.JSON(http.StatusOK, gin.H{
					"sales": decimal.New(0, -2),
					"success": true,
				})
			}
		} else {
			err := database.DB.
				Where("product_type = ?", models.PRODUCT_TYPE_SUBSCRIPTION).
				Find(&orderItems).Error
			if err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"sales": decimal.New(0, -2),
					"success": true,
				})
			}
		}


		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func sub() gin.HandlerFunc {
	return func (ctx *gin.Context) {
		var items []models.OrderItem

		startParam := ctx.Query("start")
		endParam := ctx.Query("end")

		query := database.DB.Joins("JOIN orders ON order_items.order_id = orders.id").
			Where("order_items.product_type = ?", models.PRODUCT_TYPE_SUBSCRIPTION)
		if startParam != "" {
			query = query.Where("orders.created_at >= ?", startParam)
		}

		if endParam != "" {
			query = query.Where("orders.created_at <= ?", endParam + "T23:59:59")
		}

		query = query.Where("orders.status NOT IN ?", notCountedOrderTypes)
		query.Find(&items)

		var products []KioskProduct
		for _, v := range items {
			// check product id of v.ProductID exists?
			found := false
			for index, j := range products {
				if j.ProductID == *v.ProductID {
					found = true

					fmt.Println(j)

					// update related product values
					j.Count = j.Count + v.ProductAmount
					j.Total = j.Total.Add(v.ProductPrice.Mul(decimal.NewFromInt(int64(v.ProductAmount))))

					products[index] = j
				}
			}

			// if not create a new one else append available data to it
			if !found {
				var product models.Subscription
				err := database.DB.Where("id = ?", v.ProductID).First(&product).Error
				if err != nil {
					continue
				}

				products = append(products, KioskProduct{
					ProductID: *v.ProductID,
					Name: product.Name,
					Count: v.ProductAmount,
					Total: v.ProductPrice.Mul(decimal.NewFromInt(int64(v.ProductAmount))),
				})
			}
		}

		var userSubscriptions []models.UserSubscription
		queryA := database.DB
		if startParam != "" {
			queryA = queryA.Where("created_at >= ?", startParam)
		}

		if endParam != "" {
			queryA = queryA.Where("created_at <= ?", endParam + "T23:59:59")
		}

		queryA.Preload("User").Order("id DESC").Find(&userSubscriptions)

		ctx.JSON(http.StatusOK, gin.H{
			"products": products,
			"users": userSubscriptions,
			"success": true,
		})
	}
}

type KioskProduct struct {
	ProductID	uint	`json:"product_id"`
	Name	string	`json:"name"`
	//ProductName string `json:"product_name"`
	Count		int		`json:"count"`
	Total	decimal.Decimal	`json:"total"`
}

func kiosk() gin.HandlerFunc {
	return func (ctx *gin.Context) {
		var items []models.OrderItem
		total := decimal.Zero

		startParam := ctx.Query("start")
		endParam := ctx.Query("end")

		query := database.DB.Joins("JOIN orders ON order_items.order_id = orders.id").
			Where("order_items.product_type = ?", models.PRODUCT_TYPE_SHOP)
		if startParam != "" {
			query = query.Where("orders.created_at >= ?", startParam)
		}

		if endParam != "" {
			query = query.Where("orders.created_at <= ?", endParam + "T23:59:59")
		}

		query = query.Where("orders.status NOT IN ?", notCountedOrderTypes)
		query.Find(&items)

		var products []KioskProduct
		for _, v := range items {
			// check product id of v.ProductID exists?
			found := false
			for index, j := range products {
				if j.ProductID == *v.ProductID {
					found = true

					// update related product values
					j.Count = j.Count + v.ProductAmount
					j.Total = j.Total.Add(v.ProductPrice.Mul(decimal.NewFromInt(int64(v.ProductAmount))))

					total = total.Add(v.ProductPrice.Mul(decimal.NewFromInt(int64(v.ProductAmount))))

					products[index] = j
				}
			}

			// if not create a new one else append available data to it
			if !found {
				var product models.Product
				database.DB.Where("id = ?", v.ProductID).First(&product)

				products = append(products, KioskProduct{
					ProductID: *v.ProductID,
					Name: product.Name,
					Count: v.ProductAmount,
					Total: v.ProductPrice.Mul(decimal.NewFromInt(int64(v.ProductAmount))),
				})

				total = total.Add(v.ProductPrice.Mul(decimal.NewFromInt(int64(v.ProductAmount))))
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"products": products,
			"total": total,
			"success": true,
		})
	}
}

func total() gin.HandlerFunc {
	return func (ctx *gin.Context) {
		startParam := ctx.Query("start")
		endParam := ctx.Query("end")

		query := database.DB.Table("orders")
		if startParam != "" {
			query = query.Where("orders.created_at >= ?", startParam)
		}

		if endParam != "" {
			query = query.Where("orders.created_at <= ?", endParam + "T23:59:59")
		}

		var total decimal.Decimal
		query = query.Where("orders.status NOT IN ?", notCountedOrderTypes)
		query.Select("SUM(price) as total").Scan(&total)

		ctx.JSON(http.StatusOK, gin.H{
			"total": total,
			"success": true,
		})
	}
}

type UserTotal struct {
	User	models.User	`json:"user"`
	Total	decimal.Decimal	`json:"total"`
	Count	int		`json:"count"`
}

func users() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startParam := ctx.Query("start")
		endParam := ctx.Query("end")
		userParam := ctx.Query("user")

		var users []models.User
		usersQuery := database.DB

		if userParam != "" {
			usersQuery = usersQuery.Where("name LIKE ?", userParam + "%")
		}

		usersQuery.Scopes(utils.Paginate(ctx)).Order("id desc").Find(&users)

		var result []UserTotal
		for _, user := range users {
			total := decimal.Zero
			count := 0

			var orders []models.Order
			query := database.DB

			if startParam != "" {
				query = query.Where("created_at >= ?", startParam)
			}

			if endParam != "" {
				query = query.Where("created_at <= ?", endParam + "T23:59:59")
			}

			query = query.Where("status NOT IN ?", notCountedOrderTypes)
			query.Where("user_id = ?", user.ID).Find(&orders)
			for _, order := range orders {
				total = total.Add(order.Price)
				count = count + 1
			}
			result = append(result, UserTotal{
				User: user,
				Total: total,
				Count: count,
			})
		}

		var usersTemp []models.User
		countQuery := database.DB

		if userParam != "" {
			countQuery = countQuery.Where("name LIKE ?", userParam + "%")
		}

		countQuery.Select("id").Find(&usersTemp)
		count := len(usersTemp)

		pageSize := count / 10
		if pageSize*10 < count {
			pageSize = pageSize + 1
		}

		ctx.Set("response", gin.H{
			"users":   result,
			"pagination": gin.H{
				"size": pageSize,
			},
		})
	}
}

type UserRequest struct {
	ID	uint	`uri="id" binding:"required"`
}

func user() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request GetGameRequest
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

		var items []models.OrderItem
		startParam := ctx.Query("start")
		endParam := ctx.Query("end")

		query := database.DB.Joins("JOIN orders ON order_items.order_id = orders.id").
			Where("order_items.product_type = ?", models.PRODUCT_TYPE_SHOP)
		if startParam != "" {
			query = query.Where("orders.created_at >= ?", startParam)
		}

		if endParam != "" {
			query = query.Where("orders.created_at <= ?", endParam + "T23:59:59")
		}

		total := decimal.Zero

		query = query.Where("orders.status NOT IN ?", notCountedOrderTypes)
		query.Where("orders.user_id = ?", user.ID).Find(&items)

		var products []KioskProduct
		for _, v := range items {
			// check product id of v.ProductID exists?
			found := false
			for index, j := range products {
				if j.ProductID == *v.ProductID {
					found = true

					// update related product values
					j.Count = j.Count + v.ProductAmount
					j.Total = j.Total.Add(v.ProductPrice.Mul(decimal.NewFromInt(int64(v.ProductAmount))))

					products[index] = j
				}
			}

			// if not create a new one else append available data to it
			if !found {
				var product models.Product
				database.DB.Where("id = ?", v.ProductID).First(&product)

				products = append(products, KioskProduct{
					ProductID: *v.ProductID,
					Name: product.Name,
					Count: v.ProductAmount,
					Total: v.ProductPrice.Mul(decimal.NewFromInt(int64(v.ProductAmount))),
				})
			}

			total = total.Add(v.ProductPrice.Mul(decimal.NewFromInt(int64(v.ProductAmount))))
		}

		// add game payments
		query = database.DB.Joins("JOIN orders ON order_items.order_id = orders.id").
			Where("order_items.product_type = ?", models.PRODUCT_TYPE_GAME)
		if startParam != "" {
			query = query.Where("orders.created_at >= ?", startParam)
		}
		if endParam != "" {
			query = query.Where("orders.created_at <= ?", endParam + "T23:59:59")
		}
		query = query.Where("orders.status NOT IN ?", notCountedOrderTypes)
		query.Where("orders.user_id = ?", user.ID).Find(&items)
		for _, v := range items {
			found := false
			for index, j := range products {
				if j.Name == "Game" {
					found = true

					// update related product values
					j.Count = j.Count + v.ProductAmount
					j.Total = j.Total.Add(v.ProductPrice.Mul(decimal.NewFromInt(int64(v.ProductAmount))))

					products[index] = j
				}
			}

			if !found {
				products = append(products, KioskProduct{
					ProductID: *v.ProductID,
					Name: "Game",
					Count: v.ProductAmount,
					Total: v.ProductPrice.Mul(decimal.NewFromInt(int64(v.ProductAmount))),
				})
			}

			total = total.Add(v.ProductPrice.Mul(decimal.NewFromInt(int64(v.ProductAmount))))
		}

		// add subscription payments
		query = database.DB.Joins("JOIN orders ON order_items.order_id = orders.id").
			Where("order_items.product_type = ?", models.PRODUCT_TYPE_SUBSCRIPTION)
		if startParam != "" {
			query = query.Where("orders.created_at >= ?", startParam)
		}
		if endParam != "" {
			query = query.Where("orders.created_at <= ?", endParam + "T23:59:59")
		}
		query = query.Where("orders.status NOT IN ?", notCountedOrderTypes)
		query.Where("orders.user_id = ?", user.ID).Find(&items)
		for _, v := range items {
			var product models.Subscription
			database.DB.Where("id = ?", v.ProductID).First(&product)
			found := false
			for index, j := range products {
				if j.Name == "Subscription" { //product.Name {
					found = true

					// update related product values
					j.Count = j.Count + v.ProductAmount
					j.Total = j.Total.Add(v.ProductPrice.Mul(decimal.NewFromInt(int64(v.ProductAmount))))

					products[index] = j
				}
			}

			if !found {
				products = append(products, KioskProduct{
					ProductID: *v.ProductID,
					Name: "Subscription", //product.Name,
					Count: v.ProductAmount,
					Total: v.ProductPrice.Mul(decimal.NewFromInt(int64(v.ProductAmount))),
				})
			}

			total = total.Add(v.ProductPrice.Mul(decimal.NewFromInt(int64(v.ProductAmount))))
		}

		ctx.JSON(http.StatusOK, gin.H{
			"user": user,
			"products": products,
			"total": total,
			"success": true,
		})
	}
}

func Routes(r *gin.RouterGroup) {
	r.Use(middlewares.AuthMiddleware())
	r.GET("/reports/games", allGames())
	r.GET("/reports/games/:id", getGame())
	r.GET("/reports/table-sales/all", allTableSales())
	r.GET("/reports/table-sales/:id", tableSales())
	r.GET("/reports/subscriptions", sub())
	r.GET("/reports/kiosk", kiosk())
	r.GET("/reports/users", users())
	r.GET("/reports/users/:id", user())
	r.GET("/reports/total", total())
}
