package tests

import (
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"net/url"
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

type OAuth2Mock struct {
	mock.Mock
}

func (o OAuth2Mock) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
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

func (o OAuth2Mock) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken: "AccessToken",
		Expiry:      time.Now().Add(1 * time.Hour),
	}, nil
}

func TestGoogleLogin(t *testing.T) {
	// Given
	mock := OAuth2Mock{}
	authController := web.NewAuthController(newTestConfig(), &mock)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// When
	authController.GoogleLogin(ctx)

	// Then
	m := map[string]string{}

	r, err := io.ReadAll(w.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(r, &m)
	assert.Nil(t, err)

	redirect := m["redirectUrl"]
	assert.NotNil(t, redirect)
}
