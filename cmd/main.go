package main

import (
	"Nix_trainee_practic/config"
	"Nix_trainee_practic/config/container"
	"Nix_trainee_practic/internal/http"
	"Nix_trainee_practic/internal/repository"
	"log"
)

// @title 		NIX TRAINEE PRACTIC Demo App
// @version 	V1.echo
// @description REST service for NIX TRAINEE PRACTIC PART

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @host 		localhost:8080
// @BasePath 	/
func main() {
	var conf = config.GetConfiguration()

	err := repository.Migrate(conf)
	if err != nil {
		log.Fatalf("Unable to apply migrations: %q\n", err)
	}

	cont := container.New(conf)

	// Echo Server
	srv := http.NewServer()

	http.EchoRouter(srv, cont)

	err = srv.Start()
	if err != nil {
		log.Fatal("Port already used")
	}
}
