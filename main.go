package main

import (
	"fmt"
	"log"
	"net/http"
)

var PORT = "8080"

func main() {
	log.Println("Starting webhook")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
