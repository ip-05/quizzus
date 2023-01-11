package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

func TestQuestionModel(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(t, err)

	rows := sqlmock.
		NewRows([]string{"id", "name"}).
		AddRow(uint(1), "What color are tomatoes?")

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	database, err := gorm.Open(dialector)
	if err != nil {
		t.Errorf("Failed to open gorm db, got error: %v", err)
	}

	selectById := "SELECT * FROM \"questions\" WHERE id = $1 ORDER BY \"questions\".\"id\" LIMIT 1"

	t.Run("should find question by id", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(selectById)).WithArgs(1).WillReturnRows(rows)

		var question *Question
		database.Where("id = ?", 1).First(&question)

		assert.Equal(t, uint(1), question.Id)
		assert.Equal(t, "What color are tomatoes?", question.Name)
	})
}
