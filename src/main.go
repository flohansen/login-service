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
	username := os.Getenv("FITTER_LOGIN_SERVICE_DB_USERNAME")
	region := os.Getenv("AWS_REGION")

	cfg := loginservice.NewConfigFromEnv()
	credentialsProvider := security.NewAwsCredentialsProvider(cfg.Database.Host, cfg.Database.Port, username, region)

	db, err := database.NewPostgresDatabase(cfg.Database, credentialsProvider)
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
