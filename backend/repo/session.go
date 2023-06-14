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

func (db *SessionStore) GetSessions(userId, limit int) (*[]entity.GameSession, *[]entity.Leaderboard) {
	var sessions *[]entity.GameSession
	var leaderboard *[]entity.Leaderboard

	db.DB.Select("user_id, points").Where("user_id = ?", userId).Limit(limit).Find(&sessions)
	return sessions, leaderboard
}

func (db *SessionStore) CreateSession(e *entity.GameSession) *entity.GameSession {
	db.DB.Create(&e)
	return e
}

func (db *SessionStore) EndSession(e *entity.GameSession) *entity.GameSession {
	db.DB.Where("user_id = ? and game_id = ?", e.UserId, e.GameId).Updates(&e)
	return e
}
