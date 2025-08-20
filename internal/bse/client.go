package bse

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
)

// Authenticate calls BSE SOAP API and parses the response, as per BSE docs
func Authenticate(userID, password, passKey string) (string, string, error) {
    // Construct SOAP request based on your docs
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
        </soap:Envelope>`, userID, password, passKey)

    // Send SOAP request
    resp, err := http.Post(
        "https://bsestarmfdemo.bseindia.com/MFOrderEntry/MFOrder.svc/Secure",
        "application/soap+xml; charset=utf-8",
        bytes.NewBufferString(soapPayload),
    )
    if err != nil {
        return "", "", err
    }
    defer resp.Body.Close()

    // Read response as string
    responseBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", "", fmt.Errorf("failed to read response: %v", err)
    }
    raw := string(responseBody)

    // Parse <getPasswordResult>Code|EncryptedPassword</getPasswordResult>
    start := strings.Index(raw, "<getPasswordResult>")
    end := strings.Index(raw, "</getPasswordResult>")
    if start == -1 || end == -1 {
        return "", "", fmt.Errorf("cannot find getPasswordResult in response")
    }
    result := raw[start+len("<getPasswordResult>") : end]
    parts := strings.SplitN(result, "|", 2)
    code, encrypted := "101", ""
    if len(parts) == 2 {
        code = parts[0]
        encrypted = parts[1]
    }
    return code, encrypted, nil
}
