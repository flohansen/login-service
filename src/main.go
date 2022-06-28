package main

import (
	"flhansen/fitter-login-service/src/database"
	"flhansen/fitter-login-service/src/loginservice"
	"flhansen/fitter-login-service/src/security"
	"fmt"
	"os"
)

func main() {
	os.Exit(runApplication())
}

func runApplication() int {
	cfg := loginservice.NewConfigFromEnv()
	db, err := database.New(cfg.Database)
	if err != nil {
		fmt.Printf("An error occured while creating the database connection: %v", err)
		return 1
	}
	hashEngine := security.NewBcryptEngine()

	service := loginservice.NewService(cfg, db, hashEngine)
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
