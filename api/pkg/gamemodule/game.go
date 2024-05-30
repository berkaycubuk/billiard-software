package gamemodule

import (
	"time"

	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/models"
	"github.com/berkaycubuk/billiard_software_api/pkg/gameuser"
	"github.com/shopspring/decimal"
)

func SyncUserChunksForUnpause(gameUserID uint, gameID uint) error {
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

func SyncUserChunksForPause(gameUserID uint, gameID uint) error {
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

func SyncUserChunksForTransfer(IDToExclude uint, gameID uint, oldGameID uint) error {
	// keep the current time
	now := time.Now()

	var savedPrice decimal.Decimal

	// fill price and ended_at values for current chunks
	var chunks []models.GameUserChunk
	database.DB.Where("game_id = ? AND ended_at IS NULL", oldGameID).Find(&chunks)
	for _, v := range chunks {
		price, err := gameuser.CalcUserChunkPrice(&v, now)
		if err != nil {
			return err
		}

		if v.GameUserID == IDToExclude {
			savedPrice = *price
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
			priceForChunk = savedPrice
		}

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

func SyncUserChunksForJoin(IDToExclude uint, gameID uint) error {
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

func SyncUserChunksForLeave(gameUserID uint, gameID uint) error {
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

func create(tableID uint) error {
	game := models.Game{
		TableID:   tableID,
		StartedAt: time.Now(),
	}
	err := database.DB.Create(&game).Error
	if err != nil {
		return err
	}

	return nil
}

func CreateHistory(gameUserID uint, gameID uint, action uint8) {
	database.DB.Create(&models.GameHistory{
		GameID:     gameID,
		GameUserID: gameUserID,
		Action:     action,
		CreatedAt:  time.Now(),
	})
}
