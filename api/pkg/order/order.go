package order

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/models"
	"github.com/berkaycubuk/billiard_software_api/pkg/gamemodule"
	"github.com/berkaycubuk/billiard_software_api/pkg/subscription"
	"github.com/berkaycubuk/billiard_software_api/pkg/vipps"
	"github.com/berkaycubuk/billiard_software_api/utils"
	"github.com/shopspring/decimal"
)

func CheckPaidWaitingOrders() *[]models.Order {
	var waitingOrders []models.Order

	err := database.DB.Joins("INNER JOIN order_histories oh ON oh.order_id = orders.id").Where("orders.status = ? and oh.action = ? and oh.method = ?", models.ORDER_STATUS_WAITING, models.ORDER_ACTION_INIT, models.ORDER_METHOD_VIPPS).Find(&waitingOrders).Error
	if err != nil {
		return nil
	}

	var temp []models.Order

	// lets ask vipps
	for _, v := range waitingOrders {
		resp, err := vipps.Status(v.ID, v.Reference)
		if err == nil {
			if (resp.PaymentDetails.State == "AUTHORIZED") {
				temp = append(temp, v)
				UpdateOrderStatus(&v, models.ORDER_STATUS_PAID, models.ORDER_ACTION_PAY, models.ORDER_METHOD_VIPPS)
			}
		}
	}

	return &temp
}

func CheckWaitingOrders() {
	CheckPaidWaitingOrders()

	var waitingOrders []models.Order

	err := database.DB.Where("status = ? AND created_at <= NOW() - INTERVAL 24 HOUR", models.ORDER_STATUS_WAITING).Find(&waitingOrders).Error
	if err != nil {
		return
	}

	for _, v := range waitingOrders {
		database.DB.Delete(&v)
	}
}

