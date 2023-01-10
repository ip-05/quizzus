package server

import (
	"github.com/ip-05/quizzus/controllers/web"
	"github.com/ip-05/quizzus/controllers/ws"

	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/middleware"
)

func NewRouter() *gin.Engine {
	// cfg := config.GetConfig()

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	auth := new(web.AuthController)
	game := new(web.GameController)
	ws := new(ws.CoreController)

	authGroup := router.Group("auth")

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

	wsGroup.Use(middleware.WSMiddleware())
	wsGroup.GET("/", ws.HandleWS)

	return router
}
