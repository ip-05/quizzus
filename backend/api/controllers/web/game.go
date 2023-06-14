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
	UpdateGame(body entity.UpdateGame, id int, code string, ownerId uint) (*entity.Game, error)
	DeleteGame(id int, code string, userId uint) error

	GetGame(id int, code string) (*entity.Game, error)
	GetGamesByOwner(id int, user int, limit int) (*[]entity.Game, error)
	GetFavoriteGames(user int) (*[]entity.Game, error)

	Favorite(id int, userId int) bool
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
	id, _ := strconv.Atoi(c.Param("id"))
	code := c.Param("id")

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

func (g GameController) GetMany(c *gin.Context) {
	/*

		/games?owner=1 - Get game by owner id
		If owner == authed user, then it will return all games both private and public

	*/

	owner := c.Query("owner")
	favorite, _ := strconv.ParseBool(c.Query("favorite"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit == 0 {
		limit = 10
	}

	authedUser, _ := c.Get("authedUser")
	user := authedUser.(middleware.AuthedUser)

	if favorite {
		games, err := g.game.GetFavoriteGames(int(user.Id))
		if err != nil {
			c.JSON(http.StatusOK, err.Error())
			return
		}

		c.JSON(http.StatusOK, games)
		return
	} else if owner != "" {
		id, _ := strconv.Atoi(owner)

		games, err := g.game.GetGamesByOwner(id, int(user.Id), limit)
		if err != nil {
			c.JSON(http.StatusOK, err.Error())
			return
		}

		c.JSON(http.StatusOK, games)
		return
	}

	c.JSON(http.StatusOK, []*entity.Game{})
}

func (g GameController) Favorite(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	authedUser, _ := c.Get("authedUser")
	user := authedUser.(middleware.AuthedUser)

	toggle := g.game.Favorite(id, int(user.Id))

	c.JSON(http.StatusOK, gin.H{"favorite": toggle})
}

func (g GameController) Update(c *gin.Context) {
	var body entity.UpdateGame

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
