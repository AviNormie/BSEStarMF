package bse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// UCCRegistrationClient handles BSE StAR MF Enhanced UCC Registration API calls
type UCCRegistrationClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewUCCRegistrationClient creates a new UCC Registration client
func NewUCCRegistrationClient() *UCCRegistrationClient {
	return &UCCRegistrationClient{
		BaseURL: "https://bsestarmfdemo.bseindia.com/StarMFCommonAPI/ClientMaster/Registration",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// RegisterClient performs UCC registration with BSE
func (c *UCCRegistrationClient) RegisterClient(req UCCRegistrationRequest) (*UCCRegistrationResponse, error) {
	// Convert request to JSON
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	// Create HTTP request to BSE API
	httpReq, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	
	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	
	// Send request
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	// Parse JSON response
	var uccResp UCCRegistrationResponse
	if err := json.Unmarshal(respBody, &uccResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return &uccResp, nil
}