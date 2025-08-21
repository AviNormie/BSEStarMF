package bse

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// LumpsumOrderClient handles BSE StAR MF Lumpsum Order Entry API calls
type LumpsumOrderClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewLumpsumOrderClient creates a new Lumpsum Order client
func NewLumpsumOrderClient() *LumpsumOrderClient {
	return &LumpsumOrderClient{
		BaseURL: "https://bsestarmfdemo.bseindia.com/MFOrderEntry/MFOrder.svc", // ✅ Real BSE URL
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// PlaceOrder places a lumpsum order with BSE
func (c *LumpsumOrderClient) PlaceOrder(req LumpsumOrderRequest) (*LumpsumOrderResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Construct SOAP request
	soapPayload := c.buildSOAPPayload(req)

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", c.BaseURL, bytes.NewBufferString(soapPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")
	httpReq.Header.Set("SOAPAction", "http://bsestarmf.in/MFOrderEntry/orderEntryParam")

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

	// Parse SOAP response
	orderResp, err := c.parseSOAPResponse(string(respBody))
	if err != nil {
		return nil, fmt.Errorf("failed to parse SOAP response: %w", err)
	}

	return orderResp, nil
}

// buildSOAPPayload constructs the SOAP request payload
func (c *LumpsumOrderClient) buildSOAPPayload(req LumpsumOrderRequest) string {
	// Handle optional OrderID
	orderID := ""
	if req.OrderID != nil {
		orderID = fmt.Sprintf("%d", *req.OrderID)
	}

	// Handle optional Amount and Qty
	amount := ""
	if req.Amount != nil {
		amount = fmt.Sprintf("%.2f", *req.Amount)
	}

	qty := ""
	if req.Qty != nil {
		qty = fmt.Sprintf("%.2f", *req.Qty)
	}

	return fmt.Sprintf(`
		<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope" xmlns:bses="http://bsestarmf.in/">
		  <soap:Header xmlns:wsa="http://www.w3.org/2005/08/addressing">
		    <wsa:Action>http://bsestarmf.in/MFOrderEntry/orderEntryParam</wsa:Action>
		    <wsa:To>%s</wsa:To>
		  </soap:Header>
		  <soap:Body>
		    <bses:orderEntryParam>
		      <bses:TransCode>%s</bses:TransCode>
		      <bses:TransNo>%s</bses:TransNo>
		      <bses:OrderId>%s</bses:OrderId>
		      <bses:UserId>%d</bses:UserId>
		      <bses:MemberId>%s</bses:MemberId>
		      <bses:ClientCode>%s</bses:ClientCode>
		      <bses:SchemeCd>%s</bses:SchemeCd>
		      <bses:BuySell>%s</bses:BuySell>
		      <bses:BuySellType>%s</bses:BuySellType>
		      <bses:DPTxn>%s</bses:DPTxn>
		      <bses:OrderVal>%s</bses:OrderVal>
		      <bses:Qty>%s</bses:Qty>
		      <bses:AllRedeem>%s</bses:AllRedeem>
		      <bses:FolioNo>%s</bses:FolioNo>
		      <bses:Remarks>%s</bses:Remarks>
		      <bses:KYCStatus>%s</bses:KYCStatus>
		      <bses:RefNo>%s</bses:RefNo>
		      <bses:SubBrCode>%s</bses:SubBrCode>
		      <bses:EUIN>%s</bses:EUIN>
		      <bses:EUINFlag>%s</bses:EUINFlag>
		      <bses:MinRedeem>%s</bses:MinRedeem>
		      <bses:DPC>%s</bses:DPC>
		      <bses:IPAdd>%s</bses:IPAdd>
		      <bses:Password>%s</bses:Password>
		      <bses:PassKey>%s</bses:PassKey>
		      <bses:Param1>%s</bses:Param1>
		      <bses:Param2>%s</bses:Param2>
		      <bses:Param3>%s</bses:Param3>
		      <bses:MobileNo>%s</bses:MobileNo>
		      <bses:EmailID>%s</bses:EmailID>
		      <bses:MandateID>%s</bses:MandateID>
		      <bses:Filler1>%s</bses:Filler1>
		      <bses:Filler2>%s</bses:Filler2>
		      <bses:Filler3>%s</bses:Filler3>
		      <bses:Filler4>%s</bses:Filler4>
		      <bses:Filler5>%s</bses:Filler5>
		      <bses:Filler6>%s</bses:Filler6>
		    </bses:orderEntryParam>
		  </soap:Body>
		</soap:Envelope>`,
		c.BaseURL,
		req.TransCode, req.TransNo, orderID, req.UserID, req.MemberID, req.ClientCode,
		req.SchemeCd, req.BuySell, req.BuySellType, req.DPTxn, amount, qty,
		req.AllRedeem, req.FolioNo, req.Remarks, req.KYCStatus, req.RefNo,
		req.SubBrCode, req.EUIN, req.EUINFlag, req.MinRedeem, req.DPC,
		req.IPAdd, req.Password, req.PassKey, req.Param1, req.Param2, req.Param3,
		req.MobileNo, req.EmailID, req.MandateID, req.Filler1, req.Filler2,
		req.Filler3, req.Filler4, req.Filler5, req.Filler6)
}

// parseSOAPResponse parses the SOAP response from BSE
func (c *LumpsumOrderClient) parseSOAPResponse(responseBody string) (*LumpsumOrderResponse, error) {
	// Find orderEntryParamResult in the response
	start := strings.Index(responseBody, "<orderEntryParamResult>")
	end := strings.Index(responseBody, "</orderEntryParamResult>")
	if start == -1 || end == -1 {
		return nil, fmt.Errorf("cannot find orderEntryParamResult in response")
	}

	result := responseBody[start+len("<orderEntryParamResult>"):end]

	// Parse the pipe-separated response
	// Format: TransactionCode|UniqueRefNo|OrderNumber|UserID|MemberID|ClientCode|BSERemarks|SuccessFlag
	parts := strings.Split(result, "|")
	if len(parts) < 8 {
		return nil, fmt.Errorf("invalid response format: expected 8 parts, got %d", len(parts))
	}

	// Parse UserID and OrderNumber as integers
	userID := int64(0)
	if parts[3] != "" {
		fmt.Sscanf(parts[3], "%d", &userID)
	}

	orderNumber := int64(0)
	if parts[2] != "" {
		fmt.Sscanf(parts[2], "%d", &orderNumber)
	}

	return &LumpsumOrderResponse{
		TransactionCode:       parts[0],
		UniqueRefNo:          parts[1],
		OrderNumber:          orderNumber,
		UserID:               userID,
		MemberID:             parts[4],
		ClientCode:           parts[5],
		BSERemarks:           parts[6],
		SuccessFlag:          parts[7],
		OrderEntryParamResult: result,
	}, nil
}

// PlaceOrderTest places a test order without calling BSE (for testing)
func (c *LumpsumOrderClient) PlaceOrderTest(req LumpsumOrderRequest) (*LumpsumOrderResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Return mock successful response
	return &LumpsumOrderResponse{
		TransactionCode:       req.TransCode,
		UniqueRefNo:          req.TransNo,
		OrderNumber:          123456789,
		UserID:               req.UserID,
		MemberID:             req.MemberID,
		ClientCode:           req.ClientCode,
		BSERemarks:           "Test order placed successfully",
		SuccessFlag:          FlagYes,
		OrderEntryParamResult: fmt.Sprintf("%s|%s|123456789|%d|%s|%s|Test order placed successfully|Y",
			req.TransCode, req.TransNo, req.UserID, req.MemberID, req.ClientCode),
	}, nil
}