package entity

import (
	"time"
)

type GameSession struct {
	ID          uint           `json:"id" gorm:"primary_key"`
	GameID      uint           `json:"game_id"`
	UserID      uint           `json:"user_id"`
	InstanceID  uint           `json:"-"`
	Points      float64        `json:"points"`
	Questions   int            `json:"questions"`
	Players     int            `json:"players"`
	Game        Game           `json:"game"`
	Leaderboard *[]Leaderboard `json:"leaderboard" gorm:"-"`
	StartedAt   time.Time      `json:"started_at" gorm:"default:current_timestamp"`
	EndedAt     time.Time      `json:"ended_at"`
}

type Leaderboard struct {
	Name   string  `json:"name"`
	UserID uint    `json:"user_id"`
	Points float64 `json:"points"`
}

func NewSession(gameID, userID, instID uint) *GameSession {
	session := &GameSession{
		GameID:     gameID,
		UserID:     userID,
		InstanceID: instID,
		StartedAt:  time.Now(),
	}
	return session
}
