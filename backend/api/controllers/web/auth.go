package web

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/ip-05/quizzus/api/middleware"
	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/entity"
	"golang.org/x/oauth2"

	"github.com/gin-gonic/gin"
)

type IUserService interface {
	CreateUser(body *entity.CreateUser) (*entity.User, error)
	UpdateUser()
	DeleteUser()
	GetUser(id uint) *entity.User
	GetUserByProvider(id string, provider string) *entity.User
}

type GoogleAuth interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
}

type IAuthService interface {
	AuthenticateGoogle(code string) (string, error)
}

type AuthController struct {
	Config       *config.Config
	GoogleConfig GoogleAuth
	Auth         IAuthService
	User         IUserService
}

func NewAuthController(cfg *config.Config, gcfg GoogleAuth, auth IAuthService, user IUserService) *AuthController {
	return &AuthController{Auth: auth, GoogleConfig: gcfg, Config: cfg, User: user}
}

func (a AuthController) GoogleLogin(c *gin.Context) {
	var expiration = int(20 * time.Minute)
	b := make([]byte, 16)
	rand.Read(b)
	oauthState := base64.URLEncoding.EncodeToString(b)
	c.SetCookie("oauthstate", oauthState, expiration, "*", a.Config.Server.Domain, a.Config.Server.Secure, false)

	url := a.GoogleConfig.AuthCodeURL(oauthState)

	c.JSON(http.StatusOK, gin.H{"redirectUrl": url})
}

func (a AuthController) GoogleCallback(c *gin.Context) {
	oauthState, err := c.Cookie("oauthstate")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cookie"})
		return
	}

	if c.Request.FormValue("state") != oauthState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while verifying auth token"})
		return
	}

	code := c.Request.FormValue("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code"})
		return
	}

	token, err := a.Auth.AuthenticateGoogle(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", token, 7*24*60*60, "/", a.Config.Server.Domain, a.Config.Server.Secure, false)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user", "token": token})
}

func (a AuthController) Me(c *gin.Context) {
	authedUser, _ := c.Get("authedUser")
	user := authedUser.(middleware.AuthedUser)

	dbUser := a.User.GetUser(user.Id)

	c.JSON(http.StatusOK, dbUser)
}
