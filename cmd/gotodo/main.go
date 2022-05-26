package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kcaashish/gotodo/postgres"
	"github.com/kcaashish/gotodo/web"
)

func main() {
	if er := godotenv.Load(".env"); er != nil {
		log.Fatal("Error loading .env file.")
	}
	port := os.Getenv("BACKEND_PORT")
	user := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	store, err := postgres.NewStore(fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, dbPass, host, dbPort, dbName))

	if err != nil {
		log.Fatal(err)
	}

	router := web.NewServer(store)
	fmt.Printf("Starting the server on: %v...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
