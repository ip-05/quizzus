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
	"github.com/ip-05/quizzus/api/middleware"
	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"nhooyr.io/websocket"
)

var selectGame = `SELECT * FROM "games" WHERE invite_code = $1 ORDER BY "games"."id" LIMIT 1`
var selectQuestion = `SELECT * FROM "questions" WHERE "questions"."game_id" = $1`
var selectOption = `SELECT * FROM "options" WHERE "options"."question_id" = $1`
var selectOption2 = `SELECT * FROM "options" WHERE "options"."question_id" IN ($1,$2)`

type WebSocketSuite struct {
	suite.Suite
	ctx     *gin.Context
	mock    sqlmock.Sqlmock
	db      *sql.DB
	newGame models.Game
	serv    *httptest.Server
	conn    *websocket.Conn
}

func createToken(secret string, id string, name string, email string, pfp string) (string, error) {
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

func createRows(game models.Game, two bool) (*sqlmock.Rows, *sqlmock.Rows, *sqlmock.Rows) {
	rowsGame := sqlmock.
		NewRows([]string{"id", "invite_code", "topic", "round_time", "points", "owner"}).
		AddRow(game.Id, game.InviteCode, game.Topic, game.RoundTime, game.Points, game.Owner)

	rowsQuestion := sqlmock.
		NewRows([]string{"id", "name", "game_id"}).
		AddRow(uint(1), "My question", uint(1))

	rowsOption := sqlmock.
		NewRows([]string{"id", "name", "correct", "question_id"}).
		AddRow(uint(1), "option1", false, uint(1)).
		AddRow(uint(2), "option2", false, uint(1)).
		AddRow(uint(3), "option3", true, uint(1)).
		AddRow(uint(4), "option4", false, uint(1))

	if two {
		rowsQuestion.AddRow(uint(2), "My question2", uint(1))

		rowsOption.AddRow(uint(5), "option5", true, uint(2))
		rowsOption.AddRow(uint(6), "option6", false, uint(2))
		rowsOption.AddRow(uint(7), "option7", false, uint(2))
		rowsOption.AddRow(uint(8), "option8", false, uint(2))
	}

	return rowsGame, rowsQuestion, rowsOption
}

func (w *WebSocketSuite) SetupTest() {
	w.newGame = models.Game{
		Id:         uint(1),
		InviteCode: "1234-4321",
		Topic:      "Topic",
		RoundTime:  30,
		Points:     float64(50),
		Owner:      "123123123123",
	}

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

	database, err := gorm.Open(dialector)
	assert.Nil(w.T(), err)

	gin.SetMode(gin.TestMode)

	controller := NewCoreController(database)
	controller.gameController.GameTime = 0

	ctx, engine := gin.CreateTestContext(httptest.NewRecorder())
	w.ctx = ctx

	cfg := config.Config{Secrets: &config.SecretConfig{Jwt: "secret"}}

	token, err := createToken(cfg.Secrets.Jwt, w.newGame.Owner, "Test", "test@gmail.com", "https://test.com/test.png")
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
	w.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("1234-4321").WillReturnRows(sqlmock.NewRows(nil))

	// When
	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(w.conn)

	// Then
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "GAME_NOT_FOUND")
}

func (w *WebSocketSuite) TestJoinGame_Owner() {
	// Given
	rowsGame, rowsQuestion, rowsOption := createRows(w.newGame, false)

	w.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("1234-4321").WillReturnRows(rowsGame)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)

	// When
	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(w.conn)

	// Then
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "JOINED_GAME")
	assert.Contains(w.T(), string(message), "questionCount")
	assert.Contains(w.T(), string(message), "members")
}

func (w *WebSocketSuite) TestJoinGame_NotOwner() {
	// Given
	w.newGame.Owner = "321"
	rowsGame, rowsQuestion, rowsOption := createRows(w.newGame, false)

	w.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("1234-4321").WillReturnRows(rowsGame)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)

	// When
	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(w.conn)

	// Then
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "NOT_OWNER")
}

func (w *WebSocketSuite) TestUserJoined() {
	// Given
	rowsGame, rowsQuestion, rowsOption := createRows(w.newGame, false)

	w.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("1234-4321").WillReturnRows(rowsGame)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)

	rowsGame2, rowsQuestion2, rowsOption2 := createRows(w.newGame, false)

	w.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("1234-4321").WillReturnRows(rowsGame2)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion2)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption2)

	token, err := createToken("secret", "321", "TestUser", "user@gmail.com", "https://test.com/test.png")
	assert.Nil(w.T(), err)

	conn, _, err := websocket.Dial(context.Background(), fmt.Sprintf("%s/ws?token=%s", w.serv.URL, token), &websocket.DialOptions{})
	assert.Nil(w.T(), err)

	defer conn.Close(websocket.StatusInternalError, "the sky is falling")

	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(w.conn)
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "JOINED_GAME")

	// When
	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(conn)

	// Then
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "USER_JOINED")
}

