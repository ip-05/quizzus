package session

import (
	"time"

	"github.com/ip-05/quizzus/entity"
)

type ISessionRepo interface {
	CreateSession(e *entity.GameSession) *entity.GameSession
	EndSession(e *entity.GameSession) *entity.GameSession
	GetSessions(userId, limit int) (*[]entity.GameSession, *[]entity.Leaderboard)
}

type SessionService struct {
	sessionRepo ISessionRepo
}

func NewSessionService(sessionRepo ISessionRepo) *SessionService {
	return &SessionService{
		sessionRepo: sessionRepo,
	}
}

func (ss *SessionService) NewSession(id, userId, instId int) uint {
	s := entity.NewSession(uint(id), uint(userId), uint(instId))

	session := ss.sessionRepo.CreateSession(s)
	return session.UserId
}

func (ss *SessionService) EndSession(id, userId, instId, place, questions, players int, points float64) uint {
	s := entity.NewSession(uint(id), uint(userId), uint(instId))
	s.Place = place
	s.Questions = questions
	s.Players = players
	s.Points = points
	s.EndedAt = time.Now()

	session := ss.sessionRepo.EndSession(s)
	return session.UserId
}

func (ss *SessionService) GetSessions(userId, limit int) (*[]entity.GameSession, *[]entity.Leaderboard) {
	// TODO: Return game sessions with filled game info and leaderboard
	// Leaderboard needs to be gotten by joining all game sessions and returning user id and points
	// Like this:
	// leaderboard: [{ user: 1, points: 500 }]
	sessions, leaderboard := ss.sessionRepo.GetSessions(userId, limit)
	return sessions, leaderboard
}
