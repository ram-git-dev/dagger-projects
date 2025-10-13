// File: services/service-a/main.go

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
		fmt.Fprintln(w, "Hello from Service A")
	})

	// Print message to terminal so we know it's running
	log.Println("Service A listening on port 8081")

	// Start the HTTP server
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatalf("Error starting Service A: %v", err)
	}
}
