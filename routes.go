package main

import (
	"net/http"
)

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/mutate", http.HandlerFunc(MutateResourceIDHandler))
	mux.Handle("/validate", http.HandlerFunc(ValidateDeleteResource))
	return mux
}
