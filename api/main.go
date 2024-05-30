package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/pkg/dashboard"
	"github.com/berkaycubuk/billiard_software_api/pkg/notification"
	"github.com/berkaycubuk/billiard_software_api/pkg/order"
	"github.com/berkaycubuk/billiard_software_api/pkg/payment"
	"github.com/berkaycubuk/billiard_software_api/pkg/pricing"
	"github.com/berkaycubuk/billiard_software_api/pkg/reporting"
	"github.com/berkaycubuk/billiard_software_api/pkg/schedule"
	"github.com/berkaycubuk/billiard_software_api/pkg/shop"
	"github.com/berkaycubuk/billiard_software_api/pkg/table"
	"github.com/berkaycubuk/billiard_software_api/pkg/upload"
	"github.com/berkaycubuk/billiard_software_api/routes"
	"github.com/berkaycubuk/gofusion"
	//"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func main() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}


	/*
		if err := sentry.Init(sentry.ClientOptions{
				Dsn: "",
				EnableTracing: true,
				// Set TracesSampleRate to 1.0 to capture 100%
				// of transactions for performance monitoring.
				// We recommend adjusting this value in production,
				TracesSampleRate: 1.0,
		}); err != nil {
				fmt.Printf("Sentry initialization failed: %v", err)
		}
	*/


	database.ConnectDatabase()

	// router setup
	reactor := gofusion.Init()

	reactor.Router.Use(sentrygin.New(sentrygin.Options{
			Repanic: true,
	}))

	reactor.Router.Use(gofusion.CORSMiddleware("*"))
	reactor.Router.Use(gofusion.HandleAPIResponse())

	routes.AuthRoutes(reactor.Router)
	routes.UserRoutes(reactor.Router)
	routes.ProductRoutes(reactor.Router)
	routes.RoleRoutes(reactor.Router)
	routes.SubscriptionRoutes(reactor.Router)

	reactor.Router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	})

	reactor.Router.Static("/static", "./static")

	v1 := reactor.Router.Group("/v1")
	{
		table.Routes(v1)
		payment.Routes(v1)
		shop.Routes(v1)
		order.Routes(v1)
		upload.Routes(v1)
		pricing.Routes(v1)
		reporting.Routes(v1)
		notification.Routes(v1)
		dashboard.Routes(v1)
	}

	// schedule
	schedule.RunJobs()

	// start the http server
	gofusion.Run(reactor, 4000)
}
