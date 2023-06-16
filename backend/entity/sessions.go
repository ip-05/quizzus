package entity

import (
	"time"
)

type GameSession struct {
	Id          uint           `json:"id" gorm:"primary_key"`
	GameId      uint           `json:"gameId"`
	UserId      uint           `json:"userId"`
	InstanceId  uint           `json:"-"`
	Points      float64        `json:"points"`
	Questions   int            `json:"questions"`
	Players     int            `json:"players"`
	Game        Game           `json:"game"`
	Leaderboard *[]Leaderboard `json:"leaderboard" gorm:"-"`
	StartedAt   time.Time      `json:"startedAt" gorm:"default:current_timestamp"`
	EndedAt     time.Time      `json:"endedAt"`
}

type Leaderboard struct {
	Name   string  `json:"name"`
	UserId uint    `json:"userId"`
	Points float64 `json:"points"`
}

func NewSession(gameId, userId, instId uint) *GameSession {
	session := &GameSession{
		GameId:     gameId,
		UserId:     userId,
		InstanceId: instId,
		StartedAt:  time.Now(),
	}
	return session
}
