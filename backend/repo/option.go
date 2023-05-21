package repo

import (
	"errors"

	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/entity"
	"gorm.io/gorm"
)

type OptionStore struct {
	DB *gorm.DB
}

func NewOptionStore(cfg *config.Config, db *gorm.DB) *OptionStore {
	return &OptionStore{
		DB: db,
	}
}

func (db *OptionStore) Get(id int) (*entity.Option, error) {
	var option entity.Option

	db.DB.Where("id = ?", id).First(&option)

	if option.Id == 0 {
		return nil, errors.New("option not found")
	}

	return &option, nil
}

func (db *OptionStore) GetByGameId(gameId int) (*[]entity.Option, error) {
	var options []entity.Option

	db.DB.Where("game_id = ?", gameId).Find(&options)

	if len(options) == 0 {
		return nil, errors.New("options not found")
	}

	return &options, nil
}

func (db *OptionStore) Create(e *entity.Option) (*entity.Option, error) {
	db.DB.Create(&e)
	return e, nil
}

func (db *OptionStore) Update(e *entity.Option) (*entity.Option, error) {
	db.DB.Where("id = ?", e.Id).Updates(&e)
	return e, nil
}

func (db *OptionStore) Delete(e *entity.Option) error {
	db.DB.Where("id = ?", e.Id).Delete(&e)
	return nil
}
