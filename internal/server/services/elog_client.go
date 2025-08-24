package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"  // ← ADDED: Missing import
	"time"

	"sapphirebroking.com/sapphire_mf/internal/util"
)

// ELOGRequest represents the exact ELOG request structure as per BSE specification
type ELOGRequest struct {
	UserID            string `json:"userid"`            // Login Id (mandatory)
	MemberID          string `json:"memberid"`          // Member Code (mandatory)
	Password          string `json:"password"`          // Login Password (mandatory)
	ClientCode        string `json:"clientcode"`        // Client Code (mandatory)
	Holder            string `json:"holder"`            // 1/2/3 for Multiple Holders (mandatory)
	DocumentType      string `json:"documenttype"`      // NRM/RIA (mandatory)
	IntRefNo          string `json:"intrefno"`          // Internal Reference (optional)
	LoopbackURL       string `json:"loopbackurl"`       // Loopback URL (mandatory)
	AllowLoopbackMsg  string `json:"allowloopbackmsg"`  // Y/N for additional variables (mandatory)
}

// ELOGResponse represents the exact ELOG response structure as per BSE specification
type ELOGResponse struct {
	StatusCode  string `json:"statuscode"`  // 100-success, 101-failure
	AuthURL     string `json:"authurl"`     // URL for member to open in app
	ErrorDesc   string `json:"errordesc"`   // Error/Success Description
	IntRefNo    string `json:"intrefno"`    // Internal reference sent in request
}

// ELOGClientService handles ELOG API operations
type ELOGClientService struct {
	client  *http.Client
	baseURL string
	logger  util.Logger
}

// NewELOGClientService creates a new ELOG client service
func NewELOGClientService(logger util.Logger) *ELOGClientService {
	return &ELOGClientService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		// BSE ELOG API endpoint
		baseURL: "https://bsestarmfdemo.bseindia.com/BSEMFWEBAPI/api/_2FAELOGController/_2FAELOG/w",
		logger:  logger,
	}
}

// SubmitELOGRequest submits ELOG request to BSE and returns authentication URL
func (e *ELOGClientService) SubmitELOGRequest(ctx context.Context, req *ELOGRequest) (*ELOGResponse, error) {
	e.logger.Info("Submitting ELOG request to BSE for client: %s", req.ClientCode)

	// Convert request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		e.logger.Error("Failed to marshal ELOG request: %v", err)
		return &ELOGResponse{
			StatusCode: "101",
			AuthURL:    "",
			ErrorDesc:  "Request serialization failed",
			IntRefNo:   req.IntRefNo,
		}, nil
	}

	// Create HTTP request to BSE
	httpReq, err := http.NewRequestWithContext(ctx, "POST", e.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		e.logger.Error("Failed to create HTTP request: %v", err)
		return &ELOGResponse{
			StatusCode: "101",
			AuthURL:    "",
			ErrorDesc:  "Request creation failed",
			IntRefNo:   req.IntRefNo,
		}, nil
	}

	// Set headers for BSE API
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", "BSE-StarMF-Client/1.0")

	// Make actual call to BSE
	e.logger.Info("Making HTTP request to BSE ELOG API: %s", e.baseURL)
	resp, err := e.client.Do(httpReq)
	if err != nil {
		e.logger.Error("BSE ELOG API connection failed: %v", err)
		return &ELOGResponse{
			StatusCode: "101",
			AuthURL:    "",
			ErrorDesc:  "Failed to connect to BSE service",
			IntRefNo:   req.IntRefNo,
		}, nil
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		e.logger.Error("Failed to read BSE response: %v", err)
		return &ELOGResponse{
			StatusCode: "101",
			AuthURL:    "",
			ErrorDesc:  "Failed to read BSE response",
			IntRefNo:   req.IntRefNo,
		}, nil
	}

	e.logger.Info("BSE ELOG API response status: %d, body length: %d", resp.StatusCode, len(respBody))

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		e.logger.Error("BSE ELOG API returned error status: %d, body: %s", resp.StatusCode, string(respBody))
		return &ELOGResponse{
			StatusCode: "101",
			AuthURL:    "",
			ErrorDesc:  fmt.Sprintf("BSE API error: HTTP %d", resp.StatusCode),
			IntRefNo:   req.IntRefNo,
		}, nil
	}

	// Parse BSE response
	var bseResponse ELOGResponse
	if err := json.Unmarshal(respBody, &bseResponse); err != nil {
		e.logger.Error("Failed to parse BSE ELOG response: %v, body: %s", err, string(respBody))
		
		// Try to handle non-JSON response from BSE
		responseStr := string(respBody)
		if responseStr != "" {
			// BSE might return pipe-delimited or other format
			return e.parseBSETextResponse(responseStr, req.IntRefNo), nil  // ← FIXED: Added nil error
		}
		
		return &ELOGResponse{
			StatusCode: "101",
			AuthURL:    "",
			ErrorDesc:  "Invalid response format from BSE",
			IntRefNo:   req.IntRefNo,
		}, nil
	}

	e.logger.Info("ELOG request completed for client: %s, BSE status: %s", req.ClientCode, bseResponse.StatusCode)
	return &bseResponse, nil
}

