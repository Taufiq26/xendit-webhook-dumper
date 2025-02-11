package main

import (
	"log"
	"net/http"
	"os"
	"xendit-webhook-dumper/webhooks"
)

func main() {
	// Get port from environment variable, default to 8080 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/xendit/webhook", webhooks.HandleWebhook)

	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
