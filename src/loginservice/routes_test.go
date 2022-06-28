package loginservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"flhansen/fitter-login-service/src/database"
	"flhansen/fitter-login-service/src/security"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type DatabaseMock struct {
	mock.Mock
}

func (m *DatabaseMock) CreateUser(user database.DbUser) (int, error) {
	args := m.Called(user)
	return args.Int(0), args.Error(1)
}

func (m *DatabaseMock) DeleteUserById(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *DatabaseMock) GetUserById(id int) (database.DbUser, error) {
	args := m.Called(id)
	return args.Get(0).(database.DbUser), args.Error(1)
}

func (m *DatabaseMock) GetUserByUsername(username string) (database.DbUser, error) {
	args := m.Called(username)
	return args.Get(0).(database.DbUser), args.Error(1)
}

type HashEngineMock struct {
	mock.Mock
}

func (m *HashEngineMock) HashPassword(password []byte) ([]byte, error) {
	args := m.Called(password)
	return args.Get(0).([]byte), args.Error(1)
}

func TestLoginHandlerSuccess(t *testing.T) {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("testpass"), 8)

	mockedDatabase := new(DatabaseMock)
	mockedDatabase.
		On("GetUserByUsername", "testuser").
		Return(database.DbUser{Username: "testuser", Password: string(passwordHash)}, nil).
		Once()
	mockedHashEngine := new(HashEngineMock)
	service := NewService(LoginServiceConfig{
		Jwt: security.JwtConfig{
			SignKey: "secret",
		},
	}, mockedDatabase, mockedHashEngine)

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
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("testpass"), 8)

	mockedDatabase := new(DatabaseMock)
	mockedDatabase.
		On("GetUserByUsername", "testuser").
		Return(database.DbUser{Username: "testuser", Password: string(passwordHash)}, nil).
		Once()
	mockedHashEngine := new(HashEngineMock)
	service := NewService(LoginServiceConfig{}, mockedDatabase, mockedHashEngine)

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
}

func TestLoginHandlerUserNotExist(t *testing.T) {
	mockedDatabase := new(DatabaseMock)
	mockedDatabase.
		On("GetUserByUsername", mock.Anything).
		Return(database.DbUser{}, errors.New("user does not exist")).
		Once()
	mockedHashEngine := new(HashEngineMock)
	service := NewService(LoginServiceConfig{}, mockedDatabase, mockedHashEngine)

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
}

func TestLoginHandlerBadRequest(t *testing.T) {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("testpass"), 8)

	mockedDatabase := new(DatabaseMock)
	mockedDatabase.
		On("GetUserByUsername", "testuser").
		Return(database.DbUser{Username: "testuser", Password: string(passwordHash)}, nil).
		Once()
	mockedHashEngine := new(HashEngineMock)
	service := NewService(LoginServiceConfig{}, mockedDatabase, mockedHashEngine)

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
}

func TestRegisterHandlerShouldReturnErrorIfInvalidJsonBody(t *testing.T) {
	// given
	invalidBody := []byte(`{ username: " }`)
	mockedDatabase := new(DatabaseMock)
	mockedHashEngine := new(HashEngineMock)
	service := NewService(LoginServiceConfig{}, mockedDatabase, mockedHashEngine)

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
}

func TestRegisterHandlerShouldReturnErrorIfUserAlreadyExists(t *testing.T) {
	// given
	mockedDatabase := new(DatabaseMock)
	mockedDatabase.
		On("GetUserByUsername", mock.Anything).
		Return(database.DbUser{}, errors.New("user already exists")).
		Once()
	body := []byte(`{ "username": "testuser", "password": "testpass", "email": "testmail@test.com" }`)
	mockedHashEngine := new(HashEngineMock)
	service := NewService(LoginServiceConfig{}, mockedDatabase, mockedHashEngine)

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
}

func TestRegisterHandlerShouldReturnErrorWhenCreatingUserFailed(t *testing.T) {
	// given
	mockedDatabase := new(DatabaseMock)
	mockedDatabase.
		On("GetUserByUsername", mock.Anything).
		Return(database.DbUser{}, nil).
		Once().
		On("CreateUser", mock.Anything).
		Return(-1, errors.New("could not create user")).
		Once()
	body := []byte(`{ "username": "testuser", "password": "testpass", "email": "testmail@test.com" }`)
	mockedHashEngine := new(HashEngineMock)
	mockedHashEngine.
		On("HashPassword", mock.Anything).
		Return([]byte{}, nil).
		Once()
	service := NewService(LoginServiceConfig{}, mockedDatabase, mockedHashEngine)

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
}

func TestRegisterHandlerShouldReturnErrorIfPasswordCouldNotBeHashed(t *testing.T) {
	// given
	mockedDatabase := new(DatabaseMock)
	mockedDatabase.
		On("GetUserByUsername", mock.Anything).
		Return(database.DbUser{}, nil).
		Once().
		On("CreateUser", mock.Anything).
		Return(0, nil).
		Once()
	body := []byte(`{ "username": "testuser", "password": "testpass", "email": "testmail@test.com" }`)
	mockedHashEngine := new(HashEngineMock)
	mockedHashEngine.
		On("HashPassword", mock.Anything).
		Return([]byte{}, errors.New("error while hashing password")).
		Once()
	service := NewService(LoginServiceConfig{}, mockedDatabase, mockedHashEngine)

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
}

func TestRegisterHandlerSucceeded(t *testing.T) {
	// given
	mockedDatabase := new(DatabaseMock)
	mockedDatabase.
		On("GetUserByUsername", mock.Anything).
		Return(database.DbUser{}, nil).
		Once().
		On("CreateUser", mock.Anything).
		Return(0, nil).
		Once()
	body := []byte(`{ "username": "testuser", "password": "testpass", "email": "testmail@test.com" }`)
	mockedHashEngine := new(HashEngineMock)
	mockedHashEngine.
		On("HashPassword", mock.Anything).
		Return([]byte{}, nil)
	service := NewService(LoginServiceConfig{}, mockedDatabase, mockedHashEngine)

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
}
