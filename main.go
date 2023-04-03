package main

import (
	"fmt"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	PORT          = "8080"
	clientset     *kubernetes.Clientset
	dynamicClient dynamic.Interface
)

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

func init() {
	var err error
	// Load the Kubernetes configuration from a file
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Panic(err)
	}

	// Create a Kubernetes clientset using the configuration
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err)
	}
	// Create a dynamic client for interacting with the API server
	dynamicClient, err = dynamic.NewForConfig(config)
	if err != nil {
		log.Panic(err)
	}
}
