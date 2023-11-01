package game

import (
	"testing"

	"github.com/ip-05/quizzus/entity"
	"github.com/stretchr/testify/assert"
)

var testGameBody = entity.CreateGame{
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

	repo := NewRepository(db)

	newGame, err := entity.NewGame(testGameBody, uint(1))
	assert.Nil(t, err)

	game := repo.CreateGame(newGame)
	assert.Greater(t, game.ID, uint(0))
}

func TestRepo_GetGame(t *testing.T) {
	db, cleanup := SetupIntegration(t)
	defer cleanup()

	repo := NewRepository(db)

	newGame, err := entity.NewGame(testGameBody, uint(1))
	assert.Nil(t, err)

	game := repo.CreateGame(newGame)
	assert.Greater(t, game.ID, uint(0))

	gotGame := repo.GetGame(int(game.ID), game.InviteCode)
	assert.Equal(t, gotGame.Topic, newGame.Topic)
}

func TestRepo_DeleteGame(t *testing.T) {
	db, cleanup := SetupIntegration(t)
	defer cleanup()

	repo := NewRepository(db)

	newGame, err := entity.NewGame(testGameBody, uint(1))
	assert.Nil(t, err)

	game := repo.CreateGame(newGame)
	assert.Greater(t, game.ID, uint(0))

	repo.DeleteGame(game)

	gotGame := repo.GetGame(1, "")
	assert.Equal(t, uint(0), gotGame.ID)
}

func TestRepo_UpdateGame(t *testing.T) {
	db, cleanup := SetupIntegration(t)
	defer cleanup()

	repo := NewRepository(db)

	newGame, err := entity.NewGame(testGameBody, uint(1))
	assert.Nil(t, err)

	game := repo.CreateGame(newGame)
	assert.Greater(t, game.ID, uint(0))

	game.Topic = "Updated topic"
	updatedGame := repo.UpdateGame(int(game.ID), game.InviteCode, game)
	assert.Equal(t, updatedGame.Topic, "Updated topic")
}
