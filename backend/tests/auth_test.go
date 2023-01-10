package tests

import (
	"context"
	"encoding/json"
	"github.com/ip-05/quizzus/middleware"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/controllers/web"
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
	args := h.Called(url)

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

func TestGoogleLogin(t *testing.T) {
	// Given
	mock := oAuth2Mock{}
	httpMock := httpClientMock{}

	authController := web.NewAuthController(newTestConfig(), &mock, &httpMock)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(w)

	engine.GET("/auth/me", authController.Me)
	engine.GET("/auth/google", authController.GoogleLogin)
	engine.GET("/auth/google/callback", authController.GoogleCallback)

	t.Run("should return redirect url", func(t *testing.T) {
		// When
		ctx.Request, _ = http.NewRequest(http.MethodGet, "/auth/google/callback", nil)
		engine.ServeHTTP(w, ctx.Request)

		// Then
		m := map[string]string{}

		body, err := io.ReadAll(w.Body)
		assert.Nil(t, err)

		err = json.Unmarshal(body, &m)
		assert.Nil(t, err)

		redirect := m["redirectUrl"]
		assert.NotNil(t, redirect)
	})

	t.Run("should error - missing cookie", func(t *testing.T) {
		// When
		w = httptest.NewRecorder()

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/auth/google/callback", nil)
		ctx.Request.Form, _ = url.ParseQuery("state=secondState")

		engine.ServeHTTP(w, ctx.Request)

		// Then
		body, err := io.ReadAll(w.Body)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, string(body), "Invalid cookie")
	})

	t.Run("should error - mismatched state", func(t *testing.T) {
		// When
		w = httptest.NewRecorder()

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/auth/google/callback", nil)
		ctx.Request.AddCookie(&http.Cookie{
			Name:     "oauthstate",
			Value:    "firstState",
			Path:     "/",
			Domain:   "localhost",
			MaxAge:   60,
			Secure:   false,
			HttpOnly: true,
		})
		ctx.Request.Form, _ = url.ParseQuery("state=secondState")

		engine.ServeHTTP(w, ctx.Request)

		// Then
		body, err := io.ReadAll(w.Body)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, string(body), "Error while verifying auth token")
	})

	t.Run("should error - missing code", func(t *testing.T) {
		// When
		w = httptest.NewRecorder()

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/auth/google/callback", nil)
		ctx.Request.AddCookie(&http.Cookie{
			Name:     "oauthstate",
			Value:    "firstState",
			Path:     "/",
			Domain:   "localhost",
			MaxAge:   60,
			Secure:   false,
			HttpOnly: true,
		})
		ctx.Request.Form, _ = url.ParseQuery("state=firstState")

		engine.ServeHTTP(w, ctx.Request)

		// Then
		body, err := io.ReadAll(w.Body)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, string(body), "Missing code")
	})

	t.Run("should return error on verifying auth token", func(t *testing.T) {
		// When
		w = httptest.NewRecorder()

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/auth/google/callback", nil)
		ctx.Request.AddCookie(&http.Cookie{
			Name:     "oauthstate",
			Value:    "firstState",
			Path:     "/",
			Domain:   "localhost",
			MaxAge:   60,
			Secure:   false,
			HttpOnly: true,
		})
		ctx.Request.Form, _ = url.ParseQuery("state=firstState&code=code")

		httpMock.On("Get", "https://www.googleapis.com/oauth2/v2/userinfo?access_token=AccessToken").Return(http.StatusUnauthorized, "").Times(1)

		engine.ServeHTTP(w, ctx.Request)

		// Then
		body, err := io.ReadAll(w.Body)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, string(body), "Error while verifying auth token")
	})

	userInfo := web.UserInfo{
		Id:            "123123",
		Email:         "john@doe.com",
		VerifiedEmail: true,
		Picture:       "https://john.doe.com/picture.png",
		GivenName:     "John",
	}
	userString, _ := json.Marshal(&userInfo)

	httpMock.On("Get", "https://www.googleapis.com/oauth2/v2/userinfo?access_token=AccessToken").Return(http.StatusOK, string(userString)).Times(2)

	t.Run("should return jwt token", func(t *testing.T) {
		// When
		w = httptest.NewRecorder()

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/auth/google/callback", nil)
		ctx.Request.AddCookie(&http.Cookie{
			Name:     "oauthstate",
			Value:    "firstState",
			Path:     "/",
			Domain:   "localhost",
			MaxAge:   60,
			Secure:   false,
			HttpOnly: true,
		})
		ctx.Request.Form, _ = url.ParseQuery("state=firstState&code=code")

		engine.ServeHTTP(w, ctx.Request)

		// Then
		body, err := io.ReadAll(w.Body)
		assert.Nil(t, err)

		assert.Contains(t, string(body), "Successfully authenticated user")
	})

	t.Run("should set cookie", func(t *testing.T) {
		// When
		w = httptest.NewRecorder()

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/auth/google/callback", nil)
		ctx.Request.AddCookie(&http.Cookie{
			Name:     "oauthstate",
			Value:    "firstState",
			Path:     "/",
			Domain:   "localhost",
			MaxAge:   60,
			Secure:   false,
			HttpOnly: true,
		})
		ctx.Request.Form, _ = url.ParseQuery("state=firstState&code=code")

		engine.ServeHTTP(w, ctx.Request)

		// Then
		setCookie := w.Header().Get("Set-Cookie")
		assert.NotEmpty(t, setCookie)
	})
}

func TestMe(t *testing.T) {
	// Given
	mock := oAuth2Mock{}
	authController := web.NewAuthController(newTestConfig(), &mock, &http.Client{})

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	authedUser := middleware.AuthedUser{
		Id:             "123",
		Name:           "John",
		Email:          "john@doe.com",
		ProfilePicture: "https://doe.com/profile.png",
	}

	ctx.Set("authedUser", authedUser)

	// When
	authController.Me(ctx)

	// Then
	r, err := io.ReadAll(w.Body)
	assert.Nil(t, err)

	json, err := json.Marshal(authedUser)
	assert.Nil(t, err)

	assert.Equal(t, json, r)
}
