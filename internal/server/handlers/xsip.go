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

// SOAP Envelope structure for XSIP requests
type XSIPSOAPEnvelope struct {
	XMLName xml.Name `xml:"soap:Envelope"`
	SoapNS  string   `xml:"xmlns:soap,attr"`
	Body    XSIPSOAPBody `xml:"soap:Body"`
}

type XSIPSOAPBody struct {
	XSIPRequest XSIPSOAPRequest `xml:"XsipOrderEntryParam"`
}

type XSIPSOAPRequest struct {
	TransactionCode   string `xml:"TransactionCode"`
	UniqueRefNo      string `xml:"UniqueRefNo"`
	SchemeCode       string `xml:"SchemeCode"`
	MemberCode       string `xml:"MemberCode"`
	ClientCode       string `xml:"ClientCode"`
	UserId           string `xml:"UserId"`
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
	Brokerage        string `xml:"Brokerage"`
	MandateID        string `xml:"MandateID"`
	SubberCode       string `xml:"SubberCode"`
	Euin             string `xml:"Euin"`
	EuinVal          string `xml:"EuinVal"`
	DPC              string `xml:"DPC"`
	XsipRegID        string `xml:"XsipRegID"`
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

// SOAP Response structure for XSIP
type XSIPSOAPResponse struct {
	XMLName xml.Name `xml:"soap:Envelope"`
	SoapNS  string   `xml:"xmlns:soap,attr"`
	Body    XSIPSOAPResponseBody `xml:"soap:Body"`
}

type XSIPSOAPResponseBody struct {
	XSIPResult XSIPSOAPResult `xml:"XsipOrderEntryParamResponse"`
}

type XSIPSOAPResult struct {
	Result string `xml:"XsipOrderEntryParamResult"`
}

func XSIPHandler(w http.ResponseWriter, r *http.Request) {
	// Set SOAP XML content type
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	
	// CREATE logger and SOAP service here
	logger := util.NewStandardLogger()
	xsipSoapService, err := services.NewSOAPClientService(logger)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		soapFault := createXSIPSOAPFault("Failed to initialize SOAP service")
		xml.NewEncoder(w).Encode(soapFault)
		return
	}
	
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		soapFault := createXSIPSOAPFault("Method not allowed")
		xml.NewEncoder(w).Encode(soapFault)
		return
	}
	
	// Check if SOAP service is available
	if xsipSoapService == nil {
		w.WriteHeader(http.StatusInternalServerError)
		soapFault := createXSIPSOAPFault("SOAP service unavailable")
		xml.NewEncoder(w).Encode(soapFault)
		return
	}
	
	// Read SOAP XML body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		soapFault := createXSIPSOAPFault("Failed to read request body")
		xml.NewEncoder(w).Encode(soapFault)
		return
	}
	
	// Parse SOAP XML
	var soapEnv XSIPSOAPEnvelope
	err = xml.Unmarshal(body, &soapEnv)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		soapFault := createXSIPSOAPFault("Invalid SOAP XML: " + err.Error())
		xml.NewEncoder(w).Encode(soapFault)
		return
	}
	
	// Convert SOAP request to internal XSIP request
	req := convertSOAPToXSIPRequest(&soapEnv.Body.XSIPRequest)
	
	// Validate the request
	if err := validateXSIPRequest(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		soapFault := createXSIPSOAPFault("Validation error: " + err.Error())
		xml.NewEncoder(w).Encode(soapFault)
		return
	}
	
	// Call BSE SOAP service
	xsipResp, err := xsipSoapService.XSIPOrderEntry(r.Context(), &req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		soapFault := createXSIPSOAPFault("Failed to process XSIP request with BSE service")
		xml.NewEncoder(w).Encode(soapFault)
		return
	}
	
	// Convert response to SOAP XML
	soapResponse := convertXSIPResponseToSOAP(xsipResp)
	
	// Return SOAP response
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(soapResponse)
}

// Helper function to convert SOAP request to internal XSIP format
func convertSOAPToXSIPRequest(soapReq *XSIPSOAPRequest) services.XSIPRequest {
	return services.XSIPRequest{
		TransactionCode:    soapReq.TransactionCode,
		UniqueRefNo:       soapReq.UniqueRefNo,
		SchemeCode:        soapReq.SchemeCode,
		MemberID:          soapReq.MemberCode,
		ClientCode:        soapReq.ClientCode,
		UserID:            soapReq.UserId,
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
		Brokerage:         soapReq.Brokerage,
		XSIPMandateID:     soapReq.MandateID,
		SubBrCode:         soapReq.SubberCode,
		EUIN:              soapReq.Euin,
		EUINFlag:          soapReq.EuinVal,
		DPC:               soapReq.DPC,
		XSIPRegID:         soapReq.XsipRegID,
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

// Helper function to convert XSIP response to SOAP
func convertXSIPResponseToSOAP(xsipResp *services.XSIPOrderResponse) XSIPSOAPResponse {
	// Format response as pipe-delimited string (BSE format)
	result := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s",
		xsipResp.TransactionCode,
		xsipResp.UniqueRefNo,
		xsipResp.MemberID,
		xsipResp.ClientCode,
		xsipResp.UserID,
		xsipResp.XSIPRegID,
		xsipResp.BSERemarks,
		xsipResp.SuccessFlag)
	
	return XSIPSOAPResponse{
		SoapNS: "http://schemas.xmlsoap.org/soap/envelope/",
		Body: XSIPSOAPResponseBody{
			XSIPResult: XSIPSOAPResult{
				Result: result,
			},
		},
	}
}

// Helper function for XSIP SOAP faults
func createXSIPSOAPFault(message string) interface{} {
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

// DELETE lines 235-243 (the duplicate function):
// // Helper function to parse integers (shared with SIP handler)
// func parseIntOrDefault(s string, defaultVal int) int {
// 	if s == "" {
// 		return defaultVal
// 	}
// 	// Add proper string to int conversion
// 	val := 0
// 	fmt.Sscanf(s, "%d", &val)
// 	return val
// }

// Keep existing validateXSIPRequest function
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
	
	// Validate BSE Code (Filler3) - MANDATORY for XSIP
	if strings.TrimSpace(req.Filler3) == "" {
		return fmt.Errorf("BSE Code (filler3) is mandatory")
	}
	
	// Validate BSE Code Remark (Filler4) - Conditional Mandatory
	if req.Filler3 == "13" && strings.TrimSpace(req.Filler4) == "" {
		return fmt.Errorf("BSE Code Remark (filler4) is mandatory when BSE Code is '13' (Others)")
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
	if len(req.Filler3) > 2 {
		return fmt.Errorf("BSE Code must not exceed 2 characters")
	}
	if len(req.Filler4) > 200 {
		return fmt.Errorf("BSE Code Remark must not exceed 200 characters")
	}
	
	return nil
}
