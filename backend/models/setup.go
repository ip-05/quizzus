package models

import (
	"fmt"

	"github.com/ip-05/quizzus/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	cfg := config.GetConfig()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.DbName)

	database, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&Option{})
	if err != nil {
		return
	}

	err = database.AutoMigrate(&Question{})
	if err != nil {
		return
	}

	err = database.AutoMigrate(&Game{})
	if err != nil {
		return
	}

	DB = database
}
