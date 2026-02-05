package handlers

import (
	"company_mgmt_api/middlewares"
	"company_mgmt_api/services"
	"encoding/json"
	"net/http"
)

type VerificationHandler struct {
	// Add necessary fields here, e.g., services, config, etc.
	OTP *services.OTPService
}

func (h *VerificationHandler) Request(w http.ResponseWriter, r *http.Request) {
	// Implementation for sending OTP
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(middlewares.UserIDKey).(string)

	if err := h.OTP.GenerateAndSendOTP(r.Context(), userID, req.Email, "email verification"); err != nil {
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *VerificationHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	// Implementation for verifying OTP
	var req struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	userID := r.Context().Value(middlewares.UserIDKey).(string)

	if err := h.OTP.VerifyOTP(r.Context(), userID, req.Code, "email verification"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
