package ws

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http/httptest"
	"nhooyr.io/websocket"
	"regexp"
	"testing"
	"time"
)

type WebSocketSuite struct {
	suite.Suite
	ctx        *gin.Context
	controller *CoreController
	mock       sqlmock.Sqlmock
	db         *sql.DB

	serv *httptest.Server
	conn *websocket.Conn
}

func CreateToken(secret string, id string, name string, email string, pfp string) (string, error) {
	secretKey := []byte(secret)
	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":             id,
		"name":           name,
		"email":          email,
		"profilePicture": pfp,
		"exp":            time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	return tokenJWT.SignedString(secretKey)
}

func (w *WebSocketSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.Nil(w.T(), err)
	w.db = db
	w.mock = mock

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 w.db,
		PreferSimpleProtocol: true,
	})

	database, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	assert.Nil(w.T(), err)

	gin.SetMode(gin.TestMode)

	controller := NewCoreController(database)
	ctx, engine := gin.CreateTestContext(httptest.NewRecorder())
	w.ctx = ctx

	cfg := config.Config{Secrets: &config.SecretConfig{Jwt: "secret"}}

	token, err := CreateToken(cfg.Secrets.Jwt, "123123123123", "Test", "test@gmail.com", "https://test.com/test.png")
	assert.Nil(w.T(), err)

	engine.Use(middleware.WSMiddleware(&cfg))
	engine.GET("/ws", controller.HandleWS)

	w.serv = httptest.NewServer(engine)
	w.conn, _, err = websocket.Dial(ctx, fmt.Sprintf("%s/ws?token=%s", w.serv.URL, token), &websocket.DialOptions{})
	assert.Nil(w.T(), err)
}

func (w *WebSocketSuite) TearDownTest() {
	w.conn.Close(websocket.StatusInternalError, "the sky is falling")
	w.serv.Close()
}

func (w *WebSocketSuite) TestPing() {
	// When
	MessageReply(false, "PING").Send(w.conn)

	// Then
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "PONG")
}

func (w *WebSocketSuite) TestGetGame_None() {
	// When
	MessageReply(false, "GET_GAME").Send(w.conn)

	// Then
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "NOT_IN_GAME")
}

func (w *WebSocketSuite) TestJoinGame_NotFound() {
	// Given
	selectQuery := `SELECT * FROM "games" WHERE invite_code = $1 ORDER BY "games"."id" LIMIT 1`

	w.mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).WithArgs("1234-4321").WillReturnRows(sqlmock.NewRows(nil))

	// When
	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(w.conn)

	// Then
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "GAME_NOT_FOUND")
}

func TestWebSocket(t *testing.T) {
	suite.Run(t, new(WebSocketSuite))
}
