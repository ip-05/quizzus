package web

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ip-05/quizzus/app/auth"
	"github.com/ip-05/quizzus/app/user"
	"github.com/ip-05/quizzus/entity"
	"github.com/ip-05/quizzus/repo"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

func newTestConfig() *config.Config {
	return &config.Config{
		Secrets: &config.SecretConfig{
			Jwt: "secret",
		},
		Server: &config.ServerConfig{
			Domain: "localhost",
			Secure: false,
		},
	}
}

type httpClientMock struct {
	mock.Mock
}

func (h httpClientMock) Get(url string) (resp *http.Response, err error) {
	args := h.Called()

	stringReader := strings.NewReader(args.Get(1).(string))
	stringReadCloser := io.NopCloser(stringReader)

	return &http.Response{
		StatusCode: args.Get(0).(int),
		Body:       stringReadCloser,
	}, err
}

type oAuth2Mock struct {
	mock.Mock
}

func (o oAuth2Mock) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	u := url.URL{
		Scheme: "http",
		Host:   "localhost",
		Path:   "login/oauth/authorize",
	}

	v := url.Values{}
	v.Set("state", state)

	u.RawQuery = v.Encode()
	return u.String()
}

func (o oAuth2Mock) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken: "AccessToken",
		Expiry:      time.Now().Add(1 * time.Hour),
	}, nil
}

type AuthControllerSuite struct {
	suite.Suite
	ctx        *gin.Context
	engine     *gin.Engine
	w          *httptest.ResponseRecorder
	controller *AuthController
	httpMock   httpClientMock
}

func (s *AuthControllerSuite) SetupTest() {
	oAuthMock := oAuth2Mock{}
	s.httpMock = httpClientMock{}

	db, _, err := sqlmock.New()
	assert.Nil(s.T(), err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	database, err := gorm.Open(dialector)
	userRepo := repo.NewUserStore(database)

	// Business logic layer
	userService := user.NewUserService(userRepo)
	authService := auth.NewAuthService(newTestConfig(), oAuthMock, userService, &s.httpMock)

	s.controller = NewAuthController(newTestConfig(), oAuthMock, authService, userService)

	gin.SetMode(gin.TestMode)

	s.w = httptest.NewRecorder()
	s.ctx, s.engine = gin.CreateTestContext(s.w)

	s.engine.GET("/auth/google", s.controller.GoogleLogin)
	s.engine.GET("/auth/google/callback", s.controller.GoogleCallback)
}

func (s *AuthControllerSuite) TestLogin_RedirectUrl() {
	// When
	s.ctx.Request, _ = http.NewRequest(http.MethodGet, "/auth/google/callback", nil)
	s.engine.ServeHTTP(s.w, s.ctx.Request)

	// Then
	m := map[string]string{}

	body, err := io.ReadAll(s.w.Body)
	assert.Nil(s.T(), err)

	err = json.Unmarshal(body, &m)
	assert.Nil(s.T(), err)

	redirect := m["redirectUrl"]
	assert.NotNil(s.T(), redirect)
}

func (s *AuthControllerSuite) TestLogin_MissingCookie() {
	// When
	s.ctx.Request, _ = http.NewRequest(http.MethodGet, "/auth/google/callback", nil)
	s.ctx.Request.Form, _ = url.ParseQuery("state=secondState")

	s.engine.ServeHTTP(s.w, s.ctx.Request)

	// Then
	body, err := io.ReadAll(s.w.Body)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), http.StatusBadRequest, s.w.Code)
	assert.Contains(s.T(), string(body), "Invalid cookie")
}

func (s *AuthControllerSuite) TestLogin_MismatchedState() {
	// When
	s.ctx.Request, _ = http.NewRequest(http.MethodGet, "/auth/google/callback", nil)
	s.ctx.Request.AddCookie(&http.Cookie{
		Name:     "oauthstate",
		Value:    "firstState",
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   60,
		Secure:   false,
		HttpOnly: true,
	})
	s.ctx.Request.Form, _ = url.ParseQuery("state=secondState")

	s.engine.ServeHTTP(s.w, s.ctx.Request)

	// Then
	body, err := io.ReadAll(s.w.Body)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), http.StatusBadRequest, s.w.Code)
	assert.Contains(s.T(), string(body), "Error while verifying auth token")
}

