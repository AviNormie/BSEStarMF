package handlers

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sapphirebroking.com/sapphire_mf/internal/server/services"
	"sapphirebroking.com/sapphire_mf/internal/util"
	"strings"
)

// SOAP Envelope structure for SIP requests
type SIPSOAPEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	SoapNS  string   `xml:"xmlns:soap,attr"`
	Body    SIPSOAPBody `xml:"Body"`
}

type SIPSOAPBody struct {
	SIPRequest SIPSOAPRequest `xml:"SipOrderEntryParam"`
}

type SIPSOAPRequest struct {
	TransactionCode   string `xml:"TransactionCode"`
	UniqueRefNo      string `xml:"UniqueRefNo"`
	SchemeCode       string `xml:"SchemeCode"`
	MemberCode       string `xml:"MemberCode"`
	ClientCode       string `xml:"ClientCode"`
	UserID           string `xml:"UserID"`
	InternalRefNo    string `xml:"InternalRefNo"`
	TransMode        string `xml:"TransMode"`
	DpTxnMode        string `xml:"DpTxnMode"`
	StartDate        string `xml:"StartDate"`
	FrequencyType    string `xml:"FrequencyType"`
	FrequencyAllowed string `xml:"FrequencyAllowed"`
	InstallmentAmount string `xml:"InstallmentAmount"`
	NoOfInstallment  string `xml:"NoOfInstallment"`
	Remarks          string `xml:"Remarks"`
	FolioNo          string `xml:"FolioNo"`
	FirstOrderFlag   string `xml:"FirstOrderFlag"`
	SubberCode       string `xml:"SubberCode"`
	Euin             string `xml:"Euin"`
	EuinVal          string `xml:"EuinVal"`
	DPC              string `xml:"DPC"`
	RegId            string `xml:"RegId"`
	IPAdd            string `xml:"IPAdd"`
	Password         string `xml:"Password"`
	PassKey          string `xml:"PassKey"`
	Param1           string `xml:"Param1"`
	Param2           string `xml:"Param2"`
	Param3           string `xml:"Param3"`
	Filler1          string `xml:"Filler1"`
	Filler2          string `xml:"Filler2"`
	Filler3          string `xml:"Filler3"`
	Filler4          string `xml:"Filler4"`
	Filler5          string `xml:"Filler5"`
	Filler6          string `xml:"Filler6"`
}

// SOAP Response structure
type SIPSOAPResponse struct {
	XMLName xml.Name `xml:"soap:Envelope"`
	SoapNS  string   `xml:"xmlns:soap,attr"`
	Body    SIPSOAPResponseBody `xml:"soap:Body"`
}

type SIPSOAPResponseBody struct {
	SIPResult SIPSOAPResult `xml:"SipOrderEntryParamResponse"`
}

type SIPSOAPResult struct {
	Result string `xml:"SipOrderEntryParamResult"`
}

