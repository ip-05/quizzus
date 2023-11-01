package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	authController "github.com/ip-05/quizzus/api/controllers/auth"
	gameController "github.com/ip-05/quizzus/api/controllers/game"
	sessionController "github.com/ip-05/quizzus/api/controllers/session"
	userController "github.com/ip-05/quizzus/api/controllers/user"
	ws "github.com/ip-05/quizzus/api/controllers/ws"

	"github.com/ip-05/quizzus/api/middleware"
	"github.com/ip-05/quizzus/config"
	"golang.org/x/oauth2"
)

func InitWeb(
	cfg *config.Config,
	gcfg *oauth2.Config,
	gameSvc gameController.Service,
	authSvc authController.AuthService,
	userSvc userController.Service,
	sessionSvc sessionController.Service,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{cfg.Frontend.Base}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization")
	router.Use(cors.New(corsConfig))

	userController := userController.NewController(userSvc)
	sessionsController := sessionController.NewController(sessionSvc)
	authController := authController.NewController(cfg, gcfg, authSvc, userSvc)
	gameController := gameController.NewController(gameSvc)

	ws := ws.NewCoreController(gameSvc, userSvc, sessionSvc)

	userGroup := router.Group("users")
	{
		userGroup.Use(middleware.AuthMiddleware(cfg))
		userGroup.GET("/:id", userController.Get)
		userGroup.PATCH("/me", userController.Update)
		userGroup.DELETE("/me", userController.Delete)
	}

	authGroup := router.Group("auth")
	{
		googleGroup := authGroup.Group("google")
		{
			googleGroup.GET("", authController.GoogleLogin)
			googleGroup.GET("/callback", authController.GoogleCallback)
		}
	}

	gamesGroup := router.Group("games")
	{
		gamesGroup.Use(middleware.AuthMiddleware(cfg))
		gamesGroup.GET("/:id", gameController.Get)
		gamesGroup.GET("", gameController.GetMany)
		gamesGroup.POST("/:id/favorite", gameController.Favorite)
		gamesGroup.POST("", gameController.CreateGame)
		gamesGroup.PATCH("", gameController.Update)
		gamesGroup.DELETE("", gameController.Delete)
	}

	sessionsGroup := router.Group("sessions")
	{
		sessionsGroup.Use(middleware.AuthMiddleware(cfg))
		sessionsGroup.GET("", sessionsController.GetSessions)
		sessionsGroup.GET("/:id", sessionsController.GetSession)
	}

	wsGroup := router.Group("ws")
	{
		wsGroup.Use(middleware.WSMiddleware(cfg))
		wsGroup.GET("", ws.HandleWS)
	}

	return router
}
