package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/api/middleware"
	"github.com/ip-05/quizzus/entity"
)

type Service interface {
	CreateUser(body *entity.CreateUser) (*entity.User, error)
	UpdateUser(id uint, body entity.UpdateUser) (*entity.User, error)
	DeleteUser(id uint)
	GetUser(id uint) *entity.User
	GetUserByProvider(id string, provider string) *entity.User
}

type Controller struct {
	service Service
}

func NewController(userSvc Service) *Controller {
	return &Controller{service: userSvc}
}

func (c Controller) Get(ctx *gin.Context) {
	var userId int
	var err error

	id := ctx.Param("id")
	if id == "me" {
		authedUser, _ := ctx.Get("authedUser")
		user := authedUser.(middleware.AuthedUser)
		userId = int(user.Id)
	} else {
		userId, err = strconv.Atoi(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}

	dbUser := c.service.GetUser(uint(userId))
	ctx.JSON(http.StatusOK, dbUser)
}

func (c Controller) Update(ctx *gin.Context) {
	var body entity.UpdateUser

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	authedUser, _ := ctx.Get("authedUser")
	user := authedUser.(middleware.AuthedUser)

	updatedUser, err := c.service.UpdateUser(user.Id, body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, updatedUser)
}

func (c Controller) Delete(ctx *gin.Context) {
	authedUser, _ := ctx.Get("authedUser")
	user := authedUser.(middleware.AuthedUser)

	c.service.DeleteUser(user.Id)

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully deleted"})
}
