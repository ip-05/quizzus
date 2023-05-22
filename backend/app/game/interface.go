package game

import "github.com/ip-05/quizzus/entity"

type IGame interface {
	Get(id int, code string) (*entity.Game, error)
	Create(e *entity.Game) *entity.Game
	Update(id int, code string, e *entity.Game) (*entity.Game, error)
	Delete(id int, code, userId string) error
	DeleteQuestion(id int)
}

type IService interface {
	CreateGame(body entity.CreateBody, ownerId string) (*entity.Game, error)
	UpdateGame(body entity.UpdateBody, id int, code, ownerId string) (*entity.Game, error)
	DeleteGame(id int, code, userId string) error
	GetGame(id int, code string) (*entity.Game, error)
}
