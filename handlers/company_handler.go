package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"company_mgmt_api/middlewares"
	"company_mgmt_api/models"
	"company_mgmt_api/services"
)

type CompanyHandler struct {
	Service *services.CompanyService
}

func NewCompanyHandler(s *services.CompanyService) *CompanyHandler {
	return &CompanyHandler{Service: s}
}

func (h *CompanyHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "companyId")

	company, err := h.Service.Get(r.Context(), id)
	if err != nil {
		http.Error(w, "company not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}

func (h *CompanyHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "companyId")
	userID := r.Context().Value(middlewares.UserIDKey).(string)

	var payload struct {
		Name   string `json:"name"`
		Domain string `json:"domain"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	company := &models.Company{
		ID:     id,
		Name:   payload.Name,
		Domain: payload.Domain,
		Status: payload.Status,
	}

	err := h.Service.Update(r.Context(), company, userID)
	if err != nil {
		http.Error(w, "failed to update company", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "company updated successfully"})
}
