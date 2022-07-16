package loginservice

import (
	"flhansen/fitter-login-service/src/mocks"
	"flhansen/fitter-login-service/src/testhelper"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetAddr(t *testing.T) {
	t.Cleanup(testhelper.CreateTestEnvironment(map[string]string{
		"LOGIN_SERVICE_HOST": "0.0.0.0",
		"LOGIN_SERVICE_PORT": "8000",
	}))

	mockedHashEngine := new(mocks.HashEngine)
	mockedAccountRepo := new(mocks.AccountRepository)
	service := NewService(NewConfigFromEnv(), mockedAccountRepo, mockedHashEngine)
	addr := service.GetAddr()

	assert.Equal(t, "0.0.0.0:8000", addr)
}

func TestStart(t *testing.T) {
	t.Cleanup(testhelper.CreateTestEnvironment(map[string]string{
		"LOGIN_SERVICE_HOST": "0.0.0.0",
		"LOGIN_SERVICE_PORT": "8000",
	}))

	mockedHashEngine := new(mocks.HashEngine)
	mockedAccountRepo := new(mocks.AccountRepository)
	service := NewService(NewConfigFromEnv(), mockedAccountRepo, mockedHashEngine)

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
