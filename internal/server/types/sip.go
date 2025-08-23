package types

// SIPRequest structure for both NEW and CXL operations
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
	EUIN              string `json:"euin,omitempty" validate:"max=15"`
	EUINDeclaration   string `json:"euin_declaration,omitempty" validate:"oneof=Y N"`
	DPC               string `json:"dpc,omitempty" validate:"max=10"`
	RegID             string `json:"reg_id,omitempty" validate:"max=20"`
	IPAddress         string `json:"ip_address,omitempty" validate:"max=15"`
	Password          string `json:"password,omitempty" validate:"max=10"`
	PassKey           string `json:"pass_key,omitempty" validate:"max=10"`
	Param1            string `json:"param1,omitempty" validate:"max=50"`
	Param2            string `json:"param2,omitempty" validate:"max=50"`
	Param3            string `json:"param3,omitempty" validate:"max=50"`
	Filler1           string `json:"filler1,omitempty" validate:"max=50"`
	Filler2           string `json:"filler2,omitempty" validate:"max=50"`
	Filler3           string `json:"filler3,omitempty" validate:"max=50"`
	Filler4           string `json:"filler4,omitempty" validate:"max=50"`
	Filler5           string `json:"filler5,omitempty" validate:"max=50"`
	Filler6           string `json:"filler6,omitempty" validate:"max=50"`
}

// SIPOrderResponse represents the response from SIP order entry
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

// EnhancedSIPCancellationRequest for enhanced SIP cancellation
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

// EnhancedSIPCancellationResponse for enhanced SIP cancellation response
type EnhancedSIPCancellationResponse struct {
	SIPRegID    int64  `json:"sip_reg_id"`
	BSERemarks  string `json:"bse_remarks"`
	SuccessFlag string `json:"success_flag"` // 0 - Success & 1 - failure
	IntRefNo    string `json:"int_ref_no"`
}
