package repository

import (
	"database/sql"
	"flhansen/fitter-login-service/src/database"
	"testing"
	"time"

	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
	"github.com/stretchr/testify/suite"
)

type AccountRepositoryTestSuite struct {
	suite.Suite
	database *gnomock.Container
	repo     AccountRepository
	db       *sql.DB
}

func TestAccountRepository(t *testing.T) {
	suite.Run(t, new(AccountRepositoryTestSuite))
}

func (suite *AccountRepositoryTestSuite) SetupSuite() {
	preset := postgres.Preset(
		postgres.WithUser("test", "test"),
		postgres.WithDatabase("test"),
		postgres.WithQueries(QUERY_CREATE_ACCOUNT_TABLE))
	suite.database, _ = gnomock.Start(preset)
	suite.T().Cleanup(func() { gnomock.Stop(suite.database) })

	suite.repo = NewAccountRepository(DatabaseConfig{
		Host:         suite.database.Host,
		Port:         suite.database.DefaultPort(),
		Username:     "test",
		Password:     "test",
		DatabaseName: "test",
	})

	dsn := database.DataSourceName(
		suite.database.Host,
		suite.database.DefaultPort(),
		"test", "test", "test")
	suite.db, _ = sql.Open("postgres", dsn)
}

func (suite *AccountRepositoryTestSuite) TearDownTest() {
	_, err := suite.db.Exec("DELETE FROM account")
	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *AccountRepositoryTestSuite) TestCreateAccountShouldReturnErrorIfAccountAlreadyExists() {
	suite.db.QueryRow("INSERT INTO account (username, password, email, creation_date) VALUES ($1, $2, $3, $4)",
		"test", "test", "test@test.com", time.Now())

	id, err := suite.repo.CreateAccount(Account{
		Username:     "test",
		Password:     "test",
		Email:        "test@test.com",
		CreationDate: time.Now(),
	})

	suite.Error(err)
	suite.Equal(-1, id)
}

func (suite *AccountRepositoryTestSuite) TestCreateAccountShouldReturnErrorIfScanFails() {
	suite.db.QueryRow("INSERT INTO account (username, password, email, creation_date) VALUES ($1, $2, $3, $4)",
		"test", "test", "test@test.com", time.Now())

	id, err := suite.repo.CreateAccount(Account{
		Username: "test",
	})

	suite.Error(err)
	suite.Equal(-1, id)
}

func (suite *AccountRepositoryTestSuite) TestCreateAccountShouldSucceed() {
	id, err := suite.repo.CreateAccount(Account{
		Username:     "test",
		Password:     "test",
		Email:        "test@test.com",
		CreationDate: time.Now(),
	})

	suite.NoError(err)
	suite.NotEqual(-1, id)
}

func (suite *AccountRepositoryTestSuite) TestGetAccountByIdShouldReturnErrorIfScanFails() {
	_, err := suite.repo.GetAccountById(-1)
	suite.Error(err)
}

func (suite *AccountRepositoryTestSuite) TestGetAccountByIdShouldSucceed() {
	username := "test"
	password := "test"
	email := "test"
	creationDate := time.Now()

	row := suite.db.QueryRow("INSERT INTO account (username, password, email, creation_date) VALUES ($1, $2, $3, $4) RETURNING id",
		username, password, email, creationDate)

	id := -1
	err := row.Scan(&id)
	if err != nil {
		suite.T().Fatal(err)
	}

	user, err := suite.repo.GetAccountById(id)

	suite.NoError(err)
	suite.Equal(id, user.Id)
	suite.Equal(username, user.Username)
	suite.Equal(password, user.Password)
	suite.Equal(email, user.Email)
	suite.Equal(creationDate.UnixMilli(), user.CreationDate.UnixMilli())
}

func (suite *AccountRepositoryTestSuite) TestGetAccountByUsernameShouldReturnErrorIfScanFails() {
	_, err := suite.repo.GetAccountByUsername("test")
	suite.Error(err)
}

func (suite *AccountRepositoryTestSuite) TestGetAccountByUsernameShouldSucceed() {
	username := "test"
	password := "test"
	email := "test"
	creationDate := time.Now()

	suite.db.QueryRow("INSERT INTO account (username, password, email, creation_date) VALUES ($1, $2, $3, $4) RETURNING id",
		username, password, email, creationDate)

	user, err := suite.repo.GetAccountByUsername(username)

	suite.NoError(err)
	suite.NotEqual(-1, user.Id)
	suite.Equal(username, user.Username)
	suite.Equal(password, user.Password)
	suite.Equal(email, user.Email)
	suite.Equal(creationDate.UnixMilli(), user.CreationDate.UnixMilli())
}

func (suite *AccountRepositoryTestSuite) TestDeleteAccountByIdShouldReturnErrorIfExecFails() {
	err := suite.repo.DeleteAccountById(-1)
	suite.Error(err)
}

func (suite *AccountRepositoryTestSuite) TestDeleteAccountByIdShouldSucceed() {
	username := "test"
	password := "test"
	email := "test"
	creationDate := time.Now()

	row := suite.db.QueryRow("INSERT INTO account (username, password, email, creation_date) VALUES ($1, $2, $3, $4) RETURNING id",
		username, password, email, creationDate)

	id := -1
	err := row.Scan(&id)
	if err != nil {
		suite.T().Fatal(err)
	}

	err = suite.repo.DeleteAccountById(id)

	suite.NoError(err)
}

func (suite *AccountRepositoryTestSuite) TestDeleteAccountsSucceed() {
	err := suite.repo.DeleteAccounts()
	suite.NoError(err)
}
