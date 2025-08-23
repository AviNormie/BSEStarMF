package types

// LumpsumOrderRequest represents a lumpsum mutual fund order request
type LumpsumOrderRequest struct {
	// Basic Transaction Details
	TransCode    string `json:"trans_code" validate:"required,oneof=NEW" xml:"TransCode"`           // Transaction type - NEW only
	TransNo      string `json:"trans_no" validate:"required,max=19" xml:"TransNo"`                // Unique order reference
	OrderId      string `json:"order_id,omitempty" xml:"OrderId"`                                 // BSE order number (blank for new)
	UserID       string `json:"user_id" validate:"required,max=5" xml:"UserID"`                   // BSE user code
	MemberId     string `json:"member_id" validate:"required,max=20" xml:"MemberId"`              // BSE member code
	ClientCode   string `json:"client_code" validate:"required,max=20" xml:"ClientCode"`          // Client identifier

	// Investment Details
	SchemeCd     string `json:"scheme_cd" validate:"required,max=20" xml:"SchemeCd"`              // Scheme Code as per BSE
	BuySell      string `json:"buy_sell" validate:"required,oneof=P R" xml:"BuySell"`            // P (Purchase) / R (Redemption)
	BuySellType  string `json:"buy_sell_type" validate:"required,oneof=FRESH ADDITIONAL" xml:"BuySellType"` // FRESH (new) or ADDITIONAL
	DPTxn        string `json:"dp_txn" validate:"required,oneof=CDSL NSDL PHYSICAL P" xml:"DPTxn"` // CDSL/NSDL/PHYSICAL

	// Order Amount/Quantity
	OrderVal     string `json:"order_val,omitempty" validate:"max=14" xml:"OrderVal"`             // Amount for purchase/redeem
	Qty          string `json:"qty,omitempty" validate:"max=8" xml:"Qty"`                        // Quantity of units
	AllRedeem    string `json:"all_redeem" validate:"required,oneof=Y N" xml:"AllRedeem"`         // Redeem all units? (Y/N)

	// Additional Details
	FolioNo      string `json:"folio_no,omitempty" validate:"max=20" xml:"FolioNo"`               // Folio Number
	Remarks      string `json:"remarks,omitempty" validate:"max=255" xml:"Remarks"`              // Any comments
	RefNo        string `json:"ref_no,omitempty" validate:"max=20" xml:"RefNo"`                   // Internal reference

	// Compliance & Regulatory
	KYCStatus    string `json:"kyc_status" validate:"required,oneof=Y N" xml:"KYCStatus"`         // KYC status of client (Y/N)
	SubBrCode    string `json:"sub_br_code,omitempty" validate:"max=15" xml:"SubBrCode"`          // Sub broker code
	EUIN         string `json:"euin,omitempty" validate:"max=20" xml:"EUIN"`                      // EUIN number
	EUINVal      string `json:"euin_val" validate:"required,oneof=Y N" xml:"EUINVal"`            // EUIN declaration (Y/N)
	MinRedeem    string `json:"min_redeem" validate:"required,oneof=Y N" xml:"MinRedeem"`         // Min redemption flag (Y/N)
	DPC          string `json:"dpc" validate:"required,oneof=Y N" xml:"DPC"`                      // DPC Purchase flag (Y/N)

	// Authentication
	Password     string `json:"password" validate:"required,max=250" xml:"Password"`             // Encrypted password
	PassKey      string `json:"pass_key" validate:"required,max=10" xml:"PassKey"`               // Pass key

	// Optional Parameters
	IPAdd        string `json:"ip_add,omitempty" validate:"max=15" xml:"IPAdd"`                   // IP Address
	Parma1       string `json:"parma1,omitempty" validate:"max=20" xml:"Parma1"`                  // Sub Broker ARN (Filler 1)
	Param2       string `json:"param2,omitempty" validate:"max=20" xml:"Param2"`                  // PG Reference No
	Param3       string `json:"param3,omitempty" validate:"max=20" xml:"Param3"`                  // Bank Account No (Redemption only)
	MobileNo     string `json:"mobile_no,omitempty" validate:"max=10" xml:"MobileNo"`             // 10-digit Indian mobile
	EmailID      string `json:"email_id,omitempty" validate:"email,max=50" xml:"EmailID"`         // Client email
	MandateID    string `json:"mandate_id,omitempty" validate:"max=20" xml:"MandateID"`           // For purchase via OTM

	// Filler Fields
	Filler1      string `json:"filler1,omitempty" validate:"max=30" xml:"Filler1"`               // Extra field 1
	Filler2      string `json:"filler2,omitempty" validate:"max=30" xml:"Filler2"`               // Extra field 2
	Filler3      string `json:"filler3,omitempty" validate:"max=30" xml:"Filler3"`               // Extra field 3
	Filler4      string `json:"filler4,omitempty" validate:"max=30" xml:"Filler4"`               // Extra field 4
	Filler5      string `json:"filler5,omitempty" validate:"max=30" xml:"Filler5"`               // Extra field 5
	Filler6      string `json:"filler6,omitempty" validate:"max=30" xml:"Filler6"`               // Extra field 6
}

// LumpsumOrderResponse represents the response from lumpsum order entry
type LumpsumOrderResponse struct {
	Success      bool   `json:"success"`                                                          // Overall success status
	TransCode    string `json:"trans_code" xml:"TransCode"`                                      // Transaction code (from request)
	TransNo      string `json:"trans_no" xml:"TransNo"`                                        // Transaction number (from request)
	OrderId      string `json:"order_id" xml:"OrderId"`                                        // Generated BSE Order number
	UserID       string `json:"user_id" xml:"UserID"`                                          // User ID as given
	MemberId     string `json:"member_id" xml:"MemberId"`                                      // Member code as given
	ClientCode   string `json:"client_code" xml:"ClientCode"`                                  // Client code
	Remarks      string `json:"remarks,omitempty" xml:"Remarks"`                              // BSE remarks
	SuccessFlag  string `json:"success_flag" xml:"SuccessFlag"`                                // Order success (Y/N)
	Message      string `json:"message,omitempty"`                                             // Error or success message
}