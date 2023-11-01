package game

import (
	"fmt"
	"os"
	"testing"

	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func SetupIntegration(t *testing.T) (*gorm.DB, func() error) {
	cfg := config.Init("config", "../config")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.DbName)

	db, err := gorm.Open(postgres.Open(psqlInfo))
	if err != nil {
		panic("Failed to connect to database!")
	}

	err = db.AutoMigrate(
		&entity.Option{},
		&entity.Question{},
		&entity.Game{},
	)
	if err != nil {
		return nil, nil
	}

	return db, func() error {
		res := db.Exec("SET session_replication_role = 'replica';")
		if res.Error != nil {
			return res.Error
		}

		res = db.Exec("TRUNCATE TABLE games;")
		if res.Error != nil {
			return res.Error
		}

		res = db.Exec("TRUNCATE TABLE questions;")
		if res.Error != nil {
			return res.Error
		}

		res = db.Exec("TRUNCATE TABLE options;")
		if res.Error != nil {
			return res.Error
		}

		res = db.Exec("SET session_replication_role = 'origin';")
		if res.Error != nil {
			return res.Error
		}

		return nil
	}
}