func (w *WebSocketSuite) TestJoinGame_AlreadyInGame() {
	// Given
	rowsGame, rowsQuestion, rowsOption := createRows(w.newGame, false)

	w.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("1234-4321").WillReturnRows(rowsGame)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)

	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(w.conn)
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "JOINED_GAME")

	// When
	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(w.conn)

	// Then
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "ALREADY_IN_GAME")
}

func (w *WebSocketSuite) TestGetGame() {
	// Given
	rowsGame, rowsQuestion, rowsOption := createRows(w.newGame, false)

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

func (w *WebSocketSuite) TestLeaveGame_Owner() {
	// Given
	rowsGame, rowsQuestion, rowsOption := createRows(w.newGame, false)

	w.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("1234-4321").WillReturnRows(rowsGame)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)
	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(w.conn)
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "JOINED_GAME")

	// When
	MessageReply(false, "LEAVE_GAME").Send(w.conn)

	// Then
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "GAME_DELETED")
}

func (w *WebSocketSuite) TestLeaveGame_Player() {
	// Given
	rowsGame, rowsQuestion, rowsOption := createRows(w.newGame, false)

	w.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("1234-4321").WillReturnRows(rowsGame)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)

	rowsGame2, rowsQuestion2, rowsOption2 := createRows(w.newGame, false)

	w.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("1234-4321").WillReturnRows(rowsGame2)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion2)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption2)

	token, err := createToken("secret", "321", "TestUser", "user@gmail.com", "https://test.com/test.png")
	assert.Nil(w.T(), err)

	conn, _, err := websocket.Dial(context.Background(), fmt.Sprintf("%s/ws?token=%s", w.serv.URL, token), &websocket.DialOptions{})
	assert.Nil(w.T(), err)

	defer conn.Close(websocket.StatusInternalError, "the sky is falling")

	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(w.conn)
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "JOINED_GAME")

	// When
	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(conn)
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "USER_JOINED")
	_, message, _ = conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "JOINED_GAME")

	// Then
	MessageReply(false, "LEAVE_GAME").Send(conn)
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "USER_LEFT")
	_, message, _ = conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "LEFT_GAME")
}

func (w *WebSocketSuite) TestStartGame_InProgress() {
	// Given
	rowsGame, rowsQuestion, rowsOption := createRows(w.newGame, false)

	w.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("1234-4321").WillReturnRows(rowsGame)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)

	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(w.conn)
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "JOINED_GAME")

	// When
	MessageReply(false, "START_GAME").Send(w.conn)

	// Then
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "GAME_IN_PROGRESS")
}

func (w *WebSocketSuite) TestIsOwner_True() {
	// Given
	rowsGame, rowsQuestion, rowsOption := createRows(w.newGame, false)

	w.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("1234-4321").WillReturnRows(rowsGame)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)

	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(w.conn)
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "JOINED_GAME")

	// When
	MessageReply(false, "IS_OWNER").Send(w.conn)
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "IS_OWNER")
	assert.Contains(w.T(), string(message), "true")
}

func (w *WebSocketSuite) TestIsOwner_False() {
	// Given
	rowsGame, rowsQuestion, rowsOption := createRows(w.newGame, false)

	w.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("1234-4321").WillReturnRows(rowsGame)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)

	w.newGame.Owner = "321"
	rowsGame2, rowsQuestion2, rowsOption2 := createRows(w.newGame, false)

	w.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("1234-4321").WillReturnRows(rowsGame2)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion2)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption2)

	token, err := createToken("secret", "321", "TestUser", "user@gmail.com", "https://test.com/test.png")
	assert.Nil(w.T(), err)

	conn, _, err := websocket.Dial(context.Background(), fmt.Sprintf("%s/ws?token=%s", w.serv.URL, token), &websocket.DialOptions{})
	assert.Nil(w.T(), err)

	defer conn.Close(websocket.StatusInternalError, "the sky is falling")

	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(w.conn)
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "JOINED_GAME")

	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(conn)
	_, message, _ = conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "JOINED_GAME")

	// When
	MessageReply(false, "IS_OWNER").Send(conn)

	// Then
	_, message, _ = conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "IS_OWNER")
	assert.Contains(w.T(), string(message), "false")
}

func (w *WebSocketSuite) TestPlayRounds() {
	// Given
	w.newGame.RoundTime = 0
	rowsGame, rowsQuestion, rowsOption := createRows(w.newGame, false)

	w.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("1234-4321").WillReturnRows(rowsGame)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)

	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(w.conn)
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "JOINED_GAME")

	// When
	MessageReply(false, "START_GAME").Send(w.conn)
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "GAME_IN_PROGRESS")

	// Then
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "ROUND_FINISHED")
	assert.Contains(w.T(), string(message), "option4")

	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "GAME_FINISHED")
}

