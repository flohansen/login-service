package main

import (
	"flhansen/fitter-login-service/src/loginservice"
	"fmt"
	"os"
)

func main() {
	os.Exit(runApplication())
}

func runApplication() int {
	cfg := loginservice.NewConfigFromEnv()
	service, err := loginservice.New(cfg)
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
