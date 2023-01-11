package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

func TestOptionModel(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(t, err)

	rows := sqlmock.
		NewRows([]string{"id", "name", "correct"}).
		AddRow(uint(1), "Red", true)

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

	selectById := "SELECT * FROM \"options\" WHERE id = $1 ORDER BY \"options\".\"id\" LIMIT 1"

	t.Run("should find option by id", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(selectById)).WithArgs(1).WillReturnRows(rows)

		var option *Option
		database.Where("id = ?", 1).First(&option)

		assert.Equal(t, uint(1), option.Id)
		assert.Equal(t, "Red", option.Name)
		assert.Equal(t, true, option.Correct)
	})
}
