package loginservice

import (
	"flhansen/fitter-login-service/src/database"
	"flhansen/fitter-login-service/src/testhelper"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewWithConfig(t *testing.T) {
	service, err := New(LoginServiceConfig{})
	assert.NotNil(t, service)
	assert.Nil(t, err)
}

func TestNewWithInvalidConfig(t *testing.T) {
	service, err := New(LoginServiceConfig{Host: "0.0.0.0", Port: -1, Database: database.DatabaseConfig{Port: -1}})
	assert.Nil(t, service)
	assert.NotNil(t, err)
}

func TestNewConfigFromEnv(t *testing.T) {
	t.Cleanup(testhelper.CreateTestEnvironment(map[string]string{
		"FITTER_LOGIN_SERVICE_DB_HOST":     "localhost",
		"FITTER_LOGIN_SERVICE_DB_PORT":     "5432",
		"FITTER_LOGIN_SERVICE_DB_NAME":     "testdb",
		"FITTER_LOGIN_SERVICE_DB_USER":     "testuser",
		"FITTER_LOGIN_SERVICE_HOST":        "0.0.0.0",
		"FITTER_LOGIN_SERVICE_PORT":        "8000",
		"FITTER_LOGIN_SERVICE_JWT_SIGNKEY": "secret",
		"AWS_REGION":                       "eu-central-1",
	}))

	config := NewConfigFromEnv()

	assert.Equal(t, "localhost", config.Database.Host)
	assert.Equal(t, 5432, config.Database.Port)
	assert.Equal(t, "testdb", config.Database.Name)
	assert.Equal(t, "testuser", config.Database.User)
	assert.Equal(t, "0.0.0.0", config.Host)
	assert.Equal(t, 8000, config.Port)
	assert.Equal(t, "secret", config.Jwt.SignKey)
	assert.Equal(t, "eu-central-1", config.Database.Region)
}

func TestGetAddr(t *testing.T) {
	t.Cleanup(testhelper.CreateTestEnvironment(map[string]string{
		"FITTER_LOGIN_SERVICE_HOST": "0.0.0.0",
		"FITTER_LOGIN_SERVICE_PORT": "8000",
	}))

	service, _ := New(NewConfigFromEnv())
	addr := service.GetAddr()

	assert.Equal(t, "0.0.0.0:8000", addr)
}

func TestStart(t *testing.T) {
	t.Cleanup(testhelper.CreateTestEnvironment(map[string]string{
		"FITTER_LOGIN_SERVICE_HOST": "0.0.0.0",
		"FITTER_LOGIN_SERVICE_PORT": "8000",
	}))

	service, _ := New(NewConfigFromEnv())

	done := make(chan error)
	go func() {
		done <- service.Start()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		return
	case err := <-done:
		t.Fatal(err)
	}
}

func TestServer(t *testing.T) {
	t.Cleanup(testhelper.CreateTestEnvironment(map[string]string{
		"FITTER_LOGIN_SERVICE_HOST": "0.0.0.0",
		"FITTER_LOGIN_SERVICE_PORT": "8001",
	}))

	service, _ := New(NewConfigFromEnv())
	server := service.Server()

	done := make(chan error)
	go func() {
		done <- server.ListenAndServe()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		return
	case err := <-done:
		t.Fatal(err)
	}
}
