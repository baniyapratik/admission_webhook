package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/admission/v1"
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
	var admissionReview *v1.AdmissionReview
	switch string(ar.Request.Operation) {
	case "CREATE":
		admissionReview = mutateID(ar)
	default:
		errMsg := "operation not supported by the webhook yet"
		admissionReview = reviewResponse(ar.Request.UID, false, http.StatusBadRequest, errMsg)
	}

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

func ValidateDeleteResource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logr.FromContextOrDiscard(ctx)
	// if the request is invalid return bad request
	ar, err := parseRequest(r)

	if err != nil {
		log.Error(err, "error while parsing the request")
		http.Error(w, "bad request, see server logs for details", http.StatusBadRequest)
		return
	}
	var admissionReview *v1.AdmissionReview

	switch string(ar.Request.Operation) {
	case "DELETE":
		admissionReview = validate(ar)
	default:
		errMsg := "operation not supported by the webhook yet"
		admissionReview = reviewResponse(ar.Request.UID, false, http.StatusBadRequest, errMsg)
	}
	// ensure we can marshal the admissionReview before sending it back
	respBytes, err := json.Marshal(admissionReview)
	// AdmissionReview: string(respBytes)
	if err != nil {
		log.Error(err, "admissionReview marshal error, error while converting response to byte slice")
		http.Error(w, "internal error, see server logs for details", http.StatusInternalServerError)
		return
	}
	// Ready to write response
	if _, err := w.Write(respBytes); err != nil {
		log.Error(err, "Unable to write response")
		http.Error(w, "internal error, see server logs for details", http.StatusInternalServerError)
	}
}
