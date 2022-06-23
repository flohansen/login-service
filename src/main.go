package main

import (
	"flhansen/fitter-login-service/src/routes"
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func main() {
	os.Exit(runApplication())
}

func runApplication() int {
	router := httprouter.New()
	router.GET("/hello", routes.HelloWorldHandler)

	host := os.Getenv("FITTER_LOGIN_SERVICE_HOST")
	port := os.Getenv("FITTER_LOGIN_SERVICE_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	server := &http.Server{Addr: addr, Handler: router}

	fmt.Printf("Starting service at %s", addr)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("An error occured while starting the server: %v", err)
		return 1
	}

	return 0
}
