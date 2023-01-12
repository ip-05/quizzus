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

type OptionModelSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock
}

func (s *OptionModelSuite) SetupTest() {
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

func (s *OptionModelSuite) Test_FindById() {
	// Given
	rows := sqlmock.
		NewRows([]string{"id", "name", "correct"}).
		AddRow(uint(1), "Red", true)

	selectById := `SELECT * FROM "options" WHERE "options"."id" = $1 ORDER BY "options"."id" LIMIT 1`

	s.mock.ExpectQuery(regexp.QuoteMeta(selectById)).WithArgs(1).WillReturnRows(rows)

	// When
	var option *Option
	s.db.First(&option, 1)

	// Then
	err := s.mock.ExpectationsWereMet()
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), uint(1), option.Id)
	assert.Equal(s.T(), "Red", option.Name)
	assert.Equal(s.T(), true, option.Correct)
}

func (s *OptionModelSuite) Test_Create() {
	// Given
	option := Option{
		Name:    "Red",
		Correct: true,
	}

	insert := `INSERT INTO "options" ("name","correct","question_id") VALUES ($1,$2,$3) RETURNING "id"`
	newId := 1

	s.mock.ExpectBegin()
	s.mock.
		ExpectQuery(regexp.QuoteMeta(insert)).
		WithArgs(option.Name, option.Correct, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newId))
	s.mock.ExpectCommit()

	// When
	s.db.Save(&option)

	// Then
	err := s.mock.ExpectationsWereMet()
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), uint(newId), option.Id)
}

func (s *OptionModelSuite) Test_Update() {
	// Given
	option := Option{
		Id:      1,
		Name:    "Red",
		Correct: true,
	}

	update := `UPDATE "options" SET "name"=$1,"correct"=$2,"question_id"=$3 WHERE "id" = $4`

	existingId := 1

	s.mock.ExpectBegin()
	s.mock.
		ExpectExec(regexp.QuoteMeta(update)).
		WithArgs(option.Name, false, option.QuestionID, existingId).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	// When
	option.Correct = false
	s.db.Save(&option)

	// Then
	err := s.mock.ExpectationsWereMet()
	assert.Nil(s.T(), err)
}

func (s *OptionModelSuite) Test_Delete() {
	// Given
	option := Option{
		Id:      1,
		Name:    "Red",
		Correct: true,
	}

	delete := `DELETE FROM "options" WHERE "options"."id" = $1`

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(delete)).WithArgs(option.Id).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	// When
	s.db.Delete(&option)

	// Then
	err := s.mock.ExpectationsWereMet()
	assert.Nil(s.T(), err)
}

func TestOptionModel(t *testing.T) {
	suite.Run(t, new(OptionModelSuite))
}