// parseBSETextResponse handles non-JSON responses from BSE (pipe-delimited format)
func (e *ELOGClientService) parseBSETextResponse(response, intRefNo string) *ELOGResponse {
	e.logger.Info("Parsing BSE text response: %s", response)
	
	// BSE might return format like: "100|authurl|success_message" or "101|error_message"
	parts := strings.Split(response, "|")
	if len(parts) < 2 {
		return &ELOGResponse{
			StatusCode: "101",
			AuthURL:    "",
			ErrorDesc:  "Invalid BSE response format",
			IntRefNo:   intRefNo,
		}
	}
	
	statusCode := strings.TrimSpace(parts[0])
	
	if statusCode == "100" && len(parts) >= 3 {
		// Success response with auth URL
		return &ELOGResponse{
			StatusCode: "100",
			AuthURL:    strings.TrimSpace(parts[1]),
			ErrorDesc:  strings.TrimSpace(parts[2]),
			IntRefNo:   intRefNo,
		}
	} else {
		// Error response
		errorDesc := "Unknown error"
		if len(parts) >= 2 {
			errorDesc = strings.TrimSpace(parts[1])
		}
		return &ELOGResponse{
			StatusCode: "101",
			AuthURL:    "",
			ErrorDesc:  errorDesc,
			IntRefNo:   intRefNo,
		}
	}
}

// ValidateELOGStatus validates the ELOG status from loopback
func (e *ELOGClientService) ValidateELOGStatus(status, elgStatus string) (bool, string) {
	switch status {
	case "SUCCESS":
		return e.parseELGStatus(elgStatus), "ELOG authentication successful"
	case "FAILURE":
		return false, "ELOG authentication failed"
	case "PENDING":
		return false, "ELOG authentication pending"
	default:
		return false, "Unknown ELOG status"
	}
}

// parseELGStatus parses the 3-digit ELG status
func (e *ELOGClientService) parseELGStatus(elgStatus string) bool {
	if len(elgStatus) != 3 {
		return false
	}
	return elgStatus[0] == '1'
}

// GetELGStatusDescription returns human-readable description of ELG status
func (e *ELOGClientService) GetELGStatusDescription(elgStatus string) string {
	if len(elgStatus) != 3 {
		return "Invalid ELG status format"
	}

	switch elgStatus {
	case "111":
		return "All holders ELOG approved"
	case "110":
		return "Primary and secondary holders ELOG approved"
	case "101":
		return "Primary and third holders ELOG approved"
	case "100":
		return "Only primary holder ELOG approved"
	case "011":
		return "Secondary and third holders ELOG approved"
	case "010":
		return "Only secondary holder ELOG approved"
	case "001":
		return "Only third holder ELOG approved"
	case "000":
		return "No holders ELOG approved"
	default:
		return fmt.Sprintf("Custom ELG status: %s", elgStatus)
	}
}