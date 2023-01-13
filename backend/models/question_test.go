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

type QuestionModelSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock
}

func (s *QuestionModelSuite) SetupTest() {
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

func (s *QuestionModelSuite) Test_FindById() {
	// Given
	rows := sqlmock.
		NewRows([]string{"id", "name", "game_id"}).
		AddRow(uint(1), "What color are tomatoes?", 1)

	selectById := `SELECT * FROM "questions" WHERE "questions"."id" = $1 ORDER BY "questions"."id" LIMIT 1`

	s.mock.ExpectQuery(regexp.QuoteMeta(selectById)).WithArgs(1).WillReturnRows(rows)

	// When
	var question *Question
	s.db.First(&question, 1)

	// Then
	err := s.mock.ExpectationsWereMet()
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), uint(1), question.Id)
	assert.Equal(s.T(), "What color are tomatoes?", question.Name)
	assert.Equal(s.T(), uint(1), question.GameID)
}

func (s *QuestionModelSuite) Test_Create() {
	// Given
	question := Question{
		Name:   "What color are tomatoes?",
		GameID: 1,
	}

	insert := `INSERT INTO "questions" ("name","game_id") VALUES ($1,$2) RETURNING "id"`
	newId := 1

	s.mock.ExpectBegin()
	s.mock.
		ExpectQuery(regexp.QuoteMeta(insert)).
		WithArgs(question.Name, question.GameID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(newId))
	s.mock.ExpectCommit()

	// When
	s.db.Save(&question)

	// Then
	err := s.mock.ExpectationsWereMet()
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), uint(newId), question.Id)
}

func (s *QuestionModelSuite) Test_Update() {
	// Given
	question := Question{
		Id:     1,
		Name:   "What color are tomatoes?",
		GameID: 1,
	}

	update := `UPDATE "questions" SET "name"=$1,"game_id"=$2 WHERE "id" = $3`

	existingId := 1

	s.mock.ExpectBegin()
	s.mock.
		ExpectExec(regexp.QuoteMeta(update)).
		WithArgs("What color are watermelons?", question.GameID, existingId).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	// When
	question.Name = "What color are watermelons?"
	s.db.Save(&question)

	// Then
	err := s.mock.ExpectationsWereMet()
	assert.Nil(s.T(), err)
}

func (s *QuestionModelSuite) Test_Delete() {
	// Given
	question := Question{
		Id:     1,
		Name:   "What color are tomatoes?",
		GameID: 1,
	}

	delete := `DELETE FROM "questions" WHERE "questions"."id" = $1`

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(delete)).WithArgs(question.Id).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	// When
	s.db.Delete(&question)

	// Then
	err := s.mock.ExpectationsWereMet()
	assert.Nil(s.T(), err)
}

func TestQuestionModel(t *testing.T) {
	suite.Run(t, new(QuestionModelSuite))
}
