package bse

import (
	
	"fmt"
)

// LumpsumOrderRequest represents the request structure for BSE StAR MF Lumpsum Order Entry
type LumpsumOrderRequest struct {
	// Mandatory fields
	TransCode    string  `json:"TransCode"`    // NEW/MOD/CXL
	TransNo      string  `json:"TransNo"`      // Unique reference number (YYYYMMDD000001)
	OrderID      *int64  `json:"OrderId,omitempty"` // BSE order number (blank for new orders)
	UserID       int64   `json:"UserId"`       // User ID as given by BSE
	MemberID     string  `json:"MemberId"`     // Member code as given by BSE
	ClientCode   string  `json:"ClientCode"`   // Client Code
	SchemeCd     string  `json:"SchemeCd"`     // BSE scheme code
	BuySell      string  `json:"BuySell"`      // P/R (Purchase/Redemption)
	BuySellType  string  `json:"BuySellType"`  // FRESH/ADDITIONAL
	DPTxn        string  `json:"DPTxn"`        // C/N/P (CDSL/NSDL/PHYSICAL)
	AllRedeem    string  `json:"AllRedeem"`    // Y/N
	KYCStatus    string  `json:"KYCStatus"`    // Y/N
	EUIN         string  `json:"EUIN"`         // EUIN number
	EUINFlag     string  `json:"EUINFlag"`     // Y/N
	MinRedeem    string  `json:"MinRedeem"`    // Y/N
	DPC          string  `json:"DPC"`          // Y
	IPAdd        string  `json:"IPAdd"`        // IP Address
	Password     string  `json:"Password"`     // Encrypted password
	PassKey      string  `json:"PassKey"`      // Pass Key
	Param1       string  `json:"Param1"`       // Sub Broker ARN
	MobileNo     string  `json:"MobileNo"`     // 10 digit mobile number

	// Either Amount or Qty fields
	Amount *float64 `json:"Amount,omitempty"` // Purchase/Redemption amount
	Qty    *float64 `json:"Qty,omitempty"`    // Redemption quantity

	// Optional fields
	FolioNo   string `json:"FolioNo,omitempty"`   // For physical transactions
	Remarks   string `json:"Remarks,omitempty"`   // Remarks
	RefNo     string `json:"RefNo,omitempty"`     // Internal reference
	SubBrCode string `json:"SubBrCode,omitempty"` // Sub Broker code
	Param2    string `json:"Param2,omitempty"`    // PG Reference No (Purchase only)
	Param3    string `json:"Param3,omitempty"`    // Bank Account No (Redemption only)
	EmailID   string `json:"EmailID,omitempty"`   // Email ID
	MandateID string `json:"MandateID,omitempty"` // For OTM (Purchase only)
	Filler1   string `json:"Filler1,omitempty"`   // Filler 1
	Filler2   string `json:"Filler2,omitempty"`   // Filler 2
	Filler3   string `json:"Filler3,omitempty"`   // Filler 3
	Filler4   string `json:"Filler4,omitempty"`   // Filler 4
	Filler5   string `json:"Filler5,omitempty"`   // Filler 5
	Filler6   string `json:"Filler6,omitempty"`   // Filler 6
}

// LumpsumOrderResponse represents the response structure from BSE StAR MF Lumpsum Order Entry
type LumpsumOrderResponse struct {
	TransactionCode     string `json:"TransactionCode"`     // Transaction Code as given in request
	UniqueRefNo         string `json:"UniqueRefNo"`         // Unique reference number as given in request
	OrderNumber         int64  `json:"OrderNumber"`         // BSE order number
	UserID              int64  `json:"UserId"`              // User ID as given by BSE
	MemberID            string `json:"MemberId"`            // Member code as given by BSE
	ClientCode          string `json:"ClientCode"`          // Client Code
	BSERemarks          string `json:"BSERemarks"`          // BSE Response Return remarks
	SuccessFlag         string `json:"SuccessFlag"`         // Order success flag
	OrderEntryParamResult string `json:"orderEntryParamResult"` // SOAP response result
}

