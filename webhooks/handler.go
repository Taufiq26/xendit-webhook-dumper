package webhooks

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
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

// Webhook represents a single webhook entry
type Webhook struct {
	Timestamp time.Time     `json:"timestamp"`
	Payload   XenditWebhook `json:"payload"`
}

// WebhookCollection is a thread-safe collection of webhooks
type WebhookCollection struct {
	mu       sync.Mutex
	Webhooks []Webhook `json:"webhooks"`
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

	// Parse the webhook payload
	var webhookPayload XenditWebhook
	if err := json.Unmarshal(body, &webhookPayload); err != nil {
		log.Printf("Error parsing webhook payload: %v", err)
		http.Error(w, "Error parsing webhook payload", http.StatusBadRequest)
		return
	}

	// Determine the date for file organization
	paidTime, err := time.Parse(time.RFC3339, webhookPayload.PaidAt)
	if err != nil {
		log.Printf("Error parsing paid_at time: %v. Using current time.", err)
		paidTime = time.Now()
	}
	dateFolder := paidTime.Format("2006-01-02")

	// Use home directory for storing webhooks
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error getting home directory: %v", err)
		homeDir = "/root" // Fallback for server environment
	}

	// Create full path for storing webhooks
	dirPath := filepath.Join(homeDir, "xendit-webhook-data", dateFolder)

	// Ensure directory exists
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Printf("Error creating directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create or append to the daily webhook file
	filePath := filepath.Join(dirPath, "webhooks.json")

	// Mutex to prevent concurrent file access
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	// Read existing webhooks or create new collection
	var collection WebhookCollection
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		// File exists, read existing content
		existingData, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading existing file: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if len(existingData) > 0 {
			if err := json.Unmarshal(existingData, &collection); err != nil {
				log.Printf("Error parsing existing webhooks: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
	}

	// Add new webhook to collection
	newWebhook := Webhook{
		Timestamp: time.Now(),
		Payload:   webhookPayload,
	}
	collection.Webhooks = append(collection.Webhooks, newWebhook)

	// Write updated collection back to file
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Pretty print JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(collection); err != nil {
		log.Printf("Error writing to file: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Webhook saved to %s", filePath)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received and saved"))
}
