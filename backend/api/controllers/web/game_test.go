package web

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ip-05/quizzus/app/game"
	"github.com/ip-05/quizzus/entity"
	"github.com/ip-05/quizzus/repo"

	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/api/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var insertGame = `INSERT INTO "games" ("invite_code","topic","round_time","points","owner") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`
var insertQuestion = `INSERT INTO "questions" ("name","game_id") VALUES ($1,$2) ON CONFLICT ("id") DO UPDATE SET "name"="excluded"."name","game_id"="excluded"."game_id" RETURNING "id"`
var insertOptions = `INSERT INTO "options" ("name","correct","question_id") VALUES ($1,$2,$3),($4,$5,$6),($7,$8,$9),($10,$11,$12) ON CONFLICT ("id") DO UPDATE SET "name"="excluded"."name","correct"="excluded"."correct","question_id"="excluded"."question_id" RETURNING "id"`

var selectGame = `SELECT * FROM "games" WHERE invite_code = $1 or id = $2 ORDER BY "games"."id" LIMIT 1`
var selectQuestion = `SELECT * FROM "questions" WHERE "questions"."game_id" = $1`
var selectOption = `SELECT * FROM "options" WHERE "options"."question_id" = $1`

var updateById = `UPDATE "games" SET "invite_code"=$1,"topic"=$2,"round_time"=$3,"points"=$4,"owner"=$5 WHERE "id" = $6`
var deleteById = `DELETE FROM "games" WHERE "games"."id" = $1`

type GameSuite struct {
	suite.Suite
	ctx        *gin.Context
	w          *httptest.ResponseRecorder
	controller *GameController
	mock       sqlmock.Sqlmock
	db         *sql.DB
}

func (gs *GameSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.Nil(gs.T(), err)
	gs.db = db

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 gs.db,
		PreferSimpleProtocol: true,
	})

	database, err := gorm.Open(dialector)
	repository := repo.NewGameStore(database)
	svc := game.NewGameService(repository)
	assert.Nil(gs.T(), err)
	gs.controller = NewGameController(svc)

	gin.SetMode(gin.TestMode)

	gs.w = httptest.NewRecorder()
	gs.ctx, _ = gin.CreateTestContext(gs.w)
	gs.mock = mock

	authedUser := middleware.AuthedUser{
		Id:             "123",
		Name:           "John",
		Email:          "john@doe.com",
		ProfilePicture: "https://doe.com/profile.png",
	}

	gs.ctx.Set("authedUser", authedUser)
}

func (gs *GameSuite) TearDownTest() {
	gs.db.Close()
}

func (gs *GameSuite) TestCreateGame_OK() {
	// Given
	option1 := entity.CreateOption{
		Name:    "blue",
		Correct: false,
	}

	option2 := entity.CreateOption{
		Name:    "brown",
		Correct: false,
	}

	option3 := entity.CreateOption{
		Name:    "red",
		Correct: true,
	}

	option4 := entity.CreateOption{
		Name:    "black",
		Correct: false,
	}

	createQuestion := entity.CreateQuestion{
		Name:    "What color is a tomato?",
		Options: []entity.CreateOption{option1, option2, option3, option4},
	}

	createBody := entity.CreateBody{
		Topic:     "Topic Test",
		RoundTime: 10,
		Points:    10,
		Questions: []entity.CreateQuestion{createQuestion},
	}

	wantOwner := "123"

	gs.mock.ExpectBegin()
	gs.mock.ExpectQuery(regexp.QuoteMeta(insertGame)).
		WithArgs(sqlmock.AnyArg(), createBody.Topic, createBody.RoundTime, createBody.Points, wantOwner).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(1)))

	gs.mock.ExpectQuery(regexp.QuoteMeta(insertQuestion)).
		WithArgs(createQuestion.Name, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(1)))

	gs.mock.ExpectQuery(regexp.QuoteMeta(insertOptions)).
		WithArgs(
			option1.Name, option1.Correct, sqlmock.AnyArg(),
			option2.Name, option2.Correct, sqlmock.AnyArg(),
			option3.Name, option3.Correct, sqlmock.AnyArg(),
			option4.Name, option4.Correct, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(1)))
	gs.mock.ExpectCommit()

	bodyString, _ := json.Marshal(&createBody)
	reader := strings.NewReader(string(bodyString))
	gs.ctx.Request = httptest.NewRequest("POST", "/games", reader)

	// When
	gs.controller.CreateGame(gs.ctx)

	// Then
	bodyBytes, _ := io.ReadAll(gs.w.Body)

	assert.Equal(gs.T(), gs.w.Code, http.StatusOK)
	assert.Contains(gs.T(), string(bodyBytes), createBody.Topic)
	assert.Contains(gs.T(), string(bodyBytes), createQuestion.Name)
	assert.Contains(gs.T(), string(bodyBytes), option1.Name)
	assert.Contains(gs.T(), string(bodyBytes), option2.Name)
	assert.Contains(gs.T(), string(bodyBytes), option3.Name)
	assert.Contains(gs.T(), string(bodyBytes), option4.Name)

	err := gs.mock.ExpectationsWereMet()
	assert.Nil(gs.T(), err)
}

