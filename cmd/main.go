package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gymondo/db/postgres/connection"
	"gymondo/internal/api/rest"
	"gymondo/internal/repository"
	"gymondo/internal/service"
	"log"
	"net/http"

	_ "gymondo/cmd/docs"
)

const serverPort = "80"

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	conn, err := connection.StartDB()
	if err != nil {
		log.Fatalf("Could not start DB connection: %v", err)
	}
	defer conn.Close()

	repo := repository.New(conn)
	serv := service.New(repo)

	apiRoutes := rest.New(serv)
	log.Printf("Starting balance service on port %s\n", serverPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", serverPort),
		Handler: apiRoutes.NewRoutes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