func SIPHandler(w http.ResponseWriter, r *http.Request) {
	// Create logger and SOAP service within the handler
	logger := util.NewStandardLogger()
	logger.Info("[DEBUG] SIP Handler called")
	
	sipSoapService, err := services.NewSOAPClientService(logger)
	if err != nil {
		logger.Error("[ERROR] Failed to initialize SOAP service: %v", err)
		http.Error(w, "Failed to initialize SOAP service", http.StatusInternalServerError)
		return
	}
	
	// Set SOAP XML content type
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	
	if r.Method != http.MethodPost {
		logger.Error("[ERROR] Method not allowed: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		soapFault := createSOAPFault("Method not allowed")
		xml.NewEncoder(w).Encode(soapFault)
		return
	}
	
	// Check if SOAP service is available
	if sipSoapService == nil {
		logger.Error("[ERROR] SOAP service unavailable")
		w.WriteHeader(http.StatusInternalServerError)
		soapFault := createSOAPFault("SOAP service unavailable")
		xml.NewEncoder(w).Encode(soapFault)
		return
	}
	
	// Read SOAP XML body
	logger.Info("[DEBUG] Reading SOAP XML body")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("[ERROR] Failed to read request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		soapFault := createSOAPFault("Failed to read request body")
		xml.NewEncoder(w).Encode(soapFault)
		return
	}
	
	logger.Info("[DEBUG] Request body length: %d bytes", len(body))
	logger.Info("[DEBUG] Request body preview: %.200s", string(body))
	
	// Parse SOAP XML
	logger.Info("[DEBUG] Parsing SOAP XML")
	var soapEnv SIPSOAPEnvelope
	err = xml.Unmarshal(body, &soapEnv)
	if err != nil {
		logger.Error("[ERROR] Invalid SOAP XML: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		soapFault := createSOAPFault("Invalid SOAP XML: " + err.Error())
		xml.NewEncoder(w).Encode(soapFault)
		return
	}
	
	logger.Info("[DEBUG] SOAP XML parsed successfully")
	
	// Convert SOAP request to internal SIP request
	logger.Info("[DEBUG] Converting SOAP to SIP request")
	req := convertSOAPToSIPRequest(&soapEnv.Body.SIPRequest)
	
	// Validate the request
	logger.Info("[DEBUG] Validating SIP request")
	if err := validateSIPRequest(&req); err != nil {
		logger.Error("[ERROR] Validation error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		soapFault := createSOAPFault("Validation error: " + err.Error())
		xml.NewEncoder(w).Encode(soapFault)
		return
	}
	
	logger.Info("[DEBUG] Validation passed, calling BSE SOAP service")
	
	// Call BSE SOAP service
	sipResp, err := sipSoapService.SIPOrderEntry(r.Context(), &req)
	if err != nil {
		logger.Error("[ERROR] BSE SOAP service error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		soapFault := createSOAPFault("Failed to process SIP request with BSE service")
		xml.NewEncoder(w).Encode(soapFault)
		return
	}
	
	logger.Info("[DEBUG] BSE SOAP service call successful")
	
	// Convert response to SOAP XML
	soapResponse := convertSIPResponseToSOAP(sipResp)
	
	// Return SOAP response
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(soapResponse)
	
	logger.Info("[DEBUG] SOAP response sent successfully")
}

// Helper function to convert SOAP request to internal format
func convertSOAPToSIPRequest(soapReq *SIPSOAPRequest) services.SIPRequest {
	return services.SIPRequest{
		TransactionCode:    soapReq.TransactionCode,
		UniqueRefNo:       soapReq.UniqueRefNo,
		SchemeCode:        soapReq.SchemeCode,
		MemberID:          soapReq.MemberCode,
		ClientCode:        soapReq.ClientCode,
		UserID:            soapReq.UserID,
		InternalRefNo:     soapReq.InternalRefNo,
		TransMode:         soapReq.TransMode,
		DPTransactionMode: soapReq.DpTxnMode,
		StartDate:         soapReq.StartDate,
		FrequencyType:     soapReq.FrequencyType,
		// Convert string to int for numeric fields
		FrequencyAllowed:  parseIntOrDefault(soapReq.FrequencyAllowed, 1),
		InstallmentAmount: parseIntOrDefault(soapReq.InstallmentAmount, 0),
		NoOfInstallments:  parseIntOrDefault(soapReq.NoOfInstallment, 0),
		Remarks:           soapReq.Remarks,
		FolioNo:           soapReq.FolioNo,
		FirstOrderFlag:    soapReq.FirstOrderFlag,
		SubBrCode:         soapReq.SubberCode,
		EUIN:              soapReq.Euin,
		EUINDeclaration:   soapReq.EuinVal,
		DPC:               soapReq.DPC,
		RegID:             soapReq.RegId,
		IPAddress:         soapReq.IPAdd,
		Password:          soapReq.Password,
		PassKey:           soapReq.PassKey,
		Param1:            soapReq.Param1,
		Param2:            soapReq.Param2,
		Param3:            soapReq.Param3,
		Filler1:           soapReq.Filler1,
		Filler2:           soapReq.Filler2,
		Filler3:           soapReq.Filler3,
		Filler4:           soapReq.Filler4,
		Filler5:           soapReq.Filler5,
		Filler6:           soapReq.Filler6,
	}
}

// Helper function to convert SIP response to SOAP
func convertSIPResponseToSOAP(sipResp *services.SIPOrderResponse) SIPSOAPResponse {
	// Format response as pipe-delimited string (BSE format)
	result := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s",
		sipResp.TransactionCode,
		sipResp.UniqueRefNo,
		sipResp.MemberID,
		sipResp.ClientCode,
		sipResp.UserID,
		sipResp.SIPRegID,
		sipResp.BSERemarks,
		sipResp.SuccessFlag)
	
	return SIPSOAPResponse{
		SoapNS: "http://schemas.xmlsoap.org/soap/envelope/",
		Body: SIPSOAPResponseBody{
			SIPResult: SIPSOAPResult{
				Result: result,
			},
		},
	}
}

// Helper functions
func parseIntOrDefault(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	// Add proper string to int conversion
	val := 0
	fmt.Sscanf(s, "%d", &val)
	return val
}

func createSOAPFault(message string) interface{} {
	return map[string]interface{}{
		"soap:Envelope": map[string]interface{}{
			"xmlns:soap": "http://schemas.xmlsoap.org/soap/envelope/",
			"soap:Body": map[string]interface{}{
				"soap:Fault": map[string]interface{}{
					"faultcode":   "Client",
					"faultstring": message,
				},
			},
		},
	}
}

// Keep existing validateSIPRequest function unchanged
func validateSIPRequest(req *services.SIPRequest) error {
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
	
	// Validate BSE Code (Filler2) - MANDATORY
	if strings.TrimSpace(req.Filler2) == "" {
		return fmt.Errorf("BSE Code (filler2) is mandatory")
	}
	
	// Validate BSE Code Remark (Filler3) - Conditional Mandatory
	if req.Filler2 == "13" && strings.TrimSpace(req.Filler3) == "" {
		return fmt.Errorf("BSE Code Remark (filler3) is mandatory when BSE Code is '13' (Others)")
	}
	
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
	if len(req.Filler2) > 2 {
		return fmt.Errorf("BSE Code must not exceed 2 characters")
	}
	if len(req.Filler3) > 200 {
		return fmt.Errorf("BSE Code Remark must not exceed 200 characters")
	}
	
	return nil
}