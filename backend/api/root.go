package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/api/controllers/web"
	"github.com/ip-05/quizzus/api/controllers/ws"
	"github.com/ip-05/quizzus/api/middleware"
	"github.com/ip-05/quizzus/app/auth"
	"github.com/ip-05/quizzus/config"
	"golang.org/x/oauth2"
)

func InitWeb(cfg *config.Config, gcfg *oauth2.Config, gameService web.IGameService, authService web.IAuthService, userService auth.IUserService) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{cfg.Frontend.Base}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization")
	router.Use(cors.New(corsConfig))

	userController := web.NewUserController(userService)
	authController := web.NewAuthController(cfg, gcfg, authService, userService)

	game := web.NewGameController(gameService)
	ws := ws.NewCoreController(gameService, userService)

	userGroup := router.Group("users")
	{
		userGroup.Use(middleware.AuthMiddleware(cfg))
		userGroup.GET("/:id", userController.Get)
		userGroup.PATCH("/me", userController.Update)
		userGroup.DELETE("/me", userController.Delete)
	}

	authGroup := router.Group("auth")
	{
		authGroup.GET("/google", authController.GoogleLogin)
		authGroup.GET("/google/callback", authController.GoogleCallback)
	}

	gamesGroup := router.Group("games")
	{
		gamesGroup.Use(middleware.AuthMiddleware(cfg))
		gamesGroup.GET("/:id", game.Get)
		gamesGroup.GET("", game.GetMany)
		gamesGroup.POST("/:id/favorite", game.Favorite)
		gamesGroup.POST("", game.CreateGame)
		gamesGroup.PATCH("", game.Update)
		gamesGroup.DELETE("", game.Delete)
	}

	wsGroup := router.Group("ws")
	{
		wsGroup.Use(middleware.WSMiddleware(cfg))
		wsGroup.GET("", ws.HandleWS)
	}

	return router
}
