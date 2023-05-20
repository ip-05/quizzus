package db

import (
	"fmt"

	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type store struct {
	DB *gorm.DB
}

func New(cfg *config.Config) *store {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.DbName)

	db, err := gorm.Open(postgres.Open(psqlInfo))
	if err != nil {
		panic("Failed to connect to database!")
	}

	err = db.AutoMigrate(&entity.Option{})
	if err != nil {
		return nil
	}

	err = db.AutoMigrate(&entity.Question{})
	if err != nil {
		return nil
	}

	err = db.AutoMigrate(&entity.Game{})
	if err != nil {
		return nil
	}

	return &store{
		DB: db,
	}
}
