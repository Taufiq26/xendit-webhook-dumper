package webhooks

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type XenditWebhook struct {
	ID                     string `json:"id"`
	ExternalID             string `json:"external_id"`
	UserID                 string `json:"user_id"`
	IsHigh                 bool   `json:"is_high"`
	PaymentMethod          string `json:"payment_method"`
	Status                 string `json:"status"`
	MerchantName           string `json:"merchant_name"`
	Amount                 int    `json:"amount"`
	PaidAmount             int    `json:"paid_amount"`
	BankCode               string `json:"bank_code"`
	PaidAt                 string `json:"paid_at"`
	PayerEmail             string `json:"payer_email"`
	Description            string `json:"description"`
	AdjustedReceivedAmount int    `json:"adjusted_received_amount"`
	FeesPaidAmount         int    `json:"fees_paid_amount"`
	Updated                string `json:"updated"`
	Created                string `json:"created"`
	Currency               string `json:"currency"`
	PaymentChannel         string `json:"payment_channel"`
	PaymentDestination     string `json:"payment_destination"`
}

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Log the raw payload for debugging
	log.Printf("Received webhook payload: %s", string(body))

	// Parse the webhook payload
	var webhook XenditWebhook
	if err := json.Unmarshal(body, &webhook); err != nil {
		log.Printf("Error parsing webhook payload: %v", err)
		http.Error(w, "Error parsing webhook payload", http.StatusBadRequest)
		return
	}

	// Determine the date for file organization
	paidTime, err := time.Parse(time.RFC3339, webhook.PaidAt)
	if err != nil {
		log.Printf("Error parsing paid_at time: %v. Using current time.", err)
		paidTime = time.Now()
	}
	dateFolder := paidTime.Format("2006-01-02")

	// Create date-based directory with full path
	currentDir, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	dirPath := filepath.Join(currentDir, "webhooks", "data", dateFolder)
	log.Printf("Attempting to create directory: %s", dirPath)

	// Ensure directory exists with full permissions
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Printf("Error creating directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s_%s.json", uuid.New().String(), webhook.ID)
	filePath := filepath.Join(dirPath, filename)

	// Write webhook data to file
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file %s: %v", filePath, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Pretty print JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(webhook); err != nil {
		log.Printf("Error writing to file %s: %v", filePath, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Webhook saved to %s", filePath)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received and saved"))
}
