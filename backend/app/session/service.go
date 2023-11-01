package session

import (
	"time"

	"github.com/ip-05/quizzus/entity"
)

type Repository interface {
	CreateSession(e *entity.GameSession) *entity.GameSession
	EndSession(e *entity.GameSession) *entity.GameSession
	GetSessions(userId, limit int) *[]entity.GameSession
	GetSession(id, userId int) *entity.GameSession
}

type Service struct {
	repo Repository
}

func NewSessionService(sessionRepo Repository) *Service {
	return &Service{
		repo: sessionRepo,
	}
}

func (s Service) NewSession(id, userId, instId int) uint {
	newSession := entity.NewSession(uint(id), uint(userId), uint(instId))
	session := s.repo.CreateSession(newSession)
	return session.UserId
}

func (s Service) EndSession(id, userId, instId, questions, players int, points float64) uint {
	newSession := entity.NewSession(uint(id), uint(userId), uint(instId))
	newSession.Questions = questions
	newSession.Players = players
	newSession.Points = points
	newSession.EndedAt = time.Now()

	session := s.repo.EndSession(newSession)
	return session.UserId
}

func (s Service) GetSessions(userId, limit int) *[]entity.GameSession {
	return s.repo.GetSessions(userId, limit)
}

func (s Service) GetSession(id, userId int) *entity.GameSession {
	return s.repo.GetSession(id, userId)
}
