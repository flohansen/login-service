package loginservice

import (
	"bytes"
	"encoding/json"
	"flhansen/fitter-login-service/src/repository"
	"flhansen/fitter-login-service/src/security"
	"net/http"
	"testing"
	"time"

	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type LoginServiceTestSuite struct {
	suite.Suite
	accountRepo  repository.AccountRepository
	loginService *LoginService
	database     *gnomock.Container
}

func TestLoginServiceTestSuite(t *testing.T) {
	suite.Run(t, new(LoginServiceTestSuite))
}

func (s *LoginServiceTestSuite) setupDatabase() {
	preset := postgres.Preset(
		postgres.WithUser("test", "test"),
		postgres.WithDatabase("test"),
		postgres.WithQueries(repository.QUERY_CREATE_ACCOUNT_TABLE))
	database, err := gnomock.Start(preset)
	if err != nil {
		s.T().Fatal(err)
	}

	s.database = database
	s.T().Cleanup(func() { gnomock.Stop(s.database) })
}

func (s *LoginServiceTestSuite) SetupSuite() {
	s.setupDatabase()

	s.accountRepo = repository.NewAccountRepository(repository.DatabaseConfig{
		Host:         s.database.Host,
		Port:         s.database.DefaultPort(),
		Username:     "test",
		Password:     "test",
		DatabaseName: "test",
	})

	hashEngine := security.NewBcryptEngine()
	logger := logrus.New()

	s.loginService = NewService(LoginServiceConfig{
		Host: "0.0.0.0",
		Port: 8080,
	}, s.accountRepo, hashEngine, logger)

	go s.loginService.Start()
}

func (s *LoginServiceTestSuite) TearDownSuite() {
}

func (s *LoginServiceTestSuite) SetupTest() {

}

func (s *LoginServiceTestSuite) TearDownTest() {
	_ = s.accountRepo.DeleteAccounts()
}

func (s *LoginServiceTestSuite) TestServiceShouldRegisterUser() {
	body, _ := json.Marshal(map[string]interface{}{
		"username": "test",
		"password": "test",
		"email":    "test@test.com",
	})

	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/auth/register", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		s.T().Fatal(err)
	}

	var resp map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&resp)

	s.NoError(err)
	s.Equal(http.StatusOK, res.StatusCode)
	s.NotNil(resp["userId"])
}

func (s *LoginServiceTestSuite) TestServiceShouldLoginUserAndReceiveAuthToken() {
	hashEngine := security.NewBcryptEngine()
	hashedPassword, _ := hashEngine.HashPassword([]byte("test"))
	s.accountRepo.CreateAccount(repository.Account{
		Username:     "test",
		Password:     string(hashedPassword),
		Email:        "test@test.com",
		CreationDate: time.Now(),
	})

	body, _ := json.Marshal(map[string]interface{}{
		"username": "test",
		"password": "test",
	})

	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/auth/login", bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		s.T().Fatal(err)
	}

	var resp map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&resp)

	s.NoError(err)
	s.Equal(http.StatusOK, res.StatusCode)
	s.NotNil(resp["token"])
}
