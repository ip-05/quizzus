package game

import "github.com/ip-05/quizzus/entity"

type IGame interface {
	Get(id int) (*entity.Game, error)
	GetByCode(code string) (*entity.Game, error)
	Create(e *entity.Game) (*entity.Game, error)
	Update(e *entity.Game) (*entity.Game, error)
	Delete(e *entity.Game) error
}

type IQuestion interface {
	Get(id int) (*entity.Question, error)
	GetByGameId(gameId int) (*[]entity.Question, error)
	Create(e *entity.Question) (*entity.Question, error)
	Update(e *entity.Question) (*entity.Question, error)
	Delete(e *entity.Question) error
}

type IOption interface {
	Get(id int) (*entity.Option, error)
	GetByGameId(gameId int) (*[]entity.Option, error)
	Create(e *entity.Option) (*entity.Option, error)
	Update(e *entity.Option) (*entity.Option, error)
	Delete(e *entity.Option) error
}

type IService interface {
	CreateGame()
	UpdateGame()
	DeleteGame()
	GetGame()
	GetGameByCode()
}
