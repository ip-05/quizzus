package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type WSMiddlewareSuite struct {
	suite.Suite
	ctx    *gin.Context
	engine *gin.Engine
	w      *httptest.ResponseRecorder
}

func (s *WSMiddlewareSuite) SetupTest() {
	// Given
	cfg := &config.Config{
		Secrets: &config.SecretConfig{
			Jwt: "secretsasdasdasdasdasdasdsdsddsd",
		},
	}

	wsMiddleware := WSMiddleware(cfg)

	gin.SetMode(gin.TestMode)
	s.w = httptest.NewRecorder()

	s.ctx, s.engine = gin.CreateTestContext(s.w)
	s.engine.Use(wsMiddleware)
	s.engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
}

func (s *WSMiddlewareSuite) Test_Forbidden() {
	// When
	s.ctx.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	s.engine.ServeHTTP(s.w, s.ctx.Request)

	// Then
	assert.Equal(s.T(), http.StatusForbidden, s.w.Code)
}

func (s *WSMiddlewareSuite) Test_OK() {
	// When
	s.ctx.Request, _ = http.NewRequest(http.MethodGet, "/?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImpvaG5AZG9lLmNvbSIsImV4cCI6MTY3Mzk5MDM1OCwiaWQiOiIxMjMxMjMiLCJuYW1lIjoiSm9obiIsInByb2ZpbGVQaWN0dXJlIjoiaHR0cHM6Ly9qb2huLmRvZS5jb20vcGljdHVyZS5wbmcifQ.J2Vx9KoqptH1jFvbXEP-VrPngSC4TfgYAvsj7DVy_J4", nil)
	s.engine.ServeHTTP(s.w, s.ctx.Request)

	// Then
	body, err := io.ReadAll(s.w.Body)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), http.StatusOK, s.w.Code)
	assert.Equal(s.T(), "OK", string(body))
}

func TestWSMiddleware(t *testing.T) {
	suite.Run(t, new(WSMiddlewareSuite))
}
