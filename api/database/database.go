package database

import (
	"fmt"
	"os"
	"time"

	"github.com/berkaycubuk/billiard_software_api/models"
	"github.com/shopspring/decimal"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Migrations() {
	DB.AutoMigrate(
		// User
		&models.User{},
		&models.UserSubscription{},
		&models.UserSubscriptionChunk{},
		&models.UserRole{},
		&models.PasswordReset{},

		// Notification
		&models.Notification{},

		// Role
		&models.Role{},
		&models.RoleConfig{},

		// Subscription
		&models.Subscription{},

		// Table
		&models.Table{},

		// Game
		&models.Game{},
		&models.GameUser{},
		&models.GameUserChunk{},
		&models.GameHistory{},
		&models.GameUserOrder{},

		// Pricing
		&models.Pricing{},

		// Upload
		&models.Upload{},

		// Product
		&models.Product{},
		&models.ProductImage{},

		// Order
		&models.Order{},
		&models.OrderDetail{},
		&models.OrderItem{},
		&models.OrderHistory{},
		&models.OrderDiscount{},
	)
}

func ConnectDatabase() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	// dbUrl := os.Getenv("DB_URL")
	// dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// dsn := fmt.Sprintf(
	// 	"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
	// 	dbUser,
	// 	dbPass,
	// 	dbUrl,
	// 	dbPort,
	// 	dbName,
	// )

	// docker special edition with container name "mysql"
	dsn := fmt.Sprintf(
		"%v:%v@tcp(db-mysql)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPass,
		dbName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	DB = db

	//DropTables()

	Migrations()

	//Seed()
}

func DropTables() {
	// user
	DB.Migrator().DropTable(&models.User{})
	DB.Migrator().DropTable(&models.UserSubscription{})
	DB.Migrator().DropTable(&models.UserSubscriptionChunk{})
	DB.Migrator().DropTable(&models.UserRole{})
	DB.Migrator().DropTable(&models.PasswordReset{})

	// Notification
	DB.Migrator().DropTable(&models.Notification{})

	// role
	DB.Migrator().DropTable(&models.Role{})
	DB.Migrator().DropTable(&models.RoleConfig{})

	// subscription
	DB.Migrator().DropTable(&models.Subscription{})

	// table
	DB.Migrator().DropTable(&models.Table{})

	// game
	DB.Migrator().DropTable(&models.Game{})
	DB.Migrator().DropTable(&models.GameUser{})
	DB.Migrator().DropTable(&models.GameUserChunk{})
	DB.Migrator().DropTable(&models.GameHistory{})

	// Pricing
	DB.Migrator().DropTable(&models.Pricing{})

	// Upload
	DB.Migrator().DropTable(&models.Upload{})

	// product
	DB.Migrator().DropTable(&models.Product{})
	DB.Migrator().DropTable(&models.ProductImage{})

	// order
	DB.Migrator().DropTable(&models.Order{})
	DB.Migrator().DropTable(&models.OrderDetail{})
	DB.Migrator().DropTable(&models.OrderItem{})
	DB.Migrator().DropTable(&models.OrderHistory{})
	DB.Migrator().DropTable(&models.OrderDiscount{})
}

func Seed() {
	// Roles
	roles := []models.Role{
		{
			Name:      "Admin",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Kiosk",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Guest Device",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	DB.Create(&roles)

	roleConfigs := []models.RoleConfig{
		// admin
		{
			RoleId:    1,
			Name:      "access_admin",
			Value:     "true",
			ValueType: "bool",
		},
		{
			RoleId:    1,
			Name:      "access_home",
			Value:     "true",
			ValueType: "bool",
		},
		{
			RoleId:    1,
			Name:      "header.show_notifications",
			Value:     "true",
			ValueType: "bool",
		},
		{
			RoleId:    1,
			Name:      "header.show_profile",
			Value:     "true",
			ValueType: "bool",
		},

		// user
		{
			RoleId:    2,
			Name:      "access_admin",
			Value:     "false",
			ValueType: "bool",
		},
		{
			RoleId:    2,
			Name:      "access_home",
			Value:     "true",
			ValueType: "bool",
		},
		{
			RoleId:    2,
			Name:      "header.show_notifications",
			Value:     "true",
			ValueType: "bool",
		},
		{
			RoleId:    2,
			Name:      "header.show_profile",
			Value:     "true",
			ValueType: "bool",
		},

		// kiosk
		{
			RoleId:    3,
			Name:      "access_admin",
			Value:     "false",
			ValueType: "bool",
		},
		{
			RoleId:    3,
			Name:      "access_home",
			Value:     "false",
			ValueType: "bool",
		},
		{
			RoleId:    3,
			Name:      "header.show_notifications",
			Value:     "false",
			ValueType: "bool",
		},
		{
			RoleId:    3,
			Name:      "header.show_profile",
			Value:     "false",
			ValueType: "bool",
		},

		// tablet (guest device) (tablet or computer for guest users)
		{
			RoleId:    4,
			Name:      "access_admin",
			Value:     "false",
			ValueType: "bool",
		},
		{
			RoleId:    4,
			Name:      "access_home",
			Value:     "true",
			ValueType: "bool",
		},
		{
			RoleId:    4,
			Name:      "header.show_notifications",
			Value:     "false",
			ValueType: "bool",
		},
		{
			RoleId:    4,
			Name:      "header.show_profile",
			Value:     "false",
			ValueType: "bool",
		},
	}
	DB.Create(&roleConfigs)
	now := time.Now()

	// Users
	users := []models.User{
		{
			Name:      "Admin",
			Surname:   "User",
			Email:     "admin@admin.admin",
			Phone:     "0000000000",
			Password:  "$2a$08$pUpPKofXRhkCp2bslK91veosI3F2v8tMzOMY2gFnYKb19MHUdxILS", // password
			EmailVerifiedAt: &now,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Test",
			Surname:   "User",
			Email:     "test@test.test",
			Phone:     "0000000001",
			Password:  "$2a$08$pUpPKofXRhkCp2bslK91veosI3F2v8tMzOMY2gFnYKb19MHUdxILS", // password
			EmailVerifiedAt: &now,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Kiosk",
			Surname:   "Device",
			Email:     "kiosk@kiosk.kiosk",
			Phone:     "0000000002",
			Password:  "$2a$08$pUpPKofXRhkCp2bslK91veosI3F2v8tMzOMY2gFnYKb19MHUdxILS", // password
			EmailVerifiedAt: &now,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Guest",
			Surname:   "Device",
			Email:     "guest@guest.guest",
			Phone:     "0000000003",
			Password:  "$2a$08$pUpPKofXRhkCp2bslK91veosI3F2v8tMzOMY2gFnYKb19MHUdxILS", // password
			EmailVerifiedAt: &now,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	DB.Create(&users)

	// User roles
	userRoles := []models.UserRole{
		// admin
		{
			UserID:    1,
			RoleID:    1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		// test user
		{
			UserID:    2,
			RoleID:    2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		// kiosk device
		{
			UserID:    3,
			RoleID:    3,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		// guest device
		{
			UserID:    4,
			RoleID:    4,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	DB.Create(&userRoles)

	tables := []models.Table{
		{
			Name:   "Table One",
			Status: 1,
		},
		{
			Name:   "Table Two",
			Status: 1,
		},
		{
			Name:   "Table Three",
			Status: 1,
		},
		{
			Name:   "Table Four",
			Status: 1,
		},
	}
	DB.Create(&tables)

	products := []models.Product{
		{
			Name:   "Test product 1",
			Price:  decimal.NewFromFloat(12.99),
			Order: 1,
			Status: 1,
		},
	}
	DB.Create(&products)

	subscriptions := []models.Subscription{
		{
			Name:      "Monthly",
			Price:     decimal.NewFromFloat32(100),
			Hours:     24 * 30,
			Hidden:		false,
			CreatedAt: time.Now(),
		},
		{
			Name:      "Yearly",
			Price:     decimal.NewFromFloat32(1000),
			Hours:     24 * 30 * 12,
			Hidden:		false,
			CreatedAt: time.Now(),
		},
	}
	DB.Create(&subscriptions)

	pricings := []models.Pricing{
		// admin pricing
		{
			RoleID:         &roles[0].ID,
			SubscriptionID: &subscriptions[0].ID,
			PlayerCount:    1,
			PerMinute:      decimal.New(0, -2),
		},
		{
			RoleID:         &roles[0].ID,
			SubscriptionID: &subscriptions[0].ID,
			PlayerCount:    2,
			PerMinute:      decimal.New(0, -2),
		},
		{
			RoleID:         &roles[0].ID,
			SubscriptionID: &subscriptions[0].ID,
			PlayerCount:    3,
			PerMinute:      decimal.New(0, -2),
		},
		{
			RoleID:         &roles[0].ID,
			SubscriptionID: &subscriptions[0].ID,
			PlayerCount:    4,
			PerMinute:      decimal.New(0, -2),
		},

		// admin pricing
		{
			RoleID:         &roles[0].ID,
			SubscriptionID: &subscriptions[1].ID,
			PlayerCount:    1,
			PerMinute:      decimal.New(0, -2),
		},
		{
			RoleID:         &roles[0].ID,
			SubscriptionID: &subscriptions[1].ID,
			PlayerCount:    2,
			PerMinute:      decimal.New(0, -2),
		},
		{
			RoleID:         &roles[0].ID,
			SubscriptionID: &subscriptions[1].ID,
			PlayerCount:    3,
			PerMinute:      decimal.New(0, -2),
		},
		{
			RoleID:         &roles[0].ID,
			SubscriptionID: &subscriptions[1].ID,
			PlayerCount:    4,
			PerMinute:      decimal.New(0, -2),
		},

		// admin pricing
		{
			RoleID:         &roles[0].ID,
			SubscriptionID: nil,
			PlayerCount:    1,
			PerMinute:      decimal.New(100, -2),
		},
		{
			RoleID:         &roles[0].ID,
			SubscriptionID: nil,
			PlayerCount:    2,
			PerMinute:      decimal.New(100, -2),
		},
		{
			RoleID:         &roles[0].ID,
			SubscriptionID: nil,
			PlayerCount:    3,
			PerMinute:      decimal.New(66, -2),
		},
		{
			RoleID:         &roles[0].ID,
			SubscriptionID: nil,
			PlayerCount:    4,
			PerMinute:      decimal.New(50, -2),
		},

		// guest pricing
		{
			RoleID:         nil,
			SubscriptionID: nil,
			PlayerCount:    1,
			PerMinute:      decimal.New(150, -2),
		},
		{
			RoleID:         nil,
			SubscriptionID: nil,
			PlayerCount:    2,
			PerMinute:      decimal.New(150, -2),
		},
		{
			RoleID:         nil,
			SubscriptionID: nil,
			PlayerCount:    3,
			PerMinute:      decimal.New(100, -2),
		},
		{
			RoleID:         nil,
			SubscriptionID: nil,
			PlayerCount:    4,
			PerMinute:      decimal.New(75, -2),
		},

		// user pricing
		{
			RoleID:         &roles[1].ID,
			SubscriptionID: nil,
			PlayerCount:    1,
			PerMinute:      decimal.New(100, -2),
		},
		{
			RoleID:         &roles[1].ID,
			SubscriptionID: nil,
			PlayerCount:    2,
			PerMinute:      decimal.New(100, -2),
		},
		{
			RoleID:         &roles[1].ID,
			SubscriptionID: nil,
			PlayerCount:    3,
			PerMinute:      decimal.New(66, -2),
		},
		{
			RoleID:         &roles[1].ID,
			SubscriptionID: nil,
			PlayerCount:    4,
			PerMinute:      decimal.New(50, -2),
		},

		// user pricing with monthly sub
		{
			RoleID:         &roles[1].ID,
			SubscriptionID: &subscriptions[0].ID,
			PlayerCount:    1,
			PerMinute:      decimal.New(0, -2),
		},
		{
			RoleID:         &roles[1].ID,
			SubscriptionID: &subscriptions[0].ID,
			PlayerCount:    2,
			PerMinute:      decimal.New(0, -2),
		},
		{
			RoleID:         &roles[1].ID,
			SubscriptionID: &subscriptions[0].ID,
			PlayerCount:    3,
			PerMinute:      decimal.New(0, -2),
		},
		{
			RoleID:         &roles[1].ID,
			SubscriptionID: &subscriptions[0].ID,
			PlayerCount:    4,
			PerMinute:      decimal.New(0, -2),
		},

		// user pricing with yearly sub
		{
			RoleID:         &roles[1].ID,
			SubscriptionID: &subscriptions[1].ID,
			PlayerCount:    1,
			PerMinute:      decimal.New(0, -2),
		},
		{
			RoleID:         &roles[1].ID,
			SubscriptionID: &subscriptions[1].ID,
			PlayerCount:    2,
			PerMinute:      decimal.New(0, -2),
		},
		{
			RoleID:         &roles[1].ID,
			SubscriptionID: &subscriptions[1].ID,
			PlayerCount:    3,
			PerMinute:      decimal.New(0, -2),
		},
		{
			RoleID:         &roles[1].ID,
			SubscriptionID: &subscriptions[1].ID,
			PlayerCount:    4,
			PerMinute:      decimal.New(0, -2),
		},
	}
	DB.Create(&pricings)
}
