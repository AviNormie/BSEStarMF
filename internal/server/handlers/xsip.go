package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sapphirebroking.com/sapphire_mf/internal/server/services"
	"strings"
)

var xsipSoapService *services.SOAPClientService

func init() {
	var err error
	xsipSoapService, err = services.NewSOAPClientService()
	if err != nil {
		// Log error but don't panic - handle gracefully in handler
	}
}

// XSIPHandler handles XSIP registration and cancellation requests
func XSIPHandler(w http.ResponseWriter, r *http.Request) {
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
	if xsipSoapService == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "SOAP service unavailable",
		})
		return
	}
	
	var req services.XSIPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "Invalid JSON payload",
		})
		return
	}
	
	// Validate required fields
	if err := validateXSIPRequest(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    "VALIDATION_ERROR",
		})
		return
	}
	
	// Call BSE SOAP service
	xsipResp, err := xsipSoapService.XSIPOrderEntry(r.Context(), &req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "Failed to process XSIP request with BSE service",
		})
		return
	}
	
	// Return response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(xsipResp)
}

func validateXSIPRequest(req *services.XSIPRequest) error {
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
	if strings.TrimSpace(req.EUINFlag) == "" {
		return fmt.Errorf("EUIN flag should not be blank")
	}
	if strings.TrimSpace(req.DPC) == "" {
		return fmt.Errorf("DPC should not be blank")
	}
	
	// Validate BSE Code (Filler3) - MANDATORY
	if strings.TrimSpace(req.Filler3) == "" {
		return fmt.Errorf("BSE Code (filler3) is mandatory")
	}
	
	// Validate BSE Code Remark (Filler4) - Conditional Mandatory
	if req.Filler3 == "13" && strings.TrimSpace(req.Filler4) == "" {
		return fmt.Errorf("BSE Code Remark (filler4) is mandatory when BSE Code is '13' (Others)")
	}
	
	// Validate XSIP Mandate ID for XSIP Orders
	if req.TransactionCode == "NEW" && strings.TrimSpace(req.XSIPMandateID) == "" {
		return fmt.Errorf("XSIP Mandate ID is mandatory for XSIP Orders")
	}
	
	// Validate ISIP Mandate ID for ISIP Orders (Param2)
	// Note: This would need additional logic to determine if it's an ISIP order
	
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
	if len(req.Filler3) > 2 {
		return fmt.Errorf("BSE Code must not exceed 2 characters")
	}
	if len(req.Filler4) > 200 {
		return fmt.Errorf("BSE Code Remark must not exceed 200 characters")
	}
	
	return nil
}