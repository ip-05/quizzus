package game

import "github.com/ip-05/quizzus/entity"

type IGameRepo interface {
	Get(id int, code string) (*entity.Game, error)
	Create(e *entity.Game) *entity.Game
	Update(id int, code string, e *entity.Game) (*entity.Game, error)
	Delete(id int, code, userId string) error
	DeleteQuestion(id int)
}