func (s *AuthControllerSuite) TestLogin_MissingCode() {
	// When
	s.ctx.Request, _ = http.NewRequest(http.MethodGet, "/auth/google/callback", nil)
	s.ctx.Request.AddCookie(&http.Cookie{
		Name:     "oauthstate",
		Value:    "firstState",
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   60,
		Secure:   false,
		HttpOnly: true,
	})
	s.ctx.Request.Form, _ = url.ParseQuery("state=firstState")

	s.engine.ServeHTTP(s.w, s.ctx.Request)

	// Then
	body, err := io.ReadAll(s.w.Body)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), http.StatusBadRequest, s.w.Code)
	assert.Contains(s.T(), string(body), "Missing code")
}

func (s *AuthControllerSuite) TestLogin_VerifyError() {
	// When
	s.w = httptest.NewRecorder()

	s.ctx.Request, _ = http.NewRequest(http.MethodGet, "/auth/google/callback", nil)
	s.ctx.Request.AddCookie(&http.Cookie{
		Name:     "oauthstate",
		Value:    "firstState",
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   60,
		Secure:   false,
		HttpOnly: true,
	})
	s.ctx.Request.Form, _ = url.ParseQuery("state=firstState&code=code")

	s.httpMock.On("Get").Return(http.StatusUnauthorized, "").Times(1)

	s.engine.ServeHTTP(s.w, s.ctx.Request)

	// Then
	body, err := io.ReadAll(s.w.Body)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), http.StatusBadRequest, s.w.Code)
	assert.Contains(s.T(), string(body), "Error while verifying auth token")
}

func (s *AuthControllerSuite) TestLogin_ReturnJWT() {
	// When
	userInfo := entity.GoogleUser{
		Id:            "123123",
		Email:         "john@doe.com",
		VerifiedEmail: true,
		Picture:       "https://john.doe.com/picture.png",
		GivenName:     "John",
	}
	userString, _ := json.Marshal(&userInfo)

	s.httpMock.On("Get").Return(http.StatusOK, string(userString)).Times(2)

	s.ctx.Request, _ = http.NewRequest(http.MethodGet, "/auth/google/callback", nil)
	s.ctx.Request.AddCookie(&http.Cookie{
		Name:     "oauthstate",
		Value:    "firstState",
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   60,
		Secure:   false,
		HttpOnly: true,
	})
	s.ctx.Request.Form, _ = url.ParseQuery("state=firstState&code=code")

	s.engine.ServeHTTP(s.w, s.ctx.Request)

	// Then
	body, err := io.ReadAll(s.w.Body)
	assert.Nil(s.T(), err)

	assert.Contains(s.T(), string(body), "Successfully authenticated user")
}

func (s *AuthControllerSuite) TestLogin_SetCookie() {
	// When
	userInfo := entity.GoogleUser{
		Id:            "123123",
		Email:         "john@doe.com",
		VerifiedEmail: true,
		Picture:       "https://john.doe.com/picture.png",
		GivenName:     "John",
	}
	userString, _ := json.Marshal(&userInfo)

	s.httpMock.On("Get").Return(http.StatusOK, string(userString)).Times(2)

	s.ctx.Request, _ = http.NewRequest(http.MethodGet, "/auth/google/callback", nil)
	s.ctx.Request.AddCookie(&http.Cookie{
		Name:     "oauthstate",
		Value:    "firstState",
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   60,
		Secure:   false,
		HttpOnly: true,
	})
	s.ctx.Request.Form, _ = url.ParseQuery("state=firstState&code=code")

	s.engine.ServeHTTP(s.w, s.ctx.Request)

	// Then
	setCookie := s.w.Header().Get("Set-Cookie")
	assert.NotEmpty(s.T(), setCookie)
}

// func (s *AuthControllerSuite) TestMe() {
// 	// Given
// 	authedUser := middleware.AuthedUser{
// 		Id:   uint(123),
// 		Name: "John",
// 	}

// 	s.ctx.Set("authedUser", authedUser)

// 	// When
// 	s.controller.Me(s.ctx)

// 	// Then
// 	r, err := io.ReadAll(s.w.Body)
// 	assert.Nil(s.T(), err)

// 	json, err := json.Marshal(authedUser)
// 	assert.Nil(s.T(), err)

// 	assert.Equal(s.T(), json, r)
// }

func TestGoogleLogin(t *testing.T) {
	suite.Run(t, new(AuthControllerSuite))
}
