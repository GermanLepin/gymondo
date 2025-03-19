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
	"os"
)

const webPort = "80"

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading .env file")
	}
}

func main() {
	conn, err := connection.StartDB()
	if err != nil {
		log.Println("couldn't start db")
		os.Exit(1)
	}

	productRepository := repository.NewProductRepository(conn)
	//subscriptionRepository := repository.NewSubscriptionRepository(conn)
	//voucherRepository := repository.NewVoucherRepository(conn)

	// create services
	productService := service.New(productRepository)

	// add services to api route
	apiRoutes := rest.New(productService)

	log.Printf("starting balance service on port %s\n", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: apiRoutes.NewRoutes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
