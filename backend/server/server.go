package server

import (
	"fmt"

	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/models"
)

func Init() {
	r := NewRouter()

	cfg := config.GetConfig()
	models.ConnectDatabase()

	r.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
}