func (gs *GameSuite) TestGet_NotFound() {
	// Given
	rows := sqlmock.
		NewRows([]string{"id", "invite_code", "topic", "round_time", "points", "owner"})

	gs.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("", 1).WillReturnRows(rows)

	gs.ctx.Request = httptest.NewRequest("GET", "/games?id=1", nil)

	// When
	gs.controller.Get(gs.ctx)

	// Then
	bodyBytes, _ := io.ReadAll(gs.w.Body)
	assert.Contains(gs.T(), string(bodyBytes), "game not found")

	err := gs.mock.ExpectationsWereMet()
	assert.Nil(gs.T(), err)
}

func (gs *GameSuite) TestGet_NotOwner() {
	// Given
	rowsGame := sqlmock.
		NewRows([]string{"id", "invite_code", "topic", "round_time", "points", "owner"}).
		AddRow(uint(1), "1234-4321", "My quiz", 30, float64(10), "321")

	rowsQuestion := sqlmock.
		NewRows([]string{"id", "name", "game_id"}).
		AddRow(uint(1), "My question", uint(1))

	rowsOption := sqlmock.
		NewRows([]string{"id", "name", "correct", "question_id"}).
		AddRow(uint(1), "option1", false, uint(1)).
		AddRow(uint(2), "option2", false, uint(1)).
		AddRow(uint(3), "option3", true, uint(1)).
		AddRow(uint(4), "option4", false, uint(1))

	gs.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("", 1).WillReturnRows(rowsGame)
	gs.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	gs.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)

	// When
	gs.ctx.Request = httptest.NewRequest("GET", "/games?id=1", nil)
	gs.controller.Get(gs.ctx)

	// Then
	bodyBytes, _ := io.ReadAll(gs.w.Body)

	assert.Equal(gs.T(), http.StatusOK, gs.w.Code)
	assert.Contains(gs.T(), string(bodyBytes), "My quiz")
	assert.NotContains(gs.T(), string(bodyBytes), "1234-4321")

	// err := gs.mock.ExpectationsWereMet()
	// assert.Nil(gs.T(), err)
}

