package session

import (
	"time"

	"github.com/ip-05/quizzus/entity"
)

type Repository interface {
	CreateSession(e *entity.GameSession) *entity.GameSession
	EndSession(e *entity.GameSession) *entity.GameSession
	GetSessions(userID, limit int) *[]entity.GameSession
	GetSession(ID, userID int) *entity.GameSession
}

type Service struct {
	repo Repository
}

func NewSessionService(sessionRepo Repository) *Service {
	return &Service{
		repo: sessionRepo,
	}
}

func (s Service) NewSession(ID, userID, instID int) uint {
	newSession := entity.NewSession(uint(ID), uint(userID), uint(instID))
	session := s.repo.CreateSession(newSession)
	return session.UserID
}

func (s Service) EndSession(ID, userID, instID, questions, players int, points float64) uint {
	newSession := entity.NewSession(uint(ID), uint(userID), uint(instID))
	newSession.Questions = questions
	newSession.Players = players
	newSession.Points = points
	newSession.EndedAt = time.Now()

	session := s.repo.EndSession(newSession)
	return session.UserID
}

func (s Service) GetSessions(userID, limit int) *[]entity.GameSession {
	return s.repo.GetSessions(userID, limit)
}

func (s Service) GetSession(ID, userID int) *entity.GameSession {
	return s.repo.GetSession(ID, userID)
}
