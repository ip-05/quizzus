package ws

import (
	"context"
	"database/sql"
	"fmt"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

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
	"nhooyr.io/websocket"
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

func (w *WebSocketSuite) TestLeaveGame_None() {
	// When
	MessageReply(false, "LEAVE_GAME").Send(w.conn)

	// Then
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "NOT_IN_GAME")
}

func (w *WebSocketSuite) TestIsOwner_None() {
	// When
	MessageReply(false, "IS_OWNER").Send(w.conn)

	// Then
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "NOT_IN_GAME")
}

func (w *WebSocketSuite) TestStartGame_None() {
	// When
	MessageReply(false, "START_GAME").Send(w.conn)

	// Then
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "NOT_IN_GAME")
}

func (w *WebSocketSuite) TestResetGame_None() {
	// When
	MessageReply(false, "RESET_GAME").Send(w.conn)

	// Then
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "NOT_IN_GAME")
}

func (w *WebSocketSuite) TestAnswerQuestion_None() {
	// When
	MessageReply(false, "ANSWER_QUESTION").Send(w.conn)

	// Then
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "NOT_IN_GAME")
}

func (w *WebSocketSuite) TestNextRound_None() {
	// When
	MessageReply(false, "NEXT_ROUND").Send(w.conn)

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

func (w *WebSocketSuite) TestJoinGame() {
	// Given
	rowsGame := sqlmock.
		NewRows([]string{"id", "invite_code", "topic", "round_time", "points", "owner"}).
		AddRow(uint(1), "1234-4321", "Topic", 30, float64(50), "123123123123")

	rowsQuestion := sqlmock.
		NewRows([]string{"id", "name", "game_id"}).
		AddRow(uint(1), "My question", uint(1))

	rowsOption := sqlmock.
		NewRows([]string{"id", "name", "correct", "question_id"}).
		AddRow(uint(1), "option1", false, uint(1)).
		AddRow(uint(2), "option2", false, uint(1)).
		AddRow(uint(3), "option3", true, uint(1)).
		AddRow(uint(4), "option4", false, uint(1))

	selectGame := `SELECT * FROM "games" WHERE invite_code = $1 ORDER BY "games"."id" LIMIT 1`
	selectQuestion := `SELECT * FROM "questions" WHERE "questions"."game_id" = $1`
	selectOption := `SELECT * FROM "options" WHERE "options"."question_id" = $1`

	w.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("1234-4321").WillReturnRows(rowsGame)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)

	// When
	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(w.conn)

	// Then
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "JOINED_GAME")
}

func (w *WebSocketSuite) TestJoinGame_AlreadyInGame() {

}

func (w *WebSocketSuite) TestGetGame() {
	// Given
	rowsGame := sqlmock.
		NewRows([]string{"id", "invite_code", "topic", "round_time", "points", "owner"}).
		AddRow(uint(1), "1234-4321", "Topic", 30, float64(50), "123123123123")

	rowsQuestion := sqlmock.
		NewRows([]string{"id", "name", "game_id"}).
		AddRow(uint(1), "My question", uint(1))

	rowsOption := sqlmock.
		NewRows([]string{"id", "name", "correct", "question_id"}).
		AddRow(uint(1), "option1", false, uint(1)).
		AddRow(uint(2), "option2", false, uint(1)).
		AddRow(uint(3), "option3", true, uint(1)).
		AddRow(uint(4), "option4", false, uint(1))

	selectGame := `SELECT * FROM "games" WHERE invite_code = $1 ORDER BY "games"."id" LIMIT 1`
	selectQuestion := `SELECT * FROM "questions" WHERE "questions"."game_id" = $1`
	selectOption := `SELECT * FROM "options" WHERE "options"."question_id" = $1`

	w.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("1234-4321").WillReturnRows(rowsGame)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)
	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(w.conn)
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "JOINED_GAME")

	// When
	DataReply(false, "GET_GAME", JoinGameData{GameId: "1234-4321"}).Send(w.conn)

	// Then
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "GET_GAME")
}

func (w *WebSocketSuite) TestLeaveGame() {

}

func (w *WebSocketSuite) TestStartGame() {

}

func (w *WebSocketSuite) TestIsOwner_True() {

}

func (w *WebSocketSuite) TestIsOwner_False() {

}

func (w *WebSocketSuite) TestAnswerQuestion_Standby() {

}

func (w *WebSocketSuite) TestAnswerQuestion_Success() {

}

func (w *WebSocketSuite) TestNextRound_InProgress() {

}

func (w *WebSocketSuite) TestNextRound_Success() {

}

func (w *WebSocketSuite) TestResetGame_InProgress() {

}

func (w *WebSocketSuite) TestResetGame_Success() {

}

func TestWebSocket(t *testing.T) {
	suite.Run(t, new(WebSocketSuite))
}
