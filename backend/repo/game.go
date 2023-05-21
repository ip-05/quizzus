package repo

import (
	"errors"

	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/entity"
	"gorm.io/gorm"
)

type GameStore struct {
	DB *gorm.DB
}

func NewGameStore(cfg *config.Config, db *gorm.DB) *GameStore {
	return &GameStore{
		DB: db,
	}
}

func (db *GameStore) Get(id int) (*entity.Game, error) {
	var game entity.Game

	db.DB.Where("id = ?", id).First(&game)

	if game.Id == 0 {
		return nil, errors.New("game not found")
	}

	return &game, nil
}

func (db *GameStore) GetByCode(code string) (*entity.Game, error) {
	var game entity.Game

	db.DB.Where("invite_code = ?", code).First(&game)

	if game.Id == 0 {
		return nil, errors.New("game not found")
	}

	return &game, nil
}

func (db *GameStore) Create(e *entity.Game) (*entity.Game, error) {
	db.DB.Create(&e)
	return e, nil
}

func (db *GameStore) Update(e *entity.Game) (*entity.Game, error) {
	db.DB.Where("id = ?", e.Id).Updates(&e)
	return e, nil
}

func (db *GameStore) Delete(e *entity.Game) error {
	db.DB.Where("id = ?", e.Id).Delete(&e)
	return nil
}
