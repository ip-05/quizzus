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
	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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

	database, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	assert.Nil(gs.T(), err)
	gs.controller = NewGameController(database)

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
	gs.T().Log(">>>>>> single test tear down")
	gs.db.Close()
}

func (gs *GameSuite) TestCreateGame_QuestionMin() {
	// Given
	createBody := CreateBody{
		Topic:     "Topic Test",
		RoundTime: 10,
		Points:    10,
		Questions: []CreateQuestion{},
	}

	wantOwner := "123"

	insert := `INSERT INTO "games" ("invite_code","topic","round_time","points","owner") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`

	gs.mock.ExpectBegin()
	gs.mock.ExpectQuery(regexp.QuoteMeta(insert)).
		WithArgs(sqlmock.AnyArg(), createBody.Topic, createBody.RoundTime, createBody.Points, wantOwner).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(1)))
	gs.mock.ExpectCommit()

	bodyString, _ := json.Marshal(&createBody)
	reader := strings.NewReader(string(bodyString))
	gs.ctx.Request = httptest.NewRequest("POST", "/games", reader)

	// When
	gs.controller.CreateGame(gs.ctx)

	// Then
	bodyBytes, _ := io.ReadAll(gs.w.Body)
	assert.Contains(gs.T(), string(bodyBytes), "Should be at least 1 question")

	// err := gs.mock.ExpectationsWereMet()
	// assert.Nil(gs.T(), err)
}

func (gs *GameSuite) TestCreateGame_OptionMin() {
	// Given
	createQuestion := CreateQuestion{
		Name:    "What color is a tomato?",
		Options: []CreateOption{},
	}

	createBody := CreateBody{
		Topic:     "Topic Test",
		RoundTime: 10,
		Points:    10,
		Questions: []CreateQuestion{createQuestion},
	}

	wantOwner := "123"

	insert := `INSERT INTO "games" ("invite_code","topic","round_time","points","owner") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`

	gs.mock.ExpectBegin()
	gs.mock.ExpectQuery(regexp.QuoteMeta(insert)).
		WithArgs(sqlmock.AnyArg(), createBody.Topic, createBody.RoundTime, createBody.Points, wantOwner).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(1)))
	gs.mock.ExpectCommit()

	bodyString, _ := json.Marshal(&createBody)
	reader := strings.NewReader(string(bodyString))
	gs.ctx.Request = httptest.NewRequest("POST", "/games", reader)

	// When
	gs.controller.CreateGame(gs.ctx)

	// Then
	bodyBytes, _ := io.ReadAll(gs.w.Body)
	assert.Contains(gs.T(), string(bodyBytes), "Should be 4 options")

	// err := gs.mock.ExpectationsWereMet()
	// assert.Nil(gs.T(), err)
}

