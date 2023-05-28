package repo

import (
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

func (db *GameStore) Get(id int, code string) *entity.Game {
	var game entity.Game
	db.DB.Preload("Questions.Options").Where("invite_code = ? or id = ?", code, id).First(&game)
	return &game
}

func (db *GameStore) Create(e *entity.Game) *entity.Game {
	db.DB.Session(&gorm.Session{FullSaveAssociations: true}).Create(&e)
	return e
}

func (db *GameStore) Update(id int, code string, e *entity.Game) *entity.Game {
	db.DB.Session(&gorm.Session{FullSaveAssociations: true}).Where("invite_code = ? or id = ?", code, id).Updates(&e)
	return e
}

func (db *GameStore) Delete(e *entity.Game) {
	for _, v := range e.Questions {
		db.DB.Select(clause.Associations).Unscoped().Delete(&v)
	}
	db.DB.Select(clause.Associations).Unscoped().Delete(&e)
}

func (db *GameStore) DeleteQuestion(id int) {
	db.DB.Select(clause.Associations).Unscoped().Delete(&entity.Question{}, id)
	db.DB.Exec("DELETE FROM options WHERE question_id = ?", id)
}
