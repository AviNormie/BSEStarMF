package handlers

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sapphirebroking.com/sapphire_mf/internal/server/services"
	"sapphirebroking.com/sapphire_mf/internal/util"
	
)

// SOAP SIP Cancellation Request
type SOAPSIPCancellationEnvelope struct {
	XMLName xml.Name `xml:"soap:Envelope"`
	SoapNS  string   `xml:"xmlns:soap,attr"`
	Body    SOAPSIPCancellationBody `xml:"soap:Body"`
}

type SOAPSIPCancellationBody struct {
	Request SOAPSIPCancellationRequest `xml:"SipCancellation"`
}

type SOAPSIPCancellationRequest struct {
	TransactionCode string `xml:"TransactionCode"`
	UniqueRefNo    string `xml:"UniqueRefNo"`
	SchemeCode     string `xml:"SchemeCode"`
	MemberCode     string `xml:"MemberCode"`
	ClientCode     string `xml:"ClientCode"`
	UserID         string `xml:"UserID"`
	RegId          string `xml:"RegId"`
	Password       string `xml:"Password"`
	PassKey        string `xml:"PassKey"`
	Filler2        string `xml:"Filler2"`
	// Add other required fields as needed
}

// SOAP SIP Cancellation Response
type SOAPSIPCancellationResponse struct {
	XMLName xml.Name `xml:"soap:Envelope"`
	SoapNS  string   `xml:"xmlns:soap,attr"`
	Body    SOAPSIPCancellationResponseBody `xml:"soap:Body"`
}

type SOAPSIPCancellationResponseBody struct {
	Response SOAPSIPCancellationResult `xml:"SipCancellationResponse"`
}

type SOAPSIPCancellationResult struct {
	Result string `xml:"SipCancellationResult"`
}

func SOAPSIPCancellationHandler(w http.ResponseWriter, r *http.Request) {
	// Set SOAP content type
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Header().Set("SOAPAction", "\"SipCancellation\"")
	
	// Create SOAP service
	logger := util.NewStandardLogger()
	soapService, err := services.NewSOAPClientService(logger)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		soapFault := createSOAPFault("Failed to initialize SOAP service")
		xml.NewEncoder(w).Encode(soapFault)
		return
	}
	
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		soapFault := createSOAPFault("Only POST method allowed")
		xml.NewEncoder(w).Encode(soapFault)
		return
	}
	
	// Read SOAP XML
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		soapFault := createSOAPFault("Failed to read request body")
		xml.NewEncoder(w).Encode(soapFault)
		return
	}
	
	// Parse SOAP envelope
	var soapEnv SOAPSIPCancellationEnvelope
	err = xml.Unmarshal(body, &soapEnv)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		soapFault := createSOAPFault("Invalid SOAP XML: " + err.Error())
		xml.NewEncoder(w).Encode(soapFault)
		return
	}
	
	// Convert to internal SIP request
	sipReq := convertSOAPToSIPCancellationRequest(&soapEnv.Body.Request)
	
	// Call BSE SOAP service
	sipResp, err := soapService.SIPOrderEntry(r.Context(), &sipReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		soapFault := createSOAPFault("BSE service error: " + err.Error())
		xml.NewEncoder(w).Encode(soapFault)
		return
	}
	
	// Convert response to SOAP
	soapResponse := convertSIPCancellationResponseToSOAP(sipResp)
	
	// Return SOAP response
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(soapResponse)
}

func convertSOAPToSIPCancellationRequest(soapReq *SOAPSIPCancellationRequest) services.SIPRequest {
	return services.SIPRequest{
		TransactionCode: "CXL", // Force cancellation
		UniqueRefNo:    soapReq.UniqueRefNo,
		SchemeCode:     soapReq.SchemeCode,
		MemberID:       soapReq.MemberCode,
		ClientCode:     soapReq.ClientCode,
		UserID:         soapReq.UserID,
		RegID:          soapReq.RegId,
		Password:       soapReq.Password,
		PassKey:        soapReq.PassKey,
		Filler2:        soapReq.Filler2,
		// Set required defaults for cancellation
		TransMode:         "D",
		DPTransactionMode: "C",
		StartDate:         "01/01/2024",
		FrequencyType:     "MONTHLY",
		FrequencyAllowed:  1,
		InstallmentAmount: 1000,
		NoOfInstallments:  12,
		FirstOrderFlag:    "Y",
		EUIN:              "TEST",
		EUINDeclaration:   "Y",
		DPC:               "Y",
	}
}

func convertSIPCancellationResponseToSOAP(sipResp *services.SIPOrderResponse) SOAPSIPCancellationResponse {
	result := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s",
		sipResp.TransactionCode,
		sipResp.UniqueRefNo,
		sipResp.MemberID,
		sipResp.ClientCode,
		sipResp.UserID,
		sipResp.SIPRegID,
		sipResp.BSERemarks,
		sipResp.SuccessFlag)
	
	return SOAPSIPCancellationResponse{
		SoapNS: "http://schemas.xmlsoap.org/soap/envelope/",
		Body: SOAPSIPCancellationResponseBody{
			Response: SOAPSIPCancellationResult{
				Result: result,
			},
		},
	}
}

// Remove this entire function (it's already in sip.go):
// func createSOAPFault(message string) interface{} {
//     return map[string]interface{}{
//         "soap:Envelope": map[string]interface{}{
//             "xmlns:soap": "http://schemas.xmlsoap.org/soap/envelope/",
//             "soap:Body": map[string]interface{}{
//                 "soap:Fault": map[string]interface{}{
//                     "faultcode":   "Server",
//                     "faultstring": message,
//                 },
//             },
//         },
//     }
// }