func (gs *GameSuite) TestCreateGame_OK() {
	// Given
	option1 := CreateOption{
		Name:    "blue",
		Correct: false,
	}

	option2 := CreateOption{
		Name:    "brown",
		Correct: false,
	}

	option3 := CreateOption{
		Name:    "red",
		Correct: true,
	}

	option4 := CreateOption{
		Name:    "black",
		Correct: false,
	}

	createQuestion := CreateQuestion{
		Name:    "What color is a tomato?",
		Options: []CreateOption{option1, option2, option3, option4},
	}

	createBody := CreateBody{
		Topic:     "Topic Test",
		RoundTime: 10,
		Points:    10,
		Questions: []CreateQuestion{createQuestion},
	}

	wantOwner := "123"

	insertGame := `INSERT INTO "games" ("invite_code","topic","round_time","points","owner") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`
	insertQuestion := `INSERT INTO "questions" ("name","game_id") VALUES ($1,$2) ON CONFLICT ("id") DO UPDATE SET "name"="excluded"."name","game_id"="excluded"."game_id" RETURNING "id"`
	insertOptions := `INSERT INTO "options" ("name","correct","question_id") VALUES ($1,$2,$3),($4,$5,$6),($7,$8,$9),($10,$11,$12) ON CONFLICT ("id") DO UPDATE SET "name"="excluded"."name","correct"="excluded"."correct","question_id"="excluded"."question_id" RETURNING "id"`

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

func (gs *GameSuite) TestCreateGame_TopicMax() {
	// Given
	createBody := CreateBody{
		Topic:     strings.Repeat(".", 129),
		RoundTime: 10,
		Points:    10,
		Questions: []CreateQuestion{},
	}

	wantOwner := "123"

	insert := `INSERT INTO "games" ("invite_code","topic","round_time","points","owner") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`

	gs.mock.ExpectBegin()
	gs.mock.ExpectQuery(regexp.QuoteMeta(insert)).
		WithArgs(sqlmock.AnyArg(), createBody.Topic, createBody.RoundTime, createBody.Points, wantOwner).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(1)))
	gs.mock.ExpectCommit()

	bodyString, _ := json.Marshal(&createBody)
	reader := strings.NewReader(string(bodyString))
	gs.ctx.Request = httptest.NewRequest("POST", "/games", reader)

	// When
	gs.controller.CreateGame(gs.ctx)

	// Then
	bodyBytes, _ := io.ReadAll(gs.w.Body)
	assert.Contains(gs.T(), string(bodyBytes), "Too long topic name")

	// err := gs.mock.ExpectationsWereMet()
	// assert.Nil(gs.T(), err)
}

func (gs *GameSuite) TestCreateGame_RoundMinMax() {
	// Given
	createBody := CreateBody{
		Topic:     "Topic Test",
		RoundTime: 5,
		Points:    10,
		Questions: []CreateQuestion{},
	}

	wantOwner := "123"

	insert := `INSERT INTO "games" ("invite_code","topic","round_time","points","owner") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`

	gs.mock.ExpectBegin()
	gs.mock.ExpectQuery(regexp.QuoteMeta(insert)).
		WithArgs(sqlmock.AnyArg(), createBody.Topic, createBody.RoundTime, createBody.Points, wantOwner).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(1)))
	gs.mock.ExpectCommit()

	bodyString, _ := json.Marshal(&createBody)
	reader := strings.NewReader(string(bodyString))
	gs.ctx.Request = httptest.NewRequest("POST", "/games", reader)

	// When
	gs.controller.CreateGame(gs.ctx)

	// Then
	bodyBytes, _ := io.ReadAll(gs.w.Body)
	assert.Contains(gs.T(), string(bodyBytes), "Round time should be over 10 or below 60 (seconds)")

	// err := gs.mock.ExpectationsWereMet()
	// assert.Nil(gs.T(), err)
}

func (gs *GameSuite) TestCreateGame_PointsMin() {
	// Given
	createBody := CreateBody{
		Topic:     "Topic Test",
		RoundTime: 15,
		Points:    -5,
		Questions: []CreateQuestion{},
	}

	wantOwner := "123"

	insert := `INSERT INTO "games" ("invite_code","topic","round_time","points","owner") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`

	gs.mock.ExpectBegin()
	gs.mock.ExpectQuery(regexp.QuoteMeta(insert)).
		WithArgs(sqlmock.AnyArg(), createBody.Topic, createBody.RoundTime, createBody.Points, wantOwner).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(1)))
	gs.mock.ExpectCommit()

	bodyString, _ := json.Marshal(&createBody)
	reader := strings.NewReader(string(bodyString))
	gs.ctx.Request = httptest.NewRequest("POST", "/games", reader)

	// When
	gs.controller.CreateGame(gs.ctx)

	// Then
	bodyBytes, _ := io.ReadAll(gs.w.Body)
	assert.Contains(gs.T(), string(bodyBytes), "Points should not be lower than 0")

	// err := gs.mock.ExpectationsWereMet()
	// assert.Nil(gs.T(), err)
}

