package gameuser

import (
	"time"

	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/models"
	"github.com/berkaycubuk/billiard_software_api/pkg/pricing"
	"github.com/shopspring/decimal"
)

func Create(userID uint, gameID uint, name string) error {
	gameUser := models.GameUser{
		UserID: &userID,
		GameID: gameID,
		Name:   name,
	}
	err := database.DB.Create(&gameUser).Error
	if err != nil {
	}

	return nil
}

func CalcUserChunkPrice(chunk *models.GameUserChunk, end time.Time) (*decimal.Decimal, error) {
	if chunk.Status == models.GAME_USER_PAUSED {
		zeroDecimal := decimal.New(0, -2)
		return &zeroDecimal, nil
	}

	// time span
	length := end.Sub(chunk.StartedAt)

	perMinute, err := pricing.PerMinutePricingForUser(chunk.GameUserID, chunk.Players)
	if err != nil {
		return nil, err
	}

	price := perMinute.Mul(decimal.NewFromFloat(length.Minutes()))
	return &price, nil
}

func CalcUserTotalPrice(gameID uint, gameUserID uint) decimal.Decimal {
	now := time.Now()
	total := decimal.NewFromInt(0)

	var chunks []models.GameUserChunk
	database.DB.Where("game_id = ? AND game_user_id = ?", gameID, gameUserID).Find(&chunks)
	for _, v := range chunks {
		if v.Status == models.GAME_USER_PAUSED {
			continue
		}

		if v.EndedAt == nil {
			// use role minute pricing thing
			diff := now.Sub(v.StartedAt)

			perMinute, err := pricing.PerMinutePricingForUser(gameUserID, v.Players)

			if err == nil {
				chunkPrice := perMinute.Mul(decimal.NewFromFloat(diff.Minutes()))
				total = total.Add(chunkPrice)
			}
		} else {
			total = total.Add(*v.Price)
		}
	}

	// check for custom 200kr rule
	var gameUser models.GameUser
	err := database.DB.Where("id = ?", gameUserID).First(&gameUser).Error
	if err == nil {
		if gameUser.UserID != nil { // apply limit if not guest
			if total.GreaterThan(decimal.New(200, 0)) {
				total = decimal.New(200, 0)
			}

			// if user have sub just count it 0
			// this is customer specific request
			var userSubscription models.UserSubscription
			err = database.DB.Where("user_id = ? AND status = ?", gameUser.UserID, models.USER_SUBSCRIPTION_STATUS_ACTIVE).First(&userSubscription).Error
			if err == nil {
					total = decimal.Zero
			}
		}
	}

	return total
}

func CalcUserTotalTime(gameID uint, gameUserID uint) float64 {
	now := time.Now()
	minutes := 0.0
	var chunks []models.GameUserChunk
	database.DB.Where("game_id = ? AND game_user_id = ?", gameID, gameUserID).Find(&chunks)
	for _, v := range chunks {
		if v.Status == models.GAME_USER_PAUSED {
			continue
		}

		if v.EndedAt == nil {
			diff := now.Sub(v.StartedAt)
			minutes = minutes + diff.Minutes()
		} else {
			diff := v.EndedAt.Sub(v.StartedAt)
			minutes = minutes + diff.Minutes()
		}
	}

	return minutes
}
