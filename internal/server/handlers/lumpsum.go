package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"sapphirebroking.com/sapphire_mf/internal/server/services"
	"sapphirebroking.com/sapphire_mf/internal/server/types"
)

func LumpsumHandler(lumpsumService *services.SOAPClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req types.LumpsumOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := validateLumpsumRequest(&req); err != nil {
			http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
			return
		}

		response, err := lumpsumService.LumpsumOrderEntry(r.Context(), &req)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// Basic validation without external validator library
func validateLumpsumRequest(req *types.LumpsumOrderRequest) error {
	if strings.TrimSpace(req.TransCode) == "" {
		return fmt.Errorf("trans_code is required")
	}
	if req.TransCode != "NEW" {
		return fmt.Errorf("trans_code must be NEW")
	}
	if strings.TrimSpace(req.TransNo) == "" {
		return fmt.Errorf("trans_no is required")
	}
	if len(req.TransNo) > 19 {
		return fmt.Errorf("trans_no must be max 19 characters")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return fmt.Errorf("user_id is required")
	}
	if len(req.UserID) > 5 {
		return fmt.Errorf("user_id must be max 5 characters")
	}
	if strings.TrimSpace(req.MemberId) == "" {
		return fmt.Errorf("member_id is required")
	}
	if strings.TrimSpace(req.ClientCode) == "" {
		return fmt.Errorf("client_code is required")
	}
	if strings.TrimSpace(req.SchemeCd) == "" {
		return fmt.Errorf("scheme_cd is required")
	}
	if req.BuySell != "P" && req.BuySell != "R" {
		return fmt.Errorf("buy_sell must be P (Purchase) or R (Redemption)")
	}
	if req.BuySellType != "FRESH" && req.BuySellType != "ADDITIONAL" {
		return fmt.Errorf("buy_sell_type must be FRESH or ADDITIONAL")
	}
	if strings.TrimSpace(req.Password) == "" {
		return fmt.Errorf("password is required")
	}
	if strings.TrimSpace(req.PassKey) == "" {
		return fmt.Errorf("pass_key is required")
	}
	return nil
}