package entity

import (
	"errors"

	"github.com/ip-05/quizzus/utils"
)

type Game struct {
	Id         uint        `json:"id" gorm:"primary_key"`
	InviteCode string      `json:"inviteCode"`
	Topic      string      `json:"topic"`
	RoundTime  int         `json:"roundTime"`
	Points     float64     `json:"points"`
	Questions  []*Question `json:"questions"`
	Owner      uint        `json:"ownerId"`
}

type CreateGame struct {
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

func NewGame(body CreateGame, ownerId uint) (*Game, error) {
	code := utils.GenerateCode()
	game := &Game{
		InviteCode: code,
		Topic:      body.Topic,
		RoundTime:  body.RoundTime,
		Points:     body.Points,
		Owner:      ownerId,
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
