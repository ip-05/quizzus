package server

import (
	"github.com/gin-contrib/cors"
	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/controllers/web"
	"github.com/ip-05/quizzus/controllers/ws"

	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/middleware"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	cfg := config.GetConfig()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{cfg.Frontend.Base}
	config.AllowCredentials = true
	config.AddAllowHeaders("Authorization")
	router.Use(cors.New(config))

	auth := new(web.AuthController)
	game := new(web.GameController)
	ws := new(ws.CoreController)

	authGroup := router.Group("auth")

	authGroup.GET("/google", auth.GoogleLogin)
	authGroup.GET("/google/callback", auth.GoogleCallback)

	authGroup.Use(middleware.AuthMiddleware())
	authGroup.GET("/me", auth.Me)

	gamesGroup := router.Group("games")
	gamesGroup.Use(middleware.AuthMiddleware())
	gamesGroup.GET("", game.Get)
	gamesGroup.POST("", game.CreateGame)
	gamesGroup.PATCH("", game.Update)
	gamesGroup.DELETE("", game.Delete)

	wsGroup := router.Group("ws")
	wsGroup.Use(middleware.WSMiddleware())
	wsGroup.GET("", ws.HandleWS)

	return router
}
