package handlers

import (
	"company_mgmt_api/config"
	"company_mgmt_api/middlewares"
	"company_mgmt_api/repositories"
	"company_mgmt_api/services"
	"company_mgmt_api/utils"
	"encoding/json"
	"net/http"
	"time"
)

type VerificationHandler struct {
	OTP         *services.OTPService
	UserRepo    *repositories.UserRepository
	CompanyRepo *repositories.CompanyRepository
	Sessions    *repositories.SessionRepository
	Cfg         *config.Config
}

func (h *VerificationHandler) Request(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(middlewares.UserIDKey).(string)

	if err := h.OTP.GenerateAndSendOTP(r.Context(), userID, req.Email, "email_verification"); err != nil {
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *VerificationHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"otp_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	userID, err := h.UserRepo.FindIDByEmail(ctx, req.Email)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Verify OTP
	if err := h.OTP.VerifyOTP(ctx, userID, req.Code, "email_verification"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Activate account
	if err := h.UserRepo.SetEmailVerified(ctx, userID); err != nil {
		http.Error(w, "failed to activate account", http.StatusInternalServerError)
		return
	}

	// Fetch user
	user, err := h.UserRepo.FindByID(ctx, userID)
	if err != nil {
		http.Error(w, "failed to fetch user", http.StatusInternalServerError)
		return
	}

	// Fetch company
	company, err := h.CompanyRepo.FindByID(ctx, user.CompanyID)
	if err != nil {
		http.Error(w, "failed to fetch company", http.StatusInternalServerError)
		return
	}

	// Generate tokens
	access, err := utils.GenerateAccessToken(user.ID, user.CompanyID, user.Role, h.Cfg.JWTAccessSecret)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	refresh, err := utils.GenerateRefreshToken()
	if err != nil {
		http.Error(w, "failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	if err := h.Sessions.CreateSession(ctx, user.ID, user.CompanyID, user.Role, refresh, time.Now().Add(30*24*time.Hour)); err != nil {
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Email verified successfully. Your account is now active.",
		"verified_at": time.Now().UTC().Format(time.RFC3339),
		"role":        user.Role,
		"user_id":     user.ID,
		"company": map[string]string{
			"id":         company.ID,
			"name":       company.Name,
			"email":      user.Email,
			"domain":     company.Domain,
			"created_at": company.CreatedAt.UTC().Format(time.RFC3339),
		},
		"tokens": map[string]string{
			"access_token":  access,
			"refresh_token": refresh,
		},
	})
}
