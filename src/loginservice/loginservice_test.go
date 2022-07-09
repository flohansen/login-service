package loginservice

import (
	"flhansen/fitter-login-service/src/mocks"
	"flhansen/fitter-login-service/src/testhelper"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigFromEnv(t *testing.T) {
	t.Cleanup(testhelper.CreateTestEnvironment(map[string]string{
		"FITTER_LOGIN_SERVICE_DB_HOST":     "localhost",
		"FITTER_LOGIN_SERVICE_DB_PORT":     "5432",
		"FITTER_LOGIN_SERVICE_DB_NAME":     "testdb",
		"FITTER_LOGIN_SERVICE_HOST":        "0.0.0.0",
		"FITTER_LOGIN_SERVICE_PORT":        "8000",
		"FITTER_LOGIN_SERVICE_JWT_SIGNKEY": "secret",
	}))

	config := NewConfigFromEnv()

	assert.Equal(t, "0.0.0.0", config.Host)
	assert.Equal(t, 8000, config.Port)
	assert.Equal(t, "secret", config.Jwt.SignKey)
}

func TestGetAddr(t *testing.T) {
	t.Cleanup(testhelper.CreateTestEnvironment(map[string]string{
		"FITTER_LOGIN_SERVICE_HOST": "0.0.0.0",
		"FITTER_LOGIN_SERVICE_PORT": "8000",
	}))

	mockedHashEngine := new(mocks.HashEngine)
	mockedAccountRepo := new(mocks.AccountRepository)
	service, _ := NewService(NewConfigFromEnv(), mockedAccountRepo, mockedHashEngine)
	addr := service.GetAddr()

	assert.Equal(t, "0.0.0.0:8000", addr)
}

func TestStart(t *testing.T) {
	t.Cleanup(testhelper.CreateTestEnvironment(map[string]string{
		"FITTER_LOGIN_SERVICE_HOST": "0.0.0.0",
		"FITTER_LOGIN_SERVICE_PORT": "8000",
	}))

	mockedHashEngine := new(mocks.HashEngine)
	mockedAccountRepo := new(mocks.AccountRepository)
	service, _ := NewService(NewConfigFromEnv(), mockedAccountRepo, mockedHashEngine)

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
