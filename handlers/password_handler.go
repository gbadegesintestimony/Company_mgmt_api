package handlers

import (
	"company_mgmt_api/repositories"
	"company_mgmt_api/services"
	"encoding/json"
	"log"
	"net/http"
)

type PasswordHandler struct {
	OTP         *services.OTPService
	UserRepo    *repositories.UserRepository
	CompanyRepo *repositories.CompanyRepository
	AuthSvc     *services.AuthService
}

func (h *PasswordHandler) Forgot(w http.ResponseWriter, r *http.Request) {
	// Implementation for sending password reset OTP
	var req struct {
		Email  string `json:"email"`
		Domain string `json:"domain"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("FORGOT PASSWORD — email: %s, domain: %s", req.Email, req.Domain)

	// find company by domain first

	company, err := h.CompanyRepo.FindByDomain(r.Context(), req.Domain)
	if err != nil {
		log.Printf("FORGOT PASSWORD — company not found for domain %s: %v", req.Domain, err)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	log.Printf("FORGOT PASSWORD — company found: %s", company.ID)
	user, err := h.UserRepo.FindByEmailAndCompany(r.Context(), req.Email, company.ID)
	if err != nil {
		log.Printf("FORGOT PASSWORD — user not found for email %s: %v", req.Email, err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	log.Printf("FORGOT PASSWORD — user found: %s", user.ID)

	if err := h.OTP.GenerateAndSendOTP(r.Context(), user.ID, req.Email, "password_reset"); err != nil {
		log.Printf("FORGOT PASSWORD — OTP send failed: %v", err)
		http.Error(w, "failed to send OTP", http.StatusInternalServerError)
		return
	}
	log.Printf("FORGOT PASSWORD — OTP sent successfully to %s", req.Email)
	w.WriteHeader(http.StatusNoContent)
}

func (h *PasswordHandler) Reset(w http.ResponseWriter, r *http.Request) {
	// Implementation for verifying OTP and resetting password
	var req struct {
		Email    string `json:"email"`
		Domain   string `json:"domain"`
		Code     string `json:"code"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	company, err := h.CompanyRepo.FindByDomain(r.Context(), req.Domain)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Find user by email
	user, err := h.UserRepo.FindByEmailAndCompany(r.Context(), req.Email, company.ID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Verify OTP and reset password
	if err := h.AuthSvc.ResetPassword(r.Context(), user.ID, req.Code, req.Password); err != nil {
		http.Error(w, "invalid or expired OTP", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "password reset successfully"})
}
