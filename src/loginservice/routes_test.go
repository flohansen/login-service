package loginservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"flhansen/fitter-login-service/src/mocks"
	"flhansen/fitter-login-service/src/repository"
	"flhansen/fitter-login-service/src/security"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginHandlerSuccess(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), 8)
	mockedHashEngine := new(mocks.HashEngine)
	mockedHashEngine.
		On("HashPassword", mock.Anything).
		Return(hashedPassword, nil)
	mockedAccountRepo := new(mocks.AccountRepository)
	mockedAccountRepo.
		On("GetAccountByUsername", mock.Anything).
		Return(repository.Account{Password: string(hashedPassword)}, nil)
	mockedLogger := new(mocks.Logger)
	service := NewService(LoginServiceConfig{
		Jwt: security.JwtConfig{
			SignKey: "secret",
		},
	}, mockedAccountRepo, mockedHashEngine, mockedLogger)

	requestBody, err := json.Marshal(UserLoginRequest{
		Username: "testuser",
		Password: "testpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	request, err := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	service.handler.ServeHTTP(recorder, request)

	var resp map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.NotNil(t, resp["status"])
	assert.NotNil(t, resp["message"])
	assert.NotNil(t, resp["token"])
}

func TestLoginHandlerWrongCredentials(t *testing.T) {
	mockedHashEngine := new(mocks.HashEngine)
	mockedAccountRepo := new(mocks.AccountRepository)
	mockedAccountRepo.
		On("GetAccountByUsername", mock.Anything).
		Return(repository.Account{}, nil)
	mockedLogger := new(mocks.Logger)
	mockedLogger.
		On("Errorf", mock.Anything, mock.Anything, mock.Anything)
	service := NewService(LoginServiceConfig{}, mockedAccountRepo, mockedHashEngine, mockedLogger)

	requestBody, err := json.Marshal(UserLoginRequest{
		Username: "testuser",
		Password: "wrongtestpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	request, err := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	service.handler.ServeHTTP(recorder, request)

	var resp map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.NotNil(t, resp["status"])
	assert.NotNil(t, resp["message"])
	assert.Nil(t, resp["token"])
	mockedLogger.AssertCalled(t, "Errorf", "(%s) wrong password: %s", mock.Anything, mock.Anything)
}

func TestLoginHandlerUserNotExist(t *testing.T) {
	mockedHashEngine := new(mocks.HashEngine)
	mockedAccountRepo := new(mocks.AccountRepository)
	mockedAccountRepo.On("GetAccountByUsername", mock.Anything).Return(repository.Account{}, errors.New("err"))
	mockedLogger := new(mocks.Logger)
	mockedLogger.
		On("Warnf", mock.Anything, mock.Anything, mock.Anything)
	service := NewService(LoginServiceConfig{}, mockedAccountRepo, mockedHashEngine, mockedLogger)

	requestBody, err := json.Marshal(UserLoginRequest{
		Username: "testuserxyz",
		Password: "testpass",
	})
	if err != nil {
		t.Fatal(err)
	}

	request, err := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	service.handler.ServeHTTP(recorder, request)

	var resp map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.NotNil(t, resp["status"])
	assert.NotNil(t, resp["message"])
	assert.Nil(t, resp["token"])
	mockedLogger.AssertCalled(t, "Warnf", "(%s) login of user '%s' failed", mock.Anything, mock.Anything)
}

func TestLoginHandlerBadRequest(t *testing.T) {
	mockedHashEngine := new(mocks.HashEngine)
	mockedAccountRepo := new(mocks.AccountRepository)
	mockedLogger := new(mocks.Logger)
	mockedLogger.
		On("Errorf", mock.Anything, mock.Anything, mock.Anything)
	service := NewService(LoginServiceConfig{}, mockedAccountRepo, mockedHashEngine, mockedLogger)

	requestBody := []byte("{ username: ")
	request, err := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	service.handler.ServeHTTP(recorder, request)

	var resp map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.NotNil(t, resp["status"])
	assert.NotNil(t, resp["message"])
	assert.Nil(t, resp["token"])
	mockedLogger.AssertCalled(t, "Errorf", "(%s) decode request body failed: %s", mock.Anything, mock.Anything)
}

func TestRegisterHandlerShouldReturnErrorIfInvalidJsonBody(t *testing.T) {
	// given
	invalidBody := []byte(`{ username: " }`)
	mockedHashEngine := new(mocks.HashEngine)
	mockedAccountRepo := new(mocks.AccountRepository)
	mockedLogger := new(mocks.Logger)
	mockedLogger.
		On("Errorf", mock.Anything, mock.Anything, mock.Anything)
	service := NewService(LoginServiceConfig{}, mockedAccountRepo, mockedHashEngine, mockedLogger)

	// when
	request, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(invalidBody))
	responseWriter := httptest.NewRecorder()
	service.handler.ServeHTTP(responseWriter, request)

	var response map[string]interface{}
	err := json.NewDecoder(responseWriter.Body).Decode(&response)

	// then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, responseWriter.Code)
	assert.NotNil(t, response["status"])
	assert.NotNil(t, response["message"])
	mockedLogger.AssertCalled(t, "Errorf", "(%s) decode request body failed: %s", mock.Anything, mock.Anything)
}

func TestRegisterHandlerShouldReturnErrorIfUserAlreadyExists(t *testing.T) {
	// given
	body := []byte(`{ "username": "testuser", "password": "testpass", "email": "testmail@test.com" }`)
	mockedHashEngine := new(mocks.HashEngine)
	mockedAccountRepo := new(mocks.AccountRepository)
	mockedAccountRepo.
		On("GetAccountByUsername", mock.Anything).
		Return(repository.Account{}, nil)
	mockedLogger := new(mocks.Logger)
	mockedLogger.
		On("Warnf", mock.Anything, mock.Anything, mock.Anything)
	service := NewService(LoginServiceConfig{}, mockedAccountRepo, mockedHashEngine, mockedLogger)

	// when
	request, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	responseWriter := httptest.NewRecorder()
	service.handler.ServeHTTP(responseWriter, request)

	var response map[string]interface{}
	err := json.NewDecoder(responseWriter.Body).Decode(&response)

	// then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, responseWriter.Code)
	assert.NotNil(t, response["status"])
	assert.NotNil(t, response["message"])
	mockedLogger.AssertCalled(t, "Warnf", "(%s) register user '%s' failed", mock.Anything, mock.Anything)
}

func TestRegisterHandlerShouldReturnErrorWhenCreatingUserFailed(t *testing.T) {
	// given
	body := []byte(`{ "username": "testuser", "password": "testpass", "email": "testmail@test.com" }`)
	mockedHashEngine := new(mocks.HashEngine)
	mockedHashEngine.
		On("HashPassword", mock.Anything).
		Return([]byte{}, nil).
		Once()
	mockedAccountRepo := new(mocks.AccountRepository)
	mockedAccountRepo.
		On("GetAccountByUsername", mock.Anything).
		Return(repository.Account{}, errors.New("user not found")).
		On("CreateAccount", mock.Anything).
		Return(-1, errors.New("could not create user"))
	mockedLogger := new(mocks.Logger)
	mockedLogger.
		On("Errorf", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	service := NewService(LoginServiceConfig{}, mockedAccountRepo, mockedHashEngine, mockedLogger)

	// when
	request, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	responseWriter := httptest.NewRecorder()
	service.handler.ServeHTTP(responseWriter, request)

	var response map[string]interface{}
	err := json.NewDecoder(responseWriter.Body).Decode(&response)

	// then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, responseWriter.Code)
	assert.NotNil(t, response["status"])
	assert.NotNil(t, response["message"])
	mockedLogger.AssertCalled(t, "Errorf", "(%s) insert user '%s' into database failed: %s", mock.Anything, mock.Anything, mock.Anything)
}

func TestRegisterHandlerShouldReturnErrorIfPasswordCouldNotBeHashed(t *testing.T) {
	// given
	body := []byte(`{ "username": "testuser", "password": "testpass", "email": "testmail@test.com" }`)
	mockedHashEngine := new(mocks.HashEngine)
	mockedHashEngine.
		On("HashPassword", mock.Anything).
		Return([]byte{}, errors.New("error while hashing password")).
		Once()
	mockedAccountRepo := new(mocks.AccountRepository)
	mockedAccountRepo.
		On("GetAccountByUsername", mock.Anything).
		Return(repository.Account{}, errors.New("user not found"))
	mockedLogger := new(mocks.Logger)
	mockedLogger.
		On("Errorf", mock.Anything, mock.Anything, mock.Anything)
	service := NewService(LoginServiceConfig{}, mockedAccountRepo, mockedHashEngine, mockedLogger)

	// when
	request, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	responseWriter := httptest.NewRecorder()
	service.handler.ServeHTTP(responseWriter, request)

	var response map[string]interface{}
	err := json.NewDecoder(responseWriter.Body).Decode(&response)

	// then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, responseWriter.Code)
	assert.NotNil(t, response["status"])
	assert.NotNil(t, response["message"])
	mockedLogger.AssertCalled(t, "Errorf", "(%s) hashing password failed: %s", mock.Anything, mock.Anything)
}

func TestRegisterHandlerSucceeded(t *testing.T) {
	// given
	body := []byte(`{ "username": "testuser", "password": "testpass", "email": "testmail@test.com" }`)
	mockedHashEngine := new(mocks.HashEngine)
	mockedHashEngine.
		On("HashPassword", mock.Anything).
		Return([]byte{}, nil)
	mockedAccountRepo := new(mocks.AccountRepository)
	mockedAccountRepo.
		On("GetAccountByUsername", mock.Anything).
		Return(repository.Account{}, errors.New("user not found")).
		On("CreateAccount", mock.Anything).
		Return(1, nil)
	mockedLogger := new(mocks.Logger)
	service := NewService(LoginServiceConfig{}, mockedAccountRepo, mockedHashEngine, mockedLogger)

	// when
	request, _ := http.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	responseWriter := httptest.NewRecorder()
	service.handler.ServeHTTP(responseWriter, request)

	var response map[string]interface{}
	err := json.NewDecoder(responseWriter.Body).Decode(&response)

	// then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, responseWriter.Code)
	assert.NotNil(t, response["status"])
	assert.NotNil(t, response["message"])
	assert.NotNil(t, response["userId"])
	mockedLogger.AssertNotCalled(t, "Printf")
}
