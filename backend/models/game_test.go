package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

func TestGameModel(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(t, err)

	rows := sqlmock.
		NewRows([]string{"id", "invite_code", "topic", "round_time", "points", "owner"}).
		AddRow(uint(1), "1234-4321", "My quiz", 30, float64(10), "1234432112344321")

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

	selectById := "SELECT * FROM \"games\" WHERE id = $1 ORDER BY \"games\".\"id\" LIMIT 1"

	t.Run("should find game by id", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(selectById)).WithArgs(1).WillReturnRows(rows)

		var game *Game
		database.Where("id = ?", 1).First(&game)

		assert.Equal(t, uint(1), game.Id)
		assert.Equal(t, "1234-4321", game.InviteCode)
		assert.Equal(t, "My quiz", game.Topic)
		assert.Equal(t, 30, game.RoundTime)
		assert.Equal(t, float64(10), game.Points)
		assert.Equal(t, "1234432112344321", game.Owner)
	})
}
