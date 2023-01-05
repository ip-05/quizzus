package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ip-05/quizzus/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	cfg := config.GetConfig()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.DbName)

	database, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		Logger: newLogger,
	})
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
