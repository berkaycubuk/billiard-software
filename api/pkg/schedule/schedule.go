package schedule

import (
	"fmt"
	"time"

	"github.com/berkaycubuk/billiard_software_api/pkg/order"
	"github.com/berkaycubuk/billiard_software_api/pkg/subscription"
	"github.com/berkaycubuk/billiard_software_api/pkg/payment"
	"github.com/go-co-op/gocron"
)

func RunJobs() {
	s := gocron.NewScheduler(time.UTC)

	// check subscriptions
	s.Every(5).Minutes().Do(func() {
		fmt.Println("CRON: Checking user subscriptions")
		subscription.CheckUserSubscriptions()
	})

	// check waiting orders
	s.Every(10).Minutes().Do(func() {
		fmt.Println("CRON: Checking waiting orders")
		order.CheckWaitingOrders()
	})

	// capture orders
	s.Every(20).Minutes().Do(func() {
		fmt.Println("CRON: Capturing orders")
		payment.CapturePayments()
	})

	s.StartAsync()
}
