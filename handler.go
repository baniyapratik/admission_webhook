package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func MutateResourceIDHandler(w http.ResponseWriter, r *http.Request) {
	// if the request is invalid return bad request
	ar, err := parseRequest(r)
	if err != nil {
		http.Error(w, "parse request error: "+err.Error(), http.StatusBadRequest)
		return
	}
	// admissionReview response
	admissionReview := mutateID(ar)
	// ensure we can marshal the admissionReview before sending it back
	respBytes, err := json.Marshal(admissionReview)
	log.Println("AdmissionReview: " + string(respBytes))
	if err != nil {
		log.Println(fmt.Sprintf("error %s, while converting response to byte slice", err.Error()))
		http.Error(w, "admissionReview marshal error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Ready to write response")
	if _, err := w.Write(respBytes); err != nil {
		log.Println(fmt.Sprintf("Can't write response: %v", err))
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}
