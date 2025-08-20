package bse

// UCCRegistrationRequest represents the request structure for BSE StAR MF Enhanced UCC Registration API
type UCCRegistrationRequest struct {
	UserID     string `json:"UserId"`
	MemberCode string `json:"MemberCode"`
	Password   string `json:"Password"`
	RegnType   string `json:"RegnType"`   // ✅ Supports both "NEW" and "MOD"
	Param      string `json:"Param"`
	Filler1    string `json:"Filler1,omitempty"`
	Filler2    string `json:"Filler2,omitempty"`
}

// UCCRegistrationResponse represents the response structure from BSE StAR MF Enhanced UCC Registration API
type UCCRegistrationResponse struct {
	Status  string `json:"Status"`   // 0 - Success & 1 - Failure
	Remarks string `json:"Remarks"`  // Return Remarks
	Filler1 string `json:"Filler1"`  // ✅ Always included
	Filler2 string `json:"Filler2"`  // ✅ Always included
}

// Constants for registration types and status
const (
	RegnTypeNew    = "NEW"
	RegnTypeModify = "MOD"
	StatusSuccess  = "0"
	StatusFailure  = "1"
)

// IsSuccessResponse checks if the registration was successful
func (r *UCCRegistrationResponse) IsSuccessResponse() bool {
	return r.Status == StatusSuccess
}