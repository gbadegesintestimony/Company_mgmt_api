package handlers

import (
	"company_mgmt_api/middlewares"
	"company_mgmt_api/models"
	"company_mgmt_api/repositories"
	"company_mgmt_api/services"
	"company_mgmt_api/utils"
	"encoding/json"
	"net/http"
)

type MeHandler struct {
	// Implementation goes here
	Employees *repositories.EmployeeRepository
	Sessions  *repositories.SessionRepository
	Services  *services.EmployeeService
	UserRepo  *repositories.UserRepository
}

func (h *MeHandler) Profile(w http.ResponseWriter, r *http.Request) {
	// Handler logic goes here
	var p models.Profile
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(middlewares.UserIDKey).(string)
	companyID := r.Context().Value(middlewares.CompanyIDKey).(string)

	if err := h.Employees.UpdateProfile(
		r.Context(),
		userID,
		companyID,
		&p,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := h.UserRepo.FindByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to fetch updated profile", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Profile Updated Successfully",
		"profile": map[string]string{
			"user_id":    user.ID,
			"email":      user.Email,
			"first_name": strVal(user.FirstName),
			"last_name":  strVal(user.LastName),
			"role":       user.Role,
		},
	})
}

func (h *MeHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(string)
	// companyID := r.Context().Value(middlewares.CompanyIDKey).(string)

	user, err := h.UserRepo.FindByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	type ProfileResponse struct {
		UserID     string `json:"user_id"`
		Email      string `json:"email"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Phone      string `json:"phone"`
		JobTitle   string `json:"job_title"`
		Department string `json:"department"`
		Role       string `json:"role"`
		CompanyID  string `json:"company_id"`
		IsActive   bool   `json:"is_active"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ProfileResponse{
		UserID:     user.ID,
		Email:      user.Email,
		FirstName:  strVal(user.FirstName),
		LastName:   strVal(user.LastName),
		Phone:      strVal(user.Phone),
		JobTitle:   strVal(user.JobTitle),
		Department: strVal(user.Department),
		Role:       user.Role,
		CompanyID:  user.CompanyID,
		IsActive:   user.IsActive,
	})
}

func (h *MeHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Handler logic goes here
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(middlewares.UserIDKey).(string)
	companyID := r.Context().Value(middlewares.CompanyIDKey).(string)

	// Get user to verify old password
	employee, err := h.Employees.FindByIDAndCompany(r.Context(), userID, companyID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Verify old password
	if err := utils.VerifyPassword(req.OldPassword, employee.PasswordHash); err != nil {
		http.Error(w, "incorrect current password", http.StatusUnauthorized)
		return
	}

	newHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	if err := h.Employees.UpdatePassword(r.Context(), userID, companyID, newHash); err != nil {
		http.Error(w, "failed to update password", http.StatusInternalServerError)
		return
	}
	if err := h.Sessions.RevokeAll(r.Context(), userID); err != nil {
		http.Error(w, "failed to revoke sessions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "password changed successfully. Please log in again"})
}

// Helper to safely dereference *string
func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
