package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/entity"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"time"
)

type IUserService interface {
	CreateUser(body *entity.CreateUser) (*entity.User, error)
	UpdateUser()
	DeleteUser()
	GetUser()
	GetUserByProvider(id string, provider string) *entity.User
}

type GoogleAuth interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
}

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
}

type AuthService struct {
	Config       *config.Config
	GoogleConfig GoogleAuth
	Http         HttpClient

	User IUserService
}

func NewAuthService(cfg *config.Config, gcfg GoogleAuth, user IUserService, http HttpClient) *AuthService {
	return &AuthService{
		Config:       cfg,
		GoogleConfig: gcfg,
		User:         user,
		Http:         http,
	}
}

func (u *AuthService) AuthenticateGoogle(code string) (string, error) {
	token, err := u.GoogleConfig.Exchange(context.Background(), code)
	if err != nil {
		return "", errors.New("error exchanging code for token")
	}

	response, err := u.Http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return "", errors.New("error fetching user data from google")
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", errors.New("error fetching user data from google")
	}

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("error while decoding token contents")
	}

	var userInfo entity.GoogleUser

	err = json.Unmarshal(contents, &userInfo)
	if err != nil {
		return "", errors.New("error while parsing user info")
	}

	existingUser := u.User.GetUserByProvider(userInfo.Id, "google")
	if existingUser != nil {
		// TODO: Move token generation to utility class
		secretKey := []byte(u.Config.Secrets.Jwt)
		tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":   existingUser.Id,
			"name": existingUser.Name,
			"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(),
		})

		tokenString, err := tokenJWT.SignedString(secretKey)
		if err != nil {
			return "", errors.New("error while signing JWT token")
		}

		return tokenString, nil
	}

	newUser, err := entity.NewGoogleUser(userInfo)
	if err != nil {
		return "", err
	}

	user, err := u.User.CreateUser(newUser)
	if err != nil {
		return "", err
	}

	secretKey := []byte(u.Config.Secrets.Jwt)
	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.Id,
		"name": user.Name,
		"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	tokenString, err := tokenJWT.SignedString(secretKey)
	if err != nil {
		return "", errors.New("error while signing JWT token")
	}

	return tokenString, nil
}
