package repo

import (
	"fmt"
	"log"

	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(cfg *config.Config) *gorm.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.DbName)

	database, err := gorm.Open(postgres.Open(psqlInfo))
	if err != nil {
		panic("FAIL_CONNECT_DB")
	}

	err = database.AutoMigrate(
		&entity.Option{},
		&entity.Question{},
		&entity.Game{},
		&entity.User{},
		&entity.FavoriteGame{},
		&entity.GameSession{},
		&entity.GameSession{},
	)
	if err != nil {
		log.Print("FAIL_MIGRATIONS")
	}
	return database
}
