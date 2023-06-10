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
	var user *entity.User
	db.DB.Where("id = ?", id).First(&user)
	return nil
}

func (db *UserStore) Create(e *entity.User) *entity.User {
	db.DB.Create(&e)
	return e
}

func (db *UserStore) Update(id int, e *entity.User) *entity.User {
	db.DB.Where("id = ?", id).Updates(&e)
	return e
}

func (db *UserStore) Delete(e *entity.User) {
	db.DB.Unscoped().Delete(&e)
}
