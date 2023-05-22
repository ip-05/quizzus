package game

import (
	"errors"

	"github.com/ip-05/quizzus/entity"
)

type GameService struct {
	gameRepo IGame
}

func NewGameService(gameR IGame) *GameService {
	return &GameService{
		gameRepo: gameR,
	}
}

func (gs *GameService) CreateGame(body entity.CreateBody, ownerId string) (*entity.Game, error) {
	e, err := entity.NewGame(body, ownerId)
	if err != nil {
		return nil, err
	}

	game := gs.gameRepo.Create(e)
	return game, nil
}

func (gs *GameService) UpdateGame(body entity.UpdateBody, id int, code, ownerId string) (*entity.Game, error) {
	game, err := gs.gameRepo.Get(id, code)
	if err != nil {
		return nil, err
	}

	if ownerId != game.Owner {
		return nil, errors.New("you shall not pass! (not owner)")
	}

	game.Topic = body.Topic
	game.RoundTime = body.RoundTime
	game.Points = body.Points

	err = game.Validate()
	if err != nil {
		return nil, err
	}

	ids := make(map[uint]int)
	// assign each question id from existing game a 1
	for _, y := range game.Questions {
		ids[y.Id] += 1
	}

	for i, x := range body.Questions {
		// assign each question id from update a +1
		ids[x.Id] += 1

		// if question ids match => assign new values to question and it's options
		if ids[x.Id] == 2 {
			game.Questions[i].Name = x.Name

			for j := 0; j < 4; j++ {
				game.Questions[i].Options[j].Name = x.Options[j].Name
				game.Questions[i].Options[j].Correct = x.Options[j].Correct
			}
		} else {
			// if question ids don't match (question doesn't already exist) => add a new question to game
			question := entity.Question{
				Name: x.Name,
			}

			err = question.Validate()
			if err != nil {
				return nil, err
			}

			for i := 0; i < 4; i++ {
				question.Options = append(question.Options, &entity.Option{Name: x.Options[i].Name, Correct: x.Options[i].Correct})
			}

			game.Questions = append(game.Questions, &question)
			ids[x.Id] += 1
		}
	}

	// if questions isn't in update body => delete it from the game
	for i, v := range ids {
		if v == 1 {
			for j, v2 := range game.Questions {
				if v2.Id == i {
					game.Questions = append(game.Questions[:j], game.Questions[j+1:]...)
				}
			}
			gs.gameRepo.DeleteQuestion(int(i))
		}
	}

	e, err := gs.gameRepo.Update(id, code, game)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (gs *GameService) DeleteGame(id int, code, userId string) error {
	err := gs.gameRepo.Delete(id, code, userId)
	if err != nil {
		return err
	}
	return nil
}

func (gs *GameService) GetGame(id int, code string) (*entity.Game, error) {
	e, err := gs.gameRepo.Get(id, code)
	if err != nil {
		return nil, err
	}
	return e, nil
}
