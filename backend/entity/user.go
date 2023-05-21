package entity

type UserInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
}

func NewUser() (*UserInfo, error) {
	user := &UserInfo{}
	return user, nil
}

// TODO: add users to DB
