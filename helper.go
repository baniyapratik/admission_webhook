package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	admv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"log"
	"net/http"
)

func reviewResponse(uid types.UID, allowed bool, httpCode int32,
	reason string) *admv1.AdmissionReview {
	return &admv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		// admission review response
		Response: &admv1.AdmissionResponse{
			UID:     uid,
			Allowed: allowed,
			Result: &metav1.Status{
				Code:    httpCode,
				Message: reason,
			},
		},
	}
}

// admissionError wraps error as AdmissionResponse
func admissionError(msg string, code int32) *admv1.AdmissionResponse {
	return &admv1.AdmissionResponse{
		Allowed: false,
		Result: &metav1.Status{
			Message: msg,
			Code:    code,
		},
	}
}

// patchReviewResponse builds an admission review with given json patch
func patchReviewResponse(uid types.UID, patch []byte) *admv1.AdmissionReview {
	patchType := admv1.PatchTypeJSONPatch
	return &admv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admv1.AdmissionResponse{
			UID:       uid,
			Allowed:   true,
			PatchType: &patchType,
			Patch:     patch,
		},
	}
}

// parseRequest does sanity check on the request body to ensure its valid
// and extracts the AdmissionReview if possible
func parseRequest(r *http.Request) (*admv1.AdmissionReview, error) {
	// read and populate request body
	var body []byte
	if r.Body != nil {
		// read the request body
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	// if body is empty
	if len(body) == 0 {
		errMsg := "request body is empty"
		log.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// verify the content type is of application/json
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		errMsg := fmt.Sprintf("Content-Type=%s, expect application/json", contentType)
		log.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// ensure the request body is AdmissionReview
	ar := admv1.AdmissionReview{}
	if err := json.Unmarshal(body, &ar); err != nil {
		errMsg := "could not parse admission review request: " + err.Error()
		log.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	// check if the admissionRequest exists
	if ar.Request == nil {
		errMsg := "admission review can't be used: Request field is nil"
		log.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	// return ar if request is valid
	return &ar, nil
}
