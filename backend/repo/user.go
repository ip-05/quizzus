package repo

import (
	"github.com/ip-05/quizzus/entity"
	"gorm.io/gorm"
)

type UserStore struct {
	DB *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{
		DB: db,
	}
}

func (db *UserStore) Get(id int) *entity.User {
	return nil
}

func (db *UserStore) Create() *entity.User {
	return nil
}

func (db *UserStore) Update() *entity.User {
	return nil
}

func (db *UserStore) Delete() *entity.User {
	return nil
}
