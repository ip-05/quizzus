package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

type GameModelSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock
}

func (s *GameModelSuite) SetupTest() {
	// Given
	db, mock, err := sqlmock.New()
	assert.Nil(s.T(), err)

	dialector := postgres.New(postgres.Config{
		DSN:        "sqlmock_db_0",
		DriverName: "postgres",
		Conn:       db,
	})

	s.db, err = gorm.Open(dialector)
	if err != nil {
		s.T().Errorf("Failed to open gorm db, got error: %v", err)
	}

	s.mock = mock
}

func (s *GameModelSuite) Test_FindById() {
	// Given
	rows := sqlmock.
		NewRows([]string{"id", "invite_code", "topic", "round_time", "points", "owner"}).
		AddRow(uint(1), "1234-4321", "My quiz", 30, float64(10), "1234432112344321")

	selectById := `SELECT * FROM "games" WHERE "games"."id" = $1 ORDER BY "games"."id" LIMIT 1`

	s.mock.ExpectQuery(regexp.QuoteMeta(selectById)).WithArgs(1).WillReturnRows(rows)

	// When
	var game *Game
	s.db.First(&game, 1)

	// Then
	err := s.mock.ExpectationsWereMet()
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), uint(1), game.Id)
	assert.Equal(s.T(), "1234-4321", game.InviteCode)
	assert.Equal(s.T(), "My quiz", game.Topic)
	assert.Equal(s.T(), 30, game.RoundTime)
	assert.Equal(s.T(), float64(10), game.Points)
	assert.Equal(s.T(), "1234432112344321", game.Owner)
}

func (s *GameModelSuite) Test_Create() {
	// Given
	game := Game{
		InviteCode: "1234-4321",
		Topic:      "My quiz",
		RoundTime:  30,
		Points:     10,
		Questions:  []Question{},
		Owner:      "1234432112344321",
	}

	insert := `INSERT INTO "games" ("invite_code","topic","round_time","points","owner") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`
	newId := 1

	s.mock.ExpectBegin()
	s.mock.
		ExpectQuery(regexp.QuoteMeta(insert)).
		WithArgs(game.InviteCode, game.Topic, game.RoundTime, game.Points, game.Owner).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newId))
	s.mock.ExpectCommit()

	// When
	s.db.Save(&game)

	// Then
	err := s.mock.ExpectationsWereMet()
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), uint(newId), game.Id)
}

func (s *GameModelSuite) Test_Update() {
	// Given
	game := Game{
		Id:         1,
		InviteCode: "1234-4321",
		Topic:      "My quiz",
		RoundTime:  30,
		Points:     10,
		Questions:  []Question{},
		Owner:      "1234432112344321",
	}

	update := `UPDATE "games" SET "invite_code"=$1,"topic"=$2,"round_time"=$3,"points"=$4,"owner"=$5 WHERE "id" = $6`

	existingId := 1

	s.mock.ExpectBegin()
	s.mock.
		ExpectExec(regexp.QuoteMeta(update)).
		WithArgs(game.InviteCode, "New quiz", game.RoundTime, game.Points, game.Owner, existingId).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	// When
	game.Topic = "New quiz"
	s.db.Save(&game)

	// Then
	err := s.mock.ExpectationsWereMet()
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), "New quiz", game.Topic)
}

func (s *GameModelSuite) Test_Delete() {
	// Given
	game := Game{
		Id:         1,
		InviteCode: "1234-4321",
		Topic:      "My quiz",
		RoundTime:  30,
		Points:     10,
		Questions:  []Question{},
		Owner:      "1234432112344321",
	}

	delete := `DELETE FROM "games" WHERE "games"."id" = $1`

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(delete)).WithArgs(game.Id).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	// When
	s.db.Delete(&game)

	// Then
	err := s.mock.ExpectationsWereMet()
	assert.Nil(s.T(), err)
}

func TestGameModel(t *testing.T) {
	suite.Run(t, new(GameModelSuite))
}
