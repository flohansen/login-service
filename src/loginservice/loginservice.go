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
	handler    *httprouter.Router
	config     LoginServiceConfig
	db         database.SqlDatabase
	hashEngine security.HashEngine
}

func NewService(cfg LoginServiceConfig, db database.SqlDatabase, hashEngine security.HashEngine) *LoginService {
	service := &LoginService{
		handler:    httprouter.New(),
		config:     cfg,
		db:         db,
		hashEngine: hashEngine,
	}

	service.handler.POST("/api/auth/login", service.LoginHandler)
	service.handler.POST("/api/auth/register", service.RegisterHandler)
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
	return fmt.Sprintf("%s:%d", service.config.Host, service.config.Port)
}

func (service *LoginService) Start() error {
	return http.ListenAndServe(service.GetAddr(), service.handler)
}

func (service *LoginService) Server() *http.Server {
	server := &http.Server{Addr: service.GetAddr(), Handler: service.handler}
	return server
}
