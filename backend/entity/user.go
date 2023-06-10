package entity

type User struct {
	Id          string   `json:"id"`
	Picture     string   `json:"picture"`
	GivenName   string   `json:"given_name"`
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

func NewUser() (*User, error) {
	user := &User{}
	return user, nil
}

// TODO: add users to DB