func (gs *GameSuite) TestGet_Ok() {
	// Given
	rowsGame := sqlmock.
		NewRows([]string{"id", "invite_code", "topic", "round_time", "points", "owner"}).
		AddRow(uint(1), "1234-4321", "My quiz", 30, float64(10), "123")

	rowsQuestion := sqlmock.
		NewRows([]string{"id", "name", "game_id"}).
		AddRow(uint(1), "My question", uint(1))

	rowsOption := sqlmock.
		NewRows([]string{"id", "name", "correct", "question_id"}).
		AddRow(uint(1), "option1", false, uint(1)).
		AddRow(uint(2), "option2", false, uint(1)).
		AddRow(uint(3), "option3", true, uint(1)).
		AddRow(uint(4), "option4", false, uint(1))

	gs.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("", 1).WillReturnRows(rowsGame)
	gs.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	gs.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)

	// When
	gs.ctx.Request = httptest.NewRequest("GET", "/games?id=1", nil)
	gs.controller.Get(gs.ctx)

	// Then
	bodyBytes, _ := io.ReadAll(gs.w.Body)

	assert.Equal(gs.T(), http.StatusOK, gs.w.Code)
	assert.Contains(gs.T(), string(bodyBytes), "My quiz")
	assert.Contains(gs.T(), string(bodyBytes), "1234-4321")
	assert.Contains(gs.T(), string(bodyBytes), "My question")
	assert.Contains(gs.T(), string(bodyBytes), "option1")
	assert.Contains(gs.T(), string(bodyBytes), "option2")
	assert.Contains(gs.T(), string(bodyBytes), "option3")
	assert.Contains(gs.T(), string(bodyBytes), "option4")

	// err := gs.mock.ExpectationsWereMet()
	// assert.Nil(gs.T(), err)
}

func (gs *GameSuite) TestUpdate_NotFound() {
	// Given
	option1 := entity.UpdateOption{
		Id:      1,
		Name:    "blue",
		Correct: false,
	}

	option2 := entity.UpdateOption{
		Id:      1,
		Name:    "brown",
		Correct: false,
	}

	option3 := entity.UpdateOption{
		Id:      1,
		Name:    "red",
		Correct: true,
	}

	option4 := entity.UpdateOption{
		Id:      1,
		Name:    "black",
		Correct: false,
	}

	updateQuestion := entity.UpdateQuestion{
		Id:      1,
		Name:    "What color is a tomato?",
		Options: []entity.UpdateOption{option1, option2, option3, option4},
	}

	updateBody := entity.UpdateBody{
		Topic:     "Update Topic",
		RoundTime: 10,
		Points:    10,
		Questions: []entity.UpdateQuestion{updateQuestion},
	}

	rowsGame := sqlmock.
		NewRows([]string{"id", "invite_code", "topic", "round_time", "points", "owner"})

	gs.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("", 2).WillReturnRows(rowsGame)

	bodyString, _ := json.Marshal(&updateBody)
	reader := strings.NewReader(string(bodyString))
	gs.ctx.Request = httptest.NewRequest("PATCH", "/games?id=2", reader)

	// When
	gs.controller.Update(gs.ctx)

	// Then
	bodyBytes, _ := io.ReadAll(gs.w.Body)
	assert.Equal(gs.T(), http.StatusBadRequest, gs.w.Code)
	assert.Contains(gs.T(), string(bodyBytes), "game not found")

	err := gs.mock.ExpectationsWereMet()
	assert.Nil(gs.T(), err)
}

