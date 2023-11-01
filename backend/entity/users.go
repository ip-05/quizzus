package entity

import (
	"errors"
	"fmt"
)

type User struct {
	ID         uint   `json:"id" gorm:"primary_key"`
	GoogleID   string `json:"-"`
	DiscordID  string `json:"-"`
	TelegramID string `json:"-"`
	Picture    string `json:"picture"`
	Name       string `json:"name"`
}

type CreateUser struct {
	GoogleID   string `json:"google_id"`
	DiscordID  string `json:"discord_id"`
	TelegramID string `json:"telegram_id"`
	Picture    string `json:"picture"`
	Name       string `json:"name"`
}

type UpdateUser struct {
	Picture string `json:"picture"`
	Name    string `json:"name"`
}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
}

type DiscordUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type TelegramUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	PhotoUrl string `json:"photo_url"`
}

func NewUser(body *CreateUser) (*User, error) {
	user := &User{
		GoogleID:   body.GoogleID,
		DiscordID:  body.DiscordID,
		TelegramID: body.TelegramID,
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
		GoogleID: body.ID,
		Picture:  body.Picture,
		Name:     body.GivenName,
	}

	return user, nil
}

func NewDiscordUser(body DiscordUser) (*CreateUser, error) {
	user := &CreateUser{
		DiscordID: body.ID,
		Picture:   fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", body.ID, body.Avatar),
		Name:      body.Username,
	}

	return user, nil
}

func NewTelegramUser(body TelegramUser) (*CreateUser, error) {
	user := &CreateUser{
		TelegramID: body.ID,
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
