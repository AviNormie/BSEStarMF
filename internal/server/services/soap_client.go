package services

import (
	"context"
	"fmt"
	"github.com/hooklift/gowsdl/soap"
	"sapphirebroking.com/sapphire_mf/myservice"

	"sapphirebroking.com/sapphire_mf/internal/server/types"
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

// SIP Request structure for both NEW and CXL operations
type SIPRequest struct {
	TransactionCode    string `json:"transaction_code" validate:"required,oneof=NEW CXL"`
	UniqueRefNo       string `json:"unique_ref_no" validate:"required,max=19"`
	SchemeCode        string `json:"scheme_code" validate:"required,max=20"`
	MemberID          string `json:"member_id" validate:"required,max=20"`
	ClientCode        string `json:"client_code" validate:"required,max=20"`
	UserID            string `json:"user_id" validate:"required,max=5"`
	InternalRefNo     string `json:"internal_ref_no,omitempty" validate:"max=25"`
	TransMode         string `json:"trans_mode" validate:"required,oneof=D P"`
	DPTransactionMode string `json:"dp_transaction_mode" validate:"required,oneof=C N P"`
	StartDate         string `json:"start_date" validate:"required"`
	FrequencyType     string `json:"frequency_type" validate:"required,oneof=MONTHLY QUARTERLY WEEKLY"`
	FrequencyAllowed  int    `json:"frequency_allowed" validate:"required,min=1,max=1"`
	InstallmentAmount int    `json:"installment_amount" validate:"required"`
	NoOfInstallments  int    `json:"no_of_installments" validate:"required"`
	Remarks           string `json:"remarks,omitempty" validate:"max=100"`
	FolioNo           string `json:"folio_no,omitempty" validate:"max=20"`
	FirstOrderFlag    string `json:"first_order_flag" validate:"required,oneof=Y N"`
	SubBrCode         string `json:"sub_br_code,omitempty" validate:"max=15"`
	EUIN              string `json:"euin" validate:"required,max=20"`
	EUINDeclaration   string `json:"euin_declaration" validate:"required,oneof=Y N"`
	DPC               string `json:"dpc" validate:"required,oneof=Y"`
	RegID             string `json:"reg_id,omitempty"`
	IPAddress         string `json:"ip_address,omitempty" validate:"max=20"`
	Password          string `json:"password" validate:"required,max=250"`
	PassKey           string `json:"pass_key" validate:"required,max=10"`
	Param1            string `json:"param1,omitempty" validate:"max=20"` // Sub Broker ARN
	Param2            string `json:"param2,omitempty" validate:"max=10"` // End Date for Daily SIP
	Param3            string `json:"param3,omitempty" validate:"max=10"` // Mobile No
	Filler1           string `json:"filler1,omitempty" validate:"max=30"` // Email ID
	Filler2           string `json:"filler2" validate:"required,max=2"`   // BSE Code - MANDATORY
	Filler3           string `json:"filler3,omitempty" validate:"max=200"` // BSE Code Remark - Conditional Mandatory
	Filler4           string `json:"filler4,omitempty" validate:"max=30"`
	Filler5           string `json:"filler5,omitempty" validate:"max=30"`
	Filler6           string `json:"filler6,omitempty" validate:"max=30"`
}

// SIP Response structure for SOAP service
type SIPOrderResponse struct {
	Success              bool   `json:"success"`
	TransactionCode      string `json:"transaction_code"`
	UniqueRefNo          string `json:"unique_ref_no"`
	MemberID             string `json:"member_id"`
	ClientCode           string `json:"client_code"`
	UserID               string `json:"user_id"`
	SIPRegID             string `json:"sip_reg_id,omitempty"`
	BSERemarks           string `json:"bse_remarks,omitempty"`
	SuccessFlag          string `json:"success_flag"`
	FirstOrderTodayOrderNo string `json:"first_order_today_order_no,omitempty"`
	Message              string `json:"message,omitempty"`
}

// XSIP Request structure for both NEW and CXL operations - Updated to match exact documentation
type XSIPRequest struct {
	TransactionCode    string `json:"transaction_code" validate:"required,oneof=NEW CXL"`
	UniqueRefNo       string `json:"unique_ref_no" validate:"required,max=19"`
	SchemeCode        string `json:"scheme_code" validate:"required,max=20"`
	MemberID          string `json:"member_id" validate:"required,max=20"`
	ClientCode        string `json:"client_code" validate:"required,max=20"`
	UserID            string `json:"user_id" validate:"required,max=5"`
	InternalRefNo     string `json:"internal_ref_no,omitempty" validate:"max=25"`
	TransMode         string `json:"trans_mode" validate:"required,oneof=D P"`
	DPTransactionMode string `json:"dp_transaction_mode" validate:"required,oneof=C N P"`
	StartDate         string `json:"start_date" validate:"required"`
	FrequencyType     string `json:"frequency_type" validate:"required,oneof=MONTHLY QUARTERLY WEEKLY"`
	FrequencyAllowed  int    `json:"frequency_allowed" validate:"required,min=1,max=1"`
	InstallmentAmount int    `json:"installment_amount" validate:"required"`
	NoOfInstallments  int    `json:"no_of_installments" validate:"required"`
	Remarks           string `json:"remarks,omitempty" validate:"max=100"`
	FolioNo           string `json:"folio_no,omitempty" validate:"max=20"`
	FirstOrderFlag    string `json:"first_order_flag" validate:"required,oneof=Y N"`
	Brokerage         string `json:"brokerage,omitempty"`
	XSIPMandateID     string `json:"xsip_mandate_id,omitempty"` // BSE mandate ID (XSIP/Emandate) - Mandatory for XSIP Orders
	SubBrCode         string `json:"sub_br_code,omitempty" validate:"max=15"`
	EUIN              string `json:"euin,omitempty" validate:"max=20"`
	EUINFlag          string `json:"euin_flag" validate:"required,oneof=Y N"`
	DPC               string `json:"dpc" validate:"required,oneof=Y"`
	XSIPRegID         string `json:"xsip_reg_id,omitempty"`
	IPAddress         string `json:"ip_address,omitempty" validate:"max=20"`
	Password          string `json:"password" validate:"required,max=250"`
	PassKey           string `json:"pass_key" validate:"required,max=10"`
	Param1            string `json:"param1,omitempty" validate:"max=20"` // Sub Broker ARN
	Param2            string `json:"param2,omitempty" validate:"max=15"` // ISIP Mandate ID - Mandatory for ISIP Orders
	Param3            string `json:"param3,omitempty" validate:"max=10"` // End Date for Daily XSIP - DD/MM/YYYY
	Filler1           string `json:"filler1,omitempty" validate:"max=30"` // Mobile No
	Filler2           string `json:"filler2,omitempty" validate:"max=50"` // Email ID
	Filler3           string `json:"filler3" validate:"required,max=2"`   // BSE Code - MANDATORY
	Filler4           string `json:"filler4,omitempty" validate:"max=200"` // BSE Code Remark - Conditional Mandatory
	Filler5           string `json:"filler5,omitempty" validate:"max=30"`
	Filler6           string `json:"filler6,omitempty" validate:"max=30"`
}

// XSIP Response structure - Updated to match exact documentation
type XSIPOrderResponse struct {
	Success              bool   `json:"success"`
	TransactionCode      string `json:"transaction_code"`
	UniqueRefNo          string `json:"unique_ref_no"`
	MemberID             string `json:"member_id"`
	ClientCode           string `json:"client_code"`
	UserID               string `json:"user_id"`
	XSIPRegID            string `json:"xsip_reg_id,omitempty"`
	BSERemarks           string `json:"bse_remarks,omitempty"`
	SuccessFlag          string `json:"success_flag"`
	FirstOrderTodayOrderNo string `json:"first_order_today_order_no,omitempty"`
	Message              string `json:"message,omitempty"`
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

// SIPOrderEntry handles SIP registration and cancellation
func (s *SOAPClientService) SIPOrderEntry(ctx context.Context, req *SIPRequest) (*SIPOrderResponse, error) {
	// Create context with timeout for SOAP call
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	
	// Create SOAP request using existing structs
	soapRequest := &myservice.SipOrderEntryParam{
		TransactionCode:   &req.TransactionCode,
		UniqueRefNo:      &req.UniqueRefNo,
		SchemeCode:       &req.SchemeCode,
		MemberCode:       &req.MemberID,
		ClientCode:       &req.ClientCode,
		UserID:           &req.UserID,
		InternalRefNo:    &req.InternalRefNo,
		TransMode:        &req.TransMode,
		DpTxnMode:        &req.DPTransactionMode,
		StartDate:        &req.StartDate,
		FrequencyType:    &req.FrequencyType,
		FrequencyAllowed: stringPtr(fmt.Sprintf("%d", req.FrequencyAllowed)),
		InstallmentAmount: stringPtr(fmt.Sprintf("%d", req.InstallmentAmount)),
		NoOfInstallment:  stringPtr(fmt.Sprintf("%d", req.NoOfInstallments)),
		Remarks:          &req.Remarks,
		FolioNo:          &req.FolioNo,
		FirstOrderFlag:   &req.FirstOrderFlag,
		SubberCode:       &req.SubBrCode,
		Euin:             &req.EUIN,
		EuinVal:          &req.EUINDeclaration,
		DPC:              &req.DPC,
		RegId:            &req.RegID,
		IPAdd:            &req.IPAddress,
		Password:         &req.Password,
		PassKey:          &req.PassKey,
		Param1:           &req.Param1,
		Param2:           &req.Param2,
		Param3:           &req.Param3,
		Filler1:          &req.Filler1,
		Filler2:          &req.Filler2, // BSE Code
		Filler3:          &req.Filler3, // BSE Code Remark
		Filler4:          &req.Filler4,
		Filler5:          &req.Filler5,
		Filler6:          &req.Filler6,
	}
	
	// Call SOAP service
	response, err := s.client.SipOrderEntryParamContext(ctxWithTimeout, soapRequest)
	if err != nil {
		return &SIPOrderResponse{
			Success: false,
			Message: "Failed to connect to BSE service",
		}, nil
	}
	
	// Parse response
	if response.SipOrderEntryParamResult == nil {
		return &SIPOrderResponse{
			Success: false,
			Message: "Empty response from BSE service",
		}, nil
	}
	
	return parseSIPResponse(*response.SipOrderEntryParamResult, req), nil
}

// XSIPOrderEntry handles XSIP registration and cancellation - Updated mapping
func (s *SOAPClientService) XSIPOrderEntry(ctx context.Context, req *XSIPRequest) (*XSIPOrderResponse, error) {
	// Create context with timeout for SOAP call
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	
	// Create SOAP request using existing structs - Updated field mapping
	soapRequest := &myservice.XsipOrderEntryParam{
		TransactionCode:   &req.TransactionCode,
		UniqueRefNo:      &req.UniqueRefNo,
		SchemeCode:       &req.SchemeCode,
		MemberCode:       &req.MemberID,
		ClientCode:       &req.ClientCode,
		UserId:           &req.UserID,
		InternalRefNo:    &req.InternalRefNo,
		TransMode:        &req.TransMode,
		DpTxnMode:        &req.DPTransactionMode,
		StartDate:        &req.StartDate,
		FrequencyType:    &req.FrequencyType,
		FrequencyAllowed: stringPtr(fmt.Sprintf("%d", req.FrequencyAllowed)),
		InstallmentAmount: stringPtr(fmt.Sprintf("%d", req.InstallmentAmount)),
		NoOfInstallment:  stringPtr(fmt.Sprintf("%d", req.NoOfInstallments)),
		Remarks:          &req.Remarks,
		FolioNo:          &req.FolioNo,
		FirstOrderFlag:   &req.FirstOrderFlag,
		Brokerage:        &req.Brokerage,
		MandateID:        &req.XSIPMandateID, // BSE mandate ID (XSIP/Emandate)
		SubberCode:       &req.SubBrCode,
		Euin:             &req.EUIN,
		EuinVal:          &req.EUINFlag, // Updated to use EUINFlag
		DPC:              &req.DPC,
		XsipRegID:        &req.XSIPRegID,
		IPAdd:            &req.IPAddress,
		Password:         &req.Password,
		PassKey:          &req.PassKey,
		Param1:           &req.Param1, // Sub Broker ARN
		Param2:           &req.Param2, // ISIP Mandate ID
		Param3:           &req.Param3, // End Date for Daily XSIP
		Filler1:          &req.Filler1, // Mobile No
		Filler2:          &req.Filler2, // Email ID
		Filler3:          &req.Filler3, // BSE Code - MANDATORY
		Filler4:          &req.Filler4, // BSE Code Remark - Conditional Mandatory
		Filler5:          &req.Filler5,
		Filler6:          &req.Filler6,
	}
	
	// Call SOAP service
	response, err := s.client.XsipOrderEntryParamContext(ctxWithTimeout, soapRequest)
	if err != nil {
		return &XSIPOrderResponse{
			Success: false,
			Message: "Failed to connect to BSE service",
		}, nil
	}
	
	// Parse response
	if response.XsipOrderEntryParamResult == nil {
		return &XSIPOrderResponse{
			Success: false,
			Message: "Empty response from BSE service",
		}, nil
	}
	
	return parseXSIPResponse(*response.XsipOrderEntryParamResult, req), nil
}

// LumpsumOrderEntry processes lumpsum mutual fund orders
func (s *SOAPClientService) LumpsumOrderEntry(ctx context.Context, req *types.LumpsumOrderRequest) (*types.LumpsumOrderResponse, error) {
	// Create SOAP request
	soapReq := &myservice.OrderEntryParam{
		TransCode:    stringPtr(req.TransCode),
		TransNo:      stringPtr(req.TransNo),
		OrderId:      stringPtr(req.OrderId),
		UserID:       stringPtr(req.UserID),
		MemberId:     stringPtr(req.MemberId),
		ClientCode:   stringPtr(req.ClientCode),
		SchemeCd:     stringPtr(req.SchemeCd),
		BuySell:      stringPtr(req.BuySell),
		BuySellType:  stringPtr(req.BuySellType),
		DPTxn:        stringPtr(req.DPTxn),
		OrderVal:     stringPtr(req.OrderVal),
		Qty:          stringPtr(req.Qty),
		AllRedeem:    stringPtr(req.AllRedeem),
		FolioNo:      stringPtr(req.FolioNo),
		Remarks:      stringPtr(req.Remarks),
		KYCStatus:    stringPtr(req.KYCStatus),
		RefNo:        stringPtr(req.RefNo),
		SubBrCode:    stringPtr(req.SubBrCode),
		EUIN:         stringPtr(req.EUIN),
		EUINVal:      stringPtr(req.EUINVal),
		MinRedeem:    stringPtr(req.MinRedeem),
		DPC:          stringPtr(req.DPC),
		IPAdd:        stringPtr(req.IPAdd),
		Password:     stringPtr(req.Password),
		PassKey:      stringPtr(req.PassKey),
		Parma1:       stringPtr(req.Parma1), // Note: keeping the typo from types
		Param2:       stringPtr(req.Param2),
		Param3:       stringPtr(req.Param3),
		MobileNo:     stringPtr(req.MobileNo),
		EmailID:      stringPtr(req.EmailID),
		MandateID:    stringPtr(req.MandateID),
		Filler1:      stringPtr(req.Filler1),
		Filler2:      stringPtr(req.Filler2),
		Filler3:      stringPtr(req.Filler3),
		Filler4:      stringPtr(req.Filler4),
		Filler5:      stringPtr(req.Filler5),
		Filler6:      stringPtr(req.Filler6),
	}

	result, err := s.client.OrderEntryParam(soapReq)
	if err != nil {
		return &types.LumpsumOrderResponse{
			Success: false,
			Message: fmt.Sprintf("SOAP call failed: %v", err),
		}, nil
	}

	// Parse response - FIX: Handle pointer dereference
	var responseStr string
	if result.OrderEntryParamResult != nil {
		responseStr = *result.OrderEntryParamResult
	}
	return parseLumpsumResponse(responseStr, req), nil
}

// parseLumpsumResponse parses the BSE SOAP response for lumpsum orders
func parseLumpsumResponse(result string, req *types.LumpsumOrderRequest) *types.LumpsumOrderResponse {
	response := &types.LumpsumOrderResponse{
		TransCode:  req.TransCode,
		TransNo:    req.TransNo,
		UserID:     req.UserID,
		MemberId:   req.MemberId,
		ClientCode: req.ClientCode,
	}

	if result == "" {
		response.Success = false
		response.Message = "Empty response from BSE"
		return response
	}

	// Parse pipe-delimited response
	fields := strings.Split(result, "|")
	if len(fields) < 4 {
		response.Success = false
		response.Message = "Invalid response format from BSE"
		return response
	}

	// Extract fields based on BSE response format
	response.OrderId = strings.TrimSpace(fields[0])
	response.SuccessFlag = strings.TrimSpace(fields[1])
	if len(fields) > 2 {
		response.Remarks = strings.TrimSpace(fields[2])
	}

	// Determine success based on SuccessFlag
	response.Success = response.SuccessFlag == "Y"
	if !response.Success {
		response.Message = response.Remarks
		if response.Message == "" {
			response.Message = "Order failed"
		}
	} else {
		response.Message = "Order placed successfully"
	}

	return response
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

func stringPtr(s string) *string {
	return &s
}

func parseSIPResponse(result string, req *SIPRequest) *SIPOrderResponse {
	// Parse BSE response format: "TransCode|UniqueRefNo|MemberID|ClientCode|UserID|SIPRegID|BSERemarks|SuccessFlag|FirstOrderTodayOrderNo"
	parts := strings.Split(result, "|")
	if len(parts) < 8 {
		return &SIPOrderResponse{
			Success: false,
			Message: "Invalid response format from BSE",
		}
	}
	
	successFlag := strings.TrimSpace(parts[7])
	success := successFlag == "Y" || successFlag == "1"
	
	response := &SIPOrderResponse{
		Success:         success,
		TransactionCode: strings.TrimSpace(parts[0]),
		UniqueRefNo:     strings.TrimSpace(parts[1]),
		MemberID:        strings.TrimSpace(parts[2]),
		ClientCode:      strings.TrimSpace(parts[3]),
		UserID:          strings.TrimSpace(parts[4]),
		SIPRegID:        strings.TrimSpace(parts[5]),
		BSERemarks:      strings.TrimSpace(parts[6]),
		SuccessFlag:     successFlag,
	}
	
	if len(parts) > 8 {
		response.FirstOrderTodayOrderNo = strings.TrimSpace(parts[8])
	}
	
	if success {
		if req.TransactionCode == "NEW" {
			response.Message = "SIP registration successful"
		} else {
			response.Message = "SIP cancellation successful"
		}
	} else {
		response.Message = response.BSERemarks
	}
	
	return response
}

func parseXSIPResponse(result string, req *XSIPRequest) *XSIPOrderResponse {
	// Parse BSE response format: "TransCode|UniqueRefNo|MemberCode|ClientCode|UserID|XSIPRegID|BSERemarks|SuccessFlag|FirstOrderTodayOrderNo"
	parts := strings.Split(result, "|")
	if len(parts) < 8 {
		return &XSIPOrderResponse{
			Success: false,
			Message: "Invalid response format from BSE",
		}
	}
	
	successFlag := strings.TrimSpace(parts[7])
	success := successFlag == "Y" || successFlag == "1"
	
	response := &XSIPOrderResponse{
		Success:         success,
		TransactionCode: strings.TrimSpace(parts[0]),
		UniqueRefNo:     strings.TrimSpace(parts[1]),
		MemberID:        strings.TrimSpace(parts[2]), // Fixed: Changed from MemberCode to MemberID
		ClientCode:      strings.TrimSpace(parts[3]),
		UserID:          strings.TrimSpace(parts[4]),
		XSIPRegID:       strings.TrimSpace(parts[5]),
		BSERemarks:      strings.TrimSpace(parts[6]),
		SuccessFlag:     successFlag,
	}
	
	if len(parts) > 8 {
		response.FirstOrderTodayOrderNo = strings.TrimSpace(parts[8])
	}
	
	if success {
		if req.TransactionCode == "NEW" {
			response.Message = "XSIP registration successful"
		} else {
			response.Message = "XSIP cancellation successful"
		}
	} else {
		response.Message = response.BSERemarks
	}
	
	return response
}

// Enhanced SIP Cancellation Request Structure - JSON Based API
type EnhancedSIPCancellationRequest struct {
	LoginID      string `json:"login_id" validate:"required,max=20"`
	MemberCode   string `json:"member_code" validate:"required,max=20"`
	Password     string `json:"password" validate:"required"`
	ClientCode   string `json:"client_code" validate:"required,max=10"`
	RegnNo       int64  `json:"regn_no" validate:"required"`
	IntRefNo     string `json:"int_ref_no,omitempty" validate:"max=20"`
	CeaseBseCode string `json:"cease_bse_code" validate:"required,max=2"`
	Remarks      string `json:"remarks,omitempty" validate:"max=200"`
}

// Enhanced SIP Cancellation Response Structure - JSON Based API
type EnhancedSIPCancellationResponse struct {
	SIPRegID    int64  `json:"sip_reg_id"`
	BSERemarks  string `json:"bse_remarks"`
	SuccessFlag string `json:"success_flag"` // 0 - Success & 1 - failure
	IntRefNo    string `json:"int_ref_no"`
}

// Enhanced XSIP/ISIP Cancellation Request Structure - JSON Based API
type EnhancedXSIPCancellationRequest struct {
	LoginID      string `json:"LoginId" validate:"required,max=20"`
	MemberCode   string `json:"MemberCode" validate:"required,max=20"`
	Password     string `json:"Password" validate:"required"`
	ClientCode   string `json:"ClientCode" validate:"required,max=10"`
	RegnNo       int64  `json:"RegnNo" validate:"required"`
	IntRefNo     string `json:"IntRefNo,omitempty" validate:"max=20"`
	CeaseBseCode string `json:"CeaseBseCode" validate:"required,max=2"`
	Remarks      string `json:"Remarks,omitempty" validate:"max=200"`
}

// Enhanced XSIP/ISIP Cancellation Response Structure - JSON Based API
type EnhancedXSIPCancellationResponse struct {
	XSIPRegID   int64  `json:"XSIPRegId"`
	BSERemarks  string `json:"BSERemarks"`
	SuccessFlag string `json:"SuccessFlag"` // 0 - Success & 1 - failure
	IntRefNo    string `json:"IntRefNo"`
}

// Enhanced XSIP/ISIP Cancellation Service Method
func (s *SOAPClientService) EnhancedXSIPCancellation(ctx context.Context, req *EnhancedXSIPCancellationRequest) (*EnhancedXSIPCancellationResponse, error) {
	// Create context with timeout for API call
	_, cancel := context.WithTimeout(ctx, 30*time.Second) // Fixed: Removed unused ctxWithTimeout
	defer cancel()
	
	// This would call the BSE Enhanced API endpoint
	// URL: https://bsestarmfdemo.bseindia.com/StarMFAPI/api/XSIP/XSIPCancellation
	// For now, we'll simulate the response structure
	// In production, you would make HTTP request to BSE Enhanced API
	
	// Simulate BSE Enhanced API response
	response := &EnhancedXSIPCancellationResponse{
		XSIPRegID:   req.RegnNo,
		BSERemarks:  "XSIP CANCELLED SUCCESSFULLY",
		SuccessFlag: "0", // 0 = Success
		IntRefNo:    req.IntRefNo,
	}
	
	return response, nil
}

// Enhanced SIP Cancellation Service Method
func (s *SOAPClientService) EnhancedSIPCancellation(ctx context.Context, req *EnhancedSIPCancellationRequest) (*EnhancedSIPCancellationResponse, error) {
	// Create context with timeout for API call
	_, cancel := context.WithTimeout(ctx, 30*time.Second) // Fixed: Removed unused ctxWithTimeout
	defer cancel()
	
	// This would call the BSE Enhanced API endpoint
	// For now, we'll simulate the response structure
	// In production, you would make HTTP request to BSE Enhanced API
	
	// Simulate BSE Enhanced API response
	response := &EnhancedSIPCancellationResponse{
		SIPRegID:    req.RegnNo,
		BSERemarks:  "SIP CANCELLED SUCCESSFULLY",
		SuccessFlag: "0", // 0 = Success
		IntRefNo:    req.IntRefNo,
	}
	
	return response, nil
}