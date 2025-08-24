package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sapphirebroking.com/sapphire_mf/internal/server/services"
	"sapphirebroking.com/sapphire_mf/internal/util"
	"strings"
)

// REMOVE these lines:
// var enhancedSipService *services.SOAPClientService
// func init() { ... }

func EnhancedSIPCancellationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// CREATE logger and SOAP service here
	logger := util.NewStandardLogger()
	enhancedSipService, err := services.NewSOAPClientService(logger)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to initialize SOAP service",
		})
		return
	}
	
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "Method not allowed",
		})
		return
	}
	
	// Check if service is available
	if enhancedSipService == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "Service unavailable",
		})
		return
	}
	
	var req services.EnhancedSIPCancellationRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
			"code":    "VALIDATION_ERROR",
		})
		return
	}
	
	// REMOVE this duplicate line:
	// err = json.NewDecoder(r.Body).Decode(&req)
	
	// Validate the request (use the already decoded req)
	if err := validateEnhancedSIPRequest(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
			"code":    "VALIDATION_ERROR",
		})
		return
	}
	
	// Call BSE Enhanced API service
	sipResp, err := enhancedSipService.EnhancedSIPCancellation(r.Context(), &req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "Failed to process Enhanced SIP cancellation request",
		})
		return
	}
	
	// Return response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sipResp)
}

func validateEnhancedSIPRequest(req *services.EnhancedSIPCancellationRequest) error {
	// Validate required fields
	if strings.TrimSpace(req.LoginID) == "" {
		return fmt.Errorf("login ID should not be blank")
	}
	if strings.TrimSpace(req.MemberCode) == "" {
		return fmt.Errorf("member code should not be blank")
	}
	if strings.TrimSpace(req.Password) == "" {
		return fmt.Errorf("password should not be blank")
	}
	if strings.TrimSpace(req.ClientCode) == "" {
		return fmt.Errorf("client code should not be blank")
	}
	if req.RegnNo == 0 {
		return fmt.Errorf("registration number should not be blank")
	}
	if strings.TrimSpace(req.CeaseBseCode) == "" {
		return fmt.Errorf("cease BSE code is mandatory")
	}
	
	// Validate BSE Code Remark - Conditional Mandatory
	if req.CeaseBseCode == "13" && strings.TrimSpace(req.Remarks) == "" {
		return fmt.Errorf("remarks are mandatory when cease BSE code is '13' (Others)")
	}
	
	// Validate field lengths
	if len(req.LoginID) > 20 {
		return fmt.Errorf("login ID must not exceed 20 characters")
	}
	if len(req.MemberCode) > 20 {
		return fmt.Errorf("member code must not exceed 20 characters")
	}
	if len(req.ClientCode) > 10 {
		return fmt.Errorf("client code must not exceed 10 characters")
	}
	if len(req.IntRefNo) > 20 {
		return fmt.Errorf("internal reference number must not exceed 20 characters")
	}
	if len(req.CeaseBseCode) > 2 {
		return fmt.Errorf("cease BSE code must not exceed 2 characters")
	}
	if len(req.Remarks) > 200 {
		return fmt.Errorf("remarks must not exceed 200 characters")
	}
	
	return nil
}