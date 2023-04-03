package main

import (
	"context"
	"fmt"
	admv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"k8s.io/client-go/kubernetes"

	"net/http"
)

// validate checks if the incoming object is valid, if not rejects request to enforce custom policies
func validate(ar *admv1.AdmissionReview) *admv1.AdmissionReview {

	err := isNamespaceEmpty(clientset, dynamicClient, ar.Request.Namespace)
	if err != nil {
		return reviewResponse(ar.Request.UID, false, http.StatusBadRequest, err.Error())
	}
	return reviewResponse(ar.Request.UID, true, http.StatusOK, "validation successful")
}

// isNamespaceEmpty checks if a namespace is empty
func isNamespaceEmpty(clientSet kubernetes.Interface, dynamicClient dynamic.Interface, namespaceName string) error {

	// Get a list of all resource types supported by the API server
	apiResourcesList, err := clientSet.Discovery().ServerPreferredResources()
	if err != nil {
		return err
	}

	// Create a channel to receive errors from goroutines
	errChan := make(chan error)
	count := 0
	// Iterate over each resource type and get all resources in the given namespace
	for _, apiResources := range apiResourcesList {
		for _, apiResource := range apiResources.APIResources {
			// ConfigMap, ServiceAccount and Secret has defaults when creating a namespace
			// such as kube-root-ca.crt, default-token and default respectively
			if apiResource.Kind == "ConfigMap" || apiResource.Kind == "ServiceAccount" || apiResource.Kind == "Secret" {
				continue
			}
			// Check if the resource is namespaced and has a "list" verb
			if apiResource.Namespaced && contains(apiResource.Verbs, "list") {

				// Get the API group, version, and resource kind for the resource
				groupVersion, err := schema.ParseGroupVersion(apiResources.GroupVersion)
				if err != nil {
					//errChan <- err
					return err
				}
				resource := apiResource.Name

				// Get a dynamic client for the resource type
				client := dynamicClient.Resource(
					schema.GroupVersionResource{
						Group:    groupVersion.Group,
						Version:  groupVersion.Version,
						Resource: resource,
					},
				).Namespace(namespaceName)
				count++
				go func() {
					//defer func() {
					//	fmt.Println("I am done")
					//	//wg.Done()
					//}()
					// Use the client to get a list of resources in the namespace
					fmt.Println("Calling the list api")
					objList, err := client.List(context.Background(), metav1.ListOptions{})
					if err != nil {
						errChan <- err
						return
					}
					fmt.Println("finished calling")
					// Return as soon as we find a single resource in the ns
					if len(objList.Items) > 0 {
						errChan <- fmt.Errorf("namespace has resourecs associated with it")
						return
					}
					errChan <- nil
					fmt.Println("I am Done")
				}()
			}
		}
	}
	fmt.Printf("Number of threads %d\n", count)
	// wait for the goroutines to finish

	var errReceived error
	for i := 0; i < count; i++ {
		err := <-errChan
		if err != nil {
			errReceived = err
		}
		fmt.Printf("Remaining %d", count-i)
	}

	// No resources were found in the namespace besides the defaults
	// which are created during the creation of ns
	return errReceived
}

// Helper function to check if a string is in a slice
func contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}
