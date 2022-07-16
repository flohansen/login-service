package main

import (
	"flhansen/fitter-login-service/src/loginservice"
	"flhansen/fitter-login-service/src/repository"
	"flhansen/fitter-login-service/src/security"
	"fmt"
	"os"
	"strconv"
)

func main() {
	os.Exit(runApplication())
}

func createConfigFromEnvironment() (loginservice.LoginServiceConfig, repository.DatabaseConfig, error) {
	host := os.Getenv("LOGIN_SERVICE_HOST")
	port := os.Getenv("LOGIN_SERVICE_PORT")
	jwtSignKey := os.Getenv("LOGIN_SERVICE_JWT_SIGN_KEY")
	databaseHost := os.Getenv("LOGIN_SERVICE_DATABASE_HOST")
	databasePort := os.Getenv("LOGIN_SERVICE_DATABASE_PORT")
	databaseUser := os.Getenv("LOGIN_SERVICE_DATABASE_USER")
	databasePass := os.Getenv("LOGIN_SERVICE_DATABASE_PASSWORD")
	databaseName := os.Getenv("LOGIN_SERVICE_DATABASE_NAME")

	var serviceConfig loginservice.LoginServiceConfig
	var databaseConfig repository.DatabaseConfig

	portValue, err := strconv.Atoi(port)
	if err != nil {
		return serviceConfig, databaseConfig, err
	}

	databasePortValue, err := strconv.Atoi(databasePort)
	if err != nil {
		return serviceConfig, databaseConfig, err
	}

	serviceConfig = loginservice.LoginServiceConfig{
		Host: host,
		Port: portValue,
		Jwt: security.JwtConfig{
			SignKey: jwtSignKey,
		},
	}

	databaseConfig = repository.DatabaseConfig{
		Host:         databaseHost,
		Port:         databasePortValue,
		Username:     databaseUser,
		Password:     databasePass,
		DatabaseName: databaseName,
	}

	return serviceConfig, databaseConfig, nil
}

func runApplication() int {
	serviceConfig, databaseConfig, err := createConfigFromEnvironment()
	if err != nil {
		fmt.Printf("An error occured while creating configuration: %v", err)
		return 1
	}

	hashEngine := security.NewBcryptEngine()
	accountRepo := repository.NewAccountRepository(databaseConfig)
	service := loginservice.NewService(serviceConfig, accountRepo, hashEngine)

	fmt.Printf("Starting service at %s", service.GetAddr())
	if err := service.Start(); err != nil {
		fmt.Printf("An error occured while starting the service: %v", err)
		return 1
	}

	return 0
}
