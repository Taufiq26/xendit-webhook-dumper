package main

import (
	"log"
	"net/http"
	"xendit-webhook-dumper/webhooks"
)

func main() {
	http.HandleFunc("/xendit/webhook", webhooks.HandleWebhook)

	log.Println("Server starting on :6969")
	log.Fatal(http.ListenAndServe(":6969", nil))
}
