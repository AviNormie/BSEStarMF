package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AviNormie/BSEStarMF/internal/bse"
)

// UCCRegistrationHandler handles UCC registration requests
func UCCRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var req bse.UCCRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode UCC registration request: %v", err)
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}

	// Log the registration attempt
	log.Printf("UCC Registration attempt - UserID: %s, MemberCode: %s, RegnType: %s", 
		req.UserID, req.MemberCode, req.RegnType)

	// Create BSE client and perform registration
	client := bse.NewUCCRegistrationClient()
	resp, err := client.RegisterClient(req)
	if err != nil {
		log.Printf("UCC Registration failed for UserID %s: %v", req.UserID, err)
		http.Error(w, fmt.Sprintf("Registration failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Log the response status
	if resp.IsSuccessResponse() {
		log.Printf("UCC Registration successful for UserID %s: %s", req.UserID, resp.Remarks)
	} else {
		log.Printf("UCC Registration failed for UserID %s: Status %s - %s", req.UserID, resp.Status, resp.Remarks)
	}

	// Return response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode UCC registration response: %v", err)
	}
}

// UCCHealthHandler provides a health check endpoint
func UCCHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	response := map[string]interface{}{
		"status":    "healthy",
		"message":   "BSE StAR MF UCC Registration API is running",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}