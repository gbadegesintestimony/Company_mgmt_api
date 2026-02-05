package handlers

import (
	"company_mgmt_api/config"
	"company_mgmt_api/repositories"
	"company_mgmt_api/services"
	"company_mgmt_api/utils"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	Cfg      *config.Config
	Email    *services.EmailService
	Sessions *repositories.SessionRepository
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Registration logic would go here
	var req struct {
		Company  string `json:"company"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	// Hash the password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	_ = passwordHash // Use the variable to avoid unused warning

	// DB insertion logic would go here (omitted for brevity)
	h.Email.SendOTP(r.Context(), req.Email, "123456", "Email")

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Login logic would go here
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	// DB retrieval and password verification logic would go here (omitted for brevity)
	access, err := utils.GenerateAccessToken(
		"user-id", "company-id", "admin", h.Cfg.JWTAccessSecret,
	)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	refresh, err := utils.GenerateRefreshToken(h.Cfg.JWTRefreshSecret)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  access,
		"refresh_token": refresh,
	})
}
