package game

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ip-05/quizzus/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServiceSuite struct {
	suite.Suite
	svc  *GameService
	mock sqlmock.Sqlmock
	db   *sql.DB
}

func (s *ServiceSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.Nil(s.T(), err)
	s.db = db

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 s.db,
		PreferSimpleProtocol: true,
	})

	database, err := gorm.Open(dialector)
	repository := repo.NewGameStore(database)
	s.svc = NewGameService(repository)
	assert.Nil(s.T(), err)

	s.mock = mock
}

func (s *ServiceSuite) TearDownTest() {
	s.db.Close()
}

func TestCreateGame(t *testing.T) {

}

func TestUpdateGame(t *testing.T) {

}

func TestDeleteGame(t *testing.T) {

}

func TestGetGame(t *testing.T) {

}
