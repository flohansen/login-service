package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createTestEnvironment(variables map[string]string) func() {
	oldVariables := map[string]string{}

	for name, value := range variables {
		if oldValue, ok := os.LookupEnv(name); ok {
			oldVariables[name] = oldValue
		}

		os.Setenv(name, value)
	}

	return func() {
		for name := range variables {
			oldValue, ok := oldVariables[name]
			if ok {
				os.Setenv(name, oldValue)
			} else {
				os.Unsetenv(name)
			}
		}
	}
}

func TestRunApplication(t *testing.T) {
	t.Cleanup(createTestEnvironment(map[string]string{
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
	t.Cleanup(createTestEnvironment(map[string]string{
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
