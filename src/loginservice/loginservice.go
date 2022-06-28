package loginservice

import (
	"flhansen/fitter-login-service/src/database"
	"flhansen/fitter-login-service/src/security"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type LoginServiceConfig struct {
	Host     string
	Port     int
	Database database.DatabaseConfig
	Jwt      security.JwtConfig
}

type LoginService struct {
	Handler *httprouter.Router
	Config  LoginServiceConfig
	Db      database.SqlDatabase
}

func New(cfg LoginServiceConfig) (*LoginService, error) {
	db, err := database.New(cfg.Database)
	if err != nil {
		return nil, err
	}

	return NewService(cfg, db), nil
}

func NewService(cfg LoginServiceConfig, db database.SqlDatabase) *LoginService {
	service := &LoginService{
		Handler: httprouter.New(),
		Config:  cfg,
		Db:      db,
	}

	service.Handler.POST("/api/auth/login", service.LoginHandler)
	service.Handler.POST("/api/auth/register", service.RegisterHandler)
	return service
}

func NewConfigFromEnv() LoginServiceConfig {
	dbPort, _ := strconv.Atoi(os.Getenv("FITTER_LOGIN_SERVICE_DB_PORT"))
	port, _ := strconv.Atoi(os.Getenv("FITTER_LOGIN_SERVICE_PORT"))

	return LoginServiceConfig{
		Host: os.Getenv("FITTER_LOGIN_SERVICE_HOST"),
		Port: port,
		Database: database.DatabaseConfig{
			Host:   os.Getenv("FITTER_LOGIN_SERVICE_DB_HOST"),
			Port:   dbPort,
			Name:   os.Getenv("FITTER_LOGIN_SERVICE_DB_NAME"),
			User:   os.Getenv("FITTER_LOGIN_SERVICE_DB_USER"),
			Region: os.Getenv("AWS_REGION"),
		},
		Jwt: security.JwtConfig{
			SignKey: os.Getenv("FITTER_LOGIN_SERVICE_JWT_SIGNKEY"),
		},
	}
}

func (service *LoginService) GetAddr() string {
	return fmt.Sprintf("%s:%d", service.Config.Host, service.Config.Port)
}

func (service *LoginService) Start() error {
	return http.ListenAndServe(service.GetAddr(), service.Handler)
}

func (service *LoginService) Server() *http.Server {
	server := &http.Server{Addr: service.GetAddr(), Handler: service.Handler}
	return server
}
