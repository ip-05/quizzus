package game

import (
	"errors"

	"github.com/ip-05/quizzus/entity"
)

type Repository interface {
	Get(id int, code string) *entity.Game
	GetByOwner(id int, hidePrivate bool, limit int) *[]entity.Game
	GetFavorite(user int) *[]entity.Game

	Create(e *entity.Game) *entity.Game
	Update(id int, code string, e *entity.Game) *entity.Game
	Delete(e *entity.Game)
	DeleteQuestion(id int)

	ToggleFavorite(e *entity.FavoriteGame) bool
}

type Service struct {
	repo Repository
}

func NewService(gameRepo Repository) *Service {
	return &Service{
		repo: gameRepo,
	}
}

func (s Service) CreateGame(body entity.CreateGame, ownerId uint) (*entity.Game, error) {
	e, err := entity.NewGame(body, ownerId)
	if err != nil {
		return nil, err
	}

	game := s.repo.Create(e)
	return game, nil
}

func (s Service) UpdateGame(body entity.UpdateGame, id int, code string, ownerId uint) (*entity.Game, error) {
	game, err := s.GetGame(id, code)
	if err != nil {
		return nil, err
	}

	if ownerId != game.Owner {
		return nil, errors.New("you shall not pass! (not owner)")
	}

	game.Topic = body.Topic
	game.RoundTime = body.RoundTime
	game.Points = body.Points

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

			for i := 0; i < 4; i++ {
				question.Options = append(question.Options, &entity.Option{Name: x.Options[i].Name, Correct: x.Options[i].Correct})
			}

			err = question.Validate()
			if err != nil {
				return nil, err
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
			s.repo.DeleteQuestion(int(i))
		}
	}

	err = game.Validate()
	if err != nil {
		return nil, err
	}

	e := s.repo.Update(id, code, game)
	return e, nil
}

func (s Service) DeleteGame(id int, code string, userId uint) error {
	game, err := s.GetGame(id, code)
	if err != nil {
		return err
	}

	if userId != game.Owner {
		return errors.New("you shall not pass! (not owner)")
	}

	s.repo.Delete(game)
	return nil
}

func (s Service) GetGame(id int, code string) (*entity.Game, error) {
	e := s.repo.Get(id, code)

	if e.Id == 0 {
		return nil, errors.New("game not found")
	}

	return e, nil
}

func (s Service) GetGamesByOwner(id int, user int, limit int) (*[]entity.Game, error) {
	hidePrivate := true

	if user == id {
		hidePrivate = false
	}

	games := s.repo.GetByOwner(id, hidePrivate, limit)
	return games, nil
}

func (s Service) GetFavoriteGames(user int) (*[]entity.Game, error) {
	games := s.repo.GetFavorite(user)
	return games, nil
}

func (s Service) Favorite(id int, userId int) bool {
	favorite := &entity.FavoriteGame{
		GameId: uint(id),
		UserId: uint(userId),
	}

	toggle := s.repo.ToggleFavorite(favorite)
	return toggle
}