func (gs *GameSuite) TestGet_NotFound() {
	// Given
	rows := sqlmock.
		NewRows([]string{"id", "invite_code", "topic", "round_time", "points", "owner"})

	selectById := `SELECT * FROM "games" WHERE invite_code = $1 or id = $2 ORDER BY "games"."id" LIMIT 1`

	gs.mock.ExpectQuery(regexp.QuoteMeta(selectById)).WithArgs("", 1).WillReturnRows(rows)

	gs.ctx.Request = httptest.NewRequest("GET", "/games?id=1", nil)

	// When
	gs.controller.Get(gs.ctx)

	// Then
	bodyBytes, _ := io.ReadAll(gs.w.Body)
	assert.Equal(gs.T(), http.StatusNotFound, gs.w.Code)
	assert.Contains(gs.T(), string(bodyBytes), "Game not found")

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

	selectGame := `SELECT * FROM "games" WHERE invite_code = $1 or id = $2 ORDER BY "games"."id" LIMIT 1`
	selectQuestion := `SELECT * FROM "questions" WHERE "questions"."game_id" = $1`
	selectOption := `SELECT * FROM "options" WHERE "options"."question_id" = $1`

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

	selectGame := `SELECT * FROM "games" WHERE invite_code = $1 or id = $2 ORDER BY "games"."id" LIMIT 1`
	selectQuestion := `SELECT * FROM "questions" WHERE "questions"."game_id" = $1`
	selectOption := `SELECT * FROM "options" WHERE "options"."question_id" = $1`

	gs.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("", 1).WillReturnRows(rowsGame)
	gs.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	gs.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)
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
	option1 := UpdateOption{
		Id:      1,
		Name:    "blue",
		Correct: false,
	}

	option2 := UpdateOption{
		Id:      1,
		Name:    "brown",
		Correct: false,
	}

	option3 := UpdateOption{
		Id:      1,
		Name:    "red",
		Correct: true,
	}

	option4 := UpdateOption{
		Id:      1,
		Name:    "black",
		Correct: false,
	}

	updateQuestion := UpdateQuestion{
		Id:      1,
		Name:    "What color is a tomato?",
		Options: []UpdateOption{option1, option2, option3, option4},
	}

	updateBody := UpdateBody{
		Topic:     "Update Topic",
		RoundTime: 10,
		Points:    10,
		Questions: []UpdateQuestion{updateQuestion},
	}

	rowsGame := sqlmock.
		NewRows([]string{"id", "invite_code", "topic", "round_time", "points", "owner"})

	selectGame := `SELECT * FROM "games" WHERE invite_code = $1 or id = $2 ORDER BY "games"."id" LIMIT 1`

	gs.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("", 2).WillReturnRows(rowsGame)

	bodyString, _ := json.Marshal(&updateBody)
	reader := strings.NewReader(string(bodyString))
	gs.ctx.Request = httptest.NewRequest("PATCH", "/games?id=2", reader)

	// When
	gs.controller.Update(gs.ctx)

	// Then
	bodyBytes, _ := io.ReadAll(gs.w.Body)
	assert.Equal(gs.T(), http.StatusNotFound, gs.w.Code)
	assert.Contains(gs.T(), string(bodyBytes), "Game not found")

	err := gs.mock.ExpectationsWereMet()
	assert.Nil(gs.T(), err)
}

