package entity

import (
	"errors"
	"fmt"
)

type User struct {
	Id         uint   `json:"id" gorm:"primary_key"`
	GoogleId   string `json:"google_id,omitempty"`
	DiscordId  string `json:"discord_id,omitempty"`
	TelegramId string `json:"telegram_id,omitempty"`
	Picture    string `json:"picture"`
	Name       string `json:"name"`
}

type CreateUser struct {
	GoogleId   string `json:"google_id"`
	DiscordId  string `json:"discord_id"`
	TelegramId string `json:"telegram_id"`
	Picture    string `json:"picture"`
	Name       string `json:"name"`
}

type UpdateUser struct {
	Picture string `json:"picture"`
	Name    string `json:"name"`
}

type GoogleUser struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
}

type DiscordUser struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type TelegramUser struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	PhotoUrl string `json:"photo_url"`
}

func NewUser(body *CreateUser) (*User, error) {
	user := &User{
		GoogleId:   body.GoogleId,
		DiscordId:  body.DiscordId,
		TelegramId: body.TelegramId,
		Picture:    body.Picture,
		Name:       body.Name,
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func NewGoogleUser(body GoogleUser) (*CreateUser, error) {
	user := &CreateUser{
		GoogleId: body.Id,
		Picture:  body.Picture,
		Name:     body.GivenName,
	}

	return user, nil
}

func NewDiscordUser(body DiscordUser) (*CreateUser, error) {
	user := &CreateUser{
		DiscordId: body.Id,
		Picture:   fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", body.Id, body.Avatar),
		Name:      body.Username,
	}

	return user, nil
}

func NewTelegramUser(body TelegramUser) (*CreateUser, error) {
	user := &CreateUser{
		TelegramId: body.Id,
		Picture:    body.PhotoUrl,
		Name:       body.Username,
	}

	return user, nil
}

func (u *User) Validate() error {
	if len(u.Name) < 2 || len(u.Name) > 32 {
		return errors.New("name must be between 2 and 32 characters long")
	}

	return nil
}
