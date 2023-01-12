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

type ModelsSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock
}

func (s *ModelsSuite) SetupTest() {
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

func (s *ModelsSuite) TestMigrations_Option() {
	// Given
	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

	selectSchema := `SELECT count(*) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = $1 AND table_type = $2`
	create := `CREATE TABLE "options" ("id" bigserial,"name" text,"correct" boolean,"question_id" bigint,PRIMARY KEY ("id"))`

	s.mock.ExpectQuery(regexp.QuoteMeta(selectSchema)).WithArgs("options", "BASE TABLE").WillReturnRows(rows)
	s.mock.ExpectExec(regexp.QuoteMeta(create)).WillReturnResult(sqlmock.NewResult(0, 1))

	// When
	err := s.db.AutoMigrate(&Option{})

	// Then
	assert.Nil(s.T(), err)

	err = s.mock.ExpectationsWereMet()
	assert.Nil(s.T(), err)
}

func (s *ModelsSuite) TestMigrations_Question() {
	// Given
	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

	selectSchema := `SELECT count(*) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = $1 AND table_type = $2`
	create := `CREATE TABLE "questions" ("id" bigserial,"name" text,"game_id" bigint,PRIMARY KEY ("id"))`

	s.mock.ExpectQuery(regexp.QuoteMeta(selectSchema)).WithArgs("questions", "BASE TABLE").WillReturnRows(rows)
	s.mock.ExpectExec(regexp.QuoteMeta(create)).WillReturnResult(sqlmock.NewResult(0, 1))

	// When
	err := s.db.AutoMigrate(&Question{})

	// Then
	assert.Nil(s.T(), err)

	err = s.mock.ExpectationsWereMet()
	assert.Nil(s.T(), err)
}

func (s *ModelsSuite) TestMigrations_Game() {
	// Given
	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

	selectSchema := `SELECT count(*) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = $1 AND table_type = $2`
	create := `CREATE TABLE "games" ("id" bigserial,"invite_code" text,"topic" text,"round_time" bigint,"points" decimal,"owner" text,PRIMARY KEY ("id"))`

	s.mock.ExpectQuery(regexp.QuoteMeta(selectSchema)).WithArgs("games", "BASE TABLE").WillReturnRows(rows)
	s.mock.ExpectExec(regexp.QuoteMeta(create)).WillReturnResult(sqlmock.NewResult(0, 1))

	// When
	err := s.db.AutoMigrate(&Game{})

	// Then
	assert.Nil(s.T(), err)

	err = s.mock.ExpectationsWereMet()
	assert.Nil(s.T(), err)
}

func TestModels(t *testing.T) {
	suite.Run(t, new(ModelsSuite))
}
