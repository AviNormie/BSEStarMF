package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type AuthRequest struct {
    UserID   string `json:"user_id"`
    Password string `json:"password"`
    PassKey  string `json:"pass_key"`
}

type AuthResponse struct {
    Code              string `json:"code"`
    EncryptedPassword string `json:"encrypted_password"`
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
    var req AuthRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Create SOAP Payload 
    soapPayload := fmt.Sprintf(`
        <soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope" xmlns:bses="http://bsestarmf.in/">
          <soap:Header xmlns:wsa="http://www.w3.org/2005/08/addressing">
            <wsa:Action>http://bsestarmf.in/MFOrderEntry/getPassword</wsa:Action>
            <wsa:To>https://bsestarmfdemo.bseindia.com/MFOrderEntry/MFOrder.svc/Secure</wsa:To>
          </soap:Header>
          <soap:Body>
            <bses:getPassword>
              <bses:UserId>%s</bses:UserId>
              <bses:Password>%s</bses:Password>
              <bses:PassKey>%s</bses:PassKey>
            </bses:getPassword>
          </soap:Body>
        </soap:Envelope>`, req.UserID, req.Password, req.PassKey)

    // Send SOAP Request
    resp, err := http.Post(
        "https://bsestarmfdemo.bseindia.com/MFOrderEntry/MFOrder.svc/Secure",
        "application/soap+xml; charset=utf-8",
        bytes.NewBufferString(soapPayload),
    )
    if err != nil {
        http.Error(w, "Request to BSE failed: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, "Failed to read response body: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Extract Response (structure: "Code|EncryptedPassword")
    // Find getPasswordResult and parse it
    raw := string(body)
    start := strings.Index(raw, "<getPasswordResult>")
    end := strings.Index(raw, "</getPasswordResult>")
    if start == -1 || end == -1 {
        http.Error(w, "SOAP response parsing failed", http.StatusInternalServerError)
        return
    }
    result := raw[start+len("<getPasswordResult>") : end]
    parts := strings.SplitN(result, "|", 2)
    code, password := "101", ""
    if len(parts) == 2 {
        code = parts[0]
        password = parts[1]
    }

    // Return as JSON
    respJson := AuthResponse{Code: code, EncryptedPassword: password}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(respJson)
}
