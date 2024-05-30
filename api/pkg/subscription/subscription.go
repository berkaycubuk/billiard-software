package subscription

import (
	"fmt"
	"time"

	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/models"
	"github.com/berkaycubuk/billiard_software_api/pkg/notification"
)

func ActivateSubscription(userID uint, subscriptionID uint) error {
	var userSubscription models.UserSubscription
	err := database.DB.Where("user_id = ? AND deleted_at IS NULL", userID).First(&userSubscription).Error
	if err == nil {
		//return errors.New("subscription exists")
		// delete active subscription
		now := time.Now()
		database.DB.Model(&userSubscription).Update("deleted_at", &now)
		database.DB.Model(&userSubscription).Update("status", models.USER_SUBSCRIPTION_STATUS_ENDED)
	}

	var subscription models.Subscription
	err = database.DB.Where("id = ?", subscriptionID).First(&subscription).Error
	if err != nil {
		return err
	}

	userSubscription = models.UserSubscription{
		UserID:         userID,
		SubscriptionID: subscriptionID,
		Name:           subscription.Name,
		Status:         models.USER_SUBSCRIPTION_STATUS_ACTIVE,
		EndingAt:       time.Now().Add(time.Hour * time.Duration(subscription.Hours)),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	err = database.DB.Create(&userSubscription).Error
	if err != nil {
		return err
	}

	userSubscriptionChunk := models.UserSubscriptionChunk{
		UserSubscriptionID: userSubscription.ID,
		Action:             models.USER_SUBSCRIPTION_ACTION_ACTIVATED,
		ActionBy:           &userID,
		CreatedAt:          time.Now(),
	}
	err = database.DB.Create(&userSubscriptionChunk).Error
	if err != nil {
		return err
	}

	return nil
}

func CheckUserSubscriptions() {
	now := time.Now()

	var userSubscriptions []models.UserSubscription
	err := database.DB.Where("status = ? AND deleted_at IS NULL", models.USER_SUBSCRIPTION_STATUS_ACTIVE).Find(&userSubscriptions).Error
	if err != nil {
		return
	}

	for _, v := range userSubscriptions {
		if v.EndingAt.Before(now) { // it's ended
			err := database.DB.Model(&v).Update("status", models.USER_SUBSCRIPTION_STATUS_ENDED).Error
			if err != nil {
				fmt.Println(err.Error())
			}

			notification.Create(
				v.UserID,
				models.NOTIFICATION_MEDIUM_APP,
				models.NOTIFICATION_TYPE_SIMPLE,
				fmt.Sprintf("Your %s subscription is ended.", v.Name),
			)
		}
	}
}
