package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Game struct {
	ID		uint		`json:"id" gorm:"primaryKey"`
	TableID		uint		`json:"table_id"`
	StartedAt	time.Time	`json:"started_at"`
	EndedAt		*time.Time	`json:"ended_at"`
	GameUsers	[]GameUser	`json:"game_users"`
	GameHistories	[]GameHistory	`json:"game_histories"`
}

type GameUser struct {
	ID		uint		`json:"id" gorm:"primaryKey"`
	UserID		*uint		`json:"user_id"`
	GameID		uint		`json:"game_id"`
	Name		string		`gorm:"type:varchar(255)" json:"name"`
	Status		uint8		`json:"status"`
	StartedAt	time.Time	`json:"started_at"`
	EndedAt		*time.Time	`json:"ended_at"`
	GameUserChunks	[]GameUserChunk `json:"game_user_chunks"`
}

const (
	GAME_USER_PLAYING uint8 = 1
	GAME_USER_PAUSED	= 2
)

type GameUserChunk struct {
	ID		uint		`json:"id" gorm:"primaryKey"`
	GameID		uint		`json:"game_id"`
	GameUserID	uint		`json:"game_user_id"`
	Price		*decimal.Decimal	`gorm:"type:decimal(10,2)" json:"price"`
	Status		uint8		`json:"status"`
	Players		uint8		`json:"players"`
	StartedAt	time.Time	`json:"started_at"`
	EndedAt		*time.Time	`json:"ended_at"`
}

const (
	GAME_HISTORY_JOIN uint8 = 1
	GAME_HISTORY_LEAVE      = 2
	GAME_HISTORY_PAUSE      = 3
	GAME_HISTORY_UNPAUSE    = 4
	GAME_HISTORY_TRANSFER   = 5
)

type GameHistory struct {
	ID		uint		`json:"id" gorm:"primaryKey"`
	GameID		uint		`json:"game_id"`
	GameUserID	uint		`json:"game_user_id"`
	Action		uint8		`json:"action"`
	CreatedAt	time.Time	`json:"created_at"`
	GameUser	*GameUser	`json:"game_user"`
}

type GameUserOrder struct {
	ID		uint		`json:"id" gorm:"primaryKey"`
	OrderID		uint		`json:"order_id"`
	GameUserID	uint		`json:"game_user_id"`
}
