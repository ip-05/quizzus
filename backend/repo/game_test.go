package repo

import (
	"testing"

	"github.com/ip-05/quizzus/entity"
	"github.com/stretchr/testify/assert"
)

var testGameBody = entity.CreateBody{
	Topic:     "My game",
	RoundTime: 10,
	Points:    3,
	Questions: []entity.CreateQuestion{
		{
			Name: "What color is tomato?",
			Options: []entity.CreateOption{
				{Name: "Red", Correct: true},
				{Name: "Green", Correct: false},
				{Name: "Blue", Correct: false},
				{Name: "Orange", Correct: false},
			},
		},
	},
}

func TestRepo_CreateGame(t *testing.T) {
	db, cleanup := SetupIntegration(t)
	defer cleanup()

	repo := NewGameStore(db)

	newGame, err := entity.NewGame(testGameBody, "1")
	assert.Nil(t, err)

	game := repo.Create(newGame)
	assert.Greater(t, game.Id, uint(0))
}

func TestRepo_GetGame(t *testing.T) {
	db, cleanup := SetupIntegration(t)
	defer cleanup()

	repo := NewGameStore(db)

	newGame, err := entity.NewGame(testGameBody, "1")
	assert.Nil(t, err)

	game := repo.Create(newGame)
	assert.Greater(t, game.Id, uint(0))

	gotGame, err := repo.Get(int(game.Id), game.InviteCode)
	assert.Nil(t, err)
	assert.Equal(t, gotGame.Topic, newGame.Topic)
}

func TestRepo_DeleteGame(t *testing.T) {
	db, cleanup := SetupIntegration(t)
	defer cleanup()

	repo := NewGameStore(db)

	newGame, err := entity.NewGame(testGameBody, "1")
	assert.Nil(t, err)

	game := repo.Create(newGame)
	assert.Greater(t, game.Id, uint(0))

	err = repo.Delete(int(game.Id), game.InviteCode, "1")
	assert.Nil(t, err)
}

func TestRepo_UpdateGame(t *testing.T) {
	db, cleanup := SetupIntegration(t)
	defer cleanup()

	repo := NewGameStore(db)

	newGame, err := entity.NewGame(testGameBody, "1")
	assert.Nil(t, err)

	game := repo.Create(newGame)
	assert.Greater(t, game.Id, uint(0))

	game.Topic = "Updated topic"
	updatedGame := repo.Update(int(game.Id), game.InviteCode, game)
	assert.Equal(t, updatedGame.Topic, "Updated topic")
}
