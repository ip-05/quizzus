package repo

import (
	"fmt"

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
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&entity.Option{})
	if err != nil {
		return nil
	}

	err = database.AutoMigrate(&entity.Question{})
	if err != nil {
		return nil
	}

	err = database.AutoMigrate(&entity.Game{})
	if err != nil {
		return nil
	}
	return database
}
