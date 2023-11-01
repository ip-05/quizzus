package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/entity"
	"golang.org/x/oauth2"

	"github.com/gin-gonic/gin"
)

type GoogleAuth interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
}

type AuthService interface {
	AuthenticateGoogle(code string) (string, error)
}

type UserService interface {
	CreateUser(body *entity.CreateUser) (*entity.User, error)
	UpdateUser(id uint, body entity.UpdateUser) (*entity.User, error)
	DeleteUser(id uint)
	GetUser(id uint) *entity.User
	GetUserByProvider(id string, provider string) *entity.User
}

type Controller struct {
	Config       *config.Config
	GoogleConfig GoogleAuth
	Auth         AuthService
	User         UserService
}

func NewController(cfg *config.Config, gcfg GoogleAuth, auth AuthService, user UserService) *Controller {
	return &Controller{Auth: auth, GoogleConfig: gcfg, Config: cfg, User: user}
}

func (c Controller) GoogleLogin(ctx *gin.Context) {
	var expiration = int(20 * time.Minute)
	b := make([]byte, 16)
	rand.Read(b)
	oauthState := base64.URLEncoding.EncodeToString(b)
	ctx.SetCookie("oauthstate", oauthState, expiration, "*", c.Config.Server.Domain, c.Config.Server.Secure, false)

	url := c.GoogleConfig.AuthCodeURL(oauthState)

	ctx.JSON(http.StatusOK, gin.H{"redirectUrl": url})
}

func (c Controller) GoogleCallback(ctx *gin.Context) {
	oauthState, err := ctx.Cookie("oauthstate")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cookie"})
		return
	}

	if ctx.Request.FormValue("state") != oauthState {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error while verifying auth token"})
		return
	}

	code := ctx.Request.FormValue("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing code"})
		return
	}

	token, err := c.Auth.AuthenticateGoogle(code)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.SetCookie("token", token, 7*24*60*60, "/", c.Config.Server.Domain, c.Config.Server.Secure, false)
	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user", "token": token})
}
