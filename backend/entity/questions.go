package entity

import (
	"errors"
)

type Question struct {
	ID      uint      `json:"id" gorm:"primary_key"`
	Name    string    `json:"name"`
	Options []*Option `json:"options"`
	GameID  uint      `json:"-"`
}

type CreateQuestion struct {
	Name    string         `json:"name"`
	Options []CreateOption `json:"options"`
}

type UpdateQuestion struct {
	ID      uint           `json:"id"`
	Name    string         `json:"name"`
	Options []UpdateOption `json:"options"`
}

func NewQuestion(q CreateQuestion) (*Question, error) {
	question := &Question{
		Name: q.Name,
	}

	for _, o := range q.Options {
		option, err := NewOption(o)
		if err != nil {
			return nil, err
		}
		question.Options = append(question.Options, option)
	}

	if err := question.Validate(); err != nil {
		return nil, err
	}

	return question, nil
}

func (q *Question) Validate() error {
	len := len(q.Options)
	if len != 2 && len != 4 {
		return errors.New("should be 2 or 4 options")
	}
	return nil
}
