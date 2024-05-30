package reporting

import (
	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/models"
	"github.com/shopspring/decimal"
)

type TableGameSale struct {
	TableID	uint	`json:"table_id"`
	Name string	`json:"name"`
	Sales	decimal.Decimal `json:"sales"`
	Games	int		`json:"games"`
}

func TableGameSales(tableID uint, from string, to string) (decimal.Decimal, int) {
	// get games
	var games []models.Game
	query := database.DB.Where("table_id = ?", tableID)
	if to != "" {
		query.Where("started_at <= ?", to)
	}

	if from != "" {
		query.Where("started_at >= ?", from)
	}

	// new logic
	err := query.Find(&games).Error
	if err != nil {
		return decimal.New(0, -2), 0
	}

	total := decimal.Zero
	for _, v := range games {
		var orderItems []models.OrderItem
		err = database.DB.Joins("JOIN orders ON order_items.order_id = orders.id").Where("orders.status NOT IN ?", []int{models.ORDER_STATUS_CANCELED, models.ORDER_STATUS_DELETED, models.ORDER_STATUS_TRANSFERRED}).Where("order_items.product_type = ? and order_items.product_id = ?", models.PRODUCT_TYPE_GAME, v.ID).Find(&orderItems).Error
		if err == nil {
			for _, j := range orderItems {
				total = total.Add(j.ProductPrice)
			}
		}
	}

	return total, len(games)

	/*
	err := query.Preload("GameUsers").Find(&games).Error
	if err != nil {
		return decimal.New(0, -2), 0
	}

	total := decimal.New(0, -2)
	for _, v := range games {
		for _, j := range v.GameUsers {
			userPrice := gameuser.CalcUserTotalPrice(v.ID, j.ID)
			total = total.Add(userPrice)
		}
	}

	return total, len(games)
	*/
}

func getTableName(id uint) string {
	var table models.Table
	err := database.DB.Where("id = ?", id).First(&table).Error
	if err != nil {
		return ""
	}

	return table.Name
}

func AllTableGameSales(from string, to string) []TableGameSale {
	// get tables
	var tables []models.Table
	database.DB.Find(&tables)

	var tableGameSales []TableGameSale
	for _, v := range tables {
		price, count := TableGameSales(v.ID, from, to)
		name := getTableName(v.ID)
		tableGameSales = append(tableGameSales, TableGameSale{
			TableID: v.ID,
			Name: name,
			Sales: price,
			Games: count,
		})
	}

	return tableGameSales
}
