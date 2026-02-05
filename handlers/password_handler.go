package handlers

import (
	"company_mgmt_api/services"
	"encoding/json"
	"net/http"
)

type PasswordHandler struct {
	OTP *services.OTPService
}

func (h *PasswordHandler) Forgot(w http.ResponseWriter, r *http.Request) {
	// Implementation for sending password reset OTP
	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	//find userID by email
	userID := "user_id"

	if err := h.OTP.GenerateAndSendOTP(r.Context(), userID, req.Email, "password reset"); err != nil {
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *PasswordHandler) Reset(w http.ResponseWriter, r *http.Request) {
	// Implementation for verifying OTP and resetting password
	var req struct {
		Code     string `json:"code"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	// verify OTP + update password (next phase)
	w.WriteHeader(http.StatusNoContent)
}
