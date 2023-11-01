package game

import (
	"errors"

	"github.com/ip-05/quizzus/entity"
)

type Repository interface {
	GetGame(ID int, code string) *entity.Game
	GetGameByOwner(ID int, hidePrivate bool, limit int) *[]entity.Game
	GetFavoriteGame(user int) *[]entity.Game

	CreateGame(e *entity.Game) *entity.Game
	UpdateGame(ID int, code string, e *entity.Game) *entity.Game
	DeleteGame(e *entity.Game)
	DeleteQuestion(ID int)

	ToggleFavoriteGame(e *entity.FavoriteGame) bool
}

type Service struct {
	repo Repository
}

func NewService(gameRepo Repository) *Service {
	return &Service{
		repo: gameRepo,
	}
}

func (s Service) CreateGame(body entity.CreateGame, ownerID uint) (*entity.Game, error) {
	e, err := entity.NewGame(body, ownerID)
	if err != nil {
		return nil, err
	}

	game := s.repo.CreateGame(e)
	return game, nil
}

func (s Service) UpdateGame(body entity.UpdateGame, ID int, code string, ownerID uint) (*entity.Game, error) {
	game, err := s.GetGame(ID, code)
	if err != nil {
		return nil, err
	}

	if ownerID != game.Owner {
		return nil, errors.New("you shall not pass! (not owner)")
	}

	game.Topic = body.Topic
	game.RoundTime = body.RoundTime
	game.Points = body.Points

	ids := make(map[uint]int)
	// assign each question id from existing game a 1
	for _, y := range game.Questions {
		ids[y.ID] += 1
	}

	for i, x := range body.Questions {
		// assign each question id from update a +1
		ids[x.ID] += 1

		// if question ids match => assign new values to question and it's options
		if ids[x.ID] == 2 {
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

			for i := 0; i < 4; i++ {
				question.Options = append(question.Options, &entity.Option{Name: x.Options[i].Name, Correct: x.Options[i].Correct})
			}

			err = question.Validate()
			if err != nil {
				return nil, err
			}

			game.Questions = append(game.Questions, &question)
			ids[x.ID] += 1
		}
	}

	// if questions isn't in update body => delete it from the game
	for i, v := range ids {
		if v == 1 {
			for j, v2 := range game.Questions {
				if v2.ID == i {
					game.Questions = append(game.Questions[:j], game.Questions[j+1:]...)
				}
			}
			s.repo.DeleteQuestion(int(i))
		}
	}

	err = game.Validate()
	if err != nil {
		return nil, err
	}

	e := s.repo.UpdateGame(ID, code, game)
	return e, nil
}

func (s Service) DeleteGame(ID int, code string, userID uint) error {
	game, err := s.GetGame(ID, code)
	if err != nil {
		return err
	}

	if userID != game.Owner {
		return errors.New("you shall not pass! (not owner)")
	}

	s.repo.DeleteGame(game)
	return nil
}

func (s Service) GetGame(ID int, code string) (*entity.Game, error) {
	e := s.repo.GetGame(ID, code)

	if e.ID == 0 {
		return nil, errors.New("game not found")
	}

	return e, nil
}

func (s Service) GetGamesByOwner(ID int, user int, limit int) (*[]entity.Game, error) {
	hidePrivate := true

	if user == ID {
		hidePrivate = false
	}

	games := s.repo.GetGameByOwner(ID, hidePrivate, limit)
	return games, nil
}

func (s Service) GetFavoriteGames(user int) (*[]entity.Game, error) {
	games := s.repo.GetFavoriteGame(user)
	return games, nil
}

func (s Service) Favorite(ID int, userID int) bool {
	favorite := &entity.FavoriteGame{
		GameID: uint(ID),
		UserID: uint(userID),
	}

	toggle := s.repo.ToggleFavoriteGame(favorite)
	return toggle
}
