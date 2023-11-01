package game

import (
	"github.com/ip-05/quizzus/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r Repository) GetGame(ID int, code string) *entity.Game {
	var game entity.Game
	r.DB.Preload("Questions.Options").Where("invite_code = ? or id = ?", code, ID).First(&game)
	return &game
}

func (r Repository) GetFavoriteGame(user int) *[]entity.Game {
	var games *[]entity.Game
	r.DB.
		Preload("Questions.Options").
		Select("games.*").
		Joins("INNER JOIN favorite_games ON favorite_games.game_id = games.id").
		Where("favorite_games.user_id = ?", user).Find(&games)
	return games
}

func (r Repository) GetGameByOwner(ID int, hidePrivate bool, limit int) *[]entity.Game {
	var games *[]entity.Game

	hide := ""
	if hidePrivate {
		hide = " and public = true"
	}

	r.DB.Preload("Questions.Options").Where("owner = ?"+hide, ID).Limit(limit).Find(&games)
	return games
}

func (r Repository) CreateGame(e *entity.Game) *entity.Game {
	r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Create(&e)
	return e
}

func (r Repository) UpdateGame(ID int, code string, e *entity.Game) *entity.Game {
	r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Where("invite_code = ? or id = ?", code, ID).Updates(&e)
	return e
}

func (r Repository) DeleteGame(e *entity.Game) {
	for _, v := range e.Questions {
		r.DB.Select(clause.Associations).Unscoped().Delete(&v)
	}
	r.DB.Select(clause.Associations).Unscoped().Delete(&e)
}

func (r Repository) DeleteQuestion(ID int) {
	r.DB.Select(clause.Associations).Unscoped().Delete(&entity.Question{}, ID)
	r.DB.Exec("DELETE FROM options WHERE question_id = ?", ID)
}

func (r Repository) ToggleFavoriteGame(e *entity.FavoriteGame) bool {
	favorite := entity.FavoriteGame{}
	r.DB.Where("favorite_games.game_id = ? and favorite_games.user_id = ?", e.GameID, e.UserID).First(&favorite)

	if favorite.ID == 0 {
		r.DB.Create(&e)
		return true
	}

	r.DB.Where("favorite_games.game_id = ? and favorite_games.user_id = ?", e.GameID, e.UserID).Delete(&favorite)
	return false
}
