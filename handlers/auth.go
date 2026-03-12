package handlers

import (
	"company_mgmt_api/config"
	"company_mgmt_api/models"
	"company_mgmt_api/repositories"
	"company_mgmt_api/services"
	"company_mgmt_api/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type AuthHandler struct {
	Cfg         *config.Config
	Email       *services.EmailService
	Sessions    *repositories.SessionRepository
	UserRepo    *repositories.UserRepository
	CompanyRepo *repositories.CompanyRepository
	OTPService  *services.OTPService
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		// CompanyID string `json:"company"`
		Company  string `json:"company"`
		Domain   string `json:"domain"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Check for existing company by domain
	if existingCompany, err := h.CompanyRepo.FindByDomain(ctx, req.Domain); err == nil {
		// Company exists - check if admin is verified
		existingUserID, err := h.UserRepo.FindIDByEmail(ctx, req.Email)
		if err != nil {
			// Different email trying to register same domain
			http.Error(w, "company already registered", http.StatusConflict)
			return
		}

		isVerified, err := h.UserRepo.IsVerified(ctx, existingUserID)
		if err != nil {
			http.Error(w, "failed to check account status", http.StatusInternalServerError)
			return
		}
		if isVerified {
			http.Error(w, "Company already registered", http.StatusConflict)
			return
		}
		if err := h.UserRepo.DeleteByID(ctx, existingUserID); err != nil {
			http.Error(w, "failed to reset unverified account", http.StatusInternalServerError)
			return
		}
		if err := h.CompanyRepo.DeleteByID(ctx, existingCompany.ID); err != nil {
			http.Error(w, "failed to reset unverified company", http.StatusInternalServerError)
			return
		}
	}

	// Generate IDs
	companyID := uuid.New().String()
	userID := uuid.New().String()

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	// 1️⃣ Create Company
	err = h.CompanyRepo.Create(ctx, companyID, req.Company, req.Domain)
	if err != nil {
		http.Error(w, "failed to create company", http.StatusInternalServerError)
		return
	}

	// 2️⃣ Create Admin User
	err = h.UserRepo.CreateAdmin(
		ctx,
		userID,
		req.Email,
		passwordHash,
		companyID,
	)
	if err != nil {
		http.Error(w, "failed to create admin user", http.StatusInternalServerError)
		return
	}

	if err := h.OTPService.GenerateAndSendOTP(ctx, userID, req.Email, "email_verification"); err != nil {

		log.Printf("OTP SENd FAILED FOR %q: %v", req.Email, err)
		http.Error(w, "failed to send verification email", http.StatusInternalServerError)
		return
	}
	log.Printf("OTP SENT successfully to %q", req.Email)

	// 5️⃣ Return Clean Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	type CreateResponse struct {
		Message   string `json:"message"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
	}

	json.NewEncoder(w).Encode(CreateResponse{
		Message:   "Registration initiated. An OTP has been sent to your email — please verify to activate your account.",
		Email:     req.Email,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	})

}

func (h *AuthHandler) AdminLogin(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Domain   string `json:"domain"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// find company by domain
	company, err := h.CompanyRepo.FindByDomain(ctx, req.Domain)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Fetch user from DB
	user, err := h.UserRepo.FindByEmailAndCompany(ctx, req.Email, company.ID) // companyID can be optional for login
	if err != nil {

		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if user.EmailVerifiedAt == nil {
		http.Error(w, "email not verified, please check inbox", http.StatusForbidden)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if !user.IsActive {
		http.Error(w, "account disabled", http.StatusForbidden)
		return
	}

	//  Verify password
	if err := utils.VerifyPassword(req.Password, user.PasswordHash); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	h.issueTokens(w, r, user, "Login Successful")

}

func (h *AuthHandler) EmployeeLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Domain   string `json:"domain"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	company, err := h.CompanyRepo.FindByDomain(ctx, req.Domain)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	user, err := h.UserRepo.FindByEmailAndCompany(ctx, req.Email, company.ID)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Enforce employee role
	if user.Role != "employee" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if !user.IsActive {
		http.Error(w, "account disabled", http.StatusForbidden)
		return
	}

	if err := utils.VerifyPassword(req.Password, user.PasswordHash); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if user.EmailVerifiedAt == nil {
		if err := h.OTPService.GenerateAndSendOTP(ctx, user.ID, req.Email, "email_verification"); err != nil {
			log.Printf("OTP SEND FAILED for %q: %v", req.Email, err)
			http.Error(w, "failed to send verification email", http.StatusInternalServerError)
			return
		}
		h.issueTokens(w, r, user, "Account not yet verified. An OTP has been sent to your email — please verify your account.")
		return
	}
	h.issueTokens(w, r, user, "Login Successful")
}

func (h *AuthHandler) issueTokens(w http.ResponseWriter, r *http.Request, user *models.User, message string) {
	access, err := utils.GenerateAccessToken(user.ID, user.CompanyID, user.Role, h.Cfg.JWTAccessSecret)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	refresh, err := utils.GenerateRefreshToken(h.Cfg.JWTRefreshSecret)
	if err != nil {
		http.Error(w, "failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	err = h.Sessions.CreateSession(r.Context(), user.ID, user.CompanyID, user.Role, refresh, time.Now().Add(30*24*time.Hour))
	if err != nil {
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}

	type TokenResponse struct {
		Message      string `json:"message"`
		Role         string `json:"role"`
		UserID       string `json:"user_id"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		LoggedInAt   string `json:"logged_in_at"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TokenResponse{
		Message:      message,
		Role:         user.Role,
		UserID:       user.ID,
		AccessToken:  access,
		RefreshToken: refresh,
		LoggedInAt:   time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	session, err := h.Sessions.GetByRefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	access, err := utils.GenerateAccessToken(
		session.UserID,
		session.CompanyID,
		session.Role,
		h.Cfg.JWTAccessSecret,
	)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	// Generate new refresh token
	newRefresh, err := utils.GenerateRefreshToken(h.Cfg.JWTRefreshSecret)
	if err != nil {
		http.Error(w, "failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	// Revoke old session
	if err := h.Sessions.RevokeByRefreshToken(r.Context(), req.RefreshToken); err != nil {
		http.Error(w, "failed to revoke session", http.StatusInternalServerError)
		return
	}

	// Create new session with new refresh token
	if err := h.Sessions.CreateSession(r.Context(), session.UserID, session.CompanyID, session.Role, newRefresh, time.Now().Add(30*24*time.Hour)); err != nil {
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	type RefreshResponse struct {
		Message      string `json:"message"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	json.NewEncoder(w).Encode(RefreshResponse{
		Message:      "Token refreshed successfully",
		AccessToken:  access,
		RefreshToken: newRefresh,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.Sessions.RevokeByRefreshToken(r.Context(), req.RefreshToken); err != nil {
		http.Error(w, "failed to logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
