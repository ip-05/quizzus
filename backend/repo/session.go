package repo

import (
	"github.com/ip-05/quizzus/entity"
	"gorm.io/gorm"
)

type SessionStore struct {
	DB *gorm.DB
}

func NewSessionStore(db *gorm.DB) *SessionStore {
	return &SessionStore{
		DB: db,
	}
}

func (db *SessionStore) GetSessions(userId, limit int) *[]entity.GameSession {
	var sessions []entity.GameSession

	db.DB.Preload("Game").Select("*").Where("user_id = ?", userId).Limit(limit).Order("id DESC").Find(&sessions)

	for i := 0; i < len(sessions); i++ {
		var leaderboard *[]entity.Leaderboard
		db.DB.Model(&entity.GameSession{}).Select("users.id, users.name, points").
			Joins("INNER JOIN users ON users.id = user_id").
			Where("game_sessions.instance_id = ?", sessions[0].InstanceId).
			Order("points DESC").
			Scan(&leaderboard)

		sessions[i].Leaderboard = leaderboard
	}

	return &sessions
}

func (db *SessionStore) GetSession(id, userId int) *entity.GameSession {
	var session entity.GameSession

	db.DB.Preload("Game").Select("*").Where("user_id = ? and id = ?", userId, id).First(&session)

	var leaderboard *[]entity.Leaderboard
	db.DB.Model(&entity.GameSession{}).Select("users.id, users.name, points").
		Joins("INNER JOIN users ON users.id = user_id").
		Where("game_sessions.instance_id = ?", session.InstanceId).
		Order("points DESC").
		Scan(&leaderboard)

	session.Leaderboard = leaderboard

	return &session
}

func (db *SessionStore) CreateSession(e *entity.GameSession) *entity.GameSession {
	db.DB.Create(&e)
	return e
}

func (db *SessionStore) EndSession(e *entity.GameSession) *entity.GameSession {
	db.DB.Where("user_id = ? and game_id = ?", e.UserId, e.GameId).Updates(&e)
	return e
}
