package session

import (
	"github.com/ip-05/quizzus/entity"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r Repository) GetSessions(userID, limit int) *[]entity.GameSession {
	var sessions []entity.GameSession

	r.DB.Preload("Game").Select("*").Where("user_id = ?", userID).Limit(limit).Order("id DESC").Find(&sessions)

	for i := 0; i < len(sessions); i++ {
		var leaderboard *[]entity.Leaderboard
		r.DB.Model(&entity.GameSession{}).Select("users.id, users.name, points").
			Joins("INNER JOIN users ON users.id = user_id").
			Where("game_sessions.instance_id = ?", sessions[0].InstanceID).
			Order("points DESC").
			Scan(&leaderboard)

		sessions[i].Leaderboard = leaderboard
	}

	return &sessions
}

func (r Repository) GetSession(ID, userID int) *entity.GameSession {
	var session entity.GameSession

	r.DB.Preload("Game").Select("*").Where("user_id = ? and id = ?", userID, ID).First(&session)

	var leaderboard *[]entity.Leaderboard
	r.DB.Model(&entity.GameSession{}).Select("users.id, users.name, points").
		Joins("INNER JOIN users ON users.id = user_id").
		Where("game_sessions.instance_id = ?", session.InstanceID).
		Order("points DESC").
		Scan(&leaderboard)

	session.Leaderboard = leaderboard

	return &session
}

func (r Repository) CreateSession(e *entity.GameSession) *entity.GameSession {
	r.DB.Create(&e)
	return e
}

func (r Repository) EndSession(e *entity.GameSession) *entity.GameSession {
	r.DB.Where("user_id = ? and game_id = ?", e.UserID, e.GameID).Updates(&e)
	return e
}
