package game

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/api/middleware"
	"github.com/ip-05/quizzus/entity"
)

type Service interface {
	CreateGame(body entity.CreateGame, ownerId uint) (*entity.Game, error)
	UpdateGame(body entity.UpdateGame, id int, code string, ownerId uint) (*entity.Game, error)
	DeleteGame(id int, code string, userId uint) error

	GetGame(id int, code string) (*entity.Game, error)
	GetGamesByOwner(id int, user int, limit int) (*[]entity.Game, error)
	GetFavoriteGames(user int) (*[]entity.Game, error)

	Favorite(id int, userId int) bool
}

type Controller struct {
	service Service
}

func NewController(gameSvc Service) *Controller {
	return &Controller{service: gameSvc}
}

func (c Controller) CreateGame(ctx *gin.Context) {
	var body entity.CreateGame

	authedUser, _ := ctx.Get("authedUser")
	user := authedUser.(middleware.AuthedUser)

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game, err := c.service.CreateGame(body, user.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, game)
}

func (c Controller) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	code := ctx.Param("id")

	game, err := c.service.GetGame(id, code)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authedUser, _ := ctx.Get("authedUser")
	user := authedUser.(middleware.AuthedUser)

	if user.Id != game.Owner {
		ctx.JSON(http.StatusOK, gin.H{"message": "Game found", "topic": game.Topic})
		return
	}

	ctx.JSON(http.StatusOK, game)
}

func (c Controller) GetMany(ctx *gin.Context) {
	/*

		/games?owner=1 - Get game by owner id
		If owner == authed user, then it will return all games both private and public

	*/

	owner := ctx.Query("owner")
	favorite, err := strconv.ParseBool(ctx.Query("favorite"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "value for favorite must be true (1) or false (0)"})
		return
	}
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	if limit == 0 {
		limit = 10
	}

	authedUser, _ := ctx.Get("authedUser")
	user := authedUser.(middleware.AuthedUser)

	if favorite {
		games, err := c.service.GetFavoriteGames(int(user.Id))
		if err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, games)
		return
	} else if owner != "" {
		id, _ := strconv.Atoi(owner)

		games, err := c.service.GetGamesByOwner(id, int(user.Id), limit)
		if err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, games)
		return
	}

	ctx.JSON(http.StatusOK, []*entity.Game{})
}

func (c Controller) Favorite(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	authedUser, _ := ctx.Get("authedUser")
	user := authedUser.(middleware.AuthedUser)

	toggle := c.service.Favorite(id, int(user.Id))

	ctx.JSON(http.StatusOK, gin.H{"favorite": toggle})
}

func (c Controller) Update(ctx *gin.Context) {
	var body entity.UpdateGame

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	id, _ := strconv.Atoi(ctx.Query("id"))
	code := ctx.Query("invite_code")

	authedUser, _ := ctx.Get("authedUser")
	user := authedUser.(middleware.AuthedUser)

	game, err := c.service.UpdateGame(body, id, code, user.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, game)
}

func (c Controller) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Query("id"))
	code := ctx.Query("invite_code")

	authedUser, _ := ctx.Get("authedUser")
	user := authedUser.(middleware.AuthedUser)

	err := c.service.DeleteGame(id, code, user.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully deleted"})
}
