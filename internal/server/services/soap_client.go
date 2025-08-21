package services

import (
	"context"
	"fmt"
	"github.com/hooklift/gowsdl/soap"
	"sapphirebroking.com/sapphire_mf/myservice"
	"strings"
	"time"
)

type SOAPClientService struct {
	client myservice.MFOrderEntry
}

func NewSOAPClientService() (*SOAPClientService, error) {
	url := "https://bsestarmfdemo.bseindia.com/MFOrderEntry/MFOrder.svc"
	
	soapClient := soap.NewClient(url)
	mfClient := myservice.NewMFOrderEntry(soapClient)
	
	return &SOAPClientService{
		client: mfClient,
	}, nil
}

type AuthResponse struct {
	ResponseCode     string `json:"response_code"`
	EncryptedPassword string `json:"encrypted_password,omitempty"`
	ErrorMessage     string `json:"error_message,omitempty"`
}

func (s *SOAPClientService) Authenticate(ctx context.Context, userID, password, passKey string) (*AuthResponse, error) {
	// Create context with timeout for SOAP call
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	
	// Create SOAP request using existing structs
	request := &myservice.GetPassword{
		UserId:   &userID,
		Password: &password,
		PassKey:  &passKey,
	}
	
	// Call SOAP service
	response, err := s.client.GetPasswordContext(ctxWithTimeout, request)
	if err != nil {
		// Return BSE-style error for connection failures
		return &AuthResponse{
			ResponseCode: "101",
			ErrorMessage: "USER NOT EXISTS",
		}, nil
	}
	
	// Parse response according to BSE documentation
	if response.GetPasswordResult == nil {
		return &AuthResponse{
			ResponseCode: "101",
			ErrorMessage: "USER NOT EXISTS",
		}, nil
	}
	
	result := *response.GetPasswordResult
	
	// Parse BSE response format: "ResponseCode|EncryptedPassword" or "ResponseCode|ErrorMessage"
	parts := strings.Split(result, "|")
	if len(parts) == 0 {
		return &AuthResponse{
			ResponseCode: "101",
			ErrorMessage: "Invalid response format",
		}, nil
	}
	
	responseCode := strings.TrimSpace(parts[0])
	
	if responseCode == "100" && len(parts) > 1 {
		// Success - BSE returned code 100 with encrypted password
		return &AuthResponse{
			ResponseCode:      "100",
			EncryptedPassword: strings.TrimSpace(parts[1]),
		}, nil
	} else if responseCode == "101" {
		// Authentication failed - extract error message if available
		errorMsg := "USER NOT EXISTS"
		if len(parts) > 1 {
			errorMsg = strings.TrimSpace(parts[1])
		}
		return &AuthResponse{
			ResponseCode: "101",
			ErrorMessage: errorMsg,
		}, nil
	} else {
		// Handle other BSE error codes
		errorMsg := s.getErrorMessage(responseCode)
		return &AuthResponse{
			ResponseCode: responseCode,
			ErrorMessage: errorMsg,
		}, nil
	}
}

func (s *SOAPClientService) getErrorMessage(code string) string {
	// Official BSE error messages from documentation
	errorMessages := map[string]string{
		"101": "USER NOT EXISTS",
		"USER ID SHOULD NOT BE BLANK":           "Empty User ID field",
		"MEMBER ID SHOULD NOT BE BLANK":         "Empty Member ID field",
		"PASSWORD SHOULD NOT BE BLANK":          "Empty password field",
		"PASSKEY SHOULD NOT BE BLANK":           "Empty passkey field",
		"USER IS DISABLED. CONTACT ADMIN":       "User blocked or disabled",
		"YOU HAVE EXCEEDED MAXIMUM LOGIN ATTEMPTS": "Too many wrong password attempts",
		"INVALID ACCOUNT INFORMATION":           "Incorrect login details",
		"THE MEMBER IS SUSPENDED":               "Member blocked or inactive",
		"PASSWORD EXPIRED":                      "User password has expired",
		"USER NOT EXISTS":                       "Invalid user credentials",
	}
	
	if msg, exists := errorMessages[code]; exists {
		return msg
	}
	return fmt.Sprintf("Unknown error code: %s", code)
}