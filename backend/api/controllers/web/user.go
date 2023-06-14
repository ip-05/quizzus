package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/api/middleware"
	"github.com/ip-05/quizzus/entity"
)

type IUserService interface {
	CreateUser(body *entity.CreateUser) (*entity.User, error)
	UpdateUser(id uint, body entity.UpdateUser) (*entity.User, error)
	DeleteUser(id uint)
	GetUser(id uint) *entity.User
	GetUserByProvider(id string, provider string) *entity.User
}

type UserController struct {
	user IUserService
}

func NewUserController(user IUserService) *UserController {
	return &UserController{user: user}
}

func (u UserController) Me(c *gin.Context) {
	authedUser, _ := c.Get("authedUser")
	user := authedUser.(middleware.AuthedUser)

	dbUser := u.user.GetUser(user.Id)

	c.JSON(http.StatusOK, dbUser)
}

func (u UserController) Update(c *gin.Context) {
	var body entity.UpdateUser

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	authedUser, _ := c.Get("authedUser")
	user := authedUser.(middleware.AuthedUser)

	updatedUser, err := u.user.UpdateUser(user.Id, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (u UserController) Delete(c *gin.Context) {
	authedUser, _ := c.Get("authedUser")
	user := authedUser.(middleware.AuthedUser)

	u.user.DeleteUser(user.Id)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted"})
}