func (gs *GameSuite) TestUpdate_NotOwner() {
	// Given
	option1 := entity.UpdateOption{
		Id:      1,
		Name:    "blue",
		Correct: false,
	}

	option2 := entity.UpdateOption{
		Id:      1,
		Name:    "brown",
		Correct: false,
	}

	option3 := entity.UpdateOption{
		Id:      1,
		Name:    "red",
		Correct: true,
	}

	option4 := entity.UpdateOption{
		Id:      1,
		Name:    "black",
		Correct: false,
	}

	updateQuestion := entity.UpdateQuestion{
		Id:      1,
		Name:    "What color is a tomato?",
		Options: []entity.UpdateOption{option1, option2, option3, option4},
	}

	updateBody := entity.UpdateBody{
		Topic:     "Update Topic",
		RoundTime: 10,
		Points:    10,
		Questions: []entity.UpdateQuestion{updateQuestion},
	}

	rowsGame := sqlmock.
		NewRows([]string{"id", "invite_code", "topic", "round_time", "points", "owner"}).
		AddRow(uint(1), "1234-4321", "Topic", 30, float64(50), "321")

	rowsQuestion := sqlmock.
		NewRows([]string{"id", "name", "game_id"}).
		AddRow(uint(1), "My question", uint(1))

	rowsOption := sqlmock.
		NewRows([]string{"id", "name", "correct", "question_id"}).
		AddRow(uint(1), "option1", false, uint(1)).
		AddRow(uint(2), "option2", false, uint(1)).
		AddRow(uint(3), "option3", true, uint(1)).
		AddRow(uint(4), "option4", false, uint(1))

	gs.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("", 1).WillReturnRows(rowsGame)
	gs.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	gs.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)
	gs.mock.ExpectExec(regexp.QuoteMeta(updateById)).WithArgs("1234-4321", updateBody.Topic, updateBody.RoundTime, updateBody.Points, "321", uint(1)).WillReturnResult(sqlmock.NewResult(0, 1))

	bodyString, _ := json.Marshal(&updateBody)
	reader := strings.NewReader(string(bodyString))
	gs.ctx.Request = httptest.NewRequest("PATCH", "/games?id=1", reader)

	// When
	gs.controller.Update(gs.ctx)

	// Then
	bodyBytes, _ := io.ReadAll(gs.w.Body)
	assert.Equal(gs.T(), http.StatusBadRequest, gs.w.Code)
	assert.Contains(gs.T(), string(bodyBytes), "not owner")

	// err := gs.mock.ExpectationsWereMet()
	// assert.Nil(gs.T(), err)
}

func (gs *GameSuite) TestUpdate_Ok() {
	// Given
	option1 := entity.UpdateOption{
		Id:      1,
		Name:    "blue",
		Correct: false,
	}

	option2 := entity.UpdateOption{
		Id:      1,
		Name:    "brown",
		Correct: false,
	}

	option3 := entity.UpdateOption{
		Id:      1,
		Name:    "red",
		Correct: true,
	}

	option4 := entity.UpdateOption{
		Id:      1,
		Name:    "black",
		Correct: false,
	}

	updateQuestion := entity.UpdateQuestion{
		Id:      1,
		Name:    "What color is a tomato?",
		Options: []entity.UpdateOption{option1, option2, option3, option4},
	}

	updateBody := entity.UpdateBody{
		Topic:     "Update Topic",
		RoundTime: 10,
		Points:    10,
		Questions: []entity.UpdateQuestion{updateQuestion},
	}

	rowsGame := sqlmock.
		NewRows([]string{"id", "invite_code", "topic", "round_time", "points", "owner"}).
		AddRow(uint(1), "1234-4321", "Topic", 30, float64(50), "123")

	rowsQuestion := sqlmock.
		NewRows([]string{"id", "name", "game_id"}).
		AddRow(uint(1), "My question", uint(1))

	rowsOption := sqlmock.
		NewRows([]string{"id", "name", "correct", "question_id"}).
		AddRow(uint(1), "option1", false, uint(1)).
		AddRow(uint(2), "option2", false, uint(1)).
		AddRow(uint(3), "option3", true, uint(1)).
		AddRow(uint(4), "option4", false, uint(1))

	gs.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("", 1).WillReturnRows(rowsGame)
	gs.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	gs.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)
	gs.mock.ExpectExec(regexp.QuoteMeta(updateById)).WithArgs("1234-4321", updateBody.Topic, updateBody.RoundTime, updateBody.Points, "123", uint(1)).WillReturnResult(sqlmock.NewResult(0, 1))

	bodyString, _ := json.Marshal(&updateBody)
	reader := strings.NewReader(string(bodyString))
	gs.ctx.Request = httptest.NewRequest("PATCH", "/games?id=1", reader)

	// When
	gs.controller.Update(gs.ctx)

	// Then
	bodyBytes, _ := io.ReadAll(gs.w.Body)
	assert.Equal(gs.T(), http.StatusOK, gs.w.Code)
	assert.Contains(gs.T(), string(bodyBytes), updateBody.Topic)
	assert.Contains(gs.T(), string(bodyBytes), updateQuestion.Name)
	assert.Contains(gs.T(), string(bodyBytes), option1.Name)
	assert.Contains(gs.T(), string(bodyBytes), option4.Name)

	// err := gs.mock.ExpectationsWereMet()
	// assert.Nil(gs.T(), err)
}

