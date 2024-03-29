package session

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/api/middleware"
	"github.com/ip-05/quizzus/entity"
)

type Service interface {
	GetSession(ID, userID int) *entity.GameSession
	GetSessions(userID, limit int) *[]entity.GameSession
	NewSession(ID, userID, instID int) uint
	EndSession(ID, userID, instID, questions, players int, points float64) uint
}

type Controller struct {
	service Service
}

func NewController(sessionSvc Service) *Controller {
	return &Controller{service: sessionSvc}
}

func (c Controller) GetSessions(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	if limit == 0 {
		limit = 10
	}

	authedUser, _ := ctx.Get("authedUser")
	user := authedUser.(middleware.AuthedUser)

	sessions := c.service.GetSessions(int(user.ID), limit)

	ctx.JSON(http.StatusOK, sessions)
}

func (c Controller) GetSession(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	authedUser, _ := ctx.Get("authedUser")
	user := authedUser.(middleware.AuthedUser)

	session := c.service.GetSession(id, int(user.ID))

	if session.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	ctx.JSON(http.StatusOK, session)
}
