package main

import (
	"bytes"
	"flhansen/fitter-login-service/src/database"
	"flhansen/fitter-login-service/src/loginservice"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/elgohr/go-localstack"
	"github.com/stretchr/testify/suite"
)

type LoginServiceTestSuite struct {
	suite.Suite
	loginService *loginservice.LoginService
	localstack   *localstack.Instance
	session      *session.Session
}

func TestLoginServiceTestSuite(t *testing.T) {
	suite.Run(t, new(LoginServiceTestSuite))
}

func (s *LoginServiceTestSuite) setupLocalstack() {
	var err error
	s.localstack, err = localstack.NewInstance(localstack.WithVersion("0.12.3"))
	if err != nil {
		s.T().Fatal(err)
	}

	if err := s.localstack.Start(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *LoginServiceTestSuite) setupAwsLocalstack() {
	s.session, _ = session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("test", "test", ""),
		Region:      aws.String(endpoints.EuCentral1RegionID),
		Endpoint:    aws.String(s.localstack.Endpoint(localstack.SecretsManager)),
		DisableSSL:  aws.Bool(true),
	})
}

func (s *LoginServiceTestSuite) SetupSuite() {
	s.setupLocalstack()
	s.setupAwsLocalstack()

	sqlDatabase := database.NewPostgresDatabase(database.DatabaseConfig{Host: "localhost", Port: 5432, Name: "test"})

	s.loginService = loginservice.NewService(loginservice.LoginServiceConfig{
		Host: "0.0.0.0",
		Port: 8080,
	}, sqlDatabase, hashEngine)
	go s.loginService.Start()
}

func (s *LoginServiceTestSuite) TearDownSuite() {

}

func (s *LoginServiceTestSuite) SetupTest() {

}

func (s *LoginServiceTestSuite) TearDownTest() {

}

func (s *LoginServiceTestSuite) TestServiceShouldRegisterUser() {
	body := []byte(`{"username":"testuser","password":"testpass","email":"test@test.com"}`)
	req, _ := http.NewRequest(http.MethodPost, "http://"+s.loginService.GetAddr()+"/api/auth/register", bytes.NewBuffer(body))

	client := &http.Client{}
	res, err := client.Do(req)

	s.NoError(err)
	s.Equal(http.StatusOK, res.StatusCode)
}

func (s *LoginServiceTestSuite) TestServiceShouldLoginUserAndReceiveAuthToken() {
	s.FailNow("TODO: Implement me")
}
