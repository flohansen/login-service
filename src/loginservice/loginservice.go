package loginservice

import (
	"flhansen/fitter-login-service/src/repository"
	"flhansen/fitter-login-service/src/security"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type LoginServiceConfig struct {
	Host string
	Port int
	Jwt  security.JwtConfig
}

type LoginService struct {
	handler     *httprouter.Router
	config      LoginServiceConfig
	accountRepo repository.AccountRepository
	hashEngine  security.HashEngine
}

func NewService(cfg LoginServiceConfig, accountRepo repository.AccountRepository, hashEngine security.HashEngine) *LoginService {
	service := &LoginService{
		handler:     httprouter.New(),
		config:      cfg,
		hashEngine:  hashEngine,
		accountRepo: accountRepo,
	}

	service.handler.POST("/api/auth/login", service.LoginHandler)
	service.handler.POST("/api/auth/register", service.RegisterHandler)
	return service
}

func NewConfigFromEnv() LoginServiceConfig {
	port, _ := strconv.Atoi(os.Getenv("LOGIN_SERVICE_PORT"))

	return LoginServiceConfig{
		Host: os.Getenv("LOGIN_SERVICE_HOST"),
		Port: port,
		Jwt: security.JwtConfig{
			SignKey: os.Getenv("LOGIN_SERVICE_JWT_SIGNKEY"),
		},
	}
}

func (service *LoginService) GetAddr() string {
	return fmt.Sprintf("%s:%d", service.config.Host, service.config.Port)
}

func (service *LoginService) Start() error {
	return http.ListenAndServe(service.GetAddr(), service.handler)
}
