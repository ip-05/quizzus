package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/controllers"
	"github.com/ip-05/quizzus/middleware"
)

func NewRouter() *gin.Engine {
	cfg := config.GetConfig()

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.Frontend.Base},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	auth := new(controllers.AuthController)
	game := new(controllers.GameController)
	ws := new(controllers.WebSocketController)

	authGroup := router.Group("auth")

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	authGroup.GET("/google", auth.GoogleLogin)
	authGroup.GET("/google/callback", auth.GoogleCallback)

	authGroup.Use(middleware.AuthMiddleware())
	authGroup.GET("/me", auth.Me)

	//router.GET("/game/:id", game.GetById)

	router.GET("/games", game.Get)
	router.POST("/games", game.CreateGame)
	router.PATCH("/games", game.Update)
	router.DELETE("/games", game.Delete)

	router.Use(middleware.AuthMiddleware())
	router.GET("/ws", ws.UpgradeWS)

	return router
}
