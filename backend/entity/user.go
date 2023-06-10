package entity

type User struct {
	Id          string   `json:"id"`
	GoogleId    string   `json:"google_id"`
	DiscordId   string   `json:"discord_id"`
	Picture     string   `json:"picture"`
	Name        string   `json:"name"`
	SavedGames  []string `json:"saved_games"`
	PlayedGames []string `json:"played_games"`
}

type GoogleUser struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
}

type DiscordUser struct {
}

type TelegramUser struct {
}

func NewGoogleUser(body GoogleUser) (*User, error) {
	user := &User{
		GoogleId:  body.Id,
		DiscordId: "",
		Picture:   body.Picture,
		Name:      body.GivenName,
	}
	return user, nil
}

// TODO: add users to DB
