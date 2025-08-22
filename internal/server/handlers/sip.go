package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sapphirebroking.com/sapphire_mf/internal/server/services"
	"strings"
)

var sipSoapService *services.SOAPClientService

func init() {
	var err error
	sipSoapService, err = services.NewSOAPClientService()
	if err != nil {
		// Log error but don't panic - handle gracefully in handler
	}
}

// SIPHandler handles SIP registration and cancellation requests
func SIPHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "Method not allowed",
		})
		return
	}
	
	// Check if SOAP service is available
	if sipSoapService == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "SOAP service unavailable",
		})
		return
	}
	
	var req services.SIPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "Invalid JSON payload",
		})
		return
	}
	
	// Validate required fields
	if err := validateSIPRequest(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    "VALIDATION_ERROR",
		})
		return
	}
	
	// Call BSE SOAP service
	sipResp, err := sipSoapService.SIPOrderEntry(r.Context(), &req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "Failed to process SIP request with BSE service",
		})
		return
	}
	
	// Return response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sipResp)
}

func validateSIPRequest(req *services.SIPRequest) error {
	// Validate transaction code
	if req.TransactionCode != "NEW" && req.TransactionCode != "CXL" {
		return fmt.Errorf("transaction code must be NEW or CXL")
	}
	
	// Validate required fields
	if strings.TrimSpace(req.UniqueRefNo) == "" {
		return fmt.Errorf("unique reference number should not be blank")
	}
	if strings.TrimSpace(req.SchemeCode) == "" {
		return fmt.Errorf("scheme code should not be blank")
	}
	if strings.TrimSpace(req.MemberID) == "" {
		return fmt.Errorf("member ID should not be blank")
	}
	if strings.TrimSpace(req.ClientCode) == "" {
		return fmt.Errorf("client code should not be blank")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return fmt.Errorf("user ID should not be blank")
	}
	if strings.TrimSpace(req.Password) == "" {
		return fmt.Errorf("password should not be blank")
	}
	if strings.TrimSpace(req.PassKey) == "" {
		return fmt.Errorf("pass key should not be blank")
	}
	
	// Validate BSE Code (Filler2) - MANDATORY
	if strings.TrimSpace(req.Filler2) == "" {
		return fmt.Errorf("BSE Code (filler2) is mandatory")
	}
	
	// Validate BSE Code Remark (Filler3) - Conditional Mandatory
	if req.Filler2 == "13" && strings.TrimSpace(req.Filler3) == "" {
		return fmt.Errorf("BSE Code Remark (filler3) is mandatory when BSE Code is '13' (Others)")
	}
	
	// Validate field lengths
	if len(req.UniqueRefNo) > 19 {
		return fmt.Errorf("unique reference number must not exceed 19 characters")
	}
	if len(req.UserID) > 5 {
		return fmt.Errorf("user ID must not exceed 5 characters")
	}
	if len(req.Password) > 250 {
		return fmt.Errorf("password must not exceed 250 characters")
	}
	if len(req.PassKey) > 10 {
		return fmt.Errorf("pass key must not exceed 10 characters")
	}
	if len(req.Filler2) > 2 {
		return fmt.Errorf("BSE Code must not exceed 2 characters")
	}
	if len(req.Filler3) > 200 {
		return fmt.Errorf("BSE Code Remark must not exceed 200 characters")
	}
	
	return nil
}