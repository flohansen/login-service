package loginservice

import (
	"flhansen/fitter-login-service/src/repository"
	"flhansen/fitter-login-service/src/security"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Logger interface {
	Infof(format string, v ...any)
	Warnf(format string, v ...any)
	Errorf(format string, v ...any)
}

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
	logger      Logger
}

func NewService(cfg LoginServiceConfig, accountRepo repository.AccountRepository, hashEngine security.HashEngine, logger Logger) *LoginService {
	service := &LoginService{
		handler:     httprouter.New(),
		config:      cfg,
		hashEngine:  hashEngine,
		accountRepo: accountRepo,
		logger:      logger,
	}

	service.handler.POST("/api/auth/login", service.LoginHandler)
	service.handler.POST("/api/auth/register", service.RegisterHandler)
	return service
}

func (service *LoginService) GetAddr() string {
	return fmt.Sprintf("%s:%d", service.config.Host, service.config.Port)
}

func (service *LoginService) Start() error {
	return http.ListenAndServe(service.GetAddr(), service.handler)
}
