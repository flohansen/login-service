package main

import (
	"flhansen/fitter-login-service/src/loginservice"
	"flhansen/fitter-login-service/src/repository"
	"flhansen/fitter-login-service/src/security"
	"fmt"
	"os"
)

func main() {
	os.Exit(runApplication())
}

func runApplication() int {
	cfg := loginservice.NewConfigFromEnv()
	hashEngine := security.NewBcryptEngine()
	accountRepo, err := repository.NewAccountRepository(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Database)
	if err != nil {
		fmt.Printf("An error occured while connecting to the database: %v", err)
		return 1
	}

	service, err := loginservice.NewService(cfg, accountRepo, hashEngine)
	if err != nil {
		fmt.Printf("An error occured while creating the service: %v", err)
		return 1
	}

	fmt.Printf("Starting service at %s", service.GetAddr())
	if err := service.Start(); err != nil {
		fmt.Printf("An error occured while starting the service: %v", err)
		return 1
	}

	return 0
}
