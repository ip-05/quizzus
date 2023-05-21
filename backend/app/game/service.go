package game

import "github.com/ip-05/quizzus/entity"

type GameService struct {
	gameRepo     IGame
	questionRepo IQuestion
	optionRepo   IOption
}

func NewGameService(gameR IGame, questionR IQuestion, optionR IOption) *GameService {
	return &GameService{
		gameRepo:     gameR,
		questionRepo: questionR,
		optionRepo:   optionR,
	}
}

func (gs *GameService) CreateGame() {

}

func (gs *GameService) UpdateGame(e *entity.Game) {

}

func (gs *GameService) DeleteGame(e *entity.Game) {

}

func (gs *GameService) GetGame(id int) {

}

func (gs *GameService) GetGameByCode(code string) {

}