func Create(userID *uint, price decimal.Decimal, status uint8) (*models.Order, error) {
	ref_string := ""
	if userID != nil {
		ref_string = fmt.Sprintf("billiard-%v-%d", utils.RandomStringLowercase(10), userID)
	} else {
		ref_string = fmt.Sprintf("billiard-%v-%d", utils.RandomStringLowercase(10), rand.Int())
	}

	order := models.Order{
		Reference: ref_string,
		Status:    status,
		UserID:    userID,
		Price:     price,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := database.DB.Create(&order).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func CreateOrderItem(orderID uint, productType uint8, productID *uint, productName string, productAmount int, productPrice decimal.Decimal) (*models.OrderItem, error) {
	orderItem := models.OrderItem{
		OrderID:       orderID,
		ProductType:   productType,
		ProductID:     productID,
		ProductName:   productName,
		ProductAmount: productAmount,
		ProductPrice:  productPrice,
	}

	err := database.DB.Create(&orderItem).Error
	if err != nil {
		return nil, err
	}

	return &orderItem, nil
}

func CreateOrderDetails(orderID uint, userName string, userSurname string, userPhone string, userEmail string) (*models.OrderDetail, error) {
	orderDetails := models.OrderDetail{
		OrderID:     orderID,
		UserName:    userName,
		UserSurname: userSurname,
		UserPhone:   userPhone,
		UserEmail:   userEmail,
	}

	err := database.DB.Create(&orderDetails).Error
	if err != nil {
		return nil, err
	}

	return &orderDetails, nil
}

func CreateOrderHistory(orderID uint, action uint8, method uint8) (*models.OrderHistory, error) {
	orderHistory := models.OrderHistory{
		OrderID:   orderID,
		Action:    action,
		Method:    method,
		CreatedAt: time.Now(),
	}

	err := database.DB.Create(&orderHistory).Error
	if err != nil {
		return nil, err
	}

	return &orderHistory, nil
}

func UpdateOrderStatus(order *models.Order, status uint8, action uint8, method uint8) {
	if order.Status == status {
		return
	}

	database.DB.Model(&order).Update("status", status)

	var items []models.OrderItem
	err := database.DB.Where("order_id = ?", order.ID).Find(&items).Error
	if err == nil {
		if status == models.ORDER_STATUS_PAID || status == models.ORDER_STATUS_APPROVED {
				for _, v := range items {
					HandleOrderItemAfterPayment(order, &v)
				}
		}
	}

	CreateOrderHistory(order.ID, action, method)
}

func HandleOrderItemAfterPayment(order *models.Order, item *models.OrderItem) {
	if item.ProductType == models.PRODUCT_TYPE_SUBSCRIPTION {
		subscription.ActivateSubscription(*order.UserID, *item.ProductID)
	} else if item.ProductType == models.PRODUCT_TYPE_GAME {
		var gameUser models.GameUser
		// if user_id not exists find it with order user_name
		if order.UserID != nil {
			err := database.DB.Where("user_id = ? AND ended_at IS NULL", order.UserID).First(&gameUser).Error
			if err != nil {
				return
			}
		} else {
			// get details
			var orderDetail models.OrderDetail
			err := database.DB.Where("order_id = ?", order.ID).First(&orderDetail).Error
			if err != nil {
				return
			}

			// then find the game user with order detail name
			err = database.DB.Where("name = ? AND ended_at IS NULL", orderDetail.UserName).First(&gameUser).Error
			if err != nil {
				return
			}
		}

		gamemodule.SyncUserChunksForLeave(gameUser.ID, *item.ProductID)

		err := database.DB.Model(&models.GameUser{}).Where("id = ?", gameUser.ID).
			Update("ended_at", time.Now()).Error
		if err != nil {
			return
		}

		gamemodule.CreateHistory(gameUser.ID, *item.ProductID, models.GAME_HISTORY_LEAVE)

		// check are there anyone playing
		var gameUsers []models.GameUser
		database.DB.Where("game_id = ? AND ended_at IS NULL", item.ProductID).Find(&gameUsers)
		if len(gameUsers) == 0 { // no one playing
			database.DB.Model(&models.Game{}).Where("id = ?", item.ProductID).Update("ended_at", time.Now())
		}

		// pay laters
		if order.UserID != nil {
			var payLaterOrders []models.Order
			err = database.DB.Where("user_id = ? and status = ?", order.UserID, models.ORDER_STATUS_PAY_LATER).Find(&payLaterOrders).Error
			if err == nil {
				for _, v := range payLaterOrders {
					database.DB.Model(&v).Update("status", models.ORDER_STATUS_TRANSFERRED)
				}
			}
		}

		// orders added to game user
		var gameUserOrders []models.GameUserOrder
		err = database.DB.Where("game_user_id = ?", gameUser.ID).Find(&gameUserOrders).Error
		if err == nil {
				for _, v := range gameUserOrders {
					var orderTemp models.Order
					err = database.DB.Where("id = ? and status = ?", v.OrderID, models.ORDER_STATUS_PAY_LATER).First(&orderTemp).Error
					if err == nil {
							database.DB.Model(&orderTemp).Update("status", models.ORDER_STATUS_TRANSFERRED)
					}
				}
		}
	}
}

func TransferPayLater(payLaterOrderID uint, newOrderID uint) {
	// now := time.Now()

	// check pay later order
	var payLaterOrder models.Order
	err := database.DB.Where("id = ?", payLaterOrderID).First(&payLaterOrder).Error
	if err != nil { // not found I guess
		return
	}

	// is paylater
	if payLaterOrder.Status != models.ORDER_STATUS_PAY_LATER {
		return
	}

	var newOrder models.Order
	err = database.DB.Where("id = ?", newOrderID).First(&newOrder).Error
	if err != nil {
		return
	}

	// order items
	var oldItems []models.OrderItem
	database.DB.Where("order_id = ?", payLaterOrderID).Find(&oldItems)
	for _, v := range oldItems {
		database.DB.Create(&models.OrderItem{
			OrderID:       newOrderID,
			ProductID:     v.ProductID,
			ProductName:   v.ProductName,
			ProductAmount: v.ProductAmount,
			ProductPrice:  v.ProductPrice,
			ProductType:   v.ProductType,
		})
	}

	// update price of new order
	newPrice := newOrder.Price.Add(payLaterOrder.Price)
	err = database.DB.Model(&newOrder).Update("price", newPrice).Error
	if err != nil {
		fmt.Println(err.Error())
	}

	/*
		// delete after payment is complete
		database.DB.Model(models.Order{}).
			Where("id = ?", payLaterOrderID).
			UpdateColumns(models.Order{
				Status:    models.ORDER_STATUS_TRANSFERRED,
				DeletedAt: &now,
			})

		CreateOrderHistory(payLaterOrderID, models.ORDER_ACTION_TRANSFER, models.ORDER_METHOD_SYSTEM)
	*/

	// implement discounts as well
}

func GetAll() []models.Order {
	var orders []models.Order
	err := database.DB.Where("deleted_at IS NULL").Find(&orders).Error
	if err != nil {
		return nil
	}

	return orders
}

func Get(ID uint) *models.Order {
	var order models.Order
	err := database.DB.Where("id = ?", ID).Preload("Histories").Preload("Discounts").Preload("Items").Preload("Detail").First(&order).Error
	if err != nil {
		return nil
	}

	return &order
}

func Delete(ID uint) {
	order := models.Order{
		ID: ID,
	}

	database.DB.Delete(&order)
}

func Approve(ID uint) error {
	var order models.Order
	err := database.DB.Where("id = ?", ID).First(&order).Error
	if err != nil {
		return nil
	}

	err = database.DB.Model(&order).Update("status", models.ORDER_STATUS_APPROVED).Error
	if err != nil {
		return err
	}

	var items []models.OrderItem
	err = database.DB.Where("order_id = ?", order.ID).Find(&items).Error
	if err == nil {
		for _, v := range items {
			HandleOrderItemAfterPayment(&order, &v)
		}
	}

	CreateOrderHistory(order.ID, models.ORDER_ACTION_APPROVE, models.ORDER_METHOD_SYSTEM)

	return nil
}
