package entity

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type Game struct {
	Id         uint       `json:"id" gorm:"primary_key"`
	InviteCode string     `json:"inviteCode"`
	Topic      string     `json:"topic"`
	RoundTime  int        `json:"roundTime"`
	Points     float64    `json:"points"`
	Questions  []Question `json:"questions"`
	Owner      string     `json:"ownerId"`
}

type CreateBody struct {
	Topic     string           `json:"topic"`
	RoundTime int              `json:"roundTime"`
	Points    float64          `json:"points"`
	Questions []CreateQuestion `json:"questions"`
}

type UpdateBody struct {
	Topic     string           `json:"topic"`
	RoundTime int              `json:"roundTime"`
	Points    float64          `json:"points"`
	Questions []UpdateQuestion `json:"questions"`
}

func NewGame(topic string, roundTime int, points float64, owner string) (*Game, error) {
	game := &Game{
		InviteCode: generateCode(),
		Topic:      topic,
		RoundTime:  roundTime,
		Points:     points,
		Owner:      owner,
	} // what to do with questions?

	// some checks ?
	return game, nil
}

func generateCode() string {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	code := hex.EncodeToString(bytes)

	return fmt.Sprintf("%s-%s", code[:4], code[4:])
}
