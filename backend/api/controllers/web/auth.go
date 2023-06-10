package web

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/ip-05/quizzus/api/middleware"
	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/entity"
	"golang.org/x/oauth2"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type AuthController struct {
	Config       *config.Config
	GoogleConfig GoogleAuth
	Http         HttpClient
	//	User         IUserService
}

type GoogleAuth interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
}

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
}

type IUserService interface {
	CreateUser()
	UpdateUser()
	DeleteUser()
	GetUser()
}

func NewAuthController(cfg *config.Config, gcfg GoogleAuth, http HttpClient) *AuthController {
	return &AuthController{Config: cfg, GoogleConfig: gcfg, Http: http}
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

	token, err := a.GoogleConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while exchanging auth token"})
		return
	}

	response, err := a.Http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while verifying auth token"})
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while verifying auth token"})
		return
	}

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while decoding token contents"})
		return
	}

	var userInfo entity.GoogleUser

	err = json.Unmarshal(contents, &userInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while parsing user info"})
		return
	}

	secretKey := []byte(a.Config.Secrets.Jwt)
	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":             userInfo.Id,
		"name":           userInfo.GivenName,
		"email":          userInfo.Email,
		"profilePicture": userInfo.Picture,
		"exp":            time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	tokenString, err := tokenJWT.SignedString(secretKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while signing JWT token"})
		return
	}

	// a.User.CreateUser(googleUserInfo)
	c.SetCookie("token", tokenString, 7*24*60*60, "/", a.Config.Server.Domain, a.Config.Server.Secure, false)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user", "token": tokenString})
}

func (a AuthController) Me(c *gin.Context) {
	authedUser, _ := c.Get("authedUser")
	user := authedUser.(middleware.AuthedUser)
	c.JSON(http.StatusOK, user)
}
