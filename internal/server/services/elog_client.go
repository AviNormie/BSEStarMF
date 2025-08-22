package services

import (
	
	"context"
	
	"fmt"

	"net/http"
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
		baseURL: "https://bsestarmfdemo.bseindia.com/BSEMFWEBAPI/api/_2FAELOGController/_2FAELOG/w",
		logger:  logger,
	}
}

// SubmitELOGRequest submits ELOG request to BSE and returns authentication URL
func (e *ELOGClientService) SubmitELOGRequest(ctx context.Context, req *ELOGRequest) (*ELOGResponse, error) {
	e.logger.Info("Submitting ELOG request for client: %s", req.ClientCode)

	// For demo purposes, simulate BSE responses based on client code
	// In production, this would make actual HTTP call to BSE
	response := e.simulateBSEResponse(req)

	e.logger.Info("ELOG request completed for client: %s, status: %s", req.ClientCode, response.StatusCode)
	return response, nil
}

// simulateBSEResponse simulates BSE ELOG responses for testing
func (e *ELOGClientService) simulateBSEResponse(req *ELOGRequest) *ELOGResponse {
	// Simulate different scenarios based on client code
	switch req.ClientCode {
	case "invalid123":
		// Scenario 1: Invalid Client Code
		return &ELOGResponse{
			StatusCode: "101",
			AuthURL:    "",
			ErrorDesc:  "Invalid Client Code",
			IntRefNo:   req.IntRefNo,
		}
	case "already123":
		// Scenario 2: Authentication already done
		return &ELOGResponse{
			StatusCode: "101",
			AuthURL:    "",
			ErrorDesc:  "FAILED : AUTHENTICATION IS ALREADY DONE",
			IntRefNo:   req.IntRefNo,
		}
	default:
		// Scenario 3: Success
		authURL := fmt.Sprintf("https://www.bsestarmf.in/3log/liefwbc23fq8pfg8qpwcwq8u8_%s_%d", 
			req.ClientCode, time.Now().Unix())
		return &ELOGResponse{
			StatusCode: "100",
			AuthURL:    authURL,
			ErrorDesc:  "ELOG Link Generated Successfully",
			IntRefNo:   req.IntRefNo,
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