// Constants for lumpsum order entry
const (
	// Transaction codes
	TransCodeNew    = "NEW"
	TransCodeModify = "MOD"
	TransCodeCancel = "CXL"

	// Buy/Sell types
	BuySellPurchase   = "P"
	BuySellRedemption = "R"

	// Buy/Sell sub-types
	BuySellTypeFresh      = "FRESH"
	BuySellTypeAdditional = "ADDITIONAL"

	// DP Transaction types
	DPTxnCDSL     = "C"
	DPTxnNSDL     = "N"
	DPTxnPhysical = "P"

	// Flag values
	FlagYes = "Y"
	FlagNo  = "N"
)

// Validate validates the lumpsum order request
func (req *LumpsumOrderRequest) Validate() error {
	// Validate mandatory fields
	if req.TransCode == "" {
		return fmt.Errorf("transCode is required")
	}
	if req.TransNo == "" {
		return fmt.Errorf("transNo is required")
	}
	if req.UserID == 0 {
		return fmt.Errorf("userID is required")
	}
	if req.MemberID == "" {
		return fmt.Errorf("memberID is required")
	}
	if req.ClientCode == "" {
		return fmt.Errorf("clientCode is required")
	}
	if req.SchemeCd == "" {
		return fmt.Errorf("schemeCd is required")
	}
	if req.BuySell == "" {
		return fmt.Errorf("buySell is required")
	}
	if req.BuySellType == "" {
		return fmt.Errorf("buySellType is required")
	}
	if req.DPTxn == "" {
		return fmt.Errorf("dPTxn is required")
	}
	if req.AllRedeem == "" {
		return fmt.Errorf("allRedeem is required")
	}
	if req.KYCStatus == "" {
		return fmt.Errorf("kYCStatus is required")
	}
	if req.EUIN == "" {
		return fmt.Errorf("eUIN is required")
	}
	if req.EUINFlag == "" {
		return fmt.Errorf("eUINFlag is required")
	}
	if req.MinRedeem == "" {
		return fmt.Errorf("minRedeem is required")
	}
	if req.DPC == "" {
		return fmt.Errorf("dPC is required")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	if req.PassKey == "" {
		return fmt.Errorf("passKey is required")
	}
	if req.Param1 == "" {
		return fmt.Errorf("param1 (Sub Broker ARN) is required")
	}
	if req.MobileNo == "" {
		return fmt.Errorf("mobileNo is required")
	}

	// Validate Amount or Qty (either one required)
	if req.AllRedeem != FlagYes {
		if req.Amount == nil && req.Qty == nil {
			return fmt.Errorf("either amount or qty is required when allRedeem is not Y")
		}
	}

	// Validate transaction code
	if req.TransCode != TransCodeNew && req.TransCode != TransCodeModify && req.TransCode != TransCodeCancel {
		return fmt.Errorf("transCode must be NEW, MOD, or CXL")
	}

	// Validate BuySell
	if req.BuySell != BuySellPurchase && req.BuySell != BuySellRedemption {
		return fmt.Errorf("buySell must be P or R")
	}

	// Validate BuySellType
	if req.BuySellType != BuySellTypeFresh && req.BuySellType != BuySellTypeAdditional {
		return fmt.Errorf("buySellType must be FRESH or ADDITIONAL")
	}

	// Validate DPTxn
	if req.DPTxn != DPTxnCDSL && req.DPTxn != DPTxnNSDL && req.DPTxn != DPTxnPhysical {
		return fmt.Errorf("dPTxn must be C, N, or P")
	}

	// Validate flags
	if req.AllRedeem != FlagYes && req.AllRedeem != FlagNo {
		return fmt.Errorf("allRedeem must be Y or N")
	}
	if req.KYCStatus != FlagYes && req.KYCStatus != FlagNo {
		return fmt.Errorf("kYCStatus must be Y or N")
	}
	if req.EUINFlag != FlagYes && req.EUINFlag != FlagNo {
		return fmt.Errorf("eUINFlag must be Y or N")
	}
	if req.MinRedeem != FlagYes && req.MinRedeem != FlagNo {
		return fmt.Errorf("minRedeem must be Y or N")
	}
	if req.DPC != FlagYes {
		return fmt.Errorf("dPC must be Y")
	}

	return nil
}

// IsSuccessResponse checks if the order was successful
func (resp *LumpsumOrderResponse) IsSuccessResponse() bool {
	return resp.SuccessFlag == FlagYes
}