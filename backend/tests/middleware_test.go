package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/middleware"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	// Given
	cfg := &config.Config{
		Secrets: &config.SecretConfig{
			Jwt: "secretsasdasdasdasdasdasdsdsddsd",
		},
	}

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(w)

	engine.Use(middleware.AuthMiddleware(cfg))
	engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	t.Run("should return forbidden", func(t *testing.T) {
		// When
		ctx.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
		engine.ServeHTTP(w, ctx.Request)

		// Then
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("should pass through", func(t *testing.T) {
		// When
		w = httptest.NewRecorder()

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
		ctx.Request.Header.
			Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImpvaG5AZG9lLmNvbSIsImV4cCI6MTY3Mzk5MDM1OCwiaWQiOiIxMjMxMjMiLCJuYW1lIjoiSm9obiIsInByb2ZpbGVQaWN0dXJlIjoiaHR0cHM6Ly9qb2huLmRvZS5jb20vcGljdHVyZS5wbmcifQ.J2Vx9KoqptH1jFvbXEP-VrPngSC4TfgYAvsj7DVy_J4")
		engine.ServeHTTP(w, ctx.Request)

		// Then
		body, err := io.ReadAll(w.Body)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "OK", string(body))
	})
}

func TestWSMiddleware(t *testing.T) {
	// Given
	cfg := &config.Config{
		Secrets: &config.SecretConfig{
			Jwt: "secretsasdasdasdasdasdasdsdsddsd",
		},
	}

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(w)

	engine.Use(middleware.WSMiddleware(cfg))
	engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	t.Run("should return forbidden", func(t *testing.T) {
		// When
		ctx.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
		engine.ServeHTTP(w, ctx.Request)

		// Then
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("should pass through", func(t *testing.T) {
		// When
		w = httptest.NewRecorder()

		ctx.Request, _ = http.NewRequest(http.MethodGet, "/?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImpvaG5AZG9lLmNvbSIsImV4cCI6MTY3Mzk5MDM1OCwiaWQiOiIxMjMxMjMiLCJuYW1lIjoiSm9obiIsInByb2ZpbGVQaWN0dXJlIjoiaHR0cHM6Ly9qb2huLmRvZS5jb20vcGljdHVyZS5wbmcifQ.J2Vx9KoqptH1jFvbXEP-VrPngSC4TfgYAvsj7DVy_J4", nil)
		engine.ServeHTTP(w, ctx.Request)

		// Then
		body, err := io.ReadAll(w.Body)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "OK", string(body))
	})
}
