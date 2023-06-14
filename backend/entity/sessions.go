package entity

import (
	"crypto/rand"
	"math/big"
	"time"
)

type GameSession struct {
	Id         uint      `json:"id" gorm:"primary_key"`
	GameId     uint      `json:"gameId"`
	UserId     uint      `json:"userId"`
	InstanceId uint      `json:"-"`
	Place      int       `json:"place"`
	Points     float64   `json:"points"`
	Questions  int       `json:"questions"`
	Players    int       `json:"players"`
	StartedAt  time.Time `json:"startedAt" gorm:"default:current_timestamp"`
	EndedAt    time.Time `json:"endedAt"`
}

func NewSession(gameId, userId uint) *GameSession {
	random, _ := rand.Int(rand.Reader, big.NewInt(10000000000000))
	session := &GameSession{
		GameId:     gameId,
		UserId:     userId,
		InstanceId: uint(random.Uint64()),
		StartedAt:  time.Now(),
	}
	return session
}
