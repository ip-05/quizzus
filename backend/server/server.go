package server

import (
	"fmt"
	"github.com/ip-05/quizzus/config"
)

func Init() {
	r := NewRouter()

	cfg := config.GetConfig()
	r.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
}
