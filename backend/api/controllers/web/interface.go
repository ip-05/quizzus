package web

import "github.com/ip-05/quizzus/entity"

type IGameService interface {
	CreateGame(body entity.CreateBody, ownerId string) (*entity.Game, error)
	UpdateGame(body entity.UpdateBody, id int, code, ownerId string) (*entity.Game, error)
	DeleteGame(id int, code, userId string) error
	GetGame(id int, code string) (*entity.Game, error)
}