func (gs *GameSuite) TestDelete_NotFound() {
	// Given
	rowsGame := sqlmock.
		NewRows([]string{"id", "invite_code", "topic", "round_time", "points", "owner"})

	gs.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("", 2).WillReturnRows(rowsGame)

	gs.ctx.Request = httptest.NewRequest("DELETE", "/games?id=2", nil)

	// When
	gs.controller.Delete(gs.ctx)

	// Then
	bodyBytes, _ := io.ReadAll(gs.w.Body)
	assert.Contains(gs.T(), string(bodyBytes), "game not found")

	err := gs.mock.ExpectationsWereMet()
	assert.Nil(gs.T(), err)
}

func (gs *GameSuite) TestDelete_NotOwner() {
	// Given
	rowsGame := sqlmock.
		NewRows([]string{"id", "invite_code", "topic", "round_time", "points", "owner"}).
		AddRow(uint(1), "1234-4321", "Topic", 30, float64(50), "321")

	rowsQuestion := sqlmock.
		NewRows([]string{"id", "name", "game_id"}).
		AddRow(uint(1), "My question", uint(1))

	rowsOption := sqlmock.
		NewRows([]string{"id", "name", "correct", "question_id"}).
		AddRow(uint(1), "option1", false, uint(1)).
		AddRow(uint(2), "option2", false, uint(1)).
		AddRow(uint(3), "option3", true, uint(1)).
		AddRow(uint(4), "option4", false, uint(1))

	gs.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("", 1).WillReturnRows(rowsGame)
	gs.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	gs.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)
	gs.mock.ExpectExec(regexp.QuoteMeta(deleteById)).WithArgs(uint(1)).WillReturnResult(sqlmock.NewResult(0, 1))

	gs.ctx.Request = httptest.NewRequest("DELETE", "/games?id=1", nil)

	// When
	gs.controller.Delete(gs.ctx)

	// Then
	bodyBytes, _ := io.ReadAll(gs.w.Body)
	assert.Contains(gs.T(), string(bodyBytes), "not owner")

	// err := gs.mock.ExpectationsWereMet()
	// assert.Nil(gs.T(), err)
}

func (gs *GameSuite) TestDelete_Ok() {
	// Given
	rowsGame := sqlmock.
		NewRows([]string{"id", "invite_code", "topic", "round_time", "points", "owner"}).
		AddRow(uint(1), "1234-4321", "Topic", 30, float64(50), "123")

	rowsQuestion := sqlmock.
		NewRows([]string{"id", "name", "game_id"}).
		AddRow(uint(1), "My question", uint(1))

	rowsOption := sqlmock.
		NewRows([]string{"id", "name", "correct", "question_id"}).
		AddRow(uint(1), "option1", false, uint(1)).
		AddRow(uint(2), "option2", false, uint(1)).
		AddRow(uint(3), "option3", true, uint(1)).
		AddRow(uint(4), "option4", false, uint(1))

	gs.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("", 1).WillReturnRows(rowsGame)
	gs.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	gs.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)
	gs.mock.ExpectExec(regexp.QuoteMeta(deleteById)).WithArgs(uint(1)).WillReturnResult(sqlmock.NewResult(0, 1))

	gs.ctx.Request = httptest.NewRequest("DELETE", "/games?id=1", nil)

	// When
	gs.controller.Delete(gs.ctx)

	// Then
	bodyBytes, _ := io.ReadAll(gs.w.Body)
	assert.Equal(gs.T(), http.StatusOK, gs.w.Code)
	assert.Contains(gs.T(), string(bodyBytes), "Successfully deleted")

	// err := gs.mock.ExpectationsWereMet()
	// assert.Nil(gs.T(), err)
}

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(GameSuite))
}
