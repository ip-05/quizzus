package user

import (
	"github.com/ip-05/quizzus/entity"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r Repository) GetUserById(ID uint) *entity.User {
	var user *entity.User
	r.DB.Where("id = ?", ID).First(&user)
	if user.ID == 0 {
		return nil
	}
	return user
}

func (r Repository) GetUserByGoogleId(ID string) *entity.User {
	var user *entity.User
	r.DB.Where("google_id = ?", ID).First(&user)
	if user.ID == 0 {
		return nil
	}
	return user
}

func (r Repository) GetUserByDiscordId(ID string) *entity.User {
	var user *entity.User
	r.DB.Where("discord_id = ?", ID).First(&user)
	if user.ID == 0 {
		return nil
	}
	return user
}

func (r Repository) GetUserByTelegramId(ID string) *entity.User {
	var user *entity.User
	r.DB.Where("telegram_id = ?", ID).First(&user)
	if user.ID == 0 {
		return nil
	}
	return user
}

func (r Repository) CreateUser(e *entity.User) *entity.User {
	r.DB.Create(&e)
	return e
}

func (r Repository) UpdateUser(e *entity.User) *entity.User {
	r.DB.Where("id = ?", e.ID).Updates(&e)
	return e
}

func (r Repository) DeleteUser(e *entity.User) {
	r.DB.Unscoped().Delete(&e)
}
