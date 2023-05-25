package entity

import (
	"errors"
)

type Question struct {
	Id      uint      `json:"id" gorm:"primary_key"`
	Name    string    `json:"name"`
	Options []*Option `json:"options"`
	GameID  uint      `json:"-"`
}

type CreateQuestion struct {
	Name    string         `json:"name"`
	Options []CreateOption `json:"options"`
}

type UpdateQuestion struct {
	Id      uint           `json:"id"`
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
	if len(q.Options) != 4 {
		return errors.New("should be 4 options")
	}
	return nil
}