func (w *WebSocketSuite) TestAnswerQuestion_Standby() {
	// Given
	w.newGame.RoundTime = 0
	rowsGame, rowsQuestion, rowsOption := createRows(w.newGame, false)

	w.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("1234-4321").WillReturnRows(rowsGame)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)

	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(w.conn)
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "JOINED_GAME")

	// When
	DataReply(false, "ANSWER_QUESTION", AnswerData{Option: 1}).Send(w.conn)

	// Then
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "ROUND_WAITING")
}

func (w *WebSocketSuite) TestAnswerQuestion_Success() {
	// Given
	w.newGame.RoundTime = 1
	rowsGame, rowsQuestion, rowsOption := createRows(w.newGame, false)

	w.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("1234-4321").WillReturnRows(rowsGame)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)

	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(w.conn)
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "JOINED_GAME")

	MessageReply(false, "START_GAME").Send(w.conn)
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "GAME_IN_PROGRESS")
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "ROUND_IN_PROGRESS")

	// When
	DataReply(false, "ANSWER_QUESTION", AnswerData{Option: uint(2)}).Send(w.conn)
	_, message, _ = w.conn.Read(context.Background())

	// Then
	assert.Contains(w.T(), string(message), "ANSWER_ACCEPTED")
	assert.Contains(w.T(), string(message), "leaderboard")

	DataReply(false, "USER_ANSWERED", AnswerResponse{
		UserId: "123123123123",
		Option: uint(2),
	}).Send(w.conn)
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "USER_ANSWERED")
	assert.Contains(w.T(), string(message), "option")
}

func (w *WebSocketSuite) TestNextRound_InProgress() {
	// Given
	w.newGame.RoundTime = 1
	rowsGame, rowsQuestion, rowsOption := createRows(w.newGame, true)

	w.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("1234-4321").WillReturnRows(rowsGame)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectOption2)).WithArgs(1, 2).WillReturnRows(rowsOption)

	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(w.conn)
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "JOINED_GAME")

	MessageReply(false, "START_GAME").Send(w.conn)
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "GAME_IN_PROGRESS")
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "ROUND_IN_PROGRESS")

	// When
	MessageReply(false, "NEXT_ROUND").Send(w.conn)

	// Then
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "GAME_IN_PROGRESS")
}

func (w *WebSocketSuite) TestNextRound_Success() {
	// Given
	w.newGame.RoundTime = 1
	rowsGame, rowsQuestion, rowsOption := createRows(w.newGame, true)

	w.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("1234-4321").WillReturnRows(rowsGame)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectOption2)).WithArgs(1, 2).WillReturnRows(rowsOption)

	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(w.conn)
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "JOINED_GAME")

	msg1 := []string{"GAME_IN_PROGRESS", "ROUND_IN_PROGRESS", "ROUND_FINISHED"}
	MessageReply(false, "START_GAME").Send(w.conn)
	for _, v := range msg1 {
		_, message, _ = w.conn.Read(context.Background())
		assert.Contains(w.T(), string(message), v)
	}

	msg2 := []string{"ROUND_IN_PROGRESS", "ROUND_IN_PROGRESS", "ROUND_FINISHED", "GAME_FINISHED"}
	// When
	MessageReply(false, "NEXT_ROUND").Send(w.conn)
	for _, v := range msg2 {
		// Then
		_, message, _ = w.conn.Read(context.Background())
		assert.Contains(w.T(), string(message), v)
	}
}

func (w *WebSocketSuite) TestResetGame_InProgress() {
	// Given
	w.newGame.RoundTime = 1
	rowsGame, rowsQuestion, rowsOption := createRows(w.newGame, false)

	w.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("1234-4321").WillReturnRows(rowsGame)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)

	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(w.conn)
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "JOINED_GAME")

	MessageReply(false, "START_GAME").Send(w.conn)
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "GAME_IN_PROGRESS")
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "ROUND_IN_PROGRESS")

	// When
	MessageReply(false, "RESET_GAME").Send(w.conn)

	// Then
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "GAME_IN_PROGRESS")
}

func (w *WebSocketSuite) TestResetGame_Success() {
	// Given
	w.newGame.RoundTime = 1
	rowsGame, rowsQuestion, rowsOption := createRows(w.newGame, false)

	w.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("1234-4321").WillReturnRows(rowsGame)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	w.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)

	DataReply(false, "JOIN_GAME", JoinGameData{GameId: "1234-4321"}).Send(w.conn)
	_, message, _ := w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "JOINED_GAME")

	MessageReply(false, "START_GAME").Send(w.conn)
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "GAME_IN_PROGRESS")
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "ROUND_IN_PROGRESS")
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "ROUND_FINISHED")
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "GAME_FINISHED")

	// When
	MessageReply(false, "RESET_GAME").Send(w.conn)

	// Then
	_, message, _ = w.conn.Read(context.Background())
	assert.Contains(w.T(), string(message), "RESET_GAME")
}

func TestWebSocket(t *testing.T) {
	suite.Run(t, new(WebSocketSuite))
}
