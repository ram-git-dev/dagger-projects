// File: services/service-b/main.go

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Set up a single HTTP route at "/"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Write a simple response
		fmt.Fprintln(w, "Hello from Service B")
	})

	// Log that the service has started
	log.Println("Service B listening on port 8082")

	// Start the HTTP server
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatalf("Error starting Service B: %v", err)
	}
}
