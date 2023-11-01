package entity

import (
	"errors"
	"time"

	"github.com/ip-05/quizzus/utils"
)

type Game struct {
	ID         uint        `json:"id" gorm:"primary_key"`
	InviteCode string      `json:"invite_code"`
	Topic      string      `json:"topic"`
	RoundTime  int         `json:"round_time"`
	Points     float64     `json:"points"`
	Public     bool        `json:"public"`
	Questions  []*Question `json:"questions"`
	Owner      uint        `json:"owner_id"`
	CreatedAt  time.Time   `json:"created_at" gorm:"default:current_timestamp"`
}

type FavoriteGame struct {
	ID     uint `json:"id" gorm:"primary_key"`
	GameID uint `json:"game_id"`
	UserID uint `json:"user_id"`
}

type CreateGame struct {
	Topic     string           `json:"topic"`
	RoundTime int              `json:"round_time"`
	Points    float64          `json:"points"`
	Public    bool             `json:"public"`
	Questions []CreateQuestion `json:"questions"`
}

type UpdateGame struct {
	Topic     string           `json:"topic"`
	RoundTime int              `json:"round_time"`
	Points    float64          `json:"points"`
	Public    bool             `json:"public"`
	Questions []UpdateQuestion `json:"questions"`
}

func NewGame(body CreateGame, ownerID uint) (*Game, error) {
	code := utils.GenerateCode()
	game := &Game{
		InviteCode: code,
		Topic:      body.Topic,
		RoundTime:  body.RoundTime,
		Points:     body.Points,
		Public:     body.Public,
		Owner:      ownerID,
	}

	for _, q := range body.Questions {
		question, err := NewQuestion(q)
		if err != nil {
			return nil, err
		}

		game.Questions = append(game.Questions, question)
	}

	if err := game.Validate(); err != nil {
		return nil, err
	}

	return game, nil
}

func (g *Game) Validate() error {
	if len(g.Topic) > 128 {
		return errors.New("too long topic name")
	}

	if g.RoundTime < 10 || g.RoundTime > 60 {
		return errors.New("round time should be over 10 or below 60 (seconds)")
	}

	if g.Points <= 0 {
		return errors.New("points should not be lower than 0")
	}

	if len(g.Questions) < 1 {
		return errors.New("should be at least 1 question")
	}
	return nil
}
