package handlers

import (
	"company_mgmt_api/middlewares"
	"company_mgmt_api/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type AdminEmployeeHandler struct {
	// Implementation goes here
	Service *services.EmployeeService
}

func (h *AdminEmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Handler logic goes here
	var req struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	// Assume we have companyID and actorID from context/session
	companyID := r.Context().Value("company_id").(string)
	actorID := r.Context().Value("user_id").(string)

	if err := h.Service.CreateEmployees(
		r.Context(),
		companyID,
		req.Password,
		req.Email,
		req.FirstName,
		req.LastName,
		actorID,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AdminEmployeeHandler) List(w http.ResponseWriter, r *http.Request) {
	// Handler logic goes here
	// Assume we have companyID from context/session
	companyID := r.Context().Value(middlewares.CompanyIDKey).(string)

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	search := r.URL.Query().Get("search")

	employees, total, err := h.Service.Employees.ListWithCount(
		r.Context(),
		companyID,
		search,
		nil,
		limit,
		offset,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"data":        employees,
		"total_count": total,
	})
}

func (h *AdminEmployeeHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	// Handler logic goes here
	id := chi.URLParam(r, "id")
	// Assume we have companyID from context/session
	companyID := r.Context().Value("company_id").(string)

	if err := h.Service.Employees.SetActive(
		r.Context(),
		id,
		companyID,
		false,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
