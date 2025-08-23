package types

// AuthResponse represents the response from BSE authentication service
type AuthResponse struct {
	ResponseCode     string `json:"response_code"`
	EncryptedPassword string `json:"encrypted_password,omitempty"`
	ErrorMessage     string `json:"error_message,omitempty"`
}