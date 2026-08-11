package handlers

import (
	"company_mgmt_api/middlewares"
	"company_mgmt_api/services"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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
		Role      string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	// Default role to employee if not provided
	if req.Role == "" {
		req.Role = "employee"
	}

	// Assume we have companyID and actorID from context/session
	// companyID := r.Context().Value("company_id").(string)
	companyID := chi.URLParam(r, "companyId")
	actorID := r.Context().Value(middlewares.UserIDKey).(string)

	result, err := h.Service.CreateEmployees(
		r.Context(),
		companyID,
		req.Password,
		req.Email,
		req.FirstName,
		req.LastName,
		req.Role,
		actorID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Employee created successfully",
		"employee": map[string]string{
			"user_id":            result.ID,
			"email":              result.Email,
			"first_name":         result.FirstName,
			"last_name":          result.LastName,
			"role":               result.Role,
			"temporary_password": req.Password,
			"created_at":         result.CreatedAt.UTC().Format(time.RFC3339),
		},
	})
}

func (h *AdminEmployeeHandler) List(w http.ResponseWriter, r *http.Request) {
	// Handler logic goes here
	// Assume we have companyID from context/session
	companyID := r.Context().Value(middlewares.CompanyIDKey).(string)

	limit := 10
	offset := 0
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"data":        employees,
		"total_count": total,
	})
}

func (h *AdminEmployeeHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	// Handler logic goes here
	userID := chi.URLParam(r, "userId")
	// Assume we have companyID from context/session
	// companyID := r.Context().Value("company_id").(string)
	companyID := chi.URLParam(r, "companyId")

	if err := h.Service.Employees.SetActive(
		r.Context(),
		userID,
		companyID,
		false,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AdminEmployeeHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	companyID := chi.URLParam(r, "companyId")

	employee, err := h.Service.Get(r.Context(), userID, companyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(employee)
}

func (h *AdminEmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	companyID := chi.URLParam(r, "companyId")

	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Role      string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	err := h.Service.Update(r.Context(), userID, companyID, req.FirstName, req.LastName, req.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AdminEmployeeHandler) Reactivate(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	companyID := chi.URLParam(r, "companyId")

	if err := h.Service.Reactivate(r.Context(), userID, companyID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AdminEmployeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")
	companyID := chi.URLParam(r, "companyId")

	if err := h.Service.Delete(r.Context(), userID, companyID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
