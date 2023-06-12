package web

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/api/middleware"
	"github.com/ip-05/quizzus/entity"
)

type IGameService interface {
	CreateGame(body entity.CreateGame, ownerId uint) (*entity.Game, error)
	UpdateGame(body entity.UpdateBody, id int, code string, ownerId uint) (*entity.Game, error)
	DeleteGame(id int, code string, userId uint) error
	GetGame(id int, code string) (*entity.Game, error)
}

type GameController struct {
	game IGameService
}

func NewGameController(game IGameService) *GameController {
	return &GameController{game: game}
}

func (g *GameController) CreateGame(c *gin.Context) {
	var body entity.CreateGame

	authedUser, _ := c.Get("authedUser")
	user := authedUser.(middleware.AuthedUser)

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game, err := g.game.CreateGame(body, user.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, game)
}

func (g GameController) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	code := c.Query("invite_code")

	game, err := g.game.GetGame(id, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authedUser, _ := c.Get("authedUser")
	user := authedUser.(middleware.AuthedUser)

	if user.Id != game.Owner {
		c.JSON(http.StatusOK, gin.H{"message": "Game found", "topic": game.Topic})
		return
	}

	c.JSON(http.StatusOK, game)
}

func (g GameController) Update(c *gin.Context) {
	var body entity.UpdateBody

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	id, _ := strconv.Atoi(c.Query("id"))
	code := c.Query("invite_code")

	authedUser, _ := c.Get("authedUser")
	user := authedUser.(middleware.AuthedUser)

	game, err := g.game.UpdateGame(body, id, code, user.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, game)
}

func (g GameController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	code := c.Query("invite_code")

	authedUser, _ := c.Get("authedUser")
	user := authedUser.(middleware.AuthedUser)

	err := g.game.DeleteGame(id, code, user.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted"})
}
