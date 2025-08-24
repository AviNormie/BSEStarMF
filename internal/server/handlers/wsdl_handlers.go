package handlers

import (
	"net/http"
)

// SIP Cancellation WSDL Handler
func SIPCancellationWSDLHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	
	wsdl := `<?xml version="1.0" encoding="utf-8"?>
<definitions xmlns="http://schemas.xmlsoap.org/wsdl/"
             targetNamespace="http://localhost:8080/soap/"
             xmlns:tns="http://localhost:8080/soap/">
  <message name="SipCancellationRequest">
    <part name="parameters" element="tns:SipCancellation"/>
  </message>
  <message name="SipCancellationResponse">
    <part name="parameters" element="tns:SipCancellationResponse"/>
  </message>
  <portType name="SipCancellationPortType">
    <operation name="SipCancellation">
      <input message="tns:SipCancellationRequest"/>
      <output message="tns:SipCancellationResponse"/>
    </operation>
  </portType>
  <binding name="SipCancellationBinding" type="tns:SipCancellationPortType">
    <soap:binding transport="http://schemas.xmlsoap.org/soap/http"/>
    <operation name="SipCancellation">
      <soap:operation soapAction="SipCancellation"/>
      <input><soap:body use="literal"/></input>
      <output><soap:body use="literal"/></output>
    </operation>
  </binding>
  <service name="SipCancellationService">
    <port name="SipCancellationPort" binding="tns:SipCancellationBinding">
      <soap:address location="http://localhost:8080/soap/SipCancellation"/>
    </port>
  </service>
</definitions>`
	
	w.Write([]byte(wsdl))
}

// XSIP Cancellation WSDL Handler
func XSIPCancellationWSDLHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	
	wsdl := `<?xml version="1.0" encoding="utf-8"?>
<definitions xmlns="http://schemas.xmlsoap.org/wsdl/"
             targetNamespace="http://localhost:8080/soap/"
             xmlns:tns="http://localhost:8080/soap/">
  <message name="XsipCancellationRequest">
    <part name="parameters" element="tns:XsipCancellation"/>
  </message>
  <message name="XsipCancellationResponse">
    <part name="parameters" element="tns:XsipCancellationResponse"/>
  </message>
  <portType name="XsipCancellationPortType">
    <operation name="XsipCancellation">
      <input message="tns:XsipCancellationRequest"/>
      <output message="tns:XsipCancellationResponse"/>
    </operation>
  </portType>
  <binding name="XsipCancellationBinding" type="tns:XsipCancellationPortType">
    <soap:binding transport="http://schemas.xmlsoap.org/soap/http"/>
    <operation name="XsipCancellation">
      <soap:operation soapAction="XsipCancellation"/>
      <input><soap:body use="literal"/></input>
      <output><soap:body use="literal"/></output>
    </operation>
  </binding>
  <service name="XsipCancellationService">
    <port name="XsipCancellationPort" binding="tns:XsipCancellationBinding">
      <soap:address location="http://localhost:8080/soap/XsipCancellation"/>
    </port>
  </service>
</definitions>`
	
	w.Write([]byte(wsdl))
}