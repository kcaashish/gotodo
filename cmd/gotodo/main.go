package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kcaashish/gotodo/postgres"
	"github.com/kcaashish/gotodo/web"
)

func main() {
	store, err := postgres.NewStore("postgres://postgres:secret@localhost/postgres?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	router := web.NewServer(store)
	fmt.Println("Starting the server on: 3000...")
	log.Fatal(http.ListenAndServe(":3000", router))
}
