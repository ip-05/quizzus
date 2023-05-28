package repo

import (
	"errors"

	"github.com/ip-05/quizzus/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GameStore struct {
	DB *gorm.DB
}

func NewGameStore(db *gorm.DB) *GameStore {
	return &GameStore{
		DB: db,
	}
}

func (db *GameStore) Get(id int, code string) (*entity.Game, error) {
	var game entity.Game
	db.DB.Preload("Questions.Options").Where("invite_code = ? or id = ?", code, id).First(&game)

	if game.Id == 0 {
		return nil, errors.New("game not found")
	}

	return &game, nil
}

func (db *GameStore) Create(e *entity.Game) *entity.Game {
	db.DB.Session(&gorm.Session{FullSaveAssociations: true}).Create(&e)
	return e
}

func (db *GameStore) Update(id int, code string, e *entity.Game) *entity.Game {
	db.DB.Session(&gorm.Session{FullSaveAssociations: true}).Where("invite_code = ? or id = ?", code, id).Updates(&e)
	return e
}

func (db *GameStore) Delete(id int, code string, userId string) error {
	game, err := db.Get(id, code)
	if err != nil {
		return err
	}

	if userId != game.Owner {
		return errors.New("you shall not pass! (not owner)")
	}

	for _, v := range game.Questions {
		db.DB.Select(clause.Associations).Unscoped().Delete(&v)
	}
	db.DB.Select(clause.Associations).Unscoped().Delete(&game)

	return nil
}

func (db *GameStore) DeleteQuestion(id int) {
	db.DB.Select(clause.Associations).Unscoped().Delete(&entity.Question{}, id)
	db.DB.Exec("DELETE FROM options WHERE question_id = ?", id)
}
