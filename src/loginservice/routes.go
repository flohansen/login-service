package loginservice

import (
	"encoding/json"
	"flhansen/fitter-login-service/src/repository"
	"flhansen/fitter-login-service/src/security"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func sendSimpleResponse(w http.ResponseWriter, status int, message string) {
	sendResponse(w, status, message, map[string]interface{}{})
}

func sendResponse(w http.ResponseWriter, status int, message string, props map[string]interface{}) {
	response := map[string]interface{}{
		"status":  status,
		"message": message,
	}

	for key, value := range props {
		if _, ok := props[key]; ok {
			response[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func (service *LoginService) LoginHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var request UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		sendSimpleResponse(w, http.StatusBadRequest, "Wrong request body format.")
		return
	}

	user, err := service.accountRepo.GetAccountByUsername(request.Username)
	if err != nil {
		sendSimpleResponse(w, http.StatusUnauthorized, "Wrong user credentials.")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		sendSimpleResponse(w, http.StatusUnauthorized, "Wrong user credentials.")
		return
	}

	token, _ := security.GenerateToken(user.Id, user.Username, jwt.SigningMethodHS256, []byte(service.config.Jwt.SignKey))

	sendResponse(w, 200, "User login successful.", map[string]interface{}{
		"token": token,
	})
}

func (service *LoginService) RegisterHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var request UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		service.logger.Printf("[%s][ERROR] decode request body failed: %s", r.RemoteAddr, err.Error())
		sendSimpleResponse(w, http.StatusInternalServerError, "Invalid request body.")
		return
	}

	_, err := service.accountRepo.GetAccountByUsername(request.Username)
	if err == nil {
		service.logger.Printf("[%s][WARNING] register user '%s' failed", r.RemoteAddr, request.Username)
		sendSimpleResponse(w, http.StatusBadRequest, "User already exists.")
		return
	}

	passwordHash, err := service.hashEngine.HashPassword([]byte(request.Password))
	if err != nil {
		service.logger.Printf("[%s][ERROR] hashing password failed: %s", r.RemoteAddr, err.Error())
		sendSimpleResponse(w, http.StatusInternalServerError, "Could not register user.")
		return
	}

	id, err := service.accountRepo.CreateAccount(repository.Account{
		Username:     request.Username,
		Password:     string(passwordHash),
		Email:        request.Email,
		CreationDate: time.Now(),
	})
	if err != nil {
		service.logger.Printf("[%s][ERROR] insert user '%s' into database failed: %s", r.RemoteAddr, request.Username, err.Error())
		sendSimpleResponse(w, http.StatusInternalServerError, "Could not register user.")
		return
	}

	sendResponse(w, http.StatusOK, "User registered successfully", map[string]interface{}{
		"userId": id,
	})
}
