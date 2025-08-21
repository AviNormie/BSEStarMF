package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/BrokingSapphire/BSEStarMF/internal/bse"
)

// LumpsumOrderHandler handles lumpsum order entry requests
func LumpsumOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var req bse.LumpsumOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode lumpsum order request: %v", err)
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}

	// Log the order attempt
	log.Printf("Lumpsum Order attempt - TransCode: %s, TransNo: %s, UserID: %d, ClientCode: %s, SchemeCd: %s, BuySell: %s",
		req.TransCode, req.TransNo, req.UserID, req.ClientCode, req.SchemeCd, req.BuySell)

	// Create BSE client and place order
	client := bse.NewLumpsumOrderClient()
	resp, err := client.PlaceOrder(req)
	if err != nil {
		log.Printf("Lumpsum Order failed for TransNo %s: %v", req.TransNo, err)
		http.Error(w, fmt.Sprintf("Order placement failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Log the response status
	if resp.IsSuccessResponse() {
		log.Printf("Lumpsum Order successful for TransNo %s: OrderNumber %d - %s",
			req.TransNo, resp.OrderNumber, resp.BSERemarks)
	} else {
		log.Printf("Lumpsum Order failed for TransNo %s: %s",
			req.TransNo, resp.BSERemarks)
	}

	// Return response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode lumpsum order response: %v", err)
	}
}

// LumpsumOrderHealthHandler provides a health check endpoint for lumpsum orders
func LumpsumOrderHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"status":    "healthy",
		"message":   "BSE StAR MF Lumpsum Order Entry API is running",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
		"endpoints": []string{
			"POST /api/lumpsum/order - Place lumpsum order",
			"GET /api/lumpsum/health - Health check",
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// LumpsumOrderTestHandler handles test lumpsum order requests (no BSE call)
func LumpsumOrderTestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var req bse.LumpsumOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode lumpsum order request: %v", err)
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}

	// Log the order attempt
	log.Printf("Test Lumpsum Order - TransCode: %s, TransNo: %s, UserID: %d, ClientCode: %s",
		req.TransCode, req.TransNo, req.UserID, req.ClientCode)

	// Create BSE client and place test order
	client := bse.NewLumpsumOrderClient()
	resp, err := client.PlaceOrderTest(req)
	if err != nil {
		log.Printf("Test Lumpsum Order validation failed for TransNo %s: %v", req.TransNo, err)
		http.Error(w, fmt.Sprintf("Order validation failed: %v", err), http.StatusBadRequest)
		return
	}

	// Log success
	log.Printf("Test Lumpsum Order successful for TransNo %s: %s", req.TransNo, resp.BSERemarks)

	// Return response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode test lumpsum order response: %v", err)
	}
}