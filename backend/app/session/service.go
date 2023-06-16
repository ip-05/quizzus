package session

import (
	"time"

	"github.com/ip-05/quizzus/entity"
)

type ISessionRepo interface {
	CreateSession(e *entity.GameSession) *entity.GameSession
	EndSession(e *entity.GameSession) *entity.GameSession
	GetSessions(userId, limit int) *[]entity.GameSession
	GetSession(id, userId int) *entity.GameSession
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

func (ss *SessionService) EndSession(id, userId, instId, questions, players int, points float64) uint {
	s := entity.NewSession(uint(id), uint(userId), uint(instId))
	s.Questions = questions
	s.Players = players
	s.Points = points
	s.EndedAt = time.Now()

	session := ss.sessionRepo.EndSession(s)
	return session.UserId
}

func (ss *SessionService) GetSessions(userId, limit int) *[]entity.GameSession {
	sessions := ss.sessionRepo.GetSessions(userId, limit)
	return sessions
}

func (ss *SessionService) GetSession(id, userId int) *entity.GameSession {
	session := ss.sessionRepo.GetSession(id, userId)
	return session
}
