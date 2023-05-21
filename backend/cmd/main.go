package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/api/controllers/web"
	"github.com/ip-05/quizzus/api/controllers/ws"
	"github.com/ip-05/quizzus/api/middleware"
	"github.com/ip-05/quizzus/app/game"
	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/repo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var cfg = config.GetConfig()

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{cfg.Frontend.Base}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization")
	router.Use(cors.New(corsConfig))

	auth := web.NewAuthController(cfg, &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/auth/google", cfg.Frontend.Base),
		ClientID:     cfg.Google.ClientId,
		ClientSecret: cfg.Google.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}, &http.Client{})

	// Repository layer
	db := repo.New(cfg)

	gameRepo := repo.NewGameStore(cfg, db)
	questionRepo := repo.NewQuestionStore(cfg, db)
	optionRepo := repo.NewOptionStore(cfg, db)

	// Business logic layer
	gameService := game.NewGameService(gameRepo, questionRepo, optionRepo)

	// TODO: Presentation layer
	// apiWeb := InitWeb(gameService) ?
	// apiWs := InitWs(gameService) ?

	game := web.NewGameController(db)
	ws := ws.NewCoreController(db)

	authGroup := router.Group("auth")
	{
		authGroup.GET("/google", auth.GoogleLogin)
		authGroup.GET("/google/callback", auth.GoogleCallback)

		authGroup.Use(middleware.AuthMiddleware(cfg))
		authGroup.GET("/me", auth.Me)
	}

	gamesGroup := router.Group("games")
	{
		gamesGroup.Use(middleware.AuthMiddleware(cfg))
		gamesGroup.GET("", game.Get)
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

func main() {
	r := NewRouter()

	r.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
}
