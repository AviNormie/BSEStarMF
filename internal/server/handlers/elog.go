package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"sapphirebroking.com/sapphire_mf/internal/server/services"
	"sapphirebroking.com/sapphire_mf/internal/util"
)

// ELOGHandler handles ELOG-related requests
type ELOGHandler struct {
	elogService *services.ELOGClientService
	logger      util.Logger
}

// NewELOGHandler creates a new ELOG handler
func NewELOGHandler(elogService *services.ELOGClientService, logger util.Logger) *ELOGHandler {
	return &ELOGHandler{
		elogService: elogService,
		logger:      logger,
	}
}

// ELOGRequestHandler handles ELOG authentication requests
func (h *ELOGHandler) ELOGRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendBSEErrorResponse(w, "Method not allowed")
		return
	}

	// Parse request body
	var req services.ELOGRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode ELOG request: %v", err)
		h.sendBSEErrorResponse(w, "Invalid JSON payload")
		return
	}

	// Manual validation
	if err := h.validateELOGRequest(&req); err != nil {
		h.logger.Error("ELOG request validation failed: %v", err)
		// Send validation error in BSE format
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&services.ELOGResponse{
			StatusCode: "101",
			AuthURL:    "",
			ErrorDesc:  err.Error(),
			IntRefNo:   req.IntRefNo,
		})
		return
	}

	// Submit ELOG request to BSE (this will use the simulation logic)
	resp, err := h.elogService.SubmitELOGRequest(r.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to submit ELOG request: %v", err)
		h.sendBSEErrorResponse(w, "Failed to process ELOG request")
		return
	}

	// Send BSE response directly (EXACT BSE FORMAT)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)  // Direct BSE response - NO WRAPPING
}

// ELOGCallbackHandler handles loopback callbacks from BSE
func (h *ELOGHandler) ELOGCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendBSEErrorResponse(w, "Method not allowed")
		return
	}

	// Parse query parameters
	status := r.URL.Query().Get("STATUS")
	elgStatus := r.URL.Query().Get("elgstatus")

	if status == "" {
		h.sendBSEErrorResponse(w, "Missing STATUS parameter")
		return
	}

	// Validate ELOG status
	isValid, message := h.elogService.ValidateELOGStatus(status, elgStatus)
	description := h.elogService.GetELGStatusDescription(elgStatus)

	// Create callback response
	callbackResp := map[string]interface{}{
		"status":      status,
		"elg_status":  elgStatus,
		"is_valid":    isValid,
		"message":     message,
		"description": description,
		"timestamp":   time.Now().Unix(),
	}

	h.logger.Info("ELOG callback received - Status: %s, ELG Status: %s, Valid: %v", status, elgStatus, isValid)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(callbackResp)
}

// validateELOGRequest performs manual validation on ELOG request
func (h *ELOGHandler) validateELOGRequest(req *services.ELOGRequest) error {
	// Check required fields
	if req.UserID == "" {
		return fmt.Errorf("userid is required")
	}
	if req.MemberID == "" {
		return fmt.Errorf("memberid is required")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	if req.ClientCode == "" {
		return fmt.Errorf("clientcode is required")
	}
	if req.Holder == "" {
		return fmt.Errorf("holder is required")
	}
	if req.DocumentType == "" {
		return fmt.Errorf("documenttype is required")
	}
	if req.LoopbackURL == "" {
		return fmt.Errorf("loopbackurl is required")
	}
	if req.AllowLoopbackMsg == "" {
		return fmt.Errorf("allowloopbackmsg is required")
	}

	// Validate holder value
	if req.Holder != "1" && req.Holder != "2" && req.Holder != "3" {
		return fmt.Errorf("holder must be 1, 2, or 3")
	}

	// Validate document type
	if req.DocumentType != "NRM" && req.DocumentType != "RIA" {
		return fmt.Errorf("document type must be NRM or RIA")
	}

	// Validate allowloopbackmsg
	if req.AllowLoopbackMsg != "Y" && req.AllowLoopbackMsg != "N" {
		return fmt.Errorf("allowloopbackmsg must be Y or N")
	}

	// Validate loopback URL format
	if _, err := url.Parse(req.LoopbackURL); err != nil {
		return fmt.Errorf("invalid loopback URL format")
	}

	// Ensure HTTPS for loopback URL
	if !strings.HasPrefix(req.LoopbackURL, "https://") {
		return fmt.Errorf("loopback URL must use HTTPS")
	}

	return nil
}

// sendBSEErrorResponse sends error in BSE format
func (h *ELOGHandler) sendBSEErrorResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)  // BSE always returns 200
	json.NewEncoder(w).Encode(&services.ELOGResponse{
		StatusCode: "101",
		AuthURL:    "",
		ErrorDesc:  message,
		IntRefNo:   "",
	})
}