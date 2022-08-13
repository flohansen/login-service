package main

import (
	"flhansen/fitter-login-service/src/testhelper"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunApplicationShouldReturnErrorIfParsingServicePortFailed(t *testing.T) {
	t.Cleanup(testhelper.CreateTestEnvironment(map[string]string{
		"LOGIN_SERVICE_PORT":          "a",
		"LOGIN_SERVICE_DATABASE_PORT": "0",
	}))

	done := make(chan int)
	go func() {
		done <- runApplication()
	}()

	select {
	case exitCode := <-done:
		assert.Equal(t, 1, exitCode)
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Application did not terminate")
	}
}

func TestRunApplicationShouldReturnErrorIfParsingDatabasePortFailed(t *testing.T) {
	t.Cleanup(testhelper.CreateTestEnvironment(map[string]string{
		"LOGIN_SERVICE_PORT":          "0",
		"LOGIN_SERVICE_DATABASE_PORT": "a",
	}))

	done := make(chan int)
	go func() {
		done <- runApplication()
	}()

	select {
	case exitCode := <-done:
		assert.Equal(t, 1, exitCode)
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Application did not terminate")
	}
}

func TestRunApplicationShouldReturnErrorIfStartingServiceFailed(t *testing.T) {
	t.Cleanup(testhelper.CreateTestEnvironment(map[string]string{
		"LOGIN_SERVICE_PORT":          "-1",
		"LOGIN_SERVICE_DATABASE_PORT": "0",
	}))

	done := make(chan int)
	go func() {
		done <- runApplication()
	}()

	select {
	case exitCode := <-done:
		assert.Equal(t, 1, exitCode)
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Application did not terminate")
	}
}

func TestRunApplicationShouldSucceed(t *testing.T) {
	t.Cleanup(testhelper.CreateTestEnvironment(map[string]string{
		"LOGIN_SERVICE_HOST":          "localhost",
		"LOGIN_SERVICE_PORT":          "0",
		"LOGIN_SERVICE_DATABASE_PORT": "0",
	}))

	done := make(chan int)
	go func() {
		done <- runApplication()
	}()

	select {
	case <-done:
		t.Fatal("Application did terminate")
	case <-time.After(200 * time.Millisecond):
		return
	}
}
