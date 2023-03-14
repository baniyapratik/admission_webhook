package main

import (
	admv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"log"
	"net/http"
)

var (
	scheme       = runtime.NewScheme()
	codecs       = serializer.NewCodecFactory(scheme)
	deserializer = codecs.UniversalDeserializer()
)

func mutateID(ar *admv1.AdmissionReview) *admv1.AdmissionReview {
	// create the patch
	patchBytes, err := createPatches()
	if err != nil {
		errMsg := "error getting patches. error :" + err.Error()
		return reviewResponse(ar.Request.UID, false, http.StatusBadRequest, errMsg)
	}
	log.Println("AdmissionResponse: patch=" + string(patchBytes))

	return patchReviewResponse(ar.Request.UID, patchBytes)
}
