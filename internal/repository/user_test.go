package repository

import (
	"regexp"
	"testing"
	"time"

	"chatto/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/suite"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

type Suite struct {
	suite.Suite
	Db    *gorm.DB
	mocks sqlmock.Sqlmock

	repo IUserRepository
	data string
}

func (s *Suite) SetupSuite() {
	db, mocks, err := sqlmock.New()
	require.NoError(s.T(), err, "Failed to create mock for the database connection:")

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	require.NoError(s.T(), err, "Failed to create gorm db")

	s.repo = NewUserRepository(gormDb)
	s.mocks = mocks
	s.Db = gormDb
}

func (s *Suite) TestFindUserById() {
	user := &model.User{
		Id:        uuid.NewString(),
		Name:      "mizhan",
		Password:  "asd",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	sql := regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1`)

	s.mocks.ExpectQuery(sql).
		WithArgs(user.Id).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "password", "created_at", "updated_at"}).
			AddRow(user.Id, user.Name, user.Password, user.CreatedAt, user.UpdatedAt))

	user2, err := s.repo.FindUserById(user.Id)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), user2)
	require.Equal(s.T(), user, user2)
}

func (s *Suite) TestFindUserByName() {
	user := &model.User{
		Id:        uuid.NewString(),
		Name:      "mizhan",
		Password:  "asd",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	sql := regexp.QuoteMeta(`SELECT * FROM "users" WHERE name = $1`)

	s.mocks.ExpectQuery(sql).
		WithArgs(user.Name).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "password", "created_at", "updated_at"}).
			AddRow(user.Id, user.Name, user.Password, user.CreatedAt, user.UpdatedAt))

	user2, err := s.repo.FindUserByName(user.Name)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), user2)
	require.Equal(s.T(), user, user2)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
