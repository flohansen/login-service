package loginservice

import (
	"bytes"
	"encoding/json"
	"flhansen/fitter-login-service/src/repository"
	"flhansen/fitter-login-service/src/security"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/localstack"
	"github.com/orlangure/gnomock/preset/postgres"
	"github.com/stretchr/testify/suite"
)

type LoginServiceTestSuite struct {
	suite.Suite
	accountRepo  repository.AccountRepository
	loginService *LoginService
	localstack   *gnomock.Container
	database     *gnomock.Container
	session      *session.Session
}

func TestLoginServiceTestSuite(t *testing.T) {
	suite.Run(t, new(LoginServiceTestSuite))
}

func (s *LoginServiceTestSuite) setupLocalstack() {
	preset := localstack.Preset(
		localstack.WithVersion("0.12.2"))
	s.localstack, _ = gnomock.Start(preset)

	s.T().Cleanup(func() { gnomock.Stop(s.localstack) })
}

func (s *LoginServiceTestSuite) setupAwsLocalstack() {
	s.session, _ = session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("test", "test", ""),
		Region:      aws.String(endpoints.EuCentral1RegionID),
		Endpoint:    aws.String(fmt.Sprintf("http://%s/", s.localstack.Address(localstack.APIPort))),
		DisableSSL:  aws.Bool(true),
	})
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
	s.setupLocalstack()
	s.setupAwsLocalstack()
	s.setupDatabase()

	accountRepo, err := repository.NewAccountRepository(
		s.database.Host,
		s.database.DefaultPort(),
		"test",
		"test",
		"test",
	)
	if err != nil {
		s.T().Fatal(err)
	}

	hashEngine := security.NewBcryptEngine()
	loginService, err := NewService(LoginServiceConfig{
		Host: "0.0.0.0",
		Port: 8080,
	}, accountRepo, hashEngine)
	if err != nil {
		s.T().Fatal(err)
	}

	s.accountRepo = accountRepo
	s.loginService = loginService
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
