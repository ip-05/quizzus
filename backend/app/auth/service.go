package auth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/entity"
	"github.com/ip-05/quizzus/utils"
	"golang.org/x/oauth2"
)

type Repository interface {
	CreateUser(body *entity.CreateUser) (*entity.User, error)
	UpdateUser(ID uint, body entity.UpdateUser) (*entity.User, error)
	DeleteUser(ID uint)
	GetUserById(ID uint) *entity.User
	GetUserByProvider(ID string, provider string) *entity.User
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

	Repo Repository
}

func NewService(cfg *config.Config, gcfg GoogleAuth, userRepo Repository, http HttpClient) *AuthService {
	return &AuthService{
		Config:       cfg,
		GoogleConfig: gcfg,
		Repo:         userRepo,
		Http:         http,
	}
}

func (s *AuthService) AuthenticateGoogle(code string) (string, error) {
	token, err := s.GoogleConfig.Exchange(context.Background(), code)
	if err != nil {
		return "", errors.New("error exchanging code for token")
	}

	response, err := s.Http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
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

	existingUser := s.Repo.GetUserByProvider(userInfo.ID, "google")
	if existingUser != nil {
		tokenString, err := utils.GenerateToken(existingUser.ID, existingUser.Name, s.Config.Secrets.Jwt)
		if err != nil {
			return "", err
		}
		return tokenString, nil
	}

	newUser, err := entity.NewGoogleUser(userInfo)
	if err != nil {
		return "", err
	}

	user, err := s.Repo.CreateUser(newUser)
	if err != nil {
		return "", err
	}

	tokenString, err := utils.GenerateToken(user.ID, user.Name, s.Config.Secrets.Jwt)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
