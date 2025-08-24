package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"sapphirebroking.com/sapphire_mf/internal/util"
)

// EnhancedAPIClient handles BSE Enhanced REST API operations
type EnhancedAPIClient struct {
	client  *http.Client
	baseURL string
	logger  util.Logger
}

// NewEnhancedAPIClient creates a new Enhanced API client
func NewEnhancedAPIClient(logger util.Logger) *EnhancedAPIClient {
	return &EnhancedAPIClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://bsestarmfdemo.bseindia.com/StarMFAPI/api",
		logger:  logger,
	}
}

// CallEnhancedSIPCancellation calls the BSE Enhanced SIP Cancellation API
func (e *EnhancedAPIClient) CallEnhancedSIPCancellation(ctx context.Context, req *EnhancedSIPCancellationRequest) (*EnhancedSIPCancellationResponse, error) {
	url := fmt.Sprintf("%s/SIP/SIPCancellation", e.baseURL)
	
	e.logger.Info("[DEBUG] Calling Enhanced SIP Cancellation API: %s", url)
	
	// Convert request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}
	
	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	
	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	
	// Make the request
	resp, err := e.client.Do(httpReq)
	if err != nil {
		e.logger.Error("[ERROR] Enhanced SIP API call failed: %v", err)
		return nil, fmt.Errorf("API call failed: %v", err)
	}
	defer resp.Body.Close()
	
	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	
	e.logger.Info("[DEBUG] Enhanced SIP API Response: %s", string(body))
	
	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}
	
	// Parse response
	var apiResp EnhancedSIPCancellationResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	
	return &apiResp, nil
}

// CallEnhancedXSIPCancellation calls the BSE Enhanced XSIP Cancellation API
func (e *EnhancedAPIClient) CallEnhancedXSIPCancellation(ctx context.Context, req *EnhancedXSIPCancellationRequest) (*EnhancedXSIPCancellationResponse, error) {
	url := fmt.Sprintf("%s/XSIP/XSIPCancellation", e.baseURL)
	
	e.logger.Info("[DEBUG] Calling Enhanced XSIP Cancellation API: %s", url)
	
	// Convert request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}
	
	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	
	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	
	// Make the request
	resp, err := e.client.Do(httpReq)
	if err != nil {
		e.logger.Error("[ERROR] Enhanced XSIP API call failed: %v", err)
		return nil, fmt.Errorf("API call failed: %v", err)
	}
	defer resp.Body.Close()
	
	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	
	e.logger.Info("[DEBUG] Enhanced XSIP API Response: %s", string(body))
	
	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}
	
	// Parse response
	var apiResp EnhancedXSIPCancellationResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	
	return &apiResp, nil
}