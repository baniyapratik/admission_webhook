package main

import (
	"encoding/json"

	"github.com/google/uuid"
)

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func createPatches() ([]byte, error) {
	var patches []patchOperation
	patches = append(patches, createUUIDPatch())
	return json.Marshal(patches)
}

// createUUIDPatch generates uuid and creates a json patch
func createUUIDPatch() patchOperation {
	// create the uuid and patch it in the admissionReview object
	var id string
	id = uuid.New().String()
	return patchOperation{
		Op:    "replace",
		Path:  "/spec/id",
		Value: id,
	}
}
