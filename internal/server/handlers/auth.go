package handlers

import (
	"encoding/json"
	"net/http"
	"sapphirebroking.com/sapphire_mf/internal/server/services"
	"sapphirebroking.com/sapphire_mf/internal/util"
	"strings"
)

// HTTP Request/Response structs for REST API
type GetPasswordRequest struct {
	UserID   string `json:"user_id" validate:"required,max=5"`
	Password string `json:"password" validate:"required,max=20"`
	PassKey  string `json:"pass_key" validate:"required,max=10"`
}

type GetPasswordResponse struct {
	Success           bool   `json:"success"`
	ResponseCode      string `json:"response_code"`
	EncryptedPassword string `json:"encrypted_password,omitempty"`
	Message           string `json:"message,omitempty"`
	SessionValidity   string `json:"session_validity,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
}

// Remove the global variable and init() function
// var soapService *services.SOAPClientService
// func init() { ... }

func GetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "Method not allowed",
		})
		return
	}
	
	// Create logger and SOAP service within the handler
	logger := util.NewStandardLogger()
	soapService, err := services.NewSOAPClientService(logger)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "SOAP service unavailable",
		})
		return
	}
	
	var req GetPasswordRequest
	if decodeErr := json.NewDecoder(r.Body).Decode(&req); decodeErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "Invalid JSON payload",
		})
		return
	}
	
	// Validate required fields according to BSE specs
	if strings.TrimSpace(req.UserID) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "USER ID SHOULD NOT BE BLANK",
			Code:    "VALIDATION_ERROR",
		})
		return
	}
	
	if strings.TrimSpace(req.Password) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "PASSWORD SHOULD NOT BE BLANK",
			Code:    "VALIDATION_ERROR",
		})
		return
	}
	
	if strings.TrimSpace(req.PassKey) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "PASSKEY SHOULD NOT BE BLANK",
			Code:    "VALIDATION_ERROR",
		})
		return
	}
	
	// Validate field lengths according to BSE specs
	if len(req.UserID) > 5 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "User ID must not exceed 5 characters",
			Code:    "VALIDATION_ERROR",
		})
		return
	}
	
	if len(req.Password) > 20 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "Password must not exceed 20 characters",
			Code:    "VALIDATION_ERROR",
		})
		return
	}
	
	if len(req.PassKey) > 10 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "Pass Key must not exceed 10 characters",
			Code:    "VALIDATION_ERROR",
		})
		return
	}
	
	// Validate PassKey is alphanumeric only (BSE requirement)
	for _, char := range req.PassKey {
		if !((char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9')) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{
				Success: false,
				Error:   "PassKey must be alphanumeric only (no special characters)",
				Code:    "VALIDATION_ERROR",
			})
			return
		}
	}
	
	// Call BSE SOAP service
	authResp, err := soapService.Authenticate(r.Context(), req.UserID, req.Password, req.PassKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "Failed to authenticate with BSE service",
		})
		return
	}
	
	// Prepare response based on BSE response codes
	if authResp.ResponseCode == "100" {
		// Success - BSE returned code 100
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GetPasswordResponse{
			Success:           true,
			ResponseCode:      authResp.ResponseCode,
			EncryptedPassword: authResp.EncryptedPassword,
			Message:           "Authentication successful",
			SessionValidity:   "1 Hour", // As per BSE documentation
		})
	} else {
		// Authentication failed - BSE returned code 101 or other error
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(GetPasswordResponse{
			Success:      false,
			ResponseCode: authResp.ResponseCode,
			Message:      authResp.ErrorMessage,
		})
	}
}