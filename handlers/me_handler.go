package handlers

import (
	"company_mgmt_api/middlewares"
	"company_mgmt_api/repositories"
	"encoding/json"
	"net/http"
)

type MeHandler struct {
	// Implementation goes here
	Employees *repositories.EmployeeRepository
	Sessions  *repositories.SessionRepository
}

func (h *MeHandler) Profile(w http.ResponseWriter, r *http.Request) {
	// Handler logic goes here
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	// Assume we have userID from context/session
	userID := r.Context().Value(middlewares.UserIDKey).(string)
	companyID := r.Context().Value(middlewares.CompanyIDKey).(string)

	if err := h.Employees.UpdateProfile(
		r.Context(),
		userID,
		companyID,
		req.FirstName,
		req.LastName,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *MeHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Handler logic goes here
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	// hash, _ := utils.HashPassword(req.NewPassword)

	userID := r.Context().Value(middlewares.UserIDKey).(string)

	if err := h.Sessions.RevokeAll(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)

}
