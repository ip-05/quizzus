package server

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/controllers/web"
	"github.com/ip-05/quizzus/controllers/ws"
	"github.com/ip-05/quizzus/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/middleware"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization")
	router.Use(cors.New(corsConfig))

	cfg := config.GetConfig()

	auth := web.NewAuthController(cfg, &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/auth/google", cfg.Frontend.Base),
		ClientID:     cfg.Google.ClientId,
		ClientSecret: cfg.Google.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	})
	game := web.NewGameController(models.DB)
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
