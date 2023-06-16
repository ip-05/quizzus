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

func (db *GameStore) GetFavorite(user int) *[]entity.Game {
	var games *[]entity.Game
	db.DB.
		Preload("Questions.Options").
		Select("games.*").
		Joins("INNER JOIN favorite_games ON favorite_games.game_id = games.id").
		Where("favorite_games.user_id = ?", user).Find(&games)
	return games
}

func (db *GameStore) GetByOwner(id int, hidePrivate bool, limit int) *[]entity.Game {
	var games *[]entity.Game

	hide := ""
	if hidePrivate {
		hide = " and public = true"
	}

	db.DB.Preload("Questions.Options").Where("owner = ?"+hide, id).Limit(limit).Find(&games)
	return games
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

func (db *GameStore) ToggleFavorite(e *entity.FavoriteGame) bool {
	favorite := entity.FavoriteGame{}
	db.DB.Where("favorite_games.game_id = ? and favorite_games.user_id = ?", e.GameId, e.UserId).First(&favorite)

	if favorite.Id == 0 {
		db.DB.Create(&e)
		return true
	}

	db.DB.Where("favorite_games.game_id = ? and favorite_games.user_id = ?", e.GameId, e.UserId).Delete(&favorite)
	return false
}
