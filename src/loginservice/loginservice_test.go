package loginservice

import (
	"flhansen/fitter-login-service/src/mocks"
	"flhansen/fitter-login-service/src/security"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetAddr(t *testing.T) {
	config := LoginServiceConfig{
		Host: "0.0.0.0",
		Port: 8000,
		Jwt: security.JwtConfig{
			SignKey: "testkey",
		},
	}

	mockedHashEngine := new(mocks.HashEngine)
	mockedAccountRepo := new(mocks.AccountRepository)
	service := NewService(config, mockedAccountRepo, mockedHashEngine, log.Default())
	addr := service.GetAddr()

	assert.Equal(t, "0.0.0.0:8000", addr)
}

func TestStart(t *testing.T) {
	config := LoginServiceConfig{
		Host: "0.0.0.0",
		Port: 8001,
		Jwt: security.JwtConfig{
			SignKey: "testkey",
		},
	}

	mockedHashEngine := new(mocks.HashEngine)
	mockedAccountRepo := new(mocks.AccountRepository)
	service := NewService(config, mockedAccountRepo, mockedHashEngine, log.Default())

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
