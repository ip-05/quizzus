package main

import (
	"fmt"
	"net/http"

	"github.com/ip-05/quizzus/app/auth"
	"github.com/ip-05/quizzus/app/session"
	"github.com/ip-05/quizzus/app/user"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/ip-05/quizzus/api"
	"github.com/ip-05/quizzus/app/game"
	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/repo"
)

func main() {
	cfg := config.GetConfig()

	// Google Config
	gcfg := &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/auth/google", cfg.Frontend.Base),
		ClientID:     cfg.Google.ClientId,
		ClientSecret: cfg.Google.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// Repository layer
	db := repo.New(cfg)
	gameRepo := repo.NewGameStore(db)
	userRepo := repo.NewUserStore(db)
	sessionRepo := repo.NewSessionStore(db)

	// Business logic layer
	gameService := game.NewService(gameRepo)
	userService := user.NewService(userRepo)
	authService := auth.NewService(cfg, gcfg, userService, &http.Client{})
	sessionService := session.NewSessionService(sessionRepo)

	// Presentation layer
	r := api.InitWeb(cfg, gcfg, gameService, authService, userService, sessionService)

	r.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
}
