package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"company_mgmt_api/repositories"
)

type AuditHandler struct {
	Repo *repositories.AuditRepository
}

func NewAuditHandler(r *repositories.AuditRepository) *AuditHandler {
	return &AuditHandler{Repo: r}
}

func (h *AuditHandler) List(w http.ResponseWriter, r *http.Request) {
	companyID := chi.URLParam(r, "companyId")

	logs, err := h.Repo.ListByCompany(r.Context(), companyID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(logs)
}
