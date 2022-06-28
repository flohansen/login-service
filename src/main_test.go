package main

import (
	"flhansen/fitter-login-service/src/testhelper"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunApplication(t *testing.T) {
	t.Cleanup(testhelper.CreateTestEnvironment(map[string]string{
		"FITTER_LOGIN_SERVICE_HOST": "0.0.0.0",
		"FITTER_LOGIN_SERVICE_PORT": "8080",
	}))

	done := make(chan int)

	go func() {
		done <- runApplication()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		return
	case exitCode := <-done:
		t.Fatalf("Application did terminate with code %d", exitCode)
	}
}

func TestRunApplicationError(t *testing.T) {
	t.Cleanup(testhelper.CreateTestEnvironment(map[string]string{
		"FITTER_LOGIN_SERVICE_HOST": "0.0.0.0",
		"FITTER_LOGIN_SERVICE_PORT": "-1",
	}))

	done := make(chan int)

	go func() {
		done <- runApplication()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("Application did not terminate")
	case exitCode := <-done:
		assert.Equal(t, 1, exitCode)
	}
}

func TestRunApplicationInvalidEnvVars(t *testing.T) {
	t.Cleanup(testhelper.CreateTestEnvironment(map[string]string{
		"FITTER_LOGIN_SERVICE_HOST":    "0.0.0.0",
		"FITTER_LOGIN_SERVICE_PORT":    "0",
		"FITTER_LOGIN_SERVICE_DB_PORT": "-1",
	}))

	done := make(chan int)

	go func() {
		done <- runApplication()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("Application did not terminate")
	case exitCode := <-done:
		assert.Equal(t, 1, exitCode)
	}
}