func (gs *GameSuite) TestUpdate_NotOwner() {
	// Given
	option1 := UpdateOption{
		Id:      1,
		Name:    "blue",
		Correct: false,
	}

	option2 := UpdateOption{
		Id:      1,
		Name:    "brown",
		Correct: false,
	}

	option3 := UpdateOption{
		Id:      1,
		Name:    "red",
		Correct: true,
	}

	option4 := UpdateOption{
		Id:      1,
		Name:    "black",
		Correct: false,
	}

	updateQuestion := UpdateQuestion{
		Id:      1,
		Name:    "What color is a tomato?",
		Options: []UpdateOption{option1, option2, option3, option4},
	}

	updateBody := UpdateBody{
		Topic:     "Update Topic",
		RoundTime: 10,
		Points:    10,
		Questions: []UpdateQuestion{updateQuestion},
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

	selectGame := `SELECT * FROM "games" WHERE invite_code = $1 or id = $2 ORDER BY "games"."id" LIMIT 1`
	selectQuestion := `SELECT * FROM "questions" WHERE "questions"."game_id" = $1`
	selectOption := `SELECT * FROM "options" WHERE "options"."question_id" = $1`
	updateById := `UPDATE "games" SET "invite_code"=$1,"topic"=$2,"round_time"=$3,"points"=$4,"owner"=$5 WHERE "id" = $6`

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
	assert.Equal(gs.T(), http.StatusForbidden, gs.w.Code)
	assert.Contains(gs.T(), string(bodyBytes), "not owner")

	// err := gs.mock.ExpectationsWereMet()
	// assert.Nil(gs.T(), err)
}

func (gs *GameSuite) TestUpdate_TopicMax() {
	// Given
	option1 := UpdateOption{
		Id:      1,
		Name:    "blue",
		Correct: false,
	}

	option2 := UpdateOption{
		Id:      1,
		Name:    "brown",
		Correct: false,
	}

	option3 := UpdateOption{
		Id:      1,
		Name:    "red",
		Correct: true,
	}

	option4 := UpdateOption{
		Id:      1,
		Name:    "black",
		Correct: false,
	}

	updateQuestion := UpdateQuestion{
		Id:      1,
		Name:    "What color is a tomato?",
		Options: []UpdateOption{option1, option2, option3, option4},
	}

	updateBody := UpdateBody{
		Topic:     strings.Repeat(".", 129),
		RoundTime: 10,
		Points:    10,
		Questions: []UpdateQuestion{updateQuestion},
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

	selectGame := `SELECT * FROM "games" WHERE invite_code = $1 or id = $2 ORDER BY "games"."id" LIMIT 1`
	selectQuestion := `SELECT * FROM "questions" WHERE "questions"."game_id" = $1`
	selectOption := `SELECT * FROM "options" WHERE "options"."question_id" = $1`
	updateById := `UPDATE "games" SET "invite_code"=$1,"topic"=$2,"round_time"=$3,"points"=$4,"owner"=$5 WHERE "id" = $6`

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
	assert.Contains(gs.T(), string(bodyBytes), "Too long topic name")

	// err := gs.mock.ExpectationsWereMet()
	// assert.Nil(gs.T(), err)
}

func (gs *GameSuite) TestUpdate_TimeMinMax() {
	// Given
	option1 := UpdateOption{
		Id:      1,
		Name:    "blue",
		Correct: false,
	}

	option2 := UpdateOption{
		Id:      1,
		Name:    "brown",
		Correct: false,
	}

	option3 := UpdateOption{
		Id:      1,
		Name:    "red",
		Correct: true,
	}

	option4 := UpdateOption{
		Id:      1,
		Name:    "black",
		Correct: false,
	}

	updateQuestion := UpdateQuestion{
		Id:      1,
		Name:    "What color is a tomato?",
		Options: []UpdateOption{option1, option2, option3, option4},
	}

	updateBody := UpdateBody{
		Topic:     "Update Topic",
		RoundTime: 100,
		Points:    10,
		Questions: []UpdateQuestion{updateQuestion},
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

	selectGame := `SELECT * FROM "games" WHERE invite_code = $1 or id = $2 ORDER BY "games"."id" LIMIT 1`
	selectQuestion := `SELECT * FROM "questions" WHERE "questions"."game_id" = $1`
	selectOption := `SELECT * FROM "options" WHERE "options"."question_id" = $1`
	updateById := `UPDATE "games" SET "invite_code"=$1,"topic"=$2,"round_time"=$3,"points"=$4,"owner"=$5 WHERE "id" = $6`

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
	assert.Equal(gs.T(), http.StatusBadRequest, gs.w.Code)
	assert.Contains(gs.T(), string(bodyBytes), "Round time should be over 10 or below 60 (seconds)")

	// err := gs.mock.ExpectationsWereMet()
	// assert.Nil(gs.T(), err)
}

func (gs *GameSuite) TestUpdate_PointsMin() {
	// Given
	option1 := UpdateOption{
		Id:      1,
		Name:    "blue",
		Correct: false,
	}

	option2 := UpdateOption{
		Id:      1,
		Name:    "brown",
		Correct: false,
	}

	option3 := UpdateOption{
		Id:      1,
		Name:    "red",
		Correct: true,
	}

	option4 := UpdateOption{
		Id:      1,
		Name:    "black",
		Correct: false,
	}

	updateQuestion := UpdateQuestion{
		Id:      1,
		Name:    "What color is a tomato?",
		Options: []UpdateOption{option1, option2, option3, option4},
	}

	updateBody := UpdateBody{
		Topic:     "Update Topic",
		RoundTime: 10,
		Points:    -10,
		Questions: []UpdateQuestion{updateQuestion},
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

	selectGame := `SELECT * FROM "games" WHERE invite_code = $1 or id = $2 ORDER BY "games"."id" LIMIT 1`
	selectQuestion := `SELECT * FROM "questions" WHERE "questions"."game_id" = $1`
	selectOption := `SELECT * FROM "options" WHERE "options"."question_id" = $1`
	updateById := `UPDATE "games" SET "invite_code"=$1,"topic"=$2,"round_time"=$3,"points"=$4,"owner"=$5 WHERE "id" = $6`

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
	assert.Equal(gs.T(), http.StatusBadRequest, gs.w.Code)
	assert.Contains(gs.T(), string(bodyBytes), "Points should not be lower than 0")

	// err := gs.mock.ExpectationsWereMet()
	// assert.Nil(gs.T(), err)
}

func (gs *GameSuite) TestUpdate_OptionsMin() {
	// Given
	option1 := UpdateOption{
		Name:    "blue",
		Correct: false,
	}

	option2 := UpdateOption{
		Name:    "brown",
		Correct: false,
	}

	option3 := UpdateOption{
		Name:    "red",
		Correct: true,
	}

	updateQuestion := UpdateQuestion{
		Name:    "What color is a tomato?",
		Options: []UpdateOption{option1, option2, option3},
	}

	updateBody := UpdateBody{
		Topic:     "Update Topic",
		RoundTime: 10,
		Points:    10,
		Questions: []UpdateQuestion{updateQuestion},
	}

	rowsGame := sqlmock.
		NewRows([]string{"id", "invite_code", "topic", "round_time", "points", "owner"}).
		AddRow(uint(1), "1234-4321", "Topic", 30, float64(50), "123")

	rowsQuestion := sqlmock.
		NewRows([]string{"id", "name", "game_id"})

	selectGame := `SELECT * FROM "games" WHERE invite_code = $1 or id = $2 ORDER BY "games"."id" LIMIT 1`
	selectQuestion := `SELECT * FROM "questions" WHERE "questions"."game_id" = $1`
	updateById := `UPDATE "games" SET "invite_code"=$1,"topic"=$2,"round_time"=$3,"points"=$4,"owner"=$5 WHERE "id" = $6`

	gs.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("", 1).WillReturnRows(rowsGame)
	gs.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	gs.mock.ExpectExec(regexp.QuoteMeta(updateById)).WithArgs("1234-4321", updateBody.Topic, updateBody.RoundTime, updateBody.Points, "123", uint(1)).WillReturnResult(sqlmock.NewResult(0, 1))

	bodyString, _ := json.Marshal(&updateBody)
	reader := strings.NewReader(string(bodyString))
	gs.ctx.Request = httptest.NewRequest("PATCH", "/games?id=1", reader)

	// When
	gs.controller.Update(gs.ctx)

	// Then
	bodyBytes, _ := io.ReadAll(gs.w.Body)
	assert.Equal(gs.T(), http.StatusBadRequest, gs.w.Code)
	assert.Contains(gs.T(), string(bodyBytes), "Should be 4 options")

	// err := gs.mock.ExpectationsWereMet()
	// assert.Nil(gs.T(), err)
}

func (gs *GameSuite) TestUpdate_Ok() {
	// Given
	option1 := UpdateOption{
		Id:      1,
		Name:    "blue",
		Correct: false,
	}

	option2 := UpdateOption{
		Id:      1,
		Name:    "brown",
		Correct: false,
	}

	option3 := UpdateOption{
		Id:      1,
		Name:    "red",
		Correct: true,
	}

	option4 := UpdateOption{
		Id:      1,
		Name:    "black",
		Correct: false,
	}

	updateQuestion := UpdateQuestion{
		Id:      1,
		Name:    "What color is a tomato?",
		Options: []UpdateOption{option1, option2, option3, option4},
	}

	updateBody := UpdateBody{
		Topic:     "Update Topic",
		RoundTime: 10,
		Points:    10,
		Questions: []UpdateQuestion{updateQuestion},
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

	selectGame := `SELECT * FROM "games" WHERE invite_code = $1 or id = $2 ORDER BY "games"."id" LIMIT 1`
	selectQuestion := `SELECT * FROM "questions" WHERE "questions"."game_id" = $1`
	selectOption := `SELECT * FROM "options" WHERE "options"."question_id" = $1`
	updateById := `UPDATE "games" SET "invite_code"=$1,"topic"=$2,"round_time"=$3,"points"=$4,"owner"=$5 WHERE "id" = $6`

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

	selectGame := `SELECT * FROM "games" WHERE invite_code = $1 or id = $2 ORDER BY "games"."id" LIMIT 1`

	gs.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("", 2).WillReturnRows(rowsGame)

	gs.ctx.Request = httptest.NewRequest("DELETE", "/games?id=2", nil)

	// When
	gs.controller.Delete(gs.ctx)

	// Then
	bodyBytes, _ := io.ReadAll(gs.w.Body)
	assert.Equal(gs.T(), http.StatusNotFound, gs.w.Code)
	assert.Contains(gs.T(), string(bodyBytes), "Game not found")

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

	selectGame := `SELECT * FROM "games" WHERE invite_code = $1 or id = $2 ORDER BY "games"."id" LIMIT 1`
	selectQuestion := `SELECT * FROM "questions" WHERE "questions"."game_id" = $1`
	selectOption := `SELECT * FROM "options" WHERE "options"."question_id" = $1`
	deleteById := `DELETE FROM "games" WHERE "games"."id" = $1`

	gs.mock.ExpectQuery(regexp.QuoteMeta(selectGame)).WithArgs("", 1).WillReturnRows(rowsGame)
	gs.mock.ExpectQuery(regexp.QuoteMeta(selectQuestion)).WithArgs(1).WillReturnRows(rowsQuestion)
	gs.mock.ExpectQuery(regexp.QuoteMeta(selectOption)).WithArgs(1).WillReturnRows(rowsOption)
	gs.mock.ExpectExec(regexp.QuoteMeta(deleteById)).WithArgs(uint(1)).WillReturnResult(sqlmock.NewResult(0, 1))

	gs.ctx.Request = httptest.NewRequest("DELETE", "/games?id=1", nil)

	// When
	gs.controller.Delete(gs.ctx)

	// Then
	bodyBytes, _ := io.ReadAll(gs.w.Body)
	assert.Equal(gs.T(), http.StatusForbidden, gs.w.Code)
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

	selectGame := `SELECT * FROM "games" WHERE invite_code = $1 or id = $2 ORDER BY "games"."id" LIMIT 1`
	selectQuestion := `SELECT * FROM "questions" WHERE "questions"."game_id" = $1`
	selectOption := `SELECT * FROM "options" WHERE "options"."question_id" = $1`
	deleteById := `DELETE FROM "games" WHERE "games"."id" = $1`

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

func TestGenerateCode(t *testing.T) {
	// Given
	wantLen := 9
	wantContain := "-"

	// When
	gotCode := generateCode()

	// Then
	assert.Equal(t, wantLen, len(gotCode))
	assert.Contains(t, gotCode, wantContain)
}

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(GameSuite))
}
