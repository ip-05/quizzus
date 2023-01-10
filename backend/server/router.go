package server

import (
	"github.com/ip-05/quizzus/controllers/web"
	"github.com/ip-05/quizzus/controllers/ws"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/middleware"
)

func NewRouter() *gin.Engine {
	cfg := config.GetConfig()

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	cors := cors.New(cors.Config{
		AllowOrigins:     []string{cfg.Frontend.Base},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	})

	auth := new(web.AuthController)
	game := new(web.GameController)
	ws := new(ws.CoreController)

	authGroup := router.Group("auth")
	authGroup.Use(cors)

	authGroup.GET("/google", auth.GoogleLogin)
	authGroup.GET("/google/callback", auth.GoogleCallback)

	authGroup.Use(middleware.AuthMiddleware())
	authGroup.GET("/me", auth.Me)

	//router.GET("/game/:id", game.GetById)

	gamesGroup := router.Group("games")
	gamesGroup.Use(middleware.AuthMiddleware())
	gamesGroup.GET("/", game.Get)
	gamesGroup.POST("/", game.CreateGame)
	gamesGroup.PATCH("/", game.Update)
	gamesGroup.DELETE("/", game.Delete)

	wsGroup := router.Group("ws")
	wsGroup.Use(cors)

	wsGroup.Use(middleware.WSMiddleware())
	wsGroup.GET("/", ws.HandleWS)

	return router
}
